package bubbly

import (
	"fmt"

	"github.com/valocode/bubbly/api"
	"github.com/valocode/bubbly/api/common"
	"github.com/valocode/bubbly/api/core"
	v1 "github.com/valocode/bubbly/api/v1"
	"github.com/valocode/bubbly/env"
	"github.com/zclconf/go-cty/cty"
)

func CreateResources(bCtx *env.BubblyContext, fileParser BubblyFileParser) ([]core.Resource, error) {
	var resources []core.Resource
	for _, resBlock := range fileParser.ResourceBlocks {
		resource, err := api.NewResource(resBlock)
		if err != nil {
			return nil, fmt.Errorf(`failed to create resource from resource block "%s": %w`, resBlock.String(), err)
		}
		resources = append(resources, resource)
	}
	return resources, nil
}

// runResources runs all resources of ResourceRun kind provided by the
// resource parser. On failure/success, it sends the ResourceRun kind's
// resource output to the bubbly event store
func runResources(bCtx *env.BubblyContext, allResources []core.Resource) error {
	for _, kind := range core.ResourceRunKinds() {
		bCtx.Logger.Debug().Msgf("Running resource kinds %s", kind)
		resources := resourcesByKind(allResources, kind)
		for _, resource := range resources {

			// TODO: there must be a nicer/easier place to check if the remote_run
			// property is set for a resource run
			if kind == core.RunResourceKind {
				r := resource.(*v1.Run)

				runCtx := core.NewResourceContext(cty.NilVal, api.NewResource, nil)

				if err := common.DecodeBody(bCtx, r.SpecHCL.Body, &r.Spec, runCtx); err != nil {
					return fmt.Errorf("failed to form resource from block: %w", err)
				}

				if r.Spec.Remote != nil {
					bCtx.Logger.Debug().Str("resource", r.String()).Msg("run is of type remote and therefore should only be run by a bubbly worker")
					continue
				} else {
					bCtx.Logger.Debug().Str("resource", r.String()).Msg("run is of type local")
				}
			}

			bCtx.Logger.Debug().Msgf("Running resource %s ...", resource.String())
			ctx := core.NewResourceContext(cty.NilVal, api.NewResource, nil)
			output := common.RunResource(bCtx, ctx, resource, cty.NilVal)
			if output.Error != nil {
				return output.Error
			}
		}
	}
	return nil
}

func resourcesByKind(resources []core.Resource, kind core.ResourceKind) []core.Resource {
	resByKind := []core.Resource{}
	for _, res := range resources {
		if res.Kind() == kind {
			resByKind = append(resByKind, res)
		}
	}
	return resByKind
}
