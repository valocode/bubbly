package interval

import (
	"context"
	"fmt"
	"sync"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/hako/durafmt"
	"github.com/zclconf/go-cty/cty"

	"github.com/verifa/bubbly/api"
	"github.com/verifa/bubbly/api/common"
	"github.com/verifa/bubbly/api/core"
	v1 "github.com/verifa/bubbly/api/v1"
	"github.com/verifa/bubbly/env"
	"github.com/verifa/bubbly/events"
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
	Pool           Pool
	WorkerChannels Channels
	Context        ChannelContext
}

type ChannelContext struct {
	Context context.Context
	Cancel  context.CancelFunc
}

type Pool struct {
	mu        sync.RWMutex
	Resources []core.Resource
	Runs      []Run
}

func (p Pool) Remove(i int) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.Runs = append(p.Runs[:i], p.Runs[i+1:]...)
}

type Run struct {
	Resource v1.Run
	interval time.Duration
	Kind     RunKind
	Channel  chan RunAction
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

// ParseResources takes a slice of core.Resource and appends relevant resources to the
// worker's pool of resources. This method can be used to filter resources to only those
// relevant to the worker, in the case that the slice of resources has come from
// sources external to its private pollResources method
func (w *ResourceWorker) ParseResources(bCtx *env.BubblyContext, resources []core.Resource) {
	for _, r := range resources {
		switch r.(type) {
		case *v1.Run:
			runV1 := r.(*v1.Run)
			err := common.DecodeBody(bCtx, runV1.SpecHCL.Body, &runV1.Spec, core.NewResourceContext(cty.NilVal, nil))
			if err != nil {
				continue
			}

			// if true, then the run resource should not be run remotely
			if runV1.Spec.Remote == nil {
				bCtx.Logger.Debug().Str("resource", runV1.String()).Msg("run is of type local and therefore will be ignored by the worker")
				continue
			}

			run := Run{
				Resource: *runV1,
				Kind:     IntervalRun,
			}

			var intervalDuration *durafmt.Durafmt

			intervalDuration, err = durafmt.ParseString(defaultRunInterval)

			// if an interval has been specified, we should try to use that
			if runV1.Spec.Remote.Interval != "" {
				intervalDuration, err = durafmt.ParseString(runV1.Spec.Remote.Interval)

				// if the interval has been specified incorrectly, log this
				// to the user and use the default instead
				if err != nil {
					bCtx.Logger.Error().Err(err).Str("default_interval", defaultRunInterval).Msg("incorrect interval format specified. Using default interval instead")
				}
				run.interval = intervalDuration.Duration()
			} else {
				bCtx.Logger.Debug().Str("run", runV1.String()).Msg("run has no interval specified; treating it as a one-off run")
				run.Kind = OneOffRun
			}

			w.Pool.mu.Lock()
			w.Pool.Runs = append(w.Pool.Runs, run)
			w.Pool.mu.Unlock()
		default:
			bCtx.Logger.Debug().Str("type", string(r.Kind())).Msg("resource worker does not support resources of this kind")
		}
	}
}

// Run establishes the channels and will apply the runs
func (w *ResourceWorker) Run(bCtx *env.BubblyContext) error {
	for i, run := range w.Pool.Runs {
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
		//  on with life
		eg.Go(func() error {
			switch run.Kind {
			case IntervalRun:
				return run.ApplyWithInterval(bCtx, w.WorkerChannels[run.Resource.String()], ctx)
			case OneOffRun:
				err := run.ApplyOneOff(bCtx)

				if err == nil {
					bCtx.Logger.Debug().Msg("run successful. Removing successful run from worker's pool")
					w.Pool.Remove(i)
					return nil
				}

				return err
			}
			return nil
		})
	}

	return nil
}

// ApplyWithInterval will apply the underlying run based on the defined interval
func (r *Run) ApplyWithInterval(bCtx *env.BubblyContext, ch <-chan RunAction, ctx context.Context) error {
	// TODO: fix concurrency bug in provider then enable
	//return nil
	ticker := time.NewTicker(r.interval)
	defer ticker.Stop()
mainloop:
	for {
		select {
		case <-ticker.C:
			resContext := core.NewResourceContext(cty.NilVal, api.NewResource)
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
func (r *Run) ApplyOneOff(bCtx *env.BubblyContext) error {
	// TODO: fix concurrency bug in provider then enable
	//return nil
	bCtx.Logger.Debug().Str("id", r.Resource.String()).Msg("run resource of type OneOffRun identified")
	resContext := core.NewResourceContext(cty.NilVal, api.NewResource)
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
	if err := common.LoadResourceOutput(bCtx, &runResourceOutput); err != nil {
		return fmt.Errorf(
			`failed to store the output of running resource "%s" to the store: %w`,
			r.Resource.String(),
			err,
		)
	}

	return nil
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
