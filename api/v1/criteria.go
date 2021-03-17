package v1

import (
	"fmt"

	"github.com/valocode/bubbly/api/common"
	"github.com/valocode/bubbly/api/core"
	"github.com/valocode/bubbly/env"
	"github.com/valocode/bubbly/events"

	"github.com/zclconf/go-cty/cty"
)

var _ core.Criteria = (*Criteria)(nil)

type Criteria struct {
	*core.ResourceBlock
	Spec       criteriaSpec
	Queries    core.Queries
	Conditions core.Conditions
	Operation  core.Operation
}

func NewCriteria(resBlock *core.ResourceBlock) *Criteria {
	return &Criteria{
		ResourceBlock: resBlock,
		Queries:       core.Queries{},
		Conditions:    core.Conditions{},
	}
}

// Apply returns a core.ResourceOutput, whose Value is a cty.BoolVal indicating
// success of failure of the criteria's operation
func (c *Criteria) Apply(bCtx *env.BubblyContext, ctx *core.ResourceContext) core.ResourceOutput {
	if err := common.DecodeBody(bCtx, c.SpecHCL.Body, &c.Spec, ctx); err != nil {
		return core.ResourceOutput{
			Status: events.ResourceApplyFailure,
			Error:  fmt.Errorf(`failed to decode "%s" body spec: %s`, c.String(), err.Error()),
			Value:  cty.NilVal,
		}
	}

	// We first process the queries and then add their outputs to the
	// EvalContext

	// loop through the query blocks, apply the referenced query resources
	// to obtain their cty.Value outputs, and insert this into the EvalContext
	// This allows us to later resolve condition blocks which reference these
	// values (self.query.<name>.value)
	for idx, querySpec := range c.Spec.Queries {
		bCtx.Logger.Debug().Msgf("Applying query: %s", querySpec.Name)

		// use this resourceID to get the underlying Resource
		resID := fmt.Sprintf("%s/%s", string(core.QueryResourceKind), querySpec.Name)
		resource, output := common.RunResource(bCtx, ctx, resID, cty.NilVal)

		if output.Error != nil {
			return core.ResourceOutput{
				Status: events.ResourceApplyFailure,
				Error:  fmt.Errorf(`failed to apply query "%s" with index %d in criteria "%s": %w"`, querySpec.Name, idx, c.String(), output.Error),
				Value:  cty.NilVal,
			}
		}

		// insert the output into the EvalContext
		// TODO
		// p.Scope.InsertValue(bCtx, output.Output(), []string{"self", string(core.QueryResourceKind), resource.Name()})

		bCtx.Logger.Debug().Str("query", resource.String()).Str("output_status", output.Status.String()).Str("output_value", output.Value.GoString()).Msg("query successfully processed")
		c.Queries[resource.Name()] = resource

	}

	// we loop through the condition blocks, and apply the condition.
	// The apply uses the previously calculated query values
	// that we inserted into the EvalContext to decode the condition HCL block.
	// The core.ResourceOutput of a condition is then inserted into the
	// EvalContext for later use when decoding the operation block
	for idx, conditionSpec := range c.Spec.Conditions {
		bCtx.Logger.Debug().Msgf("Evaluating condition: %s", conditionSpec.Name)

		condition := NewCondition(conditionSpec)

		output := condition.Apply(
			bCtx,
			core.NewResourceContext(cty.NilVal, ctx.NewResource),
		)

		if output.Error != nil {
			return core.ResourceOutput{
				Status: events.ResourceApplyFailure,
				Error:  fmt.Errorf(`failed to apply condition "%s" with index %d in criteria "%s": %w"`, conditionSpec.Name, idx, c.String(), output.Error),
				Value:  cty.NilVal,
			}
		}

		// insert the output into the EvalContext. This allows us to use
		// the result when decoding the operation
		// p.Scope.InsertValue(bCtx, output.Output(), []string{"self", "condition", condition.Name})
		// add the output of the task to the parser
		ctx.State.Insert(condition.Name, output.Output())

		bCtx.Logger.Debug().
			Str("condition", condition.Name).
			Str("output_status", output.Status.String()).
			Str("output_value", output.Value.GoString()).
			Msg("condition successfully processed")

		c.Conditions[condition.Name] = condition

	}

	// finally, decode the operation block to produce the final output of the
	// criteria resource
	operationSpec := c.Spec.Operation
	operation := NewOperation(operationSpec)

	output := operation.Apply(
		bCtx,
		core.NewResourceContext(cty.NilVal, ctx.NewResource),
	)

	if output.Error != nil {
		return core.ResourceOutput{
			Status: events.ResourceApplyFailure,
			Error:  fmt.Errorf(`failed to apply operation "%s" in criteria "%s": %w"`, operationSpec.Name, c.String(), output.Error),
			Value:  cty.NilVal,
		}
	}

	bCtx.Logger.Debug().
		Str("operation", operation.Name()).
		Str("output_status", output.Status.String()).
		Str("output_value", output.Value.GoString()).
		Msg("operation successfully processed")

	c.Operation = operation

	return core.ResourceOutput{
		Status: events.ResourceApplySuccess,
		Error:  nil,
		Value:  operation.Value,
	}
}

func (c *Criteria) SpecValue() core.ResourceSpec {
	return &c.Spec
}

type criteriaSpec struct {
	Queries    QueryDeclarations     `hcl:"query,block"`
	Conditions []*conditionBlockSpec `hcl:"condition,block"`
	// TODO: Consider whether we want > 1 operation (their outputs
	// could be used to reference each other as a way of splitting longer/
	// complex operations)
	Operation *operationBlockSpec `hcl:"operation,block"`
}
