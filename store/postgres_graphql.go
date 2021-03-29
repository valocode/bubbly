package store

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	sq "github.com/Masterminds/squirrel"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/language/ast"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

const (
	orderAsc  string = "ASC"
	orderDesc string = "DESC"
)

func newSQLQueryBuilder() *sqlQueryBuilder {
	return &sqlQueryBuilder{
		sql:   psql.Select(),
		depth: 0,
	}
}

// sqlQueryBuilder stores some context that is used for constructing the sql
// statement that reflects a graphql query
type sqlQueryBuilder struct {
	sql   sq.SelectBuilder
	node  *schemaNode
	depth int
}

// tableColumns is used to store the columns that are SELECT'd in a SQl
// statement, within one single table.
// This is quite a complex problem because of GraphQL queries have a hierarchy
// and SQL queries return flat rows. How do we "unpack" flat rows into a hierarchy?
// tableColumns stores the hierarchy of the graphql query and enables us to map
// the returned SQL row to the structure defined by the GraphQL query
type tableColumns struct {
	table    string
	alias    string
	fields   []string
	scalar   bool
	children []*tableColumns
}

// length returns the number of fields in this tableColumns, which includes
// all the fields in all the descendents (children of children of children...)
func (t *tableColumns) length() int {
	var count = len(t.fields)
	for _, tt := range t.children {
		count += tt.length()
	}
	return count
}

// FIXME: ofr I have a hunch that the following comment is related to GraphQL standard,
//        which is rather prescriptive about how things should be evaluated, and also
//        to the reference implementation `graphql-js`. Investigate later.
//
// What is a bit puzzling is that if you have a query with two fields, this
// method gets called twice, once for each field, but each time the
// graphql.ResolveParams contains a list of FieldASTs with one element:

// psqlResolveRootQueries is called for each top-level query and iterates
// through the fields in that root query and resolves them.
func psqlResolveRootQueries(pool *pgxpool.Pool, tenant string, graph *schemaGraph, params graphql.ResolveParams) (interface{}, error) {
	var (
		result interface{}
		err    error
	)
	for _, field := range params.Info.FieldASTs {
		result, err = psqlResolveRootQuery(pool, tenant, graph, field)
		if err != nil {
			return nil, fmt.Errorf("failed to resolve query: %s: %w", field.Name.Value, err)
		}
	}
	return result, err
}

// psqlResolveRootQuery resolves a single root graphql query
func psqlResolveRootQuery(pool *pgxpool.Pool, tenant string, graph *schemaGraph, field *ast.Field) (interface{}, error) {
	var (
		result      = make(map[string]interface{})
		qb          = newSQLQueryBuilder()
		rootTable   = field.Name.Value
		rootAlias   = tableAlias(rootTable, 0)
		rootColumns = tableColumns{
			table:  rootTable,
			alias:  rootAlias,
			scalar: false,
		}
	)

	// Set the starting node and initialize the sql statement
	qb.node = graph.NodeIndex[rootTable]
	qb.sql = qb.sql.From(tableAsAlias(psqlAbsTableName(tenant, rootTable), rootAlias))
	// qb.columns = &rootColumns

	// Recursively go through the graphql query and resolve the sub-fields
	if err := psqlSubQuery(tenant, graph, qb, &rootColumns, &rootColumns, field, nil); err != nil {
		return nil, fmt.Errorf("failed to process root query: %s: %w", rootTable, err)
	}

	// Create the sql query and any arguments
	sqlStr, sqlArgs, err := qb.sql.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to create sql query: %w", err)
	}

	// Execute the query
	rows, err := pool.Query(context.Background(), sqlStr, sqlArgs...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute SQL query: %s: %w", sqlStr, err)
	}

	defer rows.Close()

	// Iterate through the result set and append each row of results to the
	// result value we are returning
	for rows.Next() {
		if err := psqlScanRowColumns(rows, result, rootColumns); err != nil {
			return nil, fmt.Errorf("failed scanning row values: %w", err)
		}
	}
	return result[rootTable], nil
}

