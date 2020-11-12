package config

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// tests setting up of Config from a mix of viper bindings and defaults
func TestNewConfig(t *testing.T) {
	tcs := []struct {
		desc     string
		flags    map[string]string
		expected *Config
	}{
		{
			desc: "basic creation of Config from viper bindings",
			flags: map[string]string{
				"protocol": "http",
				"port":     "8070",
				"host":     "localhost",
				"auth":     "false",
				"token":    "",
			},
			expected: &Config{
				ServerConfig: &ServerConfig{
					Protocol: "http",
					Port:     "8070",
					Auth:     false,
					Host:     "localhost",
					Token:    "",
				},
			},
		},
		{
			desc: "basic creation of Config from viper bindings and defaults",
			flags: map[string]string{
				"protocol": "http",
				"port":     "8070",
				"auth":     "false",
				"token":    "",
			},
			expected: &Config{
				ServerConfig: &ServerConfig{
					Protocol: "http",
					Port:     "8070",
					Auth:     false,
					Host:     "localhost",
					Token:    "",
				},
			},
		},
		{
			desc: "basic creation of Config from limited viper bindings",
			flags: map[string]string{
				"protocol": "https",
				"port":     "8070",
			},
			// Note the lack of defaults in the expected Config:
			expected: &Config{
				ServerConfig: &ServerConfig{
					Protocol: "https",
					Port:     "8070",
				},
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.desc, func(t *testing.T) {
			t.Parallel()

			flagSet := pflag.NewFlagSet("test", pflag.ContinueOnError)

			for k, v := range tc.flags {
				flagSet.String(k, v, "test")
			}

			viper.BindPFlags(flagSet)
			c := NewConfig()

			assert.Equal(t, tc.expected.ServerConfig, c.ServerConfig)
		})
	}
}

// tests setting up of Config from a mix of viper bindings and defaults
// using config.SetupConfigs
func TestSetupConfigs(t *testing.T) {
	tcs := []struct {
		desc     string
		flags    map[string]string
		expected *Config
	}{
		{
			desc: "basic creation of Config from viper bindings",
			flags: map[string]string{
				"protocol": "http",
				"port":     "8070",
				"host":     "localhost",
				"auth":     "false",
				"token":    "",
			},
			expected: &Config{
				ServerConfig: &ServerConfig{
					Protocol: "http",
					Port:     "8070",
					Auth:     false,
					Host:     "localhost",
					Token:    "",
				},
			},
		},
		{
			desc: "basic creation of Config from viper bindings and defaults",
			flags: map[string]string{
				"protocol": "http",
				"port":     "8070",
				"auth":     "false",
				"token":    "",
			},
			expected: &Config{
				ServerConfig: &ServerConfig{
					Protocol: "http",
					Port:     "8070",
					Auth:     false,
					Host:     "localhost",
					Token:    "",
				},
			},
		},
		{
			desc: "basic creation of Config from viper bindings and defaults",
			flags: map[string]string{
				"protocol": "https",
				"port":     "8070",
			},
			// Note the inclusion of defaults in the expected Config
			// due to the merge from mergo.Merge
			expected: &Config{
				ServerConfig: &ServerConfig{
					Protocol: "https",
					Port:     "8070",
					Auth:     false,
					Host:     "localhost",
					Token:    "",
				},
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.desc, func(t *testing.T) {
			t.Parallel()

			flagSet := pflag.NewFlagSet("test", pflag.ContinueOnError)

			for k, v := range tc.flags {
				flagSet.String(k, v, "test")
			}

			viper.BindPFlags(flagSet)
			c, err := SetupConfigs()

			assert.NoError(t, err)

			assert.Equal(t, c, tc.expected)
		})
	}
}
