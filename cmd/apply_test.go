package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/verifa/bubbly/server"
	"gopkg.in/h2non/gock.v1"
)

func TestCmdApplyFilename(t *testing.T) {
	cmd, o := NewCmdApply()

	cmd.SetArgs([]string{"-f", "../parser/testdata/local-sq-json"})
	cmd.SilenceUsage = true
	cmd.Execute()

	assert.Equal(t, "../parser/testdata/local-sq-json", o.Filename)
}

// Verify that given a set of flag configurations defining the config.ServerConfig, bubbly apply makes a POST request to the correct address
func TestCmdApplyWithServerConfigsSetup(t *testing.T) {
	for _, c := range applyWithServerConfigsSetupCases {
		t.Run(c.desc, func(t *testing.T) {
			router = server.SetupRouter()

			hostURL := c.flags["host"] + ":" + c.flags["port"]
			// Create a new server route for mocking a Bubbly server response
			gock.New(hostURL).
				Post(c.route).
				Reply(c.responseCode).
				JSON(c.response)

			cmd, o := NewCmdApply()

			// set the flag inputs to `bubbly`
			for name, value := range c.flags {
				_ = rootCmd.PersistentFlags().Set(name, value)
			}
			cmd.SetArgs([]string{"-f", "../parser/testdata/sonarqube"})
			cmd.SilenceUsage = true
			cmd.Execute()

			assert.True(t, o.Result)
		})
	}
}
