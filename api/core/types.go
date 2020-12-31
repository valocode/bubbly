package core

import (
	"github.com/zclconf/go-cty/cty"
)

// Tasks stores a map of Task by name
type Tasks map[string]Task

// Queries stores a map of Query by name
type Queries map[string]Query

// Conditions stores a map of Condition by name
type Conditions map[string]Condition

// ResourceOutput represents the output from applying a resource
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

// ResourceOutputStatus represents the output statuses for a resource
type ResourceOutputStatus string

// String gets a string value of a ResourceOutputStatus
func (r *ResourceOutputStatus) String() string {
	return string(*r)
}

const (
	// ResourceOutputSuccess represents success
	ResourceOutputSuccess ResourceOutputStatus = "Success"
	// ResourceOutputFailure represents failure
	ResourceOutputFailure ResourceOutputStatus = "Failure"
	// ResourceOutputMissingInputs represents missing inputs
	ResourceOutputMissingInputs ResourceOutputStatus = "MissingInputs"
)
