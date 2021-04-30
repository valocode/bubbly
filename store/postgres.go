package store

import (
	"context"
	"errors"
	"fmt"
	"strings"

	sq "github.com/Masterminds/squirrel"
	"github.com/graphql-go/graphql"
	pgx "github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/log/zerologadapter"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/valocode/bubbly/api/core"
	"github.com/valocode/bubbly/config"
	"github.com/valocode/bubbly/env"
	"github.com/valocode/bubbly/parser"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/gocty"
)

var (
	psql = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
)

const (
	psqlBubblySchemaPrefix        = "bb_"
	psqlTableUniqueSuffix         = "_key"
	defaultStoreConnRetryAttempts = 10
	defaultStoreConnRetryTimeout  = "200ms"
)

var _ provider = (*postgres)(nil)

func newPostgres(bCtx *env.BubblyContext) (*postgres, error) {

	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s/%s",
		bCtx.StoreConfig.PostgresUser,
		bCtx.StoreConfig.PostgresPassword,
		bCtx.StoreConfig.PostgresAddr,
		bCtx.StoreConfig.PostgresDatabase,
	)

	pool, err := psqlNewPool(bCtx, connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize connection to db: %w", err)
	}

	return &postgres{
		pool: pool,
	}, nil
}

type postgres struct {
	pool *pgxpool.Pool
}

func (p *postgres) Close() {
	p.pool.Close()
}

func (p *postgres) Apply(tenant string, schema *bubblySchema) error {

	tx, err := p.pool.Begin(context.Background())
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(context.Background())

	err = psqlApplySchema(tx, tenant, schema)
	if err != nil {
		return fmt.Errorf("failed to apply tables: %w", err)
	}

	return tx.Commit(context.Background())
}

func (p *postgres) Migrate(tenant string, schema *bubblySchema, cl schemaUpdates) error {
	migration, err := psqlGenerateMigration(config.PostgresStore, tenant, schema, cl)
	if err != nil {
		return fmt.Errorf("failed to generate migration list: %w", err)
	}
	return psqlMigrate(p.pool, tenant, schema, migration)
}

func (p *postgres) Save(bCtx *env.BubblyContext, tenant string, graph *schemaGraph, tree dataTree) error {

	tx, err := p.pool.Begin(context.Background())
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(context.Background())

	// Create a callback function that wil be called for each node in the data
	// tree we visit and will save that node
	saveNode := func(bCtx *env.BubblyContext, node *dataNode, blocks *core.DataBlocks) error {
		// Check that the data node we are saving exists in the schema graph.
		// Otherwise it does not exist in our schema
		tNode, ok := graph.NodeIndex[node.Data.TableName]
		if !ok {
			return fmt.Errorf("data block refers to non-existing table: %s", node.Data.TableName)
		}
		return psqlSaveNode(tx, tenant, node, *tNode.table)
	}

	_, err = tree.traverse(bCtx, saveNode)

	if err != nil {
		return fmt.Errorf("failed to save data in postgres: %w", err)
	}

	return tx.Commit(context.Background())
}

func (p *postgres) ResolveQuery(tenant string, graph *schemaGraph, params graphql.ResolveParams) (interface{}, error) {
	return psqlResolveRootQueries(p.pool, tenant, graph, params)
}

func (p *postgres) Tenants() ([]string, error) {
	return psqlTenantSchemas(p.pool)
}

func (p *postgres) CreateTenant(name string) error {
	return psqlCreateSchema(p.pool, name)
}

func (p *postgres) HasTable(tenant string, table core.Table) (bool, error) {
	return psqlHasTable(p.pool, tenant, table)
}

func psqlNewPool(bCtx *env.BubblyContext, connStr string) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse db config: %w", err)
	}

	// TODO: log level should "scale" with bCtx.Logger level
	config.ConnConfig.Logger = zerologadapter.NewLogger(*bCtx.Logger)
	config.ConnConfig.LogLevel = pgx.LogLevelError

	pool, err := pgxpool.ConnectConfig(context.Background(), config)
	if err != nil {
		return nil, fmt.Errorf("failed to start database connection pool: %w", err)
	}

	return pool, nil
}

