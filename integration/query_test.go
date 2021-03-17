// +build integration

package integration

import (
	"fmt"
	"testing"

	"github.com/valocode/bubbly/api/core"
	"github.com/valocode/bubbly/bubbly"
	"github.com/valocode/bubbly/client"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
	"github.com/valocode/bubbly/env"
)

func eventQuery(t *testing.T, bCtx *env.BubblyContext) {
	query := fmt.Sprintf(`
			{
				%s(id: "default/extract/gotest")  {
					id
					%s {
						status
						time
					}
				}
			}
		`, core.ResourceTableName, core.EventTableName)

	client, err := client.NewHTTP(bCtx)
	require.NoError(t, err, "failed to establish a bubbly client")

	resp, err := client.Query(bCtx, query)
	require.NoError(t, err, "bubbly client failed to query the bubbly server")

	bCtx.Logger.Debug().RawJSON("response", resp).Msg("received query response from bubbly server")

	require.NotNil(t, string(resp))
}

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

	client, err := client.NewHTTP(bCtx)
	require.NoError(t, err, "failed to establish a bubbly client")

	resp, err := client.Query(bCtx, query)
	require.NoError(t, err, "bubbly client failed to query the bubbly server")

	// too verbose to include full, so just log a snippet
	bCtx.Logger.Debug().Str("response",
		string(resp)[0:100]).Msg("received query response from bubbly server." +
		" Logging a snippet.")

	require.NotNil(t, string(resp))

	eventQuery(t, bCtx)
}
