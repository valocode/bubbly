package component

import (
	"context"

	"github.com/nats-io/nats.go"

	"github.com/verifa/bubbly/env"
)

// Component provides a NATS-compatible interface for bubbly agent services.
// It acts as a thin wrapper around NATS which all bubbly services (
// API Server, data store etc) must implement in order for inter-service
// communication to be possible.
//
// Generally, bubbly services implement this Component interface by embedding
// ComponentsCore,
// which offers default method implementations for the Component interface.
type Component interface {
	// Connect to a NATS Server
	Connect(bCtx *env.BubblyContext) error
	// Subscribe to publications on a given subject
	Subscribe(bCtx *env.BubblyContext, sub DesiredSubscription) (*nats.Subscription, error)
	// BulkSubscribe to the component's DesiredSubscriptions
	BulkSubscribe(bCtx *env.BubblyContext) ([]*nats.Subscription, error)
	// Publish a publication
	Publish(bCtx *env.BubblyContext, pub Publication) error
	// Run is the main entrypoint into running the underlying bubbly
	// processes wrapped by the Component
	Run(bCtx *env.BubblyContext, agentContext context.Context) error
	// Listen on the errgroup for anything which should trigger the Component
	// to terminate
	Listen(agentContext context.Context) error
}

type DataStore interface {
	Component
}

type UI interface {
	Component
}

type Worker interface {
	Component
}

type NATSServer interface {
	Component
}

type APIServer interface {
	Component
}

// ComponentType provides string representations of bubbly Components
type ComponentType string

const (
	DataStoreComponent  ComponentType = "data_store"
	NATSServerComponent ComponentType = "nats_server"
	APIServerComponent  ComponentType = "api_server"
	WorkerComponent     ComponentType = "worker"
	UIComponent         ComponentType = "ui"
)
