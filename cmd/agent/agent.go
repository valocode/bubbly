package agent

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/fatih/color"

	"github.com/valocode/bubbly/agent"
	"github.com/valocode/bubbly/config"

	"github.com/imdario/mergo"
	"github.com/spf13/cobra"

	cmdutil "github.com/valocode/bubbly/cmd/util"
	"github.com/valocode/bubbly/env"
	"github.com/valocode/bubbly/util/normalise"
)

var (
	_         cmdutil.Options = (*AgentOptions)(nil)
	agentLong                 = normalise.LongDesc(
		`
		Starts a bubbly agent. The agent can be configured to run all components, or only some subset, 
		depending on the flags provided.

			$ bubbly agent

		Starts a bubbly agent with only the UI component
		`,
	)

	agentExample = normalise.Examples(
		`
		# Starts the bubbly agent with all components (API Server, NATS Server, Store and Worker) 
		using application defaults

		bubbly agent
		
		# Starts the bubbly agent running only the API Server components
		bubbly agent --api-server`,
	)
)

// AgentOptions holds everything necessary to run the command.
// Flag values received to the agent command are loaded into the embedded
// BubblyContext and used to run the various agent components
type AgentOptions struct {
	cmdutil.Options
	BubblyContext *env.BubblyContext
}

// NewCmdAgent creates a new cobra.Command representing "bubbly agent"
func NewCmdAgent(bCtx *env.BubblyContext) (*cobra.Command, *AgentOptions) {
	// Set the ClientType as NATS, because all agents should run internally to
	// bubbly and can talk directly to NATS
	bCtx.ClientConfig.ClientType = config.NATSClientType

	o := &AgentOptions{
		BubblyContext: bCtx,
	}

	// cmd represents the agent command
	cmd := &cobra.Command{
		Use:     "agent [flags]",
		Short:   "Start a bubbly agent",
		Long:    agentLong + "\n\n",
		Example: agentExample,
		RunE: func(cmd *cobra.Command, args []string) error {

			validationError := o.Validate(cmd)

			if validationError != nil {
				return validationError
			}

			resolveError := o.Resolve()

			if resolveError != nil {
				return resolveError
			}

			bCtx.Logger.Debug().
				Interface(
					"data_store",
					o.BubblyContext.AgentConfig.StoreConfig,
				).
				Interface(
					"nats_server",
					bCtx.AgentConfig.NATSServerConfig,
				).
				Interface(
					"enabled_components",
					bCtx.AgentConfig.EnabledComponents,
				).
				Str(
					"deployment_type",
					bCtx.AgentConfig.DeploymentType.String(),
				).
				Msg("agent configuration")

			runError := o.Run()

			if runError != nil {
				return runError
			}

			o.Print()
			return nil
		},
		PreRunE: func(cmd *cobra.Command, _ []string) error {
			// prior to running the agent, we merge defaults with the
			// config provided by command flags to make sure
			// we have a complete configuration
			if err := mergo.Merge(
				o.BubblyContext.AgentConfig,
				config.DefaultAgentConfig(),
			); err != nil {
				return fmt.Errorf("error merging agent configuration with defaults: %w", err)
			}

			return nil
		},
	}

	f := cmd.Flags()

	// agent.AgentDeploymentType's underlying type is a string,
	// so we cast to *string on bind
	f.StringVar(
		(*string)(&o.BubblyContext.AgentConfig.DeploymentType),
		"deployment-type",
		config.DefaultDeploymentType.String(),
		"the type of agent deployment. Options: single",
	)
	f.BoolVar(
		&o.BubblyContext.AgentConfig.EnabledComponents.NATSServer,
		"nats-server",
		config.DefaultNATSServerToggle,
		"whether to run the NATS Server on this agent",
	)
	f.BoolVar(
		&o.BubblyContext.AgentConfig.EnabledComponents.APIServer,
		"api-server",
		config.DefaultAPIServerToggle,
		"whether to run the api server on this agent",
	)
	f.BoolVar(
		&o.BubblyContext.AgentConfig.EnabledComponents.DataStore,
		"data-store",
		config.DefaultDataStoreToggle,
		"whether to run the data store on this agent",
	)
	f.BoolVar(
		&o.BubblyContext.AgentConfig.EnabledComponents.Worker,
		"worker",
		config.DefaultWorkerToggle,
		"whether to run a bubbly worker on this agent",
	)
	f.StringVar(
		(*string)(&o.BubblyContext.AgentConfig.StoreConfig.Provider),
		"data-store-provider",
		config.DefaultStoreProvider,
		"provider of the bubbly data store",
	)
	port, _ := strconv.Atoi(config.DefaultNATSServerPort)
	f.IntVar(
		&o.BubblyContext.AgentConfig.NATSServerConfig.Port,
		"nats-server-port",
		port,
		"port of the NATS Server",
	)
	httpPort, _ := strconv.Atoi(config.DefaultNATSServerHTTPPort)
	f.IntVar(
		&o.BubblyContext.AgentConfig.NATSServerConfig.HTTPPort,
		"nats-server-http-port",
		httpPort,
		"HTTP Port of the NATS Server",
	)
	f.StringVar(
		&o.BubblyContext.AgentConfig.StoreConfig.PostgresAddr,
		"data-store-addr",
		config.DefaultPostgresAddr,
		"address of the data store",
	)
	f.StringVar(
		&o.BubblyContext.AgentConfig.StoreConfig.PostgresUser,
		"data-store-username",
		config.DefaultPostgresUser,
		"username of the data store",
	)
	f.StringVar(
		&o.BubblyContext.AgentConfig.StoreConfig.PostgresPassword,
		"data-store-password",
		config.DefaultPostgresPassword,
		"password of the data store",
	)
	f.StringVar(
		&o.BubblyContext.AgentConfig.StoreConfig.PostgresDatabase,
		"data-store-database",
		config.DefaultPostgresDatabase,
		"database of the data store",
	)

	return cmd, o
}

