package interval

import (
	"context"
	"errors"
	"fmt"
	"os"
	"reflect"
	"sync"
	"time"

	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"

	"github.com/zclconf/go-cty/cty"

	"github.com/valocode/bubbly/agent/component"
	"github.com/valocode/bubbly/api"
	"github.com/valocode/bubbly/api/common"
	"github.com/valocode/bubbly/api/core"
	v1 "github.com/valocode/bubbly/api/v1"
	"github.com/valocode/bubbly/env"
	"github.com/valocode/bubbly/events"
	"github.com/valocode/bubbly/server"
)

const (
	defaultRunInterval = "300s"
)

type RunKind int

const (
	IntervalRun RunKind = iota
	OneOffRun
)

func (r RunKind) String() string {
	return [...]string{
		"IntervalRun",
		"OneOffRun"}[r]
}

type ResourceWorker struct {
	Pools          Pools
	WorkerChannels Channels
	Context        ChannelContext
}

type Pools struct {
	OneOff   Pool
	Interval IntervalPool
}

type IntervalPool struct {
	IsRunning bool
	Pool      Pool
}

type ChannelContext struct {
	Context context.Context
	Cancel  context.CancelFunc
}

type Pool struct {
	mu        sync.Mutex
	Resources []core.Resource
	Runs      map[uuid.UUID]Run
}

// delete a Run from a Pool.
// Mutex locking/unlocking should be handled externally
func (p *Pool) Remove(r Run) {
	delete(p.Runs, r.UUID)
}

// Append a Run to a Pool
// Mutex locking/unlocking is handled internally
func (p *Pool) Append(r Run) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.Runs[r.UUID] = r
}

type Run struct {
	UUID        uuid.UUID
	Resource    v1.Run
	interval    time.Duration
	Kind        RunKind
	Channel     chan RunAction
	RemoteInput RemoteInput
}

// RemoteInput represents the location of any input data
// required to run the remote resource
type RemoteInput struct {
	Dir string
}

type RunPool struct {
	runLoops []RunLoop
}

type RunLoop struct {
	id       string
	interval time.Duration
	resource core.Resource
}

type ActionType string
type RunAction struct {
	Action        ActionType
	ResourceBlock *core.ResourceBlock
}

const (
	UpdateRun ActionType = "update"
	StopRun   ActionType = "stop"
)

// RunIntervalRuns establishes the channels and will run runs of type IntervalRun
// TODO: how to handle additions to the Pool during runtime?
func (w *ResourceWorker) RunIntervalRuns(bCtx *env.BubblyContext, auth *component.MessageAuth) error {
	iPool := w.Pools.Interval
	for _, run := range iPool.Pool.Runs {
		eg := new(errgroup.Group)
		var ctx, cancel = context.WithCancel(context.Background())
		w.Context = ChannelContext{
			Context: ctx,
			Cancel:  cancel,
		}

		run.Channel = make(chan RunAction)
		w.WorkerChannels = make(Channels)
		w.WorkerChannels[run.Resource.String()] = run.Channel
		// TODO: the logic for this errgroup is slight flawed considered our new
		//  worker redesign. Somewhere here we probably want to notify NATS on
		//  run failure, retry/purge the run from the worker pool and then move
		//  on
		eg.Go(func() error {
			return run.ApplyWithInterval(bCtx, w.WorkerChannels[run.Resource.String()], ctx, auth)
		})
	}

	return nil
}