func psqlTenantSchemas(pool *pgxpool.Pool) ([]string, error) {
	var (
		sql = psql.Select("schema_name").
			From("information_schema.schemata")
		schemas = make([]string, 0)
	)

	sqlStr, _, err := sql.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to create sql query: %w", err)
	}

	rows, err := pool.Query(context.Background(), sqlStr)
	if err != nil {
		return nil, fmt.Errorf("failed to execute SQL query: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var schema string
		if err := rows.Scan(&schema); err != nil {
			return nil, fmt.Errorf("failed to scan schema value: %w", err)
		}
		// Check that the schema is a bubbly schema
		if strings.HasPrefix(schema, psqlBubblySchemaPrefix) {
			// Append the bubbly schema and remove the prefix to get the tenant name
			schemas = append(schemas, schema[len(psqlBubblySchemaPrefix):])
		}
	}

	return schemas, nil
}

func psqlCreateSchema(pool *pgxpool.Pool, name string) error {
	var (
		schemaName = psqlBubblySchemaPrefix + name
		sqlStr     = "CREATE SCHEMA IF NOT EXISTS " + schemaName
	)

	_, err := pool.Exec(context.Background(), sqlStr)
	if err != nil {
		return fmt.Errorf("failed to execute SQL: %w", err)
	}
	return nil
}

func psqlHasTable(pool *pgxpool.Pool, tenant string, table core.Table) (bool, error) {
	var (
		sql = psql.Select("1").
			Prefix("SELECT EXISTS (").
			From("information_schema.tables").
			Where(sq.Eq{"table_schema": psqlSchemaName(tenant)}).
			Where(sq.Eq{"table_name": table.Name}).
			Suffix(");")
		exists bool
	)
	sqlStr, sqlArgs, err := sql.ToSql()
	if err != nil {
		return false, fmt.Errorf("failed to create sql query: %w", err)
	}

	row := pool.QueryRow(context.Background(), sqlStr, sqlArgs...)
	if err := row.Scan(&exists); err != nil {
		return false, fmt.Errorf("failed to get table status: %s: %w", table.Name, err)
	}

	return exists, nil
}

