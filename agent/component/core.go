package component

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	// "sync"

	"github.com/nats-io/nats.go"

	"github.com/valocode/bubbly/env"
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
	// EConn contains the an encoded connection to the NATS server
	EConn *nats.EncodedConn
	// DesiredSubscriptions represents the pre-configured subscriptions that
	// the Component should _attempt_ to subscribe to.
	DesiredSubscriptions DesiredSubscriptions
	// Subscriptions holds subscriptions that the Component was successfully
	// able to subscribe to via NATS.
	Subscriptions []*nats.Subscription
}

// Connect connects to a NATS server and attaches the nats.Conn
// for use when the Component communicates (via NATS pub/request) with other
// Components
func (c *ComponentCore) Connect(bCtx *env.BubblyContext) error {
	nc, err := nats.Connect(bCtx.ClientConfig.NATSAddr,
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
			bCtx.ClientConfig.NATSAddr,
			err,
		)
	}

	bCtx.Logger.Debug().
		Str("component", string(c.Type)).
		Interface("nats_server", bCtx.ClientConfig.NATSAddr).
		Msg("successfully connected to NATS Server")

	c.EConn, err = nats.NewEncodedConn(nc, nats.DEFAULT_ENCODER)
	if err != nil {
		return fmt.Errorf("failed to create encoded NATS connection: %w", err)
	}

	return nil
}

func (c *ComponentCore) Close() {
	c.EConn.Close()
}

// Request makes a Request-Reply publication to the NATS server.
// The provided request is updated with the Reply, and the function returns an
// error indicating if something went wrong
func (c ComponentCore) Request(bCtx *env.BubblyContext, req *Request) error {
	bCtx.Logger.Debug().
		Str("component", string(c.Type)).
		Interface("subject", req.Subject).
		Msg("making a NATS request")

	// Make sure the pointer where we will put the reply is initialized
	// otherwise nats will fail when decoding
	if req.Reply == nil {
		req.Reply = &Reply{}
	}

	var reply []byte
	// Publish the data containing within the Publication
	if err := c.EConn.Request(
		string(req.Subject),
		req.Data,
		&reply,
		defaultRequestTimeout*time.Second,
	); err != nil {
		// just return the err unwrapped: this lets us assert the nats.Err type upstream
		return err
	}

	if err := json.Unmarshal(reply, req.Reply); err != nil {
		return fmt.Errorf("failed to decode reply from request: %w", err)
	}
	if req.Reply.Error != "" {
		return fmt.Errorf("received error from handler: %s", req.Reply.Error)
	}
	return nil
}

// Subscribe subscribes, via the NATS connection associated with the component,
// to a given NATS queue and subject
func (c ComponentCore) Subscribe(bCtx *env.BubblyContext, sub DesiredSubscription) (*nats.Subscription, error) {
	bCtx.Logger.Debug().
		Str("subject", string(sub.Subject)).
		Str("queue", string(sub.Queue)).
		Str("component", string(c.Type)).
		Msg("subscribing")

	// Create a queue subscription
	nSub, err := c.EConn.QueueSubscribe(
		string(sub.Subject),
		string(sub.Queue),
		func(m *nats.Msg) {
			val, err := sub.Handler(bCtx, m)
			if err != nil {
				bCtx.Logger.Debug().
					Err(err).
					Str("component", string(c.Type)).
					Str("subject", string(sub.Subject)).
					Str("queue", string(sub.Queue)).
					Msg("failed to handle subscription")
				// Check if we should reply indicating an error
				if sub.Reply {
					reply, _ := json.Marshal(Reply{Data: nil, Error: fmt.Errorf("failed to handle suscription: %w", err).Error()})
					c.EConn.Publish(m.Reply, reply)
					return
				}
			}
			// If we should reply to the subscription, create the reply and
			// send it
			if sub.Reply {
				// First convert the returned value into bytes, and then into
				// RawBytes so that it doesn't get double-encoded by JSON
				b, err := json.Marshal(val)
				if err != nil {
					bCtx.Logger.Debug().
						Err(err).
						Str("component", string(c.Type)).
						Str("subject", string(sub.Subject)).
						Str("queue", string(sub.Queue)).
						Msg("failed to marshal reply into raw bytes")
					reply, _ := json.Marshal(Reply{Data: nil, Error: fmt.Errorf("failed to marshal reply into raw bytes: %w", err).Error()})
					c.EConn.Publish(m.Reply, reply)
					return
				}
				reply, _ := json.Marshal(Reply{Data: b, Error: ""})
				if err := c.EConn.Publish(m.Reply, reply); err != nil {
					bCtx.Logger.Debug().
						Err(err).
						Str("component", string(c.Type)).
						Str("subject", string(sub.Subject)).
						Str("queue", string(sub.Queue)).
						Msg("failed to publish reply")

				}
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

// Listen takes a context and listens for its closure.
func (c ComponentCore) Listen(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}
