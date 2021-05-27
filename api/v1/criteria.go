package v1

import (
	"errors"
	"fmt"

	"github.com/hashicorp/hcl/v2"
	"github.com/valocode/bubbly/api/common"
	"github.com/valocode/bubbly/api/core"
	"github.com/valocode/bubbly/env"
	"github.com/valocode/bubbly/events"
	"github.com/valocode/bubbly/parser"

	"github.com/zclconf/go-cty/cty"
)

var _ core.Criteria = (*Criteria)(nil)

type Criteria struct {
	*core.ResourceBlock
	Spec criteriaSpec
}

func NewCriteria(resBlock *core.ResourceBlock) *Criteria {
	return &Criteria{
		ResourceBlock: resBlock,
	}
}

// Run returns a core.ResourceOutput, whose Value is a cty.Object created from
// core.CriteriaOutput, which specifies the result (pass/fail) and a possible
// reason for the failure
func (c *Criteria) Run(bCtx *env.BubblyContext, ctx *core.ResourceContext) core.ResourceOutput {
	if err := common.DecodeBodyWithInputs(bCtx, c.SpecHCL.Body, &c.Spec, ctx); err != nil {
		return core.ResourceOutput{
			ID:     c.ID(),
			Status: events.ResourceRunFailure,
			Error:  fmt.Errorf(`failed to decode "%s" body spec: %s`, c.String(), err.Error()),
			Value:  cty.NilVal,
		}
	}

	queryVal, err := common.QueryToCtyValue(bCtx, ctx, c.Spec.Query)
	if err != nil {
		return core.ResourceOutput{
			ID:     c.ID(),
			Status: events.ResourceRunFailure,
			Error:  err,
			Value:  cty.NilVal,
		}
	}

	ctx.State.Insert("query", queryVal)
	queryOutput := ctx.State.ValueWithPath(nil)
	condInputs := core.AppendInputObjects(queryOutput, ctx.Inputs)
	if len(c.Spec.Conditions) == 0 {
		return core.ResourceOutput{
			ID:     c.ID(),
			Status: events.ResourceRunFailure,
			Error:  errors.New("criteria must have at least one condition"),
			Value:  cty.NilVal,
		}
	}

	// Create the result and set it to true/success by default
	var result = core.CriteriaResult{Result: true}
	for _, condition := range c.Spec.Conditions {
		resultVal, err := parser.ExpressionValue(condition.Value, condInputs)
		if err != nil {
			return core.ResourceOutput{
				ID:     c.ID(),
				Status: events.ResourceRunFailure,
				Error:  fmt.Errorf("error evaluating condition: %w", err),
				Value:  cty.NilVal,
			}
		}
		// If the result is not a boolean, we have an issue
		if !resultVal.Type().Equals(cty.Bool) {
			return core.ResourceOutput{
				ID:     c.ID(),
				Status: events.ResourceRunFailure,
				Error:  fmt.Errorf("condition is not boolean type: %s", resultVal.Type().FriendlyName()),
				Value:  cty.NilVal,
			}
		}

		// If the condition was not met, then we can stop here and return the reason
		if resultVal.False() {
			result.Result = false
			result.Reason = condition.Message
			break
		}
	}
	// If all conditions have been evaluated return the criteria result
	ctyResult, err := result.Value()
	if err != nil {
		return core.ResourceOutput{
			ID:     c.ID(),
			Status: events.ResourceRunFailure,
			Error:  err,
			Value:  cty.NilVal,
		}
	}
	return core.ResourceOutput{
		ID:     c.ID(),
		Status: events.ResourceRunSuccess,
		Value:  ctyResult,
	}
}

type criteriaSpec struct {
	Inputs     core.InputDeclarations `hcl:"input,block"`
	Query      string                 `hcl:"query,attr"`
	Conditions []conditionSpec        `hcl:"condition,block"`
}

type conditionSpec struct {
	Name    string         `hcl:",label"`
	Message string         `hcl:"message,optional"`
	Value   hcl.Expression `hcl:"value,attr"`
}
