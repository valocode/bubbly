package store

import (
	"database/sql"
	"errors"
	"fmt"
	"path/filepath"
	"testing"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/rs/zerolog"

	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zclconf/go-cty/cty"

	"github.com/valocode/bubbly/api/core"
	"github.com/valocode/bubbly/config"
	"github.com/valocode/bubbly/env"
	"github.com/valocode/bubbly/events"
	testData "github.com/valocode/bubbly/store/testdata"
	"github.com/valocode/bubbly/test"
)

const (
	// postgresVersionTag defines the tag in
	// the Postgres Docker image repository.
	// The tests are run using that image.
	postgresVersionTag string = "13.0"
)

var queryTests = []struct {
	name     string
	query    string
	expected interface{}
}{
	{
		name: "root query",
		query: `
			{
				root(name: "first_root") {
					name
				}
			}
			`,
		expected: map[string]interface{}{
			"root": []interface{}{
				map[string]interface{}{
					"name": "first_root",
				},
			},
		},
	},
	{
		name: "root with child query",
		query: `
			{
				root(name: "first_root") {
					name
					child_a(name: "first_child") {
						name
					}
				}
			}
			`,
		expected: map[string]interface{}{
			"root": []interface{}{
				map[string]interface{}{
					"name": "first_root",
					"child_a": []interface{}{
						map[string]interface{}{
							"name": "first_child",
						},
					},
				},
			},
		},
	},
	{
		name: "root with grandchild query",
		query: `
			{
				root(name: "first_root") {
					name
					child_a(name: "first_child") {
						name
						grandchild_a(name: "second_grandchild") {
							name
							child_a {
								name
							}
						}
					}
				}
			}
			`,
		expected: map[string]interface{}{
			"root": []interface{}{
				map[string]interface{}{
					"name": "first_root",
					"child_a": []interface{}{
						map[string]interface{}{
							"name": "first_child",
							"grandchild_a": []interface{}{
								map[string]interface{}{
									"child_a": map[string]interface{}{
										"name": "first_child",
									},
									"name": "second_grandchild",
								},
							},
						},
					},
				},
			},
		},
	},
	{
		name: "cyclic query from rhs join",
		query: `
			{
				grandchild_a(name: "second_grandchild") {
					name
					child_a {
						name
						grandchild_a {
							name
						}
					}
				}
			}
			`,
		expected: map[string]interface{}{
			"grandchild_a": []interface{}{
				map[string]interface{}{
					"child_a": map[string]interface{}{
						"name": "first_child",
						"grandchild_a": []interface{}{
							map[string]interface{}{
								"name": "first_grandchild",
							},
							map[string]interface{}{
								"name": "second_grandchild",
							},
						},
					},
					"name": "second_grandchild",
				},
			},
		},
	},
	// TODO not jumping nodes anymore, move test to a new "must fail" test set
	/*
		{
			name: "jump node query",
			query: `
			{
				root(name: "first_root") {
					name
					grandchild_a {
						name
					}
				}
			}
			`,
			expected: map[string]interface{}{
				"root": []interface{}{
					map[string]interface{}{
						"name": "first_root",
						"grandchild_a": []interface{}{
							map[string]interface{}{
								"name": "first_grandchild",
							},
							map[string]interface{}{
								"name": "second_grandchild",
							},
						},
					},
				},
			},
		},
	*/
	{
		name: "two sibling blocks query",
		query: `
			{
				root(name: "first_root") {
					name
					child_a(name: "first_child") {
						name
					}
					child_c(name: "sibling_child") {
						name
					}
				}
			}
			`,
		expected: map[string]interface{}{
			"root": []interface{}{
				map[string]interface{}{
					"name": "first_root",
					"child_a": []interface{}{
						map[string]interface{}{
							"name": "first_child",
						},
					},
					"child_c": []interface{}{
						map[string]interface{}{
							"name": "sibling_child",
						},
					},
				},
			},
		},
	},
}

