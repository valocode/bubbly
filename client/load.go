package client

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/valocode/bubbly/agent/component"
	"github.com/valocode/bubbly/env"
)

// Load takes data blocks and saves them to the bubbly server
func (c *httpClient) Load(bCtx *env.BubblyContext, _ *component.MessageAuth, data []byte) error {

	_, err := c.handleRequest(http.MethodPost, "/upload", bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("failed to save data: %w", err)
	}

	return nil
}

// Upload uses the bubbly NATS client to upload arbitrary data to be saved
// into the data store.
func (n *natsClient) Load(bCtx *env.BubblyContext, auth *component.MessageAuth, data []byte) error {
	bCtx.Logger.Debug().
		Msg("Uploading data to the data store")

	req := component.Request{
		Subject: component.StoreUpload,
		Data: component.MessageData{
			Auth: auth,
			Data: data,
		},
	}
	// reply is a Publication received from a bubbly store
	if err := n.request(bCtx, &req); err != nil {
		return fmt.Errorf("failed during upload: %w", err)
	}

	return nil
}
