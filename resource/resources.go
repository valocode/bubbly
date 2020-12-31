package resource

import (
	"fmt"

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
