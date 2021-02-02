package store

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/verifa/bubbly/api/core"
)

// some debugging function
func printTree(nodes []*dataNode, depth int) {
	indent := strings.Repeat("\t", depth)
	for _, node := range nodes {
		fmt.Printf("%s%d: %s: %#v\n", indent, depth, node.Data.TableName, node.RefFields)
		printTree(node.Children, depth+1)
	}
}

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
					Fields:    core.DataFields{},
					Data: core.DataBlocks{
						{
							TableName: "root_nested",
							Fields:    core.DataFields{},
						},
					},
				},
				{
					TableName: "root_join",
					Fields:    core.DataFields{},
					Joins:     []string{"root"},
				},
			},
			out: dataTree{
				&dataNode{
					Data: &core.Data{TableName: "root"},
					Children: []*dataNode{
						{Data: &core.Data{TableName: "root_nested"}},
						{Data: &core.Data{TableName: "root_join"}},
					},
				},
			},
		},
	}

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
			err = tree.traverse(func(node *dataNode) error {
				nodeList = append(nodeList, node.Data.TableName)
				return nil
			})
			assert.NoErrorf(t, err, "failed traverse tree")
			err = c.out.traverse(func(node *dataNode) error {
				expNodeList = append(expNodeList, node.Data.TableName)
				return nil
			})
			assert.NoErrorf(t, err, "failed traverse expected tree")
			assert.Equalf(t, expNodeList, nodeList, "data trees not equal")
		})
	}
}
