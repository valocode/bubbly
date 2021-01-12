package env

import (
	"fmt"
	"os"

	"github.com/rs/zerolog"
	"github.com/verifa/bubbly/config"
)

// BubblyContext holds global bubbly state that is required to be injected into
// functions throughout the codebase.
type BubblyContext struct {
	// Logger stores the global bubbly logger
	Logger *zerolog.Logger
	// Config stores global bubbly configuration,
	// such as bubbly server configuration
	ServerConfig *config.ServerConfig
	// Store provider configuration
	StoreConfig *config.StoreConfig
	AgentConfig *config.AgentConfig
	// stores configurations for accessing resources
	ResourceConfig *config.ResourceConfig
	// TODO: Could also contain a client.Client... consider.
}

// NewBubblyContext sets up a default Bubbly Context
func NewBubblyContext() *BubblyContext {
	return &BubblyContext{
		Logger:         NewDefaultLogger(),
		ServerConfig:   config.DefaultServerConfig(),
		StoreConfig:    config.DefaultStoreConfig(),
		AgentConfig:    config.DefaultAgentConfig(),
		ResourceConfig: config.DefaultResourceConfig(),
	}
}

// NewDefaultLogger sets up a default logger
func NewDefaultLogger() *zerolog.Logger {
	// Initialize Logger
	logger := zerolog.New(zerolog.ConsoleWriter{
		Out:     os.Stderr,
		NoColor: false,
	}).With().Timestamp().Logger().Level(zerolog.InfoLevel)

	return &logger
}

// GetServerConfig is a convenience method to extract the bubbly server
// configuration from a BubblyContext.
func (bCtx *BubblyContext) GetServerConfig() *config.ServerConfig {
	return bCtx.ServerConfig
}

// UpdateLogLevel is a convenience method for updating the log level of
// the zerolog.Logger managed by a BubblyContext instance
func (bCtx *BubblyContext) UpdateLogLevel(level zerolog.Level) error {
	logger := bCtx.Logger.Level(level)
	bCtx.Logger = &logger

	if bCtx.Logger.GetLevel() != level {
		return fmt.Errorf("failed to update log level")
	}

	return nil
}
