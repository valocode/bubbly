package common

import (
	"encoding/json"
	"fmt"

	"github.com/valocode/bubbly/api/core"
	"github.com/valocode/bubbly/client"
	"github.com/valocode/bubbly/env"
)

// LoadResourceOutput sends a ResourceOutput to the store for saving
func LoadResourceOutput(bCtx *env.BubblyContext, resourceOutput *core.ResourceOutput) error {
	c, err := client.New(bCtx)
	if err != nil {
		return fmt.Errorf("failed to initialize the bubbly client: %w", err)
	}
	defer c.Close()

	dataBlocks, err := resourceOutput.DataBlocks()
	if err != nil {
		return fmt.Errorf("failed to construct datablocks from provided ResourceOutput: %w", err)
	}
	dBytes, err := json.Marshal(dataBlocks)
	if err != nil {
		return fmt.Errorf("failed to marshal dBytes")
	}
	if err := c.Load(bCtx, dBytes); err != nil {
		return fmt.Errorf("failed to load ResourceOutput: %w", err)
	}

	return nil
}
