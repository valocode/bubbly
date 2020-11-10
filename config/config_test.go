package config

import (
	"testing"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/likexian/gokit/assert"
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
			desc: "basic creation of Config from viper bindings and defaults",
			flags: map[string]string{
				"protocol": "https",
				"port":     "8070",
			},
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

			flagSet := pflag.NewFlagSet("test", pflag.ContinueOnError)

			for k, v := range tc.flags {
				flagSet.String(k, v, "test")
			}

			viper.BindPFlags(flagSet)
			c := NewConfig()

			assert.Equal(t, c.ServerConfig, tc.expected.ServerConfig)
		})
	}
}
