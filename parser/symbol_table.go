package parser

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/rs/zerolog/log"
	"github.com/zclconf/go-cty/cty"
)

// SymbolTable is our wrapper around the hcl.EvalContext type.
// It is capable of insert and lookup of cty.Values, and can return the a
// hcl.EvalContext based on its contents
type SymbolTable struct {
	EvalContext *hcl.EvalContext
}

// NewSymbolTable creates a new SymbolTable
func NewSymbolTable() *SymbolTable {
	return &SymbolTable{
		EvalContext: &hcl.EvalContext{
			Variables: make(map[string]cty.Value),
			Functions: stdfunctions(),
		},
	}
}

// NestedSymbolTable creates a new SymbolTable which includes a new child
// EvalContext, and also assigns the given inputs to the "self" value.
func (s *SymbolTable) NestedSymbolTable(inputs cty.Value) *SymbolTable {
	childEvalContext := s.EvalContext.NewChild()
	childEvalContext.Variables = map[string]cty.Value{
		"self": inputs,
	}
	return &SymbolTable{
		EvalContext: childEvalContext,
	}
}

// SetInputs sets the input value
func (s *SymbolTable) SetInputs(inputs cty.Value) {
	s.EvalContext.Variables["input"] = inputs
}

// SetOutputs sets the output from a module (moduleName).
// A traversal needs to be created and module outputs are referenced with a
// root traverser "module" followed by the moduleName and then the value
func (s *SymbolTable) SetOutputs(moduleName string, outputs cty.Value) {
	traversal := hcl.Traversal{
		hcl.TraverseRoot{Name: "module"},
		hcl.TraverseAttr{Name: moduleName},
	}
	s.insert(outputs, traversal)
}

// insert takes a cty.Value and a hcl.Traversal and adds the given value at the
// given hcl.Traversal path
func (s *SymbolTable) insert(val cty.Value, traversal hcl.Traversal) {
	// log.Debug().Msgf("Inserting into symbol table: %#v --> %s", val, traversalString(traversal))
	if len(traversal) < 1 {
		panic("Cannot insert in symbol table with an empty traversal")
	}
	rootName := traversal.RootName()
	// get the root value in the Variables map
	rootVal := s.EvalContext.Variables[rootName]

	s.EvalContext.Variables[rootName] = s.insertCtyValue(rootVal, val, traversal[1:])
}

// inserCtyValue does the heavy lifting with the insert of a value.
// cty.Values are immutable, and as such, we have to create them in a functional
// way
func (s *SymbolTable) insertCtyValue(pathVal cty.Value, val cty.Value, traversal hcl.Traversal) cty.Value {
	// if path length is 0 we have traversed all the way down
	if len(traversal) == 0 {
		return val
	}

	if pathVal.IsNull() {
		// if the path value is null, then make an empty object by default
		pathVal = cty.EmptyObjectVal
	}

	switch typ := pathVal.Type(); {
	case typ.IsObjectType():
		stepAttrName := traverserName(traversal[0])
		// convert the current path value to a map so that we can modify it
		mapVal := pathVal.AsValueMap()
		// get the next path value from the map
		nextVal := mapVal[stepAttrName]
		// if map is not initialized, we need to do that
		if len(mapVal) == 0 {
			mapVal = make(map[string]cty.Value)
		}
		// assign the new cty value
		mapVal[stepAttrName] = s.insertCtyValue(nextVal, val, traversal[1:])
		return cty.ObjectVal(mapVal)
	default:
		log.Fatal().Msgf(`Unable to get next step in path with pathVal "%s" and remaining traversal "%s"`, pathVal.GoString(), traversalString(traversal))
		return cty.NilVal
	}
}

// loopup returns the value at the given traversal in the SymbolTable.
// If the value does not exist, it returns the NilValue and false
func (s *SymbolTable) lookup(traversal hcl.Traversal) (cty.Value, bool) {
	if len(traversal) < 1 {
		log.Warn().Msg("SymbolTable.lookup() received an empty traversal")
		return cty.NilVal, false
	}
	// get the base value
	rootVal, exists := s.EvalContext.Variables[traverserName(traversal[0])]
	if !exists {
		return cty.NilVal, false
	}
	path := pathFromTraversal(traversal[1:])

	val, error := path.Apply(rootVal)
	if error != nil {
		return cty.NilVal, false
	}

	return val, true
}

// pathFromTraversal returns a cty.Path representation of a traversal
func pathFromTraversal(traversal hcl.Traversal) cty.Path {
	path := cty.Path{}
	for _, tr := range traversal {
		path = path.GetAttr(traverserName(tr))
	}
	return path
}
