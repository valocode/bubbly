package core

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/zclconf/go-cty/cty"
)

// HCLMainType represents the top-level structure of HCL.
// This basically describes the entire schema of parseable HCL.
type HCLMainType struct {
	ResourceBlocks ResourceBlocks `hcl:"resource,block"`
	Locals         Locals         `hcl:"local,block"`
}

// Locals is a wrapper for a slice of Local
type Locals []*Local

// Local is the type representing any "local {...}" blocks in HCL
type Local struct {
	Name        string    `hcl:",label"`
	Description string    `hcl:"description,optional"`
	Value       cty.Value `hcl:"value,attr"`
}

// Refernce returns a local's traversal to refernce this local in HCL, together
// with its associated cty.Value. This is used so that locals can be added
// to an EvalContext
func (l *Local) Reference() (hcl.Traversal, cty.Value) {
	return hcl.Traversal{
		hcl.TraverseRoot{Name: "local"},
		hcl.TraverseAttr{Name: l.Name},
	}, l.Value
}

// Data will reference a Table name, and assign the Field values into the
// corresponding Field values in the Table
type Data struct {
	Name   string
	Fields []Field
	Data   []Data
}

type Field struct {
	Name  string
	Value cty.Value
}

type ResourceOutput struct {
	ID     string
	Status ResourceOutputStatus
	Error  error
	Value  cty.Value
}

// Output returns a cty.Value which can be used inside an HCL EvalContext
// to resolve variables/traversals
func (r *ResourceOutput) Output() cty.Value {
	return cty.ObjectVal(
		map[string]cty.Value{
			"id":     cty.StringVal(r.ID),
			"status": cty.StringVal(r.Status.String()),
			"value":  r.Value,
		},
	)
}

type ResourceOutputStatus string

func (r *ResourceOutputStatus) String() string {
	return string(*r)
}

const (
	ResourceOutputSuccess ResourceOutputStatus = "Success"
	ResourceOutputFailure ResourceOutputStatus = "Failure"
)

// DecodeResourceFn represents the function that will decode any HCL Bodies.
type DecodeResourceFn func(resource Resource, body hcl.Body, val interface{}) error
