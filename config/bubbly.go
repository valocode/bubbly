package config

import (
	"embed"
	"os"

	"github.com/rs/zerolog"
	"github.com/valocode/bubbly/auth"
)

// BubblyConfig holds global bubbly state that is required to be injected into
// functions throughout the codebase.
type BubblyConfig struct {
	// UI stores the embedded bubbly frontend
	UI *embed.FS
	// Logger stores the global bubbly logger
	Logger        zerolog.Logger
	ReleaseConfig *ReleaseConfig
	ServerConfig  *ServerConfig
	// Store provider configuration
	StoreConfig  *StoreConfig
	ClientConfig *ClientConfig
	CLIConfig    *CLIConfig
	AuthConfig   *auth.Config
	Version      *Version
}

// NewBubblyConfig sets up a default Bubbly Config
func NewBubblyConfig(opts ...func(*BubblyConfig)) *BubblyConfig {
	b := &BubblyConfig{
		Logger:        NewDefaultLogger(),
		ReleaseConfig: DefaultReleaseConfig(),
		ServerConfig:  DefaultServerConfig(),
		StoreConfig:   DefaultStoreConfig(),
		ClientConfig:  DefaultClientConfig(),
		CLIConfig:     DefaultCLIConfig(),
		AuthConfig:    DefaultAuthConfig(),
		Version:       NewVersionInfo(),
	}

	for _, opt := range opts {
		opt(b)
	}
	return b
}

func WithBubblyUI(fs *embed.FS) func(*BubblyConfig) {
	return func(bCtx *BubblyConfig) {
		bCtx.UI = fs
	}
}

func WithVersion(version *Version) func(*BubblyConfig) {
	return func(bCtx *BubblyConfig) {
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
// the zerolog.Logger managed by a BubblyConfig instance
func (bCtx *BubblyConfig) UpdateLogLevel(level zerolog.Level) {
	bCtx.Logger = bCtx.Logger.Level(level)
}
