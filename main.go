package main

import (
	"embed"
	"os"

	"github.com/valocode/bubbly/cmd"
	"github.com/valocode/bubbly/env"
)

//go:embed ui/build/*
var bubblyUI embed.FS

func main() {
	// Set up the initial BubblyContext with config.Config defaults
	bCtx := env.NewBubblyContext(env.WithBubblyUI(&bubblyUI))

	// Create initial root command, binding global flags
	rootCmd := cmd.NewCmdRoot(bCtx)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
