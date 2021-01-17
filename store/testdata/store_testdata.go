package store

import (
	"testing"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/stretchr/testify/require"
	"github.com/verifa/bubbly/api/core"
	"github.com/verifa/bubbly/env"
	"github.com/verifa/bubbly/parser"
	"github.com/zclconf/go-cty/cty"
)

func Tables(t *testing.T) core.Tables {

	// Parse tables
	bCtx := env.NewBubblyContext()
	tableWrapper := struct {
		Tables core.Tables `hcl:"table,block"`
	}{}
	body, err := parser.MergedHCLBodies(bCtx, "testdata/tables.hcl")
	require.NoErrorf(t, err, "failed to parse tables")
	err = parser.DecodeExpandBody(bCtx, body, &tableWrapper, cty.NilVal)
	require.NoErrorf(t, err, "failed to decode tables")
	return tableWrapper.Tables
}

func DataBlocks(t *testing.T) core.DataBlocks {
	// Parse data blocks
	bCtx := env.NewBubblyContext()
	dataWrapper := struct {
		Data core.DataBlocks `hcl:"data,block"`
	}{}
	body, err := parser.MergedHCLBodies(bCtx, "testdata/data.hcl")
	require.NoErrorf(t, err, "failed to parse data blocks")
	err = parser.DecodeExpandBody(bCtx, body, &dataWrapper, cty.NilVal)
	require.NoErrorf(t, err, "failed to decode data blocks")
	return dataWrapper.Data
}
