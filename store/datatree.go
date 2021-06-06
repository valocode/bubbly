package store

import (
	"fmt"
	"sort"
	"strings"

	"github.com/valocode/bubbly/api/core"
	"github.com/valocode/bubbly/env"
	"github.com/valocode/bubbly/parser"
	"github.com/zclconf/go-cty/cty"
)

// dataTree stores a slice of "root" dataNodes.
// The idea is that the "root" data nodes are those nodes which can be saved
// independently (i.e. they have no data references).
// This way, we have a starting point when traversing the tree
type dataTree []*dataNode

// traverse takes a callback function fn and visits each node in the tree
func (t dataTree) traverse(bCtx *env.BubblyContext, fn visitFn) (core.DataBlocks, error) {
	blocks := core.DataBlocks{}
	for _, d := range t {
		if err := visitNode(bCtx, d, fn, &blocks); err != nil {
			return nil, err
		}
	}
	// Reset the tree afterwards (the visited boolean)
	for _, d := range t {
		resetDataNode(d)
	}
	return blocks, nil
}

// prepare performs two main tasks over a dataTree:
// 1. validates the fields and the joins that they exist in the schema
// 2. add any "builtin" fields/data-joins that should exist (e.g. handling lifecycle)
func (t dataTree) prepare(bCtx *env.BubblyContext, schema *SchemaGraph) error {
	_, err := t.traverse(bCtx, func(bCtx *env.BubblyContext, node *dataNode, blocks *core.DataBlocks) error {
		// The
		n, ok := schema.NodeIndex[node.Data.TableName]
		if !ok {
			return fmt.Errorf("data block refers to non-existent table %s", node.Data.TableName)
		}
		// Validate the fields fist
		for field, value := range node.Data.Fields.Values {
			var exists bool
			if value.Type().IsCapsuleType() {
				// TODO: do we need to check DataRefs or is it enough to check
				// the joins? Probably enough to check the joins as the DataRefs
				// are created via the joins... Just leaving this here in case!
				continue
			}
			// Check whether it's the ID field
			if field == tableIDField {
				continue
			}
			// Otherwise check all the fields for the table and make sure it's
			// one of those
			for _, tf := range n.Table.Fields {
				if tf.Name == field {
					exists = true
					break
				}
			}
			if !exists {
				return fmt.Errorf("field does not exist for data block %s: %s", node.Data.TableName, field)
			}
		}
		// Validate relationships/joins
		// TODO: is it necessary to validate both parents and children for EVERY
		// node? Or can we just do one of them?
		for _, parent := range node.Parents {
			// Check if it's the same, and refers to itself
			if parent.Data.TableName == node.Data.TableName {
				continue
			}
			if _, ok := n.Edges[parent.Data.TableName]; !ok {
				return fmt.Errorf("data joins does not exist in the schema: %s --> %s", n.Table.Name, parent.Data.TableName)
			}
		}
		for _, child := range node.Children {
			// Check if it's the same, and refers to itself
			if child.Data.TableName == node.Data.TableName {
				continue
			}
			if _, ok := n.Edges[child.Data.TableName]; !ok {
				return fmt.Errorf("data joins does not exist in the schema: %s --> %s", child.Data.TableName, n.Table.Name)
			}
		}
		// Check whether an edge to lifecycle exists and that the relationship
		// is what we expect - that this node "belongs to" lifecycle, meaning this
		// node should have a lifecycle foreign key field
		if edge, ok := n.Edges["lifecycle"]; ok && edge.Rel == BelongsTo {
			// TODO: rewrite this creating data blocks and calling dataBlocksToNodes
			nodeContext := map[string]*dataNode{
				node.Data.TableName: node,
			}
			dataBlocks := core.DataBlocks{
				core.Data{
					TableName: "lifecycle",
					Fields: &core.DataFields{Values: map[string]cty.Value{
						tableIDField: cty.CapsuleVal(parser.DataRefType, &parser.DataRef{
							TableName: n.Table.Name, Field: "lifecycle" + tableJoinSuffix,
						}),
						// TODO: initial value
						"status": cty.StringVal("unresolved"),
					}},
				},
				core.Data{
					TableName: n.Table.Name,
					Fields: &core.DataFields{Values: map[string]cty.Value{
						// Reference it's own _id field so that we get the correct
						// record to update
						tableIDField: cty.CapsuleVal(parser.DataRefType, &parser.DataRef{
							TableName: n.Table.Name, Field: tableIDField,
						}),
					}},
					Joins: []string{"lifecycle"},
					// Specify that we should update only
					Policy: core.UpdatePolicy,
				},
			}
			_, err := dataBlocksToNodes(dataBlocks, nil, nodeContext)
			if err != nil {
				return fmt.Errorf("error creating lifecycle sub-tree: %w", err)
			}
			// Set the nodes we just created to visited = true so that we don't
			// visit them after this, and in doing so, create another lifecycle
			// and end up in an infinite loop
			for _, nc := range nodeContext {
				nc.Visited = true
			}
		}
		return nil
	})
	return err
}