func psqlSubQuery(tenant string, graph *schemaGraph, qb *sqlQueryBuilder, root *tableColumns, tc *tableColumns, field *ast.Field, path schemaPath) error {

	// GraphQL fields are conceptually functions which return values,
	// and occasionally accept arguments which alter their behaviour.
	//
	// This recursive function compiles a SQL query for the `field`,
	// taking into account (optional) arguments for the `field`.
	//
	// It is recursive because the `field` may have a selection set
	// associated with it, which in turn also requires a SQL (sub)query.
	//
	// Relevant parts of the GraphQL spec:
	//   http://spec.graphql.org/June2018/#sec-Language.Fields
	//   http://spec.graphql.org/June2018/#sec-Language.Arguments
	//

	// Create the tableColumns for this type/table in the query
	var (
		node      = qb.node
		subFields = make([]*ast.Field, 0)

		// The `order_by` GraphQL argument can only be processed after all table aliases
		// are known. That includes the tables referenced by the GraphQL subfields.
		// Therefore, upon encountering the `order_by` argument in one of the root
		// types, it is stored in this variable for processing after all the table
		// aliases are known, that is after the rest of the GraphQL query had been processed.
		orderByArg *ast.Argument
	)

	// FIXME: Why is the depth being increased here?
	qb.depth++

	// Always return the ID field of a table as the first row as we need it when
	// we aggregate the results up into the returned value
	tc.fields = append(tc.fields, tableIDField)
	qb.sql = qb.sql.Column(tc.alias + "." + tableIDField)

	// GraphQL arguments are processed here
	for _, arg := range field.Arguments {

		// The arguments can be of different kinds, from simple column names
		// to our custom function names. We'll work through the possibilities,
		// but if at the end the argument has not been resolved, raise an error.
		argIsResolved := false

		// Argument name equal to one of the column names for the current node (table)
		// adds an equality predicate in the WHERE clause.
		// Multiple expressions are `AND`ed together in the generated SQL.
		for _, tf := range node.table.Fields {
			if arg.Name.Value == tf.Name {
				qb.sql = qb.sql.Where(sq.Eq{tc.alias + "." + arg.Name.Value: arg.Value.GetValue()})
				argIsResolved = true
				break
			}
		}

		if argIsResolved {
			continue
		}

		// FIXME: how do I do process the args in a more abstract way, without having to fiddle with AST objects?

		// Process the arguments that are not GraphQL/DB field/column names...

		// FIXME: maybe use a switch for what follows, instead of a set of ifs?

		// The order_by argument is allowed only at the top level. Futhermore, it cannot be processed until
		// all subfields had been processed, because at the top level the alias names of tables are not known.
		// Therefore, defer the processing of this argument by saving a pointer to it for later processing.
		if arg.Name.Value == "order_by" {

			if qb.depth != 1 {
				return fmt.Errorf("the `order_by` argument is supported only for the root types")
			}

			orderByArg = arg

			argIsResolved = true
		}

		// TODO: distinct_on
		if arg.Name.Value == "distinct_on" {
			return fmt.Errorf("argument not implemented: %s", arg.Name.Value)
		}

		// The argument name which is not a column name is a mistake, raise error.
		if !argIsResolved {
			return fmt.Errorf("unknown identifier %s.%s", tc.table, arg.Name.Value)
		}
	}

	// Iterate over the fields in the selection set (if any) for the current `field`
	for _, selection := range field.SelectionSet.Selections {
		// Only GraphQL `Field`s are supported at this point. http://spec.graphql.org/June2018/#sec-Language.Fields
		// The `Selection` interface is implemented by the `ast.Field` type in this supported case.
		subField, ok := selection.(*ast.Field)
		if !ok {
			return fmt.Errorf("graphql query selection type not supported: %s", selection.GetSelectionSet().Kind)
		}

		fieldName := subField.Name.Value

		// Types and fields required by the GraphQL introspection system that are used
		// in the same context as user-defined types and fields are prefixed
		// with "__" two underscores. This in order to avoid naming collisions
		// with user-defined GraphQL types. http://spec.graphql.org/June2018/#sec-Reserved-Names

		// FIXME shouldn't we raise an error instead of skipping quietly?
		// We skip this field if it has a reserved name
		if strings.HasPrefix(fieldName, "__") {
			continue
		}

		// A non-nil selection set implies that the subField refers to another
		// table in our schema.
		// We need to process the columns/fields for this table first, before we
		// process any subFields, so simply append these to a slice and process
		// at the end of the function
		if subField.SelectionSet != nil {
			subFields = append(subFields, subField)
			continue
		} else {
			// If subField did not have a selection set this it is just a column
			// within the current table, so add it to the columns
			tc.fields = append(tc.fields, fieldName)
			qb.sql = qb.sql.Column(tableColumn(tc.alias, fieldName))
		}
	}

	// Once we have processed this fields columns, proceed to the subFields.
	// This is to ensure the correct order of columns in the SQL SELECT statement
	for _, subField := range subFields {

		// TODO: instead of searching the graph, check the node's edges for their destinations, and use those dest. in JOINs
		// In the following example, the current qb.node is `A`, and subField is `B`, and we have just discovered, that `B`
		// refers to another table.
		//
		// A {
		//     name
		//     B {
		//        name
		//     }
		// }

		// Are the parent field and the subfield connected in the graph at all?
		var (
			fieldName         = subField.Name.Value
			edgeToRelatedNode *schemaEdge
		)
		for _, p := range node.edges {
			if p.node.table.Name == fieldName {
				edgeToRelatedNode = p
			}
		}
		if edgeToRelatedNode == nil {
			return fmt.Errorf("no relationship found between tables: '%s', '%s'", node.table.Name, fieldName)
		}

		var (
			leftTable       = tc.table
			leftTableAlias  = tc.alias
			rightTable      = edgeToRelatedNode.node.table.Name
			rightTableAlias = tableAlias(rightTable, qb.depth)
		)
		switch edgeToRelatedNode.rel {
		case oneToOne, oneToMany:
			qb.sql = qb.sql.LeftJoin(joinOn(
				tableAsAlias(psqlAbsTableName(tenant, rightTable), rightTableAlias),
				tableColumn(leftTableAlias, tableIDField),
				tableColumn(rightTableAlias, foreignKeyField(leftTable))))
		case belongsTo:
			qb.sql = qb.sql.LeftJoin(
				joinOn(
					tableAsAlias(psqlAbsTableName(tenant, rightTable), rightTableAlias),
					tableColumn(rightTableAlias, tableIDField),
					tableColumn(leftTableAlias, foreignKeyField(rightTable))))
		}

		// Recursively resolve for the subField `B`, which may contain further nested fields.
		qb.node = graph.NodeIndex[fieldName]
		subColumns := &tableColumns{
			table:  fieldName,
			alias:  tableAlias(fieldName, qb.depth),
			scalar: edgeToRelatedNode.isScalar(),
		}
		if err := psqlSubQuery(tenant, graph, qb, root, subColumns, subField, path); err != nil {
			return err
		}
		tc.children = append(tc.children, subColumns)
	}

	// TODO: make error messages more helpful (and more systematic) in the following code
	// Process order_by, if available. At this point all table aliases are known.
	if orderByArg != nil {

		arg := orderByArg

		// This argument's value is a slice of values representing (table, field, order) tuples.
		orderByItems, ok := (arg.Value.GetValue()).([]ast.Value)
		if !ok {
			return fmt.Errorf("argument type is wrong in `order_by`: %s : %#v", arg.Name.Value, arg.Value.GetValue())
		}

		// Process each of the {table, field, order} objects...
		for _, item := range orderByItems {

			fields, ok := item.GetValue().([]*ast.ObjectField)
			if !ok {
				return fmt.Errorf("failed to type cast argument value for `order_by`, step 2, value: %#v", item.GetValue())
			}

			// FIXME probably best to make it a proper type {table,field, order}
			ob := map[string]string{}

			// Process a single {table, field, order} entry
			for _, objectField := range fields {

				switch objectField.Name.Value {
				case "table":
					s, ok := objectField.GetValue().(*ast.StringValue)
					if !ok {
						return fmt.Errorf("failed to type cast argument value for `order_by`, step 3: %#v: %#v", objectField.Name.Value, objectField)
					}
					ob["table"], _ = s.GetValue().(string) // FIXME: is this safe?

				case "field":
					s, ok := objectField.GetValue().(*ast.StringValue)
					if !ok {
						return fmt.Errorf("failed to type cast argument value for `order_by`, step 3, value: %#v", objectField)
					}
					ob["field"] = s.GetValue().(string) // FIXME: is this safe?

				case "order":
					s, ok := objectField.GetValue().(*ast.StringValue)
					if !ok {
						return fmt.Errorf("failed to type cast argument value for `order_by`, step 3, value: %#v", objectField)
					}
					ob["order"] = s.GetValue().(string) // FIXME: is this safe?
				}
			}

			tableName := ob["table"]
			f := ob["field"]
			order := strings.ToUpper(ob["order"])

			if order != orderAsc && order != orderDesc {
				return fmt.Errorf("invalid sorting order for: %s", ob["order"])
			}

			// Lookup the table alias by table name
			ta := findAliasFor(tableName, root)
			if ta == "" {
				return fmt.Errorf("could not find the table alias for: %s", tableName)
			}

			// Augment the SQL query with the ORDER BY statement.
			qb.sql = qb.sql.OrderBy(tableColumn(ta, f) + " " + order)
		}
	}

	return nil
}

