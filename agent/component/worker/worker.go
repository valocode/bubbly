package worker

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/valocode/bubbly/agent/component"
	"github.com/valocode/bubbly/env"
	"github.com/valocode/bubbly/interval"
)

func New(bCtx *env.BubblyContext) *Worker {
	w := &Worker{
		ComponentCore: &component.ComponentCore{
			Type:                 component.WorkerComponent,
			DesiredSubscriptions: nil,
		},
		ResourceWorker: &interval.ResourceWorker{
			Pools: interval.Pools{
				OneOff: interval.Pool{
					Runs: make(map[uuid.UUID]interval.Run),
				},
				Interval: interval.IntervalPool{
					IsRunning: false,
					Pool: interval.Pool{
						Runs: make(map[uuid.UUID]interval.Run),
					},
				},
			},
			WorkerChannels: nil,
			Context:        interval.ChannelContext{},
		},
	}

	w.DesiredSubscriptions = w.defaultSubscriptions()

	bCtx.Logger.Debug().Msg("successfully initialised a worker")

	return w
}

// TODO: describe more about the Worker
type Worker struct {
	*component.ComponentCore
	ResourceWorker *interval.ResourceWorker
}

// Run runs the interval.ResourceWorker
func (w *Worker) Run(bCtx *env.BubblyContext, agentContext context.Context) error {
	bCtx.Logger.Debug().
		Str(
			"component",
			string(w.Type)).
		Msg("running component")

	nSubs, err := w.BulkSubscribe(bCtx)

	if err != nil {
		return fmt.Errorf("error during bulk subscription: %w", err)
	}

	w.Subscriptions = nSubs
	bCtx.Logger.Debug().
		Str("component", string(w.Type)).
		Interface("subscriptions", w.Subscriptions).
		Msg("component is listening for subscriptions")

	return w.Listen(agentContext)
}

// a list of DesiredSubscriptions that the data store attempts to subscribe to
func (w *Worker) defaultSubscriptions() component.DesiredSubscriptions {
	return component.DesiredSubscriptions{
		component.DesiredSubscription{
			Subject: component.WorkerPostRunResource,
			Queue:   component.WorkerQueue,
			Reply:   false,
			Handler: w.postRunResourceHandler,
		},
	}
}
