package store

import (
	"database/sql"
	"fmt"
	"testing"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/ory/dockertest/v3"
	"github.com/stretchr/testify/require"

	"github.com/verifa/bubbly/api/core"
	"github.com/verifa/bubbly/config"
	"github.com/verifa/bubbly/env"
)

// TestEventTrigger tests that saving a new resource to the store results in
// both _resource and _event population using the event store trigger
func TestEventTrigger(t *testing.T) {
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

	s, err := New(bCtx)

	require.NoError(t, err)

	// save the blocks to the store
	err = s.Save(bCtx, core.DataBlocks{d})

	require.NoError(t, err)

	resQuery := fmt.Sprintf(`
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

	// query to make sure that the default trigger responsible for loading data into the _event table has worked
	result := s.Query(resQuery)

	t.Logf("%v", result.Data)

	require.NotEmpty(t, result)
}

// func TestDockerTest(t *testing.T) {
// 	pool, err := dockertest.NewPool("")
//
// 	require.NoErrorf(t, err, "failed to create dockertest pool")
//
// 	resource, err := pool.RunWithOptions(
// 		&dockertest.RunOptions{
// 			Repository: "cockroachdb/cockroach",
// 			Tag:        "v20.1.10",
// 			Cmd:        []string{"start-single-node", "--insecure"},
// 		},
// 	)
//
// 	storeProvider := "BUBBLY_STORE_PROVIDER=cockroachdb"
// 	cockroachAddr := fmt.Sprintf("COCKROACH_ADDR=localhost:%s", resource.GetPort("26257/tcp"))
//
// 	// bubbly, err2 := pool.BuildAndRun("bubbly", "../Dockerfile", []string{})
//
// 	name := "bubbly" + strconv.Itoa((rand.Intn(3000)))
//
// 	bubbly, err2 := pool.BuildAndRunWithOptions("../Dockerfile", &dockertest.RunOptions{
// 		Name: name,
// 		Env: []string{
// 			storeProvider,
// 			cockroachAddr,
// 		},
// 		Cmd: []string{"agent"},
// 	})
//
// 	t.Cleanup(func() {
// 		if err := pool.Purge(bubbly); err != nil {
// 			t.Fatalf("Could not purge bubbly: %s", err)
// 		}
// 	})
//
// 	require.NoErrorf(t, err2, "failed to start bubbly")
//
// 	bPort := bubbly.GetPort("8111/tcp")
// 	bAddr := fmt.Sprintf("http://localhost:%s", bPort)
//
// 	resp, err := http.Get(bAddr)
//
// 	require.NoError(t, err)
// 	require.NotNil(t, resp)
//
// }

// // TestRunTrigger tests that saving a new run resource to the store results
// // in a NATS publication
// func TestRunTrigger(t *testing.T) {
// 	bCtx := env.NewBubblyContext()
// 	bCtx.StoreConfig.Provider = config.CockroachDBStore
// 	bCtx.StoreConfig.CockroachAddr = "localhost:26257"
// 	bCtx.StoreConfig.CockroachDatabase = "defaultdb"
// 	bCtx.StoreConfig.CockroachUser = "root"
// 	bCtx.StoreConfig.CockroachPassword = "admin"
//
// 	// Test a simple resource
// 	resJSON := core.ResourceBlockJSON{
// 		ResourceBlockAlias: core.ResourceBlockAlias{
// 			ResourceKind:       "run",
// 			ResourceName:       "test_run",
// 			ResourceAPIVersion: "v1",
// 		},
// 		SpecRaw: "data {}",
// 	}
//
// 	d, err := resJSON.Data()
//
// 	require.NoError(t, err)
//
// 	s, err := New(bCtx)
//
// 	require.NoError(t, err)
//
// 	c, err := client.NewHTTP(bCtx)
//
// 	require.NoError(t, err)
//
// 	// resB, err := json.Marshal(d)
//
// 	require.NoError(t, err)
//
// 	err = c.Load(bCtx, core.DataBlocks{d})
//
// 	require.NoError(t, err)
//
// 	resQuery := fmt.Sprintf(`
// 			{
// 				%s(id: "run/test_run")  {
// 					id
// 					%s {
// 						status
// 						time
// 						error
// 					}
// 				}
// 			}
// 		`, core.ResourceTableName, core.EventTableName)
//
// 	// query to make sure that the default trigger responsible for loading data into the _event table has worked
// 	result := s.Query(resQuery)
//
// 	t.Logf("%v", result.Data)
//
// 	require.NotEmpty(t, result)
// }
