package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/rs/zerolog/log"
	"github.com/verifa/bubbly/api/core"
)

// Load takes the output from a load resource and POSTs it to the Bubbly
// server.
// Returns an error if loading was unsuccessful
func (c *Client) Load(data core.DataBlocks) error {

	// We must wrap the data with "data" key such that it can be unmarshalled
	// correctly by server.upload into an uploadStruct
	loadData := map[string]interface{}{
		"data": data,
	}

	log.Debug().Interface("data", loadData).Str("host", c.HostURL).Msg("Making POST to Bubbly server from client.Load to load data.")

	json.Marshal(loadData)
	jsonReq, err := json.Marshal(loadData)

	if err != nil {
		return fmt.Errorf("failed to marshal data for loading: %w", err)
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/alpha1/upload", c.HostURL), bytes.NewBuffer(jsonReq))

	if err != nil {
		return fmt.Errorf("failed to create POST request for data loading: %w", err)
	}

	_, err = c.doRequest(req)

	if err != nil {
		return fmt.Errorf("failed to make POST request for data loading: %w", err)
	}

	return nil
}
