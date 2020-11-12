package util

import (
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/verifa/bubbly/client"
	"github.com/verifa/bubbly/config"
)

// ClientSetup is a convenience function for setting up a new bubbly client.
func ClientSetup(sc config.ServerConfig) (*client.Client, error) {
	if sc.Host != "" && sc.Port != "" {
		c, err := client.NewClient(sc)
		if err != nil {
			log.Error().Msg("Unable to create Bubbly client")
			return nil, err
		}

		return c, nil
	}

	return nil, fmt.Errorf("Unable to create Bubbly client: missing required arguments.")
}
