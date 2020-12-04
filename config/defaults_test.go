package config

import (
	"testing"

	"github.com/likexian/gokit/assert"
)

func TestNewDefaultConfig(t *testing.T) {
	tcs := []struct {
		desc     string
		expected *ServerConfig
	}{
		{
			desc: "basic creation of Config from defaults",
			expected: &ServerConfig{
				Protocol: "http",
				Port:     "8111",
				Auth:     false,
				Host:     "localhost",
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.desc, func(t *testing.T) {
			c := NewDefaultConfig()

			assert.Equal(t, tc.expected, c)
		})
	}
}
