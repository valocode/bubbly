package v1

import (
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/verifa/bubbly/api/core"
	"github.com/zclconf/go-cty/cty"
)

var _ core.Task = (*Task)(nil)

type TaskBlockSpec struct {
	Name string   `hcl:",label"`
	Body hcl.Body `hcl:",remain"`
}

type Tasks map[string]*Task

type Task struct {
	Name       string
	ResourceID string           `hcl:"resource,attr"`
	Inputs     InputDefinitions `hcl:"input,block"`
}

func (t *Task) ResourceKind() core.ResourceKind {
	// TODO this is absolutely revolting...
	return core.ResourceKind(strings.Split(t.ResourceID, ".")[0])
}

func (t *Task) ResourceName() string {
	// TODO this is absolutely revolting...
	return strings.Split(t.ResourceID, ".")[1]
}

func (t *Task) TaskName() string {
	return t.Name
}

func (t *Task) Output() cty.Value {
	return cty.ObjectVal(map[string]cty.Value{
		"value": cty.StringVal("MOCK"),
	})
}
