// +build integration

package integration

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"os"
	"testing"
	"time"

	"github.com/rs/zerolog"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/valocode/bubbly/api/core"
	"github.com/valocode/bubbly/bubbly"
	"github.com/valocode/bubbly/env"
	"github.com/valocode/bubbly/events"
	integration "github.com/valocode/bubbly/integration/testdata"
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

	t.Run("snyk", func(t *testing.T) {
		// Create a new server route for mocking a Bubbly server response
		bCtx := env.NewBubblyContext()
		bCtx.UpdateLogLevel(zerolog.DebugLevel)
		err := bubbly.Apply(bCtx, "./testdata/snyk")
		assert.NoError(t, err, "Failed to apply resource")

		// test that `bubbly get all` returns valid resources from the apply
		testGet(t, bCtx, []string{"extract/snyk"})
	})

	t.Run("gosec", func(t *testing.T) {
		bCtx := env.NewBubblyContext()
		bCtx.UpdateLogLevel(zerolog.DebugLevel)
		err := bubbly.Apply(bCtx, "./testdata/gosec")
		assert.NoError(t, err, "failed to apply resource")

		testGet(t, bCtx, []string{"extract/gosec"})
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
	//t.Run("criteria", func(t *testing.T) {
	//	bCtx := env.NewBubblyContext()
	//	bCtx.UpdateLogLevel(zerolog.DebugLevel)
	//
	//	err := bubbly.Apply(bCtx, "./testdata/resources/v1/criteria/criteria.bubbly")
	//	assert.NoError(t, err, "Failed to apply resource")
	//})
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
	t.Run("remote_run_with_remote_input", func(t *testing.T) {
		bCtx := env.NewBubblyContext()
		bCtx.UpdateLogLevel(zerolog.DebugLevel)
		id := "run/sonarqube_remote"
		client := &http.Client{}

		// TODO: applying a remote resource which requires remote inputs
		//  will always fail initially. Might be valuable to filter these run
		//  resources and not auto-run them after apply to bubbly
		err := bubbly.Apply(bCtx, "./testdata/resources/v1/run/remote_run_with_remote_input.bubbly")
		require.NoError(t, err, "Failed to apply remote run resource")

		t.Run("json", func(t *testing.T) {
			filePath := "./testdata/sonarqube/sonarqube-example.json"

			req, err := formRemoteInputRequest(
				t,
				fmt.Sprintf("http://127.0.0.1:8111/api/v1/%s", id),
				nil,
				"file",
				"application/json",
				filePath,
			)

			require.NoError(t, err)

			client := &http.Client{}
			res, err := client.Do(req)

			require.NoError(t, err)
			require.NotNil(t, res)

			time.Sleep(2 * time.Second)
			r := getResource(t, bCtx, id)

			latestEvent := r.Events[len(r.Events)-1]

			// if the Worker is enabled with remote running, then we expect it to have
			// run the resource successfully
			require.Equal(t, events.ResourceRunSuccess.String(), latestEvent.Status, latestEvent.Error)
		})

		t.Run("invalid_json", func(t *testing.T) {
			// now send a POST request with invalid payload for the named run resource.
			// That is, a request with content not matching the server-side run resource's
			// input path
			filePath := "./testdata/testautomation/golang/test-report.json"

			req, err := formRemoteInputRequest(
				t,
				fmt.Sprintf("http://127.0.0.1:8111/api/v1/%s", id),
				nil,
				"file",
				"application/json",
				filePath,
			)

			require.NoError(t, err)

			res, err := client.Do(req)

			require.NoError(t, err)
			require.NotNil(t, res)

			time.Sleep(2 * time.Second)
			r := getResource(t, bCtx, id)

			latestEvent := r.Events[len(r.Events)-1]

			// the run should fail, because the file uploaded is not a valid input
			// to the resource
			require.Equal(t, events.ResourceRunFailure.String(), latestEvent.Status, latestEvent.Error)
		})

		t.Run("zip", func(t *testing.T) {

			// now send a POST request with a valid .zip payload
			filePath := "./testdata/sonarqube/sonarqube-example.zip"

			req, err := formRemoteInputRequest(
				t,
				fmt.Sprintf("http://127.0.0.1:8111/api/v1/%s", id),
				nil,
				"file",
				"application/zip",
				filePath,
			)

			require.NoError(t, err)

			res, err := client.Do(req)

			require.NoError(t, err)
			require.NotNil(t, res)

			time.Sleep(2 * time.Second)
			r := getResource(t, bCtx, id)

			latestEvent := r.Events[len(r.Events)-1]

			// the run should succeed, because the .zip file uploaded contains the
			// json file required by the run resource
			require.Equal(t, events.ResourceRunSuccess.String(), latestEvent.Status, latestEvent.Error)
		})
	})
}

// helper function to get a bubbly.Resource, by id, from Bubbly via a graphQL query
func getResource(t *testing.T, bCtx *env.BubblyContext, id string) bubbly.Resource {
	t.Helper()
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
		`, core.ResourceTableName, "id", id,
		core.EventTableName)

	resources, err := bubbly.QueryResources(bCtx, getQuery)

	require.NoError(t, err)
	require.Len(t, resources, 1)

	return resources[0]
}

// helper function which creates a new file upload via HTTP request
func formRemoteInputRequest(t *testing.T, url string, params map[string]string, paramName, contentType string, path string) (*http.Request, error) {
	t.Helper()
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition",
		fmt.Sprintf(`form-data; name="%s"; filename="%s"`,
			paramName, path))
	h.Set("Content-Type", contentType)

	part, err := writer.CreatePart(h)

	if err != nil {
		return nil, err
	}
	_, err = io.Copy(part, file)

	if err != nil {
		return nil, err
	}

	for key, val := range params {
		_ = writer.WriteField(key, val)
	}
	err = writer.Close()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req, err
}
