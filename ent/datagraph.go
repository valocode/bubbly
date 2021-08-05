// Code written by human, you are free to modify

package ent

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/ext/dynblock"
	"github.com/zclconf/go-cty/cty"
)

type DataGraph struct {
	RootNodes []*DataNode
}

type DataNode struct {
	Name      string
	Fields    []DataField
	Edges     []DataEdge
	Operation NodeOperation

	Resolved bool
	Value    interface{}
}

type DataField struct {
	Name  string
	Value cty.Value
}

type DataEdge struct {
	Name    string
	Node    *DataNode
	Inverse bool // to if false, from if true
}

func (d *DataNode) Graph() *DataGraph {
	return &DataGraph{
		RootNodes: []*DataNode{
			d,
		},
	}
}

func (d *DataNode) AddField(name string, value cty.Value) *DataNode {
	d.Fields = append(d.Fields, DataField{
		Name:  name,
		Value: value,
	})
	return d
}

func (d *DataNode) AddEdge(name string, node *DataNode) *DataNode {
	d.Edges = append(d.Edges, DataEdge{
		Name: name,
		Node: node,
	})
	return d
}

func (d *DataNode) AddInverseEdge(name string, node *DataNode) *DataNode {
	d.Edges = append(d.Edges, DataEdge{
		Name:    name,
		Node:    node,
		Inverse: true,
	})
	return d
}

// Save is traverses the graph and saves the nodes using the provided client
func (g *DataGraph) Save(client *Client) error {
	ctx := context.Background()
	tx, err := client.Tx(ctx)
	if err != nil {
		return fmt.Errorf("error creating transaction: %w", err)
	}
	txErr := g.Traverse(func(node *DataNode) error {
		if err := ProcessDataNode(tx, node); err != nil {
			return err
		}
		return nil
	})

	if txErr != nil {
		if rerr := tx.Rollback(); rerr != nil {
			return fmt.Errorf("error saving data graph triggered a rollback, which also failed: %w", err)
		}
		return txErr
	}
	return tx.Commit()
}

// traverse takes a callback function fn and visits each node in the tree
func (g *DataGraph) Traverse(fn visitFn) error {
	for _, n := range g.RootNodes {
		if err := visitNode(n, fn); err != nil {
			return fmt.Errorf("error occurred when traversing data tree: %w", err)
		}
	}
	return nil
}

// visitNode is a wrapper around the callback visit function to make sure that
// all the nodes are visited at least once, and at most once
func visitNode(node *DataNode, fn visitFn) error {
	// Don't visit this node twice
	if node.Resolved {
		return nil
	}
	// fmt.Println("visiting: " + node.Name)
	// First visit the inverse edges, which will then visit this node again
	for _, edge := range node.Edges {
		if edge.Inverse && !edge.Node.Resolved {
			// fmt.Println("visiting edge: " + node.Name + " -- " + edge.Name + " --> " + edge.Node.Name)
			return visitNode(edge.Node, fn)
		}
	}
	// Visit the node with the callback method
	if err := fn(node); err != nil {
		return fmt.Errorf("error visiting data node %s: %w", node.Name, err)
	}
	// If no error, mark the node as visited
	node.Resolved = true
	for _, edge := range node.Edges {
		if edge.Inverse {
			continue
		}
		if err := visitNode(edge.Node, fn); err != nil {
			return fmt.Errorf("error visiting data node %s: %w", edge.Node.Name, err)
		}
	}
	return nil
}

func (d *DataGraph) Print() {
	for _, n := range d.RootNodes {
		printDataGraph(n, 0)
	}
}

func printDataGraph(n *DataNode, depth int) {
	fmt.Printf("%s%s\n", strings.Repeat("  ", depth), n.Name)

	for _, e := range n.Edges {
		if e.Inverse {
			fmt.Printf("%sInverse Edge %s --> %s\n", strings.Repeat("  ", depth), e.Name, e.Node.Name)
			continue
		}
		printDataGraph(e.Node, depth+1)
	}
}

