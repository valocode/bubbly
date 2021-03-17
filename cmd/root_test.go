package cmd

import (
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/valocode/bubbly/config"
	"github.com/valocode/bubbly/env"
)

func TestAgentConfig(t *testing.T) {
	tcs := []struct {
		desc     string
		inputCtx *env.BubblyContext
		args     []string
		expected *env.BubblyContext
	}{
		{
			desc:     "basic: default agent configuration is unmodified after root command execution",
			inputCtx: env.NewBubblyContext(),
			args:     []string{},
			expected: &env.BubblyContext{
				AgentConfig: config.DefaultAgentConfig(),
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.desc, func(t *testing.T) {

			require.Equal(t, tc.expected.AgentConfig, tc.inputCtx.AgentConfig)
			bCtx := tc.inputCtx
			rootCmd := NewCmdRoot(bCtx)
			rootCmd.SetArgs(tc.args)
			rootCmd.SilenceUsage = true

			rootCmd.Execute()

			require.Equal(t, tc.expected.AgentConfig, bCtx.AgentConfig)
		})
	}
}

// TestBubblyContextLogLevel verifies that following the creation of a
// new root cmd, the bubbly logger is set to the correct log level
func TestBubblyContextLogLevel(t *testing.T) {
	tcs := []struct {
		desc          string
		inputCtx      *env.BubblyContext
		inputLogLevel zerolog.Level
		args          []string
		expected      *env.BubblyContext
	}{
		{
			desc:     "basic: log level info",
			inputCtx: env.NewBubblyContext(),
			expected: &env.BubblyContext{
				ServerConfig: &config.ServerConfig{
					Protocol: "http",
				},
				Logger: env.NewDefaultLogger(),
			},
			inputLogLevel: zerolog.InfoLevel,
		},
		{
			desc:     "basic: log level debug",
			inputCtx: env.NewBubblyContext(),
			args: []string{
				"--debug",
			},
			expected: &env.BubblyContext{
				ServerConfig: &config.ServerConfig{
					Protocol: "http",
				},
				Logger: env.NewDebugLogger(),
			},
			inputLogLevel: zerolog.DebugLevel,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.desc, func(t *testing.T) {

			bCtx := tc.inputCtx
			rootCmd := NewCmdRoot(bCtx)
			rootCmd.SetArgs(tc.args)
			rootCmd.SilenceUsage = true

			rootCmd.Execute()

			bCtx.UpdateLogLevel(tc.inputLogLevel)

			// Verify that the argument(s) bind to the BubblyContext
			// correctly
			assert.Equal(t, tc.expected.Logger, bCtx.Logger)
		})
	}
}
