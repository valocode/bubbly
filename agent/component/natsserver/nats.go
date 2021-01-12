package natsserver

import (
	"fmt"

	natsd "github.com/nats-io/nats-server/server"
	"github.com/nats-io/nats.go"

	"github.com/verifa/bubbly/agent/component"
	"github.com/verifa/bubbly/config"
	"github.com/verifa/bubbly/env"
)

var _ component.NATSServer = (*NATSServer)(nil)

type NATSServer struct {
	*component.ComponentCore

	Config     *config.NATSServerConfig
	Server     *natsd.Server
	ServerConn *nats.Conn
}

// New creates a new NATS server
func (n *NATSServer) New(bCtx *env.BubblyContext) error {
	bCtx.Logger.Debug().Msg("creating a NATS server")
	// This configure the NATS Server using natsd package
	nopts := &natsd.Options{}

	nopts.HTTPPort = n.Config.HTTPPort
	nopts.Port = n.Config.Port

	// Create the NATS Server
	ns := natsd.New(nopts)

	n.Server = ns

	return nil
}

// Connect connects to a NATS server and attaches the nats.Conn
// for use in 1:N goroutines established by the other agent components
func (n *NATSServer) Connect(bCtx *env.BubblyContext) error {
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
			`failed to establish a connection to the NATS
			server at address "%s": %w`,
			n.Config.Addr,
			err,
		)
	}

	n.ServerConn = nc

	return nil
}

// Run runs a NATS server
func (n *NATSServer) Run(bCtx *env.BubblyContext) error {
	bCtx.Logger.Debug().Str(
		"component",
		string(n.Type),
	).Msg("running component")

	if err := n.BulkSubscribe(bCtx); err != nil {
		return fmt.Errorf("error during bulk subscription: %w", err)
	}

	n.Server.Start()

	return nil
}
