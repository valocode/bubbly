package worker

import (
	"context"
	"fmt"

	"github.com/nats-io/nats.go"

	"github.com/valocode/bubbly/agent/component"
	"github.com/valocode/bubbly/env"
	"github.com/valocode/bubbly/interval"
)

const (
	defaultPollTimeout = 60
)

func New(bCtx *env.BubblyContext) *Worker {
	w := &Worker{
		ComponentCore: &component.ComponentCore{
			Type: component.WorkerComponent,
			NATSServer: &component.NATS{
				Config: bCtx.AgentConfig.NATSServerConfig,
			},
			DesiredSubscriptions: nil,
		},
		ResourceWorker: &interval.ResourceWorker{},
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
			Handler: w.PostRunResourceHandler,
			Encoder: nats.DEFAULT_ENCODER,
		},
	}
}