// RunOneOffRuns runs all resources within the resource worker's OneOff Pool.
// That is, all of its one-off run resources
func (w *ResourceWorker) RunOneOffRuns(bCtx *env.BubblyContext, auth *component.MessageAuth) error {
	w.Pools.OneOff.mu.Lock()
	bCtx.Logger.Debug().Int("pool", len(w.Pools.OneOff.Runs)).Msg("number of one-off runs to run")
	for _, run := range w.Pools.OneOff.Runs {
		// run has been triggered from a POST to /api/v1/run/:name and
		// therefore should be run from the root tmp directory associated
		// with the remote input
		if run.RemoteInput.Dir != "" {
			if err := os.Chdir(run.RemoteInput.Dir); err != nil {
				w.Pools.OneOff.mu.Unlock()
				return fmt.Errorf("unable to chdir to %v: %v", run.RemoteInput.Dir, err)
			}
		}

		dir, _ := os.Getwd()

		bCtx.Logger.Debug().Str("dir", dir).Msg("running one-off run resource")

		err := run.ApplyOneOff(bCtx, auth)

		// regardless of outcome, purge the one-off resource from the worker
		// pool to prevent run build up
		bCtx.Logger.Debug().
			Int("pool_size", len(w.Pools.OneOff.Runs)).
			Str("name", w.Pools.OneOff.Runs[run.UUID].Resource.ResourceName).
			Msg("one-off run complete. Removing run from worker's one-off pool")

		w.Pools.OneOff.Remove(run)

		bCtx.Logger.Debug().
			Int("pool_size", len(w.Pools.OneOff.Runs)).
			Msg("run removed")

		if err != nil {
			bCtx.Logger.Error().
				Err(err).
				Str("run", run.Resource.ResourceName).
				Msg("failed to run one-off run resource")
		} else {
			bCtx.Logger.Debug().
				Str("run", run.Resource.ResourceName).
				Msg("ran one-off run resource successfully")
		}

		// remove the now-redundant temp directory from the Worker's local filesystem
		if err := os.RemoveAll(run.RemoteInput.Dir); err != nil {
			bCtx.Logger.Error().Err(err).Msg("failed to purge remote input temporary directory")
		}
	}

	w.Pools.OneOff.mu.Unlock()

	return nil
}

// ParseResource writes data in the server.RemoteInput to the local
// filesystem, then decodes and validates the associated resource,
// adding it to the Worker's resource pool on successful validation
func (w *ResourceWorker) ParseResource(bCtx *env.BubblyContext, r core.Resource, input server.RemoteInput) error {
	var i RemoteInput

	// if this parser has been triggered from a POST to /api/v1/run/:name,
	// then we
	// parse the content sent alongside the request into RemoteInput
	if !reflect.DeepEqual(input, server.RemoteInput{}) {
		dir, err := w.parseInputData(input)

		if err != nil {
			return err
		}

		i = RemoteInput{
			Dir: dir,
		}
	}

	switch r.(type) {
	case *v1.Run:
		runV1 := r.(*v1.Run)
		err := common.DecodeBody(bCtx, runV1.SpecHCL.Body, &runV1.Spec, core.NewResourceContext(cty.NilVal, api.NewResource, nil))
		if err != nil {
			bCtx.Logger.Error().Str("resource", runV1.String()).Msg("worker failed to decode resource")
			return fmt.Errorf("worker failed to decode resource: %w", err)
		}

		// if true, then the run resource should not be run remotely
		if runV1.Spec.Remote == nil {
			bCtx.Logger.Debug().Str("resource", runV1.String()).Msg("run is of type local and therefore will be ignored by the worker")
			return nil
		}

		run := Run{
			Resource:    *runV1,
			Kind:        IntervalRun,
			RemoteInput: i,
		}

		// TODO: disabling interval runs until dedicated time can be invested
		//  to support them in a stable manner.

		// var intervalDuration *durafmt.Durafmt

		// intervalDuration, err = durafmt.ParseString(defaultRunInterval)

		// if an interval has been specified, we should try to use that
		if runV1.Spec.Remote.Interval != "" {
			bCtx.Logger.Info().Str("run", runV1.String()).Msg("resource worker does not currently support remote runs of type interval; treating it as one-off run instead")
			// intervalDuration, err = durafmt.ParseString(runV1.Spec.Remote.Interval)
			//
			// // if the interval has been specified incorrectly, log this
			// // to the user and use the default instead
			// if err != nil {
			// 	bCtx.Logger.Error().Err(err).Str("default_interval", defaultRunInterval).Msg("incorrect interval format specified. Using default interval instead")
			// }
			// run.interval = intervalDuration.Duration()
			// w.Pools.Interval.Pool.Append(run)
		} else {
			bCtx.Logger.Debug().Str("run", runV1.String()).Msg("run has no interval specified; treating it as a one-off run")
		}
		run.Kind = OneOffRun

		run.UUID, err = uuid.NewRandom()

		if err != nil {
			return fmt.Errorf("failed to create UUID for run: %w", err)
		}

		w.Pools.OneOff.Append(run)

	default:
		bCtx.Logger.Debug().Str("type", string(r.Kind())).Msg("resource worker does not support resources of this kind")
	}
	return nil
}