// resetDataNode takes a node, resets it, and then calls itself for the node's children
func resetDataNode(node *dataNode) {
	node.Visited = false
	for _, c := range node.Children {
		if c.Visited {
			resetDataNode(c)
		}
	}
}

// visitNode is a wrapper around the callback visit function to make sure that
// all the nodes are visited at least once, and at most once
func visitNode(bCtx *env.BubblyContext, node *dataNode, fn visitFn, blocks *core.DataBlocks) error {
	// First check that all the parents have been visited because we cannot
	// solve a node until all its parents have been solved
	var parentsVisited = true
	for _, parent := range node.Parents {
		if !parent.Visited {
			parentsVisited = false
			break
		}
	}
	// Check that parents have been visited and that the node is not being
	// visited more than once
	if parentsVisited && !node.Visited {
		// Visit the node with the callback method
		if err := fn(bCtx, node, blocks); err != nil {
			return fmt.Errorf("error traversing data tree at node %s: %w", node.Data.TableName, err)
		}
		// If no error, mark the node as visited
		node.Visited = true
		for _, child := range node.Children {
			if err := visitNode(bCtx, child, fn, blocks); err != nil {
				return err
			}
		}
	}
	return nil
}

// visitFn is a type declaration for the tree traversal visit callback function
type visitFn func(bCtx *env.BubblyContext, node *dataNode, blocks *core.DataBlocks) error

func newDataNode(d *core.Data) *dataNode {
	return &dataNode{
		Data:      d,
		Return:    make(map[string]interface{}),
		RefFields: make(map[string]struct{}),
		Parents:   make(map[string]*dataNode),
		Children:  make([]*dataNode, 0),
	}
}

// dataNode is the definition of a node in our data tree.
// Primarily, it wraps a core.Data block, and then provides additional fields
// to handle the relationship within the tree as well as storing the returned
// values from when the data block is processed by the database provider.
type dataNode struct {
	// Data is the underlying data block
	Data *core.Data
	// Return stores the values that are returned when the Data is stored in the
	// providers database. This is then used by any children to resolve the data
	// references to this Data
	Return map[string]interface{}
	// RefFields contains the fields that the children of this node have a
	// reference to. When the provider stores this node's Data, the Return
	// values are populated based on the RefFields
	RefFields map[string]struct{}
	// Visited stores the state as to whether this node has already been visited
	// when traversing the tree
	Visited bool
	// Parents is a map of dataNode with the key being the table name.
	// A node will only have one dataNode per parent table. If there are multiple
	// data references on a table name, they are grouped together.
	Parents map[string]*dataNode
	// Children stores a slice of child nodes, that have this node as their
	// parent.
	// Unlike the Parents field, a node can have multiple Children with the same
	// table name, and thus, we store them as a slice and not a map.
	Children []*dataNode
}

func (d *dataNode) Describe() string {
	var str string
	str += "data \"" + d.Data.TableName + "\" {\n"
	str += "  fields {\n"
	for name, val := range d.Data.Fields.Values {
		v, err := psqlValue(d, val)
		if err != nil {
			str += "    " + name + " = " + err.Error() + "\n"
			continue
		}
		str += "    " + name + " = " + fmt.Sprintf("%#v\n", v)
	}
	str += "  }\n"
	str += "  joins = [ " + strings.Join(d.Data.Joins, ", ") + " ]\n"
	str += "  policy = " + string(d.Data.Policy) + "\n"
	str += "}\n\n"
	return str
}

// orderedFields returns the node's Data field names in order
func (d *dataNode) orderedFields() []string {
	var fieldNames = make([]string, 0, len(d.Data.Fields.Values))

	for k := range d.Data.Fields.Values {
		// Skip the ID field because we never want to insert this ourselves
		if k == tableIDField {
			continue
		}
		fieldNames = append(fieldNames, k)
	}
	sort.Strings(fieldNames)
	return fieldNames
}

// orderedRefFields returns the node's RefFields in order
func (d *dataNode) orderedRefFields() []string {
	// At the minimum return the id field to return
	if _, ok := d.RefFields[tableIDField]; !ok {
		d.RefFields[tableIDField] = struct{}{}
	}

	var fields = make([]string, 0, len(d.RefFields))
	for field := range d.RefFields {
		fields = append(fields, field)
	}
	sort.Strings(fields)
	return fields
}

