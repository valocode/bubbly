package client

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/verifa/bubbly/env"
)

// GetResource uses the bubbly api endpoint to get a resource
// TODO: with some of the new architecture it might be possible for client
// to return an actual resource, and not just a byte... Don't want to create
// a merge hell so making a note here
func (c *Client) GetResource(bCtx *env.BubblyContext, id string) ([]byte, error) {

	resp, err := handleResponse(
		http.Get(fmt.Sprintf("%s/api/resource/%s", c.HostURL, id)),
	)
	if err != nil {
		return nil, fmt.Errorf(`failed to get resource "%s": %w`, id, err)
	}

	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

// PostResource uses the bubbly api endpoint to get a resource
func (c *Client) PostResource(bCtx *env.BubblyContext, resource []byte) error {

	_, err := handleResponse(
		http.Post(fmt.Sprintf("%s/api/resource", c.HostURL), "application/json", bytes.NewBuffer(resource)),
	)

	if err != nil {
		return fmt.Errorf(`failed to post resource: %w`, err)
	}

	return nil
}
