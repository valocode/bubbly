package client

import (
	"fmt"

	"github.com/mitchellh/mapstructure"
	"github.com/valocode/bubbly/env"
	"github.com/valocode/bubbly/store/api"
)

func GetEvents(bCtx *env.BubblyContext, req *api.EventGetRequest) (*api.EventGetResponse, error) {
	var (
		a      api.EventGetResponse
		params = make(map[string]string)
	)

	if err := mapstructure.Decode(req, &params); err != nil {
		return nil, fmt.Errorf("decoding request into params: %w", err)
	}

	if err := handleGetRequest(bCtx, &a, "events", params); err != nil {
		return nil, err
	}
	return &a, nil
}