var sqlGenTests = []struct {
	name   string
	schema string
	data   string
	query  string
	want   interface{}
}{
	{
		name:   "unrelated tables",
		schema: "tables0.hcl",
		data:   "data0.hcl",
		query: `
			{
				A(whaat: "va1") {
					whaat
				}
				B(whbbt: "vb1") {
					whbbt
				}
			}
			`,
		want: map[string]interface{}{
			"A": []interface{}{
				map[string]interface{}{
					"whaat": "va1",
				},
			},
			"B": []interface{}{
				map[string]interface{}{
					"whbbt": "vb1",
				},
			},
		},
	},
	{
		name:   "single level nested field",
		schema: "tables1.hcl",
		data:   "data1.hcl",
		query: `
			{
				A(whaat: "va1") {
					whaat
					B(whbbt: "vb1") {
						whbbt
					}
				}
			}
			`,
		want: map[string]interface{}{
			"A": []interface{}{
				map[string]interface{}{
					"whaat": "va1",
					"B": []interface{}{
						map[string]interface{}{
							"whbbt": "vb1",
						},
					},
				},
			},
		},
	},
	{
		name:   "single level nested field inverse simple arguments",
		schema: "tables1.hcl",
		data:   "data1.hcl",
		query: `
			{
				B(whbbt: "vb1") {
					whbbt
					A(whaat: "va1") {
						whaat
					}
				}
			}
			`,
		want: map[string]interface{}{
			"B": []interface{}{
				map[string]interface{}{
					"whbbt": "vb1",
					"A": map[string]interface{}{
						"whaat": "va1",
					},
				},
			},
		},
	},
	{
		name:   "single level sibling fields",
		schema: "tables2.hcl",
		data:   "data2.hcl",
		query: `
			{
				A {
					whaat
					B {
						whbbt
					}
					C {
						whcct
					}
				}
			}
			`,
		want: map[string]interface{}{
			"A": []interface{}{
				map[string]interface{}{
					"whaat": "va1",
					"B": []interface{}{
						map[string]interface{}{
							"whbbt": "vb1",
						},
					},
					"C": []interface{}{
						map[string]interface{}{
							"whcct": "vc1",
						},
					},
				},
			},
		},
	},
	{
		name:   "multi level nested fields",
		schema: "tables3.hcl",
		data:   "data3.hcl",
		query: `
			{
				A {
					whaat
					B {
						whbbt
						C {
							whcct
						}
					}
				}
			}
			`,
		want: map[string]interface{}{
			"A": []interface{}{
				map[string]interface{}{
					"whaat": "va1",
					"B": []interface{}{
						map[string]interface{}{
							"whbbt": "vb1",
							"C": []interface{}{
								map[string]interface{}{
									"whcct": "vc1",
								},
							},
						},
					},
				},
			},
		},
	},
	{
		name:   "multi level nested and sibling fields",
		schema: "tables4.hcl",
		data:   "data4.hcl",
		query: `
			{
				A {
					whaat
					B {
						whbbt
						C {
							whcct
						}
					}
					D {
						whddt
					}
				}
			}
			`,
		want: map[string]interface{}{
			"A": []interface{}{
				map[string]interface{}{
					"whaat": "va1",
					"B": []interface{}{
						map[string]interface{}{
							"whbbt": "vb1",
							"C": []interface{}{
								map[string]interface{}{
									"whcct": "vc1",
								},
							},
						},
					},
					"D": []interface{}{
						map[string]interface{}{
							"whddt": "vd1",
						},
					},
				},
			},
		},
	},
	{
		name:   "graphql field order 1",
		schema: "tables5.hcl",
		data:   "data5.hcl",
		query: `
		{
			testrun(limit: 1) {
				configuration {
					name
				}
				location(name: "Deep Dark Wood") {
					name
				}
				ok
				version {
					name
				}
			}
		}`,
		want: map[string]interface{}{
			"testrun": []interface{}{
				map[string]interface{}{
					"ok": true,
					"location": map[string]interface{}{
						"name": "Deep Dark Wood",
					},
					"configuration": map[string]interface{}{
						"name": "Primitive",
					},
					"version": map[string]interface{}{
						"name": "v1.0.1",
					},
				},
			},
		},
	},
	{
		name:   "graphql field order 2",
		schema: "tables5.hcl",
		data:   "data5.hcl",
		query: `
		{
			testrun(limit: 1) {
				ok
				location(name: "Deep Dark Wood") {
					name
				}
				version {
					name
				}
				configuration {
					name
				}
			}
		}`,
		want: map[string]interface{}{
			"testrun": []interface{}{
				map[string]interface{}{
					"ok": true,
					"location": map[string]interface{}{
						"name": "Deep Dark Wood",
					},
					"configuration": map[string]interface{}{
						"name": "Primitive",
					},
					"version": map[string]interface{}{
						"name": "v1.0.1",
					},
				},
			},
		},
	},
	{
		name:   "graphql field order 3",
		schema: "tables5.hcl",
		data:   "data5.hcl",
		query: `
		{
			testrun(limit: 1) {
				configuration {
					name
				}
				location(name: "Deep Dark Wood") {
					name
				}
			}
		}`,
		want: map[string]interface{}{
			"testrun": []interface{}{
				map[string]interface{}{
					"location": map[string]interface{}{
						"name": "Deep Dark Wood",
					},
					"configuration": map[string]interface{}{
						"name": "Primitive",
					},
				},
			},
		},
	},
	{
		name:   "graphql order_by 1",
		schema: "tables5.hcl",
		data:   "data5.hcl",
		query: `
		{
			testrun(order_by: {_id: "DESC"}) {
				configuration(order_by: {name: "DESC"}) {
					name
				}
				location(order_by: {name: "DESC"}) {
					name
				}
			}
		}`,
		want: map[string]interface{}{
			"testrun": []interface{}{
				map[string]interface{}{
					"location": map[string]interface{}{
						"name": "Secret Underground Facility on the Moon",
					},
					"configuration": map[string]interface{}{
						"name": "Magic",
					},
				},
				map[string]interface{}{
					"location": map[string]interface{}{
						"name": "Gold Coast City Skyline",
					},
					"configuration": map[string]interface{}{
						"name": "Approved by Wombats",
					},
				},
				map[string]interface{}{
					"location": map[string]interface{}{
						"name": "Deep Dark Wood",
					},
					"configuration": map[string]interface{}{
						"name": "Primitive",
					},
				},
			},
		},
	},
	{
		name:   "graphql order_by totoro asc",
		schema: "tables6.hcl",
		data:   "data6.hcl",
		query: `
		{
			hideaways(location: "Deep Dark Wood") {
				location
				ready
				distance_from_x
				crew {
					count
					characters(order_by: {name: "ASC"}) {
						name
					}
				}
			}
		}`,
		want: map[string]interface{}{
			"hideaways": []interface{}{
				map[string]interface{}{
					"location":        "Deep Dark Wood",
					"distance_from_x": 1500,
					"ready":           true,
					"crew": []interface{}{
						map[string]interface{}{
							"count": 1,
							"characters": map[string]interface{}{
								"name": "Chibi-Totoro",
							},
						},
						map[string]interface{}{
							"count": 1,
							"characters": map[string]interface{}{
								"name": "Chuu-Totoro",
							},
						},
						map[string]interface{}{
							"count": 1,
							"characters": map[string]interface{}{
								"name": "Oh-Totoro",
							},
						},
					},
				},
			},
		},
	},
	{
		name:   "graphql order_by totoro desc",
		schema: "tables6.hcl",
		data:   "data6.hcl",
		query: `
		{
			hideaways(location: "Deep Dark Wood") {
				location
				ready
				distance_from_x
				crew {
					count
					characters(order_by: {name: "DESC"}) {
						name
					}
				}
			}
		}`,
		want: map[string]interface{}{
			"hideaways": []interface{}{
				map[string]interface{}{
					"location":        "Deep Dark Wood",
					"distance_from_x": 1500,
					"ready":           true,
					"crew": []interface{}{
						map[string]interface{}{
							"count": 1,
							"characters": map[string]interface{}{
								"name": "Oh-Totoro",
							},
						},
						map[string]interface{}{
							"count": 1,
							"characters": map[string]interface{}{
								"name": "Chuu-Totoro",
							},
						},
						map[string]interface{}{
							"count": 1,
							"characters": map[string]interface{}{
								"name": "Chibi-Totoro",
							},
						},
					},
				},
			},
		},
	},
	{
		name:   "graphql order_by full house",
		schema: "tables6.hcl",
		data:   "data6.hcl",
		query: `
		{
			hideaways(order_by: {location: "DESC"}){
				location
				ready
				distance_from_x
				crew(order_by: {count: "DESC"}) {
					count
					characters(order_by: {name: "ASC"}) {
						name
					}
				}
			}
		}`,
		want: map[string]interface{}{
			"hideaways": []interface{}{
				map[string]interface{}{
					"location":        "Secret Underground Facility on the Moon",
					"distance_from_x": 384400,
					"ready":           true,
					"crew": []interface{}{
						map[string]interface{}{
							"count": 42,
							"characters": map[string]interface{}{
								"name": "Little Green Man",
							},
						},
						map[string]interface{}{
							"count": 1,
							"characters": map[string]interface{}{
								"name": "Little Green Man Leader",
							},
						},
					},
				},
				map[string]interface{}{
					"location":        "Gold Coast Villa",
					"distance_from_x": 12500,
					"ready":           false,
					"crew": []interface{}{
						map[string]interface{}{
							"count": 2,
							"characters": map[string]interface{}{
								"name": "Smol European Mouse",
							},
						},
						map[string]interface{}{
							"count": 1,
							"characters": map[string]interface{}{
								"name": "Old Horse with Shady Past",
							},
						},
						map[string]interface{}{
							"count": 1,
							"characters": map[string]interface{}{
								"name": "Wombat",
							},
						},
					},
				},
				map[string]interface{}{
					"location":        "Deep Dark Wood",
					"distance_from_x": 1500,
					"ready":           true,
					"crew": []interface{}{
						map[string]interface{}{
							"count": 1,
							"characters": map[string]interface{}{
								"name": "Chibi-Totoro",
							},
						},
						map[string]interface{}{
							"count": 1,
							"characters": map[string]interface{}{
								"name": "Chuu-Totoro",
							},
						},
						map[string]interface{}{
							"count": 1,
							"characters": map[string]interface{}{
								"name": "Oh-Totoro",
							},
						},
					},
				},
			},
		},
	},
	{
		name:   "graphql multiple order by levels",
		schema: "tables7.hcl",
		data:   "data7.hcl",
		query: `
		{
			testrun(order_by: {finish_epoch: "DESC"}) {
				configuration(order_by: {name: "ASC"}) {
					name
				}
				location(name: "Oliver's NEON workstation", order_by: {name: "ASC"}) {
					name
				}
				version(order_by: {name: "ASC"}) {
					name
				}
				ok
				finish_epoch
			}
		}`,
		want: map[string]interface{}{
			"testrun": []interface{}{
				map[string]interface{}{
					"ok":           true,
					"finish_epoch": "1611111111",
					"location": map[string]interface{}{
						"name": "Oliver's NEON workstation",
					},
					"configuration": map[string]interface{}{
						"name": "So-so",
					},
					"version": map[string]interface{}{
						"name": "v1.0.1",
					},
				},
				map[string]interface{}{
					"ok":           true,
					"finish_epoch": "1411111111",
					"location": map[string]interface{}{
						"name": "Oliver's NEON workstation",
					},
					"configuration": map[string]interface{}{
						"name": "So-so",
					},
					"version": map[string]interface{}{
						"name": "v1.0.1",
					},
				},
				map[string]interface{}{
					"ok":           false,
					"finish_epoch": "1311111111",
					"location": map[string]interface{}{
						"name": "Oliver's NEON workstation",
					},
					"configuration": map[string]interface{}{
						"name": "So-so",
					},
					"version": map[string]interface{}{
						"name": "v1.0.1",
					},
				},
			},
		},
	},
	{
		name:   "graphql limit on root",
		schema: "tables8.hcl",
		data:   "data8.hcl",
		query: `
		{
			events(
				order_by: {timestamp: "DESC"},
				limit: 2
			) {
				timestamp
			}
		}`,
		want: map[string]interface{}{
			"events": []interface{}{
				map[string]interface{}{
					"timestamp": 1010,
				},
				map[string]interface{}{
					"timestamp": 30,
				},
			},
		},
	},
}

