package parser

import (
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/likexian/gokit/assert"
	"github.com/zclconf/go-cty/cty"
)

var simpleCtyObject = cty.ObjectVal(map[string]cty.Value{
	"answer": cty.NumberIntVal(42),
})

func TestInsertBasic(t *testing.T) {
	traversal := hcl.Traversal{}
	traversal = append(traversal, hcl.TraverseRoot{
		Name: "root",
	})
	sym := NewSymbolTable()
	sym.insert(simpleCtyObject, traversal)

	assert.Equal(t, sym.EvalContext.Variables["root"], simpleCtyObject)
}

func TestInsertNested(t *testing.T) {
	traversal := hcl.Traversal{}
	traversal = append(traversal, hcl.TraverseRoot{
		Name: "root",
	})
	sym := NewSymbolTable()

	traversal = append(traversal, hcl.TraverseAttr{
		Name: "nested",
	})

	first := append(traversal, hcl.TraverseAttr{
		Name: "first",
	})
	firstVal := cty.NumberIntVal(42)
	sym.insert(firstVal, first)

	second := append(traversal, hcl.TraverseAttr{
		Name: "second",
	})
	secondVal := cty.StringVal("yoohoo!")
	sym.insert(secondVal, second)

	val, exists := sym.lookup(first)
	assert.Equal(t, exists, true)
	assert.Equal(t, val, firstVal)

	val, exists = sym.lookup(second)
	assert.Equal(t, exists, true)
	assert.Equal(t, val, secondVal)
}
