package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/valocode/bubbly/env"
)

type request struct {
	bCtx              *env.BubblyConfig
	apiV1             bool
	graphql           bool
	method            string
	requestURL        string
	queryParams       map[string]string
	queryParamsStruct interface{}
	response          interface{}
	payload           interface{}
}

func WithBubblyContext(bCtx *env.BubblyConfig) func(r *request) {
	return func(r *request) {
		r.bCtx = bCtx
	}
}

func WithMethod(method string) func(r *request) {
	return func(r *request) {
		r.method = method
	}
}

func WithRequestURL(url string) func(r *request) {
	return func(r *request) {
		r.requestURL = url
	}
}

func WithQueryParamsStruct(params interface{}) func(r *request) {
	return func(r *request) {
		r.queryParamsStruct = params
	}
}

func WithQueryParams(params map[string]string) func(r *request) {
	return func(r *request) {
		r.queryParams = params
	}
}

func WithAPIV1(toggle bool) func(r *request) {
	return func(r *request) {
		r.apiV1 = toggle
	}
}

func WithGraphQL(toggle bool) func(r *request) {
	return func(r *request) {
		r.graphql = toggle
	}
}

func WithResponse(response interface{}) func(r *request) {
	return func(r *request) {
		r.response = response
	}
}

func WithPayload(payload interface{}) func(r *request) {
	return func(r *request) {
		r.payload = payload
	}
}

func (r *request) url() string {
	switch {
	case r.graphql:
		return r.bCtx.ClientConfig.GraphQL()
	// case r.apiV1 is the default
	default:
		return r.bCtx.ClientConfig.V1() + "/" + r.requestURL
	}
}

func handleRequest(opts ...func(r *request)) error {
	r := request{}
	for _, opt := range opts {
		opt(&r)
	}
	if r.method == "" {
		return errors.New("request method must be provided")
	}
	if r.bCtx == nil {
		return errors.New("bubbly context must be provided")
	}

	if r.method == http.MethodGet {
		if r.response == nil {
			return errors.New("cannot call http GET with no response - it's pointless")
		}
	}
	var b io.Reader
	if r.payload != nil {
		b = new(bytes.Buffer)
		if err := json.NewEncoder(b.(*bytes.Buffer)).Encode(r.payload); err != nil {
			return err
		}
	}
	httpReq, err := http.NewRequest(r.method, r.url(), b)
	if err != nil {
		return err
	}
	q := httpReq.URL.Query()
	// If a struct was given that should be decoded into query parameters
	if r.queryParamsStruct != nil {
		queryParams, err := structToStringMap(r.queryParamsStruct)
		if err != nil {
			return fmt.Errorf("decoding query struct into query params: %w", err)
		}
		for name, value := range queryParams {
			// Skip empty values
			if value != "" {
				q.Add(name, value)
			}
		}
	}
	for name, value := range r.queryParams {
		// Skip empty values
		if value != "" {
			q.Add(name, value)
		}
	}
	httpReq.URL.RawQuery = q.Encode()
	httpReq.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	r.bCtx.Logger.Debug().
		Str("query", q.Encode()).
		Str("url", r.url()).
		Str("method", http.MethodGet).
		Msg("Client making request")

	httpResp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return err
	}
	if err := handleResponseError(httpResp); err != nil {
		return err
	}
	if r.response != nil {
		if err := json.NewDecoder(httpResp.Body).Decode(r.response); err != nil {
			return fmt.Errorf("decoding HTTP response: %w", err)
		}
	}
	return nil
}

func handleResponseError(resp *http.Response) error {
	if resp.StatusCode == http.StatusOK {
		return nil
	}
	var httpErr echo.HTTPError
	if err := json.NewDecoder(resp.Body).Decode(&httpErr); err != nil {
		return fmt.Errorf("decoding HTTP error from the server: %w", err)
	}
	return fmt.Errorf("HTTP Status: %d, Message: %s", resp.StatusCode, httpErr.Message)
}

// structToStringMap converts a struct to a map[string]string by JSON marshalling
// and unmarshalling
func structToStringMap(req interface{}) (map[string]string, error) {
	b, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("encoding request to JSON: %w", err)
	}
	var params map[string]string
	if err := json.Unmarshal(b, &params); err != nil {
		return nil, fmt.Errorf("decoding request to string map: %w", err)
	}
	return params, nil
}
