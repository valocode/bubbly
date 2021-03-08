package store

import (
	"fmt"

	"github.com/verifa/bubbly/api/core"
)

// schemaGraph creates a graph from the bubbly schema
type schemaGraph struct {
	nodes []*schemaNode
	// nodeIndex stores an index to the nodes using the schema table name.
	// This is probably not the best for performance, but our schemas probably
	// will not become huge and this is a great convienice for consumers of the
	// graph, to not have to traverse the graph to find a node
	nodeIndex map[string]*schemaNode
}

func (g *schemaGraph) traverse(fnVisit func(node *schemaNode) error) error {
	var visited = make(map[string]struct{})
	for _, n := range g.nodes {
		if err := visitSchemaNode(n, visited, fnVisit); err != nil {
			return fmt.Errorf("failed to traverse schema graph: %w", err)
		}
	}
	return nil
}

func visitSchemaNode(node *schemaNode, visited map[string]struct{}, fnVisit func(node *schemaNode) error) error {
	if err := fnVisit(node); err != nil {
		return err
	}
	visited[node.table.Name] = struct{}{}

	for _, e := range node.edges {
		// If we have already visited the node, then continue the loop
		if _, ok := visited[e.node.table.Name]; ok {
			continue
		}
		if err := visitSchemaNode(e.node, visited, fnVisit); err != nil {
			return err
		}
	}
	return nil
}

// schemaNode represents a node in the schema graph.
// A node is a wrapper around core.Table with the edges for explicit
// relationships to other nodes (and therefore tables)
type schemaNode struct {
	table *core.Table
	edges []*schemaEdge
}

// shortestPath uses breadth-first search to find the shortest path between
// two nodes in the graph
func (n *schemaNode) shortestPath(dst string) schemaPath {
	var (
		visited   = make(map[string]struct{})
		pathQueue = make([]schemaPath, 0)
	)

	// If the dst node is the root node, then just return an empty path
	if n.table.Name == dst {
		return schemaPath{}
	}

	// Iterate over the root node's edges and create the initial pathQueue
	for _, e := range n.edges {
		if _, ok := visited[e.node.table.Name]; !ok {
			pathQueue = append(pathQueue, schemaPath{e})
			visited[e.node.table.Name] = struct{}{}
		}
	}

	for len(pathQueue) > 0 {
		var (
			// Get the latest path in the queue
			path = pathQueue[0]
			// From the path, get the last node in the path, as we want to
			// traverse that node's edges
			tail = path[len(path)-1]
		)
		if tail.node.table.Name == dst {
			return path
		}
		for _, e := range tail.node.edges {
			if _, ok := visited[e.node.table.Name]; !ok {
				pathQueue = append(pathQueue, append(path, e))
				visited[e.node.table.Name] = struct{}{}
			}
		}
		// Remove first element from the queue
		pathQueue = pathQueue[1:]
	}

	return nil
}

// neighbours takes a distance and returns a slice of paths to all of the nodes
// within that distance in the graph
func (n *schemaNode) neighbours(distance int) []schemaPath {
	var visited = make(map[string]*schemaNode)
	return nodeNeighbours(n, schemaPath{}, visited, distance)
}

// addEdgeFromJoin takes a node and creates bi-directional edges between the
// nodes. Noteworthy is the relationship that the edges describe
func (n *schemaNode) addEdgeFromJoin(child *schemaNode, unique bool) {
	var (
		// This node (parent) has a oneToMany or oneToOne relationship with the
		// child node
		edgeToChild = &schemaEdge{node: child, rel: oneToMany}
		// The child "belongsTo" the parent (this nodes)
		edgeToParent = &schemaEdge{node: n, rel: belongsTo}
	)
	if unique {
		// If unique, then it's a oneToOne relationship, not oneToMany
		edgeToChild.rel = oneToOne
	}
	// Add the edge to the child to this node
	n.edges = append(n.edges, edgeToChild)
	// Also add the reverse relationship
	child.edges = append(child.edges, edgeToParent)
}

