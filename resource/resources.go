package resource

import (
	"fmt"
	"strings"

	"github.com/verifa/bubbly/api/core"
	"github.com/verifa/bubbly/config"
	"github.com/verifa/bubbly/env"
)

// New creates a new Resources manager
func New(bCtx *env.BubblyContext) (*Resources, error) {
	var (
		P   provider
		err error
	)

	switch bCtx.ResourceConfig.Provider {
	case config.BuntdbResourceProvider:
		P, err = newBuntdb(*bCtx.ResourceConfig)
	case config.EtcdResourceProvider:
		P, err = newEtcd(*bCtx.ResourceConfig)
	default:
		return nil, fmt.Errorf("invalid provider: %s", bCtx.ResourceConfig.Provider)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to create Provider: %w", err)
	}

	return &Resources{Provider: P}, nil
}

// Resources acts as an in-memory store for resources, using a simple map
// using a resource "id" (made up of namespace,)
type Resources struct {
	Provider provider
}

func (r Resources) Get(bCtx *env.BubblyContext, id string) (core.Resource, error) {
	// get the resources based on the id
	//
	// if id does not contains a namespace, e.g. "namespace/kind/name" ==> "kind/name", then make it "default/kind/name"
	// TODO: handle namespaces
	if len(strings.Split(id, "/")) < 3 {
		id = fmt.Sprintf("%s/%s", core.DefaultNamespace, id)
	}

	// otherwise try and get the resource from the provider
	_, err := r.Provider.Query(bCtx, id)
	if err != nil {
		return nil, fmt.Errorf(`failed to get resource "%s": %w`, id, err)
	}

	// convert the resJSON to a resource
	// p := parser.EmptyParser()
	// return p.JSONToResource(bCtx, []byte(resJSON))
	return nil, nil
}

func (r *Resources) Put(bCtx *env.BubblyContext, resource core.Resource) {
	// TODO: save to db
	// can do some extra validation on resources here
	// r.Provider.Save(bCtx, resource, id)
}
