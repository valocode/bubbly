package client

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/verifa/bubbly/config"
)

// HostURL - Default bubbly server URL
const HostURL string = "http://localhost:8080"

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

// NewUnauthClient -
func NewUnauthClient(host *string) (*Client, error) {
	c := Client{
		HTTPClient: &http.Client{Timeout: 10 * time.Second},
		// Default bubbly server URL
		HostURL: HostURL,
	}

	if host != nil {
		c.HostURL = *host
	}

	return &c, nil
}

// NewClient -
func NewClient(sc config.ServerConfig) (*Client, error) {
	c := Client{
		HTTPClient: &http.Client{Timeout: 10 * time.Second},
		// Default bubbly server URL
		HostURL: HostURL,
	}

	if sc.Host != "" && sc.Port != "" {
		c.HostURL = sc.Protocol + "://" + sc.Host + ":" + sc.Port
		u, err := url.Parse(c.HostURL)

		if err != nil {
			log.Info().Msg("Invalid bubbly host URL. Using default `http://localhost:8080`")
			c.HostURL = HostURL
		}
		log.Debug().Msgf("Client established for host %s on port %s", u.Hostname(), u.Port())
	}

	if !sc.Auth {
		return &c, nil
	}

	// TODO: support authenticated clients
	return &c, nil

	// // form request body
	// rb, err := json.Marshal(AuthStruct{
	// 	Token: sc.Token,
	// })
	// if err != nil {
	// 	return nil, err
	// }

	// // authenticate
	// req, err := http.NewRequest("POST", fmt.Sprintf("%s/authenticate", c.HostURL), strings.NewReader(string(rb)))
	// if err != nil {
	// 	return nil, err
	// }

	// body, err := c.doRequest(req)

	// // parse response body
	// ar := AuthResponse{}
	// err = json.Unmarshal(body, &ar)
	// if err != nil {
	// 	return nil, err
	// }

	// c.Token = ar.Token

	// return &c, nil
}

func (c *Client) doRequest(req *http.Request) ([]byte, error) {
	req.Header.Set("Authorization", c.Token)

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status: %d, body: %s", res.StatusCode, body)
	}

	return body, err
}
