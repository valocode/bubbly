package store

import (
	"github.com/graphql-go/graphql"
	"github.com/verifa/bubbly/api/core"
	"github.com/zclconf/go-cty/cty"
)

// newGraphQLSchema creates a new GraphQL schema
// wrapping the given memDB with a schmea that
// corresponds to the given set of tables.
func newGraphQLSchema(tables []core.Table, p provider) (graphql.Schema, error) {
	// These are the top-level query fields. Each of these fields
	// will correspond to each of the tables in the entire hierarchy.
	queryFields := make(graphql.Fields)

	// Recuresively walk the table hierarchy, appending each table
	// to qf and also creating a relationship between the parent table
	// and all its subtables, if they exist.
	for _, t := range tables {
		// These are top-level tables so we can ignore their graphQL types.
		addTableToGraphQL(t, p, queryFields)
	}

	// This config is used to create a new query type
	// that will be used to create the GraphQL schema.
	// Note that this config only contains a query, and
	// no corresponding mutation since this data is readonly.
	cfg := graphql.SchemaConfig{
		Query: graphql.NewObject(
			graphql.ObjectConfig{
				Name:   "query",
				Fields: queryFields,
			},
		),
	}

	return graphql.NewSchema(cfg)
}

// addTableToGraphQL adds as table to queryFields and returns the GraphQL type
// corresponding to the table so it can be included in parent tables (if they exist).
func addTableToGraphQL(t core.Table, p provider, queryFields graphql.Fields) (graphql.Type, graphql.FieldConfigArgument) {
	var (
		// These are the fields for this specific table
		// which will correspond to fields on the GraphQL
		// type, created dynamically below.
		fields = make(graphql.Fields)
		// These are args for this specific table as well.
		// These are args that can be used for querying on
		// this table/type.
		args = make(graphql.FieldConfigArgument)
	)

	// Set fields and args for the current table/type.
	for _, f := range t.Fields {
		ft := graphQLFieldType(f)
		fields[f.Name] = &graphql.Field{Type: ft}
		args[f.Name] = &graphql.ArgumentConfig{Type: ft}
	}

	args[filterName] = &graphql.ArgumentConfig{
		Type: graphQLFilterType(t.Name, args),
	}

	// Each sub table represents a distinct GraphQL type.
	// In order for the current type to know the GraphQL type
	// we need to resolve the types of all the subtables before
	// we can continue.
	for _, sub := range t.Tables {
		// Recursively add the subtable first so we can get its type.
		subType, subArgs := addTableToGraphQL(sub, p, queryFields)

		// Each sub type is a list on the current type.
		fields[sub.Name] = &graphql.Field{
			Type:    graphql.NewList(subType),
			Args:    subArgs,
			Resolve: p.ResolveList,
		}
	}

	// Create a GraphQL type for the current table so that we
	// can set it in the query fields and return it to be used
	// by the parent table (if there is one).
	tType := graphql.NewObject(
		graphql.ObjectConfig{
			Name:   t.Name,
			Fields: fields,
		},
	)

	// Register a type for our table corresponding to the
	// name of the table itself (note: this name is not namespaced).
	queryFields[t.Name] = &graphql.Field{
		Type:    tType,
		Args:    args,
		Resolve: p.ResolveScalar,
	}

	return tType, args
}

func graphQLFieldType(f core.TableField) *graphql.Scalar {
	switch f.Type {
	case cty.Bool:
		return graphql.Boolean
	case cty.Number:
		return graphql.Int
	case cty.String:
		return graphql.String
	}

	return nil
}

const (
	filterName = "filter"
)

const (
	filterGreaterThan          = "_gt"
	filterLessThan             = "_lt"
	filterGreaterThanOrEqualTo = "_gte"
	filterLessThanOrEqualTo    = "_lte"
	filterIn                   = "_in"
	filterNotIn                = "_not_in"
)

var scalarFilters = []string{
	filterGreaterThan,
	filterLessThan,
	filterGreaterThanOrEqualTo,
	filterLessThanOrEqualTo,
}

var listFilters = []string{
	filterIn,
	filterNotIn,
}

func graphQLFilterType(typeName string, args graphql.FieldConfigArgument) *graphql.InputObject {
	var (
		// Micro-opt: we know the size of the field map is the total number
		// of filter ops times the number of args we are given.
		numFields = (len(scalarFilters) + len(listFilters)) * len(args)
		fields    = make(graphql.InputObjectConfigFieldMap, numFields)
	)
	for n, a := range args {
		for _, f := range scalarFilters {
			fields[n+f] = &graphql.InputObjectFieldConfig{
				Type: a.Type,
			}
		}
		for _, f := range listFilters {
			fields[n+f] = &graphql.InputObjectFieldConfig{
				Type: graphql.NewList(a.Type),
			}
		}
	}

	return graphql.NewInputObject(
		graphql.InputObjectConfig{
			Name:   typeName + "_filter",
			Fields: fields,
		},
	)
}

func isValidParent(parent interface{}) bool {
	if parent == nil {
		return false
	}

	if _, ok := parent.(map[string]interface{}); ok {
		// graphql-go passes nil maps as empty values
		// and we don't work with maps so a map is not
		// a valid parent.
		return false
	}

	return true
}