func psqlApplySchema(tx pgx.Tx, tenant string, schema *bubblySchema) error {
	for _, table := range schema.Tables {
		if err := psqlApplyTable(tx, tenant, table); err != nil {
			return err
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

	return nil
}

func psqlApplyTable(tx pgx.Tx, tenant string, table core.Table) error {
	sql, err := psqlTableCreate(tenant, table)
	if err != nil {
		return fmt.Errorf("failed to prepare SQL statement: %w", err)
	}
	// Create the table
	_, err = tx.Exec(context.Background(), sql)
	if err != nil {
		return fmt.Errorf("failed to create table: %s: %w", table.Name, err)
	}
	// Apply the unique constraints
	sql = psqlTableUniqueConstraints(tenant, table)
	_, err = tx.Exec(context.Background(), sql)
	if err != nil {
		return fmt.Errorf("failed to add constraints on table: %s: %w", table.Name, err)
	}
	return nil
}

func psqlSaveNode(tx pgx.Tx, tenant string, node *dataNode, table core.Table) error {
	var (
		sqlStr  string
		sqlArgs []interface{}
		err     error
	)

	// Before generating the SQL query, we want to make sure we have default
	// values for all the unique fields (and joins) in a data block.
	// Postgres treats null as a unique value, so the unique constraint would be
	// violated (according to the bubbly schema) as we can have multiple records
	// with the same values (if any of those values is null)
	if node.Data.Policy != core.ReferencePolicy {
		err = psqlAddUniqueDataFields(table, node.Data)
		if err != nil {
			return fmt.Errorf("failed to add default unique values to data block %s: %w", node.Data.TableName, err)
		}
	}

	switch node.Data.Policy {
	case core.CreateUpdatePolicy, core.EmptyPolicy:
		// We should perform an Upsert (Insert/Update)
		sqlStr, sqlArgs, err = psqlDataCreateUpdate(tenant, node, table)
	case core.CreatePolicy:
		// We should perform an Insert
		sqlStr, sqlArgs, err = psqlDataCreate(tenant, node, table)
	case core.ReferencePolicy:
		// We should perform a Select
		sqlStr, sqlArgs, err = psqlDataReference(tenant, node, table)
	default:
		return fmt.Errorf("data block refers to unsupported policy %s: %s", node.Data.TableName, node.Data.Policy)
	}
	if err != nil {
		return fmt.Errorf("failed to create SQL statement for data block %s with policy %s: %w", node.Data.TableName, node.Data.Policy, err)
	}

	rows, err := tx.Query(context.Background(), sqlStr, sqlArgs...)
	// Execute the query
	if err != nil {
		return fmt.Errorf("failed to execute query for data block %s: %w", node.Data.TableName, err)
	}
	defer rows.Close()

	var (
		retValues = make(map[string]interface{})
		numRows   = 0
	)
	for rows.Next() {
		if numRows > 0 {
			return fmt.Errorf("received more than one row data block %s: %#v", node.Data.TableName, node.Data)
		}
		numRows += 1
		retValues, err = psqlRowValues(rows, table.Name, node.orderedRefFields())
		if err != nil {
			return fmt.Errorf("failed to insert data block: %s: %w", table.Name, err)
		}
		if rows.Err() != nil {
			return fmt.Errorf("error while reading SQL row for data block %s: %w", node.Data.TableName, rows.Err())
		}
	}
	if numRows == 0 {
		return fmt.Errorf("received no rows for data block %s: %#v", node.Data.TableName, node.Data)
	}
	// Asign the returned values so that if the child nodes need to resolve
	// their data references they have values to do so
	node.Return = retValues

	return nil
}

// psqlRowValues takes a row and returns the fields/values from the query in a map.
// To return the list of fields, they are added to a slice of pointers to
// interface{} in order to get the returned values from a pgx.Row.
// The fields are returned in a map with the field name as the key and returned
// value as the value.
func psqlRowValues(row pgx.Row, tableName string, fields []string) (map[string]interface{}, error) {
	var (
		retVal        = make(map[string]interface{}, len(fields)+1)
		scanValues    = make([]interface{}, len(fields))
		scanValuePtrs = make([]interface{}, len(fields))
	)
	for i := 0; i < len(fields); i++ {
		scanValuePtrs[i] = &scanValues[i]
	}

	if err := row.Scan(scanValuePtrs...); err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("failed to scan fields: %s: %w", strings.Join(fields, ","), err)
		}
		return nil, nil
	}

	for i, field := range fields {
		// Get value from pointer to interface{} by casting to pointer of
		// interface{} and then dereference
		retVal[field] = scanValues[i]
	}

	return retVal, nil
}

func psqlTableUniqueConstraints(tenant string, table core.Table) string {
	var (
		uniqueFields = make([]string, 0)
	)
	for _, field := range table.Fields {
		if field.Unique {
			uniqueFields = append(uniqueFields, field.Name)
		}
	}
	// Add the joins as fields to the SQL table
	for _, join := range table.Joins {
		fieldName := join.Table + "_id"
		if join.Unique {
			uniqueFields = append(uniqueFields, fieldName)
		}
	}

	// First drop the existing constraint (IF EXISTS)
	sql := "ALTER TABLE " + psqlAbsTableName(tenant, table.Name) +
		" DROP CONSTRAINT IF EXISTS " + table.Name + psqlTableUniqueSuffix
	// If we have some unique fields on which to add a constraint, then add it
	if len(uniqueFields) > 0 {
		sql += ", ADD CONSTRAINT " + table.Name + psqlTableUniqueSuffix +
			" UNIQUE (" + strings.Join(uniqueFields, ",") + ")"
	}
	return sql + ";"
}

func psqlTableCreate(tenant string, table core.Table) (string, error) {
	var (
		fieldLen    = len(table.Fields) + len(table.Joins)
		tableFields = make([]string, 0, fieldLen)
		// uniqueFields = make([]string, 0)
	)

	tableFields = append(tableFields, tableIDField+" SERIAL PRIMARY KEY")
	// Add the fields to the SQL table
	for _, field := range table.Fields {
		sqlType, err := psqlType(field.Type)
		if err != nil {
			return "", fmt.Errorf("failed to create SQL statement for table: %s: %w", table.Name, err)
		}
		tableFields = append(tableFields, field.Name+" "+sqlType)

		// if field.Unique {
		// 	uniqueFields = append(uniqueFields, field.Name)
		// }
	}
	// Add the joins as fields to the SQL table
	for _, join := range table.Joins {
		fieldName := join.Table + "_id"
		tableFields = append(tableFields, fieldName+" INT8")
		// if join.Unique {
		// 	uniqueFields = append(uniqueFields, fieldName)
		// }
	}

	// if len(uniqueFields) > 0 {
	// 	tableFields = append(tableFields, "UNIQUE ("+strings.Join(uniqueFields, ",")+")")
	// }

	return "CREATE TABLE IF NOT EXISTS " + psqlAbsTableName(tenant, table.Name) + " ( " + strings.Join(tableFields, ",") + " );", nil
}

