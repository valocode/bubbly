package bubbly

import (
	"io/ioutil"
	"net/http"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/zclconf/go-cty/cty"

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

	// Subtest
	t.Run("sonarqube", func(t *testing.T) {
		// Create a new server route for mocking a Bubbly server response
		bCtx := env.NewBubblyContext()
		bCtx.UpdateLogLevel(zerolog.DebugLevel)

		gock.New(bCtx.ServerConfig.HostURL()).
			Post("/alpha1/upload").
			Reply(http.StatusOK).
			JSON(map[string]interface{}{"status": "uploaded"})

		err := Apply(bCtx, "./testdata/sonarqube")

		assert.NoError(t, err, "Failed to apply resource")
	})

	// Subtest
	t.Run("rest2", func(t *testing.T) {

		bCtx := env.NewBubblyContext()
		bCtx.UpdateLogLevel(zerolog.DebugLevel)

		gock.New(bCtx.ServerConfig.HostURL()).
			Post("/alpha1/upload").
			Reply(http.StatusOK).
			JSON(map[string]interface{}{"status": "uploaded"})

		username := "the_mouse"

		password, perr := ioutil.ReadFile("./testdata/rest2/the_mouse_password")
		require.Nil(t, perr, "test fixture: password containing file")

		gock.New("https://api.cloud84.dev/users").
			Get("/").
			BasicAuth(username, string(password)).
			MatchParam("per_page", "12").
			Reply(http.StatusOK).
			JSON(map[string]interface{}{"answer": 42})

		err := Apply(bCtx, "./testdata/rest2")

		assert.NoError(t, err, "Failed to apply resource")
	})

	// Subtest
	t.Run("dynamic_source", func(t *testing.T) {
		bCtx := env.NewBubblyContext()
		bCtx.UpdateLogLevel(zerolog.DebugLevel)

		// Mock the top-level list (of three items),
		// and details for each of the three items.
		var items = []struct {
			url  string
			file string
		}{
			{
				"https://raw.githubusercontent.com/spdx/license-list-data/master/json/licenses.json",
				"./testdata/resources/v1/extract/multisource/licenses.json",
			},
			{
				"http://spdx.org/licenses/0BSD.json",
				"./testdata/resources/v1/extract/multisource/0BSD.json",
			},
			{
				"http://spdx.org/licenses/MPL-1.1.json",
				"./testdata/resources/v1/extract/multisource/MPL-1.1.json",
			},
			{
				"http://spdx.org/licenses/zlib-acknowledgement.json",
				"./testdata/resources/v1/extract/multisource/zlib-acknowledgement.json",
			},
		}

		for _, i := range items {
			gock.New(i.url).Get("/").Reply(http.StatusOK).File(i.file)
		}

		err := Apply(bCtx, "./testdata/resources/v1/extract/multisource/multisource.bubbly")

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

func TestApplyQuery(t *testing.T) {
	// Subtest
	t.Run("apply basic query", func(t *testing.T) {
		bCtx := env.NewBubblyContext()
		bCtx.UpdateLogLevel(zerolog.DebugLevel)

		gock.New(bCtx.ServerConfig.HostURL()).
			Post("/api/graphql").
			Reply(http.StatusOK).
			JSON(`{"data":{"test_run":{"name":"run 1","repo_version_id":0,"test_set":[{"name":"set 1","test_case":[{"ID":1,"name":"case 1.1","status":"PASS","test_set_id":1},{"ID":2,"name":"case 1.2","status":"PASS","test_set_id":1},{"ID":3,"name":"case 1.3","status":"FAIL","test_set_id":1}]},{"name":"set 2","test_case":[{"ID":4,"name":"case 2.1","status":"FAIL","test_set_id":2},{"ID":5,"name":"case 2.2","status":"FAIL","test_set_id":2}]}]}}}`)

		err := ApplyQueries(bCtx, "./testdata/resources/v1/query/query.bubbly")

		assert.NoError(t, err, "Failed to apply resource")
	})
}

func TestApplyCriteria(t *testing.T) {
	bCtx := env.NewBubblyContext()
	bCtx.UpdateLogLevel(zerolog.DebugLevel)

	tcs := []struct {
		desc     string
		input    string
		expected map[string]cty.Value
	}{
		{
			desc:  "all conditions true",
			input: `{"data":{"quote_run":{"name":"run 1","quote_set":[{"quote":"Men in rage strike those that wish them best.","book":"Othello","speaker":"Iago"},{"quote":"Look like the innocent flower, but be the serpent under ‘t.","book":"Macbeth","speaker":"Lady Macbeth"},{"quote":"Brevity is the soul of wit.","book":"Hamlet","speaker":"Polonius"}]}}}`,
			expected: map[string]cty.Value{
				"shakespeare": cty.True,
			},
		},
		{
			desc:  "all conditions false",
			input: `{"data":{"quote_run":{"name":"run 1","quote_set":[{"quote":"You speak an infinite deal of nothing.","book":"The Merchant of Venice","speaker":"Bassanio"},{"quote":"Sometimes it's about playing a poor hand well.","book":"The Goldfinch","speaker":"unknown"},{"quote":"The portrait of a blinking idiot","book":"The Merchant of Venice","speaker":"Aragon"}]}}}`,
			expected: map[string]cty.Value{
				"shakespeare": cty.False,
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.desc, func(t *testing.T) {

			gock.New(bCtx.ServerConfig.HostURL()).
				Post("/api/graphql").
				Reply(http.StatusOK).
				JSON(tc.input)

			criterionResults, err := ApplyCriterion(bCtx, filepath.FromSlash("testdata/resources/v1/criteria/criteria.bubbly"))

			require.NoError(t, err, "Failed to apply criterion")

			require.Equal(t, tc.expected, criterionResults, "one or more criteria returned an unexpected output value")
		})
	}
}