type DataNodeBody struct {
	ForLoop   *ForLoop `hcl:"for,block"`
	Leftovers hcl.Body `hcl:",remain"`
}

type ForLoop struct {
	Values   cty.Value      `hcl:"values,attr"`
	Iterator hcl.Expression `hcl:"iterator,optional"`
}

func (d *DataNodeBody) Body() hcl.Body {
	return d.Leftovers
}

func (d *DataNodeBody) ForEach(fn func() hcl.Diagnostics) hcl.Diagnostics {
	diags := fn()
	if diags.HasErrors() {
		return diags
	}
	return nil
}

func (d *DataNodeBody) Values() []cty.Value {
	if d.ForLoop == nil {
		return nil
	}
	return nil
}

// visitFn is a type declaration for the tree traversal visit callback function
type visitFn func(node *DataNode) error

type DecodeError struct {
	Diags hcl.Diagnostics
}

func (e *DecodeError) Error() string {
	// TODO: do something more meaningful here...
	return e.Diags.Error()
}

func DecodeDataGraph(body hcl.Body, eCtx *hcl.EvalContext) (*DataGraph, error) {
	var (
		graph DataGraph
		nCtx  = make(map[string]*DataNode)
		err   error
	)
	dynBody := dynblock.Expand(body, eCtx)
	graph.RootNodes, err = decodeDataNodes(dynBody, eCtx, nCtx)
	if err != nil {
		return nil, err
	}
	return &graph, nil
}

func decodeDataNodes(body hcl.Body, eCtx *hcl.EvalContext, nCtx map[string]*DataNode) ([]*DataNode, error) {
	content, diags := body.Content(AllDataNodesSchema.Schema())
	diags = append(diags, validateBlocks(body, content.Blocks, &AllDataNodesSchema)...)
	if diags.HasErrors() {
		return nil, &DecodeError{
			Diags: diags,
		}
	}

	var (
		nodes []*DataNode
		errs  *multierror.Error
	)
	for _, block := range content.Blocks {
		node, err := decodeDataNode(block, eCtx, nCtx, &nodes)
		if err != nil {
			errs = multierror.Append(errs, err)
		}
		nodes = append(nodes, node)
	}
	if errs != nil {
		return nil, errs
	}
	return nodes, nil
}
func decodeDataNode(block *hcl.Block, eCtx *hcl.EvalContext, nCtx map[string]*DataNode, rootNodes *[]*DataNode) (*DataNode, error) {
	// It can happen if there is an empty block that the body is nil.
	// In this case return a nil node and diags
	if block.Body == nil {
		return nil, nil
	}
	// TODO: FORLOOP
	// body, diags := decodeForLoop(block.Body, eCtx)
	// if diags.HasErrors() {
	// 	return nil, diags
	// }

	var (
		schema, schemaExists = DataNodeSchemas[block.Type]
		node                 = DataNode{Name: block.Type}
	)
	if !schemaExists {
		return nil, &DecodeError{
			Diags: hcl.Diagnostics{
				&hcl.Diagnostic{
					Severity: hcl.DiagError,
					Summary:  fmt.Sprintf("Unknown block: %s", block.Type),
					Detail:   fmt.Sprintf("Block is not a known schema type: %s", block.Type),
					Subject:  block.Body.MissingItemRange().Ptr(),
				},
			},
		}
	}
	content, diags := block.Body.Content(schema.Schema())
	diags = append(diags, validateBlocks(block.Body, content.Blocks, schema)...)
	// Add to the node context so that it can be referenced
	nCtx[block.Type] = &node
	// Set a default node operation
	opAttr, ok := content.Attributes[attributeNodeOperation]
	if ok {
		value, exprDiags := opAttr.Expr.Value(eCtx)
		if value.Type() != cty.String {
			diags = append(diags, &hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  fmt.Sprintf("Operation must be a string not %s", value.Type().GoString()),
				Detail:   fmt.Sprintf("Operation must be a string not %s", value.Type().GoString()),
				Subject:  opAttr.Expr.Range().Ptr(),
			})
		} else {
			// TODO: validate NodeOperation
			node.Operation = NodeOperation(value.AsString())
		}
		// Delete the attribute so that it doesn't get added as a field
		delete(content.Attributes, attributeNodeOperation)
		diags = append(diags, exprDiags...)
	} else {
		node.Operation = NodeOperation(DefaultNodeOperation)
	}
	for _, attr := range content.Attributes {
		var (
			exprDiags hcl.Diagnostics
			value     cty.Value
		)
		value, exprDiags = attr.Expr.Value(eCtx)
		node.AddField(attr.Name, value)
		diags = append(diags, exprDiags...)
	}
	var errs error
	for _, nodeBlock := range content.Blocks {
		if nodeBlock.Type == blockNodeDisjoint {

			// Add disjoint nodes as root nodes, because they don't have a "parent"
			disjointNodes, err := decodeDataNodes(nodeBlock.Body, eCtx, nCtx)
			*rootNodes = append(*rootNodes, disjointNodes...)
			errs = multierror.Append(errs, err)
			continue
		}

		// Other blocks are edges which should be decoded
		errs = multierror.Append(decodeDataEdge(nodeBlock, &node, eCtx, nCtx, rootNodes))
	}
	if diags.HasErrors() {
		errs = multierror.Append(errs, &DecodeError{Diags: diags})
	}
	if errs != nil {
		return nil, errs
	}
	return &node, nil
}

