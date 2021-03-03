package integration

import (
	"errors"
	"testing"

	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"
	"github.com/verifa/bubbly/api/core"
	"github.com/verifa/bubbly/cmd/get"
	"github.com/verifa/bubbly/env"
	"github.com/verifa/bubbly/parser"
	"github.com/zclconf/go-cty/cty"
)

func parseHCLFile(file string, val interface{}) error {
	hclFile, diags := hclparse.NewParser().ParseHCLFile(file)
	if diags != nil {
		return errors.New(diags.Error())
	}

	bCtx := env.NewBubblyContext()
	bCtx.UpdateLogLevel(zerolog.DebugLevel)

	return parser.DecodeExpandBody(bCtx, hclFile.Body, val, cty.NilVal)
}

// TestAutomationData reads and parses the data to be loaded into Bubbly Store
func TestAutomationData(dataFile string) (core.DataBlocks, error) {
	dataWrapper := struct {
		DataBlocks core.DataBlocks `hcl:"data,block"`
	}{}

	err := parseHCLFile(dataFile, &dataWrapper)
	return dataWrapper.DataBlocks, err
}

func TestBubblyCmd(t *testing.T, bCtx *env.BubblyContext, cmdString string, args []string) {
	var cmd *cobra.Command
	switch cmdString {
	case "get":
		cmd, _ = get.NewCmdGet(bCtx)
	}

	cmd.SetArgs(args)
	require.NoError(t, cmd.Execute())
}
