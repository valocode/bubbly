package v1

import (
	"fmt"

	"github.com/hashicorp/hcl/v2"
	"github.com/verifa/bubbly/api/core"
	"github.com/verifa/bubbly/env"
)

var _ core.Task = (*task)(nil)

type taskBlockSpec struct {
	Name string   `hcl:",label"`
	Body hcl.Body `hcl:",remain"`
}

type task struct {
	*taskBlockSpec
	ResourceID string           `hcl:"resource,attr"`
	Inputs     InputDefinitions `hcl:"input,block"`
}

func NewTask(taskBlock *taskBlockSpec) *task {
	return &task{
		taskBlockSpec: taskBlock,
	}
}

func (t *task) Apply(bCtx *env.BubblyContext, ctx *core.ResourceContext) error {
	fmt.Printf("body: %v\ttask: %v\n", t.taskBlockSpec.Body, t)
	if err := ctx.DecodeBody(bCtx, nil, t.taskBlockSpec.Body, t); err != nil {
		return fmt.Errorf(`Failed to decode task "%s": %s`, t.Name, err.Error())
	}
	// task should append its output which is created by applying the
	// underlying resource
	resID := core.NewResourceIDFromString(t.ResourceID)
	res, err := ctx.GetResource(resID.Kind, resID.Name)
	if err != nil {
		return fmt.Errorf("Could not find resource %s in task %s: %s", resID.String(), t.Name, err.Error())
	}

	// create a new context for applying the associated resource
	resCtx := ctx.NewContext(t.Inputs.Value())
	output := res.Apply(bCtx, resCtx)

	if output.Error != nil {
		return fmt.Errorf("Failed to apply task %s: %s", t.String(), output.Error.Error())
	}
	// set the output of task into the current context
	ctx.InsertValue(bCtx, output.Output(), []string{"self", "task", t.Name})

	return nil
}

func (t *task) String() string {
	return t.Name
}
