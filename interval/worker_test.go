package interval

import (
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"

	"github.com/valocode/bubbly/api"
	"github.com/valocode/bubbly/api/core"
	"github.com/valocode/bubbly/env"
	"github.com/valocode/bubbly/parser"
)

// parseBubblyFile parses a .bubbly file and returns a slice of core.Resource
func parseBubblyFile(t *testing.T, bCtx *env.BubblyContext, filename string) []core.Resource {
	resParser := api.NewParserType()
	require.NoError(t, parser.ParseFilename(bCtx, filename, resParser))

	require.NoError(t, resParser.CreateResources(bCtx))

	return resParser.Resources
}

// TestWorkerParseInvalidResource tests that the worker only appends correct
// resource kinds to its internal pool
func TestWorkerParseInvalidResource(t *testing.T) {
	bCtx := env.NewBubblyContext()
	bCtx.UpdateLogLevel(zerolog.DebugLevel)

	resources := parseBubblyFile(t, bCtx, "./testdata/extract.bubbly")

	worker := ResourceWorker{}

	worker.ParseResources(bCtx, resources)

	require.Nil(t, worker.Pool.Runs)
}

// TestWorkerParseLocalRun tests that a worker ignores a run resource that has not
// defined a remote{} block
func TestWorkerParseLocalRun(t *testing.T) {
	bCtx := env.NewBubblyContext()
	bCtx.UpdateLogLevel(zerolog.DebugLevel)

	resources := parseBubblyFile(t, bCtx, "./testdata/run_local.bubbly")

	worker := ResourceWorker{}

	worker.ParseResources(bCtx, resources)
	require.Nil(t, worker.Pool.Resources)
}

// TestWorkerParseRemoteIntervalRun tests that a worker identifies a remote run resource
// and adds it to its pool of resources.
// The Run has an interval specified, and must therefore be of type IntervalRun
func TestWorkerParseRemoteIntervalRun(t *testing.T) {
	bCtx := env.NewBubblyContext()
	bCtx.UpdateLogLevel(zerolog.DebugLevel)

	resources := parseBubblyFile(t, bCtx, "./testdata/run_remote_interval.bubbly")

	worker := ResourceWorker{}

	worker.ParseResources(bCtx, resources)

	require.NotNil(t, worker.Pool.Runs)
	require.Equal(t, IntervalRun, worker.Pool.Runs[0].Kind)
}

// TestWorkerParseRemoteOneOffRun tests that a worker identifies a remote run
// resource and adds it to its pool of resources.
// The Run has no interval specified, and must therefore be of type OneOffRun
func TestWorkerParseRemoteOneOffRun(t *testing.T) {
	bCtx := env.NewBubblyContext()
	bCtx.UpdateLogLevel(zerolog.DebugLevel)

	resources := parseBubblyFile(t, bCtx, "./testdata/run_remote_one_off.bubbly")

	worker := ResourceWorker{}

	worker.ParseResources(bCtx, resources)

	require.NotNil(t, worker.Pool.Runs)
	require.Equal(t, OneOffRun, worker.Pool.Runs[0].Kind)
}
