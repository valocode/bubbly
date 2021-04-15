package datastore

import (
	"fmt"

	"github.com/valocode/bubbly/agent/component"
	"github.com/valocode/bubbly/env"
	"github.com/valocode/bubbly/store"
)

var _ component.Component = (*DataStore)(nil)

// New returns a new DataStore Component initialised with a new *store.Store,
// NATSServer configuration and default Subscriptions and Publications.
// Returns an error if unable to create the underlying store.
// Store using configuration provided from the bubbly context.
func New(bCtx *env.BubblyContext) (*DataStore, error) {
	bCtx.Logger.Debug().Msg("initializing the data store")
	store, err := store.New(bCtx)
	if err != nil {
		return nil, fmt.Errorf("failed to initialise data store: %w", err)
	}

	d := &DataStore{
		ComponentCore: &component.ComponentCore{
			Type:                 component.DataStoreComponent,
			DesiredSubscriptions: nil,
		},
		Store: store,
	}

	d.DesiredSubscriptions = d.defaultSubscriptions()

	bCtx.Logger.Debug().Msg("successfully initialised the data store")
	return d, nil
}

type DataStore struct {
	*component.ComponentCore
	Store *store.Store
}

// Close overrides the ComponentCore Close() so that it can also close the server
func (d *DataStore) Close() {
	// Close the core connection
	d.ComponentCore.Close()
	// Also close the server's connection
	d.Store.Close()
}

// a list of DesiredSubscriptions that the data store attempts to subscribe to
func (d *DataStore) defaultSubscriptions() component.DesiredSubscriptions {
	return component.DesiredSubscriptions{
		component.DesiredSubscription{
			Subject: component.StoreCreateTenant,
			Queue:   component.StoreQueue,
			Reply:   true,
			Handler: d.createTenant,
		},
		component.DesiredSubscription{
			Subject: component.StoreGetResourcesByKind,
			Queue:   component.StoreQueue,
			Reply:   true,
			Handler: d.getResourcesByKindHandler,
		},
		component.DesiredSubscription{
			Subject: component.StorePostSchema,
			Queue:   component.StoreQueue,
			Reply:   true,
			Handler: d.postSchemaHandler,
		},
		component.DesiredSubscription{
			Subject: component.StoreQuery,
			Queue:   component.StoreQueue,
			Reply:   true,
			Handler: d.queryHandler,
		},
		component.DesiredSubscription{
			Subject: component.StoreUpload,
			Queue:   component.StoreQueue,
			Reply:   true,
			Handler: d.uploadHandler,
		},
	}
}