func applySchemaOrDie(t *testing.T, bCtx *env.BubblyContext, s *Store, fromFile string) {
	t.Helper()

	tables := testData.Tables(t, bCtx, fromFile)

	err := s.Apply(DefaultTenantName, tables, true)
	require.NoErrorf(t, err, "failed to apply schema")
}

func loadTestDataOrDie(t *testing.T, bCtx *env.BubblyContext, s *Store, fromFile string) {
	t.Helper()

	data := testData.DataBlocks(t, bCtx, fromFile)

	err := s.Save(DefaultTenantName, data)
	require.NoErrorf(t, err, "failed to save test data into the store")
}

func createResJSONOrDie(t *testing.T) core.Data {
	t.Helper()

	res := core.ResourceBlock{
		ResourceKind:       "kind",
		ResourceName:       "name",
		ResourceAPIVersion: "some version",
		Metadata: &core.Metadata{
			Labels: map[string]string{"label": "is a label"},
		},
		SpecRaw: "data {}",
	}
	d, err := res.Data()
	require.NoError(t, err)
	return d
}

// runQueryTestsOrDie runs all basic query tests, or fails hard on error.
func runQueryTestsOrDie(t *testing.T, bCtx *env.BubblyContext, s *Store) {
	t.Helper()

	for _, tt := range queryTests {
		t.Run(tt.name, func(t *testing.T) {
			actual, err := s.Query(DefaultTenantName, tt.query)
			require.NoError(t, err)
			require.Emptyf(t, actual.Errors, "failed to execute query %s", tt.name)
			require.Equal(t, tt.expected, actual.Data, "query response is equal")
		})
	}
}

