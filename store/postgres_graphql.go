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

func newSQLQueryBuilder() *sqlQueryBuilder {
	return &sqlQueryBuilder{
		sql:     psql.Select(),
		columns: make([]*tableColumns, 0),
		depth:   0,
	}
}

type sqlQueryBuilder struct {
	sql     sq.SelectBuilder
	columns rowColumns
	node    *schemaNode
	depth   int
}

// rowColumns stores a slice of table columns, which forms a complete SQL row
// of SELECT'd values
type rowColumns []*tableColumns

func (r rowColumns) length() int {
	var count = 0
	for _, t := range r {
		count += len(t.fields)
	}
	return count
}

// tableColumns is used to store the columns that are SELECT'd in a SQl
// statement, within one single table.
// This is used so that the "structure" of a returned SQL row is maintained
type tableColumns struct {
	table  string
	alias  string
	fields []string
	scalar bool
}

// psqlResolveRootQueries gets called for each top-level query and iterates
// through the fields in that root query and resolves them.
// What is a bit puzzling is that if you have a query with two fields, this
// method gets called twice, once for each field, but each time the
// graphql.ResolveParams contains a list of FieldASTs with one element
func psqlResolveRootQueries(pool *pgxpool.Pool, graph *schemaGraph, params graphql.ResolveParams) (interface{}, error) {
	var (
		result interface{}
		err    error
	)
	for _, field := range params.Info.FieldASTs {
		result, err = psqlResolveRootQuery(pool, graph, field)
		if err != nil {
			return nil, fmt.Errorf("failed to resolve query: %s: %w", field.Name.Value, err)
		}
	}
	return result, err
}

// psqlResolveRootQuery resolves a single root graphql query
func psqlResolveRootQuery(pool *pgxpool.Pool, graph *schemaGraph, field *ast.Field) (interface{}, error) {
	var (
		result    = make(map[string]interface{})
		qb        = newSQLQueryBuilder()
		rootTable = field.Name.Value
		rootAlias = field.Name.Value + "_" + strconv.Itoa(qb.depth)
	)

	// Set the starting node and initialize the sql statement
	qb.node = graph.nodeIndex[rootTable]
	qb.sql = qb.sql.From(tableAsAlias(rootTable, rootAlias))

	// Recursively go through the graphql query and resolve the sub-fields
	if err := psqlSubQuery(graph, qb, field, nil); err != nil {
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
		if err := psqlScanRowColumns(rows, result, qb.columns); err != nil {
			return nil, fmt.Errorf("failed scanning row values: %w", err)
		}
	}
	return result[rootTable], nil
}

