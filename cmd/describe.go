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
	"errors"
	"fmt"

	"github.com/spf13/viper"

	"github.com/davecgh/go-spew/spew"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	client "github.com/verifa/bubbly/client"
	cmdutil "github.com/verifa/bubbly/cmd/util"
	"github.com/verifa/bubbly/config"
	normalise "github.com/verifa/bubbly/util/normalise"
)

var (
	// a mock return for getting list of available resources
	// TODO: Think about if/how this pre-validation should take place. Perhaps this is instead a server-side task?
	availableResourceTypes                 = []string{"publish", "pipeline", "importer", "translator"}
	_                      cmdutil.Options = (*DescribeOptions)(nil)
	describeLong                           = normalise.LongDesc(`
		Show details of a specific resource or group of resources

		Print a detailed description of the selected resources, including related resources such
		as events or controllers. You may select all resources of a given
		type by providing only the type, or additionally provide a name prefix to describe a single resource by name. For example:

		    $ bubbly describe TYPE NAME_PREFIX`)

	describeExample = normalise.Examples(`
		# Describe an importer with name 'default'
		bubbly describe importer default

		# Describe an importer with name 'default'
		bubbly describe importer/default

		# Describe the latest version of an publish with name 'sonarqube'
		bubbly describe publish sonarqube

		# Describe the latest version of an publish with name 'sonarqube'
		bubbly describe publish/sonarqube

		# Describe the version '239iq0wi' of an publish with name 'sonarqube' 
		bubbly describe publish sonarqube --version 239iq0wi

		# Describe all publish resources
		bubbly describe publish

		# Describe a pipeline with name 'simple_pipeline'
		bubbly describe pipeline simple_pipeline`)
)

type DescribeResource struct {
	Name string
	Type string
}

// DescribeOptions -
type DescribeOptions struct {
	o      cmdutil.Options //embedding
	Config *config.Config

	// sc ServerConfig

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
func NewCmdDescribe() (*cobra.Command, *DescribeOptions) {
	o := &DescribeOptions{
		Command:  "describe",
		Config:   config.NewDefaultConfig(),
		Resource: DescribeResource{},
		Version:  "v1",
		Result:   make(map[string]client.DescribeResourceReturn, 1),
	}

	// cmd represents the describe command
	cmd := &cobra.Command{
		Use:     "describe (TYPE [NAME_PREFIX] | TYPE/NAME)",
		Short:   "Show details of a specific resource or group of resources",
		Long:    describeLong + "\n\n" + cmdutil.SuggestBubblyResources(),
		Example: describeExample,
		RunE: func(cmd *cobra.Command, args []string) error {
			log.Debug().Msgf("Args provided to describe: args: %+v, length: %d", args, len(args))
			config, err := config.SetupConfigs()

			if err != nil {
				return err
			}

			o.Config = config

			spew.Dump("Merged configuration:", o.Config)
			o.Args = args

			validationError := o.Validate(cmd)

			if validationError != nil {
				return validationError
			}

			resolveError := o.Resolve(cmd)

			if resolveError != nil {
				return resolveError
			}
			if len(o.Errs) != 0 {
				return errors.New("Non zero number of errors with DescribeOptions")
			}
			runError := o.Run()

			if runError != nil {
				return runError
			}

			o.Print(cmd)
			return nil
		},
		PreRun: func(cmd *cobra.Command, _ []string) {
			viper.BindPFlags(rootCmd.PersistentFlags())
			viper.BindPFlags(cmd.PersistentFlags())
			for _, v := range viper.AllKeys() {
				log.Debug().Msgf("Key: %s, Value: %v\n", v, viper.Get(v))
			}
		},
	}

	f := cmd.Flags()

	f.StringVarP(&o.Version, "version", "v", o.Version, "Version of resource to filter on. (e.g. -v 239iq0wi")
	// cmd.Flags().StringVarP(&o.Config.ServerConfig.Port, "port", "p", o.Config.ServerConfig.Port, "bubbly server port")
	// if err := viper.BindPFlag("port", cmd.Flags().Lookup("port")); err != nil {
	// 	log.Error().Msg(err.Error())
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
	c, err := cmdutil.ClientSetup(*o.Config.ServerConfig)
	if err != nil {
		return err
	}

	if !o.Group {
		resourceDescription, err := c.DescribeResource(o.Resource.Type, o.Resource.Name, o.Version)

		if err != nil {
			return err
		}

		o.Result[o.Resource.Name] = resourceDescription

		return nil
	}

	resourceDescriptions, err := c.DescribeResourceGroup(o.ResourceGroup, o.Version)

	if err != nil {
		return err
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

func init() {
	describeCmd, _ := NewCmdDescribe()
	rootCmd.AddCommand(describeCmd)
}
