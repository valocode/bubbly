package bubbly

import (
	"encoding/json"
	"fmt"

	"github.com/valocode/bubbly/client"
	"github.com/valocode/bubbly/env"
	"github.com/valocode/bubbly/parser"
)

// Apply applies the resources in the file/directory filename
func Apply(bCtx *env.BubblyContext, filename string) error {

	var fileParser BubblyFileParser
	if err := parser.ParseFilename(bCtx, filename, &fileParser); err != nil {
		return fmt.Errorf("failed to run parser: %w", err)
	}
	resources, err := CreateResources(bCtx, fileParser)
	if err != nil {
		return fmt.Errorf("failed to parse resources: %w", err)
	}

	client, err := client.New(bCtx)
	if err != nil {
		return fmt.Errorf("failed to create bubbly client: %w", err)
	}
	defer client.Close()

	for _, res := range resources {
		bCtx.Logger.Debug().Msgf("Applying resource %s", res.String())
		resByte, err := json.Marshal(res)
		if err != nil {
			return fmt.Errorf("failed to convert resource %s to json: %w", res.String(), err)
		}
		err = client.PostResource(bCtx, nil, resByte)
		if err != nil {
			return fmt.Errorf("failed to post resource: %w", err)
		}
		// Print the name of the resource that has just been applied to give
		// user feedback
		fmt.Println(res.ID())
	}

	if err := runResources(bCtx, resources); err != nil {
		return fmt.Errorf(`failed to run resources in file/directory "%s": %w`, filename, err)
	}

	return nil
}
