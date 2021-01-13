package datastore

import (
	"fmt"

	"github.com/nats-io/nats.go"

	"github.com/verifa/bubbly/agent/component"
	"github.com/verifa/bubbly/env"
	"github.com/verifa/bubbly/store"
)

var _ component.DataStore = (*DataStore)(nil)

type DataStore struct {
	*component.ComponentCore
	Store *store.Store
}

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
			Type: component.DataStoreComponent,
			NATSServer: &component.NATS{
				Config: bCtx.AgentConfig.NATSServerConfig,
			},
			DesiredSubscriptions: nil,
		},
		Store: store,
	}

	d.DesiredSubscriptions = d.defaultSubscriptions()

	bCtx.Logger.Debug().Msg("successfully initialised the data store")
	return d, nil
}

// a list of DesiredSubscriptions that the data store attempts to subscribe to
func (d *DataStore) defaultSubscriptions() component.DesiredSubscriptions {
	return component.DesiredSubscriptions{
		component.DesiredSubscription{
			Subject: component.StoreGetResource,
			Queue:   component.StoreQueue,
			Handler: d.GetResourceHandler,
			Encoder: nats.DEFAULT_ENCODER,
		},
		component.DesiredSubscription{
			Subject: component.StorePostResource,
			Queue:   component.StoreQueue,
			Handler: d.PostResourceHandler,
			Encoder: nats.DEFAULT_ENCODER,
		},
		component.DesiredSubscription{
			Subject: component.StoreGetResourcesByKind,
			Queue:   component.StoreQueue,
			Handler: d.GetResourcesByKindHandler,
			Encoder: nats.DEFAULT_ENCODER,
		},
		component.DesiredSubscription{
			Subject: component.StorePostSchema,
			Queue:   component.StoreQueue,
			Handler: d.PostSchemaHandler,
			Encoder: nats.DEFAULT_ENCODER,
		},
		component.DesiredSubscription{
			Subject: component.StoreQuery,
			Queue:   component.StoreQueue,
			Handler: d.QueryHandler,
			Encoder: nats.DEFAULT_ENCODER,
		},
		component.DesiredSubscription{
			Subject: component.StoreUpload,
			Queue:   component.StoreQueue,
			Handler: d.UploadHandler,
			Encoder: nats.DEFAULT_ENCODER,
		},
	}
}
