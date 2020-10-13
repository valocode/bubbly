package api

import (
	"github.com/verifa/bubbly/api/core"
	v1 "github.com/verifa/bubbly/api/v1"
)

func NewResource(resBlock *core.ResourceBlock) core.Resource {
	switch core.ResourceKind(resBlock.Kind) {
	// TODO: use resBlock.apiVersion to get version of resource...
	case core.ImporterResourceKind:
		return v1.NewImporter(resBlock)
	}

	return nil
}
