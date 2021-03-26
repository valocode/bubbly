package client

import (
	"github.com/valocode/bubbly/config"

	"github.com/valocode/bubbly/env"
)

var (
	_ Client = (*httpClient)(nil)
	_ Client = (*natsClient)(nil)
)

const (
	defaultHTTPClientTimeout = 5
	defaultNATSClientTimeout = 2
)

// Every Client must implement the Client interface's methods
type Client interface {
	// Resources
	GetResource(*env.BubblyContext, string) ([]byte, error)
	PostResource(*env.BubblyContext, []byte) error
	PostResourceToWorker(*env.BubblyContext, []byte) error
	// Data blocks
	Load(*env.BubblyContext, []byte) error
	// GraphQL Queries
	Query(*env.BubblyContext, string) ([]byte, error)
	// Applying a schema
	PostSchema(*env.BubblyContext, []byte) error
	Close()
}

// New creates a new bubbly Client.
// It checks whether the client will be run internally in the bubbly deployment
// (meaning it has direct access to the NATS server), or whether it is being
// used externally (e.g. from the command line) and should therefore use the
// HTTP client
func New(bCtx *env.BubblyContext) (Client, error) {
	// If the client is being used internally to bubbly we can talk directly
	// with the NATS server. Otherwise we need to use the HTTP client which
	// talks with bubbly's API server via HTTP
	if bCtx.ClientConfig.ClientType == config.NATSClientType {
		return newNATS(bCtx)
	}
	// Else we need the HTTP client
	return newHTTP(bCtx)
}
