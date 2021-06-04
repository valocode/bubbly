package store

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	sq "github.com/Masterminds/squirrel"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/language/ast"
	"github.com/graphql-go/graphql/language/kinds"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/valocode/bubbly/parser"
)

const (
	orderAsc     string = "ASC"
	orderDesc    string = "DESC"
	defaultLimit uint64 = 100
)

// tableColumns is used to store the columns that are SELECT'd in a SQl
// statement, within one single table.
// This is quite a complex problem because of GraphQL queries have a hierarchy
// and SQL queries return flat rows. How do we "unpack" flat rows into a hierarchy?
// tableColumns stores the hierarchy of the graphql query and enables us to map
// the returned SQL row to the structure defined by the GraphQL query
type tableColumns struct {
	table   string
	alias   string
	columns []string
	scalar  bool
	// The GraphQL Field for this table
	field    *ast.Field
	children []*tableColumns
}

// length returns the number of fields in this tableColumns, which includes
// all the fields in all the descendents (children of children of children...)
func (t *tableColumns) length() int {
	var count = len(t.columns)
	for _, tt := range t.children {
		count += tt.length()
	}
	return count
}

// psqlResolveRootQueries is called for each top-level query and iterates
// through the fields in that root query and resolves them.
func psqlResolveRootQueries(pool *pgxpool.Pool, tenant string, graph *SchemaGraph, params graphql.ResolveParams) (interface{}, error) {
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
func psqlResolveRootQuery(pool *pgxpool.Pool, tenant string, graph *SchemaGraph, field *ast.Field) (interface{}, error) {
	var (
		result      = make(map[string]interface{})
		rootTable   = field.Name.Value
		rootAlias   = tableAlias(rootTable, 0)
		rootColumns = tableColumns{
			table:  rootTable,
			alias:  rootAlias,
			field:  field,
			scalar: false,
		}
		rootSQL = sq.Select()
	)

	// Recursively go through the graphql query and resolve the sub-fields
	err := psqlSubQuery(tenant, graph, &rootSQL, nil, &rootColumns, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to process root query: %s: %w", rootTable, err)
	}

	// Create the sql query and any arguments
	sqlStr, sqlArgs, err := rootSQL.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to create sql query: %w", err)
	}

	// Change the default placeholder with $ for postgres
	sqlStr, err = sq.Dollar.ReplacePlaceholders(sqlStr)
	if err != nil {
		return nil, fmt.Errorf("error replacing the SQL (squirrel) placeholders: %w", err)
	}

	fmt.Println("SQL: " + sqlStr)

	// Execute the query
	rows, err := pool.Query(context.Background(), sqlStr, sqlArgs...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute SQL query: %s: %w", sqlStr, err)
	}
	defer rows.Close()

	// Iterate through the result set and append each row of results to the
	// result value we are returning. We should check if there are no rows
	// in which case we want to return at least an empty slice
	var hasRows bool
	for rows.Next() {
		hasRows = true
		if err := psqlScanRowColumns(rows, result, rootColumns); err != nil {
			return nil, fmt.Errorf("failed scanning row values: %w", err)
		}
	}
	if !hasRows {
		// Initialize with an empty slice to avoid returning just null
		result[rootTable] = make([]interface{}, 0)
	}
	return result[rootTable], nil
}

