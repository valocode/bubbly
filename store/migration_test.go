package store

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/verifa/bubbly/config"
	"github.com/verifa/bubbly/env"

	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/verifa/bubbly/api/core"

	"github.com/stretchr/testify/require"
)

const (
	postgresDatabase = "bubbly"
	postgresUser     = "postgres"
	postgresPassword = "secret"

	cockroachDatabase = "defaultdb"
	cockroachUser     = "root"
	cockroachPassword = "admin"
)

func TestApplyMigrationSchemaPostgres(t *testing.T) {
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

	// Get a handle to the bubbly store
	s, err := New(bCtx)
	assert.NoErrorf(t, err, "failed to initialize store")

	// Prepare definitions of db tables
	schema := core.Tables{}
	for _, table := range schema1.Tables {
		schema = append(schema, table)
	}

	// Create tables in db
	err = s.Apply(bCtx, schema)
	assert.NoError(t, err)

	// run tests
	m, err := s.p.GenerateMigration(bCtx, expectedChanges)
	require.NoError(t, err)

	err = s.p.Migrate(m)
	assert.NoError(t, err)
}

// FIXME: CockroachDB support suspended due to lack or required features
func TestApplyMigrationSchemaCockroach(t *testing.T) {

	t.SkipNow()
	// FIXME: this test is out of date, see TestPostgres

	pool, err := dockertest.NewPool("")
	require.NoErrorf(t, err, "failed to create dockertest pool")

	resource, err := pool.RunWithOptions(
		&dockertest.RunOptions{
			Repository: "cockroachdb/cockroach-unstable",
			Tag:        "v21.1.0-alpha.3", // Warning! This is using an alpha version
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
	bCtx.StoreConfig.CockroachDatabase = cockroachDatabase
	bCtx.StoreConfig.CockroachUser = cockroachUser
	bCtx.StoreConfig.CockroachPassword = cockroachPassword

	s, err := New(bCtx)
	assert.NoErrorf(t, err, "failed to initialize store")
	schema := core.Tables{}
	for _, table := range schema1.Tables {
		schema = append(schema, table)
	}
	err = s.Apply(bCtx, schema)
	assert.NoError(t, err)
	// run tests
	m, err := s.p.GenerateMigration(bCtx, expectedChanges)
	require.NoError(t, err)
	err = s.p.Migrate(m)
	assert.NoError(t, err)

	t.Cleanup(func() {
		if err := pool.Purge(resource); err != nil {
			t.Fatalf("Could not purge resource: %s", err)
		}
	})
}
