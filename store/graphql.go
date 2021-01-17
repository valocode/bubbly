package store

import (
	"fmt"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/language/ast"
	"github.com/verifa/bubbly/api/core"
	"github.com/zclconf/go-cty/cty"
)

var mapScalar = graphql.NewScalar(graphql.ScalarConfig{
	Name:        "Map",
	Description: "The `Map` scalar type represents a Map for storing key/value pairs",
	Serialize: func(value interface{}) interface{} {
		return value
	},
	ParseValue: func(value interface{}) interface{} {
		return value
	},
	ParseLiteral: func(value ast.Value) interface{} {
		// TODO: not sure exactly what to do here...
		fmt.Printf("Map ParseLiteral: %v\n", value)
		return nil
	},
})

// gqlField is our custom Graphql Field type so that we can store a field in
// it's simplest form and iteratively add to it, before we convert it into a
// real graphql field.
// One challenge is that we need to "reuse" fields inside joins, and in some
// cases the end field type might be a List or a Scalar, so it is just easier
// to encapsulate that inside this struct rather than add lots of complexity.
type gqlField struct {
	Type *graphql.Object
	Args graphql.FieldConfigArgument
	// Store the type of a field as an InputObject which can be used to add
	// this type as an argument based on filters
	FilterType *graphql.InputObject
	// Resolve graphql.FieldResolveFn
}