// schemaPath stores a slice of edges
type schemaPath []*schemaEdge

// isScalar returns true if the path from one node to another is scalar.
// Scalar means that the return type is a single instance, rather than a slice,
// and is primarily used by the graphql API to return the correct value
func (p *schemaPath) isScalar() bool {
	var isScalar = true
	for _, e := range *p {
		switch e.rel {
		case oneToMany:
			isScalar = false
		default:
			// oneToOne and hasOne both return scalar values, so do nothing
		}
	}
	return isScalar
}

// relType desribes the relationship type of a directed edge from a --> b
type relType int

const (
	oneToOne relType = iota
	oneToMany
	belongsTo
)

// schemaEdge represents an edge in the node graph
type schemaEdge struct {
	node *schemaNode
	rel  relType
}

// newSchemaGraphFromMap is a wrapper around newSchemaGraph for backwards
// compatability with the current way the schema is stored in the provider
func newSchemaGraphFromMap(tables map[string]core.Table) (*schemaGraph, error) {
	var ts = make([]core.Table, 0, len(tables))
	for _, t := range tables {
		ts = append(ts, t)
	}
	return newSchemaGraph(ts)
}

// newSchemaGraph takes a list of tables and creates a schemaGraph
func newSchemaGraph(tables core.Tables) (*schemaGraph, error) {
	var (
		nodes = make(map[string]*schemaNode)
		graph = &schemaGraph{nodeIndex: nodes}
	)

	collectTables(nodes, tables)
	// First iterate over the top-level tables to extract the root nodes in the
	// graph. Tables at the top-level, without any joins, do not have any edges
	// going to them, so they are root nodes.
	for _, table := range tables {
		if len(table.Joins) == 0 {
			graph.nodes = append(graph.nodes, nodes[table.Name])
		}
	}

	if err := tablesToGraph(nodes, tables, nil); err != nil {
		return graph, fmt.Errorf("failed to create graph: %w", err)
	}
	return graph, nil
}

func collectTables(nodes map[string]*schemaNode, tables core.Tables) {
	for index, t := range tables {
		nodes[t.Name] = &schemaNode{table: &tables[index]}
		collectTables(nodes, t.Tables)
	}
}

func tablesToGraph(nodes map[string]*schemaNode, tables core.Tables, parent *schemaNode) error {
	for _, table := range tables {
		var node = nodes[table.Name]
		// Handle the explicit joins
		for _, join := range table.Joins {
			// A join indicates that this table "belongs to" another talbe,
			// i.e. this table is a child of that table
			parent, ok := nodes[join.Table]
			if !ok {
				return fmt.Errorf("join refers to unknown table: %s --> %s", table.Name, join.Table)
			}
			// Create the edge from parent to node
			parent.addEdgeFromJoin(node, join.Single)
		}
		// Handle the implicit joins, i.e. a table nested within a table
		if parent != nil {
			parent.addEdgeFromJoin(node, table.Unique)
		}
		// Recurse
		tablesToGraph(nodes, table.Tables, node)

		// Clear unnecessary data
		table.Tables = nil
		// TODO: cannot remove joins because it breaks schema diff tests
		// table.Joins = nil
	}
	return nil
}

func nodeNeighbours(node *schemaNode, path schemaPath, visited map[string]*schemaNode, remaining int) []schemaPath {
	if remaining == 0 {
		return []schemaPath{}
	}
	visited[node.table.Name] = node

	var paths = make([]schemaPath, 0, 1)
	for _, e := range node.edges {
		if _, ok := visited[e.node.table.Name]; !ok {
			var childPath = append(path, e)
			paths = append(paths, childPath)
			paths = append(paths, nodeNeighbours(e.node, childPath, visited, remaining-1)...)
		}
	}
	return paths
}
