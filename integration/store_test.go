// +build integration

package integration

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/verifa/bubbly/api/core"
	"github.com/verifa/bubbly/store"
	"github.com/zclconf/go-cty/cty"
)

func TestStore(t *testing.T) {
	s, err := store.New(store.Config{
		Provider:         store.ProviderType(os.Getenv("PROVIDER")),
		PostgresAddr:     os.Getenv("POSTGRES_ADDR"),
		PostgresUser:     os.Getenv("POSTGRES_USER"),
		PostgresPassword: os.Getenv("POSTGRES_PASSWORD"),
		PostgresDatabase: os.Getenv("POSTGRES_DATABASE"),
	})

	assert.NoError(t, err)

	// Stores that don't have type information can't save anything.
	assert.Contains(t, s.Save(nil).Error(), "no type information")

	tables := []core.Table{
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

	err = s.Create(tables)
	assert.NoError(t, err)

	// Create baseline data.
	err = s.Save(core.DataBlocks{
		{
			TableName: "product",
			Fields: []core.DataField{
				{
					Name:  "name",
					Value: cty.StringVal("Super 111"),
				},
			},
		},
		{
			TableName: "project",
			Fields: []core.DataField{
				{
					Name:  "name",
					Value: cty.StringVal("Flagship"),
				},
				{
					Name:  "product_id",
					Value: cty.NumberIntVal(1),
				},
			},
		},
	})

	assert.NoError(t, err)

	err = s.Save(core.DataBlocks{
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
	})

	assert.NoError(t, err)

	n, err := s.Query(`{
	   	test_set(name: "set 1") {
	   			name
				test_case(status: "PASS") {
					name
	   				status
	   			}
	   		}
		}`)

	assert.NoError(t, err)

	b, err := json.MarshalIndent(n, "", "\t")
	assert.NoError(t, err)

	fmt.Println(string(b))
}
