package v1

import (
	"fmt"

	"github.com/verifa/bubbly/api/common"
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
	if err := common.DecodeBody(bCtx, t.SpecHCL.Body, &t.Spec, ctx); err != nil {
		return core.ResourceOutput{
			Status: core.ResourceOutputFailure,
			Error:  fmt.Errorf(`failed to decode "%s" body spec: %s`, t.String(), err.Error()),
			Value:  cty.NilVal,
		}
	}

	_, output := common.RunResource(bCtx, ctx, t.Spec.ResourceID, t.Spec.Inputs.Value())
	return output
}

type TaskRunSpec struct {
	ResourceID string                `hcl:"resource,attr"`
	Inputs     core.InputDefinitions `hcl:"input,block"`
}
