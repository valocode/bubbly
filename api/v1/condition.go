package v1

import (
	"fmt"

	"github.com/hashicorp/hcl/v2"
	"github.com/verifa/bubbly/api/core"
	"github.com/verifa/bubbly/env"
	"github.com/verifa/bubbly/parser"
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
	p := parser.WithInputs(bCtx, ctx.Inputs)
	if err := p.Scope.DecodeExpandBody(bCtx, c.conditionBlockSpec.Body, c); err != nil {
		return core.ResourceOutput{
			Status: core.ResourceOutputFailure,
			Error:  fmt.Errorf(`failed to decode Condition "%s": %w`, c.Name, err),
			Value:  cty.NilVal,
		}
	}

	return core.ResourceOutput{
		Status: core.ResourceOutputSuccess,
		Error:  nil,
		Value:  c.Value,
	}

}
