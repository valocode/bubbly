package store

import (
	"testing"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/stretchr/testify/require"
	"github.com/valocode/bubbly/api/core"
	"github.com/valocode/bubbly/env"
	"github.com/valocode/bubbly/parser"
	"github.com/zclconf/go-cty/cty"
)

func Tables(t *testing.T, fromFile string) core.Tables {
	t.Helper()

	bCtx := env.NewBubblyContext()
	tableWrapper := struct {
		Tables core.Tables `hcl:"table,block"`
	}{}
	body, err := parser.MergedHCLBodies(bCtx, fromFile)
	require.NoErrorf(t, err, "failed to parse tables")
	err = parser.DecodeExpandBody(bCtx, body, &tableWrapper, cty.NilVal)
	require.NoErrorf(t, err, "failed to decode tables")
	return tableWrapper.Tables
}

func DataBlocks(t *testing.T, fromFile string) core.DataBlocks {
	t.Helper()

	bCtx := env.NewBubblyContext()
	dataWrapper := struct {
		Data core.DataBlocks `hcl:"data,block"`
	}{}
	body, err := parser.MergedHCLBodies(bCtx, fromFile)
	require.NoErrorf(t, err, "failed to parse data blocks")
	err = parser.DecodeExpandBody(bCtx, body, &dataWrapper, cty.NilVal)
	require.NoErrorf(t, err, "failed to decode data blocks")
	return dataWrapper.Data
}
