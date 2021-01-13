package bubbly

import (
	"encoding/json"
	"fmt"

	"github.com/verifa/bubbly/api"
	"github.com/verifa/bubbly/api/core"
	"github.com/verifa/bubbly/client"
	"github.com/verifa/bubbly/env"
	"github.com/verifa/bubbly/parser"
	"github.com/zclconf/go-cty/cty"
)

/*
TODO: We should come up with a more clear naming convention for Apply.
Apply itself should not care about the underlying resource type; the application
of any resource should be valid.
*/

// Apply calls resource-specific apply functions
func Apply(bCtx *env.BubblyContext, filename string) error {

	resParser := api.NewParserType()
	// err = p.Parse(bCtx, resParser)
	if err := parser.ParseFilename(bCtx, filename, resParser); err != nil {
		return fmt.Errorf("failed to run parser: %w", err)
	}
	if err := resParser.CreateResources(bCtx); err != nil {
		return fmt.Errorf("failed to parse resources: %w", err)
	}

	client, err := client.NewHTTP(bCtx)
	if err != nil {
		return fmt.Errorf("failed to create bubbly client: %w", err)
	}

	for _, res := range resParser.Resources {
		bCtx.Logger.Debug().Msgf(`Applying resource "%s"`, res.String())
		resByte, err := json.Marshal(res)
		if err != nil {
			return fmt.Errorf("failed to convert resource %s to json: %w", res.String(), err)
		}
		err = client.PostResource(bCtx, resByte)
		if err != nil {
			return fmt.Errorf("failed to post resource: %w", err)
		}
	}

	if err := runResources(bCtx, resParser); err != nil {
		return fmt.Errorf(`failed to run resources in file/directory "%s": %w`, filename, err)
	}

	return nil
}

func runResources(bCtx *env.BubblyContext, resParser *api.ResourcesParserType) error {
	for _, kind := range core.ResourceRunKinds() {
		bCtx.Logger.Debug().Msgf("Running resource kinds %s", kind)
		resources := resParser.ByKind(kind)
		for _, resource := range resources {
			bCtx.Logger.Debug().Msgf("Running resource %s ...", resource.String())
			out := resource.Apply(
				bCtx,
				core.NewResourceContext(resource.Namespace(), cty.NilVal, api.NewResource),
			)
			if out.Error != nil {
				return fmt.Errorf(
					`Failed to apply resource "%s": %w`,
					resource.String(), out.Error,
				)
			}
		}
	}
	return nil
}
