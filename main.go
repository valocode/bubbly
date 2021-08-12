package main

import (
	"os"

	"github.com/valocode/bubbly/cmd"
	"github.com/valocode/bubbly/env"
)

func main() {
	// Set up the initial BubblyContext with config.Config defaults
	bCtx := env.NewBubblyContext()

	// Create initial root command, binding global flags
	rootCmd := cmd.NewCmdRoot(bCtx)

	// // Parse the global flags set up by the rootCmd against
	// if err := rootCmd.ParseFlags(os.Args); err != nil {
	// 	log.Fatal("error parsing flags: ", err)
	// }

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
