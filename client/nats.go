package client

import (
	"fmt"
	"time"

	natsd "github.com/nats-io/nats-server/v2/server"
	"github.com/nats-io/nats.go"

	"github.com/valocode/bubbly/agent/component"
	"github.com/valocode/bubbly/config"
	"github.com/valocode/bubbly/env"
)

// newNATS returns a new *client.natsClient bubbly client, using the natsClient server configuration embedded
// within the bubbly context.
func newNATS(bCtx *env.BubblyContext) (*natsClient, error) {
	bCtx.Logger.Debug().
		Interface("client_config", bCtx.AgentConfig.NATSServerConfig).
		Msg("creating a natsClient client")

	sc := bCtx.AgentConfig.NATSServerConfig

	// This configure the natsClient Server using natsd package
	nopts := &natsd.Options{}

	nopts.HTTPPort = sc.HTTPPort
	nopts.Port = sc.Port

	// Create the natsClient Server
	s := natsd.New(nopts)
	c := &natsClient{
		Config: sc,
		Server: s,
	}

	if err := c.connect(bCtx); err != nil {
		return nil, fmt.Errorf("failed to create NATS client when connecting to nats server: %w", err)
	}

	return c, nil
}

type natsClient struct {
	// The configuration of the natsClient server this client should attempt to
	// connect to
	Config *config.NATSServerConfig
	Server *natsd.Server
	Conn   *nats.Conn
	EConn  *nats.EncodedConn
}

// Request requests data on a given subject.
// It differs to Publish in that this requires a response from a subscriber.
// The response is decoded into a new Publication,
// which is returned to the caller.
func (n *natsClient) request(bCtx *env.BubblyContext, req *component.Publication) *component.Publication {
	// Connect to the natsClient Server if a connection has not already been
	// established by this client
	if n.Conn == nil || n.EConn == nil {
		bCtx.Logger.Debug().
			Msg("client is missing required connection to the natsClient Server. " +
				"Attemping to connect")

		if err := n.EncodedConnect(bCtx, req.Encoder); err != nil {
			return &component.Publication{
				Subject: req.Subject,
				Error: fmt.Errorf(
					"failed to connect to the natsClient Server: %w",
					err),
			}
		}
	}

	defer n.Conn.Close()
	defer n.EConn.Close()

	var reply component.Publication

	bCtx.Logger.Debug().
		Interface("nats_client", n.Config).
		Str("subject", string(req.Subject)).
		Msg("sending request")

	// Send a request.
	// The response from the request should always be a []byte,
	// which we can easily decode into our `reply.Data`.
	if err := n.EConn.Request(string(req.Subject), req.Data, &reply.Data,
		defaultNATSClientTimeout*time.Second); err != nil {
		return &component.Publication{
			Subject: req.Subject,
			Error:   fmt.Errorf("failed to make request: %w", err),
		}
	}

	return &reply
}

// publish sends a Publication (https://docs.nats.io/nats-concepts/pubsub)
// over a natsClient server. It returns an error if a connection to the natsClient server
// could not be established or if it was not possible to publish a message on
// the given subject.
func (n *natsClient) publish(bCtx *env.BubblyContext, pub *component.Publication) error {

	// Connect to the natsClient Server if a connection has not already been
	// established by this client.
	if n.Conn == nil || n.EConn == nil {
		bCtx.Logger.Debug().
			Msg("client is missing required connection to the natsClient Server. " +
				"Attemping to connect")

		if err := n.EncodedConnect(bCtx, pub.Encoder); err != nil {
			return fmt.Errorf("failed to connect to the natsClient Server: %w", err)
		}
	}

	defer n.Conn.Close()
	defer n.EConn.Close()

	if err := n.EConn.Publish(string(pub.Subject), pub.Data); err != nil {
		return fmt.Errorf(
			`failed to publish subject "%s" with value "%s": %w`,
			pub.Subject,
			pub.Data,
			err,
		)
	}

	return nil
}

// Connect connects to a natsClient server.
// It attaches the nats.Conn to
// the natsClient Client on a successful connection,
// or an error if it was not possible to establish such.
func (n *natsClient) connect(bCtx *env.BubblyContext) error {
	nc, err := nats.Connect(n.Config.Addr,
		nats.Name("Bubbly natsClient Server"),
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
		return fmt.Errorf(
			`client failed to establish a connection to the natsClient
			server at address "%s": %w`,
			n.Config.Addr,
			err,
		)
	}

	n.Conn = nc

	return nil
}

// ConnectEncoded will wrap an existing Connection with an EncodedConn using
// the appropriate encoder type. If a connection does not exist,
// it will attempt to establish one.
// The reason this is not done at the same time as establishing a normal Conn
// is that the EncodedConn is Publication-scoped,
// since it relies on an encoder type specific to the type of data being
// published.
func (n *natsClient) EncodedConnect(bCtx *env.BubblyContext,
	encoderType string) error {

	// make sure a valid encoder has been provided.
	// Defaults to a Json encoder,
	// since this enables Publications to be encoded and decoded
	if encoderType == "" {
		bCtx.Logger.Debug().Msg("encoder type not provided. Using default")
		encoderType = nats.DEFAULT_ENCODER
	}

	// make sure that the underlying natsClient server connection has already been
	// established
	if n.Conn == nil {
		if err := n.connect(bCtx); err != nil {
			return fmt.Errorf("client failed to establish unencoded connection to"+
				" to the natsClient Server: %w", err)
		}
	}

	ec, err := nats.NewEncodedConn(n.Conn, string(encoderType))
	if err != nil {
		return fmt.Errorf(
			"client unable to establish encoded connection to the natsClient server"+
				": %w",
			err,
		)
	}

	bCtx.Logger.Debug().
		Interface("nats_server", n.Config).
		Msg("client successfully established encoded connected to natsClient Server")

	n.EConn = ec

	return nil
}
