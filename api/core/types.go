package core

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/zclconf/go-cty/cty"
)

type Locals []*Local

// Local is the type representing any "local {...}" blocks in HCL
type Local struct {
	Name        string    `hcl:",label"`
	Description string    `hcl:"description,optional"`
	Value       cty.Value `hcl:"value,attr"`
}

type Inputs []*Input

// Input is the type representing any "input {...}" blocks in HCL
type Input struct {
	Name        string         `hcl:",label"`
	Description string         `hcl:"description,optional"`
	Default     cty.Value      `hcl:"default,optional"`
	Type        hcl.Expression `hcl:"type,optional"`
}

type Outputs []*Output

// Output is the type representing any "output {...}" blocks in HCL
type Output struct {
	Name        string    `hcl:",label"`
	Description string    `hcl:"description,optional"`
	Value       cty.Value `hcl:"value,attr"`
}

type DecodeBodyFn func(body hcl.Body, val interface{}) error
