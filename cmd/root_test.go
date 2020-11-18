package cmd

import (
	"testing"

	"github.com/rs/zerolog"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/verifa/bubbly/config"
	"github.com/verifa/bubbly/env"
)

// TODO: Remove when viper is purged
func TestRootWithPortArg(t *testing.T) {
	t.Run("test root", func(t *testing.T) {
		// Create a new describe command
		bCtx := env.NewBubblyContext()
		rootCmd := NewCmdRoot(bCtx)
		rootCmd.SetArgs([]string{"--port", "4040"})
		viper.BindPFlags(rootCmd.PersistentFlags())
		rootCmd.SilenceUsage = true

		rootCmd.Execute()
		// Verify that viper binds the port argument correctly
		assert.Equal(t, "4040", viper.Get("port"))
	})
}

// TODO: Remove when viper is purged
func TestRootWithDebugArg(t *testing.T) {
	t.Run("test root with debug", func(t *testing.T) {
		// Create a new describe command
		bCtx := env.NewBubblyContext()
		rootCmd := NewCmdRoot(bCtx)
		rootCmd.SetArgs([]string{"--debug"})
		// viper.BindPFlags(rootCmd.PersistentFlags())
		rootCmd.SilenceUsage = true

		rootCmd.Execute()

		// Verify that viper binds the port argument correctly
		assert.Equal(t, true, viper.Get("debug"))
	})
}

// TODO: Remove when viper is purged
func TestRootWithInvalidArg(t *testing.T) {
	t.Run("test root", func(t *testing.T) {
		bCtx := env.NewBubblyContext()
		rootCmd := NewCmdRoot(bCtx)
		// Create a new describe command
		rootCmd.SetArgs([]string{"--port", "4040", "--invalid-arg", "invalid"})
		// rootCmd.SetArgs([]string{"--test", "example"})
		// rootCmd.PersistentFlags().Set("test", "hi")
		viper.BindPFlags(rootCmd.PersistentFlags())
		rootCmd.SilenceUsage = true

		rootCmd.Execute()

		// Verify that viper binds the port argument correctly, and ignores an invalid argument
		assert.Equal(t, "4040", viper.Get("port"))
		assert.Nil(t, viper.Get("invalid-arg"))
	})

}

// TestBubblyContext verifies that rootCmds correctly map global server
// configuration flags like to the env.BubblyContext
func TestBubblyContext(t *testing.T) {
	tcs := []struct {
		desc     string
		inputCtx *env.BubblyContext
		args     []string
		expected *env.BubblyContext
	}{
		{
			desc:     "basic: port",
			inputCtx: env.NewBubblyContext(),
			args: []string{
				"--port", "4040",
			},
			expected: &env.BubblyContext{
				Config: &config.Config{
					ServerConfig: &config.ServerConfig{
						Port:     "4040",
						Protocol: "http",
					},
				},
				Logger: env.NewDefaultLogger(),
			},
		},
		{
			desc:     "basic: log level info",
			inputCtx: env.NewBubblyContext(),
			expected: &env.BubblyContext{
				Config: &config.Config{
					ServerConfig: &config.ServerConfig{
						Protocol: "http",
					},
				},
				Logger: env.NewDefaultLogger(),
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.desc, func(t *testing.T) {

			bCtx := tc.inputCtx
			rootCmd := NewCmdRoot(bCtx)
			rootCmd.SetArgs(tc.args)
			rootCmd.SilenceUsage = true

			rootCmd.Execute()

			// Verify that the argument(s) bind to the BubblyContext
			// correctly
			assert.Equal(t, tc.expected, bCtx)
		})
	}
}

// TestBubblyContextLogLevel verifies that following the creation of a
// new root cmd, the logger configuration of the env.BubblyContext, that we
// need to update as a part of the main.go set up process, is the only
// thing modified.
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
				Config: &config.Config{
					ServerConfig: &config.ServerConfig{
						Protocol: "http",
					},
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
				Config: &config.Config{
					ServerConfig: &config.ServerConfig{
						Protocol: "http",
					},
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
			assert.Equal(t, tc.expected, bCtx)
		})
	}
}
