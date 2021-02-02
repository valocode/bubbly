package store

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"

	sq "github.com/Masterminds/squirrel"
	"github.com/graphql-go/graphql"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/log/zerologadapter"
	"github.com/verifa/bubbly/env"
)

func psqlNewConn(bCtx *env.BubblyContext, connStr string) (*pgx.Conn, error) {
	config, err := pgx.ParseConfig(connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse db config: %w", err)
	}

	config.Logger = zerologadapter.NewLogger(*bCtx.Logger)
	// TODO: set LogLevel properly... not just debug
	config.LogLevel = pgx.LogLevelDebug

	// config.ConnConfig.LogLevel
	conn, err := pgx.ConnectConfig(context.Background(), config)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to db: %w", err)
	}

	return conn, nil
}

func psqlSaveNode(tx pgx.Tx, node *dataNode, schema *bubblySchema) error {
	table, ok := schema.Tables[node.Data.TableName]
	if !ok {
		return fmt.Errorf("data block refers to non-existing table: %s", node.Data.TableName)
	}

	sql, err := psqlDataNodeUpsert(node, table)
	if err != nil {
		return err
	}
	// Generate the sql string and args
	sqlStr, sqlArgs, err := sql.ToSql()
	if err != nil {
		return fmt.Errorf("failed to create SQL statement: %w", err)
	}
	row := tx.QueryRow(context.Background(), sqlStr, sqlArgs...)

	retValues, err := psqlRowValues(row, table.Name, node.orderedRefFields())
	if err != nil {
		return fmt.Errorf("failed to insert data block: %s: %w", table.Name, err)
	}

	// Asign the returned values so that if the child nodes need to resolve
	// their data references they have values to do so
	node.Return = retValues
	return nil
}

func psqlApplySchema(tx pgx.Tx, schema *bubblySchema) error {
	for _, table := range schema.Tables {
		sql, err := sqlTableCreate(table)
		if err != nil {
			return fmt.Errorf("failed to prepare SQL statement: %w", err)
		}
		_, err = tx.Exec(context.Background(), sql)
		if err != nil {
			return fmt.Errorf("failed to create table: %s: %w", table.Name, err)
		}
	}

	// Store the new schema by converting it to core.Data and preparing a
	// saveContext including the schema itself
	d, err := schema.Data()
	if err != nil {
		return fmt.Errorf("failed to create data block from schema: %w", err)
	}
	node := newDataNode(&d)
	if err := psqlSaveNode(tx, node, schema); err != nil {
		return fmt.Errorf("failed to insert latest tables: %w", err)
	}

	return nil
}

func psqlResolveScalar(conn *pgx.Conn, params graphql.ResolveParams) (interface{}, error) {
	result, err := psqlResolveParams(conn, params)
	if err != nil {
		return nil, err
	}
	if val, ok := result.([]map[string]interface{}); ok {
		if len(val) > 0 {
			return val[0], nil
		}
		return nil, nil
	}
	return nil, fmt.Errorf("received non-list type: %s", reflect.TypeOf(result).String())
}

func psqlResolveParams(conn *pgx.Conn, params graphql.ResolveParams) (interface{}, error) {
	var (
		retVal     = make([]map[string]interface{}, 0)
		tableName  = params.Info.FieldName
		parent     = params.Source
		args       = params.Args
		fieldNames = queryFieldSelectionSet(params)
	)

	sql := psql.Select(fieldNames...).From(tableName).OrderBy(tableIDField + " DESC")
	sql = applyGraphQLArgs(sql, args)

	if isValidParent(parent) {
		// Get the parent type name, which is also the core.Table name of the
		// parent, and hence also the SQL table name
		var (
			parentName = params.Info.ParentType.Name()
			parentVal  = parent.(map[string]interface{})
		)

		// In this case, the tableName is the child
		if childBelongsToParent(params.Info, parentName, tableName) {
			parentID := parentVal[tableIDField]
			sql = sql.Where(sq.Eq{parentName + tableJoinSuffix: parentID})
		} else {
			childID := parentVal[tableName+tableJoinSuffix]
			sql = sql.Where(sq.Eq{tableIDField: childID})
		}
	}

	sqlStr, sqlArgs, err := sql.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build sql query: %w", err)
	}

	rows, err := conn.Query(context.Background(), sqlStr, sqlArgs...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute SQL query: %w", err)
	}

	defer rows.Close()

	// Iterate through the result set
	for rows.Next() {
		rowVals, err := psqlRowValues(rows, tableName, fieldNames)
		if err != nil {
			return nil, err
		}
		retVal = append(retVal, rowVals)
	}

	// Any errors encountered by rows.Next or rows.Scan will be returned here
	if rows.Err() != nil {
		return nil, fmt.Errorf("failed to scan row values: %w", rows.Err())
	}

	return retVal, nil
}

// psqlRowValues takes a row and returns the fields/values from the query in a map.
// To return the list of fields, they are added to a slice of pointers to
// interface{} in order to get the returned values from a pgx.Row.
// The fields are returned in a map with the field name as the key and returned
// value as the value.
func psqlRowValues(row pgx.Row, tableName string, fields []string) (map[string]interface{}, error) {
	retVal := make(map[string]interface{}, len(fields)+1)
	valuePtrs := make([]interface{}, len(fields))
	for i := range fields {
		var tmp interface{}
		valuePtrs[i] = &tmp
	}

	if err := row.Scan(valuePtrs...); err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("failed to scan fields: %s: %w", strings.Join(fields, ","), err)
		}
		return nil, nil
	}

	for i, field := range fields {
		// Get value from pointer to interface{} by casting to pointer of
		// interface{} and then dereference
		retVal[field] = *(valuePtrs[i].(*interface{}))
	}

	return retVal, nil
}
