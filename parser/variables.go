package parser

import (
	"fmt"
	"reflect"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/ext/dynblock"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/zclconf/go-cty/cty"
)

// walkVariables takes a body and a type, and returns the complete list of
// traversals (i.e. variables) that are referenced in the body
// func walkVariables(body hcl.Body, ty reflect.Type) []hcl.Traversal {
// 	node := dynblock.WalkVariables(body)
// 	traversals := walkVariablesWithImpliedType(node, ty)
// 	return traversals
// }

// walkExpandVariables takes a body and a type, and returns the list of
// traversals (i.e. variables) that are referenced in any for_each attribute
// that is used to expand any hcl blocks.
// Essentially, if the returned list of traversals are resolved, then there
// is no reason why the body cannot be expanded using dynblock
func walkExpandVariables(body hcl.Body, ty reflect.Type) []hcl.Traversal {
	node := dynblock.WalkExpandVariables(body)
	traversals := walkVariablesWithImpliedType(node, ty)
	return traversals
}

// walkVariablesWithImpliedType takes a node and a type, and walks the hcl
// node graph to collect all the traversals (i.e. variables).
// It creates a schema using the type provided
func walkVariablesWithImpliedType(node dynblock.WalkVariablesNode, ty reflect.Type) []hcl.Traversal {

	zeroVal := reflect.Zero(ty)
	schema, _ := gohcl.ImpliedBodySchema(zeroVal.Interface())

	ty = nestedElem(ty)
	tags := getFieldTags(ty)
	// cycle through attributes by tags in struct
	for attr, fieldIdx := range tags.Attributes {
		// get the nested elem type
		field := nestedElem(ty.Field(fieldIdx).Type)
		var ctyType cty.Type
		// if the field is an hcl.Expression or a cty.Type
		if field.Implements(reflect.TypeOf((*hcl.Expression)(nil)).Elem()) ||
			field == reflect.TypeOf(ctyType) {
			// find the schema attribute that matches this type, and remove it
			for i, attrName := range schema.Attributes {
				if attr == attrName.Name {
					// remove that attribute
					schema.Attributes = append(schema.Attributes[:i], schema.Attributes[i+1:]...)
				}
			}
		}
	}

	if tags.Remain != nil {
		return nil
	}

	// Now let's process the children and recurse this function.
	// It is noteworthy that when we recurse we need to pass the reflect.Type
	// that represents that child block, and that means getting the field
	// from the struct, which we do by matching the child block type name to
	// the corresponding hcl tag in the struct (e.g. `hcl:"myblock,block"` will
	// be matched by any child block type called "myblock")
	vars, children := node.Visit(schema)

	// Make sure there is not ",remain" hcl tag
	if tags.Remain == nil {
		if len(children) > 0 {
			for _, child := range children {
				fieldIdx, exists := tags.Blocks[child.BlockTypeName]
				if !exists {
					panic(fmt.Sprintf(`Could not find HCL block type "%s" in type "%s"`, child.BlockTypeName, ty.String()))
				}
				field := nestedElem(ty.Field(fieldIdx).Type)
				vars = append(vars, walkVariablesWithImpliedType(child.Node, field)...)
			}
		}
	}

	return vars
}

// nestedElem is a helper function to resolve the reflect.Type to it's
// underlying type. E.g. if it is a pointer of a pointer of a pointer of a slice
// of a pointer to an int, it will return the int type
func nestedElem(ty reflect.Type) reflect.Type {
	switch ty.Kind() {
	case reflect.Ptr:
		return nestedElem(ty.Elem())
	case reflect.Slice:
		return nestedElem(ty.Elem())
	default:
	}
	return ty
}
