package store

import (
	"encoding/json"
	"fmt"
	"log"
	"testing"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/language/ast"
	"github.com/stretchr/testify/require"
)

var enumTableType1 = graphql.NewEnum(graphql.EnumConfig{
	Name:        "Table1",
	Description: "The `Table` type is either `asc` or `desc`",
	Values: graphql.EnumValueConfigMap{
		"cve_patch": &graphql.EnumValueConfig{
			Value: "woaaa",
		},
	},
})

var enumTableType2 = graphql.NewEnum(graphql.EnumConfig{
	Name:        "Table2",
	Description: "The `Table` type is either `asc` or `desc`",
	Values: graphql.EnumValueConfigMap{
		"cve_patch": &graphql.EnumValueConfig{
			Value: "whatever",
		},
		"my_table": &graphql.EnumValueConfig{
			Value: "yoyoyo",
		},
	},
})

func TestGraphQL(t *testing.T) {
	// Schema
	fields := graphql.Fields{
		"cve": &graphql.Field{
			Args: graphql.FieldConfigArgument{
				"is_null": &graphql.ArgumentConfig{
					Type: graphql.NewList(enumTableType1),
				},
			},
			Type: graphql.NewObject(
				graphql.ObjectConfig{
					Name: "cve",
					Fields: graphql.Fields{
						"id": &graphql.Field{
							Name: "id",
							Type: graphql.String,
						},
						"cve_patch": &graphql.Field{
							Args: graphql.FieldConfigArgument{
								"is_null": &graphql.ArgumentConfig{
									Type: graphql.NewList(enumTableType2),
								},
							},
							Name: "cve_patch",
							Type: graphql.NewObject(graphql.ObjectConfig{
								Name: "cve_patch",
								Fields: graphql.Fields{
									"name": &graphql.Field{
										Name: "name",
										Type: graphql.String,
									},
								},
							}),
						},
					},
				},
			),
			// Args: graphql.FieldConfigArgument{
			// 	notNullID: &graphql.ArgumentConfig{
			// 		Type: enumNullFilters,
			// 	},
			// },
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				for a, b := range p.Args {
					fmt.Printf("%s: %#v\n", a, b)
				}
				// fmt.Printf("ARG: %#v\n", p.Info.FieldASTs[0].Arguments[0].Value.GetValue())
				for _, b := range p.Info.FieldASTs[0].Arguments {
					fmt.Printf("Arg: %s:  %#v\n", b.Name.Value, b.Value.GetKind())
				}
				for _, selection := range p.Info.FieldASTs[0].SelectionSet.Selections {
					// fmt.Printf("selection: %#v\n", selection)
					subField, ok := selection.(*ast.Field)
					require.True(t, ok)

					for _, b := range subField.Arguments {
						fmt.Printf("Arg: %s:  %#v\n", b.Name.Value, b.Value.GetKind())
					}
					// sel.GetSelectionSet().Selections
				}
				return "world", nil
			},
		},
	}
	rootQuery := graphql.ObjectConfig{Name: "RootQuery", Fields: fields}
	schemaConfig := graphql.SchemaConfig{
		Query: graphql.NewObject(rootQuery),
	}
	schema, err := graphql.NewSchema(schemaConfig)
	if err != nil {
		log.Fatalf("failed to create new schema, error: %v", err)
	}

	// Query
	query := `
		{
			cve(is_null: cve_patch) {
				id
			}
		}
	`
	// query := `
	// 	{
	// 		cve {
	// 			cve_patch {
	// 				name
	// 			}
	// 		}
	// 	}
	// `
	params := graphql.Params{Schema: schema, RequestString: query}
	r := graphql.Do(params)
	if len(r.Errors) > 0 {
		log.Fatalf("failed to execute graphql operation, errors: %+v", r.Errors)
	}
	rJSON, _ := json.Marshal(r)
	fmt.Printf("%s \n", rJSON) // {"data":{"hello":"world"}}
}
