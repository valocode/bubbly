package v1

import (
	"fmt"

	"github.com/hashicorp/hcl/v2"

	"github.com/valocode/bubbly/api/common"
	"github.com/valocode/bubbly/api/core"
	"github.com/valocode/bubbly/env"
	"github.com/valocode/bubbly/events"

	"github.com/zclconf/go-cty/cty"
)

var _ core.Condition = (*Condition)(nil)

type conditionBlockSpec struct {
	Name string   `hcl:",label"`
	Body hcl.Body `hcl:",remain"`
}

type Condition struct {
	*conditionBlockSpec
	Value cty.Value `hcl:"value,attr"`
}

func NewCondition(conditionBlock *conditionBlockSpec) *Condition {
	return &Condition{
		conditionBlockSpec: conditionBlock,
	}
}

// Apply returns a core.ResourceOutput, whose Value is a cty.BoolVal indicating
// success of failure of the condition
func (c *Condition) Apply(bCtx *env.BubblyContext, ctx *core.ResourceContext) core.ResourceOutput {
	if err := common.DecodeBody(bCtx, c.conditionBlockSpec.Body, c, ctx); err != nil {
		return core.ResourceOutput{
			Status: events.ResourceApplyFailure,
			Error:  fmt.Errorf(`failed to decode condition "%s" body spec: %s`, c.Name, err.Error()),
			Value:  cty.NilVal,
		}
	}

	return core.ResourceOutput{
		Status: events.ResourceApplySuccess,
		Error:  nil,
		Value:  c.Value,
	}

}
