package store

import (
	"context"
	"fmt"
	"reflect"

	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/valocode/bubbly/config"

	"github.com/valocode/bubbly/api/core"
	"github.com/zclconf/go-cty/cty"
)

type migration []string

// generateMigration creates a list of sql statements to be executed based on a schemaUpdates
func psqlGenerateMigration(provider config.StoreProviderType, tenant string, ch schemaUpdates) (migration, error) {
	var m migration
	for _, change := range ch {
		tableName := change.TableInfo.TableName
		switch change.Action {
		case remove:
			switch change.TableInfo.ElementType {
			case tableType:
				m = append(m, "DROP TABLE IF EXISTS "+psqlAbsTableName(tenant, tableName))
			case fieldType:
				m = append(m, "ALTER TABLE IF EXISTS "+psqlAbsTableName(tenant, tableName)+" DROP COLUMN IF EXISTS "+change.TableInfo.ElementName)
			case joinType, uniqueType:
				stmt, err := removeJoinStatement(tenant, change.TableInfo, change.From)
				if err != nil {
					return nil, err
				}
				m = append(m, stmt)
			default:
				return nil, fmt.Errorf("unsupported element type: %s", change.TableInfo.ElementType)
			}
		case update:
			switch change.TableInfo.ElementType {
			case fieldType:
				stmts, err := alterColumnStatement(provider, tenant, change.TableInfo, change.To)
				if err != nil {
					return nil, err
				}
				m = append(m, stmts...)
			case joinType:
				// If a join has been changed, it means it's Single attribute
				// has changed, and this does not affect anything in Postgres
				// but only affects the GraphQL resolvers
			case uniqueType:
				// TODO: we should update the column to add a unique constraint...
			default:
				return nil, fmt.Errorf("unsupported element type: %s", change.TableInfo.ElementType)
			}
		case create:
			switch change.TableInfo.ElementType {
			case tableType:
				table, ok := (change.To).(core.Table)
				if !ok {
					return nil, fmt.Errorf("tableInterface not assignable to core.Table: %s", change.TableInfo.TableName)
				}
				stmt, err := psqlTableCreate(tenant, table)
				if err != nil {
					return nil, fmt.Errorf("failed to create SQL statement to create table %s: %w", table.Name, err)
				}
				m = append(m, stmt)
			case fieldType:
				stmts, err := createFieldStatement(tenant, change.TableInfo, change.To)
				if err != nil {
					return nil, err
				}
				m = append(m, stmts...)
			case joinType:
				stmt, err := createJoinStatement(tenant, change.TableInfo, change.To)
				if err != nil {
					return nil, err
				}
				m = append(m, stmt)
			case uniqueType:
				m = append(m, "ALTER TABLE IF EXISTS "+psqlAbsTableName(tenant, tableName)+
					" ADD CONSTRAINT IF NOT EXISTS "+change.TableInfo.TableName+"_"+change.TableInfo.ElementName+"_key"+
					" UNIQUE("+change.TableInfo.ElementName+");")
			}
		}
	}

	return m, nil
}

func psqlMigrate(conn *pgxpool.Pool, tenant string, schema *bubblySchema, migr migration) error {
	tx, err := conn.Begin(context.Background())
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(context.Background())

	for _, m := range migr {
		_, err = tx.Exec(context.Background(), m)
		if err != nil {
			return fmt.Errorf("failed to execute SQL: %w", err)
		}
	}

	// Store the new schema by converting it to core.Data and preparing a
	// saveContext including the schema itself
	d, err := schema.Data()
	if err != nil {
		return fmt.Errorf("failed to create data block from schema: %w", err)
	}
	node := newDataNode(&d)
	// Save the data block node to the schemaTable
	if err := psqlSaveNode(tx, tenant, node, schemaTable); err != nil {
		return fmt.Errorf("failed to save schema data block: %w", err)
	}

	return tx.Commit(context.Background())
}

