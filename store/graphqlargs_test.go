package store

import (
	"encoding/json"
	"fmt"
	"log"
	"testing"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/language/ast"
	"github.com/graphql-go/graphql/language/kinds"
)

// ########################################################################
// OrderBy - RECOMMENDED!
// This is more structured, and therefore much easier on our part, and the below
// is the only "custom" stuff you need.
// It adds a little overhead on the query side, but it's just so simple you cannot
// fault it :)
// ########################################################################
var orderByItem = graphql.NewInputObject(graphql.InputObjectConfig{
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

// ########################################################################
// OrderBy2 - Original, more difficult to parse
// This is less structured by Graphql standards, and therefore more effort on our end.
// We need to know the fields in a field map, and this will not work with aliases...
// ########################################################################
func orderByArgType(fields []string) *graphql.InputObject {
	fieldMap := graphql.InputObjectConfigFieldMap{}
	for _, field := range fields {
		fieldMap[field] = &graphql.InputObjectFieldConfig{
			Type: keyValue,
		}
	}
	return graphql.NewInputObject(graphql.InputObjectConfig{
		Name:   "order_by2",
		Fields: fieldMap,
	})
}

var keyValue = graphql.NewScalar(graphql.ScalarConfig{
	Name:        "KeyValue",
	Description: "TODO",
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
		var (
			objFields = astValue.GetValue().([]*ast.ObjectField)
			obj       = make(map[string]interface{}, len(objFields))
		)
		for _, v := range objFields {
			switch v.Value.GetKind() {
			case kinds.StringValue, kinds.BooleanValue, kinds.IntValue, kinds.FloatValue:
				obj[v.Name.Value] = v.Value.GetValue()
			default:
				return nil
			}
		}
		return obj
	},
})

// ########################################################################
// OrderBy3 - Original, just a variation this time with no structure
// We could post validate this, but yeah, extra work
// ########################################################################
func orderByItemParseLiteral(astValue ast.Value) interface{} {
	kind := astValue.GetKind()

	switch kind {
	case kinds.StringValue:
		return astValue.GetValue()
	case kinds.BooleanValue:
		return astValue.GetValue()
	case kinds.IntValue:
		return astValue.GetValue()
	case kinds.FloatValue:
		return astValue.GetValue()
	case kinds.ObjectValue:
		obj := make(map[string]interface{})
		for _, v := range astValue.GetValue().([]*ast.ObjectField) {
			obj[v.Name.Value] = orderByItemParseLiteral(v.Value)
		}
		return obj
	default:
		// Indicates an error by returning nil, e.g Lists are not supported
		return nil
	}
}

var orderByItemUnstructured = graphql.NewScalar(graphql.ScalarConfig{
	Name:        "OrderByStruct",
	Description: "TODO",
	Serialize: func(value interface{}) interface{} {
		return value
	},
	ParseValue: func(value interface{}) interface{} {
		return value
	},
	ParseLiteral: orderByItemParseLiteral,
})

// ########################################################################
// TEST CASES
// ########################################################################

var queries = []struct {
	query string
}{
	{
		query: `
		query {
			test_run(order_by: [{table: "network", field: "name", order: "ASC"}])
		}
	`,
	},
	{
		query: `
		query {
			test_run(order_by2: [{network: {asd: "asd"}}])
		}
	`,
	},
	{
		query: `
		query {
			test_run(order_by3: [{network: {asd: "asd"}}])
		}
	`,
	},
}

func TestGraphql(t *testing.T) {
	// Schema
	fields := graphql.Fields{
		"test_run": &graphql.Field{
			Args: graphql.FieldConfigArgument{
				"order_by": &graphql.ArgumentConfig{Type: graphql.NewList(orderByItem)},
				// TODO: you will need to create a list of the possible fields here that you could order by...
				// What if there are custom aliases?? This will not work, right?
				"order_by2": &graphql.ArgumentConfig{Type: graphql.NewList(orderByArgType([]string{"network", "test_case"}))},
				"order_by3": &graphql.ArgumentConfig{Type: graphql.NewList(orderByItemUnstructured)},
			},
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				fmt.Printf("Args: %#v\n", p.Args["order_by"])
				return "world", nil
			},
		},
	}
	queryType := graphql.NewObject(graphql.ObjectConfig{
		Name:   "query",
		Fields: fields,
	})
	schemaConfig := graphql.SchemaConfig{Query: queryType}
	schema, err := graphql.NewSchema(schemaConfig)
	if err != nil {
		t.Fatalf("failed to create new schema, error: %v", err)
	}

	for _, query := range queries {
		params := graphql.Params{Schema: schema, RequestString: query.query}
		r := graphql.Do(params)
		if len(r.Errors) > 0 {
			log.Fatalf("failed to execute graphql operation, errors: %+v", r.Errors)
		}
		rJSON, _ := json.Marshal(r)
		fmt.Printf("%s \n", rJSON)
	}
}
