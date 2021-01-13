package component

import (
	"context"
	"fmt"
	"sync"
	"time"

	// "sync"

	"github.com/nats-io/nats.go"

	"github.com/verifa/bubbly/config"
	"github.com/verifa/bubbly/env"
)

const (
	// wait a maximum of 2 seconds for responses to NATS Requests.
	// A timeout is indication that there is no subscriber listening on the
	// given subject.
	defaultRequestTimeout = 2
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
	Type ComponentType
	// NATSServer contains the connections to the NATS Server.
	// This is nested within the ComponentCore as Components are
	// individually responsible for managing their own connection
	// to a NATS Server.
	NATSServer *NATS
	// DesiredSubscriptions represents the pre-configured subscriptions that
	// the Component should _attempt_ to subscribe to.
	DesiredSubscriptions DesiredSubscriptions
	// Subscriptions holds subscriptions that the Component was successfully
	// able to subscribe to via NATS.
	Subscriptions []*nats.Subscription
	// Publications represents the list of publications that a Component
	// should publish at runtime.
}

type NATS struct {
	mu     sync.Mutex
	Config *config.NATSServerConfig
	// Connection to the NATS Server.
	// We do not store the EncodedConn because that is Publication-scoped and
	// will therefore vary depending on the type of data being published.
	Conn *nats.Conn
}

// Connect connects to a NATS server and attaches the nats.Conn
// for use when the Component communicates (via NATS pub/request) with other
// Components
func (c *ComponentCore) Connect(bCtx *env.BubblyContext) error {
	nc, err := nats.Connect(c.NATSServer.Config.Addr,
		nats.Name(fmt.Sprintf("Bubbly Agent Component: %s", string(c.Type))),
		nats.ErrorHandler(func(nc *nats.Conn, s *nats.Subscription, err error) {
			if s != nil {
				bCtx.Logger.Error().
					Err(err).
					Str("subject", s.Subject).
					Str("queue", s.Queue).
					Msg("Async error")
			} else {
				bCtx.Logger.Error().
					Err(err).
					Msg("Async error outside subscription")
			}
		}),
		nats.ClosedHandler(func(nc *nats.Conn) {
			bCtx.Logger.Debug().
				Err(nc.LastError()).
				Msg("Closed connection to NATS Server")
		}),
		nats.DisconnectErrHandler(func(nc *nats.Conn, err error) {
			bCtx.Logger.Debug().
				Err(err).
				Msg("Disconnected from NATS server")
		}),
		nats.ReconnectHandler(func(nc *nats.Conn) {
			bCtx.Logger.Debug().
				Str("addr", nc.ConnectedAddr()).
				Msg("Reconnected to NATS server")
		}))
	if err != nil {
		return fmt.Errorf(
			`failed to establish a connection to the NATS
			server at address "%s": %w`,
			c.NATSServer.Config.Addr,
			err,
		)
	}

	bCtx.Logger.Debug().
		Str("component", string(c.Type)).
		Interface("nats_server", c.NATSServer.Config).
		Msg("successfully connected to NATS Server")

	c.NATSServer.Conn = nc

	return nil
}

// Publish publishes the data field of the given Publication to a NATS server.
// It establishes an encoded connection using the provided pub.Encoder string,
// allowing the encoding to be scoped to Publications
func (c ComponentCore) Publish(bCtx *env.BubblyContext, pub Publication) error {
	bCtx.Logger.Debug().
		Str("component", string(c.Type)).
		Interface("subject", pub.Subject).
		Msg("publishing")

	if err := c.checkConnection(bCtx); err != nil {
		return fmt.Errorf("failed during connection check: %w", err)
	}

	// the component's existing encoded connection could use the wrong
	// encoder type. Therefore,
	// we always create a new Publication-scoped encoded connection
	ec, err := nats.NewEncodedConn(c.NATSServer.Conn, string(pub.Encoder))
	if err != nil {
		return fmt.Errorf(
			"unable to establish encoded connection to the NATS server: %w",
			err,
		)
	}

	// Publish the data containing within the Publication
	if err := ec.Publish(
		string(pub.Subject),
		pub.Data,
	); err != nil {
		return fmt.Errorf(
			`unable to publish message (subject "%s", value "%v") over encoded channel: %w`,
			pub.Subject,
			pub.Data,
			err,
		)
	}

	return nil
}

// Request requests the data field of the given Publication to a NATS server.
// It returns a non-nil Publication if unsuccessful and a Publication with
// non-nil Data field if successful.
func (c ComponentCore) Request(bCtx *env.BubblyContext, pub Publication) (*Publication, error) {
	bCtx.Logger.Debug().
		Str("component", string(c.Type)).
		Interface("subject", pub.Subject).
		Msg("publishing")

	if err := c.checkConnection(bCtx); err != nil {
		return nil, fmt.Errorf("failed during connection check: %w", err)
	}

	// the component's existing encoded connection could use the wrong
	// encoder type. Therefore,
	// we always create a new Publication-scoped encoded connection
	ec, err := nats.NewEncodedConn(c.NATSServer.Conn, string(pub.Encoder))
	if err != nil {
		return nil, fmt.Errorf(
			"unable to establish encoded connection to the NATS server: %w",
			err,
		)
	}

	var reply Publication

	// Publish the data containing within the Publication
	if err := ec.Request(
		string(pub.Subject),
		pub.Data,
		&reply.Data,
		defaultRequestTimeout*time.Second,
	); err != nil {
		// just return the err unwrapped: this lets us assert the nats.Err type upstream
		return nil, err
	}

	return &reply, nil
}

