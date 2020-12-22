// +build integration

package integration

import (
	"testing"

	"github.com/verifa/bubbly/bubbly"
	"github.com/verifa/bubbly/client"

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

	query := `{
		test_case(status: "pass") {
			name
			status
		}
	}`

	client, err := client.New(bCtx)
	require.NoError(t, err, "failed to establish a bubbly client")

	resp, err := client.Query(bCtx, query)
	require.NoError(t, err, "bubbly client failed to query the bubbly server")

	bCtx.Logger.Debug().RawJSON("response", resp).Msg("received query response from bubbly server")

	require.NotNil(t, string(resp))
}

// TestHCLQuery tests that, given an injection of golang test data,
// bubbly is able to apply a Query resource to retrieve a subset of the
// data
func TestHCLQuery(t *testing.T) {

	bCtx := env.NewBubblyContext()
	bCtx.UpdateLogLevel(zerolog.DebugLevel)

	// inject initial bubbly test data
	err := bubbly.Apply(bCtx, "./testdata/testautomation/golang/pipeline.bubbly")
	require.NoError(t, err, "failed to apply golang pipeline")

	err = bubbly.Apply(bCtx, "./testdata/testautomation/golang/query.bubbly")
	require.NoError(t, err, "failed to apply query resource")
}
