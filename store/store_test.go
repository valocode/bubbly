package store

import (
	"database/sql"
	"errors"
	"fmt"
	"testing"

	_ "github.com/jackc/pgx/v4/stdlib"

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
}

func storeTests(t *testing.T, bCtx *env.BubblyContext, s *Store) {
	tables := testData.Tables(t)
	data := testData.DataBlocks(t)

	err := s.Apply(bCtx, tables)
	require.NoErrorf(t, err, "failed to apply schema")

	err = s.Save(bCtx, data)
	require.NoErrorf(t, err, "failed to save data")

	// Run the query tests
	for _, tt := range queryTests {
		t.Run(tt.name, func(t *testing.T) {
			actual := s.Query(tt.query)
			require.Emptyf(t, actual.Errors, "failed to execute query %s", tt.name)
			require.Equal(t, tt.expected, actual.Data, "query response is equal")
		})
	}

	// Test a simple resource
	resJSON := core.ResourceBlockJSON{
		ResourceBlockAlias: core.ResourceBlockAlias{
			ResourceKind:       "kind",
			ResourceName:       "name",
			ResourceAPIVersion: "some version",
			Metadata: &core.Metadata{
				Labels: map[string]string{"label": "is a label"},
			},
		},
		SpecRaw: "data {}",
	}
	d, err := resJSON.Data()

	require.NoError(t, err)

	err = s.Save(bCtx, core.DataBlocks{d})

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
	result := s.Query(resQuery)
	require.Empty(t, result.Errors)
	eventTests(t, bCtx, s, &d)
}

func eventTests(t *testing.T, bCtx *env.BubblyContext, s *Store, d *core.Data) {
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

	err := s.Save(bCtx, d2)

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
	result := s.Query(resQuery)
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

	err = s.Save(bCtx, d3)

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

	result = s.Query(resQuery)
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

	err = s.Save(bCtx, dataBlocks)

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

	result = s.Query(resQuery)

	assert.Empty(t, result.Errors)
	assert.NotNil(t, result)
}

// FIXME: CockroachDB support suspended due to lack or required features
func TestCockroach(t *testing.T) {

	t.SkipNow()
	// FIXME: this test is out of date, see TestPostgres

	pool, err := dockertest.NewPool("")

	require.NoErrorf(t, err, "failed to create dockertest pool")

	resource, err := pool.RunWithOptions(
		&dockertest.RunOptions{
			Repository: "cockroachdb/cockroach",
			Tag:        "v20.1.10",
			Cmd:        []string{"start-single-node", "--insecure"},
		},
	)

	t.Cleanup(func() {
		if err := pool.Purge(resource); err != nil {
			t.Fatalf("Could not purge resource: %s", err)
		}
	})

	require.NoErrorf(t, err, "failed to start docker")
	err = pool.Retry(func() error {
		db, err := sql.Open("pgx", fmt.Sprintf("postgresql://root@localhost:%s/events?sslmode=disable", resource.GetPort("26257/tcp")))
		if err != nil {
			return err
		}
		return db.Ping()
	})
	require.NoErrorf(t, err, "failed to connect to docker container")

	bCtx := env.NewBubblyContext()
	bCtx.StoreConfig.Provider = config.CockroachDBStore
	bCtx.StoreConfig.CockroachAddr = fmt.Sprintf("localhost:%s", resource.GetPort("26257/tcp"))
	bCtx.StoreConfig.CockroachDatabase = "defaultdb"
	bCtx.StoreConfig.CockroachUser = "root"
	bCtx.StoreConfig.CockroachPassword = "admin"

	s, err := New(bCtx)
	require.NoErrorf(t, err, "failed to initialize store")

	storeTests(t, bCtx, s)
}

func TestPostgres(t *testing.T) {
	pool, err := dockertest.NewPool("")
	require.NoErrorf(t, err, "failed to connect to Docker")

	resource, err := pool.RunWithOptions(
		&dockertest.RunOptions{
			Repository: "postgres",
			Tag:        "13.0",
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

	// Wait on the database
	var db *sql.DB
	err = pool.Retry(func() error {
		pgConnStr := fmt.Sprintf("postgresql://%s:%s@localhost:%s/%s?sslmode=disable", postgresUser, postgresPassword, resource.GetPort("5432/tcp"), postgresDatabase)
		db, err = sql.Open("pgx", pgConnStr)
		if err != nil {
			return err
		}
		return db.Ping()
	})
	require.NoErrorf(t, err, "failed to connect to the database container")

	// We only used the connection for waiting on database, bubbly will manage its own
	err = db.Close()
	require.NoError(t, err, "failed to close the connection to database")

	// Set up complete. Now for the test:
	bCtx := env.NewBubblyContext()
	bCtx.StoreConfig.Provider = config.PostgresStore
	bCtx.StoreConfig.PostgresAddr = fmt.Sprintf("localhost:%s", resource.GetPort("5432/tcp"))
	bCtx.StoreConfig.PostgresDatabase = postgresDatabase
	bCtx.StoreConfig.PostgresUser = postgresUser
	bCtx.StoreConfig.PostgresPassword = postgresPassword

	s, err := New(bCtx)
	assert.NoErrorf(t, err, "failed to initialize store")

	storeTests(t, bCtx, s)
}
