package api

import (
	"fmt"

	"github.com/valocode/bubbly/api/core"
	v1 "github.com/valocode/bubbly/api/v1"
	"github.com/valocode/bubbly/env"
)

func NewParserType() *ResourcesParserType {
	return &ResourcesParserType{
		Resources: []core.Resource{},
	}
}

// ResourcesParserType is used with the parser to get blocks of resources and
// convert those to actual resources.
type ResourcesParserType struct {
	Blocks    core.ResourceBlocks `hcl:"resource,block"`
	Resources []core.Resource
}

func (r *ResourcesParserType) CreateResources(bCtx *env.BubblyContext) error {
	for _, resBlock := range r.Blocks {
		resource, err := NewResource(resBlock)
		if err != nil {
			return fmt.Errorf(`failed to create resource from resource block "%s": %w`, resBlock.String(), err)
		}
		r.Resources = append(r.Resources, resource)
	}
	return nil
}

func (r ResourcesParserType) ByKind(kind core.ResourceKind) []core.Resource {
	resByKind := []core.Resource{}
	for _, res := range r.Resources {
		if res.Kind() == kind {
			resByKind = append(resByKind, res)
		}
	}
	return resByKind
}

// NewResource creates a new resource from the given ResourceBlock
// If successful, returns a pointer to the new resource
// If unsuccessful, returns an error
func NewResource(resBlock *core.ResourceBlock) (core.Resource, error) {
	var resource core.Resource
	switch resBlock.Kind() {
	// TODO: use resBlock.APIVersion to get version of resource...
	case core.ExtractResourceKind:
		resource = v1.NewExtract(resBlock)
	case core.TransformResourceKind:
		resource = v1.NewTransform(resBlock)
	case core.LoadResourceKind:
		resource = v1.NewLoad(resBlock)
	case core.PipelineResourceKind:
		resource = v1.NewPipeline(resBlock)
	case core.RunResourceKind:
		resource = v1.NewRun(resBlock)
	case core.QueryResourceKind:
		resource = v1.NewQuery(resBlock)
	case core.CriteriaResourceKind:
		resource = v1.NewCriteria(resBlock)
	default:
		return nil, fmt.Errorf(`resource not supported: "%s"`, resBlock.Kind())
	}

	return resource, nil
}
