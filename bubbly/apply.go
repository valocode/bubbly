package bubbly

import (
	"fmt"

	"github.com/verifa/bubbly/api/core"
	"github.com/verifa/bubbly/env"
	"github.com/verifa/bubbly/parser"
	"github.com/zclconf/go-cty/cty"
)

// Apply uses a parser to get the defined resources in the given location and
// applies those resources
func Apply(bCtx *env.BubblyContext, filename string) error {
	p, err := parser.NewParserFromFilename(filename)
	if err != nil {
		return fmt.Errorf("Failed to create parser: %s", err.Error())
	}

	if err := p.Parse(bCtx); err != nil {
		return fmt.Errorf("Failed to decode parser: %s", err.Error())
	}

	// TODO: resources should be uploaded to the server

	pipelineRunKinds := p.Resources[core.PipelineRunResourceKind]
	for _, resource := range pipelineRunKinds {
		bCtx.Logger.Debug().Msgf("Processing pipeline_run %s", resource.String())
		pipelineRun := resource.(core.PipelineRun)
		out := pipelineRun.Apply(bCtx, p.Context(cty.NilVal))
		if out.Error != nil {
			return fmt.Errorf(`Failed to apply pipeline_run "%s": %s`, pipelineRun.String(), out.Error.Error())
		}
	}

	return nil
}
