package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/verifa/bubbly/env"
)

// Query takes the query string from a query resource spec and POSTs it
// to the bubbly server for querying against a bubbly store
// Returns a []byte representing the interface{} returned from the graphql-go
// request if successful
// Returns an error if querying was unsuccessful
func (c *Client) Query(bCtx *env.BubblyContext, query string) ([]byte, error) {

	// We must wrap the data with a "query" key such that it can be
	// unmarshalled correctly by server.Query into a queryReq
	queryData := map[string]string{
		"query": query,
	}

	jsonReq, err := json.Marshal(queryData)

	bCtx.Logger.Debug().RawJSON("request", jsonReq).Str("host", c.HostURL).Msg("sending query request to the bubbly server")

	if err != nil {
		return nil, fmt.Errorf("failed to marshal query data for loading: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/api/graphql", c.HostURL), bytes.NewBuffer(jsonReq))

	if err != nil {
		return nil, fmt.Errorf("failed to create %s request for query: %w", http.MethodPost, err)
	}

	rc, err := c.do(req)

	if err != nil {
		return nil, fmt.Errorf("failed to make %s request for query: %w", http.MethodPost, err)
	}

	defer rc.Close()
	body, err := ioutil.ReadAll(rc)

	return body, nil
}
