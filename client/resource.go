package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/graphql-go/graphql"
	"github.com/hashicorp/go-multierror"

	"github.com/valocode/bubbly/agent/component"
	"github.com/valocode/bubbly/api/core"
	"github.com/valocode/bubbly/env"
)

// GetResource uses the bubbly api endpoint to get a resource
// TODO: with some of the new architecture it might be possible for client
//  to return an actual resource, and not just a byte... Don't want to create
//  a merge hell so making a note here
func (c *httpClient) GetResource(bCtx *env.BubblyContext, _ *component.MessageAuth, id string) ([]byte, error) {

	bCtx.Logger.Debug().Str("resource_id", id).Msg("Getting resource from bubbly API.")

	resp, err := c.handleRequest(http.MethodGet, "/resource/"+id, nil)
	if err != nil {
		return nil, fmt.Errorf(`failed to get resource "%s": %w`, id, err)
	}

	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

// PostResource uses the bubbly api endpoint to post a resource
func (c *httpClient) PostResource(bCtx *env.BubblyContext, _ *component.MessageAuth, resource []byte) error {

	_, err := c.handleRequest(http.MethodPost, "/resource", bytes.NewBuffer(resource))
	if err != nil {
		return fmt.Errorf(`failed to post resource: %w`, err)
	}

	return nil
}

// PostResourceToWorker is not supported by the HTTP
func (h *httpClient) PostResourceToWorker(bCtx *env.BubblyContext, _ *component.MessageAuth, data []byte) error {
	return errors.New("unsupported operation for the HTTP client: PostResourceToWorker")
}

// GetResource uses the bubbly NATS client to get a resource from the data
// store.
// Takes a resource ID as input, returns a []byte representing the
// core.ResourceBlock of the resource or an error if
// the client was unable to get the resource.
func (n *natsClient) GetResource(bCtx *env.BubblyContext, auth *component.MessageAuth, resID string) ([]byte,
	error) {

	// for the graphQL query
	resQuery := fmt.Sprintf(`
		{
			%s(id: "%s") {
				name
				kind
				api_version
				metadata
				spec
			}
		}
	`, core.ResourceTableName, resID)

	bCtx.Logger.Debug().
		Str("resource_id", resID).
		Msg("Getting resource from store")

	req := component.Request{
		Subject: component.StoreQuery,
		Data: component.MessageData{
			Auth: auth,
			Data: []byte(resQuery),
		},
	}
	if err := n.request(bCtx, &req); err != nil {
		return nil, fmt.Errorf("failed to get resource from query: %w", err)
	}

	var result graphql.Result
	if err := json.Unmarshal(req.Reply.Data, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response from query to store: %w", err)
	}
	if result.HasErrors() {
		var graphqlErrors error
		for _, qlError := range result.Errors {
			graphqlErrors = multierror.Append(graphqlErrors, qlError)
		}
		return nil, fmt.Errorf("failed to get resource: %w", graphqlErrors)
	}

	if result.Data == nil || result.Data.(map[string]interface{})[core.ResourceTableName] == nil {
		return nil, fmt.Errorf(`no resource found matching ID "%s"`, resID)
	}

	// the result is always []interface{} because this is what the graphql
	// resolver will return
	resourceJSONSlice := result.Data.(map[string]interface{})[core.ResourceTableName].([]interface{})

	// ...which we presume to be of length 1, since resources with identical
	// IDs are upserted. Here we extract the first valid core.ResourceBlockJSON
	// from the slice
	return json.Marshal(resourceJSONSlice[0])
}

// PostResource uses the bubbly natsClient client to publish a resource to the data
// store.
func (n *natsClient) PostResource(bCtx *env.BubblyContext, auth *component.MessageAuth, data []byte) error {
	bCtx.Logger.Debug().
		Str("subject", string(component.StoreUpload)).
		Msg("Posting resource to store")

	req := component.Request{
		Subject: component.StoreUpload,
		Data: component.MessageData{
			Auth: auth,
			Data: data,
		},
	}

	if err := n.request(bCtx, &req); err != nil {
		return fmt.Errorf("failed to post resource: %w", err)
	}

	return nil
}

// PostResource uses the bubbly natsClient client to publish a resource to a worker
// The data is marshalled from a core.DataBlocks
func (n *natsClient) PostResourceToWorker(bCtx *env.BubblyContext, auth *component.MessageAuth, data []byte) error {
	bCtx.Logger.Debug().
		Str("subject", string(component.WorkerPostRunResource)).
		Msg("Posting resource to worker")

	// publish the resource to the worker queue to be picked up and run.
	// we offload responsibility for updating future state of the run resource to
	// the worker that picks it up. What this means is that the worker should
	// update the data store with the success/failure of the run
	if err := n.EConn.Publish(
		string(component.WorkerPostRunResource),
		component.MessageData{
			Auth: auth,
			Data: data,
		}); err != nil {
		return fmt.Errorf("failed to publish run resource to worker: %w", err)
	}

	return nil
}