// runResourceTestsOrDie runs all resource-related tests, or fails hard on error.
func runResourceTestsOrDie(t *testing.T, bCtx *env.BubblyContext, s *Store) {
	t.Helper()

	t.Run("resource", func(t *testing.T) {

		data := createResJSONOrDie(t)

		err := s.Save(DefaultTenantName, core.DataBlocks{data})
		require.NoError(t, err)

		resQuery := `
				{
					_resource(id: "kind/name") {
						name
						kind
						api_version
						metadata
						spec
					}
				}
			`

		result, err := s.Query(DefaultTenantName, resQuery)
		require.NoError(t, err)
		require.Empty(t, result.Errors)
	})
}

// runEventTestsOrDie runs all event-related tests, or fails hard on error.
func runEventTestsOrDie(t *testing.T, bCtx *env.BubblyContext, s *Store) {
	t.Helper()

	d := createResJSONOrDie(t)

	dSource := core.Data{
		TableName: core.ResourceTableName,
		Fields: core.DataFields{
			"id": d.Fields["id"],
		},
	}
	// add an event entry to the core.Event
	d2 := core.DataBlocks{dSource,
		{
			TableName: core.EventTableName,
			Fields: map[string]cty.Value{
				"status": cty.StringVal(events.ResourceCreatedUpdated.String()),
				"time":   cty.StringVal(events.TimeNow()),
			},
			// the _id value of the resource's row entry in _resource will be
			// mapped to the _resource_id column value in _event
			Joins: []string{core.ResourceTableName},
		},
	}

	err := s.Save(DefaultTenantName, d2)

	require.NoError(t, err)

	// First, check that we can query the resource under _resource and _event
	// tables
	resQuery := fmt.Sprintf(`
			{
				_resource(id: "kind/name")  {
					id
					%s {
						status
					}
				}
			}
		`, core.EventTableName)
	result, err := s.Query(DefaultTenantName, resQuery)
	require.NoError(t, err)
	require.Empty(t, result.Errors)

	a := result.Data.(map[string]interface{})[core.ResourceTableName].([]interface{})

	// Go through the returned query data and validate that the resource
	// exists in the _resource table,
	// there is a valid FK association with the resource in the status table
	for _, v := range a {
		x := v.(map[string]interface{})
		for k, v := range x {
			switch k {
			case "id":
				require.Equal(t, "kind/name", v)
			case core.EventTableName:
				et := v.([]interface{})
				for _, v := range et {
					e := v.(map[string]interface{})
					for k, v := range e {
						switch k {
						case "status":
							require.Equal(t, events.ResourceCreatedUpdated.String(), v)
						}
					}
				}
			}
		}
	}

	// test adding an event to the "kind/name" resource
	d3 := core.DataBlocks{dSource,
		{
			TableName: core.EventTableName,
			Fields: map[string]cty.Value{
				"status": cty.StringVal(events.ResourceDestroyed.String()),
				"time":   cty.StringVal(events.TimeNow()),
			},
			// Join says "this is the table from which I want to JOIN to".
			// As a result,
			// the _id pulled from the above datablock will be mapped to the
			// _resource_id column of this row entry in _event
			Joins: []string{core.ResourceTableName},
		},
	}

	err = s.Save(DefaultTenantName, d3)

	require.NoError(t, err)

	resQuery = fmt.Sprintf(`
			{
				%s(id: "kind/name")  {
					id
					%s {
						status
						time
					}
				}
			}
		`, core.ResourceTableName, core.EventTableName)

	result, err = s.Query(DefaultTenantName, resQuery)
	require.NoError(t, err)
	assert.Empty(t, result.Errors)

	resEvents := result.Data.(map[string]interface{})[core.ResourceTableName].([]interface{})

	// verify that the number of events stored for the "kind/name"
	// resource is 2
	for _, v := range resEvents {
		require.Equal(t, 2, len(v.(map[string]interface{})))
	}

	// check that we can load a resourceOutput to the store
	resOutput := core.ResourceOutput{
		ID:     "kind/name",
		Status: events.ResourceApplyFailure,
		Error:  errors.New("cannot get output of a null extract"),
		Value:  cty.NilVal,
	}

	dataBlocks, err := resOutput.DataBlocks()

	require.NoError(t, err)

	err = s.Save(DefaultTenantName, dataBlocks)

	require.NoError(t, err)

	resQuery = fmt.Sprintf(`
			{
				%s(id: "kind/name") {
					id
					%s(status: "ApplyFailure") {
						status
						time
						error
					}
				}
			}
		`, core.ResourceTableName, core.EventTableName)

	result, err = s.Query(DefaultTenantName, resQuery)
	require.NoError(t, err)
	assert.Empty(t, result.Errors)
	assert.NotNil(t, result)
}

