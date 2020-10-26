package cmd

import (
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestRootWithPortArg(t *testing.T) {
	t.Run("test root", func(t *testing.T) {
		// Create a new describe command
		rootCmd.SetArgs([]string{"--port", "4040"})
		viper.BindPFlags(rootCmd.PersistentFlags())
		rootCmd.SilenceUsage = true

		rootCmd.Execute()

		// spew.Dump(o.Config)
		// Verify that viper binds the port argument correctly
		assert.Equal(t, "4040", viper.Get("port"))
	})

}
func TestRootWithInvalidArg(t *testing.T) {
	t.Run("test root", func(t *testing.T) {
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
