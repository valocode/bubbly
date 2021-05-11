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

func Tables(t *testing.T, bCtx *env.BubblyContext, fromFile string) core.Tables {
	t.Helper()

	tableWrapper := struct {
		Tables core.Tables `hcl:"table,block"`
	}{}

	body, err := parser.MergedHCLBodies(bCtx, []string{fromFile})
	require.NoErrorf(t, err, "failed to parse tables")

	err = parser.DecodeBody(body, &tableWrapper, cty.NilVal)
	require.NoErrorf(t, err, "failed to decode tables")

	return tableWrapper.Tables
}

func DataBlocks(t *testing.T, bCtx *env.BubblyContext, fromFile string) core.DataBlocks {
	t.Helper()

	dataWrapper := struct {
		Data core.DataBlocks `hcl:"data,block"`
	}{}

	body, err := parser.MergedHCLBodies(bCtx, []string{fromFile})
	require.NoErrorf(t, err, "failed to parse data blocks")

	err = parser.DecodeExpandBody(body, &dataWrapper, cty.NilVal)
	require.NoErrorf(t, err, "failed to decode data blocks")

	return dataWrapper.Data
}
