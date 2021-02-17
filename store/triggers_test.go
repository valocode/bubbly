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

// TestTriggers tests that saving a new resource to the store results in both _resource and _event population using
// the default trigger
func TestTriggers(t *testing.T) {
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
				Labels:    map[string]string{"label": "is a label"},
				Namespace: "namespace",
			},
		},
		SpecRaw: "data {}",
	}
	d, err := resJSON.Data()

	require.NoError(t, err)

	s, err := New(bCtx)

	require.NoError(t, err)

	// save the blocks to the store
	err = s.Save(core.DataBlocks{d})

	require.NoError(t, err)

	resQuery := fmt.Sprintf(`
			{
				%s(id: "namespace/kind/name")  {
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
