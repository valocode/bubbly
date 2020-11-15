package main

import (
	"log"
	"net/http"
	"os"

	"github.com/graphql-go/handler"
	"github.com/verifa/bubbly/api/core"
	"github.com/verifa/bubbly/store"
	"github.com/zclconf/go-cty/cty"
)

func main() {
	s, err := store.New(store.Config{
		Provider:         store.ProviderType(os.Getenv("PROVIDER")),
		PostgresAddr:     os.Getenv("POSTGRES_ADDR"),
		PostgresUser:     os.Getenv("POSTGRES_USER"),
		PostgresPassword: os.Getenv("POSTGRES_PASSWORD"),
		PostgresDatabase: os.Getenv("POSTGRES_DATABASE"),
	})
	if err != nil {
		log.Fatal(err)
	}

	if err := s.Create(tables); err != nil {
		log.Fatal(err)
	}
	if err := s.Save(data); err != nil {
		log.Fatal(err)
	}

	schema := s.Schema()
	h := handler.New(&handler.Config{
		Schema: &schema,
		Pretty: true,
	})

	fs := http.FileServer(http.Dir("static"))

	http.Handle("/graphql", h)
	http.Handle("/", fs)
	log.Println("storefront listening on :8111")
	log.Fatal(http.ListenAndServe(":8111", nil))
}

var tables = []core.Table{
	{
		Name: "product",
		Fields: []core.TableField{
			{
				Name: "name",
				Type: cty.String,
			},
		},
	},
	{
		Name: "project",
		Fields: []core.TableField{
			{
				Name: "name",
				Type: cty.String,
			},
			{
				Name: "product_id",
				Type: cty.Number,
			},
		},
	},
	{
		Name: "repo",
		Fields: []core.TableField{
			{
				Name: "name",
				Type: cty.String,
			},
			{
				Name: "project_id",
				Type: cty.Number,
			},
		},
		Tables: []core.Table{
			{
				Name: "repo_version",
				Fields: []core.TableField{
					{
						Name: "name",
						Type: cty.String,
					},
					{
						Name: "version",
						Type: cty.String,
					},
				},
			},
		},
	},
	{
		Name: "test_run",
		Fields: []core.TableField{
			{
				Name: "name",
				Type: cty.String,
			},
			{
				Name: "repo_version_id",
				Type: cty.Number,
			},
		},
		Tables: []core.Table{
			{
				Name: "test_set",
				Fields: []core.TableField{
					{
						Name: "name",
						Type: cty.String,
					},
				},
				Tables: []core.Table{
					{
						Name: "test_case",
						Fields: []core.TableField{
							{
								Name: "test_set_id",
								Type: cty.Number,
							},
							{
								Name: "name",
								Type: cty.String,
							},
							{
								Name: "status",
								Type: cty.String,
							},
						},
					},
				},
			},
		},
	},
}

var data = core.DataBlocks{
	{
		TableName: "test_run",
		Fields: []core.DataField{
			{
				Name:  "name",
				Value: cty.StringVal("run 1"),
			},
		},
		Data: core.DataBlocks{
			{
				TableName: "test_set",
				Fields: []core.DataField{
					{
						Name:  "name",
						Value: cty.StringVal("set 1"),
					},
				},
				Data: core.DataBlocks{
					{
						TableName: "test_case",
						Fields: []core.DataField{
							{
								Name:  "name",
								Value: cty.StringVal("case 1.1"),
							},
							{
								Name:  "status",
								Value: cty.StringVal("PASS"),
							},
						},
					},
					{
						TableName: "test_case",
						Fields: []core.DataField{
							{
								Name:  "name",
								Value: cty.StringVal("case 1.2"),
							},
							{
								Name:  "status",
								Value: cty.StringVal("PASS"),
							},
						},
					},
					{
						TableName: "test_case",
						Fields: []core.DataField{
							{
								Name:  "name",
								Value: cty.StringVal("case 1.3"),
							},
							{
								Name:  "status",
								Value: cty.StringVal("FAIL"),
							},
						},
					},
				},
			},
			{
				TableName: "test_set",
				Fields: []core.DataField{
					{
						Name:  "name",
						Value: cty.StringVal("set 2"),
					},
				},
				Data: core.DataBlocks{
					{
						TableName: "test_case",
						Fields: []core.DataField{
							{
								Name:  "name",
								Value: cty.StringVal("case 2.1"),
							},
							{
								Name:  "status",
								Value: cty.StringVal("FAIL"),
							},
						},
					},
					{
						TableName: "test_case",
						Fields: []core.DataField{
							{
								Name:  "name",
								Value: cty.StringVal("case 2.2"),
							},
							{
								Name:  "status",
								Value: cty.StringVal("FAIL"),
							},
						},
					},
				},
			},
		},
	},
}
