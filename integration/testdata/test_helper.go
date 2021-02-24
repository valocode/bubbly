package integration

import (
	"errors"

	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/rs/zerolog"
	"github.com/verifa/bubbly/api/core"
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
