package common

import (
	"fmt"

	"github.com/valocode/bubbly/api/core"
	"github.com/valocode/bubbly/client"
	"github.com/valocode/bubbly/env"
)

// LoadResourceOutput sends a ResourceOutput to the store for saving
func LoadResourceOutput(bCtx *env.BubblyContext, resourceOutput *core.ResourceOutput) error {
	// TODO: support check for internal vs external running
	c, err := client.New(bCtx)

	if err != nil {
		return fmt.Errorf("failed to initialize the bubbly client: %w", err)
	}

	dataBlocks, err := resourceOutput.DataBlocks()

	if err != nil {
		return fmt.Errorf("failed to construct datablocks from provided ResourceOutput: %w", err)
	}

	if err := c.Load(bCtx, dataBlocks); err != nil {
		return fmt.Errorf("failed to load ResourceOutput: %w", err)
	}

	return nil
}
