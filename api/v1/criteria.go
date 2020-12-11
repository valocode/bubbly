package v1

import (
	"fmt"

	"github.com/verifa/bubbly/api/core"
	"github.com/verifa/bubbly/env"
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
	if err := ctx.DecodeBody(bCtx, c.SpecHCL.Body, &c.Spec); err != nil {
		return core.ResourceOutput{
			Status: core.ResourceOutputFailure,
			Error:  fmt.Errorf("Failed to decode criteria body spec: %w", err),
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

		// get the resource ID of the underlying query referenced
		// by the QueryDeclaration
		queryResID := core.NewResourceIDFromString(fmt.Sprint(string(core.QueryResourceKind), "/", querySpec.Name))

		// use this resourceID to get the underlying Resource
		queryRes, err := ctx.GetResource(queryResID.Kind, queryResID.Name)

		if err != nil {
			return core.ResourceOutput{
				Status: core.ResourceOutputFailure,
				Error:  fmt.Errorf(`Could not find query resource "%s" referenced by Criteria "%s": %w`, queryResID.String(), c.Name(), err),
				Value:  cty.NilVal,
			}
		}

		// applying a query resource updates the ctx, such that the
		// EvalContext includes a "self.query.<name>.value" []Traversal.
		// We then use this when decoding the condition blocks
		output := queryRes.Apply(bCtx, ctx)

		if output.Error != nil {
			return core.ResourceOutput{
				Status: core.ResourceOutputFailure,
				Error:  fmt.Errorf(`Failed to apply query "%s" with index %d in criteria "%s": %w"`, querySpec.Name, idx, c.String(), output.Error),
				Value:  cty.NilVal,
			}
		}

		// insert the output into the EvalContext
		ctx.InsertValue(bCtx, output.Output(), []string{"self", string(core.QueryResourceKind), queryRes.Name()})

		bCtx.Logger.Debug().Str("query", queryRes.String()).Str("output_status", string(output.Status)).Str("output_value", output.Value.GoString()).Msg("query successfully processed")
		c.Queries[queryRes.Name()] = queryRes

	}

	// we loop through the condition blocks, and apply the condition.
	// The apply uses the previously calculated query values
	// that we inserted into the EvalContext to decode the condition HCL block.
	// The core.ResourceOutput of a condition is then inserted into the
	// EvalContext for later use when decoding the operation block
	for idx, conditionSpec := range c.Spec.Conditions {
		bCtx.Logger.Debug().Msgf("Evaluating condition: %s", conditionSpec.Name)

		condition := NewCondition(conditionSpec)

		output := condition.Apply(bCtx, ctx)

		if output.Error != nil {
			return core.ResourceOutput{
				Status: core.ResourceOutputFailure,
				Error:  fmt.Errorf(`Failed to apply condition "%s" with index %d in criteria "%s": %w"`, conditionSpec.Name, idx, c.String(), output.Error),
				Value:  cty.NilVal,
			}
		}

		// insert the output into the EvalContext. This allows us to use
		// the result when decoding the operation
		ctx.InsertValue(bCtx, output.Output(), []string{"self", "condition", condition.Name()})

		bCtx.Logger.Debug().
			Str("condition", condition.String()).
			Str("output_status", string(output.Status)).
			Str("output_value", output.Value.GoString()).
			Msg("condition successfully processed")

		c.Conditions[condition.Name()] = condition

	}

	// finally, decode the operation block to produce the final output of the
	// criteria resource
	operationSpec := c.Spec.Operation
	operation := NewOperation(operationSpec)

	output := operation.Apply(bCtx, ctx)

	if output.Error != nil {
		return core.ResourceOutput{
			Status: core.ResourceOutputFailure,
			Error:  fmt.Errorf(`Failed to apply operation "%s" in criteria "%s": %w"`, operationSpec.Name, c.String(), output.Error),
			Value:  cty.NilVal,
		}
	}

	bCtx.Logger.Debug().
		Str("operation", operation.String()).
		Str("output_status", string(output.Status)).
		Str("output_value", output.Value.GoString()).
		Msg("operation successfully processed")

	c.Operation = operation

	return core.ResourceOutput{
		Status: core.ResourceOutputSuccess,
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
