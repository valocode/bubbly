package core

import (
	"github.com/valocode/bubbly/agent/component"
	"github.com/zclconf/go-cty/cty"
)

// NewResourceContext creates a new ResourceContext when an existing ResourceContext
// does not exist
func NewResourceContext(inputs cty.Value, newRes NewResourceFn, auth *component.MessageAuth) *ResourceContext {
	return &ResourceContext{
		Inputs:      inputs,
		State:       make(ResourceState),
		NewResource: newRes,
		Auth:        auth,
	}
}

// SubResourceContext should be used to create a new ResourceContext from an
// existing ResourceContext.
// Some of the values in ResourceContext want to be carried across to the sub
// ResourceContext and other values should not be. SubResourceContext takes care
// of that for you
func SubResourceContext(inputs cty.Value, ctx *ResourceContext) *ResourceContext {
	return &ResourceContext{
		Inputs:      inputs,
		DataBlocks:  ctx.DataBlocks,
		State:       make(ResourceState),
		NewResource: ctx.NewResource,
		Auth:        ctx.Auth,
	}
}

// ResourceContext provides a context for Resources to apply themselves
// independently.
type ResourceContext struct {
	// Inputs contains any cty.Values that should be input to a resource
	Inputs cty.Value
	// DataBlocks contains Data blocks that provide contextual information, such
	// as when looging a release_entry, the DataBlocks will contain data about the
	// release that is being logged
	DataBlocks  DataBlocks
	State       ResourceState
	NewResource NewResourceFn
	Auth        *component.MessageAuth
}

type ResourceState map[string]cty.Value

// TODO: should make this a bit more "feature rich"
func (r *ResourceState) Insert(key string, value cty.Value) {
	(*r)[key] = cty.ObjectVal(
		map[string]cty.Value{
			"value": value,
		},
	)
}

// ValueWithPath returns the ResourceState as an object at the given path.
// Example: if path ["a", "b"] is given and Resource state holds {"x": "y"},
// the result is an object {"a": {"b": {"x":"y"}}}, which can be resolved in HCL
// with a.b.x
func (r ResourceState) ValueWithPath(path []string) cty.Value {
	// make a copy of ResourceState
	retVal := make(map[string]cty.Value)
	for key, value := range r {
		retVal[key] = value
	}

	// iterate over path, backwards
	for i := len(path) - 1; i >= 0; i-- {
		// pack retVal within each string in path
		retVal[path[i]] = cty.ObjectVal(retVal)
	}
	return cty.ObjectVal(retVal)
}

// AppendInputObjects takes a list of inputs that should be object types,
// and combines them
func AppendInputObjects(inputs ...cty.Value) cty.Value {
	retVal := make(map[string]cty.Value)
	// for each provided input, add it to retVal
	for _, input := range inputs {
		for key := range input.Type().AttributeTypes() {
			retVal[key] = input.GetAttr(key)
		}
	}
	return cty.ObjectVal(retVal)
}

// NewResourceFn represents the function to create a new resource from a
// ResourceBlock. This functionality is handled by the api package, and needs
// to be passed from a higher-level dependency... This is the way
type NewResourceFn func(*ResourceBlock) (Resource, error)
