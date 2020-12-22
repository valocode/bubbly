package v1

import (
	"fmt"

	"github.com/hashicorp/hcl/v2"
	"github.com/verifa/bubbly/api/common"
	"github.com/verifa/bubbly/api/core"
	"github.com/verifa/bubbly/env"
	"github.com/verifa/bubbly/parser"
	"github.com/zclconf/go-cty/cty"
)

var _ core.Task = (*Task)(nil)

type taskBlockSpec struct {
	Name string   `hcl:",label"`
	Body hcl.Body `hcl:",remain"`
}

type Task struct {
	*taskBlockSpec
	ResourceID string           `hcl:"resource,attr"`
	Inputs     InputDefinitions `hcl:"input,block"`
}

func NewTask(taskBlock *taskBlockSpec) *Task {
	return &Task{
		taskBlockSpec: taskBlock,
	}
}

// Apply returns the output from applying the task's underlying resource
func (t *Task) Apply(bCtx *env.BubblyContext, ctx *core.ResourceContext) core.ResourceOutput {
	p := parser.WithInputs(bCtx, ctx.Inputs)
	if err := p.Scope.DecodeExpandBody(bCtx, t.taskBlockSpec.Body, t); err != nil {
		return core.ResourceOutput{
			Status: core.ResourceOutputFailure,
			Error:  fmt.Errorf(`failed to decode Task "%s": %w`, t.Name(), err),
			Value:  cty.NilVal,
		}
	}

	_, output := common.RunResource(bCtx, ctx, t.ResourceID, t.Inputs.Value())
	return output
}

// Name returns the name of the task
func (t *Task) Name() string {
	return t.taskBlockSpec.Name
}
