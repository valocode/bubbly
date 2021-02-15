package natsserver

import (
	"context"
	"fmt"
	"time"

	natsd "github.com/nats-io/nats-server/server"

	"github.com/verifa/bubbly/agent/component"
	"github.com/verifa/bubbly/config"
	"github.com/verifa/bubbly/env"
)

const (
	defaultNATSServerWait       = 10
	defaultLogStatisticsTimeout = 60
)

var _ component.Component = (*NATSServer)(nil)

type NATSServer struct {
	*component.ComponentCore
	Config *config.NATSServerConfig
	Server *natsd.Server
}

// New instances a new NATSServer.
// Aside from standard ComponentCore and Config set up,
// it populates a *natsd.Server instance and attaches it to the NATSServer
// for use at runtime.
func New(bCtx *env.BubblyContext) *NATSServer {
	bCtx.Logger.Debug().Msg("creating a NATS server")
	// This configure the NATS Server using natsd package
	nopts := &natsd.Options{}

	nopts.HTTPPort = bCtx.AgentConfig.NATSServerConfig.HTTPPort
	nopts.Port = bCtx.AgentConfig.NATSServerConfig.Port

	if e := bCtx.Logger.Debug(); e.Enabled() {
		nopts.Debug = true
	}

	// TODO: Add support for nopts.Trace
	// nopts.Trace = true

	// Create the NATS Server
	ns := natsd.New(nopts)

	// enable the logger for the nats.Server. Currently, I cannot see a
	// good way of using the zerolog.Logger.
	ns.ConfigureLogger()

	return &NATSServer{
		ComponentCore: &component.ComponentCore{
			Type: component.NATSServerComponent,
		},
		Config: bCtx.AgentConfig.NATSServerConfig,
		Server: ns,
	}
}

// Run runs a NATS server in single-server agent deployment
func (n *NATSServer) Run(bCtx *env.BubblyContext, ctx context.Context) error {
	bCtx.Logger.Debug().Str(
		"component",
		string(n.Type),
	).Msg("running component")

	nSubs, err := n.BulkSubscribe(bCtx)

	if err != nil {
		return fmt.Errorf("error during bulk subscription: %w", err)
	}

	n.Subscriptions = nSubs
	bCtx.Logger.Debug().Str("component", string(n.Type)).Interface("subscriptions", n.Subscriptions).Msg("component is listening for subscriptions")

	go n.Server.Start()

	go n.logStats(bCtx)

	// Wait for the NATS Server to be able to accept connections before
	// connecting and setting up the other bubbly agent components
	if !n.Server.ReadyForConnections(defaultNATSServerWait * time.Second) {
		return fmt.Errorf("NATS server took too long to allow client connection")
	}

	return n.Listen(ctx)
}

func (n *NATSServer) logStats(bCtx *env.BubblyContext) {
	for {
		bCtx.Logger.Debug().
			Int("num_routes", n.Server.NumRoutes()).
			Int("num_clients", n.Server.NumClients()).
			Uint32("num_subscriptions", n.Server.NumSubscriptions()).
			Msg("NATS server statistics")

		time.Sleep(defaultLogStatisticsTimeout * time.Second)
	}
}
