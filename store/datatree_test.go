package store

import (
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/valocode/bubbly/api/core"
	"github.com/valocode/bubbly/env"
	"github.com/valocode/bubbly/parser"
	"github.com/zclconf/go-cty/cty"
)

func TestDataTree(t *testing.T) {
	cases := []struct {
		desc string
		in   core.DataBlocks
		out  dataTree
	}{
		{
			desc: "simple data tree",
			in: core.DataBlocks{
				{
					TableName: "root",
					Fields:    &core.DataFields{Values: map[string]cty.Value{"foo": cty.BoolVal(true)}},
					Data: core.DataBlocks{
						{
							TableName: "root_nested",
							Fields:    &core.DataFields{Values: map[string]cty.Value{"foo": cty.BoolVal(true)}},
						},
					},
				},
				{
					TableName: "root_join",
					Fields:    &core.DataFields{},
					Joins:     []string{"root"},
				},
			},
			out: dataTree{
				&dataNode{
					Data: &core.Data{TableName: "root", Fields: &core.DataFields{Values: map[string]cty.Value{"foo": cty.BoolVal(true)}}},
					Children: []*dataNode{
						{Data: &core.Data{TableName: "root_nested", Fields: &core.DataFields{Values: map[string]cty.Value{"foo": cty.BoolVal(true)}}}},
						{Data: &core.Data{TableName: "root_join"}},
					},
				},
			},
		},
		{
			desc: "cyclic join",
			in: core.DataBlocks{
				{
					TableName: "root",
					Fields:    &core.DataFields{Values: map[string]cty.Value{"foo": cty.BoolVal(true)}},
					Data: core.DataBlocks{
						{
							TableName: "root_nested",
							Fields:    &core.DataFields{Values: map[string]cty.Value{"foo": cty.BoolVal(true)}},
						},
					},
				},
				{
					TableName: "other",
					Fields: &core.DataFields{Values: map[string]cty.Value{
						"_id": cty.CapsuleVal(parser.DataRefType, &parser.DataRef{
							TableName: "root",
							Field:     "other_id",
						}),
					}},
				},
				{
					TableName: "root",
					Fields: &core.DataFields{Values: map[string]cty.Value{
						"_id": cty.CapsuleVal(parser.DataRefType, &parser.DataRef{
							TableName: "root",
							Field:     "_id",
						}),
					}},
					Joins:  []string{"other"},
					Policy: core.UpdatePolicy,
				},
			},
			out: dataTree{
				&dataNode{
					Data: &core.Data{TableName: "root", Fields: &core.DataFields{Values: map[string]cty.Value{"foo": cty.BoolVal(true)}}},
					Children: []*dataNode{
						{Data: &core.Data{TableName: "root_nested", Fields: &core.DataFields{Values: map[string]cty.Value{"foo": cty.BoolVal(true)}}}},
						{Data: &core.Data{TableName: "other"}, Children: []*dataNode{
							{Data: &core.Data{TableName: "root"}},
						}},
					},
				},
			},
		},
	}

	bCtx := env.NewBubblyContext()
	bCtx.UpdateLogLevel(zerolog.DebugLevel)

	for _, c := range cases {
		t.Run(c.desc, func(t *testing.T) {
			tree, err := createDataTree(c.in)
			assert.NoErrorf(t, err, "failed to create data tree")

			var (
				nodeList    = make([]string, 0)
				expNodeList = make([]string, 0)
			)
			// If you want to visualize the tree, uncomment this
			// printTree(tree, 0)

			// The way we handle equality is by traversing the tree and adding
			// each node.Data.TableName to a list, and then doing equality on
			// that list. Otherwise we would need to duplicate the tree and
			// that's a bit of a PITA.
			_, err = tree.traverse(bCtx, func(bCtx *env.BubblyContext, node *dataNode,
				blocks *core.DataBlocks) error {
				t.Logf("traversing: %s", node.Data.TableName)
				nodeList = append(nodeList, node.Data.TableName)
				return nil
			})
			assert.NoErrorf(t, err, "failed traverse tree")
			_, err = c.out.traverse(bCtx, func(bCtx *env.BubblyContext, node *dataNode,
				blocks *core.DataBlocks) error {
				expNodeList = append(expNodeList, node.Data.TableName)
				return nil
			})
			assert.NoErrorf(t, err, "failed traverse expected tree")
			assert.Equalf(t, expNodeList, nodeList, "data trees not equal")
		})
	}
}
