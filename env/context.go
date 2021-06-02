package env

import (
	"fmt"
	"os"

	"github.com/rs/zerolog"
	"github.com/valocode/bubbly/config"
)

// BubblyContext holds global bubbly state that is required
// to be injected into functions throughout the codebase.
type BubblyContext struct {
	Logger       *zerolog.Logger
	AuthConfig   *config.AuthConfig
	ServerConfig *config.ServerConfig
	StoreConfig  *config.StoreConfig
	AgentConfig  *config.AgentConfig
	ClientConfig *config.ClientConfig
	CLIConfig    *config.CLIConfig
	Version      *Version
}

// NewBubblyContext sets up a default Bubbly Context
func NewBubblyContext() *BubblyContext {
	return &BubblyContext{
		Logger:       NewDefaultLogger(),
		AuthConfig:   config.DefaultAuthConfig(),
		ServerConfig: config.DefaultServerConfig(),
		StoreConfig:  config.DefaultStoreConfig(),
		AgentConfig:  config.DefaultAgentConfig(),
		ClientConfig: config.DefaultClientConfig(),
		CLIConfig:    config.DefaultCLIConfig(),
		Version:      NewVersionInfo(),
	}
}

// Version contains the SHA1 value for the commit
// in the Bubbly Git repo from which the current running
// binary was built. If the commit was tagged, the tag
// is also included.
type Version struct {
	SHA1 string
	Tag  string
}

// These will be set by the compiler using "-X env.sha1"
// and "-X env.tag". They have to be simple variables in,
// the current "env" package, because there is no way
// to set the fields of a local struct. So, the compiler
// will set these, and then the values are going to be
// used to create an instance of struct, which can then
// be included into the Bubbly context.
var (
	sha1 string
	tag  string
)

// NewVersionInfo returns a reference to the Version structure,
// populated with information provided at compile time.
func NewVersionInfo() *Version {
	return &Version{
		SHA1: sha1,
		Tag:  tag,
	}
}

// NewDefaultLogger sets up a default logger
func NewDefaultLogger() *zerolog.Logger {

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
