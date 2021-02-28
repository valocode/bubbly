// +build integration

package integration

import (
	"fmt"
	"testing"
	"time"

	"github.com/rs/zerolog"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/verifa/bubbly/api/core"
	"github.com/verifa/bubbly/bubbly"
	"github.com/verifa/bubbly/env"
	"github.com/verifa/bubbly/events"
	integration "github.com/verifa/bubbly/integration/testdata"
)

func testGet(t *testing.T, bCtx *env.BubblyContext, args []string) {
	t.Helper()
	integration.TestBubblyCmd(t, bCtx, "get", args)
}

// TestApply simply validates that a given directory containing bubbly
// configuration including a pipeline_run will result in a POST of data to
// the bubbly server.
// See client/load_test.go for actual evaluation of the loading using
// the gofight package.
func TestApply(t *testing.T) {

	// Subtest
	t.Run("sonarqube", func(t *testing.T) {
		// Create a new server route for mocking a Bubbly server response
		bCtx := env.NewBubblyContext()
		bCtx.UpdateLogLevel(zerolog.DebugLevel)

		err := bubbly.Apply(bCtx, "./testdata/sonarqube")
		assert.NoError(t, err, "Failed to apply resource")

		// test that `bubbly get all` returns valid resources from the apply
		testGet(t, bCtx, []string{"all"})
	})

	// Subtest
	t.Run("spdx-licenses.bubbly", func(t *testing.T) {
		bCtx := env.NewBubblyContext()
		bCtx.UpdateLogLevel(zerolog.DebugLevel)

		err := bubbly.Apply(bCtx, "./testdata/resources/v1/extract/spdx-licenses.bubbly")
		assert.NoError(t, err, "Failed to apply resource")
	})

	// Subtest
	t.Run("query", func(t *testing.T) {
		bCtx := env.NewBubblyContext()
		bCtx.UpdateLogLevel(zerolog.DebugLevel)

		// inject initial bubbly test data
		err := bubbly.Apply(bCtx, "./testdata/testautomation/golang/pipeline.bubbly")
		require.NoError(t, err, "failed to apply golang pipeline")
		err = bubbly.Apply(bCtx, "./testdata/resources/v1/query/query.bubbly")
		assert.NoError(t, err, "Failed to apply resource")
	})

	// TODO: need to create a criteria test for something real...
	// This fails because the data does not exist.

	// Subtest
	// t.Run("criteria", func(t *testing.T) {
	// 	bCtx := env.NewBubblyContext()
	// 	bCtx.UpdateLogLevel(zerolog.DebugLevel)

	// 	err := bubbly.Apply(bCtx, "./testdata/resources/v1/criteria/criteria.bubbly")
	// 	assert.NoError(t, err, "Failed to apply resource")
	// })
}

func TestApplyRun(t *testing.T) {
	// Subtest
	t.Run("sonarqube_run", func(t *testing.T) {
		bCtx := env.NewBubblyContext()
		bCtx.UpdateLogLevel(zerolog.DebugLevel)

		err := bubbly.Apply(bCtx, "./testdata/resources/v1/run/resources.bubbly")
		assert.NoError(t, err, "Failed to apply resource")
	})
	t.Run("remote_run", func(t *testing.T) {
		bCtx := env.NewBubblyContext()
		bCtx.UpdateLogLevel(zerolog.DebugLevel)

		err := bubbly.Apply(bCtx, "./testdata/resources/v1/run/remote_resources.bubbly")
		require.NoError(t, err, "Failed to apply remote resource")

		getQuery := fmt.Sprintf(`
			{
				%s(%s: "%s") {
					id
					%s {
						status
						time
						error
					}
				}
			}
		`, core.ResourceTableName, "id", "run/licenses_remote",
			core.EventTableName)

		// wait for the run resource to actually be saved to the store
		// TODO: consider a better mechanism for this "wait"
		time.Sleep(2 * time.Second)

		resources, err := bubbly.QueryResources(bCtx, getQuery)

		require.NoError(t, err)

		r := resources[0]

		latestEvent := r.Events[len(r.Events)-1]

		// if the Worker is enabled with remote running, then we expect it to have
		// run the resource successfully
		require.Equal(t, events.ResourceRunSuccess.String(), latestEvent.Status, latestEvent.Error)
	})
}
