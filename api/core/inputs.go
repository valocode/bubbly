package core

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/zclconf/go-cty/cty"
)

// InputDeclarationHCLWrapper is a wrapper around InputDeclarations so that
// they can be decoded from a HCL body and verified against the given inputs.
// Thus, this is used for validation that all the necessary inputs for a
// resources have been given.
//
// The Leftovers field is not actually needed, but we "trick" the gohcl package
// into thinking that there are more things, and then it does not complain about
// blocks or attributes that do not exist in this struct (which is would provide
// because of its strict parsing).
type InputDeclarationHCLWrapper struct {
	InputDeclarations InputDeclarations `hcl:"input,block"`
	Leftovers         hcl.Body          `hcl:",remain"`
}

// InputDeclarations is a wrapper for a slice of InputDeclaration
type InputDeclarations []*InputDeclaration

// InputDeclaration is the type representing any "input {...}" declaration
// blocks in HCL
type InputDeclaration struct {
	Name        string    `hcl:",label"`
	Description string    `hcl:"description,optional"`
	Default     cty.Value `hcl:"default,optional"`
	Type        cty.Type  `hcl:"type,optional"`
}

// InputDefinitions is a wrapper for a slice of InputDefinition
type InputDefinitions []*InputDefinition

// Value returns a cty.Value of the InputDefinitions that can be passed into
// an EvalContext to resolve variables in HCL
func (i *InputDefinitions) Value() cty.Value {
	inputs := map[string]cty.Value{}
	for _, input := range *i {
		inputs[input.Name] = input.Value
	}
	return cty.ObjectVal(
		map[string]cty.Value{
			"input": cty.ObjectVal(inputs),
		},
	)
}

// InputDefinition is the type representing any "input {...}" definition
// blocks in HCL
type InputDefinition struct {
	Name  string    `hcl:",label"`
	Value cty.Value `hcl:"value,attr"`
}
