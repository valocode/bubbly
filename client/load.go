package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/verifa/bubbly/api/core"
	"github.com/verifa/bubbly/env"
)

// Load takes the output from a load resource and POSTs it to the Bubbly
// server.
// Returns an error if loading was unsuccessful
func (c *HTTP) Load(bCtx *env.BubblyContext, data core.DataBlocks) error {

	jsonReq, err := json.Marshal(data)

	if err != nil {
		return fmt.Errorf("failed to marshal data for loading: %w", err)
	}

	// req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/alpha1/upload", c.HostURL), bytes.NewBuffer(jsonReq))
	_, err = c.handleResponse(
		http.Post(fmt.Sprintf("%s/api/v1/upload", c.HostURL), "application/json", bytes.NewBuffer(jsonReq)),
	)

	if err != nil {
		return fmt.Errorf(`failed to post resource: %w`, err)
	}

	return nil
}