// Validate checks the AgentOptions to see if there is sufficient information run the command.
func (o *AgentOptions) Validate(cmd *cobra.Command) error {
	return nil
}

// Resolve resolves various AgentOptions attributes from the
// provided arguments to the Command
func (o *AgentOptions) Resolve() error {
	// Resolve the agent components

	// if the user has specified specific components, only run those
	if !reflect.DeepEqual(
		o.BubblyContext.AgentConfig.EnabledComponents,
		config.DefaultAgentComponentsEnabled(),
	) {
		o.BubblyContext.Logger.Debug().
			Interface(
				"enabled_components",
				o.BubblyContext.AgentConfig.EnabledComponents,
			).
			Msg("one or more bubbly agent component explicitly set. Only enabling those specific agent components")
		return nil
	}
	o.BubblyContext.Logger.Debug().Msg("no agent components explicitly set. Enabling all agent components")

	// If no specific component has been set,
	// the agent should run all components
	o.BubblyContext.AgentConfig.EnabledComponents = &config.AgentComponentsToggle{
		APIServer:  true,
		DataStore:  true,
		Worker:     true,
		NATSServer: true,
	}

	return nil
}

// Run takes the validated and resolved AgentOptions, creates an agent.Agent
// and runs any activated bubbly agent component
func (o *AgentOptions) Run() error {

	a := agent.New(o.BubblyContext)

	if err := a.Run(o.BubblyContext); err != nil {
		return fmt.Errorf("error while running bubbly agent: %w", err)
	}
	return nil
}

// Print the success of instancing a new bubbly agent to the user
func (o *AgentOptions) Print() {
	if o.BubblyContext.CLIConfig.Color {
		color.Green("agent provisioned successfully")
	} else {
		fmt.Println("agent provisioned successfully")
	}
}
