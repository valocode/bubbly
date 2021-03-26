package client

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/valocode/bubbly/agent/component"
	"github.com/valocode/bubbly/env"
)

// PostSchema uses the bubbly api to post a schema
func (c *httpClient) PostSchema(bCtx *env.BubblyContext, schema []byte) error {

	_, err := c.handleRequest(http.MethodPost, "/schema", bytes.NewBuffer(schema))
	return err
}

func (n *natsClient) PostSchema(bCtx *env.BubblyContext, schema []byte) error {
	bCtx.Logger.Debug().
		Str("subject", string(component.StorePostSchema)).
		Msg("Posting schema to data store")

	req := component.Request{
		Subject: component.StorePostSchema,
		Data:    schema,
	}
	// reply is a Publication received from a bubbly store
	if err := n.request(bCtx, &req); err != nil {
		return fmt.Errorf("failed to post schema: %w", err)
	}

	return nil
}
