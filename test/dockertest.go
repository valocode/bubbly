package test

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/ory/dockertest/v3"
	"github.com/stretchr/testify/require"
)

func RunPostgresDocker(t *testing.T) *dockertest.Resource {
	pool, err := dockertest.NewPool("")
	require.NoErrorf(t, err, "failed to create dockertest pool")

	resource, err := pool.RunWithOptions(
		&dockertest.RunOptions{
			Repository: "postgres",
			Tag:        "13.0",
			Env: []string{
				"POSTGRES_USER=postgres",
				"POSTGRES_PASSWORD=postgres",
				"POSTGRES_DB=bubbly",
			},
		},
	)
	require.NoErrorf(t, err, "failed to start docker")

	t.Cleanup(func() {
		if err := pool.Purge(resource); err != nil {
			t.Fatalf("Could not purge resource: %s", err)
		}
	})

	err = waitUntilDatabaseIsReady(t, pool, resource)
	require.NoErrorf(t, err, "error waiting for database to be ready")

	return resource
}

func waitUntilDatabaseIsReady(t *testing.T, pool *dockertest.Pool, resource *dockertest.Resource) error {
	pgConnStr := fmt.Sprintf("postgresql://postgres:postgres@localhost:%s/bubbly?sslmode=disable",
		resource.GetPort("5432/tcp"))

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