// Subscribe subscribes, via the NATS connection associated with the component,
// to a given NATS queue and subject
func (c ComponentCore) Subscribe(bCtx *env.BubblyContext, sub DesiredSubscription) (*nats.Subscription, error) {
	bCtx.Logger.Debug().
		Str("subject", string(sub.Subject)).
		Str("queue", string(sub.Queue)).
		Str("component", string(c.Type)).
		Interface("connection_addr", c.NATSServer.Conn.ConnectedAddr()).
		Msg("subscribing")

	// make sure that the component's NATS server connection is available. If not,
	// re-establish it or error
	if err := c.checkConnection(bCtx); err != nil {
		return nil, fmt.Errorf("failed during connection check: %w", err)
	}

	// encoded connections are Subscription-scoped, since they depend on a specific
	// encoder. So we create a new one.
	ec, err := nats.NewEncodedConn(c.NATSServer.Conn, sub.Encoder)
	if err != nil {
		return nil, fmt.Errorf(
			"unable to establish encoded connection to the NATS server: %w",
			err,
		)
	}

	// Create a queue subscription
	nSub, err := ec.QueueSubscribe(
		string(sub.Subject),
		string(sub.Queue),
		func(m *nats.Msg) {
			if err := sub.Handler(bCtx, m); err != nil {
				bCtx.Logger.Debug().Err(err).
					Str("component", string(c.Type)).
					Str("subject", string(sub.Subject)).
					Str("queue", string(sub.Queue)).
					Msg("failed to handle publish/request to subscription")
			}
		},
	)

	if err != nil {
		return nil, fmt.Errorf(
			`failed to create queue subscription for subject "%s" and queue "%s": %w`,
			sub.Subject,
			sub.Queue,
			err,
		)
	}

	return nSub, nil
}

func (c ComponentCore) Run(bCtx *env.BubblyContext, agentContext context.Context) error {
	bCtx.Logger.Debug().Str(
		"component",
		string(c.Type),
	).Msg("running component")

	nSubs, err := c.BulkSubscribe(bCtx)

	if err != nil {
		return fmt.Errorf("error during bulk subscription: %w", err)
	}

	c.Subscriptions = nSubs
	bCtx.Logger.Debug().Str("component", string(c.Type)).Interface("subscriptions", c.Subscriptions).Msg("component is listening for subscriptions")

	return c.Listen(agentContext)
}

// BulkSubscribe subscribes to all Subscriptions stored by the Component.
// It is typically performed at Component runtime prior to the running of any
// long-running component process.
func (c ComponentCore) BulkSubscribe(bCtx *env.BubblyContext) ([]*nats.Subscription, error) {
	bCtx.Logger.Debug().Int("sub_count", len(c.DesiredSubscriptions)).Msg("bulk subscribing")

	wg := sync.WaitGroup{}

	wg.Add(len(c.DesiredSubscriptions))

	var nSubs []*nats.Subscription

	for _, s := range c.DesiredSubscriptions {
		// if there is a problem subscribing to any of the component's
		// provided subscriptions, error.
		nSub, err := c.Subscribe(bCtx, s)

		if err != nil {
			return nil, fmt.Errorf(
				`failed to create subscription for subject "%s" and queue "%s": %w`,
				s.Subject,
				s.Queue,
				err,
			)
		}

		bCtx.Logger.Debug().Str("component", string(c.Type)).Interface("subscription", nSub).Msg("successfully subscribed")
		nSubs = append(nSubs, nSub)

		wg.Done()

	}

	wg.Wait()

	bCtx.Logger.Debug().Str("component", string(c.Type)).Interface("subscriptions", nSubs).Msg("finished bulk subscribe")
	return nSubs, nil
}

// checkConnection checks that the Component has a valid connection to the
// NATS server prior to a subscription/publication. If not,
// it attempts to connect to the NATS server.
func (c ComponentCore) checkConnection(bCtx *env.BubblyContext) error {
	c.NATSServer.mu.Lock()
	defer c.NATSServer.mu.Unlock()

	switch c.NATSServer.Conn.Status() {
	case nats.CONNECTED:
		bCtx.Logger.Debug().
			Interface("connection_addr", c.NATSServer.Conn.ConnectedAddr()).
			Msg("already connected to a NATS server")
		return nil
	case nats.CLOSED:
		if err := c.Connect(bCtx); err != nil {
			return fmt.Errorf("failed to connect to the NATS Server: %w", err)
		} else {
			bCtx.Logger.Debug().
				Interface("connection_addr", c.NATSServer.Conn.ConnectedAddr()).
				Msg("successfully connected to a NATS server")
		}
	}
	return nil
}

// Listen takes a context and listens for its closure.
func (c ComponentCore) Listen(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}
