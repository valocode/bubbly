// Package resource is a data access layer that creates functionality for
// storing and retrieving Resources.
package resource

import (
	"fmt"
)

// New creates a new Dal
func New(cfg Config) (*DAL, error) {
	var (
		P   provider
		err error
	)

	switch cfg.Provider {
	case Buntdb:
		P, err = newBuntdb(cfg)
	case Etcd:
		P, err = newEtcd(cfg)
	default:
		return nil, fmt.Errorf("invalid provider: %s", cfg.Provider)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to create Provider: %w", err)
	}

	return &DAL{P: P}, nil
}

// DAL provides access to the data access layer
type DAL struct {
	P provider
}
