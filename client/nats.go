package client

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/nats-io/nats.go"

	"github.com/valocode/bubbly/agent/component"
	"github.com/valocode/bubbly/env"
)

// newNATS returns a new *client.natsClient bubbly client, using the NATS server configuration embedded
// within the bubbly context.
func newNATS(bCtx *env.BubblyContext) (*natsClient, error) {
	bCtx.Logger.Debug().
		Interface("client_config", bCtx.AgentConfig.NATSServerConfig).
		Msg("creating a NATS client")

	c := &natsClient{}
	nc, err := nats.Connect(bCtx.ClientConfig.NATSAddr,
		nats.Name("Bubbly NATS Server"),
		nats.ErrorHandler(func(nc *nats.Conn, s *nats.Subscription, err error) {
			if s != nil {
				bCtx.Logger.Error().
					Err(err).
					Str("subject", s.Subject).
					Str("queue", s.Queue).
					Msg("Async error")
			} else {
				bCtx.Logger.Error().
					Err(err).
					Msg("Async error outside subscription")
			}
		}))

	if err != nil {
		return nil, fmt.Errorf(
			`client failed to establish a connection to the NATS server at address "%s": %w`,
			bCtx.ClientConfig.NATSAddr,
			err,
		)
	}

	c.EConn, err = nats.NewEncodedConn(nc, nats.DEFAULT_ENCODER)
	if err != nil {
		return nil, fmt.Errorf("failed to create encoded connection to NATS server: %w", err)
	}

	return c, nil
}

type natsClient struct {
	EConn *nats.EncodedConn
}

func (n *natsClient) Close() {
	n.EConn.Close()
}

// Request publishes a Request-Reply message on a given subject.
// It differs to Publish in that this requires a response from a subscriber.
// The Reply is added to the given request.
func (n *natsClient) request(bCtx *env.BubblyContext, req *component.Request) error {

	bCtx.Logger.Debug().
		Str("subject", string(req.Subject)).
		Msg("sending request")

	// Make sure the pointer where we will put the reply is initialized
	// otherwise nats will fail when decoding
	if req.Reply == nil {
		req.Reply = &component.Reply{}
	}

	// Send a request.
	// The response from the request should always be a []byte,
	// which we can easily decode into our `reply.Data`.
	var reply []byte
	if err := n.EConn.Request(string(req.Subject), req.Data, &reply,
		defaultNATSClientTimeout*time.Second); err != nil {
		return fmt.Errorf("failed to make request: %w", err)
	}

	if err := json.Unmarshal(reply, req.Reply); err != nil {
		return fmt.Errorf("failed to decode reply from request: %w", err)
	}

	if req.Reply.Error != "" {
		return fmt.Errorf("received error from request: %s", req.Reply.Error)
	}
	return nil
}
