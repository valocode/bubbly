package core

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/zclconf/go-cty/cty"
)

// ResourceContext provides a context for Resources to apply themselves
// independently. This is relatively complicated and probably a cleaner
// architecture exists, but for now it is serving our purpose well.
//
// The process of decoding and Applying a resource is very dependent on context,
// not just decoding HCL but also resolving resources from *.bubbly files and
// the bubbly server.
// As such, the parser.Parser needs to be quite involved but there should not
// be a coupling on the Parser from the core package, and hence this type
// provides a "connector" for a resource to use functionality from the parser,
// that is exposed in the form of method pointers defined in thi type.
type ResourceContext struct {
	GetResource GetResourceFn
	DecodeBody  DecodeBodyFn
	InsertValue InsertValueFn
	NewContext  NewContextFn
	BodyToJSON  BodyToJSONFn

	Debug DebugCtx
}

// GetResourceFn represents the function that will decode any HCL Bodies.
type GetResourceFn func(kind ResourceKind, name string) (Resource, error)

// DecodeBodyFn represents the function that will decode any HCL Bodies.
type DecodeBodyFn func(resource Resource, body hcl.Body, val interface{}) error

// NewContextFn represents the function that will provide a new context
// for resources that applies nested resources (e.g. like a pipeline applying
// a task that applies an importer)
type NewContextFn func(inputs cty.Value) *ResourceContext

// InsertValueFn takes a cty.Value and a path (represented as a slice of string)
// and adds the value at that path to the EvalContext
type InsertValueFn func(value cty.Value, path []string)

// BodyToJSONFn takes an hclsyntax.Body and returns a JSON representation of
// that body, which could be converted back into an hcl.Body.
// It also takes care of resolving any local values, so that the resource
// can be decoded later with only the necessary inputs provided.
type BodyToJSONFn func(body *hclsyntax.Body) (interface{}, error)

// DebugCtx is a debug method for providing the EvalContext
type DebugCtx func() *hcl.EvalContext
