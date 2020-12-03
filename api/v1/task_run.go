package v1

import (
	"fmt"

	"github.com/verifa/bubbly/api/core"
	"github.com/verifa/bubbly/env"
	"github.com/zclconf/go-cty/cty"
)

var _ core.TaskRun = (*TaskRun)(nil)

type TaskRun struct {
	*core.ResourceBlock
	Spec TaskRunSpec
}

func NewTaskRun(resBlock *core.ResourceBlock) *TaskRun {
	return &TaskRun{
		ResourceBlock: resBlock,
	}
}

func (t *TaskRun) SpecValue() core.ResourceSpec {
	return &t.Spec
}

// Apply returns the output from applying the TaskRun's underlying resource
func (t *TaskRun) Apply(bCtx *env.BubblyContext, ctx *core.ResourceContext) core.ResourceOutput {
	// decode the hcl body into the taskrun's Spec
	if err := ctx.DecodeBody(bCtx, t, t.SpecHCL.Body, &t.Spec); err != nil {
		return core.ResourceOutput{
			Status: core.ResourceOutputFailure,
			Error:  fmt.Errorf(`Failed to decode "%s" body spec: %s`, t.String(), err.Error()),
			Value:  cty.NilVal,
		}
	}

	// get the resource ID of the underlying resource referenced
	// by the taskRun's spec
	resID := core.NewResourceIDFromString(t.Spec.ResourceID)

	// use this resourceID to get the underlying Resource
	res, err := ctx.GetResource(resID.Kind, resID.Name)
	if err != nil {
		return core.ResourceOutput{
			Status: core.ResourceOutputFailure,
			Error:  fmt.Errorf("Could not find resource %s referenced by TaskRun %s: %w", resID.String(), t.Name(), err),
			Value:  cty.NilVal,
		}
	}

	// a TaskRun should contain all necessary inputs for applying a
	// resource. Therefore we create a new context for applying the
	// resource from the inputs provided to the TaskRun
	taskRunCtx := ctx.NewContext(t.Spec.Inputs.Value())
	output := res.Apply(bCtx, taskRunCtx)

	if output.Error != nil {
		return core.ResourceOutput{
			Status: core.ResourceOutputFailure,
			Error:  fmt.Errorf("Failed to apply resource %s referenced by taskRun %s: %w", res.String(), t.String(), output.Error),
			Value:  cty.NilVal,
		}
	}

	bCtx.Logger.Debug().Interface("output", output.Output().GoString()).Interface("resource context", taskRunCtx.Debug()).Str("taskRun", t.String()).Str("referenced resource", res.String()).Msg("successfully processed TaskRun")

	// the actual cty.Value output is not meaningful to us, so we just output
	// a success status and cty.NilVal
	return core.ResourceOutput{
		Status: core.ResourceOutputSuccess,
		Error:  nil,
		Value:  cty.NilVal,
	}
}

type TaskRunSpec struct {
	ResourceID string           `hcl:"resource,attr"`
	Inputs     InputDefinitions `hcl:"input,block"`
}
