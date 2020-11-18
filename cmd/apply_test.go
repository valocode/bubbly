package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/verifa/bubbly/env"
	"github.com/verifa/bubbly/server"
	"gopkg.in/h2non/gock.v1"
)

func TestCmdApplyFilename(t *testing.T) {
	bCtx := env.NewBubblyContext()
	cmd, o := NewCmdApply(bCtx)

	cmd.SetArgs([]string{"-f", "../parser/testdata/local-sq-json"})
	cmd.SilenceUsage = true
	cmd.Execute()

	assert.Equal(t, "../parser/testdata/local-sq-json", o.Filename)
}

// Verify that given a set of flag configurations defining the config.ServerConfig, bubbly apply makes a POST request to the correct address
func TestCmdApplyWithServerConfigsSetup(t *testing.T) {
	for _, c := range applyWithServerConfigsSetupCases {
		t.Run(c.desc, func(t *testing.T) {
			bCtx := env.NewBubblyContext()
			bCtx.Config.ServerConfig = c.serverConfig
			router = server.SetupRouter(bCtx)

			// Create a new server route for mocking a Bubbly server response
			gock.New(c.address).
				Post(c.route).
				Reply(c.responseCode).
				JSON(c.response)

			cmd, o := NewCmdApply(bCtx)

			// // set the flag inputs to `bubbly`
			// for name, value := range c.flags {
			// 	_ = rootCmd.PersistentFlags().Set(name, value)
			// }
			cmd.SetArgs([]string{"-f", "../bubbly/testdata/sonarqube"})
			cmd.SilenceUsage = true
			cmd.Execute()

			assert.Equal(t, c.expected, o.Result)
		})
	}
}
