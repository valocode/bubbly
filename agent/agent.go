package agent

import (
	"context"
	"fmt"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/valocode/bubbly/agent/component"
	"github.com/valocode/bubbly/agent/component/apiserver"
	"github.com/valocode/bubbly/agent/component/datastore"
	"github.com/valocode/bubbly/agent/component/natsserver"
	"github.com/valocode/bubbly/agent/component/worker"
	"github.com/valocode/bubbly/config"
	"github.com/valocode/bubbly/env"
)

const (
	defaultNATSServerWait = 10
)

// Components represents the components of an Agent.
type Components struct {
	// Which components have been enabled for a given Agent
	Enabled *config.AgentComponentsToggle

	Worker     *worker.Worker
	NATSServer *natsserver.NATSServer
	APIServer  *apiserver.APIServer
	DataStore  *datastore.DataStore
}

// Agent is responsible for running bubbly as a collection of
// microservices, termed Components, with inter-service communication enabled by
// NATS.
type Agent struct {
	// Config stores the agent configuration determined from CLI inputs
	Config *config.AgentConfig
	// Components stores which Components are enabled on the agent as well as
	// the Component instances themselves.
	// Each Agent component represents an independent agent service.
	// They are termed Components due to the Component interface which they
	// implement, which enables inter-service communication through NATS.
	Components *Components
	// The type of deployment to deploy the Agent into.
	DeploymentType config.AgentDeploymentType
}

// New returns a new bubbly Agent initialised with the provided configuration
// and components
func New(bCtx *env.BubblyContext) *Agent {
	return &Agent{
		Config: bCtx.AgentConfig,
		Components: &Components{
			Enabled: bCtx.AgentConfig.EnabledComponents,
			NATSServer: &natsserver.NATSServer{
				ComponentCore: &component.ComponentCore{
					Type: component.NATSServerComponent,
				},
				Config: bCtx.AgentConfig.NATSServerConfig,
			},
		},
		DeploymentType: bCtx.AgentConfig.DeploymentType,
	}
}

// Run runs any enabled agent components.
func (a *Agent) Run(bCtx *env.BubblyContext) error {
	switch a.DeploymentType {
	case config.SingleDeployment:
		return a.runAsSingle(bCtx)
	default:
		return fmt.Errorf(
			`deployment type "%s" not implemented`,
			a.DeploymentType,
		)
	}
}

// runAsSingle runs agent components in an errgroup.
// Should any of the long-running
// processes go down, everything will be stopped, because every component is
// critical.
//
// In a distributed-server deployment we will probably want to improve upon
// this, because in
// this deployment strategy even if a single component goes
// down the other components will still be valuable to the rest of the bubbly fleet.
func (a *Agent) runAsSingle(bCtx *env.BubblyContext) error {
	g, agentContext := errgroup.WithContext(context.TODO())

	// If single, the we must also run the NATS server as a part of the
	// bubbly agent. Therefore we instance a new NATS server
	ns := natsserver.New(bCtx)
	g.Go(func() error {
		return ns.Run(bCtx, agentContext)
	})

	// Wait for the NATS Server to be able to accept connections before
	// connecting and setting up the other bubbly agent components
	if !ns.Server.ReadyForConnections(
		defaultNATSServerWait * time.Second) {
		return fmt.Errorf("NATS server took too long to allow client connection")
	}

	// Update the NATSServer Component of this agent
	a.Components.NATSServer = ns

	// Now set up any activated bubbly component within the errgroup
	// If enabled, initialise and run the data store
	if a.Config.EnabledComponents.DataStore {
		g.Go(func() error {
			dStore, err := datastore.New(bCtx)
			if err != nil {
				return fmt.Errorf("failed to create data store: %w", err)
			}

			// Components are responsible for maintaining individual connections
			// to the NATS Server.
			// Connect provides us with an unencoded connection which we later
			// use to establish encoded connections when publishing/subscribing
			if err := dStore.Connect(bCtx); err != nil {
				return fmt.Errorf("failed to connect to the NATS Server: %w", err)
			}
			defer dStore.Close()

			a.Components.DataStore = dStore
			return dStore.Run(bCtx, agentContext)
		})
	}

	// If enabled, initialise and run the API Server
	if a.Config.EnabledComponents.APIServer {
		g.Go(func() error {
			// establish the APIServer Component, which consists of a server.
			// Server instance responsible for handling HTTP requests and
			// publishing/requesting data via NATS to other Components.
			s, err := apiserver.New(bCtx)
			if err != nil {
				return fmt.Errorf("failed to create API server agent: %w", err)
			}

			// Components are responsible for maintaining individual connections
			// to the NATS Server. Connect connects to the NATSServer defined by
			// the bubbly context and saves the connection into the Component for
			// later use
			if err := s.Connect(bCtx); err != nil {
				return fmt.Errorf("failed to connect to the NATS Server: %w", err)
			}
			defer s.Close()
			return s.Run(bCtx, agentContext)
		})
	}

	// If enabled, initialise and run the worker
	if a.Config.EnabledComponents.Worker {
		g.Go(func() error {
			worker := worker.New(bCtx)
			if err := worker.Connect(bCtx); err != nil {
				return fmt.Errorf("failed to connect to the NATS Server: %w", err)
			}
			defer worker.Close()

			a.Components.Worker = worker
			return worker.Run(bCtx, agentContext)
		})
	}

	// Wait for any of the agent components to error. If so, we exit,
	// because in a single-server deployment every component is critical
	if err := g.Wait(); err != nil {
		return fmt.Errorf("error running one of the agent's components: %w", err)
	}
	return nil

}