func decodeDataEdge(block *hcl.Block, fromNode *DataNode, eCtx *hcl.EvalContext, nCtx map[string]*DataNode, rootNodes *[]*DataNode) error {
	schema, schemaOk := DataEdgeSchemas[fromNode.Name][block.Type]
	if !schemaOk {
		return &DecodeError{
			Diags: hcl.Diagnostics{
				&hcl.Diagnostic{
					Severity: hcl.DiagError,
					Summary:  fmt.Sprintf("Unknown data edge block: %s.%s", fromNode.Name, block.Type),
					Detail:   fmt.Sprintf("Unknown data edge block: %s.%s", fromNode.Name, block.Type),
					Subject:  block.Body.MissingItemRange().Ptr(),
				},
			},
		}
	}
	var errs error
	content, diags := block.Body.Content(schema.Schema())
	diags = append(diags, validateBlocks(block.Body, content.Blocks, schema)...)

	// Check for node references and if they exist, create an edge
	refExpr, refExists := content.Attributes[attributeNodeRef]
	if refExists {
		refVal, refDiags := refExpr.Expr.Value(eCtx)
		diags = append(diags, refDiags...)
		if refVal.Type() != cty.Bool {
			diags = append(diags, &hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  fmt.Sprintf("Attribute %s value must be a boolean", attributeNodeRef),
				Detail:   fmt.Sprintf("Attribute %s value must be a boolean", attributeNodeRef),
				Subject:  block.Body.MissingItemRange().Ptr(),
			})
		}

		if refVal.Type() == cty.Bool && refVal.True() {
			// By the schema design, the edge spec has ONLY ONE block that is accepted
			// because an edge always points to a single type. So we can do this
			// kind of risky operation because it's by design
			var (
				edgeNodeName = schema.Blocks[0].BlockHeaderSchema.Type
				edgeNode, ok = nCtx[edgeNodeName]
			)
			if !ok {
				diags = append(diags, &hcl.Diagnostic{
					Severity: hcl.DiagError,
					Summary:  fmt.Sprintf("Reference to node that has not been defined int the graph: %s", edgeNodeName),
					Detail:   fmt.Sprintf("Cannot reference a data node that has not been created by earlier in the data graph. Either create an instance of the data node here, or earlier in the graph for node %s", edgeNodeName),
					Subject:  block.Body.MissingItemRange().Ptr(),
				})
			} else {
				diags = append(diags, createInverseEdge(block, fromNode, edgeNode)...)
			}
		}
	}

	for _, edgeBlock := range content.Blocks {
		var (
			edgeNode *DataNode
			err      error
		)
		edgeNode, err = decodeDataNode(edgeBlock, eCtx, nCtx, rootNodes)
		if err != nil {
			errs = multierror.Append(errs, err)
			continue
		}
		// If the edge is required, create an inverse edge instead to make sure
		// the nodes get resolved in the correct order and we don't create a node
		// without a required edge
		if _, required := RequiredEdges[fromNode.Name][block.Type]; required {
			diags = append(diags, createInverseEdge(block, fromNode, edgeNode)...)
		} else {
			diags = append(diags, createEdge(block, fromNode, edgeNode)...)
		}
	}

	if diags.HasErrors() {
		errs = multierror.Append(errs, &DecodeError{Diags: diags})
	}
	if errs != nil {
		return errs
	}
	return nil
}

