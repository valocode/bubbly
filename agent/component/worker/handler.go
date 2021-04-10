package worker

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"

	"github.com/graphql-go/graphql"

	"github.com/valocode/bubbly/agent/component"
	"github.com/valocode/bubbly/api"
	"github.com/valocode/bubbly/api/core"
	"github.com/valocode/bubbly/env"
	"github.com/valocode/bubbly/server"
)

// postRunResourceHandler receives a server.WorkerRun containing a run name and
// (optionally) input data needed to run the resource.
// It gets the resource from the store, adds it to the worker's one-off pool,
// saves any remote inputs to the Worker's local filesystem and then runs it.
func (w *Worker) postRunResourceHandler(bCtx *env.BubblyContext, subject string, reply string, data component.MessageData) (interface{}, error) {
	var wr server.WorkerRun
	if err := json.Unmarshal(data.Data, &wr); err != nil {
		return nil, fmt.Errorf("failed to unmarshal data into WorkerRun: %w", err)
	}
	// Validate WorkerRun
	if reflect.DeepEqual(server.WorkerRun{}, wr) {
		return nil, errors.New("worker run cannot be empty")
	}
	if wr.Name == "" {
		return nil, errors.New("worker run resource name required")
	}

	bCtx.Logger.Debug().
		Str("subject", subject).
		Str("component", string(w.Type)).
		Str("resource", string(wr.Name)).
		Msg("processing request to run resource")

	res, err := w.getRunResource(bCtx, data.Auth, wr.Name)
	if err != nil {
		return nil, fmt.Errorf("interval worker failed to get resource: %w", err)
	}

	// parse the resource and add it to the worker's pool
	err = w.ResourceWorker.ParseResource(bCtx, res, wr.RemoteInput)
	if err != nil {
		return nil, fmt.Errorf("interval worker failed to parse resource: %w", err)
	}

	// TODO: Support Interval Runs
	err = w.ResourceWorker.RunOneOffRuns(bCtx, data.Auth)
	if err != nil {
		return nil, fmt.Errorf("interval worker failure: %w", err)
	}

	return nil, nil
}

// sends a NATS publication querying the Bubbly Store for a named run resource.
// Returns the fetched core.Resource or and error if unsuccessful.
func (w *Worker) getRunResource(bCtx *env.BubblyContext, auth *component.MessageAuth, name string) (core.Resource, error) {
	// We want to fetch all resource of type pipeline run from the data
	// store. So form a graphql query representing such
	resQuery := fmt.Sprintf(`
		{
			%s(id: "%s/%s") {
				name
				kind
				api_version
				metadata
				spec
			}
		}
	`, core.ResourceTableName, core.RunResourceKind, name)

	// embed the query into a Request
	req := component.Request{
		Subject: component.StoreQuery,
		Data: component.MessageData{
			Auth: auth,
			Data: []byte(resQuery),
		},
	}

	// reply is a Publication received from a bubbly store
	if err := w.Request(bCtx, &req); err != nil {
		return nil, fmt.Errorf(
			`failed to get resource "%s" from store: %w`,
			name,
			err,
		)
	}
	if req.Reply.Error != "" {
		return nil, fmt.Errorf(
			`failed to get resource "%s" from store: %s`,
			name,
			req.Reply.Error,
		)
	}

	var result graphql.Result
	if err := json.Unmarshal(req.Reply.Data, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal get resp: %w", err)
	}

	resources := result.Data.(map[string]interface{})[core.ResourceTableName]
	if resources == nil {
		return nil, errors.New("no resource found")
	}

	// extract the resource (singular) from the graphql.Result response
	resourceBytes, err := json.Marshal(resources.([]interface{})[0])
	if err != nil {
		return nil, fmt.Errorf("failed to marshal graphql query response: %w", err)
	}

	var resJSON core.ResourceBlockJSON
	if err := json.Unmarshal(resourceBytes, &resJSON); err != nil {
		return nil, fmt.Errorf("failed to unmarshal resource: %w", err)
	}

	resBlock, err := resJSON.ResourceBlock()
	if err != nil {
		return nil, fmt.Errorf("failed to form resource block from JSON: %w", err)
	}

	res, err := api.NewResource(&resBlock)
	if err != nil {
		return nil, fmt.Errorf("failed to form resource from block: %w", err)
	}

	return res, nil
}
