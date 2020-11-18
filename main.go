package main

import (
	"os"

	"github.com/swaggo/swag/example/basic/docs"

	"github.com/imdario/mergo"

	"github.com/rs/zerolog"
	"github.com/verifa/bubbly/cmd"
	"github.com/verifa/bubbly/config"
	"github.com/verifa/bubbly/env"
)

// @title Bubbly
// @version 1.0  // Change version here
// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@bubbly.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @termsOfService http://bubbly.io/terms/

// @host localhost:8080
func main() {
	setSwaggerInfo()
	// TODO: remove once migrated fully to bCtx.Logger
	// bCtx.Logger.Logger = bCtx.Logger.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	// 1. set up the initial BubblyContext with config.Config defaults
	// (we mergo.Merge later) and a default logger
	bCtx := env.NewBubblyContext()

	// 2. create initial root command, binding global flags
	rootCmd := cmd.NewCmdRoot(bCtx)

	// 3. parse the global flags set up by the rootCmd against
	// command line inputs in order to pull out the data that
	// makes up the BubblyContext
	// TODO: I believe this renders Viper reduntant. Remove.
	if err := rootCmd.ParseFlags(os.Args); err != nil {
		os.Exit(1)
	}

	fs := rootCmd.Flags()

	// 4. update the log level of the bubblyContext.Logger
	// if --debug flag is specified
	if debug, _ := fs.GetBool("debug"); debug {
		bCtx.UpdateLogLevel(zerolog.DebugLevel)
		bCtx.Logger.Debug().Str("log_level", bCtx.Logger.GetLevel().String()).
			Msg("updated log level")
	}

	// Because several of rootCmd's flags are mapped to BubblyContext.Config
	// attributes (and therefore reset when calling NewCmdRoot),
	// we need to merge the default configuration with any flags
	// provided by CLI.
	defaultConfig := config.NewDefaultConfig()

	if err := mergo.Merge(bCtx.Config, defaultConfig); err != nil {
		bCtx.Logger.Error().Err(err).Msg("error when merging configs")
		os.Exit(1)
	}

	// finally, print the final configuration to be used by bubbly
	bCtx.Logger.Info().
		Interface("final_config", bCtx.Config).
		Str("final_log_level", bCtx.Logger.GetLevel().String()).
		Msg("final bubbly context")

	rootCmd.Execute()
}

func setSwaggerInfo() {
	docs.SwaggerInfo.Title = "Bubbly Api"
	docs.SwaggerInfo.Description = "API schema and information for the bubbly server"
	docs.SwaggerInfo.Version = "1.0"
	// TODO(server): Have host be defined by environment variables
	docs.SwaggerInfo.Host = "localhost:8080"
	docs.SwaggerInfo.BasePath = "/"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}
}
