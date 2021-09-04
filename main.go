package main

import (
	"embed"
	"os"

	"github.com/valocode/bubbly/cmd"
	"github.com/valocode/bubbly/env"
)

var (
	//go:embed ui/build ui/build/_app ui/build/_app/pages/*.js ui/build/_app/assets/pages/*.css
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
