package interval

import (
	"sync"
	"testing"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"

	"github.com/valocode/bubbly/api"
	"github.com/valocode/bubbly/api/core"
	v1 "github.com/valocode/bubbly/api/v1"
	"github.com/valocode/bubbly/env"
	"github.com/valocode/bubbly/parser"
	"github.com/valocode/bubbly/server"
)

func newTestWorker(t *testing.T) ResourceWorker {
	t.Helper()
	return ResourceWorker{
		Pools: Pools{
			OneOff: Pool{
				mu:        sync.Mutex{},
				Resources: nil,
				Runs:      make(map[uuid.UUID]Run),
			},
			Interval: IntervalPool{},
		},
		WorkerChannels: nil,
		Context:        ChannelContext{},
	}
}

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

	worker := newTestWorker(t)

	for _, r := range resources {
		worker.ParseResource(bCtx, r, server.RemoteInput{})
	}

	require.Len(t, worker.Pools.OneOff.Runs, 0)
	require.Len(t, worker.Pools.Interval.Pool.Runs, 0)
}

// TestWorkerParseLocalRun tests that a worker ignores a run resource that has not
// defined a remote{} block
func TestWorkerParseLocalRun(t *testing.T) {
	bCtx := env.NewBubblyContext()
	bCtx.UpdateLogLevel(zerolog.DebugLevel)

	resources := parseBubblyFile(t, bCtx, "./testdata/run_local.bubbly")

	worker := newTestWorker(t)

	for _, r := range resources {
		worker.ParseResource(bCtx, r, server.RemoteInput{})
	}

	require.Len(t, worker.Pools.OneOff.Runs, 0)
	require.Len(t, worker.Pools.Interval.Pool.Runs, 0)
}

// TestWorkerParseRemoteIntervalRun tests that a worker identifies a
// remote run resource with valid interval and, while this type of run is
// disabled, appends it to the worker's one-off pool instead
func TestWorkerParseRemoteIntervalRun(t *testing.T) {
	bCtx := env.NewBubblyContext()
	bCtx.UpdateLogLevel(zerolog.DebugLevel)

	resources := parseBubblyFile(t, bCtx, "./testdata/run_remote_interval.bubbly")

	worker := newTestWorker(t)

	for _, r := range resources {
		err := worker.ParseResource(bCtx, r, server.RemoteInput{})
		require.Nil(t, err)
	}

	require.Len(t, worker.Pools.OneOff.Runs, 1)
	require.Len(t, worker.Pools.Interval.Pool.Runs, 0)
}

// TestWorkerParseRemoteOneOffRun tests that a worker identifies a remote run
// resource and adds it to its pool of resources.
// The Run has no interval specified, and must therefore be of type OneOffRun
func TestWorkerParseRemoteOneOffRun(t *testing.T) {

	bCtx := env.NewBubblyContext()
	bCtx.UpdateLogLevel(zerolog.DebugLevel)

	resources := parseBubblyFile(t, bCtx, "./testdata/run_remote_one_off.bubbly")

	worker := newTestWorker(t)

	for _, r := range resources {
		worker.ParseResource(bCtx, r, server.RemoteInput{})
	}

	require.Len(t, worker.Pools.OneOff.Runs, 1)
	require.Len(t, worker.Pools.Interval.Pool.Runs, 0)
}

// TestWorkerPoolAddRemove tests a worker's ability to add and remove a
// Run from its pool
func TestWorkerPoolAddRemove(t *testing.T) {

	id := uuid.New()

	run := Run{
		UUID:     id,
		Resource: v1.Run{},
	}

	worker := newTestWorker(t)

	worker.Pools.OneOff.Append(run)
	require.Len(t, worker.Pools.OneOff.Runs, 1)

	worker.Pools.OneOff.Remove(run)
	require.Len(t, worker.Pools.OneOff.Runs, 0)

}
