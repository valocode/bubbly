package interim

import (
	"github.com/graphql-go/graphql"
	"github.com/hashicorp/go-memdb"
	"github.com/zclconf/go-cty/cty"
)

// newGraphQL creates a new GraphQL schema
// wrapping the given memDB with a schmea that
// corresponds to the given set of tables.
func newGraphQL(tables []Table, memDB *memdb.MemDB) (graphql.Schema, error) {
	// These are the top-level query fields. Each of these fields
	// will correspond to each of the tables in the entire hierarchy.
	queryFields := make(graphql.Fields)

	// Recuresively walk the table hierarchy, appending each table
	// to qf and also creating a relationship between the parent table
	// and all its subtables, if they exist.
	for _, t := range tables {
		// These are top-level tables so we can ignore their graphQL types.
		addTableToGraphQL(t, memDB, queryFields)
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
func addTableToGraphQL(t Table, memDB *memdb.MemDB, queryFields graphql.Fields) graphql.Type {
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
		t := graphQLFieldType(f)
		fields[f.Name] = &graphql.Field{Type: t}
		args[f.Name] = &graphql.ArgumentConfig{Type: t}
	}

	// Each sub table represents a distinct GraphQL type.
	// In order for the current type to know the GraphQL type
	// we need to resolve the types of all the subtables before
	// we can continue.
	for _, sub := range t.Tables {
		// Recursively add the subtable first so we can get its type.
		subType := addTableToGraphQL(sub, memDB, queryFields)

		// Each sub type is a list on the current type.
		fields[sub.Name] = &graphql.Field{
			Type: graphql.NewList(subType),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return resolveList(memDB, t, p)
			},
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
		Type: tType,
		Args: args,
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			return resolveScalar(memDB, t, p)
		},
	}

	return tType
}

func resolveScalar(memDB *memdb.MemDB, t Table, p graphql.ResolveParams) (interface{}, error) {
	txn := memDB.Txn(false)
	defer txn.Abort()

	var (
		n   interface{}
		err error
	)
	for k, v := range p.Args {
		n, err = txn.First(t.Name, k, v)
		break
	}

	return n, err
}

func resolveList(memDB *memdb.MemDB, t Table, p graphql.ResolveParams) (interface{}, error) {
	txn := memDB.Txn(false)
	defer txn.Abort()

	var (
		n   interface{}
		err error
	)
	for k, v := range p.Args {
		n, err = txn.First(t.Name, k, v)
		break
	}

	return n, err
}

func graphQLFieldType(f Field) *graphql.Scalar {
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