// findAliasFor takes a table name and traverses the tableColumns tree recursively, looking for the first match
// on the table name. On match, it returns the stored table name alias for that table name.
// If the the table name cannot be found in the tree, the empty string is returned.
func findAliasFor(tableName string, tc *tableColumns) string {
	if tc.table == tableName {
		return tc.alias
	} else {
		for _, subtc := range tc.children {
			if alias := findAliasFor(tableName, subtc); alias != "" {
				return alias
			}
		}
	}
	return ""
}

func joinOn(table, leftColumn, rightColumn string) string {
	return table + " ON ( " + leftColumn + " = " + rightColumn + " ) "
}

func foreignKeyField(table string) string {
	return table + tableJoinSuffix
}

func tableColumn(table, column string) string {
	return table + "." + column
}

func tableAlias(table string, depth int) string {
	return table + "_" + strconv.Itoa(depth)
}

func tableAsAlias(table, alias string) string {
	return table + " AS " + alias
}

// psqlScanRowColumns takes a single row (that is returned from a SQL query),
// an existing result map that the row should be aggregated into, and a tableColumns
// type which contains the columns in the row, grouped by their tables.
// The rowColumns are used to unpack the row values into the result
func psqlScanRowColumns(row pgx.Row, result map[string]interface{}, columns tableColumns) error {
	var (
		columnLen     = columns.length()
		index         = 0
		scanValues    = make([]interface{}, columnLen)
		scanValuePtrs = make([]interface{}, columnLen)
	)
	// Initialize scanValues to store the values returned from the SQL row
	for i := 0; i < columnLen; i++ {
		scanValuePtrs[i] = &scanValues[i]
	}

	// Scan the values from the row into scanValues
	if err := row.Scan(scanValuePtrs...); err != nil {
		return fmt.Errorf("failed to scan values: %w", err)
	}

	psqlScanTableColumns(result, columns, scanValues, &index)

	return nil
}

