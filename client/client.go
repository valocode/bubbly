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
	sc := bCtx.GetServerConfig()

	c := &Client{
		HTTPClient: &http.Client{Timeout: 10 * time.Second},
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

	if !sc.Auth {
		return c, nil
	}

	// TODO: support authenticated clients
	return c, nil
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
