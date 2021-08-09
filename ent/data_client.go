// Code written by human, you are free to modify

package ent

import (
	"errors"
	"fmt"
	"time"

	"entgo.io/ent/schema/field"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/gocty"
)

type NodeOperation string

// StatusInProgress is the default value of the Status enum.
const DefaultNodeOperation = NodeOperationCreateUpdate

// Status values.
const (
	NodeOperationDefault      NodeOperation = ""
	NodeOperationCreateUpdate NodeOperation = "create_update"
	NodeOperationCreate       NodeOperation = "create"
	NodeOperationUpdate       NodeOperation = "update"
	NodeOperationQuery        NodeOperation = "query"
)

func ProcessDataNode(tx *Tx, node *DataNode) error {

	// Check that all inverse edge nodes have been resolved
	for _, edge := range node.Edges {
		if edge.Inverse && !edge.Node.Resolved {
			return fmt.Errorf("inverse edge has not been solved for edge: %s --> %s", edge.Name, edge.Node.Name)
		}
	}

	switch node.Operation {
	case NodeOperationCreateUpdate, NodeOperationDefault:
		// First query the data node and store the results in node
		if err := QueryNode(tx, node); err != nil {
			return err
		}
		if node.Value != nil {
			// If query returned one value, then we should update that node
			if err := UpdateNode(tx, node); err != nil {
				return err
			}
			return nil
		}
		if node.Value == nil {
			err := CreateNode(tx, node)
			if err != nil {
				return err
			}
		}
		return nil
		// case NodeOperationCreate, NodeOperationCreateUpdate, NodeOperationDefault:
	case NodeOperationCreate:
		err := CreateNode(tx, node)
		// If there is no error, great!
		if err == nil {
			return nil
		}
		// If there was an error, check if it was a constraint error
		var e *ConstraintError
		if errors.As(err, &e) {
			// If operation was to create, then fail as a unique constraint was
			// violated
			if node.Operation == NodeOperationCreate {
				return fmt.Errorf("cannot create node: unique constraint violated for table %s", node.Name)
			}
			// Otherwise, perform an update
			// First query the data node and store the results in node
			if err := QueryNode(tx, node); err != nil {
				return err
			}
			if node.Value == nil {
				return fmt.Errorf("unique constraint when create but query returns no records for table %s", node.Name)
			}
			// If query returned one value, then we should update that node
			if err := UpdateNode(tx, node); err != nil {
				return err
			}
			return nil
		}
		return err

	case NodeOperationUpdate:
		// Query the data node and store the results in node
		if err := QueryNode(tx, node); err != nil {
			return err
		}
	case NodeOperationQuery:
		// Query the data node and store the results in node
		if err := QueryNode(tx, node); err != nil {
			return err
		}
		if node.Value == nil {
			return fmt.Errorf("cannot query node: no records returned")
		}
	default:
		return fmt.Errorf("unknown operation for node %s: %s", node.Name, node.Operation)
	}

	return nil
}

// ctyToEntValue takes a cty.Value and a wanted ent type, and tries to convert
// the cty.Value into the ent type and assign it to the interface pointer given
func ctyToEntValue(value cty.Value, ty field.Type, v interface{}) error {
	// Handle all the numbers
	if ty.Numeric() {
		err := gocty.FromCtyValue(value, v)
		if err != nil {
			return err
		}
		return nil
	}
	switch ty {
	case field.TypeBool:
		err := gocty.FromCtyValue(value, v)
		if err != nil {
			return err
		}
		return nil
	case field.TypeString, field.TypeEnum:
		// Enums are of string type in entgo
		err := gocty.FromCtyValue(value, v)
		if err != nil {
			return err
		}
		return nil

	case field.TypeTime:
		switch value.Type() {
		case cty.String:
			t, err := time.Parse(time.RFC3339, value.AsString())
			if err != nil {
				return err
			}
			// Assign the value
			*v.(*time.Time) = t
			v = &t
			return nil
			// TODO
			// case parser.TimeType
		default:
			return fmt.Errorf("cannot convert cty value of type %s to time.Time", value.Type().FriendlyName())
		}

	case field.TypeJSON:
		return fmt.Errorf("unsupported type conversion to type %s for cty type %s", ty.String(), value.Type().FriendlyName())
	case field.TypeUUID:
		return fmt.Errorf("unsupported type conversion to type %s for cty type %s", ty.String(), value.Type().FriendlyName())
	case field.TypeBytes:
		return fmt.Errorf("unsupported type conversion to type %s for cty type %s", ty.String(), value.Type().FriendlyName())
	case field.TypeOther:
		return fmt.Errorf("unsupported type conversion to type %s for cty type %s", ty.String(), value.Type().FriendlyName())
	default:
		return fmt.Errorf("unsupported type conversion to type %s for cty type %s", ty.String(), value.Type().FriendlyName())
	}
}
