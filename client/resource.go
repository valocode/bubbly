package client

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/nats-io/nats.go"

	"github.com/verifa/bubbly/agent/component"
	"github.com/verifa/bubbly/env"
)

// GetResource uses the bubbly api endpoint to get a resource
// TODO: with some of the new architecture it might be possible for client
//  to return an actual resource, and not just a byte... Don't want to create
//  a merge hell so making a note here
func (h *HTTP) GetResource(bCtx *env.BubblyContext, id string) ([]byte, error) {

	bCtx.Logger.Debug().Str("resource_id", id).Msg("Getting resource from bubbly API.")

	resp, err := h.handleResponse(
		http.Get(fmt.Sprintf("%s/api/resource/%s", h.HostURL, id)),
	)
	if err != nil {
		return nil, fmt.Errorf(`failed to get resource "%s": %w`, id, err)
	}

	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

// PostResource uses the bubbly api endpoint to post a resource
func (h *HTTP) PostResource(bCtx *env.BubblyContext, resource []byte) error {

	_, err := h.handleResponse(
		http.Post(fmt.Sprintf("%s/api/resource", h.HostURL),
			"application/json", bytes.NewBuffer(resource)),
	)

	if err != nil {
		return fmt.Errorf(`failed to post resource: %w`, err)
	}

	return nil
}

// GetResource uses the bubbly NATS client to get a resource from the data
// store.
// Returns a []byte representation of the requested resource or an error if
// the client was unable to get the resource.
func (n *NATS) GetResource(bCtx *env.BubblyContext, resQuery string) ([]byte,
	error) {
	bCtx.Logger.Debug().
		Interface("nats_client", n.Config).
		Str("resource_query", resQuery).
		Msg("Getting resource from store")

	request := component.Publication{
		Subject: component.StoreGetResource,
		Data:    []byte(resQuery),
		Encoder: nats.DEFAULT_ENCODER,
	}

	// reply is a Publication received from a bubbly store
	reply := n.Request(bCtx, &request)

	if reply.Error != nil {
		return nil, fmt.Errorf(
			`failed to get resource from query %s: %w`,
			resQuery,
			reply.Error,
		)
	}

	return reply.Data, nil
}

// PostResource uses the bubbly NATS client to publish a resource to the data
// store.
func (n *NATS) PostResource(bCtx *env.BubblyContext, data []byte) error {
	bCtx.Logger.Debug().
		Interface("client", n.Config).
		Msg("Posting resource to store")

	request := component.Publication{
		Subject: component.StorePostResource,
		Data:    data,
		Encoder: nats.DEFAULT_ENCODER,
	}

	reply := n.Request(bCtx, &request)

	if reply.Error != nil {
		return fmt.Errorf(
			`failed to post resource: %w`,
			reply.Error,
		)
	}

	return nil
}

// Upload uses the bubbly NATS client to upload arbitrary data to be saved
// into the data store.
func (n *NATS) Upload(bCtx *env.BubblyContext, data []byte) error {
	bCtx.Logger.Debug().
		Interface("nats_client", n.Config).
		Msg("Uploading data to the data store")

	request := component.Publication{
		Subject: component.StoreUpload,
		Data:    data,
		Encoder: nats.DEFAULT_ENCODER,
	}

	// reply is a Publication received from a bubbly store
	reply := n.Request(bCtx, &request)

	if reply.Error != nil {
		return fmt.Errorf("failed during upload: %w", reply.Error)
	}

	return nil
}
