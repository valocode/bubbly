package store

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/valocode/bubbly/env"
	testData "github.com/valocode/bubbly/store/testdata"
)

func printSchemaGraph(graph *schemaGraph) {
	for _, n := range graph.Nodes {
		printSchemaNode(n, 0)
	}
}

func printSchemaNode(node *schemaNode, depth int) {
	indent := strings.Repeat("\t", depth)
	fmt.Printf("%s%d: %s: %#v\n", indent, depth, node.table.Name, node.edges)
	for _, e := range node.edges {
		printSchemaNode(e.node, depth+1)
	}
}

func TestSchemaGraph(t *testing.T) {

	bCtx := env.NewBubblyContext()
	bCtx.UpdateLogLevel(zerolog.DebugLevel)

	tables := testData.Tables(t, bCtx, filepath.FromSlash("testdata/tables.hcl"))
	graph, err := newSchemaGraph(tables)
	require.NoErrorf(t, err, "failed to create schema graph")

	rootNode := graph.NodeIndex["root"]
	path := rootNode.shortestPath("grandchild_a")
	assert.NotNilf(t, path, "there should be a path between the nodes")
	paths := rootNode.neighbours(2)
	assert.NotEmptyf(t, paths, "the node should have neighbours")
}
