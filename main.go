package main

import (
	"os"

	"github.com/valocode/bubbly/cmd"
	"github.com/valocode/bubbly/config"
	"github.com/valocode/bubbly/ui"
)

var (
	version = "dev"
	commit  = "dev"
	date    = "dev"
)

//go:generate go run docs/gen.go

func main() {
	// Set up the initial BubblyConfig with config.Config defaults
	bCtx := config.NewBubblyConfig(
		config.WithBubblyUI(&ui.Build),
		config.WithVersion(&config.Version{
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
