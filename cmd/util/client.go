package util

import (
	"fmt"

	"github.com/verifa/bubbly/client"
	"github.com/verifa/bubbly/env"
)

// ClientSetup is a convenience function for setting up a new bubbly client.
func ClientSetup(bCtx *env.BubblyContext) (*client.Client, error) {
	sc, err := bCtx.GetServerConfig()

	if err != nil {
		return nil, fmt.Errorf("unable to get server configuration from the bubbly context: %w", err)
	}

	if sc.Host != "" && sc.Port != "" {
		c, err := client.NewClient(bCtx)
		if err != nil {
			bCtx.Logger.Error().Msg("Unable to create Bubbly client")
			return nil, err
		}

		return c, nil
	}

	return nil, fmt.Errorf("Unable to create Bubbly client: missing required arguments.")
}
