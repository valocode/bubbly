package client

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/verifa/bubbly/env"
)

// Client -
type Client struct {
	HostURL    string
	HTTPClient *http.Client
	Token      string
}

// AuthStruct -
type AuthStruct struct {
	Token string `json:"token"`
}

// AuthResponse -
type AuthResponse struct {
	Valid bool   `json:"valid"`
	Token string `json:"token"`
}

func New(bCtx *env.BubblyContext) (*Client, error) {
	sc, err := bCtx.GetServerConfig()

	if err != nil {
		return nil, fmt.Errorf("unable to get server configuration from the bubbly context: %w", err)
	}

	c := Client{
		HTTPClient: &http.Client{Timeout: 10 * time.Second},
		// Default bubbly server URL
		HostURL: sc.HostURL(),
	}

	if sc.Protocol != "" && sc.Host != "" && sc.Port != "" {
		us := sc.Protocol + "://" + sc.Host + ":" + sc.Port
		if u, err := url.Parse(us); err == nil {
			c.HostURL = u.String()
			bCtx.Logger.Info().Str("url", u.String()).Msg("custom bubbly host set")
		}
	}

	if !sc.Auth {
		return &c, nil
	}

	// TODO: support authenticated clients
	return &c, nil
}

func (c *Client) do(req *http.Request) (io.ReadCloser, error) {
	req.Header.Set("Authorization", c.Token)

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make HTTP request: %w", err)
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status: %d", res.StatusCode)
	}

	return res.Body, nil
}
