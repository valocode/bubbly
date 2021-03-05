package worker

import (
	"encoding/json"
	"fmt"

	"github.com/nats-io/nats.go"

	"github.com/verifa/bubbly/api"
	"github.com/verifa/bubbly/api/core"
	"github.com/verifa/bubbly/env"
)

// PostRunResourceHandler receives a resourceBlockJSON matching a remote run
// resource from the store (specifically, the remote run trigger)
// It adds the resource to the worker's pool and then runs it.
func (w *Worker) PostRunResourceHandler(bCtx *env.BubblyContext, m *nats.Msg) error {
	var resJSON core.ResourceBlockJSON

	err := json.Unmarshal(m.Data, &resJSON)

	bCtx.Logger.Debug().
		Interface("subscription", m.Sub).
		Str("res", resJSON.SpecRaw).
		Str("component", string(w.Type)).
		Msg("processing resource")

	if err != nil {
		return fmt.Errorf("failed to unmarshal resource: %w", err)
	}

	resBlock, err := resJSON.ResourceBlock()

	if err != nil {
		return fmt.Errorf("failed to form resource from block: %w", err)
	}

	res, err := api.NewResource(&resBlock)

	w.ResourceWorker.ParseResources(bCtx, []core.Resource{res})

	err = w.ResourceWorker.Run(bCtx)
	if err != nil {
		return fmt.Errorf("interval worker failure: %w", err)
	}

	return nil
}
