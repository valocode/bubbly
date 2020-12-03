package bubbly

import (
	"fmt"

	"github.com/verifa/bubbly/api/core"
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
	err := ApplyPipelineRuns(bCtx, filename)

	if err != nil {
		return fmt.Errorf("failed to apply pipeline runs in file/directory `%s`: %w", filename, err)
	}
	err = ApplyTaskRuns(bCtx, filename)
	if err != nil {
		return fmt.Errorf("failed to apply task runs in file/directory `%s`: %w", filename, err)
	}
	return err
}

// ApplyPipelineRuns uses a parser to get the defined resources in the given
// location and applies any pipeline_run in those resources
func ApplyPipelineRuns(bCtx *env.BubblyContext, filename string) error {
	p, err := parser.NewParserFromFilename(filename)
	if err != nil {
		return fmt.Errorf("Failed to create parser: %s", err.Error())
	}

	if err := p.Parse(bCtx); err != nil {
		return fmt.Errorf("Failed to decode parser: %s", err.Error())
	}

	// TODO: resources should be uploaded to the server

	pipelineRunResources := p.Resources[core.PipelineRunResourceKind]
	for _, resource := range pipelineRunResources {
		bCtx.Logger.Debug().Msgf("Processing pipeline_run %s", resource.String())
		pipelineRun := resource.(core.PipelineRun)
		out := pipelineRun.Apply(bCtx, p.Context(cty.NilVal))
		if out.Error != nil {
			return fmt.Errorf(`Failed to apply pipeline_run "%s": %w`, pipelineRun.String(), out.Error)
		}
	}

	return nil
}

// ApplyTaskRuns uses a parser to get the defined resources in the given location
// and applies any task in those resources
func ApplyTaskRuns(bCtx *env.BubblyContext, filename string) error {
	p, err := parser.NewParserFromFilename(filename)
	if err != nil {
		return fmt.Errorf("Failed to create parser: %s", err.Error())
	}

	if err := p.Parse(bCtx); err != nil {
		return fmt.Errorf("Failed to decode parser: %s", err.Error())
	}

	// TODO: resources should be uploaded to the server

	taskRunResources := p.Resources[core.TaskRunResourceKind]
	for _, resource := range taskRunResources {
		bCtx.Logger.Debug().Msgf("Processing taskRun %s", resource.String())
		taskRun := resource.(core.TaskRun)
		out := taskRun.Apply(bCtx, p.Context(cty.NilVal))
		if out.Error != nil {
			return fmt.Errorf(`Failed to apply TaskRun "%s": %w`, taskRun.String(), out.Error)
		}
	}

	return nil
}