func createInverseEdge(block *hcl.Block, fromNode *DataNode, toNode *DataNode) hcl.Diagnostics {
	// Add the edge from node to edge node
	fromNode.AddInverseEdge(block.Type, toNode)
	// Add the inverse edge, from edge node to parent
	refEdge, edgeExists := RefDataEdges[fromNode.Name][block.Type]
	if !edgeExists {
		return hcl.Diagnostics{&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  fmt.Sprintf("Unknown data edge for nodes: %s --> %s", fromNode.Name, block.Type),
			Detail:   fmt.Sprintf("Unknown data edge for nodes: %s --> %s", fromNode.Name, block.Type),
			Subject:  block.Body.MissingItemRange().Ptr(),
		}}
	}
	toNode.AddEdge(refEdge, fromNode)
	return nil
}

func createEdge(block *hcl.Block, fromNode *DataNode, toNode *DataNode) hcl.Diagnostics {
	// Add the edge from node to edge node
	fromNode.AddEdge(block.Type, toNode)
	// Add the inverse edge, from edge node to parent
	refEdge, edgeExists := RefDataEdges[fromNode.Name][block.Type]
	if !edgeExists {
		return hcl.Diagnostics{&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  fmt.Sprintf("Unknown data edge for nodes: %s --> %s", fromNode.Name, block.Type),
			Detail:   fmt.Sprintf("Unknown data edge for nodes: %s --> %s", fromNode.Name, block.Type),
			Subject:  block.Body.MissingItemRange().Ptr(),
		}}
	}
	toNode.AddInverseEdge(refEdge, fromNode)
	return nil
}

// func decodeForLoop(body hcl.Body, ctx *hcl.EvalContext) (hcl.Body, hcl.Diagnostics) {
// 	var nodeBody DataNodeBody
// 	diags := gohcl.DecodeBody(body, ctx, &nodeBody)
// 	if diags.HasErrors() {
// 		return nil, diags
// 	}

// 	return nodeBody.Leftovers, nil
// }

func validateBlocks(body hcl.Body, blocks hcl.Blocks, spec *HCLBodySpec) hcl.Diagnostics {
	var diags hcl.Diagnostics
	typeBlocks := blocks.ByType()
	for _, block := range spec.Blocks {
		bls, ok := typeBlocks[block.Type]
		if block.MaxItems != 0 {
			if ok && len(bls) > block.MaxItems {
				diags = append(diags, &hcl.Diagnostic{
					Severity: hcl.DiagError,
					Summary:  fmt.Sprintf("Too many instances of block: %s", block.Type),
					Detail:   fmt.Sprintf("Too many instances of block: %s. Maximum allowed is %d", block.Type, block.MaxItems),
					Subject:  body.MissingItemRange().Ptr(),
				})
			}
		}
		if block.MinItems != 0 {
			if !ok || len(bls) < block.MinItems {
				diags = append(diags, &hcl.Diagnostic{
					Severity: hcl.DiagError,
					Summary:  fmt.Sprintf("At least %d block(s) required for: %s", block.MinItems, block.Type),
					Subject:  body.MissingItemRange().Ptr(),
				})
			}
		}
	}
	return diags
}