// psqlDataCreateUpdate generates a sql query for performing an insert/update
// (i.e. "upsert") query for a given data node
func psqlDataCreateUpdate(tenant string, node *dataNode, table core.Table) (string, []interface{}, error) {
	var (
		data           = node.Data
		fieldNames     = node.orderedFields()
		conflictValues = make([]string, 0, len(data.Fields))
		uniqueFields   = make([]string, 0)
		sqlOnConflict  = ""
		sqlReturning   = ""
	)

	for _, field := range table.Fields {
		if field.Unique {
			uniqueFields = append(uniqueFields, field.Name)
		}
	}
	for _, join := range table.Joins {
		if join.Unique {
			uniqueFields = append(uniqueFields, join.Table+tableJoinSuffix)
		}
	}

	for _, name := range fieldNames {
		conflictValues = append(conflictValues, name+"=EXCLUDED."+name)
	}

	// Create the UPSERT / ON CONFLICT part of the SQL statement, if any.
	if len(uniqueFields) > 0 {
		sqlOnConflict = fmt.Sprintf(
			"ON CONFLICT ON CONSTRAINT %s DO UPDATE SET %s",
			table.Name+psqlTableUniqueSuffix,
			strings.Join(conflictValues, ","),
		)
	}

	// Create the RETURNING part of the SQL statement, if any.
	sqlReturning = "RETURNING " + strings.Join(node.orderedRefFields(), ",")
	values, err := psqlArgValues(node)
	if err != nil {
		return "", nil, fmt.Errorf("failed to get SQL arguments: %w", err)
	}
	sql := psql.Insert(psqlAbsTableName(tenant, data.TableName)).
		Columns(fieldNames...).
		Values(values...).
		Suffix(sqlOnConflict).
		Suffix(sqlReturning)
	return sql.ToSql()
}

// psqlDataCreate generates a sql query for performing an insert, which will
// error if any uniqueness constraints are violated
func psqlDataCreate(tenant string, node *dataNode, table core.Table) (string, []interface{}, error) {
	var (
		data         = node.Data
		fieldNames   = node.orderedFields()
		sqlReturning = ""
	)

	// Create the RETURNING part of the SQL statement, if any.
	sqlReturning = "RETURNING " + strings.Join(node.orderedRefFields(), ",")
	values, err := psqlArgValues(node)
	if err != nil {
		return "", nil, fmt.Errorf("failed to get SQL arguments: %w", err)
	}
	sql := psql.Insert(psqlAbsTableName(tenant, data.TableName)).
		Columns(fieldNames...).
		Values(values...).
		Suffix(sqlReturning)
	return sql.ToSql()
}

// psqlDataReference generates a sql query for performing an select, returning
// the row which the given data node refers to. This is used so that joins can
// be made to a piece of data
func psqlDataReference(tenant string, node *dataNode, table core.Table) (string, []interface{}, error) {
	var (
		data      = node.Data
		refFields = node.orderedRefFields()
	)

	// Return the fields which are being referenced by this data node
	sql := psql.Select(refFields...).
		From(psqlAbsTableName(tenant, data.TableName))

	// Iterate over the field values that have been provided and create the SQL
	// WHERE clause so that we get the correct record back
	for name, value := range node.Data.Fields {
		v, err := psqlValue(node, value)
		if err != nil {
			return "", nil, fmt.Errorf("failed to get SQL value from data block field %s.%s: %w", node.Data.TableName, name, err)
		}
		sql = sql.Where(sq.Eq{name: v})
	}
	return sql.ToSql()
}

