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
	Config *config.Config
	// TODO: Could also contain a client.Client... consider.
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

// NewBubblyContext sets up a default Bubbly Context
func NewBubblyContext() *BubblyContext {
	return &BubblyContext{
		Config: config.NewDefaultConfig(),
		Logger: NewDefaultLogger(),
	}
}

// GetServerConfig is a convenience method to extract the bubbly server
// configuration from a BubblyContext.
func (bCtx *BubblyContext) GetServerConfig() (*config.ServerConfig, error) {
	if sc := bCtx.Config.ServerConfig; sc != nil {
		return sc, nil
	}
	// Should never be reached as ServerConfig should always at least be
	// be instanced with default values (due to mergo.Merge)
	bCtx.Logger.Error().Msg("bubbly server not configured. This error indicates an internal bubbly error")
	return nil, fmt.Errorf("bubbly server not configured")
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