// TODO: This was copied from TestPostgres for now. Extract the common func later
func TestPostgresSQLGen(t *testing.T) {

	pool, err := dockertest.NewPool("")
	require.NoErrorf(t, err, "failed to connect to Docker")

	// For SQL generation tests we run a new container for every test,
	// because we'd like sensible table names without much clutter
	// in the namespace. This also allows us to extend each test Schema
	// and add more test cases for more complex GraphQL queries in the future.

	for _, tt := range sqlGenTests {
		t.Run(tt.name, func(t *testing.T) {

			resource, err := pool.RunWithOptions(
				&dockertest.RunOptions{
					Repository: "postgres",
					Tag:        postgresVersionTag,
					Env: []string{
						"POSTGRES_PASSWORD=" + postgresPassword,
						"POSTGRES_DB=" + postgresDatabase,
					},
					Cmd: []string{
						"-c",
						"log_statement=all",
					},
				},
				func(hc *docker.HostConfig) {
					hc.AutoRemove = true
					hc.RestartPolicy = docker.RestartPolicy{
						Name: "no",
					}
				},
			)
			require.NoErrorf(t, err, "failed to start a container")

			t.Cleanup(func() {
				if err := pool.Purge(resource); err != nil {
					t.Fatal("failed to remove a container or a volume:", err)
				}
			})
			resource.Expire(360) // Tell docker to hard kill the container in X seconds

			// Wait until the database is ready
			err = waitUntilDatabaseIsReady(t, pool, resource)
			require.NoError(t, err, "failed database connection smoke test")

			// Initialise the Bubbly context
			bCtx := env.NewBubblyContext()
			bCtx.UpdateLogLevel(zerolog.DebugLevel)

			// Configure the Bubbly Store
			bCtx.StoreConfig.Provider = config.PostgresStore
			bCtx.StoreConfig.PostgresAddr = fmt.Sprintf("localhost:%s", resource.GetPort("5432/tcp"))
			bCtx.StoreConfig.PostgresDatabase = postgresDatabase
			bCtx.StoreConfig.PostgresUser = postgresUser
			bCtx.StoreConfig.PostgresPassword = postgresPassword

			// Initialise the Bubbly Store
			s, err := New(bCtx)
			require.NoErrorf(t, err, "failed to initialize store")

			// Apply the Bubbly Schema to the Bubbly Store
			applySchemaOrDie(t, bCtx, s, filepath.Join("testdata", "sqlgen", tt.schema))

			// Load the test data into the Bubbly Store
			loadTestDataOrDie(t, bCtx, s, filepath.Join("testdata", "sqlgen", tt.data))

			// Run the test
			have, err := s.Query(DefaultTenantName, tt.query)
			require.NoError(t, err)
			require.Emptyf(t, have.Errors, "failed to execute query %s", tt.name)
			require.Equal(t, tt.want, have.Data, "query response is equal")
		})
	}
}