// psqlScanTableColumns is called recursively for each child in tableColumns.
// A child in tableColumns means a graphql field was nested within another field,
// and hence the returned value should reflect this structure and the children
// of tableColumns should be nested within this tableColumns
func psqlScanTableColumns(parentVal map[string]interface{}, tc tableColumns, scanValues []interface{}, index *int) {
	var tColVal map[string]interface{}
	// Check if the value for this table of columns already exists.
	// If not, initialize it
	tVal, ok := parentVal[tc.table]
	if !ok {
		// Initialize the value from the scanned results for the group of
		// columns in this table
		tColVal = make(map[string]interface{}, len(tc.fields))
		for _, field := range tc.fields {
			tColVal[field] = scanValues[*index]
			*index++
		}
		// Check if we expect the result to be a scalar value or a list.
		// It is scalar depending on the relationship between the tables
		if tc.scalar {
			tVal = tColVal
		} else {
			tVal = make([]map[string]interface{}, 0, 1)
			tVal = append(tVal.([]map[string]interface{}), tColVal)
		}
		// Set the value for this table back into parent after it has been
		// initialized, and set the new parent value to this table column
		// value
		parentVal[tc.table] = tVal
	} else if tc.scalar {
		// If the table value should be scalar and the value is already
		// initialized, then we do not need to do anything.
		// Set the new parentVal and continue to the next table
		tColVal = tVal.(map[string]interface{})
		*index += len(tc.fields)
	} else {
		// If the parentVal already contained a value for this table.
		// Get the ID value from this result row
		var (
			tableIDVal = scanValues[*index]
			tListVal   = tVal.([]map[string]interface{})
		)
		// If the value for this table already exists, and we expect a list,
		// then we need to check if the value already exists in the list or if
		// we should add it
		for _, val := range tListVal {
			// If the ID of the table already exists, then we should append to
			// this value
			if tableIDVal == val[tableIDField] {
				tColVal = val
				break
			}
		}
		// If the value did not yet exist, we need to initialize it and append
		if tColVal == nil {
			tColVal = make(map[string]interface{})
			for _, field := range tc.fields {
				tColVal[field] = scanValues[*index]
				*index++
			}
			tListVal = append(tListVal, tColVal)

		} else {
			// Make sure we increment the index
			*index += len(tc.fields)
		}
		parentVal[tc.table] = tListVal
	}
	// Iterate through the children and unpack the remaining scanValues (starting
	// from the given index) into the given tColVal (which holds the value for
	// this tableColumns)
	for _, child := range tc.children {
		psqlScanTableColumns(tColVal, *child, scanValues, index)
	}
}
