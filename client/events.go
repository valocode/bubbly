package client

import (
	"net/http"

	"github.com/valocode/bubbly/env"
	"github.com/valocode/bubbly/store/api"
)

func GetEvents(bCtx *env.BubblyConfig, req *api.EventGetRequest) (*api.EventGetResponse, error) {
	var r api.EventGetResponse
	if err := handleRequest(
		WithBubblyContext(bCtx), WithAPIV1(true), WithMethod(http.MethodGet),
		WithResponse(&r), WithRequestURL("events"), WithQueryParamsStruct(req),
	); err != nil {
		return nil, err
	}
	return &r, nil
}
