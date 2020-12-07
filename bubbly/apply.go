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
		return fmt.Errorf("Failed to create parser: %w", err)
	}

	if err := p.Parse(bCtx); err != nil {
		return fmt.Errorf("Failed to decode parser: %s", err)
	}

	// TODO: resources should be uploaded to the server

	pipelineRunResources := p.Resources[core.PipelineRunResourceKind]
	for _, resource := range pipelineRunResources {
		bCtx.Logger.Debug().Str("id", resource.String()).Msg("processing resource")
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
		return fmt.Errorf("Failed to create parser: %w", err)
	}

	if err := p.Parse(bCtx); err != nil {
		return fmt.Errorf("Failed to decode parser: %w", err)
	}

	// TODO: resources should be uploaded to the server

	taskRunResources := p.Resources[core.TaskRunResourceKind]
	for _, resource := range taskRunResources {
		bCtx.Logger.Debug().Str("id", resource.String()).Msg("processing resource")
		taskRun := resource.(core.TaskRun)
		out := taskRun.Apply(bCtx, p.Context(cty.NilVal))
		if out.Error != nil {
			return fmt.Errorf(`Failed to apply task_run "%s": %w`, taskRun.String(), out.Error)
		}
	}

	return nil
}

// ApplyQueries uses a parser to get the defined resources in the given location
// and applies any query in those resources
func ApplyQueries(bCtx *env.BubblyContext, filename string) error {
	p, err := parser.NewParserFromFilename(filename)
	if err != nil {
		return fmt.Errorf("Failed to create parser: %w", err)
	}

	if err := p.Parse(bCtx); err != nil {
		return fmt.Errorf("Failed to decode parser: %w", err)
	}

	// TODO: resources should be uploaded to the server

	queryResources := p.Resources[core.QueryResourceKind]
	for _, resource := range queryResources {
		bCtx.Logger.Debug().Str("id", resource.String()).Msg("processing resource")
		query := resource.(core.Query)
		out := query.Apply(bCtx, p.Context(cty.NilVal))
		if out.Error != nil {
			return fmt.Errorf(`Failed to apply query "%s": %w`, query.String(), out.Error)
		}
	}

	return nil
}