func psqlSubQuery(graph *schemaGraph, qb *sqlQueryBuilder, field *ast.Field, path schemaPath) error {
	// Create the tableColumns for this type/table in the query
	var (
		node = qb.node
		tb   = tableColumns{
			table:  node.table.Name,
			alias:  tableAlias(node.table.Name, qb.depth),
			scalar: isTableResultScalar(node, path),
		}
	)
	// Increment the depth
	qb.depth++
	// Add the tableColumns to the list of table columns
	qb.columns = append(qb.columns, &tb)

	// Always return the ID field of a table as the first row as we need it when
	// we aggregate the results up into the returned value
	tb.fields = append(tb.fields, tableIDField)
	qb.sql = qb.sql.Column(tb.alias + "." + tableIDField)

	// TODO: this currently does not support all the graphql arguments, but
	// more work will be done on this to enable things like calculating the
	// cartesian product, so no need to solve it all right now
	for _, arg := range field.Arguments {
		qb.sql = qb.sql.Where(sq.Eq{tb.alias + "." + arg.Name.Value: arg.Value.GetValue()})
	}

	// Iterate over the selections/fields inside this field
	for _, selection := range field.SelectionSet.Selections {
		// Make sure we don't miss any fields. Only ast.Field types should be
		// supported right now.
		subField, ok := selection.(*ast.Field)
		if !ok {
			return fmt.Errorf("graphql query selection type not supported: %s", selection.GetSelectionSet().Kind)
		}
		var fieldName = subField.Name.Value
		// We do not care about "meta fields" or whatever they are called in
		// graphql that might get added by clients to help with caching and such
		// and they typically start with __
		if strings.HasPrefix(fieldName, "__") {
			continue
		}
		// If subField has a selection set then it is not a scalar type, and is
		// another table in our schema.
		// TODO: the logic here is a little tricky and should be cleaned up once
		// we understand more about the different types of queries we want to make
		// and this is coming in future versions so no need to spend lots of time
		// now
		if subField.SelectionSet != nil {
			// If the current node has a oath to the child type in the graphql
			// query
			path := qb.node.shortestPath(fieldName)
			if path == nil {
				return fmt.Errorf("no path was found between tables: %s --> %s", qb.node.table.Name, fieldName)
			}
			var (
				leftTable      = tb.table
				leftTableAlias = tb.alias
			)
			for _, e := range path {
				var (
					rightTable      = e.node.table.Name
					rightTableAlias = tableAlias(rightTable, qb.depth)
				)
				switch e.rel {
				case oneToOne, oneToMany:
					qb.sql = qb.sql.LeftJoin(joinOn(
						tableAsAlias(rightTable, rightTableAlias),
						tableColumn(leftTableAlias, tableIDField),
						tableColumn(rightTableAlias, foreignKeyField(leftTable))))
				case belongsTo:
					qb.sql = qb.sql.LeftJoin(
						joinOn(
							tableAsAlias(rightTable, rightTableAlias),
							tableColumn(rightTableAlias, tableIDField),
							tableColumn(leftTableAlias, foreignKeyField(rightTable))))
				}

				leftTable = e.node.table.Name
				leftTableAlias = tableAlias(e.node.table.Name, qb.depth)
			}

			// Recursively resolve
			qb.node = graph.nodeIndex[fieldName]
			if err := psqlSubQuery(graph, qb, subField, path); err != nil {
				return err
			}
			continue
		}

		// If subField did not have a selection set this it is just a column
		// within the current table, so add it to the columns
		tb.fields = append(tb.fields, fieldName)
		qb.sql = qb.sql.Column(tableColumn(tb.alias, fieldName))
	}

	return nil
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
// an existing result map that the row should be aggregated into, and a rowColumns
// type which contains the columns in the row, grouped by their tables.
// The rowColumns are used to unpack the row values into the result
func psqlScanRowColumns(row pgx.Row, result map[string]interface{}, columns rowColumns) error {
	var (
		columnLen     = columns.length()
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

	var (
		parentVal = result
		index     = 0
	)
	// For each table in columns, get the ID value from the scanned values and
	// check whether a record with that ID already exists in the result values
	for _, tc := range columns {
		// Check if the value for this table of columns already exists.
		// If not, initialize it
		tVal, ok := parentVal[tc.table]
		if !ok {
			// Initialize the value from the scanned results for the group of
			// columns in this table
			var tColVal = make(map[string]interface{})
			for _, field := range tc.fields {
				tColVal[field] = scanValues[index]
				index++
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
			parentVal = tColVal

			// Skip the rest and continue to the next table
			continue
		}
		// If the parentVal already contained a value for this table.
		// Get the ID value from this result row
		var tableIDVal = scanValues[index]

		if tc.scalar {
			// If the table value should be scalar and the value is already
			// initialized, then we do not need to do anything.
			// Set the new parentVal and continue to the next table
			parentVal = tVal.(map[string]interface{})
			index += len(tc.fields)
			continue
		}
		var tColVal map[string]interface{}
		var tListVal = tVal.([]map[string]interface{})
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
				tColVal[field] = scanValues[index]
				index++
			}
			tListVal = append(tListVal, tColVal)

		} else {
			// Make sure we increment the index
			index += len(tc.fields)
		}
		parentVal[tc.table] = tListVal
		parentVal = tColVal
	}

	return nil
}

func isTableResultScalar(node *schemaNode, path schemaPath) bool {
	// Check if nil, as this will happen when processing the root query, in
	// which case we want to return a list, not scalar, so return false
	if path == nil {
		return false
	}
	// Else, delegate to the path to figure out
	return path.isScalar()
}
