/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"github.com/spf13/cobra"

	"github.com/verifa/bubbly/env"

	normalise "github.com/verifa/bubbly/util/normalise"
)

var (
	globalConfigFile string
	rootShort        = normalise.LongDesc(`
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
	applyCmd, _ := NewCmdApply(bCtx)
	cmd.AddCommand(applyCmd)

	agentCmd, _ := NewCmdAgent(bCtx)
	cmd.AddCommand(agentCmd)

	// help topics
	cmd.AddCommand(NewHelpTopic("environment"))
}

func initFlags(bCtx *env.BubblyContext, cmd *cobra.Command) {

	f := cmd.PersistentFlags()

	f.StringVar(&bCtx.ServerConfig.Host, "host", "", "bubbly API server host")
	f.StringVar(&bCtx.ServerConfig.Port, "port", "", "bubbly API server port")
	f.Bool("debug", false, "set log level to debug")

	cmd.InitDefaultHelpFlag()
}
