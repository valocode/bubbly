package client

import (
	"fmt"

	"github.com/nats-io/nats.go"

	"github.com/valocode/bubbly/env"
)

// Connect connects to a NATS server.
// It attaches the nats.Conn to
// the NATS Client on a successful connection,
// or an error if it was not possible to establish such.
func (n *NATS) Connect(bCtx *env.BubblyContext) error {
	nc, err := nats.Connect(n.Config.Addr,
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
		return fmt.Errorf(
			`client failed to establish a connection to the NATS
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
func (n *NATS) EncodedConnect(bCtx *env.BubblyContext,
	encoderType string) error {

	// make sure a valid encoder has been provided.
	// Defaults to a Json encoder,
	// since this enables Publications to be encoded and decoded
	if encoderType == "" {
		bCtx.Logger.Debug().Msg("encoder type not provided. Using default")
		encoderType = nats.DEFAULT_ENCODER
	}

	// make sure that the underlying NATS server connection has already been
	// established
	if n.Conn == nil {
		if err := n.Connect(bCtx); err != nil {
			fmt.Errorf("client failed to establish unencoded connection to"+
				" to the NATS Server: %w", err)
		}
	}

	ec, err := nats.NewEncodedConn(n.Conn, string(encoderType))
	if err != nil {
		return fmt.Errorf(
			"client unable to establish encoded connection to the NATS server"+
				": %w",
			err,
		)
	}

	bCtx.Logger.Debug().
		Interface("nats_server", n.Config).
		Msg("client successfully established encoded connected to NATS Server")

	n.EConn = ec

	return nil
}
