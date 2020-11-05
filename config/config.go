package config

import (
	"errors"

	"github.com/imdario/mergo"
	"github.com/spf13/viper"
)

// NewServerConfig creates a ServerConfig struct from viper bindings
func (c *Config) NewServerConfig() *ServerConfig {
	return &ServerConfig{
		Port:  viper.GetString("port"),
		Host:  viper.GetString("host"),
		Auth:  viper.GetBool("auth"),
		Token: viper.GetString("token"),
	}
}

// NewConfig creates a Config struct from viper bindings
func NewConfig() *Config {
	c := &Config{}
	c.ServerConfig = c.NewServerConfig()
	return c
}

// NewDefaultConfig creates a Config struct from default configurations
func NewDefaultConfig() *Config {
	return &Config{
		ServerConfig: &ServerConfig{
			Protocol: "http",
			Port:     "8080",
			Auth:     false,
			Host:     "localhost",
		},
	}
}

// SetupConfigs creates a merged Config struct from defaults and viper bindings.
func SetupConfigs() (*Config, error) {
	defaultConfig := NewDefaultConfig()

	actualConfig := NewConfig()
	// spew.Dump("Default config: ", defaultConfig, "Actual config:", actualConfig)

	// merge default configuration on top of actual configuration
	if err := mergo.Merge(actualConfig, defaultConfig); err != nil {
		return nil, errors.New(err.Error())
	}
	// spew.Dump("new default config: ", defaultConfig, "new actual config:", actualConfig)
	return actualConfig, nil
}
