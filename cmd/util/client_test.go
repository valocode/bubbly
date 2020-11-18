package util

import (
	"testing"

	"github.com/verifa/bubbly/env"

	"github.com/stretchr/testify/assert"
	"github.com/verifa/bubbly/config"
)

// tests util.ClientSetup
func TestClientSetup(t *testing.T) {
	tcs := []struct {
		desc            string
		input           *env.BubblyContext
		expectedSuccess bool
		expectedErr     string
	}{
		{
			desc: "basic valid config.ServerConfig",
			input: &env.BubblyContext{
				Config: &config.Config{
					ServerConfig: &config.ServerConfig{
						Host: "localhost",
						Port: "8080",
					},
				},
			},
			expectedSuccess: true,
		},
		{
			desc: "empty config.ServerConfig",
			input: &env.BubblyContext{
				Config: &config.Config{
					ServerConfig: &config.ServerConfig{},
				},
			},
			expectedErr:     "Unable to create Bubbly client: missing required arguments.",
			expectedSuccess: false,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.desc, func(t *testing.T) {
			bCtx := env.NewBubblyContext()

			bCtx.Config.ServerConfig = tc.input.Config.ServerConfig
			_, err := ClientSetup(bCtx)

			if !tc.expectedSuccess {
				assert.Error(t, err)
				assert.Equal(t, err.Error(), tc.expectedErr)
				t.SkipNow()
			}

			assert.NoError(t, err)

		})
	}
}
