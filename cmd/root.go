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

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/verifa/bubbly/config"
	normalise "github.com/verifa/bubbly/util/normalise"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var (
	globalServerConfig config.ServerConfig
	globalConfigFile   string
	rootShort          = normalise.LongDesc(`
		bubbly provides a single binary for controlling both the bubbly server and bubbly client.
		
		Find more information: https://verifa.io/products/bubbly`)
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "bubbly",
	Short: rootShort,
	Run: func(cmd *cobra.Command, args []string) {
	},
}

// Execute is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	initLogger()
	cobra.OnInitialize(initConfig)

	f := rootCmd.PersistentFlags()

	f.StringVar(&globalConfigFile, "config", "", "config file (default is $HOME/.bubbly.yaml)")

	f.StringVar(&globalServerConfig.Host, "host", "", "bubbly server host")
	f.StringVar(&globalServerConfig.Port, "port", "", "bubbly server port")
	f.BoolVar(&globalServerConfig.Auth, "auth", false, "bubbly server auth")
	f.StringVar(&globalServerConfig.Token, "token", "", "bubbly server token")

	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	viper.BindPFlags(f)
}

type BubblyContext struct {
	Logger zerolog.Logger
	Config config.Config
}

func initLogger() {
	// Initialize Logger
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if gin.IsDebugging() {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	log.Logger = log.Output(
		zerolog.ConsoleWriter{
			Out:     os.Stderr,
			NoColor: false,
		},
	)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
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
		log.Debug().Str("config_file", viper.ConfigFileUsed()).Msg("loading configuration from config file")
		log.Debug().Interface("configuration", viper.AllSettings()).Msg("bubbly configuration")
	}
}
