package datastore

import (
	"fmt"

	"github.com/verifa/bubbly/agent/component"
	"github.com/verifa/bubbly/env"
	"github.com/verifa/bubbly/store"
)

var _ component.DataStore = (*DataStore)(nil)

type DataStore struct {
	*component.ComponentCore
	Store *store.Store
}

func (d *DataStore) New(bCtx *env.BubblyContext) error {
	bCtx.Logger.Debug().Msg("initializing the data store")
	store, err := store.New(bCtx)
	if err != nil {
		return fmt.Errorf("failed to initialise data store: %w", err)
	}

	d.Store = store
	bCtx.Logger.Debug().Msg("successfully initialised the data store")
	return nil
}
