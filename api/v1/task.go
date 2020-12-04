package v1

import (
	"fmt"

	"github.com/hashicorp/hcl/v2"
	"github.com/verifa/bubbly/api/core"
	"github.com/verifa/bubbly/env"
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
	bCtx.Logger.Debug().Msgf("body: %v\ttask: %v\n", t.taskBlockSpec.Body, t)
	if err := ctx.DecodeBody(bCtx, t.taskBlockSpec.Body, t); err != nil {
		return core.ResourceOutput{
			Status: core.ResourceOutputFailure,
			Error:  fmt.Errorf(`Failed to decode Task "%s": %w`, t.Name(), err),
			Value:  cty.NilVal,
		}
	}

	// Task should append its output which is created by applying the
	// underlying resource
	resID := core.NewResourceIDFromString(t.ResourceID)
	res, err := ctx.GetResource(resID.Kind, resID.Name)
	if err != nil {
		return core.ResourceOutput{
			Status: core.ResourceOutputFailure,
			Error:  fmt.Errorf("Could not find resource %s in Task %s: %w", resID.String(), t.Name(), err),
			Value:  cty.NilVal,
		}
	}

	// create a new context for applying the associated resource
	resCtx := ctx.NewContext(t.Inputs.Value())
	output := res.Apply(bCtx, resCtx)

	if output.Error != nil {
		return core.ResourceOutput{
			Status: core.ResourceOutputFailure,
			Error:  fmt.Errorf("Failed to apply Task %s: %w", t.String(), output.Error),
			Value:  cty.NilVal,
		}
	}
	// set the output of Task into the current context
	ctx.InsertValue(bCtx, output.Output(), []string{"self", "task", t.Name()})

	return core.ResourceOutput{
		Status: core.ResourceOutputSuccess,
		Error:  nil,
		Value:  output.Output(),
	}
}

// String returns a human-friendly string ID for the task resource
func (t *Task) String() string {
	return fmt.Sprintf(
		"%s.%s.%s",
		t.APIVersion(), t.Kind(), t.Name(),
	)
}

// Kind returns the resource kind
func (t *Task) APIVersion() core.APIVersion {
	return core.APIVersion("v1")
}

// Kind returns the resource kind
func (t *Task) Kind() core.ResourceKind {
	return core.ResourceKind("task")
}

// Name returns the name of the task
func (t *Task) Name() string {
	return t.taskBlockSpec.Name
}

// JSON returns a JSON representation of this task block using the given
// ResourceContext.
// TODO
func (t *Task) JSON(ctx *core.ResourceContext) ([]byte, error) {
	return nil, nil
}