// newGraphQLSchema creates a new GraphQL schema wrapping the given provider
// with a schema that corresponds to the given set of tables.
func newGraphQLSchema(schema *bubblySchema, p provider) (graphql.Schema, error) {
	var (
		fields = make(map[string]gqlField)
		// These are the top-level query fields. Each of these fields
		// will correspond to each of the tables in the entire hierarchy.
		queryFields = make(graphql.Fields)
	)

	// Iterate over the schema of tables, appending each table
	// to queryFields.
	for _, t := range schema.Tables {
		addGraphFields(t, fields)
	}

	// Iterate over the schema of tables again, this time adding the joins
	// to the graphql schema.
	for _, t := range schema.Tables {
		addGraphJoins(t, fields, p)
	}

	// Finally, we want to populate the queryFields using the graphql types
	// we have created
	for _, field := range fields {
		queryFields[field.Type.Name()] = &graphql.Field{
			Type:    graphql.NewList(field.Type),
			Args:    field.Args,
			Resolve: p.ResolveList,
		}
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

// addGraphFields takes all the tables in the schema and creates our custom
// graphql fields which we use for later processing.
func addGraphFields(t core.Table, fields map[string]gqlField) {
	// These are the fields for this specific table
	// which will correspond to fields on the GraphQL
	// type, created dynamically below.
	var (
		// typeFields are the fields that will be nested inside this type that
		// we are creating.
		typeFields = make(graphql.Fields)
		// gqlField is the graphql field which we are populating now
		gqlField = fields[t.Name]
	)
	// Initialize the args
	gqlField.Args = make(graphql.FieldConfigArgument)

	// Set fields and args for the current table/field
	for _, f := range t.Fields {
		ft := graphQLFieldType(f)
		typeFields[f.Name] = &graphql.Field{Type: ft}
		gqlField.Args[f.Name] = &graphql.ArgumentConfig{Type: ft}
	}

	// Create the FilterType for this graphql field
	gqlField.FilterType = graphQLFilterType(t.Name, gqlField.Args)

	gqlField.Args[filterID] = &graphql.ArgumentConfig{
		Type: gqlField.FilterType,
	}
	gqlField.Args[firstID] = &graphql.ArgumentConfig{
		Type: graphql.Int,
	}

	// Create a GraphQL type for the current table so that we
	// can set it in the query fields and return it to be used
	// by the parent table (if there is one).
	gqlField.Type = graphql.NewObject(
		graphql.ObjectConfig{
			Name:   t.Name,
			Fields: typeFields,
		},
	)

	// Assign the gqlField back to the map
	fields[t.Name] = gqlField
}

func addGraphJoins(t core.Table, fields map[string]gqlField, p provider) {
	// Get the type corresponding to the core.Table
	var field = fields[t.Name]
	for _, join := range t.Joins {

		var (
			// Get the field of the table that the join refers to
			joinField = fields[join.Table]
			// Set the type of the nested type to "has many" be default, by
			// making it a graphql List
			fieldType    graphql.Output = graphql.NewList(field.Type)
			fieldResolve                = p.ResolveList
		)

		// populateGraphQLTypeArgs(joinArgs, joinType)
		// Add the join "belongs to"
		field.Type.AddFieldConfig(join.Table, &graphql.Field{
			Type:    joinField.Type,
			Args:    joinField.Args,
			Resolve: p.ResolveScalar,
		})
		// Add the relationship field so that we know there is a join here
		field.Type.AddFieldConfig(join.Table+"_id", &graphql.Field{
			Type: graphql.Int,
		})

		// If the join is unique, then it only has one, and fieldType should
		// not return a list but a single instance
		if join.Unique {
			fieldType = field.Type
			fieldResolve = p.ResolveScalar
		}
		// populateGraphQLTypeArgs(fieldArgs, tType)
		joinField.Type.AddFieldConfig(t.Name, &graphql.Field{
			Type:    fieldType,
			Args:    field.Args,
			Resolve: fieldResolve,
		})
		fields[join.Table] = joinField
	}
	fields[t.Name] = field
}

func graphQLFieldType(f core.TableField) *graphql.Scalar {
	switch ty := f.Type; {
	case ty == cty.Bool:
		return graphql.Boolean
	case ty == cty.Number:
		return graphql.Int
	case ty == cty.String:
		return graphql.String
	case ty.IsObjectType():
		return mapScalar
	default:
		panic(fmt.Sprintf("Unsupported GraphQL conversion from cty.Type: %s", f.Type.GoString()))
	}
}

const (
	filterID = "filter"
	firstID  = "first"
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

	if val, ok := parent.(map[string]interface{}); ok {
		if len(val) == 0 {
			// graphql-go passes nil maps as empty values
			// and we don't work with maps so a map is not
			// a valid parent.
			return false
		}
	}

	return true
}

// childBelongsToParent takes a parent and a child table name, and checks if
// the child type in the graphql schema "belongs to" (as a relationship) the
// parent type.
// E.g. if we have a type repo and repo_version, where repo_version stores a
// field "repo_id" which is a foreign key to repo, then repo_version belongs to
// repo.
// Hence if given the following graphl query then this function returns true
// because the child (repo_version) belongs to the parent (repo)
// {
// 	 repo {
// 		repo_version {
// 			...
// 		}
// 	 }
// }
func childBelongsToParent(info graphql.ResolveInfo, parent string, child string) bool {
	var (
		queryFields = info.Schema.QueryType().Fields()
		childField  = queryFields[child]
		parentFK    = parent + tableJoinSuffix
	)
	// All top-level fields in the query type are graphql.Lists of
	// of type graphql.Object, which has it's type fields
	obj := childField.Type.(*graphql.List).OfType.(*graphql.Object)
	// Check if the parent foreign key field exists in the child.
	// Consider if "b" belongs to "a". Then "b" will have a field "a_id" which
	// is the foreigh key.
	_, ok := obj.Fields()[parentFK]
	if !ok {
		return false
	}

	return true
}

// queryFieldSelectionSet returns the list of fields in a query that should be
// returned (because they are listed in the graphql query).
// This also adds some fields that might not be listed, like the _id
func queryFieldSelectionSet(params graphql.ResolveParams) []string {
	var (
		fieldNames      = make([]string, 0)
		addTableIDField = true
	)
	for _, parent := range params.Info.FieldASTs {
		for _, selection := range parent.SelectionSet.Selections {
			if child, ok := selection.(*ast.Field); ok {
				if child.SelectionSet != nil {
					// If the childField selection set is not nil, it means that
					// the childField is another object, which is another
					// core.Table and therefore should be resolved separately.
					// The relationship between parent and child can be either:
					// 1. child belongs to parent, do nothing
					// 2. parent belongs to child, return the parent's foreign
					// key to the child
					if !childBelongsToParent(params.Info, parent.Name.Value, child.Name.Value) {
						fieldNames = append(fieldNames, child.Name.Value+tableJoinSuffix)
					}
					continue
				}
				if child.Name.Value == tableIDField {
					addTableIDField = false
				}
				fieldNames = append(fieldNames, child.Name.Value)
			}
		}
	}
	// We likely want the ID field in case it needs to be referenced in the
	// graphql query
	if addTableIDField {
		fieldNames = append(fieldNames, tableIDField)
	}
	return fieldNames
}
