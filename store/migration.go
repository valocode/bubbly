package store

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/valocode/bubbly/config"

	"github.com/valocode/bubbly/env"

	"github.com/valocode/bubbly/api/core"
	"github.com/zclconf/go-cty/cty"
)

type migration []statement

// generateMigration creates a list of sql statements to be executed based on a changelog
func psqlGenerateMigration(bCtx *env.BubblyContext, ch Changelog) (migration, error) {
	var m migration
	for _, change := range ch {
		switch change.Action {
		case remove:
			switch change.TableInfo.ElementType {
			case tableType:
				err := deleteTableStatement(&m, change.TableInfo.TableName)
				if err != nil {
					return nil, err
				}
				break
			case fieldType:
				err := dropColumnStatement(&m, change.TableInfo.TableName, change.TableInfo.ElementName)
				if err != nil {
					return nil, err
				}
				break
			case joinType:
				err := dropConstraintStatement(&m, change.TableInfo.TableName, change.TableInfo.ElementName)
				if err != nil {
					return nil, err
				}
				break
			case uniqueType:
				err := dropConstraintStatement(&m, change.TableInfo.TableName, change.TableInfo.ElementName)
				if err != nil {
					return nil, err
				}
				break
			default:
				break
			}
			break
		case update:
			switch change.TableInfo.ElementType {
			case fieldType:
				// update column type
				err := alterColumnStatement(bCtx, &m, change.TableInfo, change.To)
				if err != nil {
					return nil, err
				}
				break
			case joinType:
				// changing if the join is unique or not
				err := alterJoinStatement(&m, change.TableInfo, change.To)
				if err != nil {
					return nil, err
				}
				break
			default:
				break
			}
			break
		case create:
			switch change.TableInfo.ElementType {
			case tableType:
				err := createTableStatement(&m, change.TableInfo, change.To)
				if err != nil {
					return nil, err
				}
				break
			case fieldType:
				err := createFieldStatement(&m, change.TableInfo, change.To)
				if err != nil {
					return nil, err
				}
				break
			case joinType:
				err := createJoinStatement(&m, change.TableInfo, change.To)
				if err != nil {
					return nil, err
				}
				break
			case uniqueType:
				err := createUniqueConstraint(&m, change.TableInfo)
				if err != nil {
					return nil, err
				}
				break
			default:
				break
			}
			break
		default:
			break
		}
	}

	return m, nil
}

// drop a table if it exists
func deleteTableStatement(m *migration, table string) error {
	s := statement{
		sql:       "DROP TABLE IF EXISTS %s;",
		arguments: []string{table},
	}
	*m = append(*m, s)
	return nil
}

type statement struct {
	sql       string
	arguments []string
}

// sanitize will sanitize sql input to protect against injections
func sanitize(input string) string {
	return pgx.Identifier.Sanitize([]string{input})
}

// formatQuery will take a query string and add the sanitized arguments
func formatQuery(query string, args []string) string {
	s := make([]interface{}, len(args))
	for i, v := range args {
		arg := sanitize(v)
		// this will remove the extra "" that gets added around a string from sanitize()
		s[i] = strings.Trim(arg, "\"")
	}
	return fmt.Sprintf(query, s...)
}

func psqlMigrate(conn *pgxpool.Pool, migrationList migration) error {
	tx, err := conn.Begin(context.Background())
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(context.Background())

	for _, m := range migrationList {
		_, err = tx.Exec(context.Background(), formatQuery(m.sql, m.arguments))
		if err != nil {
			return fmt.Errorf("failed to execute SQL: %w, QUERY: %s", err, m.sql)
		}
	}

	return tx.Commit(context.Background())
}

// drop a column if it exists
func dropColumnStatement(m *migration, table string, column string) error {
	s := statement{
		sql:       "ALTER TABLE IF EXISTS %s DROP COLUMN IF EXISTS %s;",
		arguments: []string{table, column},
	}
	*m = append(*m, s)
	return nil
}

// drop a constraint if it exists
func dropConstraintStatement(m *migration, table string, constraint string) error {
	s := statement{
		sql:       "ALTER TABLE IF EXISTS %s DROP CONSTRAINT IF EXISTS %s;",
		arguments: []string{table, constraint},
	}
	*m = append(*m, s)
	return nil
}

// alters a column if it exists. This will attempt to convert existing values with the
// USING columnName::newType
// If the conversion cannot be completed, it will return an error
func alterColumnStatement(bCtx *env.BubblyContext, m *migration, info tableInfo, columnType interface{}) error {
	t, ok := (columnType).(cty.Type)
	if !ok {
		return fmt.Errorf("not able to parse value to cty.Type: %s", columnType)
	}
	sqlType, err := psqlType(t)
	if err != nil {
		return err
	}
	// this query checks to see if a particular column exists and then edits it in one go
	switch bCtx.StoreConfig.Provider {
	case config.PostgresStore:
		s := statement{
			sql: "DO $$ " +
				"BEGIN " +
				"IF EXISTS(SELECT * FROM information_schema.columns WHERE table_name='%s' and column_name='%s') " +
				"THEN " +
				"ALTER TABLE \"public\".\"%s\" ALTER COLUMN \"%s\" TYPE %s USING \"%s\"::%s;" +
				"END IF;" +
				"END $$;",
			arguments: []string{info.TableName, info.ElementName, info.TableName, info.ElementName, sqlType, info.ElementName, sqlType},
		}

		*m = append(*m, s)
		break
	case config.CockroachDBStore:
		// FIXME
		// https://github.com/valocode/bubbly/issues/201
		// cockroach doesn't support altering column types in a transaction, so this horrible workaround
		// has to be used instead. This will completely remove data from the original column

		dropStatement := statement{
			sql:       "ALTER TABLE IF EXISTS %s DROP COLUMN IF EXISTS %s;",
			arguments: []string{info.TableName, info.ElementName},
		}
		s := statement{
			sql:       "ALTER TABLE IF EXISTS %s ADD COLUMN IF NOT EXISTS %s %s;",
			arguments: []string{info.TableName, info.TableName, sqlType},
		}
		*m = append(*m, dropStatement, s)
		break
	default:
		return fmt.Errorf("provider type %s not supported", bCtx.StoreConfig.Provider)
	}

	return nil
}

