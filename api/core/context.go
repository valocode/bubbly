package core

import (
	"github.com/zclconf/go-cty/cty"
)

// NewResourceContext creates a new ResourceContext
func NewResourceContext(namespace string, inputs cty.Value, newRes NewResourceFn) *ResourceContext {
	return &ResourceContext{
		Inputs:      inputs,
		State:       make(ResourceState),
		Namespace:   namespace,
		NewResource: newRes,
	}
}

// ResourceContext provides a context for Resources to apply themselves
// independently.
type ResourceContext struct {
	// Inputs contains any cty.Values that should be input to a resource
	Inputs      cty.Value
	State       ResourceState
	Namespace   string
	NewResource NewResourceFn
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

func (r ResourceState) Value(path []string, inputs ...cty.Value) cty.Value {
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
