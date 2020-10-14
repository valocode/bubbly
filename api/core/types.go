package core

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/zclconf/go-cty/cty"
)

// HCLMainType represents the top-level structure of HCL.
// This basically describes the entire schema of parseable HCL.
type HCLMainType struct {
	ResourceBlocks ResourceBlocks `hcl:"resource,block"`
	ModuleBlocks   ModuleBlocks   `hcl:"module,block"`
	Locals         Locals         `hcl:"local,block"`
	Inputs         Inputs         `hcl:"input,block"`
	Outputs        Outputs        `hcl:"output,block"`
}

// Locals is a wrapper for a slice of Local
type Locals []*Local

// Local is the type representing any "local {...}" blocks in HCL
type Local struct {
	Name        string    `hcl:",label"`
	Description string    `hcl:"description,optional"`
	Value       cty.Value `hcl:"value,attr"`
}

// Inputs is a wrapper for a slice of Input
type Inputs []*Input

// Input is the type representing any "input {...}" blocks in HCL
type Input struct {
	Name        string         `hcl:",label"`
	Description string         `hcl:"description,optional"`
	Default     cty.Value      `hcl:"default,optional"`
	Type        hcl.Expression `hcl:"type,optional"`
}

// Outputs is a wrapper for a slice of Output
type Outputs []*Output

// Output is the type representing any "output {...}" blocks in HCL
type Output struct {
	Name        string    `hcl:",label"`
	Description string    `hcl:"description,optional"`
	Value       cty.Value `hcl:"value,attr"`
}

// DecodeBodyFn represents the function that will decode any HCL Bodies.
type DecodeBodyFn func(body hcl.Body, val interface{}) error
