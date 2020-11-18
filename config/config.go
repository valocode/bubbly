package config

import (
	"errors"

	"github.com/imdario/mergo"
	"github.com/spf13/viper"
)

// NewServerConfig creates a ServerConfig struct from viper bindings
func (c *Config) newServerConfig() *ServerConfig {
	return &ServerConfig{
		Protocol: viper.GetString("protocol"),
		Port:     viper.GetString("port"),
		Host:     viper.GetString("host"),
		Auth:     viper.GetBool("auth"),
		Token:    viper.GetString("token"),
	}
}

// NewConfig creates a Config struct from viper bindings
func NewConfig() *Config {
	c := &Config{}
	c.ServerConfig = c.newServerConfig()
	return c
}

// SetupConfigs creates a merged Config struct from defaults and viper bindings.
func SetupConfigs() (*Config, error) {
	defaultConfig := NewDefaultConfig()

	actualConfig := NewConfig()

	// merge default configuration on top of actual configuration
	if err := mergo.Merge(actualConfig, defaultConfig); err != nil {
		return nil, errors.New(err.Error())
	}
	return actualConfig, nil
}