func TestSchemaSync(t *testing.T) {
	bCtx := env.NewBubblyContext()
	resource := test.RunPostgresDocker(bCtx, t)
	bCtx.StoreConfig.PostgresAddr = fmt.Sprintf("localhost:%s", resource.GetPort("5432/tcp"))

	_, err := New(bCtx)
	assert.NoError(t, err)

}

func TestPostgres(t *testing.T) {

	pool, err := dockertest.NewPool("")
	require.NoErrorf(t, err, "failed to connect to Docker")

	resource, err := pool.RunWithOptions(
		&dockertest.RunOptions{
			Repository: "postgres",
			Tag:        postgresVersionTag,
			Env: []string{
				"POSTGRES_PASSWORD=" + postgresPassword,
				"POSTGRES_DB=" + postgresDatabase,
			},
		},
		func(hc *docker.HostConfig) {
			hc.AutoRemove = true
			hc.RestartPolicy = docker.RestartPolicy{
				Name: "no",
			}
		},
	)
	require.NoErrorf(t, err, "failed to start a container")

	t.Cleanup(func() {
		if err := pool.Purge(resource); err != nil {
			t.Fatal("failed to remove a container or a volume:", err)
		}
	})
	resource.Expire(360) // Tell docker to hard kill the container in X seconds

	// Wait until the database is ready
	err = waitUntilDatabaseIsReady(t, pool, resource)
	require.NoError(t, err, "failed database connection smoke test")

	// Initialise the Bubbly context
	bCtx := env.NewBubblyContext()
	bCtx.StoreConfig.Provider = config.PostgresStore
	bCtx.StoreConfig.PostgresAddr = fmt.Sprintf("localhost:%s", resource.GetPort("5432/tcp"))
	bCtx.StoreConfig.PostgresDatabase = postgresDatabase
	bCtx.StoreConfig.PostgresUser = postgresUser
	bCtx.StoreConfig.PostgresPassword = postgresPassword

	// Initialise the Bubbly Store
	s, err := New(bCtx)
	assert.NoErrorf(t, err, "failed to initialize store")

	// Apply the Bubbly Schema to the Bubbly Store
	applySchemaOrDie(t, bCtx, s, filepath.Join("testdata", "tables.hcl"))

	// Load the test data into the Bubbly Store
	loadTestDataOrDie(t, bCtx, s, filepath.Join("testdata", "data.hcl"))

	// Run (sub)tests
	runQueryTestsOrDie(t, bCtx, s)
	runResourceTestsOrDie(t, bCtx, s)
	runEventTestsOrDie(t, bCtx, s)
}

