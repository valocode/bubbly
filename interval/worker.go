package interval

import (
	"context"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/hako/durafmt"
	"github.com/verifa/bubbly/api"
	"github.com/verifa/bubbly/api/common"
	"github.com/verifa/bubbly/api/core"
	v1 "github.com/verifa/bubbly/api/v1"
	"github.com/verifa/bubbly/env"
	"github.com/zclconf/go-cty/cty"
)

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
	Resources    []core.Resource
	PipelineRuns []PipelineRun
}

type PipelineRun struct {
	Resource v1.PipelineRun
	interval time.Duration
	Channel  chan PipelineAction
}

type PipelinePool struct {
	pipelineRunLoops []PipelineRunLoop
}

type PipelineRunLoop struct {
	id       string
	interval time.Duration
	resource core.Resource
}

type ActionType string
type PipelineAction struct {
	Action        ActionType
	ResourceBlock *core.ResourceBlock
}

const (
	UpdatePipeline ActionType = "update"
	StopPipeline   ActionType = "stop"
)

// Run parses the supplied resources into pipelineRuns, establishes the channels and will apply
// the pipeline runs
func (w *ResourceWorker) Run(bCtx *env.BubblyContext, resources []core.Resource) error {
	for _, pr := range resources {
		prV1 := pr.(*v1.PipelineRun)
		err := common.DecodeBody(bCtx, prV1.SpecHCL.Body, &prV1.Spec, core.NewResourceContext(cty.NilVal, nil))
		if err != nil {
			continue
		}

		intervalDuration, err := durafmt.ParseString(prV1.Spec.Interval)
		if err != nil {
			bCtx.Logger.Error().Err(err).Msg("malformed date")
			continue
		}

		pipelineRun := PipelineRun{
			Resource: *prV1,
			interval: intervalDuration.Duration(),
		}

		w.Pool.PipelineRuns = append(w.Pool.PipelineRuns, pipelineRun)

		eg := new(errgroup.Group)
		var ctx, cancel = context.WithCancel(context.Background())
		w.Context = ChannelContext{
			Context: ctx,
			Cancel:  cancel,
		}

		pipelineRun.Channel = make(chan PipelineAction)
		w.WorkerChannels = make(Channels)
		w.WorkerChannels[pipelineRun.Resource.String()] = pipelineRun.Channel
		eg.Go(func() error {
			return pipelineRun.ApplyWithInterval(bCtx, w.WorkerChannels[pipelineRun.Resource.String()], ctx)
		})
	}

	return nil
}

// ApplyWithInterval will apply the underlying pipelineRun based on the defined interval
func (pr *PipelineRun) ApplyWithInterval(bCtx *env.BubblyContext, ch <-chan PipelineAction, ctx context.Context) error {
	ticker := time.NewTicker(pr.interval)
	defer ticker.Stop()
mainloop:
	for {
		select {
		case <-ticker.C:
			resContext := core.NewResourceContext(cty.NilVal, api.NewResource)
			output := pr.Resource.Apply(bCtx, resContext)
			if output.Error != nil {
				bCtx.Logger.Error().Err(output.Error).Msg("error applying pipeline_run")
			}
			// TODO send output to NATS
		case msg := <-ch:
			switch msg.Action {
			case UpdatePipeline:
				// The underlying resource has been changed, update it
				err := pr.update(*msg.ResourceBlock)
				if err != nil {
					return err
				}
			case StopPipeline:
				pr.DeleteChannel()
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

func (pr *PipelineRun) update(resBlock core.ResourceBlock) error {
	newRes, err := api.NewResource(&resBlock)
	if err != nil {
		return err
	}
	// update the PipelineRunLoop's resource
	prV1 := newRes.(*v1.PipelineRun)
	pr.Resource = *prV1
	return nil
}
