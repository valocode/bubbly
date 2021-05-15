package v1

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/graphql-go/graphql"
	"github.com/hashicorp/hcl/v2"
	"github.com/valocode/bubbly/api/common"
	"github.com/valocode/bubbly/api/core"
	"github.com/valocode/bubbly/client"
	"github.com/valocode/bubbly/env"
	"github.com/valocode/bubbly/events"
	"github.com/valocode/bubbly/parser"

	"github.com/zclconf/go-cty/cty"
	ctyjson "github.com/zclconf/go-cty/cty/json"
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

// Apply returns a core.ResourceOutput, whose Value is a cty.Object created from
// core.CriteriaOutput, which specifies the result (pass/fail) and a possible
// reason for the failure
func (c *Criteria) Apply(bCtx *env.BubblyContext, ctx *core.ResourceContext) core.ResourceOutput {
	if err := common.DecodeBodyWithInputs(bCtx, c.SpecHCL.Body, &c.Spec, ctx); err != nil {
		return core.ResourceOutput{
			ID:     c.ID(),
			Status: events.ResourceRunFailure,
			Error:  fmt.Errorf(`failed to decode "%s" body spec: %s`, c.String(), err.Error()),
			Value:  cty.NilVal,
		}
	}

	queryVal, err := c.queryToCtyValue(bCtx, ctx)
	if err != nil {
		return core.ResourceOutput{
			ID:     c.ID(),
			Status: events.ResourceRunFailure,
			// The returned error should contain enough context about what went wrong
			Error: err,
			Value: cty.NilVal,
		}
	}
	fmt.Printf("queryVal: %s\n", queryVal.GoString())

	ctx.State.Insert("query", queryVal)
	queryOutput := ctx.State.ValueWithPath(nil)
	condInputs := core.AppendInputObjects(queryOutput, ctx.Inputs)
	fmt.Printf("condInputs: %s\n", condInputs.GoString())
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
func (c *Criteria) SpecValue() core.ResourceSpec {
	return &c.Spec
}

// queryToCtyValue executes the criteria query and creates a cty.Value containing
// the results, which can then be used to evaluate the criteria conditions
func (c *Criteria) queryToCtyValue(bCtx *env.BubblyContext, ctx *core.ResourceContext) (cty.Value, error) {
	client, err := client.New(bCtx)
	if err != nil {
		return cty.NilVal, fmt.Errorf("error creating bubbly client: %w", err)
	}
	defer client.Close()

	bytes, err := client.Query(bCtx, ctx.Auth, c.Spec.Query)
	if err != nil {
		return cty.NilVal, fmt.Errorf("error executing query: %w", err)
	}

	var result graphql.Result
	if err := json.Unmarshal(bytes, &result); err != nil {
		return cty.NilVal, fmt.Errorf("error unmarshalling query result: %w", err)
	}
	if result.HasErrors() {
		return cty.NilVal, fmt.Errorf("received errors from query: %v", result.Errors)
	}

	// Operation: we are given a result from the GraphQL query, which comes as
	// an interface{} but is a map[string]interface{}. The structure of the map
	// and slices therein, depends on the GraphQL query and also the bubbly
	// schema that has been applied.
	// This data needs to be converted into a cty.Value so that the conditions
	// for this criteria can be evaluated.
	// Right now the easiest way seems to be to get the implied type using JSON,
	// but this has limitations/implications... such as missing values in JSON
	// may lead to the wrong type. Probably we need to write some of our own
	// logic here in the future, or have the ability to create a cty.Type from
	// the GraphQL query and Bubbly schema. Anyway, until such an issue arises...
	dBytes, err := json.Marshal(result.Data)
	if err != nil {
		return cty.NilVal, fmt.Errorf("error marshalling query result data: %w", err)
	}
	dType, err := ctyjson.ImpliedType(dBytes)
	if err != nil {
		return cty.NilVal, fmt.Errorf("could not imply type from query result: %w", err)
	}
	queryVal, err := ctyjson.Unmarshal(dBytes, dType)
	if err != nil {
		return cty.NilVal, fmt.Errorf("could not unmarshal result data into cty value: %w", err)
	}
	return queryVal, nil
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
