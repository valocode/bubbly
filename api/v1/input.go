package v1

import (
	"github.com/zclconf/go-cty/cty"
)

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

// InputDefinition is the type representing any "input {...}" definition
// blocks in HCL
type InputDefinition struct {
	Name  string    `hcl:",label"`
	Value cty.Value `hcl:"value,attr"`
}
