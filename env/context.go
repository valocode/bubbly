package env

import (
	"embed"
	"os"

	"github.com/rs/zerolog"
	"github.com/valocode/bubbly/config"
)

// BubblyContext holds global bubbly state that is required to be injected into
// functions throughout the codebase.
type BubblyContext struct {
	// UI stores the embedded bubbly frontend
	UI *embed.FS
	// Logger stores the global bubbly logger
	Logger        zerolog.Logger
	ReleaseConfig *config.ReleaseConfig
	ServerConfig  *config.ServerConfig
	// Store provider configuration
	StoreConfig  *config.StoreConfig
	ClientConfig *config.ClientConfig
	CLIConfig    *config.CLIConfig
}

// NewBubblyContext sets up a default Bubbly Context
func NewBubblyContext(opts ...func(*BubblyContext)) *BubblyContext {
	bCtx := BubblyContext{
		Logger:        NewDefaultLogger(),
		ReleaseConfig: config.DefaultReleaseConfig(),
		ServerConfig:  config.DefaultServerConfig(),
		StoreConfig:   config.DefaultStoreConfig(),
		ClientConfig:  config.DefaultClientConfig(),
		CLIConfig:     config.DefaultCLIConfig(),
	}
	for _, opt := range opts {
		opt(&bCtx)
	}

	return &bCtx
}

func WithBubblyUI(fs *embed.FS) func(*BubblyContext) {
	return func(bCtx *BubblyContext) {
		bCtx.UI = fs
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
func (bCtx *BubblyContext) UpdateLogLevel(level zerolog.Level) {
	bCtx.Logger = bCtx.Logger.Level(level)
}