// Tests that should bubbly go down, on reinitialisation the Store correctly
// pulls the most recently applied schema
func TestPostgresReinitialisation(t *testing.T) {

	pool, err := dockertest.NewPool("")
	require.NoErrorf(t, err, "failed to connect to Docker")

	resource, err := pool.RunWithOptions(
		&dockertest.RunOptions{
			Repository: "postgres",
			Tag:        postgresVersionTag,
			Env: []string{
				"POSTGRES_PASSWORD=" + postgresPassword,
				"POSTGRES_DB=" + postgresDatabase,
			},
		},
		func(hc *docker.HostConfig) {
			hc.AutoRemove = true
			hc.RestartPolicy = docker.RestartPolicy{
				Name: "no",
			}
		},
	)
	require.NoErrorf(t, err, "failed to start a container")

	t.Cleanup(func() {
		if err := pool.Purge(resource); err != nil {
			t.Fatal("failed to remove a container or a volume:", err)
		}
	})
	resource.Expire(360) // Tell docker to hard kill the container in X seconds

	// Wait until the database is ready
	err = waitUntilDatabaseIsReady(t, pool, resource)
	require.NoError(t, err, "failed database connection smoke test")

	// Initialise the Bubbly context
	bCtx := env.NewBubblyContext()
	bCtx.StoreConfig.Provider = config.PostgresStore
	bCtx.StoreConfig.PostgresAddr = fmt.Sprintf("localhost:%s", resource.GetPort("5432/tcp"))
	bCtx.StoreConfig.PostgresDatabase = postgresDatabase
	bCtx.StoreConfig.PostgresUser = postgresUser
	bCtx.StoreConfig.PostgresPassword = postgresPassword

	// Initialise the Bubbly Store
	s, err := New(bCtx)
	assert.NoErrorf(t, err, "failed to initialize store")

	// grab the "base" bubbly schema, set up from the internalTables
	baseSchema, err := s.currentBubblySchema(DefaultTenantName)

	require.NoError(t, err)

	// Apply a Bubbly Schema to the Bubbly Store
	applySchemaOrDie(t, bCtx, s, filepath.Join("testdata", "tables.hcl"))

	// feign re-initialisation of the Store, which should fetch the latest
	// applied schema (from applySchemaOrDie)
	s, err = New(bCtx)
	require.NoErrorf(t, err, "failed to re-initialise the store")

	newSchema, err := s.currentBubblySchema(DefaultTenantName)
	require.NoErrorf(t, err, "failed to get current bubbly schema")

	// this should return the schema that was formed from applySchemaOrDie,
	// _not_ the baseSchema at row 0 in the _schema table
	require.NotEqual(t, baseSchema, newSchema)

}

