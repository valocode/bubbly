package client

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/valocode/bubbly/env"
)

func newHTTP(bCtx *env.BubblyContext) (*httpClient, error) {
	return &httpClient{
		client: &http.Client{Timeout: defaultHTTPClientTimeout * time.Second},
		// Default bubbly server URL
		url:  bCtx.ClientConfig.BubblyAddr,
		bCtx: bCtx,
	}, nil
}

type httpClient struct {
	url    string
	client *http.Client
	bCtx   *env.BubblyContext
}

func (h *httpClient) handleRequest(method string, path string, body io.Reader) (*http.Response, error) {
	url := h.url + path
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create new request: %w", err)
	}
	req.Header.Set(echo.HeaderContentType, "application/json")
	if h.bCtx.ClientConfig.AuthToken != "" {
		// Copy the received header into the request
		req.Header.Add(echo.HeaderAuthorization, h.bCtx.ClientConfig.AuthToken)
	}

	h.bCtx.Logger.Debug().Str("url", url).Str("method", method).Msg("Making HTTP client request")

	return h.handleResponse(h.client.Do(req))
}

func (h *httpClient) handleResponse(resp *http.Response, err error) (*http.Response,
	error) {
	if err != nil {
		return resp, err
	}

	// check the status code
	if resp.StatusCode != http.StatusOK {
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)

		if err != nil {
			return nil, fmt.Errorf(`failed to read body of respose with status "%s": %w`, resp.Status, err)
		}

		var httpError echo.HTTPError
		json.Unmarshal(body, &httpError)
		return nil, fmt.Errorf(`%s: %s`, resp.Status, httpError.Message)
	}

	return resp, nil
}
