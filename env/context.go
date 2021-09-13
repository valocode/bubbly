package env

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/valocode/bubbly/auth"
	"github.com/valocode/bubbly/config"
)

// BubblyConfig holds global bubbly state that is required to be injected into
// functions throughout the codebase.
type BubblyConfig struct {
	// Logger stores the global bubbly logger
	Logger        zerolog.Logger
	ReleaseConfig *config.ReleaseConfig
	ServerConfig  *config.ServerConfig
	// Store provider configuration
	StoreConfig  *config.StoreConfig
	ClientConfig *config.ClientConfig
	CLIConfig    *config.CLIConfig
	AuthConfig   *auth.Config
}

// NewBubblyContext sets up a default Bubbly Context
func NewBubblyContext() *BubblyConfig {
	return &BubblyConfig{
		Logger:        NewDefaultLogger(),
		ReleaseConfig: config.DefaultReleaseConfig(),
		ServerConfig:  config.DefaultServerConfig(),
		StoreConfig:   config.DefaultStoreConfig(),
		ClientConfig:  config.DefaultClientConfig(),
		CLIConfig:     config.DefaultCLIConfig(),
		AuthConfig:    config.DefaultAuthConfig(),
	}
}

// NewDefaultLogger sets up a default logger
func NewDefaultLogger() zerolog.Logger {
	// return zerolog.New(os.Stderr).With().Timestamp().Logger()
	return zerolog.New(zerolog.ConsoleWriter{
		Out: os.Stderr,
	}).
		With().
		Timestamp().
		Logger().
		Level(zerolog.InfoLevel)
}

// UpdateLogLevel is a convenience method for updating the log level of
// the zerolog.Logger managed by a BubblyContext instance
func (bCtx *BubblyConfig) UpdateLogLevel(level zerolog.Level) {
	bCtx.Logger = bCtx.Logger.Level(level)
}
