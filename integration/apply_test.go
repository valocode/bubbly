// +build integration

package integration

import (
	"testing"

	"github.com/rs/zerolog"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/verifa/bubbly/bubbly"
	"github.com/verifa/bubbly/env"
)

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
	})

	// Subtest
	t.Run("rest2", func(t *testing.T) {

		bCtx := env.NewBubblyContext()
		bCtx.UpdateLogLevel(zerolog.DebugLevel)

		err := bubbly.Apply(bCtx, "./testdata/rest2")
		assert.NoError(t, err, "Failed to apply resource")
	})

	// Subtest
	t.Run("dynamic_source", func(t *testing.T) {
		bCtx := env.NewBubblyContext()
		bCtx.UpdateLogLevel(zerolog.DebugLevel)

		err := bubbly.Apply(bCtx, "./testdata/resources/v1/extract/multisource/multisource.bubbly")
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

func TestApplyTaskRun(t *testing.T) {
	// Subtest
	t.Run("task_run_sonarqube_extract", func(t *testing.T) {
		bCtx := env.NewBubblyContext()
		bCtx.UpdateLogLevel(zerolog.DebugLevel)

		err := bubbly.Apply(bCtx, "./testdata/resources/v1/taskrun/extract_sonarqube.bubbly")
		assert.NoError(t, err, "Failed to apply resource")
	})
}
