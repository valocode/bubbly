package store

import (
	"database/sql"
	"fmt"
	"testing"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/stretchr/testify/require"

	"github.com/valocode/bubbly/api/core"
	"github.com/valocode/bubbly/config"
	"github.com/valocode/bubbly/env"
)

// TestEventTrigger tests that saving a new resource to the store results in
// both _resource and _event population using the event store trigger
func TestEventTrigger(t *testing.T) {
	pool, err := dockertest.NewPool("")
	require.NoErrorf(t, err, "failed to connect to Docker")

	// TODO: probably best to extract image version to a var, and maybe even read from env var, if available?
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
	// TODO: the hard limit for expiration probably is best set via env var, if available

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

	s, err := New(bCtx)
	require.NoError(t, err)

	// save the blocks to the store
	err = s.Save(core.DataBlocks{resJSON.Data()})
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
