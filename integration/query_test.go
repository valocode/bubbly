// +build integration

package integration

import (
	"testing"

	"github.com/verifa/bubbly/bubbly"
	"github.com/verifa/bubbly/client"
	"github.com/verifa/bubbly/server"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
	"github.com/verifa/bubbly/env"
)

// TestClientQuery tests that, given an injection of golang test data,
// the bubbly client is able to use its Query method to query the data
func TestClientQuery(t *testing.T) {

	bCtx := env.NewBubblyContext()
	bCtx.UpdateLogLevel(zerolog.DebugLevel)

	// inject initial bubbly test data
	err := bubbly.Apply(bCtx, "./testdata/testautomation/golang/pipeline.bubbly")

	require.NoError(t, err, "failed to apply golang pipeline")

	// get the store instanced by TestMain
	s := server.GetStore()

	query := `{
		test_case(status: "pass") {
			name
			status
		}
	}`
	// query the store using a graphql query. This is useful as a means
	// of verifying that the query itself is valid over the given testdata
	_, err = s.Query(query)

	require.NoError(t, err)

	// instance a new bubbly client
	c, err := client.New(bCtx)

	require.NoError(t, err, "failed to establish a bubbly client")

	// using the same graphql query, we now validate that the bubbly
	// client processes the query and queries the store correctly (knowing
	// that the query itself is proven valid)
	actual, err := c.Query(bCtx, query)
	require.NoError(t, err, "bubbly client failed to query the bubbly server")

	bCtx.Logger.Debug().RawJSON("response", actual).Msg("received query response from bubbly server")

	require.NotNil(t, string(actual))
}

// TestHCLQuery tests that, given an injection of golang test data,
// bubbly is able to apply a Query resource to retrieve a subset of the
// data
func TestHCLQuery(t *testing.T) {

	bCtx := env.NewBubblyContext()
	bCtx.UpdateLogLevel(zerolog.DebugLevel)

	err := bubbly.Apply(bCtx, "./testdata/testautomation/golang/pipeline.bubbly")

	require.NoError(t, err, "failed to apply golang pipeline")

	err = bubbly.ApplyQueries(bCtx, "./testdata/testautomation/golang/query.bubbly")

	require.NoError(t, err, "failed to apply query resource")
}
