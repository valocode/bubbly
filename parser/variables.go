package parser

import (
	"reflect"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/ext/dynblock"
	"github.com/hashicorp/hcl/v2/gohcl"
)

func walkVariables(node dynblock.WalkVariablesNode, ty reflect.Type) []hcl.Traversal {

	zeroVal := reflect.Zero(ty)
	schema, _ := gohcl.ImpliedBodySchema(zeroVal.Interface())

	fieldByTagName := make(map[string]reflect.Type)

	ty = nestedElem(ty)

	for i := 0; i < ty.NumField(); i++ {
		field := ty.Field(i)
		tag, exists := field.Tag.Lookup("hcl")
		// let's do some filtering
		if !exists {
			continue
		}
		// Hacky: we use the typeexpr package which HCL sees as variables, e.g.
		// the literal `type = string` will see `string` as a variable.
		// These are all decoded into hcl.Expression types, which we don't use
		// for anything else (right now), so ignore them when finding all the
		// variables
		if field.Type.Implements(reflect.TypeOf((*hcl.Expression)(nil)).Elem()) {
			continue
		}
		name := strings.Split(tag, ",")[0]
		fieldByTagName[name] = nestedElem(field.Type)
	}

	// let's remove the attributes that we do not want to get variables of
	for i, attr := range schema.Attributes {
		_, exists := fieldByTagName[attr.Name]
		if !exists {
			// remove that attribute
			schema.Attributes = append(schema.Attributes[:i], schema.Attributes[i+1:]...)
		}
	}

	vars, children := node.Visit(schema)
	for _, child := range children {
		fieldElement := fieldByTagName[child.BlockTypeName]
		vars = append(vars, walkVariables(child.Node, fieldElement)...)
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
	case reflect.Struct:
		return ty
	default:
		return ty
	}
}
