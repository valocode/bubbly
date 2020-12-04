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
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	cmdutil "github.com/verifa/bubbly/cmd/util"
	"github.com/verifa/bubbly/env"
	"github.com/verifa/bubbly/server"
	normalise "github.com/verifa/bubbly/util/normalise"
)

var (
	_          cmdutil.Options = (*ApplyOptions)(nil)
	serverLong                 = normalise.LongDesc(`
		Start a bubbly server

			$ bubbly server [flags]
		`)

	serverExample = normalise.Examples(`
		# Start the bubbly server using defaults (http://localhost:8080)
		bubbly server
		
		# Start the bubbly server at port 8040
		bubbly --port 8040 server`)
	router        *gin.Engine
	bubblyVersion = "0.0.1"
)

// ServerOptions -
type ServerOptions struct {
	o             cmdutil.Options //embedding
	BubblyContext *env.BubblyContext
	Command       string
	Args          []string

	// Result from o.Run() - success / failure for the server
	Result bool
}

// NewCmdServer creates a new cobra.Command representing "bubbly server"
func NewCmdServer(bCtx *env.BubblyContext) (*cobra.Command, *ServerOptions) {
	o := &ServerOptions{
		Command:       "server",
		BubblyContext: bCtx,
		Result:        false,
	}

	// cmd represents the apply command
	cmd := &cobra.Command{
		Use:     "server [flags]",
		Short:   "Start a bubbly server",
		Long:    serverLong + "\n\n" + cmdutil.SuggestBubblyResources(),
		Example: serverExample,
		RunE: func(cmd *cobra.Command, args []string) error {
			bCtx.Logger.Debug().Strs("arguments", args).
				Msg("server arguments")
			bCtx.Logger.Debug().Interface("bubbly_context", o.BubblyContext).
				Msg("bubbly context")

			o.Args = args

			validationError := o.Validate(cmd)

			if validationError != nil {
				return validationError
			}

			resolveError := o.Resolve(cmd)

			if resolveError != nil {
				return resolveError
			}
			runError := o.Run()

			if runError != nil {
				return runError
			}

			o.Print(cmd)
			return nil
		},
		PreRun: func(cmd *cobra.Command, _ []string) {
			bCtx.Logger.Debug().Interface("configuration", viper.AllSettings()).Msg("bubbly viper configuration")
		},
	}

	return cmd, o
}

// Validate checks the ServerOptions to see if there is sufficient information run the command.
func (o *ServerOptions) Validate(cmd *cobra.Command) error {
	if len(o.Args) != 0 {
		return cmdutil.UsageErrorf(cmd, "Unexpected args: %v", o.Args)
	}
	// This should never be reached if we have set the ServerOptions.Config correctly with defaults.
	if (o.BubblyContext.ServerConfig.Host == "") || (o.BubblyContext.ServerConfig.Port == "") {
		return fmt.Errorf("Internal Error: Server configs missing.")
	}
	return nil
}

// Resolve resolves various ApplyOptions attributes from the provided arguments to cmd
func (o *ServerOptions) Resolve(cmd *cobra.Command) error {
	return nil
}

// Run runs the server command over the validated ServerOptions configuration
func (o *ServerOptions) Run() error {
	server.SetVersion(bubblyVersion)
	// initialize the router's endpoints

	err := server.ListenAndServe(o.BubblyContext)

	if err != nil {
		o.Result = false
		return fmt.Errorf("error while serving: %w", err)
	}
	o.Result = true

	return nil
}

// Print formats and prints the ServerOptions.Result from o.Run()
func (o *ServerOptions) Print(cmd *cobra.Command) {
	fmt.Fprintf(cmd.OutOrStdout(), "Server result: %t\n", o.Result)
}
