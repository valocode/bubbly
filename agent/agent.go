package agent

import (
	"fmt"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/verifa/bubbly/agent/component"
	"github.com/verifa/bubbly/agent/component/apiserver"
	"github.com/verifa/bubbly/agent/component/datastore"
	"github.com/verifa/bubbly/agent/component/natsserver"
	"github.com/verifa/bubbly/agent/component/ui"
	"github.com/verifa/bubbly/agent/component/worker"
	"github.com/verifa/bubbly/config"
	"github.com/verifa/bubbly/env"
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
	UI         *ui.UI
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
		},
	}
}

// Init initialises a bubbly agent struct with the provided deploymentType
// and initial NATS Server component configuration
func (a *Agent) Init(bCtx *env.BubblyContext) {

	a.DeploymentType = bCtx.AgentConfig.DeploymentType

	a.Components.NATSServer = &natsserver.NATSServer{
		ComponentCore: &component.ComponentCore{
			Type: component.NATSServerComponent,
		},
		Config:     a.Config.NATSServerConfig,
		Server:     nil,
		ServerConn: nil,
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
	g := new(errgroup.Group)

	// If single, the we must also run the NATS server as a part of the
	// bubbly agent. Therefore we instance a new NATS server
	n := a.Components.NATSServer

	if err := n.New(bCtx); err != nil {
		return fmt.Errorf("unable to initialise NATS Server: %w", err)
	}

	g.Go(func() error {
		return n.Run(bCtx)
	})

	// Wait for the NATS Server to be able to accept connections before
	// connecting and setting up the other bubbly agent components
	if !n.Server.ReadyForConnections(defaultNATSServerWait * time.Second) {
		return fmt.Errorf("NATS server took too long to allow client connection")
	}

	// Update the NATSServer associated with this agent with the newly
	// established *nats.Server and ServerConn
	a.Components.NATSServer = n

	// Connect to the NATS Server. We load the connection into
	// agent.NATSServerConn, thereby allowing us to pass an establish
	// NATS server connection to other components.
	// TODO: Consider whether it would be more appropriate for
	//  components to individually connect.
	if err := a.Components.NATSServer.Connect(bCtx); err != nil {
		return fmt.Errorf("failed to connect to the NATS Server: %w", err)
	}

	// Now, just set up any activated bubbly component within the errgroup

	// if enabled, initialise and run the data store
	if a.Config.EnabledComponents.DataStore {
		dStore := &datastore.DataStore{
			ComponentCore: &component.ComponentCore{
				Type:           component.DataStoreComponent,
				NATSServerConn: n.ServerConn,
				Subscriptions:  nil,
				Publications:   nil,
			},
		}

		if err := dStore.New(bCtx); err != nil {
			return fmt.Errorf("failed to create data store: %w", err)
		}
		a.Components.DataStore = dStore
		g.Go(func() error {
			return dStore.Run(bCtx)
		})
	}

	// if enabled, initialise and run the API Server
	if a.Config.EnabledComponents.APIServer {
		server := &apiserver.APIServer{
			ComponentCore: &component.ComponentCore{
				Type:           component.APIServerComponent,
				NATSServerConn: n.ServerConn,
				Subscriptions:  nil,
				Publications:   nil,
			},
		}
		g.Go(func() error {
			return server.Run(bCtx)
		})
	}

	// if enabled, initialise and run the worker
	if a.Config.EnabledComponents.Worker {
		worker := &worker.Worker{
			ComponentCore: &component.ComponentCore{
				Type:           component.WorkerComponent,
				NATSServerConn: n.ServerConn,
				Subscriptions: component.Subscriptions{
					// We know that workers are responsible for running
					// pipeline_run interval resources. Therefore,
					// we subscribe to its corresponding subject
					component.Subscription{
						Subject: component.WorkerPipelineRunIntervalSubject,
						Queue:   component.QueueWorker,
					},
				},
				Publications: nil,
			},
		}

		a.Components.Worker = worker
		g.Go(func() error {
			return worker.Run(bCtx)
		})
	}

	// if enabled, initialise and run the UI
	if a.Config.EnabledComponents.UI {
		ui := &ui.UI{
			ComponentCore: &component.ComponentCore{
				Type:           component.UIComponent,
				NATSServerConn: n.ServerConn,
				Subscriptions:  nil,
				Publications:   nil,
			},
		}

		a.Components.UI = ui
		g.Go(func() error {
			return ui.Run(bCtx)
		})
	}

	// wait for any of the agent components to error. If so, we exit,
	// because in a single-server deployment every component is critical
	if err := g.Wait(); err != nil {
		return fmt.Errorf("error running one of the agent's components: %w", err)
	}
	return nil
}
