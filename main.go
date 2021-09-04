package main

import (
	"os"

	"github.com/valocode/bubbly/cmd"
	"github.com/valocode/bubbly/env"
	"github.com/valocode/bubbly/ui"
)

var (
	version = "dev"
	commit  = "dev"
	date    = "dev"
)

func main() {
	// Set up the initial BubblyContext with config.Config defaults
	bCtx := env.NewBubblyContext(
		env.WithBubblyUI(&ui.Build),
		env.WithVersion(&env.Version{
			Version: version,
			Commit:  commit,
			Date:    date,
		}),
	)

	// Create initial root command, binding global flags
	rootCmd := cmd.NewCmdRoot(bCtx)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
