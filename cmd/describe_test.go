package cmd

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/verifa/bubbly/env"
)

// TODO: gock requires explicit client intercept. Rewrite tests.

// func TestDescribeResourceReturn(t *testing.T) {

// 	for _, c := range describeResourceReturnCases {
// 		t.Run(c.desc, func(t *testing.T) {
// 			defer gock.Off() // Flush pending mocks after test execution
// 			hostURL := c.sc.Host + ":" + c.sc.Port
// 			// Create a new server route for mocking a Bubbly server response
// 			gock.New(hostURL).
// 				Get(c.route).
// 				Reply(c.responseCode).
// 				JSON(c.response)

// 			// Create a new describe command
// 			cmd, o := NewCmdDescribe()
// 			b := bytes.NewBufferString("")
// 			cmd.SetOut(b)
// 			cmd.SetArgs([]string{c.rType, c.rName})
// 			cmd.SilenceUsage = true
// 			err := cmd.Execute()

// 			assert.NoError(t, err)

// 			_, err = ioutil.ReadAll(b)

// 			// got := string(out)

// 			assert.NoError(t, err)

// 			assert.Len(t, o.Result, 1)
// 			// assert.Equal(t, got, c.expected)

// 			assert.True(t, gock.IsDone())

// 			gock.Flush()
// 		})
// 	}
// }

// func TestDescribeResourceGroupReturn(t *testing.T) {

// 	for _, c := range describeResourceGroupReturnCases {
// 		t.Run(c.desc, func(t *testing.T) {
// 			defer gock.Off() // Flush pending mocks after test execution
// 			hostURL := c.sc.Host + ":" + c.sc.Port
// 			// Create a new server route for mocking a Bubbly server response
// 			gock.New(hostURL).
// 				Get(c.route).
// 				Reply(c.responseCode).
// 				JSON(c.response)

// 			// Create a new describe command
// 			cmd, _ := NewCmdDescribe()
// 			b := bytes.NewBufferString("")
// 			cmd.SetOut(b)
// 			cmd.SetArgs([]string{c.rType})
// 			cmd.SilenceUsage = true
// 			cmd.Execute()

// 			out, err := ioutil.ReadAll(b)

// 			got := string(out)

// 			assert.NoError(t, err)

// 			// To keep it simple, just check that all expected resources are present in the output
// 			for _, contain := range c.expectedContains {
// 				assert.Contains(t, got, contain)
// 			}

// 			assert.True(t, gock.IsDone())
// 		})
// 	}
// }

//TODO: convert to integration test that validates retrieving information on a
// valid resource.
func TestDescribeValidResourceType(t *testing.T) {
	for _, c := range describeValidResourceCases {
		t.Run(c.desc, func(t *testing.T) {
			bCtx := env.NewBubblyContext()
			cmd, _ := NewCmdDescribe(bCtx)
			b := bytes.NewBufferString("")
			// cmd.SetOut(b)
			cmd.SetErr(b)
			cmd.SetArgs([]string{c.rType, c.rName})
			// Silence usage output as that makes testing error outputs more painful
			cmd.SilenceUsage = true
			err := cmd.Execute()

			out, err := ioutil.ReadAll(b)
			assert.NoError(t, err)

			got := string(out)

			assert.Contains(t, got, c.expectedContains)
		})
	}
}

func TestDescribeInvalidResourceType(t *testing.T) {
	for _, c := range describeInvalidResourceCases {
		t.Run(c.desc, func(t *testing.T) {
			bCtx := env.NewBubblyContext()
			cmd, _ := NewCmdDescribe(bCtx)
			b := bytes.NewBufferString("")
			// cmd.SetOut(b)
			cmd.SetErr(b)
			cmd.SetArgs([]string{c.rType, c.rName})
			// Silence usage output as that makes testing error outputs more painful
			cmd.SilenceUsage = true
			cmd.Execute()

			out, err := ioutil.ReadAll(b)
			assert.NoError(t, err)

			got := string(out)

			assert.Equal(t, got, c.expected)
		})
	}
}

func TestDescribeHelpMessage(t *testing.T) {
	bCtx := env.NewBubblyContext()
	cmd, _ := NewCmdDescribe(bCtx)
	b := bytes.NewBufferString("")
	cmd.SetOut(b)
	cmd.SetArgs([]string{"-h"})
	cmd.SilenceUsage = true
	cmd.Execute()
	// Silence usage output as that makes testing error outputs more painful

	out, err := ioutil.ReadAll(b)

	assert.NoError(t, err)

	expected := `describe (TYPE [NAME_PREFIX] | TYPE/NAME) [flags]`

	got := string(out)

	assert.Contains(t, got, expected)

}

func TestDescribeCmdSetup(t *testing.T) {
	bCtx := env.NewBubblyContext()
	cmd, _ := NewCmdDescribe(bCtx)
	assert.NotNil(t, cmd)
}

func TestDescribeWithVersionArg(t *testing.T) {

	for _, c := range describeWithVersionArgCases {
		t.Run(c.desc, func(t *testing.T) {
			bCtx := env.NewBubblyContext()
			// Create a new describe command
			cmd, o := NewCmdDescribe(bCtx)
			cmd.SetOut(bytes.NewBufferString(""))
			cmd.SetArgs([]string{c.rType, c.rName, "--version", c.version})
			cmd.SilenceUsage = true
			cmd.Execute()

			// Verify that the version flag is bound by viper correctly.
			assert.Equal(t, c.expected, viper.Get("version"))

			// Alternatively, if we remove Viper we can assert that the
			// variable associated to the "version" flag was assigned
			// correctly by Cobra:
			assert.Equal(t, c.expected, o.Version)
		})
	}
}

// Verify that ServerConfigs are populated from rootCmd inputs
func TestDescribeWithServerConfigsSetup(t *testing.T) {
	for _, c := range describeWithServerConfigsSetupCases {
		t.Run(c.desc, func(t *testing.T) {
			bCtx := env.NewBubblyContext()
			rootCmd := NewCmdRoot(bCtx)
			// Create a new describe command
			cmd, o := NewCmdDescribe(bCtx)
			// set the flag inputs to `bubbly`
			for name, value := range c.flags {
				_ = rootCmd.PersistentFlags().Set(name, value)
			}

			cmd.SetOut(bytes.NewBufferString(""))
			cmd.SetArgs([]string{c.rType, c.rName})
			cmd.SilenceUsage = true
			cmd.Execute()

			assert.Equal(t, c.expected.ServerConfig, o.BubblyContext.Config.ServerConfig)
		})
	}
}