// alters a column if it exists. This will attempt to convert existing values with the
// USING columnName::newType
// If the conversion cannot be completed, it will return an error
func alterColumnStatement(provider config.StoreProviderType, tenant string, info tableInfo, columnType interface{}) ([]string, error) {
	t, ok := (columnType).(cty.Type)
	if !ok {
		return nil, fmt.Errorf("cannot assign type to cty.Type: %s", reflect.TypeOf(columnType).String())
	}
	sqlType, err := psqlType(t)
	if err != nil {
		return nil, fmt.Errorf("failed to get postgres type for cty type: %w", err)
	}
	// this query checks to see if a particular column exists and then edits it in one go
	switch provider {
	case config.PostgresStore:
		return []string{"DO $$ " +
			"BEGIN " +
			"IF EXISTS(" +
			"SELECT * FROM information_schema.columns " +
			"WHERE table_schema='" + psqlSchemaName(tenant) + "' AND table_name='" + info.TableName + "' AND column_name='" + info.ElementName +
			"') " +
			"THEN " +
			"ALTER TABLE " + psqlAbsTableName(tenant, info.TableName) + " ALTER COLUMN \"" + info.ElementName + "\" TYPE " + sqlType + " USING \"" + info.ElementName + "\"::" + sqlType + ";" +
			"END IF;" +
			"END $$;"}, nil

	case config.CockroachDBStore:
		// FIXME
		// https://github.com/valocode/bubbly/issues/201
		// cockroach doesn't support altering column types in a transaction, so this horrible workaround
		// has to be used instead. This will completely remove data from the original column

		return []string{
			"ALTER TABLE IF EXISTS " + psqlAbsTableName(tenant, info.TableName) + " DROP COLUMN IF EXISTS " + info.ElementName,
			"ALTER TABLE IF EXISTS " + psqlAbsTableName(tenant, info.TableName) + " ADD COLUMN IF NOT EXISTS " + info.ElementName + " " + sqlType,
		}, nil

	default:
		return nil, fmt.Errorf("unsupported provider: %s", provider)
	}
}

// this will create a column in a table, and then if specified add the UNIQUE constraint
func createFieldStatement(tenant string, info tableInfo, fieldInterface interface{}) ([]string, error) {
	field, ok := (fieldInterface).(core.TableField)
	if !ok {
		return nil, fmt.Errorf("cannot assign type to core.TableField: %s", reflect.TypeOf(fieldInterface).String())
	}
	fieldType, err := psqlType(field.Type)
	if err != nil {
		return nil, fmt.Errorf("could not get postgres type for field %s: %w", field.Name, err)
	}

	var statements = make([]string, 0, 1)
	statements = append(statements, "ALTER TABLE IF EXISTS "+psqlAbsTableName(tenant, info.TableName)+" ADD COLUMN IF NOT EXISTS "+info.ElementName+" "+fieldType)
	if field.Unique {
		statements = append(statements,
			"ALTER TABLE IF EXISTS "+psqlAbsTableName(tenant, info.TableName)+" ADD CONSTRAINT "+info.TableName+"_"+info.ElementName+"_key UNIQUE ("+info.ElementName+");",
		)
	}

	return statements, nil
}

// joins are currently not managed by FK constraints and just share the same name
// because of this, adding a join, is simply adding a column with the convention "parentTableName_id"
// with the type of SERIAL
func createJoinStatement(tenant string, info tableInfo, joinInterface interface{}) (string, error) {
	join, ok := (joinInterface).(core.TableJoin)
	if !ok {
		return "", fmt.Errorf("cannot assign type to core.TableJoin: %s", reflect.TypeOf(joinInterface).String())
	}
	return "ALTER TABLE IF EXISTS " + psqlAbsTableName(tenant, info.TableName) + " ADD COLUMN IF NOT EXISTS " + join.Table + tableJoinSuffix + " SERIAL;", nil
}

// removeJoinStatement removes a the column which acts as a join for a table
func removeJoinStatement(tenant string, info tableInfo, joinInterface interface{}) (string, error) {
	join, ok := (joinInterface).(core.TableJoin)
	if !ok {
		return "", fmt.Errorf("cannot assign type to core.TableJoin: %s", reflect.TypeOf(joinInterface).String())
	}
	return "ALTER TABLE IF EXISTS " + psqlAbsTableName(tenant, info.TableName) + " DROP COLUMN IF EXISTS " + join.Table + tableJoinSuffix, nil
}
