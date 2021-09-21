package client

import (
	"net/http"

	"github.com/valocode/bubbly/env"
)

func RunQuery(bCtx *env.BubblyConfig, query string) (map[string]interface{}, error) {
	var resp map[string]interface{}

	req := struct {
		Query *string `json:"query,omitempty"`
	}{
		Query: &query,
	}
	if err := handleRequest(
		WithBubblyContext(bCtx), WithGraphQL(true), WithMethod(http.MethodPost),
		WithPayload(req), WithRequestURL("graphql"), WithResponse(&resp),
	); err != nil {
		return nil, err
	}
	return resp, nil
}
