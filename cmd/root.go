package cmd

import (
	"github.com/spf13/cobra"

	agentCmd "github.com/valocode/bubbly/cmd/agent"
	applyCmd "github.com/valocode/bubbly/cmd/apply"
	getCmd "github.com/valocode/bubbly/cmd/get"
	queryCmd "github.com/valocode/bubbly/cmd/query"
	releaseCmd "github.com/valocode/bubbly/cmd/release"
	schemaCmd "github.com/valocode/bubbly/cmd/schema"
	versionCmd "github.com/valocode/bubbly/cmd/version"
	"github.com/valocode/bubbly/cmd/topics"
	"github.com/valocode/bubbly/cmd/util"
	"github.com/valocode/bubbly/config"
	"github.com/valocode/bubbly/env"
)

var (
	rootShort = util.LongDesc(`
		bubbly: release readiness in a bubble
		
		Find more information: https://bubbly.dev`)
)

// NewCmdRoot creates a new cobra.Command representing "bubbly"
func NewCmdRoot(bCtx *env.BubblyContext) *cobra.Command {
	// cmd represents the apply command
	cmd := &cobra.Command{
		Use:   "bubbly",
		Short: rootShort,
		// Do not print usage on error
		SilenceUsage: true,
		// Do not print errors on error (we will do that ourselves)
		SilenceErrors: true,
		FParseErrWhitelist: cobra.FParseErrWhitelist{
			// Allow unknown flags for parsing to proceed in cases
			// where flags for child commands are provided.
			UnknownFlags: true,
		},
	}

	initFlags(bCtx, cmd)

	initCommands(bCtx, cmd)

	return cmd
}

func initCommands(bCtx *env.BubblyContext, cmd *cobra.Command) {
	applyCmd, _ := applyCmd.NewCmdApply(bCtx)
	cmd.AddCommand(applyCmd)

	agentCmd, _ := agentCmd.NewCmdAgent(bCtx)
	cmd.AddCommand(agentCmd)

	// help topics
	cmd.AddCommand(topics.NewHelpTopic("environment"))

	getCmd, _ := getCmd.NewCmdGet(bCtx)
	cmd.AddCommand(getCmd)

	cmd.AddCommand(releaseCmd.New(bCtx))
	cmd.AddCommand(queryCmd.New(bCtx))
	cmd.AddCommand(schemaCmd.NewCmdSchema(bCtx))
	cmd.AddCommand(versionCmd.New(bCtx))
}

func initFlags(bCtx *env.BubblyContext, cmd *cobra.Command) {

	f := cmd.PersistentFlags()

	f.StringVar(&bCtx.ServerConfig.Host, "host", config.DefaultAPIServerHost, "bubbly API server host")
	f.StringVar(&bCtx.ServerConfig.Port, "port", config.DefaultAPIServerPort, "bubbly API server port")

	f.Bool("debug", config.DefaultDebugToggle, "specify whether to enable debug logging")

	cmd.InitDefaultHelpFlag()
}