// addChild handles the creation of relationships between parent and child in
// both directions, and updates the parent's RefFields with the childs data
// references
func (d *dataNode) addChild(child *dataNode, fields map[string]struct{}) {
	// Add the child to the list of children
	d.Children = append(d.Children, child)
	// Add the referenced fields to the list of RefFields
	for field := range fields {
		d.RefFields[field] = struct{}{}
	}
	// Finally add the parent in the reverse order
	child.Parents[d.Data.TableName] = d
}

// createDataTree is the top-level function for creating a data tree.
// It takes the data blocks, as received by the store, and returns the data tree
func createDataTree(data core.DataBlocks) (dataTree, error) {
	// nodeContext will store a map of the "most recent" nodes to be visited
	// when traversing the data blocks.
	// Joins across data blocks are resolved by pulling the most recent data block
	// of that table, and creating a data reference to that field in that table.
	// So nodeContext is very important
	var nodeContext = make(map[string]*dataNode)

	dataNodes, err := dataBlocksToNodes(data, nil, nodeContext)
	if err != nil {
		return nil, fmt.Errorf("failed to create data node tree: %w", err)
	}

	return dataNodes, nil
}

// dataBlocksToNodes is recursively called to convert all data blocks into nodes
func dataBlocksToNodes(data core.DataBlocks, parent *core.Data, nodeContext map[string]*dataNode) (dataTree, error) {
	var dataNodes = make(dataTree, 0)
	for index := range data {
		// Store reference to the data block so that we can update it
		d := &data[index]
		// Ensure the data block has initialized fields. If the fields are not
		// initialized, it's because none were provided and so this might be a
		// data block that only does joins, in which case we will add those joins
		// as fields, which means this needs to be init'd
		if d.Fields == nil {
			d.Fields = &core.DataFields{}
		}
		if d.Fields.Values == nil {
			d.Fields.Values = make(map[string]cty.Value)
		}
		// dataRefs stores a map of table references containing a map of fields
		var dataRefs = make(map[string]map[string]struct{})

		// Check if the current data blocks have a parent. If they do, the
		// easiest thing to do is to just add a join, and then the join will
		// be processed as if it were explicitly provided.
		// Make sure the ignore_nesting property was not set, otherwise the
		// implicit join should be ignored
		if parent != nil && !d.IgnoreNesting {
			d.Joins = append(d.Joins, parent.TableName)
		}

		// It is important that we process the joins first because they will
		// create some data refs, that are stored as fields. Thuse we need to
		// process the joins before the fields.
		for _, join := range d.Joins {
			var (
				fieldName = join + tableJoinSuffix
				// Create a DataRef to the table ID field of the specified join
				fieldValue = cty.CapsuleVal(
					parser.DataRefType,
					&parser.DataRef{
						TableName: join,
						Field:     tableIDField,
					},
				)
			)
			// Create a join for the field
			d.Fields.Values[fieldName] = fieldValue
		}

		for _, fVal := range d.Fields.Values {
			// We are interested in DataRefs which are capsule types, so ignore
			// other field types
			if fVal.Type() != parser.DataRefType {
				continue
			}
			// Add the data ref to the map of references
			var (
				ref       = fVal.EncapsulatedValue().(*parser.DataRef)
				tableRefs = dataRefs[ref.TableName]
			)
			// Ensure the slice is initialized if this is the first field to be
			// added
			if tableRefs == nil {
				tableRefs = make(map[string]struct{})
			}
			// Add the field and then re-assign back to the map
			tableRefs[ref.Field] = struct{}{}
			dataRefs[ref.TableName] = tableRefs
		}

		// Create a node for the current data block
		node := newDataNode(d)

		// If there are no data refs, then it's easy, just add this data block
		// to the root data nodes
		if len(dataRefs) == 0 {
			dataNodes = append(dataNodes, node)
		}

		// If there are data references then we need to add this node as a child
		for tableName, fields := range dataRefs {
			// Get the node which the data ref refers to
			parentNode, ok := nodeContext[tableName]
			if !ok {
				return nil, fmt.Errorf("join refers to data block that does not exist: %s", tableName)
			}
			// Add the child node to the parent
			parentNode.addChild(node, fields)
		}
		// Add the node to the nodeContext so that it can be referenced by the
		// next dataRef to this table, if one exists.
		// It's important that we do this AFTER we handle the dataRefs for this
		// node in case this node depends a table with the same name
		nodeContext[d.TableName] = node

		childNodes, err := dataBlocksToNodes(d.Data, d, nodeContext)
		if err != nil {
			return nil, err
		}
		dataNodes = append(dataNodes, childNodes...)

		// Clear the unnecessary data
		d.Joins = nil
		d.Data = nil
	}

	return dataNodes, nil
}
