package api

import (
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

// NewResource creates a new resource from the given ResourceBlock, and adds
// it to Resources and also returns a pointer to it for convenience.
func (r *Resources) NewResource(resBlock *core.ResourceBlock) core.Resource {
	var resource core.Resource
	switch resBlock.Kind() {
	// TODO: use resBlock.APIVersion to get version of resource...
	// TODO: automate decoding of resBlock.SpecHCL.Body into Resource.Spec()
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
		log.Fatalf("Resource not supported: %s", resBlock.Kind())
	}
	// add the resource to the map
	if _, exists := (*r)[resource.Kind()][resource.Name()]; exists {
		log.Fatalf("Resource %s already exists!", resource.String())
	}
	// add the new resource to resources
	(*r)[resource.Kind()][resource.Name()] = resource

	return resource
}

// Resource returns the desired resource based on the ResourceKind and the name
// of the resource.
// It returns nil if the resource does not exist.
func (r *Resources) Resource(kind core.ResourceKind, name string) core.Resource {
	if resource, exists := (*r)[kind][name]; exists {
		return resource
	}
	return nil
}