// parse the []byte from a NATS publication and returns the path to the
// temporary directory containing the written file(s)
func (w *ResourceWorker) parseInputData(input server.RemoteInput) (string, error) {
	var (
		err error
		dir string
	)
	switch input.Format {
	case "json":
		// write the bytes to a json file
		dir, err = createJSONFromBytes(input.Filename, input.Data)

		return dir, nil
	case "zip":
		// write the bytes as an unzipped collection of files
		dir, err = unzipFromBytes(input.Filename, input.Data)

		if err != nil {
			return "", fmt.Errorf("failed to unzip from bytes: %w", err)
		}
	}
	if err != nil {
		return "", fmt.Errorf("failed to process bytes: %w", err)
	}

	if dir == "" {
		return "", errors.New("failed to process bytes into valid directory")
	}

	return dir, nil
}

// ApplyWithInterval will apply the underlying run based on the defined interval
// TODO: enable once functionality is stabilised
func (r *Run) ApplyWithInterval(bCtx *env.BubblyContext, ch <-chan RunAction, ctx context.Context, auth *component.MessageAuth) error {
	ticker := time.NewTicker(r.interval)
	defer ticker.Stop()
mainloop:
	for {
		select {
		case <-ticker.C:
			resContext := core.NewResourceContext(cty.NilVal, api.NewResource, auth)
			output := r.Resource.Apply(bCtx, resContext)
			if output.Error != nil {
				bCtx.Logger.Error().Err(output.Error).Msg("error applying run")
			}
			// TODO send output to NATS
		case msg := <-ch:
			switch msg.Action {
			case UpdateRun:
				// The underlying resource has been changed, update it
				err := r.update(*msg.ResourceBlock)
				if err != nil {
					return err
				}
			case StopRun:
				r.DeleteChannel()
				break mainloop
			default:
				// Unsupported message type -> Skip
			}
		case <-ctx.Done():
			// receiving shutdown signal
			break mainloop
		}
	}

	return nil
}

// ApplyOneOff is used to apply remote Run resources that lack a specified
// interval. The resource is applied only once.
func (r *Run) ApplyOneOff(bCtx *env.BubblyContext, auth *component.MessageAuth) error {
	bCtx.Logger.Debug().Str("id", r.Resource.String()).Msg("run resource of type OneOffRun identified")
	resContext := core.NewResourceContext(cty.NilVal, api.NewResource, auth)
	subResourceOutput := r.Resource.Apply(bCtx, resContext)

	runResourceOutput := core.ResourceOutput{
		ID:     r.Resource.String(),
		Status: events.ResourceRunSuccess,
		Error:  nil,
	}

	// if any child resource of the run resource has failed, then
	// mark the run resource has failed and attach the error message
	if subResourceOutput.Error != nil {
		runResourceOutput.Status = events.ResourceRunFailure
		runResourceOutput.Error = fmt.Errorf(`failed to run resource "%s": %w`, r.Resource.String(), subResourceOutput.Error)
	}

	// TODO: better consider what action to perform if the worker fails to
	//  run the current run resource

	// load the output of the run resource to the event store
	if err := common.LoadResourceOutput(bCtx, &runResourceOutput, auth); err != nil {
		return fmt.Errorf(
			`failed to store the output of running resource "%s" to the store: %w`,
			r.Resource.String(),
			err,
		)
	}

	return runResourceOutput.Error
}

func (r *Run) update(resBlock core.ResourceBlock) error {
	newRes, err := api.NewResource(&resBlock)
	if err != nil {
		return err
	}
	// update the Run's resource
	runV1 := newRes.(*v1.Run)
	r.Resource = *runV1
	return nil
}
