package cmd

import (
	"github.com/MakeNowJust/heredoc/v2"
	"github.com/fatih/color"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
	"github.com/valocode/bubbly/cmd/adapter"
	"github.com/valocode/bubbly/cmd/demo"
	"github.com/valocode/bubbly/cmd/release"
	"github.com/valocode/bubbly/cmd/server"
	"github.com/valocode/bubbly/config"
	"github.com/valocode/bubbly/env"
)

// NewCmdRoot creates a new cobra.Command representing "bubbly"
func NewCmdRoot(bCtx *env.BubblyContext) *cobra.Command {
	// cmd represents the apply command
	cmd := &cobra.Command{
		Use:   "bubbly",
		Short: "Bubbly - The Release Readiness Platform",
		Long: heredoc.Doc(`Bubbly is a release readiness platform helping you
			release software with confidence, continuously.

			To learn more visit https://bubbly.dev
			Documentation available at https://docs.bubbly.dev`),
		// Do not print usage on error
		SilenceUsage: true,
		// Do not print errors on error (we will do that ourselves)
		// SilenceErrors: true,
		// Run: ,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if bCtx.CLIConfig.Debug {
				bCtx.UpdateLogLevel(zerolog.DebugLevel)
			}
			if bCtx.CLIConfig.NoColor {
				color.NoColor = true
			}
		},
	}

	f := cmd.PersistentFlags()

	var test bool
	f.BoolVar(&test, "test", config.DefaultCLIDebugToggle, "specify whether to enable debug logging")
	f.BoolVar(&bCtx.CLIConfig.Debug, "debug", config.DefaultCLIDebugToggle, "specify whether to enable debug logging")
	f.BoolVar(&bCtx.CLIConfig.NoColor, "no-color", config.DefaultCLINoColorToggle, "specify whether to disable colorful logging")
	cmd.InitDefaultHelpFlag()

	cmd.AddCommand(adapter.New(bCtx))
	cmd.AddCommand(demo.New(bCtx))
	cmd.AddCommand(release.New(bCtx))
	cmd.AddCommand(server.New(bCtx))

	// finally, print the final configuration to be used by bubbly
	bCtx.Logger.Debug().
		Interface("server_config", bCtx.ServerConfig).
		Interface("store_config", bCtx.StoreConfig).
		Msg("bubbly config")
	return cmd
}