// alter a join constraint if it exists. This will first drop an existing constraint and then add it again
// with the updated values. This is because there is not a simple way to alter a constraint in SQL
func alterJoinStatement(m *migration, info tableInfo, uniqueInterface interface{}) error {
	unique, ok := (uniqueInterface).(bool)
	if !ok {
		return fmt.Errorf("uniqueInterface not assignable to bool: %s", uniqueInterface)
	}
	resetSQL := statement{
		sql:       "ALTER TABLE IF EXISTS %s DROP CONSTRAINT IF EXISTS %s_%s_unique;",
		arguments: []string{info.TableName, info.TableName, info.ElementName},
	}
	*m = append(*m, resetSQL)
	if unique {
		s := statement{
			sql:       "ALTER TABLE IF EXISTS %s ADD CONSTRAINT %s_%s_unique UNIQUE (%s);",
			arguments: []string{info.TableName, info.TableName, info.ElementName, info.ElementName},
		}
		*m = append(*m, s)
	}
	return nil
}

func createTableStatement(m *migration, info tableInfo, tableInterface interface{}) error {
	table, ok := (tableInterface).(core.Table)
	if !ok {
		return fmt.Errorf("tableInterface not assignable to core.Table: %s", tableInterface)
	}
	var columnSQL string
	for _, column := range table.Fields {
		columnType, err := psqlType(column.Type)
		if err != nil {
			return err
		}
		columnSQL += fmt.Sprintf(", %s %s", column.Name, columnType)
		if column.Unique {
			columnSQL += " UNIQUE"
		}
	}
	// For now, we are not using a FK constraint, so this is just a "soft" reference to a parent table
	for _, join := range table.Joins {
		columnSQL += fmt.Sprintf(", %s SERIAL", join.Table)
		if join.Single {
			columnSQL += " UNIQUE"
		}
	}

	var uniqueConstraint string
	if table.Unique {
		uniqueConstraint = " UNIQUE"
	} else {
		uniqueConstraint = ""
	}

	s := statement{
		sql:       "CREATE TABLE IF NOT EXISTS %s (name text%s%s);",
		arguments: []string{info.TableName, uniqueConstraint, columnSQL},
	}
	*m = append(*m, s)

	for _, subTable := range table.Tables {
		err := createTableStatement(m, info, subTable)
		if err != nil {
			return err
		}
	}
	return nil
}

// this will create a column in a table, and then if specified add the UNIQUE constraint
func createFieldStatement(m *migration, info tableInfo, fieldInterface interface{}) error {
	field, ok := (fieldInterface).(core.TableField)
	if !ok {
		return fmt.Errorf("fieldInterface not assignable to core.TableField: %s", fieldInterface)
	}
	tableType, err := psqlType(field.Type)
	if err != nil {
		return err
	}

	s := statement{
		sql:       "ALTER TABLE IF EXISTS %s ADD COLUMN IF NOT EXISTS %s %s;",
		arguments: []string{info.TableName, info.ElementName, tableType},
	}
	*m = append(*m, s)

	if field.Unique {
		uniqueStatement := statement{
			sql:       "ALTER TABLE IF EXISTS %s ADD CONSTRAINT %s_%s_key UNIQUE (%s);",
			arguments: []string{info.TableName, info.TableName, info.ElementName, info.ElementName},
		}
		*m = append(*m, uniqueStatement)
	}

	return nil
}

// joins are currently not managed by FK constraints and just share the same name
// because of this, adding a join, is simply adding a column with the convention "parentTableName_id"
// with the type of SERIAL
func createJoinStatement(m *migration, info tableInfo, joinInterface interface{}) error {
	join, ok := (joinInterface).(core.TableJoin)
	if !ok {
		return fmt.Errorf("uniqueInterface not assignable to bool: %s", joinInterface)
	}
	s := statement{
		sql:       "ALTER TABLE IF EXISTS %s ADD COLUMN IF NOT EXISTS %s_id SERIAL;",
		arguments: []string{info.TableName, info.ElementName},
	}
	if join.Single {
		uniqueSQL := statement{
			sql:       "ALTER TABLE IF EXISTS %s ADD CONSTRAINT %s_%s_unique UNIQUE (%s_id);",
			arguments: []string{info.TableName, info.TableName, info.ElementName, info.ElementName},
		}
		*m = append(*m, s, uniqueSQL)
	} else {
		*m = append(*m, s)
	}
	return nil
}

func createUniqueConstraint(m *migration, info tableInfo) error {
	s := statement{
		sql:       "ALTER TABLE IF EXISTS %s ADD CONSTRAINT IF NOT EXISTS %s_%s_key UNIQUE(id);",
		arguments: []string{info.TableName, info.TableName, info.ElementName},
	}
	*m = append(*m, s)
	return nil
}
