package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/graphql-go/graphql"
	"github.com/hashicorp/go-multierror"
	"github.com/nats-io/nats.go"

	"github.com/valocode/bubbly/agent/component"
	"github.com/valocode/bubbly/api/core"
	"github.com/valocode/bubbly/env"
)

// GetResource uses the bubbly api endpoint to get a resource
// TODO: with some of the new architecture it might be possible for client
//  to return an actual resource, and not just a byte... Don't want to create
//  a merge hell so making a note here
func (h *httpClient) GetResource(bCtx *env.BubblyContext, id string) ([]byte, error) {

	bCtx.Logger.Debug().Str("resource_id", id).Msg("Getting resource from bubbly API.")

	resp, err := h.handleResponse(
		http.Get(fmt.Sprintf("%s/api/v1/resource/%s", h.HostURL, id)),
	)
	if err != nil {
		return nil, fmt.Errorf(`failed to get resource "%s": %w`, id, err)
	}

	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

// PostResource uses the bubbly api endpoint to post a resource
func (h *httpClient) PostResource(bCtx *env.BubblyContext, resource []byte) error {

	_, err := h.handleResponse(
		http.Post(fmt.Sprintf("%s/api/v1/resource", h.HostURL),
			"application/json", bytes.NewBuffer(resource)),
	)

	if err != nil {
		return fmt.Errorf(`failed to post resource: %w`, err)
	}

	return nil
}

// PostResourceToWorker is not supported by the HTTP
func (h *httpClient) PostResourceToWorker(bCtx *env.BubblyContext, data []byte) error {
	return errors.New("unsupported operation for the HTTP client: PostResourceToWorker")
}

// GetResource uses the bubbly NATS client to get a resource from the data
// store.
// Takes a resource ID as input, returns a []byte representing the
// core.ResourceBlock of the resource or an error if
// the client was unable to get the resource.
func (n *natsClient) GetResource(bCtx *env.BubblyContext, resID string) ([]byte,
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
		Interface("nats_client", n.Config).
		Str("resource_query", resQuery).
		Msg("Getting resource from store")

	request := component.Publication{
		Subject: component.StoreQuery,
		Data:    []byte(resQuery),
		Encoder: nats.DEFAULT_ENCODER,
	}

	reply := n.request(bCtx, &request)

	if reply.Error != nil {
		return nil, fmt.Errorf(
			`failed to get resource from query %s: %w`,
			resQuery,
			reply.Error,
		)
	}

	var result graphql.Result

	err := json.Unmarshal(reply.Data, &result)

	if err != nil {
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

	// the result is a []interface{}
	resourceJSONSlice := result.Data.(map[string]interface{})[core.ResourceTableName].([]interface{})

	// ...which we presume to be of length 1, since resources with identical
	// IDs are upserted. Here we extract the first valid core.ResourceBlockJSON
	// from the slice
	resJSON, err := json.Marshal(resourceJSONSlice[0])

	return resJSON, nil
}

// PostResource uses the bubbly natsClient client to publish a resource to the data
// store.
func (n *natsClient) PostResource(bCtx *env.BubblyContext, data []byte) error {
	bCtx.Logger.Debug().
		Interface("client", n.Config).
		Msg("Posting resource to store")

	request := component.Publication{
		Subject: component.StorePostResource,
		Data:    data,
		Encoder: nats.DEFAULT_ENCODER,
	}

	reply := n.request(bCtx, &request)

	if reply.Error != nil {
		return fmt.Errorf(
			`failed to post resource: %w`,
			reply.Error,
		)
	}

	return nil
}

// PostResource uses the bubbly natsClient client to publish a resource to a worker
// The data is marshalled from a core.DataBlocks
func (n *natsClient) PostResourceToWorker(bCtx *env.BubblyContext, data []byte) error {
	bCtx.Logger.Debug().
		Interface("client", n.Config).
		Interface("resource", data).
		Str("subject", string(component.WorkerPostRunResource)).
		Msg("Posting resource to worker")

	request := component.Publication{
		Subject: component.WorkerPostRunResource,
		Data:    data,
		Encoder: nats.DEFAULT_ENCODER,
	}

	// publish the resource to the worker queue to be picked up and run.
	// we offload responsibility for updating future state of the run resource to
	// the worker that picks it up. What this means is that the worker should
	// update the data store with the success/failure of the run

	err := n.publish(bCtx, &request)

	if err != nil {
		return fmt.Errorf("failed to publish run resource to worker: %w", err)
	}

	return nil
}
