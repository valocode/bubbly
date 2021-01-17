package store

import (
	"database/sql"
	"fmt"
	"testing"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/ory/dockertest/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/verifa/bubbly/api/core"
	"github.com/verifa/bubbly/config"
	"github.com/verifa/bubbly/env"
	testData "github.com/verifa/bubbly/store/testdata"
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
			root(name: "test_value") {
				name
			}
		}
		`,
		expected: map[string]interface{}{
			"root": []interface{}{
				map[string]interface{}{
					"name": "test_value",
				},
			},
		},
	},
	{
		name: "root with child query",
		query: `
		{
			root {
				name
				child_a(name: "test_value") {
					name
				}
			}
		}
		`,
		expected: map[string]interface{}{
			"root": []interface{}{
				map[string]interface{}{
					"name": "test_value",
					"child_a": []interface{}{
						map[string]interface{}{
							"name": "test_value",
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
			root {
				name
				child_a(name: "test_value") {
					name
					grandchild_a(name: "join_value") {
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
					"name": "test_value",
					"child_a": []interface{}{
						map[string]interface{}{
							"name": "test_value",
							"grandchild_a": []interface{}{
								map[string]interface{}{
									"child_a": map[string]interface{}{
										"name": "test_value",
									},
									"name": "join_value",
								},
							},
						},
					},
				},
			},
		},
	},
}

func storeTests(t *testing.T, s *Store) {
	tables := testData.Tables(t)
	data := testData.DataBlocks(t)

	err := s.Apply(tables)
	require.NoErrorf(t, err, "failed to apply schema")

	err = s.Save(data)
	require.NoErrorf(t, err, "failed to save data")

	// Run the query tests
	for _, tt := range queryTests {
		t.Run(tt.name, func(t *testing.T) {
			actual, err := s.Query(tt.query)
			assert.NoErrorf(t, err, "failed to execute query %s", tt.name)
			assert.Equal(t, tt.expected, actual, "query response is equal")
		})
	}

	// Test a simple resource
	resJSON := core.ResourceBlockJSON{
		ResourceBlockAlias: core.ResourceBlockAlias{
			ResourceKind:       "kind",
			ResourceName:       "name",
			ResourceAPIVersion: "some version",
			Metadata: &core.Metadata{
				Labels:    map[string]string{"label": "is a label"},
				Namespace: "namespace",
			},
		},
		SpecRaw: "data {}",
	}
	d, err := resJSON.Data()
	err = s.Save(core.DataBlocks{d})
	require.NoError(t, err)

	resQuery := `
			{
				_resource(id: "namespace/kind/name") {
					name
					kind
					api_version
					metadata
					spec
				}
			}
		`
	_, err = s.Query(resQuery)
	assert.NoError(t, err)
}

func TestCockroach(t *testing.T) {
	pool, err := dockertest.NewPool("")
	require.NoErrorf(t, err, "failed to create dockertest pool")

	resource, err := pool.RunWithOptions(
		&dockertest.RunOptions{
			Repository: "cockroachdb/cockroach",
			Tag:        "v20.1.10",
			Cmd:        []string{"start-single-node", "--insecure"},
		},
	)
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
	assert.NoErrorf(t, err, "failed to initialize store")

	storeTests(t, s)
	// You can't defer this because os.Exit doesn't care for defer
	if err := pool.Purge(resource); err != nil {
		t.Fatalf("Could not purge resource: %s", err)
	}
}

func TestPosgres(t *testing.T) {
	pool, err := dockertest.NewPool("")
	require.NoErrorf(t, err, "failed to create dockertest pool")

	resource, err := pool.RunWithOptions(
		&dockertest.RunOptions{
			Repository: "postgres",
			Tag:        "13.0",
			Env: []string{
				"POSTGRES_PASSWORD=secret",
				"POSTGRES_DB=bubbly",
			},
		},
	)
	require.NoErrorf(t, err, "failed to start docker")
	err = pool.Retry(func() error {
		db, err := sql.Open("pgx", fmt.Sprintf("postgresql://postgres:secret@localhost:%s/bubbly?sslmode=disable", resource.GetPort("5432/tcp")))
		if err != nil {
			return err
		}
		return db.Ping()
	})
	require.NoErrorf(t, err, "failed to connect to docker container")

	bCtx := env.NewBubblyContext()
	bCtx.StoreConfig.Provider = config.PostgresStore
	bCtx.StoreConfig.PostgresAddr = fmt.Sprintf("localhost:%s", resource.GetPort("5432/tcp"))
	bCtx.StoreConfig.PostgresDatabase = "bubbly"
	bCtx.StoreConfig.PostgresUser = "postgres"
	bCtx.StoreConfig.PostgresPassword = "secret"

	s, err := New(bCtx)
	assert.NoErrorf(t, err, "failed to initialize store")

	storeTests(t, s)
	// You can't defer this because os.Exit doesn't care for defer
	if err := pool.Purge(resource); err != nil {
		t.Fatalf("Could not purge resource: %s", err)
	}
}
