package apiserver

import (
	"context"
	"fmt"

	"github.com/verifa/bubbly/agent/component"
	"github.com/verifa/bubbly/env"
	"github.com/verifa/bubbly/server"
)

var _ component.APIServer = (*APIServer)(nil)

type APIServer struct {
	*component.ComponentCore
	Server *server.Server
}

type Result struct {
	Message string
	Error   error
}

// Run runs the bubbly API Server.
func (a *APIServer) Run(bCtx *env.BubblyContext, agentContext context.Context) error {
	bCtx.Logger.Debug().Str(
		"component",
		string(a.Type),
	).Msg("running component")

	nSubs, err := a.BulkSubscribe(bCtx)

	if err != nil {
		return fmt.Errorf("error during bulk subscription: %w", err)
	}

	a.Subscriptions = nSubs

	bCtx.Logger.Debug().Str("component", string(a.Type)).Interface("subscriptions", a.Subscriptions).Msg("component is listening for subscriptions")

	ch := make(chan error, 1)
	defer close(ch)

	// run the actual api server in a separate goroutine, but track its
	// performance using a channel
	go a.run(bCtx, ch)

	select {
	// if the api server fails, error
	case err := <-ch:
		return fmt.Errorf("error while running API server: %w", err)
	// if another agent component fails, error
	case <-agentContext.Done():
		return agentContext.Err()
	}
}

func (a *APIServer) run(bCtx *env.BubblyContext, ch chan error) {
	if err := a.Server.ListenAndServe(bCtx); err != nil {
		bCtx.Logger.Debug().Err(err).Msg("API Server finished due to error")
		ch <- fmt.Errorf("API Server failed: %w", err)
	}
}
