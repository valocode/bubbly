package client

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/nats-io/nats.go"
	"github.com/valocode/bubbly/agent/component"
	"github.com/valocode/bubbly/api/core"
	"github.com/valocode/bubbly/env"
)

// Load takes the output from a load resource and POSTs it to the Bubbly
// server.
// Returns an error if loading was unsuccessful
func (c *httpClient) Load(bCtx *env.BubblyContext, data core.DataBlocks) error {

	dBytes, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal data for loading: %w", err)
	}

	// req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/alpha1/upload", c.HostURL), bytes.NewBuffer(jsonReq))
	_, err = c.handleResponse(
		c.Client.Post(fmt.Sprintf("%s/api/v1/upload", c.HostURL), "application/json", bytes.NewBuffer(dBytes)),
	)

	if err != nil {
		return fmt.Errorf(`failed to post resource: %w`, err)
	}

	return nil
}

// Upload uses the bubbly NATS client to upload arbitrary data to be saved
// into the data store.
func (n *natsClient) Load(bCtx *env.BubblyContext, data core.DataBlocks) error {
	bCtx.Logger.Debug().
		Interface("nats_client", n.Config).
		Msg("Uploading data to the data store")

	dBytes, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal data for loading: %w", err)
	}

	request := component.Publication{
		Subject: component.StoreUpload,
		Data:    dBytes,
		Encoder: nats.DEFAULT_ENCODER,
	}

	// reply is a Publication received from a bubbly store
	reply := n.request(bCtx, &request)
	if reply.Error != nil {
		return fmt.Errorf("failed during upload: %w", reply.Error)
	}

	return nil
}
