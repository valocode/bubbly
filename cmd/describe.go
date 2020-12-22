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

	"github.com/spf13/viper"

	"github.com/spf13/cobra"
	client "github.com/verifa/bubbly/client"
	cmdutil "github.com/verifa/bubbly/cmd/util"
	"github.com/verifa/bubbly/env"
	normalise "github.com/verifa/bubbly/util/normalise"
)

var (
	// a mock return for getting list of available resources
	// TODO: Think about if/how this pre-validation should take place. Perhaps this is instead a server-side task?
	availableResourceTypes                 = []string{"load", "pipeline", "extract", "transform"}
	_                      cmdutil.Options = (*DescribeOptions)(nil)
	describeLong                           = normalise.LongDesc(`
		Show details of a specific resource or group of resources

		Print a detailed description of the selected resources, including related resources such
		as events or controllers. You may select all resources of a given
		type by providing only the type, or additionally provide a name prefix to describe a single resource by name. 
		
		For example:

		    $ bubbly describe TYPE NAME_PREFIX`)

	describeExample = normalise.Examples(`
		# Describe an extract with name 'default'
		bubbly describe extract default

		# Describe an extract with name 'default'
		bubbly describe extract/default

		# Describe the latest version of an load with name 'sonarqube'
		bubbly describe load sonarqube

		# Describe the latest version of an load with name 'sonarqube'
		bubbly describe load/sonarqube

		# Describe the version '239iq0wi' of an load with name 'sonarqube' 
		bubbly describe load sonarqube --version 239iq0wi

		# Describe all load resources
		bubbly describe load

		# Describe a pipeline with name 'simple_pipeline'
		bubbly describe pipeline simple_pipeline`)
)

type DescribeResource struct {
	Name string
	Type string
}

// DescribeOptions -
type DescribeOptions struct {
	o             cmdutil.Options //embedding
	BubblyContext *env.BubblyContext

	Command string
	Args    []string

	// results of arg parsing by the Resolve function
	Resource      DescribeResource
	ResourceGroup string

	// describe a specific resource (false) or group of resources (true)
	Group   bool
	Version string
	Errs    []error

	// Results from o.Run()
	Result map[string]client.DescribeResourceReturn
}

// NewCmdDescribe creates a new cobra.Command representing "bubbly describe"
func NewCmdDescribe(bCtx *env.BubblyContext) (*cobra.Command, *DescribeOptions) {
	o := &DescribeOptions{
		Command:       "describe",
		BubblyContext: bCtx,
		Resource:      DescribeResource{},
		Version:       "v1",
		Result:        make(map[string]client.DescribeResourceReturn, 1),
	}

	// cmd represents the describe command
	cmd := &cobra.Command{
		Use:     "describe (TYPE [NAME_PREFIX] | TYPE/NAME)",
		Short:   "Show details of a specific resource or group of resources",
		Long:    describeLong + "\n\n" + cmdutil.SuggestBubblyResources(),
		Example: describeExample,
		RunE: func(cmd *cobra.Command, args []string) error {
			bCtx.Logger.Debug().Strs("arguments", args).
				Msg("describe arguments")

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
			// viper.BindPFlags(rootCmd.PersistentFlags())
			viper.BindPFlags(cmd.PersistentFlags())
			bCtx.Logger.Debug().Interface("configuration", viper.AllSettings()).Msg("bubbly configuration")
		},
	}

	f := cmd.Flags()

	f.StringVarP(&o.Version, "version", "v", o.Version, "Version of resource to filter on. (e.g. -v 239iq0wi")
	// cmd.Flags().StringVarP(&o.Config.ServerConfig.Port, "port", "p", o.Config.ServerConfig.Port, "bubbly server port")
	// if err := viper.BindPFlag("port", cmd.Flags().Lookup("port")); err != nil {
	// 	bCtx.Logger.Error().Msg(err.Error())
	// }

	viper.BindPFlags(f)

	return cmd, o
}

// Validate checks the DescribeOptions to see if there is sufficient information run the command.
func (o *DescribeOptions) Validate(cmd *cobra.Command) error {
	if len(o.Args) == 0 {
		return fmt.Errorf("you must specify the type of resource to describe. %s", cmdutil.SuggestBubblyResources())
	}

	if !validResourceType(o.Args[0]) {
		return fmt.Errorf("Invalid resource type %s. %s", o.Args[0], cmdutil.SuggestBubblyResources())
	}

	if len(o.Args) == 1 {
		o.Group = true
	}

	return nil
}

// Resolve resolves various DescribeOptions attributes from the provided arguments to cmd
func (o *DescribeOptions) Resolve(cmd *cobra.Command) error {
	if !o.Group {
		o.Resource = DescribeResource{
			Type: o.Args[0],
			Name: o.Args[1],
		}
		return nil
	}
	o.ResourceGroup = o.Args[0]
	return nil
}

// Run runs the describe command over the validated DescribeOptions configuration
func (o *DescribeOptions) Run() error {
	c, err := client.New(o.BubblyContext)
	if err != nil {
		return fmt.Errorf("failed to set up client: %w", err)
	}

	if !o.Group {
		resourceDescription, err := c.DescribeResource(o.BubblyContext, o.Resource.Type, o.Resource.Name, o.Version)

		if err != nil {
			return fmt.Errorf("failed to describe resource: %w", err)
		}

		o.Result[o.Resource.Name] = *resourceDescription

		return nil
	}

	resourceDescriptions, err := c.DescribeResourceGroup(o.BubblyContext, o.ResourceGroup, o.Version)

	if err != nil {
		return fmt.Errorf("failed to describe resource group: %w", err)
	}

	o.Result = resourceDescriptions

	return nil
}

// Print formats and prints the outputs the response from o.Run()
func (o *DescribeOptions) Print(cmd *cobra.Command) {
	if !o.Group {
		for _, r := range o.Result {
			// Push printing to cmd's Out, making it testable.
			fmt.Fprintf(cmd.OutOrStdout(), "EXISTS: %t, STATUS: %s, EVENTS:\n", r.Exists, r.Status)
			for _, e := range r.Events {
				fmt.Fprintf(cmd.OutOrStdout(), "	status: %s, age: %s, message %s\n", e.Status, e.Age, e.Message)
			}
		}
		return
	}
	for rName, r := range o.Result {
		// Push printing to cmd's Out, making it testable.
		fmt.Fprintf(cmd.OutOrStdout(), "RESOURCE: %s, EXISTS: %t, STATUS: %s, EVENTS:\n", rName, r.Exists, r.Status)
		for _, e := range r.Events {
			fmt.Fprintf(cmd.OutOrStdout(), "	status: %s, age: %s, message %s\n", e.Status, e.Age, e.Message)
		}
	}
}

func validResourceType(rTy string) bool {
	for _, ty := range availableResourceTypes {
		if rTy == ty {
			return true
		}
	}
	return false
}
