package integration

import (
	"errors"

	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/rs/zerolog"
	"github.com/verifa/bubbly/api/core"
	"github.com/verifa/bubbly/env"
	"github.com/verifa/bubbly/parser"
)

func parseHCLFile(file string, val interface{}) error {
	hclFile, diags := hclparse.NewParser().ParseHCLFile(file)
	if diags != nil {
		return errors.New(diags.Error())
	}

	bCtx := env.NewBubblyContext()
	bCtx.UpdateLogLevel(zerolog.DebugLevel)

	s := parser.NewScope()
	return s.DecodeExpandBody(bCtx, hclFile.Body, val)
}

func TestSchema(baseDir string) (core.Tables, error) {

	tableWrapper := struct {
		Tables core.Tables `hcl:"table,block"`
	}{}

	err := parseHCLFile(baseDir+"/testdata/schema/schema.bubbly", &tableWrapper)
	return tableWrapper.Tables, err
}

func TestAutomationData(baseDir string) (core.DataBlocks, error) {
	dataWrapper := struct {
		DataBlocks core.DataBlocks `hcl:"data,block"`
	}{}

	err := parseHCLFile(baseDir+"/testdata/testautomation/data.bubbly", &dataWrapper)
	return dataWrapper.DataBlocks, err
}
