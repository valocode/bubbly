package api

import (
	"fmt"

	"github.com/valocode/bubbly/api/core"
	v1 "github.com/valocode/bubbly/api/v1"
)

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
