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
	"os"

	"github.com/spf13/cobra"
	"github.com/verifa/bubbly/env"
	normalise "github.com/verifa/bubbly/util/normalise"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var (
	globalConfigFile string
	// debug            bool
	rootShort = normalise.LongDesc(`
		bubbly provides a single binary for controlling both the bubbly server and bubbly client.
		
		Find more information: https://verifa.io/products/bubbly`)
	// rootCmd *cobra.Command
)

// NewCmdRoot creates a new cobra.Command representing "bubbly"
func NewCmdRoot(bCtx *env.BubblyContext) *cobra.Command {
	// cmd represents the apply command
	cmd := &cobra.Command{
		Use:   "bubbly",
		Short: rootShort,
		Run: func(cmd *cobra.Command, args []string) {
		},
		PreRun: func(cmd *cobra.Command, _ []string) {
			viper.BindPFlags(cmd.PersistentFlags())
		},
	}

	initFlags(bCtx, cmd)

	cobra.OnInitialize(func() {
		initConfig(bCtx)
	})
	initCommands(bCtx, cmd)

	return cmd
}

func initCommands(bCtx *env.BubblyContext, cmd *cobra.Command) {
	applyCmd, _ := NewCmdApply(bCtx)
	cmd.AddCommand(applyCmd)

	describeCmd, _ := NewCmdDescribe(bCtx)
	cmd.AddCommand(describeCmd)

	serverCmd, _ := NewCmdServer(bCtx)
	cmd.AddCommand(serverCmd)
}

func initFlags(bCtx *env.BubblyContext, cmd *cobra.Command) {

	f := cmd.PersistentFlags()

	f.StringVar(&globalConfigFile, "config", "", "config file (default is $HOME/.bubbly.yaml)")

	f.StringVar(&bCtx.ServerConfig.Host, "host", "", "bubbly server host")
	f.StringVar(&bCtx.ServerConfig.Port, "port", "", "bubbly server port")
	f.BoolVar(&bCtx.ServerConfig.Auth, "auth", false, "bubbly server auth")
	f.StringVar(&bCtx.ServerConfig.Token, "token", "", "bubbly server token")
	// Option 1: just bind normally, then parse the flags in main.go
	// to determine the value
	f.Bool("debug", false, "set log level to debug")

	// Option 2: bind to a specific part of the bubbly context's config.Config
	// downside: the Config then contains log settings, which is
	// counter-intuitive given the context already has a Logger field.
	// advantage: cleaner in main.go, as no need for separate flag parsing.
	// f.BoolVar(&bCtx.Config.LoggerConfig.Debug, "debug", false, "set log level to debug")

	cmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	cmd.InitDefaultHelpFlag()

	viper.BindPFlags(f)
}

// initConfig reads in config file and ENV variables if set.
func initConfig(bCtx *env.BubblyContext) {
	if globalConfigFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(globalConfigFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".bubbly" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".bubbly")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		bCtx.Logger.Debug().Str("config_file", viper.ConfigFileUsed()).Msg("loading configuration from config file")
		bCtx.Logger.Debug().Interface("configuration", viper.AllSettings()).Msg("bubbly configuration")
	}
}
