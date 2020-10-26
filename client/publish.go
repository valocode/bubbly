package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/rs/zerolog/log"
	"github.com/verifa/bubbly/api/core"
)

// Publish takes the output from a publish resource and POSTs it to the Bubbly
// server.
// Returns an error if publishing was unsuccessful
func (c *Client) Publish(data core.DataBlocks) error {

	// We must wrap the data with "data" key such that it can be unmarshalled
	// correctly by server.upload into an uploadStruct
	publishData := map[string]interface{}{
		"data": data,
	}

	log.Debug().Msgf("Making POST to Bubbly server from client.Publish to publish data %+v to %s\n", publishData, c.HostURL)

	json.Marshal(publishData)
	jsonReq, err := json.Marshal(publishData)

	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/alpha1/upload", c.HostURL), bytes.NewBuffer(jsonReq))

	if err != nil {
		return err
	}

	_, err = c.doRequest(req)

	if err != nil {
		return err
	}

	return nil
}
