package component

import (
	"github.com/nats-io/nats.go"

	"github.com/valocode/bubbly/env"
)

type DesiredSubscriptions []DesiredSubscription

// DesiredSubscription represents a simple encapsulation of all that is
// required to subscribe to a NATS channel for cross-component communication.
// Every DesiredSubscription should produce a *nats.Subscription at
// subscription-time.
type DesiredSubscription struct {
	Subject Subject
	Queue   Queue
	Encoder string
	// We pass a handler function in with any subscription.
	// TODO: consider the alternative of a single / generic handler with a
	//  switch for handling the different subject types.
	//  With the changes to the store (graphql interface) this might
	//  remove duplication that currently exists in a few handlers
	Handler SubscriptionHandlerFn
}

// SubscriptionHandlerFn represents a function that handles a particular
// Subscription subject. It is by bubbly components when receiving request
// /publish messages in order to process (and optionally reply)
// to these messages
type SubscriptionHandlerFn func(*env.BubblyContext, *nats.Msg) error

type Subjects []Subject

// Subject represents a string matching a NATS Subject.
// Bubbly components use Subjects within their Publish and
// Subscribe method signatures in order to communicate with one another
// Therefore, any cross-component communication requires a Subject
type Subject string

// Any Subjects that components use to communicate with one another should be
// defined centrally here
const (
	StoreGetResourcesByKind Subject = "store.GetResourcesByKind"
	StorePostResource       Subject = "store.PostResource"
	StorePostSchema         Subject = "store.PostSchema"
	StoreQuery              Subject = "store.Query"
	StoreUpload             Subject = "store.Upload"
	WorkerPostRunResource   Subject = "worker.PostRunResource"
)

type Queues []Queue

// Queue represents a string matching a NATS Queue group,
// which NATS uses to load balance messages across groups of subscribers
// Bubbly components use Queues when subscribing on a given subject.
type Queue string

// Any Queue that components use as a part of their Subjects should be
// defined centrally here
const (
	WorkerQueue Queue = "worker"
	StoreQueue  Queue = "store"
)
