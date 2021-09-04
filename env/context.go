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
	Version      *Version
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
		Version:       NewVersionInfo(),
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

func WithVersion(version *Version) func(*BubblyContext) {
	return func(bCtx *BubblyContext) {
		bCtx.Version = version
	}
}

// Version contains the SHA1 value for the commit
// in the Bubbly Git repo from which the current running
// binary was built. If the commit was tagged, the tag
// is also included.
type Version struct {
	Commit  string
	Version string
	Date    string
}

// NewVersionInfo returns a reference to the Version structure,
// populated with information provided at compile time.
func NewVersionInfo() *Version {
	return &Version{
		Version: "dev",
		Commit:  "dev",
		Date:    "dev",
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
