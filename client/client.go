package client

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	natsd "github.com/nats-io/nats-server/server"
	"github.com/nats-io/nats.go"

	"github.com/verifa/bubbly/agent/component"
	"github.com/verifa/bubbly/config"

	"github.com/verifa/bubbly/env"
)

var (
	_ HTTPClient = (*HTTP)(nil)
	_ NATSClient = (*NATS)(nil)
)

const (
	defaultHTTPClientTimeout = 5
	defaultNATSClientTimeout = 2
)

// Every Client must implement the Client interface's methods
type Client interface {
	GetResource(*env.BubblyContext, string) ([]byte, error)
	PostResource(*env.BubblyContext, []byte) error
}

type ClientType string

const (
	NATSClientType ClientType = "NATS"
	HTTPClientType ClientType = "HTTP"
)

type ClientCore struct {
	Type ClientType
}

type HTTPClient interface {
	Client
	do(r *http.Request) (io.ReadCloser, error)
}

type NATSClient interface {
	Client
	Publish(*env.BubblyContext, *component.Publication) error
	Request(*env.BubblyContext, *component.Publication) *component.Publication
}

type NATS struct {
	*ClientCore
	// The configuration of the NATS server this client should attempt to
	// connect to
	Config *config.NATSServerConfig
	Server *natsd.Server
	Conn   *nats.Conn
	EConn  *nats.EncodedConn
}

type HTTP struct {
	*ClientCore
	HostURL    string
	HTTPClient *http.Client
}

func NewHTTP(bCtx *env.BubblyContext) (*HTTP, error) {
	sc := bCtx.GetServerConfig()

	c := &HTTP{
		ClientCore: &ClientCore{
			Type: HTTPClientType,
		},
		HTTPClient: &http.Client{Timeout: defaultHTTPClientTimeout * time.Second},
		// Default bubbly server URL
		HostURL: sc.HostURL(),
	}

	if sc.Protocol != "" && sc.Host != "" && sc.Port != "" {
		us := sc.Protocol + "://" + sc.Host + ":" + sc.Port
		u, err := url.Parse(us)
		if err != nil {
			return nil, fmt.Errorf("failed to create client host: %w", err)
		}
		bCtx.Logger.Debug().Str("url", u.String()).Msg("custom bubbly host set")
		c.HostURL = u.String()
	}

	// TODO: support authenticated clients
	return c, nil
}

func (c *HTTP) do(req *http.Request) (io.ReadCloser, error) {
	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make HTTP request: %w", err)
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status: %d", res.StatusCode)
	}

	return res.Body, nil
}

// NewNATS returns a new *client.NATS bubbly client, using the NATS server configuration embedded
// within the bubbly context.
func NewNATS(bCtx *env.BubblyContext) *NATS {
	bCtx.Logger.Debug().
		Interface("client_config", bCtx.AgentConfig.NATSServerConfig).
		Msg("creating a NATS client")

	sc := bCtx.AgentConfig.NATSServerConfig

	// This configure the NATS Server using natsd package
	nopts := &natsd.Options{}

	nopts.HTTPPort = sc.HTTPPort
	nopts.Port = sc.Port

	// Create the NATS Server
	s := natsd.New(nopts)

	return &NATS{
		ClientCore: &ClientCore{
			Type: NATSClientType,
		},
		Config: sc,
		Server: s,
	}
}

// Request requests data on a given subject.
// It differs to Publish in that this requires a response from a subscriber.
// The response is decoded into a new Publication,
// which is returned to the caller.
func (n *NATS) Request(bCtx *env.BubblyContext, req *component.Publication) *component.Publication {
	// Connect to the NATS Server if a connection has not already been
	// established by this client
	if n.Conn == nil || n.EConn == nil {
		bCtx.Logger.Debug().
			Str("client", string(n.Type)).
			Msg("client is missing required connection to the NATS Server. " +
				"Attemping to connect")

		if err := n.EncodedConnect(bCtx, req.Encoder); err != nil {
			return &component.Publication{
				Subject: req.Subject,
				Error: fmt.Errorf(
					"failed to connect to the NATS Server: %w",
					err),
			}
		}
	}

	defer n.Conn.Close()
	defer n.EConn.Close()

	var reply component.Publication

	bCtx.Logger.Debug().
		Interface("nats_client", n.Config).
		Str("subject", string(req.Subject)).
		Msg("sending request")

	// Send a request.
	// The response from the request should always be a []byte,
	// which we can easily decode into our `reply.Data`.
	if err := n.EConn.Request(string(req.Subject), req.Data, &reply.Data,
		defaultNATSClientTimeout*time.Second); err != nil {
		return &component.Publication{
			Subject: req.Subject,
			Error:   fmt.Errorf("failed to make request: %w", err),
		}
	}

	return &reply
}

// publish sends a Publication (https://docs.nats.io/nats-concepts/pubsub)
// over a NATS server. It returns an error if a connection to the NATS server
// could not be established or if it was not possible to publish a message on
// the given subject.
func (n *NATS) Publish(bCtx *env.BubblyContext, pub *component.Publication) error {

	// Connect to the NATS Server if a connection has not already been
	// established by this client.
	if n.Conn == nil || n.EConn == nil {
		bCtx.Logger.Debug().
			Str("client", string(n.Type)).
			Msg("client is missing required connection to the NATS Server. " +
				"Attemping to connect")

		if err := n.EncodedConnect(bCtx, pub.Encoder); err != nil {
			return fmt.Errorf("failed to connect to the NATS Server: %w", err)
		}
	}

	defer n.Conn.Close()
	defer n.EConn.Close()

	if err := n.EConn.Publish(string(pub.Subject), pub); err != nil {
		return fmt.Errorf(
			`failed to publish subject "%s" with value "%s": %w`,
			pub.Subject,
			pub.Data,
			err,
		)
	}

	return nil
}
