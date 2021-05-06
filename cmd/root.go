package cmd

import (
	"github.com/spf13/cobra"

	agentCmd "github.com/valocode/bubbly/cmd/agent"
	applyCmd "github.com/valocode/bubbly/cmd/apply"
	getCmd "github.com/valocode/bubbly/cmd/get"
	schemaCmd "github.com/valocode/bubbly/cmd/schema"
	"github.com/valocode/bubbly/cmd/topics"
	"github.com/valocode/bubbly/config"
	"github.com/valocode/bubbly/env"

	"github.com/valocode/bubbly/util/normalise"
)

var (
	rootShort = normalise.LongDesc(`
		bubbly: release readiness in a bubble
		
		Find more information: https://bubbly.dev`)
)

// NewCmdRoot creates a new cobra.Command representing "bubbly"
func NewCmdRoot(bCtx *env.BubblyContext) *cobra.Command {
	// cmd represents the apply command
	cmd := &cobra.Command{
		Use:   "bubbly",
		Short: rootShort,
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

	cmd.AddCommand(schemaCmd.NewCmdSchema(bCtx))
}

func initFlags(bCtx *env.BubblyContext, cmd *cobra.Command) {

	f := cmd.PersistentFlags()

	f.StringVar(&bCtx.ServerConfig.Host, "host", config.DefaultAPIServerHost, "bubbly API server host")
	f.StringVar(&bCtx.ServerConfig.Port, "port", config.DefaultAPIServerPort, "bubbly API server port")

	f.Bool("debug", config.DefaultDebugToggle, "specify whether to enable debug logging")

	cmd.InitDefaultHelpFlag()
}
