package client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/valocode/bubbly/env"
)

func newHTTP(bCtx *env.BubblyContext) (*httpClient, error) {
	sc := bCtx.GetServerConfig()

	c := &httpClient{
		// ClientCore: &ClientCore{
		// 	Type: HTTPClientType,
		// },
		Client: &http.Client{Timeout: defaultHTTPClientTimeout * time.Second},
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

type httpClient struct {
	// *ClientCore
	HostURL string
	Client  *http.Client
}

func (h *httpClient) newRequest() (*http.Request, error) {

	return nil, nil
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