// TODO: extract into a helper as a similar block of code is used elsewhere in store (?) tests
// waitUntilDatabaseIsReady employs exponential backoff retry algo to verify that the containerised database is up.
// It returns an error if the connection to the database had failed.
func waitUntilDatabaseIsReady(t *testing.T, pool *dockertest.Pool, resource *dockertest.Resource) error {

	pgConnStr := fmt.Sprintf("postgresql://%s:%s@localhost:%s/%s?sslmode=disable",
		postgresUser, postgresPassword, resource.GetPort("5432/tcp"), postgresDatabase)

	return pool.Retry(func() error {

		// Open does not necessarily establish the connection,
		// it may just validate the arguments.
		db, err := sql.Open("pgx", pgConnStr)

		// If Open failed miserably though, that's enough
		// to conclude that a connection cannot be established
		// right now. Defer to the caller to decide how to proceed.
		if err != nil {
			return err
		}

		// So, Open had succeeded. Now can attempt Ping which
		// will actually connect to the database.
		err = db.Ping()

		// Ping must have failed, defer to the caller
		// to decide how to proceed.
		if err != nil {
			return err
		}

		// Ping was successfull, close the connection as bubbly
		// will manage its own.
		err = db.Close()

		// Defer the decision to retry to the caller. If the DB
		// was closed successfully though, the caller would not retry.
		return err
	})
}