// psqlArgValues takes a data node and returns the values of for the fields
// that have been provided
func psqlArgValues(node *dataNode) ([]interface{}, error) {
	var (
		data   = node.Data
		values = make([]interface{}, 0, len(data.Fields))
	)
	// We need to order the fields to make sure the list of values we give
	// match up to the list of fields names
	for _, f := range node.orderedFields() {
		val, err := psqlValue(node, data.Fields[f])
		if err != nil {
			return nil, fmt.Errorf("failed to get SQL value from cty.Value for field: %s: %w", f, err)
		}
		values = append(values, val)
	}
	return values, nil
}

func psqlValue(node *dataNode, val cty.Value) (interface{}, error) {
	// Check if the value is a capsule value, in which case it needs special
	// treatment
	if val.Type().IsCapsuleType() {
		// Get the underlying DataRef type
		ref := val.EncapsulatedValue().(*parser.DataRef)
		if parent, ok := node.Parents[ref.TableName]; ok {
			if ret, ok := parent.Return[ref.Field]; ok {
				return ret, nil
			}
		}
		return nil, fmt.Errorf("could not find data ref: %s.%s", ref.TableName, ref.Field)
	}

	// If not a capsule type, it is a regular cty.Value
	return valueFromCty(val)
}

func valueFromCty(val cty.Value) (interface{}, error) {
	switch ty := val.Type(); {
	case ty == cty.Bool:
		return val.True(), nil
	case ty == cty.Number:
		var number int
		if err := gocty.FromCtyValue(val, &number); err != nil {
			return nil, fmt.Errorf("failed to convert cty.Value to int: %s: %w", val.GoString(), err)
		}
		return number, nil
	case ty == cty.String:
		return val.AsString(), nil
	case ty.IsObjectType():
		var (
			m   = val.AsValueMap()
			ret = make(map[string]interface{}, len(m))
		)
		for k, v := range m {
			var err error
			ret[k], err = valueFromCty(v)
			if err != nil {
				return nil, err
			}
		}
		return ret, nil
	default:
		return nil, fmt.Errorf("unsupported cty type: %s", ty.GoString())
	}
}

// sqlType takes a cty.Type and returns a string representation of the
// corresponding SQL type
func psqlType(ty cty.Type) (string, error) {
	switch {
	case ty == cty.Bool:
		return "BOOL", nil
	case ty == cty.Number:
		return "INT8", nil
	case ty == cty.String:
		return "TEXT", nil
	case ty.IsObjectType():
		return "JSONB", nil
	default:
		return "", fmt.Errorf("unsupported SQL type: %s", ty.GoString())
	}
}

func psqlSchemaName(tenant string) string {
	return psqlBubblySchemaPrefix + tenant
}

func psqlAbsTableName(tenant string, table string) string {
	return psqlBubblySchemaPrefix + tenant + "." + table
}

func psqlAddUniqueDataFields(table core.Table, data *core.Data) error {
	for _, field := range table.Fields {
		if field.Unique {
			if _, ok := data.Fields[field.Name]; !ok {
				val, err := psqlDefaultFieldValue(field.Type)
				if err != nil {
					return fmt.Errorf("failed to get default value for unique table field %s.%s: %w", table.Name, field.Name, err)
				}
				data.Fields[field.Name] = val
			}
		}
	}
	for _, join := range table.Joins {
		if join.Unique {
			fieldName := join.Table + tableJoinSuffix
			if _, ok := data.Fields[fieldName]; !ok {
				// The forgeign key can never be -1, and so if we are generating
				// a default value we want to make sure this will not actually
				// create an unintential join!
				// We need this default value because null is a unique value in
				// postgres... which screws with our unique constraints
				data.Fields[fieldName] = cty.NumberIntVal(-1)
			}
		}
	}
	return nil
}

func psqlDefaultFieldValue(ty cty.Type) (cty.Value, error) {
	switch {
	case ty == cty.Bool:
		return cty.BoolVal(false), nil
	case ty == cty.Number:
		return cty.NumberIntVal(0), nil
	case ty == cty.String:
		return cty.StringVal(""), nil
	case ty.IsObjectType():
		return cty.EmptyObjectVal, nil
	default:
		return cty.NilVal, fmt.Errorf("unsupported cty.Type: %s", ty.GoString())
	}
}
