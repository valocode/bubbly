package component

import (
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
//
type Component interface {
	Subscribe(bCtx *env.BubblyContext, sub Subscription) error
	BulkSubscribe(bCtx *env.BubblyContext) error
	Process(bCtx *env.BubblyContext, m *nats.Msg) error
	Publish(bCtx *env.BubblyContext, pub Publication) error
	BulkPublish(bCtx *env.BubblyContext) error
	Run(bCtx *env.BubblyContext) error
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
