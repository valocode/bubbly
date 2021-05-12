package integration

import (
	"errors"
	"testing"

	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"
	"github.com/valocode/bubbly/api/core"
	"github.com/valocode/bubbly/cmd/get"
	"github.com/valocode/bubbly/env"
	"github.com/valocode/bubbly/parser"
	"github.com/zclconf/go-cty/cty"
)

func parseHCLFile(file string, val interface{}) error {
	hclFile, diags := hclparse.NewParser().ParseHCLFile(file)
	if diags != nil {
		return errors.New(diags.Error())
	}

	bCtx := env.NewBubblyContext()
	bCtx.UpdateLogLevel(zerolog.DebugLevel)

	return parser.DecodeExpandBody(hclFile.Body, val, cty.NilVal)
}

// readTestAutomationData reads and parses the data to be loaded into Bubbly Store
func readTestAutomationData(dataFile string) (core.DataBlocks, error) {
	dataWrapper := struct {
		DataBlocks core.DataBlocks `hcl:"data,block"`
	}{}

	err := parseHCLFile(dataFile, &dataWrapper)
	return dataWrapper.DataBlocks, err
}

// FIXME: not sure a Test... case can be in testdata folder?
func TestBubblyCmd(t *testing.T, bCtx *env.BubblyContext, cmdString string, args []string) {
	t.Helper()

	var cmd *cobra.Command
	switch cmdString {
	case "get":
		cmd, _ = get.NewCmdGet(bCtx)
	}

	cmd.SetArgs(args)
	require.NoError(t, cmd.Execute())
}
