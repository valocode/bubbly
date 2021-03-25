package client

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/nats-io/nats.go"

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
		Interface("nats_client", n.Config).
		Str("subject", string(component.StorePostResource)).
		Msg("Posting schema to data store")

	request := component.Publication{
		Subject: component.StorePostSchema,
		Data:    schema,
		Encoder: nats.DEFAULT_ENCODER,
	}

	// reply is a Publication received from a bubbly store
	reply := n.request(bCtx, &request)

	if reply.Error != nil {
		return fmt.Errorf("failed during schema post: %w", reply.Error)
	}

	return nil
}