func psqlSubQuery(tenant string, graph *SchemaGraph, sql *sq.SelectBuilder, parent *tableColumns, tc *tableColumns, depth int) error {

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
		node = graph.NodeIndex[tc.table]
		// nodeQuery is the subquery for the current node
		nodeQuery = sq.Select().From(tableAsAlias(psqlAbsTableName(tenant, tc.table), tc.alias))
		subFields = make([]*ast.Field, 0)

		// In SQL terms, notNull says whether to perform a LEFT or an INNER JOIN
		notNull bool
		// filterIsNull contains a list of tables that we should filter for null on
		filterIsNullArg *ast.Argument
		// The `order_by` GraphQL argument can only be processed after all table aliases
		// are known. That includes the tables referenced by the GraphQL subfields.
		// Therefore, upon encountering the `order_by` argument in one of the root
		// types, it is stored in this variable for processing after all the table
		// aliases are known, that is after the rest of the GraphQL query had been processed.
		orderByArg *ast.Argument

		// The `first` arg is a limit on the results in ASC order
		firstArg *ast.Argument
		// The `last` arg is a limit on the results in DESC order
		lastArg *ast.Argument
	)

	// Always return the ID field of a table as the first row as we need it when
	// we aggregate the results up into the returned value
	tc.columns = append(tc.columns, tableIDField)
	nodeQuery = nodeQuery.Column(tc.alias + "." + tableIDField)
	*sql = sql.Column(tc.alias + "." + tableIDField)

	// GraphQL arguments are processed here
	for _, arg := range tc.field.Arguments {

		// The arguments can be of different kinds, from simple column names
		// to our custom function names. We'll work through the possibilities,
		// but if at the end the argument has not been resolved, raise an error.
		argIsResolved := false

		// Argument name equal to one of the column names for the current node (table)
		// adds an equality predicate in the WHERE clause.
		// Multiple expressions are `AND`ed together in the generated SQL.
		for _, tf := range node.Table.Fields {
			if arg.Name.Value == tf.Name {
				nodeQuery = nodeQuery.Where(sq.Eq{tc.alias + "." + arg.Name.Value: arg.Value.GetValue()})
				argIsResolved = true
				break
			}
		}
		// Resolve the id field
		if arg.Name.Value == tableIDField {
			nodeQuery = nodeQuery.Where(sq.Eq{tc.alias + "." + arg.Name.Value: arg.Value.GetValue()})
			argIsResolved = true
		}

		if argIsResolved {
			continue
		}

		// Process the arguments that are not GraphQL/DB field/column names...
		switch arg.Name.Value {
		case notNullID:
			// The notNullID argument is used on sub-fields and should be
			// processed by the parent. E.g.
			// table_a {
			// 	table_b(not_null: true) {...}
			// }
			notNull = true
			argIsResolved = true
		case filterIsNullID:
			filterIsNullArg = arg
			argIsResolved = true
		case orderByID:
			// The order_by argument is allowed only at the top level. Futhermore, it cannot be processed until
			// all subfields had been processed, because at the top level the alias names of tables are not known.
			// Therefore, defer the processing of this argument by saving a pointer to it for later processing.
			orderByArg = arg
			argIsResolved = true
		case firstID:
			firstArg = arg
			argIsResolved = true
		case lastID:
			lastArg = arg
			argIsResolved = true
		}

		if firstArg != nil && lastArg != nil {
			return fmt.Errorf("cannot provide both 'first' and 'last' arguments for table %s", tc.table)
		}

		// The argument name which is not a column name is a mistake, raise error.
		if !argIsResolved {
			return fmt.Errorf("unknown argument identifier for table %s: %s", tc.table, arg.Name.Value)
		}
	}

	// Iterate over the fields in the selection set (if any) for the current `field`
	for _, selection := range tc.field.SelectionSet.Selections {
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

		// Skip the _id field as we will add it ourselves
		if fieldName == tableIDField {
			continue
		}

		// A non-nil selection set implies that the subField refers to another
		// table in our schema.
		// We need to process the columns/fields for this table first, before we
		// process any subFields, so simply append these to a slice and process
		// at the end of the function
		if subField.SelectionSet != nil {
			subFields = append(subFields, subField)
		} else {
			// If subField did not have a selection set this it is just a column
			// within the current table, so add it to the columns.
			// There are exceptions to adding the column to the SQL query,
			// e.g. time fields should be to_json(time)
			var colVal string
			for _, tf := range node.Table.Fields {
				if tf.Name == fieldName {
					switch tf.Type {
					case parser.TimeType:
						// colVal = tableColumn(tc.alias, fieldName)
						colVal = "to_json(" + tableColumn(tc.alias, fieldName) + ") as " + fieldName
					default:
						colVal = tableColumn(tc.alias, fieldName)
					}
				}
			}
			tc.columns = append(tc.columns, fieldName)
			nodeQuery = nodeQuery.Column(colVal)
			*sql = sql.Column(colVal)
		}
	}

	// Once we have processed this fields columns, proceed to the subFields.
	// To ensure the correct ordering of both columns in the nodeQuery and also
	// the order of JOINs in the root SQL query, we need to be a bit careful.
	// nodeQuery must be added to the root SQL query BEFORE we recurse and process
	// the children/subFields, but we also need to know if subFields have a
	// dependency on a foreign key in nodeQuery
	var subColumns []*tableColumns
	for _, subField := range subFields {
		var (
			fieldName         = subField.Name.Value
			edgeToRelatedNode *SchemaEdge
		)

		edgeToRelatedNode, err := node.Edge(fieldName)
		if err != nil {
			return fmt.Errorf("no relationship found between tables: '%s', '%s'", node.Table.Name, fieldName)
		}

		// Recursively resolve for the subField `B`, which may contain further nested fields.
		subCol := &tableColumns{
			table:  fieldName,
			alias:  tableAlias(fieldName, depth),
			field:  subField,
			scalar: edgeToRelatedNode.isScalar(),
		}
		tc.children = append(tc.children, subCol)
		subColumns = append(subColumns, subCol)

		// There is a strange circumstance where we have something like:
		// a {
		// 	b {
		// 	  a { ...}
		// 	}
		// }
		// In this scenario, we need ot prevent the join column from being added
		// twice to either a or b
		if parent != nil && fieldName == parent.table {
			continue
		}
		var (
			rhsJoinOn      string
			leftTableAlias = tc.alias
			rightTable     = edgeToRelatedNode.Node.Table.Name
		)

		if edgeToRelatedNode.Rel == BelongsTo {
			rhsJoinOn = tableColumn(leftTableAlias, foreignKeyField(rightTable))
			// Make sure we select the foreign key field from the parent
			nodeQuery = nodeQuery.Column(rhsJoinOn)
		}
	}

	//
	// Is Null
	//
	// If the filter_is_null argument is used, we need to filter on null for the
	// tables provided
	if filterIsNullArg != nil {
		var filterIsNull []string
		switch {
		case filterIsNullArg.Value.GetKind() == kinds.EnumValue:
			filterIsNull = append(filterIsNull, filterIsNullArg.Value.GetValue().(string))
		case filterIsNullArg.Value.GetKind() == kinds.ListValue:
			for _, arg := range filterIsNullArg.Value.GetValue().([]ast.Value) {
				filterIsNull = append(filterIsNull, arg.GetValue().(string))
			}
		default:
			return fmt.Errorf("unknown argument kind for argument %s: %s", filterIsNullID, filterIsNullArg.Value.GetKind())
		}
		// Goal: create a subquery that will filter the current node/table when
		// the provided argument table is null.
		// Nicest way was using a WHERE NOT EXISTS query, e.g.
		// WHERE NOT EXISTS( SELECT 1 FROM abc WHERE abc.id = node.abc_id)
		for _, nullTable := range filterIsNull {
			e, err := node.Edge(nullTable)
			if err != nil {
				return fmt.Errorf("no edge found for filter_is_null argument on table %s: %s", node.Table.Name, nullTable)
			}
			subQueryFilter := sq.Select("1").From(psqlAbsTableName(tenant, nullTable))
			// Depending on the edge direction (relationship) we know if the left table
			// or right table has the foreign key
			if e.Rel == BelongsTo {
				subQueryFilter = subQueryFilter.Where(tableColumn(nullTable, tableIDField) + " = " + tableColumn(tc.alias, nullTable+tableJoinSuffix))
			} else {
				subQueryFilter = subQueryFilter.Where(tableColumn(nullTable, tc.table+tableJoinSuffix) + " = " + tableColumn(tc.alias, tableIDField))
			}
			sqStr, sqArgs := subQueryFilter.MustSql()
			nodeQuery = nodeQuery.Where("NOT EXISTS( "+sqStr+" )", sqArgs...)
		}
	}

	//
	// Order
	//
	// By default we want to preserve the "natural" order, unless an order_by
	// is specified
	//
	if orderByArg != nil {
		orderByFields, ok := orderByArg.Value.GetValue().([]*ast.ObjectField)
		if !ok {
			return fmt.Errorf("invalid format for 'order_by' argument")
		}
		for _, orderBy := range orderByFields {
			var (
				field = orderBy.Name.Value
				order = strings.ToUpper(orderBy.Value.GetValue().(string))
			)
			if !(order == orderAsc || order == orderDesc) {
				return fmt.Errorf("unknown order for 'order_by': %s", order)
			}
			// Add the ORDER BY to both the nodeQuery and the root SQL query
			nodeQuery = nodeQuery.OrderBy(tableColumn(tc.alias, field) + " " + order)
			*sql = sql.OrderBy(tableColumn(tc.alias, field) + " " + order)
		}
	}

	//
	// Limit
	//
	// Limiting is handled by the `first` and `last` values. This is only
	// added to the subquery for this node, and needs to come after the ordering
	// of other fields to make sure we respect the wishes of the user and then
	// get first/last based on the given order
	//
	if firstArg != nil {
		limitStr, ok := firstArg.Value.GetValue().(string)
		if !ok {
			return fmt.Errorf("could not convert the value of the argument `first`: %#v", firstArg.Value.GetValue())
		}
		n, err := strconv.ParseUint(limitStr, 10, 64)
		if err != nil {
			return fmt.Errorf("could not convert the value to unsigned integer: %s", limitStr)
		}
		// Order by ASC and then limit
		nodeQuery = nodeQuery.
			OrderBy(tableColumn(tc.alias, tableIDField) + " " + orderAsc).
			Limit(n)
	}
	if lastArg != nil {
		limitStr, ok := lastArg.Value.GetValue().(string)
		if !ok {
			return fmt.Errorf("could not convert the value of the argument `last`: %#v", lastArg.Value.GetValue())
		}
		n, err := strconv.ParseUint(limitStr, 10, 64)
		if err != nil {
			return fmt.Errorf("could not convert the value to unsigned integer: %s", limitStr)
		}
		// Order by DESC and then limit
		nodeQuery = nodeQuery.
			OrderBy(tableColumn(tc.alias, tableIDField) + " " + orderDesc).
			Limit(n)
	}
	// Add default orderBy and limit if there isn't one already
	if lastArg == nil && firstArg == nil {
		if orderByArg == nil {
			nodeQuery = nodeQuery.OrderBy(tableColumn(tc.alias, tableIDField) + " " + orderDesc)
		}
		nodeQuery = nodeQuery.Limit(defaultLimit)
	}

	// Before processing any subFields (which are like "children" in GraphQL),
	// we need to add nodeQuery to the rootSQL query.
	// If parent is nil then we are at the root table and this should be the FROM
	// part of the SQL statement.
	// Else, things are a bit more involved as the nodeQuery is part of a JOIN
	if parent == nil {
		*sql = sql.FromSelect(nodeQuery, tc.alias)
	} else {

		edgeToParent, err := node.Edge(parent.table)
		if err != nil {
			return err
		}

		//
		// SQL JOIN
		//
		// The next part we need to add to the query is the SQL JOIN.
		// This is rather complex, as the current node SQL query needs to filter
		// on the parent, and we don't know what the parent _id field is until
		// we resolve it, which happens in it's own subquery.
		// Fortunately, Postgres has a killer feature called *LATERAL* subqueries.
		// Essentially, they let you refer to table columns on the left hand side
		// of the query.
		// This means we can resolve parent (on the left hand side) and then
		// use the value of the parent table to resolve this node's subquery
		// with a WHERE on the foreign key... Pretty sweet! (but complex)
		//
		var (
			joinStr         string
			lhsJoinOn       string
			rhsJoinOn       string
			leftTable       = parent.table
			leftTableAlias  = parent.alias
			rightTable      = tc.table
			rightTableAlias = tc.alias
		)
		switch edgeToParent.Rel {
		case BelongsTo:
			lhsJoinOn = tableColumn(leftTableAlias, tableIDField)
			rhsJoinOn = tableColumn(rightTableAlias, foreignKeyField(leftTable))
			// Make sure we select the foreign key field from the nodeQuery
			nodeQuery = nodeQuery.Column(rhsJoinOn)
		case OneToOne, OneToMany:
			lhsJoinOn = tableColumn(rightTableAlias, tableIDField)
			rhsJoinOn = tableColumn(leftTableAlias, foreignKeyField(rightTable))
		}

		// Add the WHERE condition for this subquery
		nodeQuery = nodeQuery.Where(lhsJoinOn + " = " + rhsJoinOn)
		// Generate the SQL query for this node
		sqlStr, sqlArgs, err := nodeQuery.ToSql()
		if err != nil {
			return fmt.Errorf("error creating SQL query for node %s: %w", node.Table.Name, err)
		}
		sqlStr = " ( " + sqlStr + " ) AS " + rightTableAlias
		joinStr = "LATERAL " + sqlStr + " ON true"

		if notNull {
			*sql = sql.InnerJoin(joinStr, sqlArgs...)
		} else {
			*sql = sql.LeftJoin(joinStr, sqlArgs...)
		}
	}

	// Create and add sub queries for the children to the root SQL query
	for _, subCol := range subColumns {
		err := psqlSubQuery(tenant, graph, sql, tc, subCol, depth+1)
		if err != nil {
			return err
		}
	}

	// After we have processed the sub fields, if there was not orderBy given
	// for this field then add a default one to "preserve" the natural order.
	// IMPORTANT: this has to come AFTER we handle sub fields, so that we honour
	// the requests made by sub children
	if orderByArg == nil {
		*sql = sql.OrderBy(tableColumn(tc.alias, tableIDField) + " " + orderAsc)
	}
	return nil
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
		var isNilTable = true
		tColVal = make(map[string]interface{}, len(tc.columns))
		for _, field := range tc.columns {
			val := scanValues[*index]
			tColVal[field] = val
			*index++
			// If isNilTable is false, or val is not nil, then set isNilTable to
			// false indicating that there was a field with a value.
			// If all field values were null, then we don't want to add this value
			// to the result
			isNilTable = isNilTable && val == nil
		}

		// If all values were nil and there are no children (hence no need for
		// any values at all), then set the tColVal to nil as it's just a map
		// of nil values with no children
		if isNilTable && len(tc.children) == 0 {
			tColVal = nil
		}
		// Check if we expect the result to be a scalar value or a list.
		// It is scalar depending on the relationship between the tables in the
		// schema.
		// If scalar, simply assign the value
		if tc.scalar {
			tVal = tColVal
		}
		// If not scalar, we should create a list. Then it depends on whether
		// the value is nil or not, and if there are children
		if !tc.scalar {
			tVal = make([]map[string]interface{}, 0)
			// Only append tColVal if it's not nil
			if tColVal != nil {
				tVal = append(tVal.([]map[string]interface{}), tColVal)
			}
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
		*index += len(tc.columns)
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
			var isNilTable = true
			tColVal = make(map[string]interface{})
			for _, field := range tc.columns {
				val := scanValues[*index]
				tColVal[field] = val
				*index++
				// If isNilTable is false, or val is not nil, then set isNilTable to
				// false indicating that there was a field with a value.
				// If all field values were null, then we don't want to add this value
				// to the result
				isNilTable = isNilTable && val == nil
			}
			// If all values were nil we want to ignore this entry.
			// If not, then append to the list
			if !isNilTable {
				tListVal = append(tListVal, tColVal)
			}

		} else {
			// Make sure we increment the index
			*index += len(tc.columns)
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
