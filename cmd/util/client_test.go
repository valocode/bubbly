package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/verifa/bubbly/config"
)

// tests util.ClientSetup
func TestClientSetup(t *testing.T) {
	tcs := []struct {
		desc            string
		input           *config.ServerConfig
		expectedSuccess bool
		expectedErr     string
	}{
		{
			desc: "basic valid config.ServerConfig",
			input: &config.ServerConfig{
				Host: "localhost",
				Port: "8080",
			},
			expectedSuccess: true,
		},
		{
			desc:            "empty config.ServerConfig",
			input:           &config.ServerConfig{},
			expectedErr:     "Unable to create Bubbly client: missing required arguments.",
			expectedSuccess: false,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.desc, func(t *testing.T) {
			_, err := ClientSetup(*tc.input)

			if !tc.expectedSuccess {
				assert.Error(t, err)
				assert.Equal(t, err.Error(), tc.expectedErr)
				t.SkipNow()
			}

			assert.NoError(t, err)

		})
	}
}
