package api

import (
	"fmt"
	"log"

	"github.com/verifa/bubbly/api/core"
	v1 "github.com/verifa/bubbly/api/v1"
)

// Resources is a map of a map, with the first map for the kind, and the second
// map for the name
type Resources map[core.ResourceKind]map[string]core.Resource

// NewResources returns a new instance of the Resources type
func NewResources() *Resources {
	resources := Resources{}
	for _, kind := range core.ResourceKindPriority() {
		resources[kind] = make(map[string]core.Resource)
	}
	return &resources
}

// NewResourcesFromBlocks takes a list of ResourceBlock and creates a Resources
// container for them, by converting each of them into a Resource.
func NewResourcesFromBlocks(blocks core.ResourceBlocks) *Resources {
	res := NewResources()
	for _, block := range blocks {
		_, err := res.NewResource(block)

		if err != nil {
			log.Fatalf(err.Error())
		}

	}
	return res
}

// NewResource creates a new resource from the given ResourceBlock, and adds
// it to Resources
// If successful, returns a pointer to the new resource
// If unsuccessful, returns an error
func (r *Resources) NewResource(resBlock *core.ResourceBlock) (core.Resource, error) {
	var resource core.Resource
	switch resBlock.Kind() {
	// TODO: use resBlock.APIVersion to get version of resource...
	case core.ImporterResourceKind:
		resource = v1.NewImporter(resBlock)
	case core.TranslatorResourceKind:
		resource = v1.NewTranslator(resBlock)
	case core.PublishResourceKind:
		resource = v1.NewPublish(resBlock)
	case core.PipelineResourceKind:
		resource = v1.NewPipeline(resBlock)
	case core.PipelineRunResourceKind:
		resource = v1.NewPipelineRun(resBlock)
	default:
		return nil, fmt.Errorf("resource not supported: %s", resBlock.Kind())
	}
	// add the resource to the map
	if _, exists := (*r)[resource.Kind()][resource.Name()]; exists {
		return nil, fmt.Errorf("resource %s already exists", resource.String())
	}
	// add the new resource to resources
	(*r)[resource.Kind()][resource.Name()] = resource

	return resource, nil
}

// Get returns the desired resource based on the ResourceKind and the name
// of the resource.
// It returns nil if the resource does not exist.
func (r *Resources) Get(kind core.ResourceKind, name string) core.Resource {
	if resource, exists := (*r)[kind][name]; exists {
		return resource
	}
	return nil
}
