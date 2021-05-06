package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/valocode/bubbly/agent/component"
	"github.com/valocode/bubbly/env"
)

// Query takes the query string from a query resource spec and POSTs it
// to the bubbly server for querying against a bubbly store
// Returns a []byte representing the interface{} returned from the graphql-go
// request if successful
// Returns an error if querying was unsuccessful
func (c *httpClient) Query(bCtx *env.BubblyContext, _ *component.MessageAuth, query string) ([]byte, error) {

	// We must wrap the data with a "query" key such that it can be
	// unmarshalled correctly by server.Query into a queryReq
	queryData := map[string]string{
		"query": query,
	}

	jsonReq, err := json.Marshal(queryData)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal query data for loading: %w", err)
	}

	resp, err := c.handleRequest(http.MethodPost, "/graphql", bytes.NewBuffer(jsonReq))
	if err != nil {
		return nil, fmt.Errorf("failed to make %s request for query: %w", http.MethodPost, err)
	}

	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

func (n *natsClient) Query(bCtx *env.BubblyContext, auth *component.MessageAuth, query string) ([]byte, error) {

	req := &component.Request{
		Subject: component.StoreQuery,
		Data: component.MessageData{
			Auth: auth,
			Data: []byte(query),
		},
	}

	if err := n.request(bCtx, req); err != nil {
		return nil, fmt.Errorf("NATS client failed to query: %w", err)
	}
	if req.Reply.Error != "" {
		return nil, fmt.Errorf("NATS client failed to query: %s", req.Reply.Error)
	}

	return req.Reply.Data, nil
}
