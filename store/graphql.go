package store

import (
	"fmt"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/language/ast"
	"github.com/graphql-go/graphql/language/kinds"
	"github.com/valocode/bubbly/api/core"
	"github.com/zclconf/go-cty/cty"
)

//
// The GraphQL Schema is a representation of the Bubbly Schema Graph,
// enabling GraphQL access to Bubbly.
//

// graphJoinDistance ???
const graphJoinDistance = 3

// gqlField is our custom Graphql Field type so that we can store a field in
// it's simplest form and iteratively add to it, before we convert it into a
// real graphql field.
// One challenge is that we need to "reuse" fields inside joins, and in some
// cases the end field type might be a List or a Scalar, so it is just easier
// to encapsulate that inside this struct rather than add lots of complexity.
type gqlField struct {
	Type *graphql.Object
	Args graphql.FieldConfigArgument
}

// newGraphQLSchema creates a new GraphQL schema wrapping the given provider
// with a schema that corresponds to the given set of tables.
func newGraphQLSchema(graph *schemaGraph, resolveFn graphql.FieldResolveFn) (graphql.Schema, error) {
	var (
		fields = make(map[string]gqlField)
		// These are the top-level query fields. Each of these fields
		// will correspond to each of the tables in the entire hierarchy.
		queryFields = make(graphql.Fields)
	)

	if len(graph.Nodes) == 0 {
		return graphql.Schema{}, nil
	}

	// Traverse the schema graph and add each node/table to the graphql fields
	graph.traverse(func(node *schemaNode) error {
		addGraphFields(*node.table, fields)
		return nil
	})

	// Create the relationships among the types using graph neighbours within
	// a certain distance of each other
	graph.traverse(func(node *schemaNode) error {
		paths := node.neighbours(graphJoinDistance)
		addGraphEdges(*node.table, paths, fields)
		return nil
	})

	// Finally, we want to populate the queryFields using the graphql types
	// we have created
	for _, field := range fields {
		queryFields[field.Type.Name()] = &graphql.Field{
			Type:    graphql.NewList(field.Type),
			Args:    field.Args,
			Resolve: resolveFn,
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

// Support for order_by argument
var orderByType = graphql.NewInputObject(graphql.InputObjectConfig{
	Name: "order_by",
	Fields: graphql.InputObjectConfigFieldMap{
		"table": &graphql.InputObjectFieldConfig{
			Type: graphql.String,
		},
		"field": &graphql.InputObjectFieldConfig{
			Type: graphql.String,
		},
		"order": &graphql.InputObjectFieldConfig{
			Type: graphql.String,
		},
	},
})

// Support for distinct_on argument
var distinctOnType = graphql.NewInputObject(graphql.InputObjectConfig{
	Name: "distinct_on",
	Fields: graphql.InputObjectConfigFieldMap{
		"table": &graphql.InputObjectFieldConfig{
			Type: graphql.String,
		},
		"field": &graphql.InputObjectFieldConfig{
			Type: graphql.String,
		},
	},
})

// addGraphFields updates the `gqlField` map containing GraphQL Field definitions
// with information for every field of the Table `t`, which is a table coming
// from the Bubbly Schema.
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

	// Add the _id field to the schema
	typeFields[tableIDField] = &graphql.Field{Type: graphql.String}
	gqlField.Args[tableIDField] = &graphql.ArgumentConfig{Type: graphql.String}

	gqlField.Args[filterID] = &graphql.ArgumentConfig{
		Type: graphQLFilterType(t.Name, gqlField.Args),
	}
	gqlField.Args[limitID] = &graphql.ArgumentConfig{
		Type: graphql.Int,
	}
	gqlField.Args[orderByID] = &graphql.ArgumentConfig{
		Type: graphql.NewList(orderByType),
	}
	// FIXME it would be great to define HERE that `distinct_on` requires `order_by`
	gqlField.Args[distinctOnID] = &graphql.ArgumentConfig{
		Type: graphql.NewList(distinctOnType),
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

// addGraphEdges ???
func addGraphEdges(t core.Table, paths []schemaPath, fields map[string]gqlField) {
	var field = fields[t.Name]
	for _, path := range paths {
		// We only care about the destination in the path and whether it is scalar.
		// The middle or passing edges will be included as their own paths
		var (
			edge                        = path[len(path)-1]
			dstField                    = fields[edge.node.table.Name]
			dstFieldType graphql.Output = dstField.Type
		)

		if !path.isScalar() {
			dstFieldType = graphql.NewList(dstFieldType)
		}
		field.Type.AddFieldConfig(edge.node.table.Name, &graphql.Field{
			Type: dstFieldType,
			Args: dstField.Args,
		})
	}
}

// graphQLFieldType ???
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
	filterID     = "filter"
	limitID      = "limit"
	orderByID    = "order_by"
	distinctOnID = "distinct_on"
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

// graphQLFilterType ???
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

// isValidValue ???
func isValidValue(value interface{}) bool {

	if value == nil {
		return false
	}

	// graphql-go passes nil maps as empty values
	if val, ok := value.(map[string]interface{}); ok {
		if len(val) == 0 {
			return false
		}
	}

	return true
}

// parseValueToMap ???
func parseValueToMap(astValue ast.Value) interface{} {
	switch astValue.GetKind() {
	case kinds.StringValue:
		return astValue.GetValue()
	case kinds.BooleanValue:
		return astValue.GetValue()
	case kinds.IntValue:
		return astValue.GetValue()
	case kinds.FloatValue:
		return astValue.GetValue()
	case kinds.ObjectValue:
		var (
			objFields = astValue.GetValue().([]*ast.ObjectField)
			obj       = make(map[string]interface{}, len(objFields))
		)
		for _, v := range objFields {
			obj[v.Name.Value] = parseValueToMap(v.Value)
		}
		return obj
	case kinds.ListValue:
		var (
			astList = astValue.GetValue().([]ast.Value)
			list    = make([]interface{}, 0, len(astList))
		)
		for _, v := range astList {
			list = append(list, parseValueToMap(v))
		}
		return list
	default:
		return nil
	}
}

// FIXME: what's going on here?
var mapScalar = graphql.NewScalar(graphql.ScalarConfig{
	Name:        "Map",
	Description: "The `Map` scalar type represents a Map for storing key/value pairs",
	Serialize: func(value interface{}) interface{} {
		return value
	},
	ParseValue: func(value interface{}) interface{} {
		return value
	},
	ParseLiteral: func(astValue ast.Value) interface{} {
		if astValue.GetKind() != kinds.ObjectValue {
			return nil
		}
		return parseValueToMap(astValue)
	},
})
