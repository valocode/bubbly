package client

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/nats-io/nats.go"

	"github.com/verifa/bubbly/agent/component"
	"github.com/verifa/bubbly/env"
)

// PostSchema uses the bubbly api to post a schema
func (c *HTTP) PostSchema(bCtx *env.BubblyContext, schema []byte) error {

	_, err := c.handleResponse(
		http.Post(fmt.Sprintf("%s/api/v1/schema", c.HostURL), "application/json", bytes.NewBuffer(schema)),
	)
	return err
}

func (n *NATS) PostSchema(bCtx *env.BubblyContext, schema []byte) error {
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
	reply := n.Request(bCtx, &request)

	if reply.Error != nil {
		return fmt.Errorf("failed during schema post: %w", reply.Error)
	}

	return nil
}
