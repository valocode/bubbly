package bubbly

import (
	"net/http"
	"testing"

	"github.com/rs/zerolog"

	"github.com/stretchr/testify/assert"

	"github.com/verifa/bubbly/env"
	"gopkg.in/h2non/gock.v1"
)

// TestApply simply validates that a given directory containing bubbly
// configuration including a pipeline_run will result in a POST of data to
// the bubbly server.
// See client/load_test.go for actual evaluation of the loading using
// the gofight package.
func TestApply(t *testing.T) {

	defer gock.Off()

	host := "localhost"
	port := "8080"
	hostURL := host + ":" + port

	// Subtest
	t.Run("sonarqube", func(t *testing.T) {
		// Create a new server route for mocking a Bubbly server response
		gock.New(hostURL).
			Post("/alpha1/upload").
			Reply(http.StatusOK).
			JSON(map[string]interface{}{"status": "uploaded"})

		bCtx := env.NewBubblyContext()
		bCtx.UpdateLogLevel(zerolog.DebugLevel)

		err := Apply(bCtx, "./testdata/sonarqube")

		assert.NoError(t, err, "Failed to apply resource")
	})
}

func TestApplyTaskRun(t *testing.T) {
	// Subtest
	t.Run("task_run_sonarqube_extract", func(t *testing.T) {
		bCtx := env.NewBubblyContext()
		bCtx.UpdateLogLevel(zerolog.DebugLevel)

		err := Apply(bCtx, "./testdata/resources/v1/taskrun/extract_sonarqube.bubbly")

		assert.NoError(t, err, "Failed to apply resource")
	})
}
