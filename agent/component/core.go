package component

import (
	"fmt"
	"sync"

	"github.com/nats-io/nats.go"

	"github.com/verifa/bubbly/env"
)

const (
	defaultWaitGroupDelta = 1
)

// ComponentCore provides a minimal-viable implementation of a bubbly
// component. It:
// - provides a default implementation of the various methods
// required by a component to interact with other
// bubbly components across a NATS server,
// - stores the component's Subscriptions and Publications,
// - provides access to the underlying NATS Server connected to by the
// bubbly agent.
//
// It is recommended that all bubbly agent components embed ComponentCore.
type ComponentCore struct {
	Type           ComponentType
	NATSServerConn *nats.Conn
	Subscriptions  Subscriptions
	Publications   Publications
}

// Subscribe subscribes, via the NATS connection associated with the component,
// to a given NATS queue and subject
func (c ComponentCore) Subscribe(bCtx *env.BubblyContext, sub Subscription) error {
	bCtx.Logger.Debug().
		Interface("subscription", sub).
		Str("component", string(c.Type)).
		Msg("subscribing to subject")

	defer c.NATSServerConn.Close()

	ec, err := nats.NewEncodedConn(c.NATSServerConn, nats.JSON_ENCODER)
	if err != nil {
		return fmt.Errorf(
			"unable to establish encoded connection to the NATS server: %w",
			err,
		)
	}
	defer ec.Close()

	wg := sync.WaitGroup{}
	// For current bubbly use cases,
	// any component that subscribes to a subject will want to continue
	// listening indefinitely. Therefore, we add an initial positive delta to
	// the waitgroup, which we never decrement.
	// TODO: Integrate a toggle for this logic into the Subscription struct
	wg.Add(defaultWaitGroupDelta)

	// Create a queue subscription
	if _, err := ec.QueueSubscribe(
		string(sub.Subject),
		string(sub.Queue),
		func(m *nats.Msg) error {
			if err := c.Process(bCtx, m); err != nil {
				return fmt.Errorf(`failed to process queue subscription for subject "%s
" and queue "%s": %w`, sub.Subject, sub.Queue, err)
			}
			// TODO: Subscriptions which are non-infinite should mark wg.
			//  Done() here.
			// wg.Done()
			return nil
		},
	); err != nil {
		return fmt.Errorf(
			`failed to create queue subscription for subject "%s" and queue "%s": %w`,
			sub.Subject,
			sub.Queue,
			err,
		)
	}

	// Wait for messages to come in
	wg.Wait()
	return nil
}

// Process processes a message received as a result of a subscription
func (c ComponentCore) Process(bCtx *env.BubblyContext, m *nats.Msg) error {
	bCtx.Logger.Debug().
		Interface("nats_message", m).
		Str("component", string(c.Type)).
		Msg("processing message")

	return nil
}

// Publish publishes data on a given subject to a NATS server
func (c ComponentCore) Publish(bCtx *env.BubblyContext, pub Publication) error {
	bCtx.Logger.Debug().
		Str("component", string(c.Type)).
		Interface("publication", pub).
		Msg("publishing message")

	ec, err := nats.NewEncodedConn(c.NATSServerConn, nats.JSON_ENCODER)
	if err != nil {
		return fmt.Errorf(
			"unable to establish encoded connection to the NATS server: %w",
			err,
		)
	}

	// Publish the message
	if err := ec.Publish(
		string(pub.Subject),
		pub.Value,
	); err != nil {
		return fmt.Errorf(
			`unable to publish message (subject "%s", value "%v") over encoded channel: %w`,
			pub.Subject,
			pub.Value,
			err,
		)
	}

	return nil
}

func (c ComponentCore) Run(bCtx *env.BubblyContext) error {
	bCtx.Logger.Debug().Str(
		"component",
		string(c.Type),
	).Msg("running component")

	if err := c.BulkSubscribe(bCtx); err != nil {
		return fmt.Errorf("error during bulk subscription: %w", err)
	}

	return nil
}

// BulkSubscribe subscribes to all Subscriptions stored by the Component.
// It is typically performed at Component runtime prior to the running of any
// long-running component process.
func (c ComponentCore) BulkSubscribe(bCtx *env.BubblyContext) error {
	for _, s := range c.Subscriptions {
		// if there is a problem subscribing to any of the component's
		// provided subscriptions, error.
		if err := c.Subscribe(bCtx, s); err != nil {
			return fmt.Errorf(
				`failed to create subscription for subject "%s" and queue "%s": %w`,
				s.Subject,
				s.Queue,
				err,
			)
		}
	}
	return nil
}

// BulkPublish publishes all Publications stored by the Component.
func (c ComponentCore) BulkPublish(bCtx *env.BubblyContext) error {
	for _, p := range c.Publications {
		// if there is a problem publishing to any of the component's
		// provided publications, error.
		if err := c.Publish(bCtx, p); err != nil {
			return fmt.Errorf(
				`failed to publish subject "%s" with value "%s": %w`,
				p.Subject,
				p.Value,
				err,
			)
		}
	}

	return nil
}
