package common

import (
	"encoding/json"
	"fmt"

	"github.com/valocode/bubbly/agent/component"
	"github.com/valocode/bubbly/api/core"
	"github.com/valocode/bubbly/client"
	"github.com/valocode/bubbly/env"
)

// LogResourceRun sends an event to the bubbly store containing the
// ResourceOuput from a resource that has been run. In short, log a resource run
func LogResourceRun(bCtx *env.BubblyContext, resourceOutput core.ResourceOutput, auth *component.MessageAuth) error {
	c, err := client.New(bCtx)
	if err != nil {
		return fmt.Errorf("failed to initialize the bubbly client: %w", err)
	}
	defer c.Close()

	dataBlocks, err := resourceOutput.EventData()
	if err != nil {
		return fmt.Errorf("error creating event data for resource output: %w", err)
	}
	dBytes, err := json.Marshal(dataBlocks)
	if err != nil {
		return fmt.Errorf("failed to marshal dBytes")
	}
	if err := c.Load(bCtx, auth, dBytes); err != nil {
		return fmt.Errorf("failed to load ResourceOutput: %w", err)
	}

	return nil
}
