package main

import (
	"embed"
	"os"

	"github.com/valocode/bubbly/cmd"
	"github.com/valocode/bubbly/env"
)

//   go:embed ui/build/* ui/build/_app/pages/__layout.*
var (
	// For some strange reason the __layout.* files in the svelte tree do not
	// get embedded unless specified explicitly... Bug?
	bubblyUI embed.FS
	version  = "dev"
	commit   = "dev"
	date     = "dev"
)

func main() {
	// Set up the initial BubblyContext with config.Config defaults
	bCtx := env.NewBubblyContext(
		env.WithBubblyUI(&bubblyUI),
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
