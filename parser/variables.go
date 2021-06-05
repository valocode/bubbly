package parser

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/ext/dynblock"
	"github.com/hashicorp/hcl/v2/gohcl"
)

func walkVariables(node dynblock.WalkVariablesNode, ty reflect.Type, nodeBlock *dynblock.WalkVariablesChild) []hcl.Traversal {

	zeroVal := reflect.Zero(ty)
	schema, partial := gohcl.ImpliedBodySchema(zeroVal.Interface())
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
		vars = append(vars, walkVariables(child.Node, fieldElement, &child)...)
	}

	// We have to handle the strange case of `fields { ... }` inside the data
	// blocks because of the `"hcl:,remain"` tag, which means we get no attributes
	// and no blocks...
	// We know the ImpliedBodySchema method returns partial
	if partial {
		// This is VERY HACKY and we need a better way of dealing with this.
		// Initially the approach for data blocks didn't look so troublesome,
		// but turns out using ",remain" is a PITA!
		if nodeBlock.BlockTypeName == "fields" {
			attrs, diags := nodeBlock.Body().JustAttributes()
			if diags.HasErrors() {
				fmt.Printf("ERROR: %#v\n", diags.Errs())
			}
			for _, a := range attrs {
				trs := a.Expr.Variables()
				for _, tr := range trs {
					if tr.RootName() == "self" {
						vars = append(vars, tr)
					}
				}
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
	case reflect.Struct:
		return ty
	default:
		return ty
	}
}
