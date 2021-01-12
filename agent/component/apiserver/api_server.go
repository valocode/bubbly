package apiserver

import (
	"fmt"

	"github.com/verifa/bubbly/agent/component"
	"github.com/verifa/bubbly/env"
	"github.com/verifa/bubbly/server"
)

var _ component.APIServer = (*APIServer)(nil)

type APIServer struct {
	*component.ComponentCore
}

// Run runs the bubbly API Server.
func (a *APIServer) Run(bCtx *env.BubblyContext) error {
	bCtx.Logger.Debug().Str(
		"component",
		string(a.Type),
	).Msg("running component")

	if err := a.BulkSubscribe(bCtx); err != nil {
		return fmt.Errorf("error during bulk subscription: %w", err)
	}

	err := server.ListenAndServe(bCtx)

	if err != nil {
		return fmt.Errorf("error while serving: %w", err)
	}

	return nil
}
