package v1

import (
	"fmt"

	"github.com/hashicorp/hcl/v2"
	"github.com/verifa/bubbly/api/core"
	"github.com/verifa/bubbly/env"
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
	if err := ctx.DecodeBody(bCtx, c.conditionBlockSpec.Body, c); err != nil {
		return core.ResourceOutput{
			Status: core.ResourceOutputFailure,
			Error:  fmt.Errorf(`Failed to decode Condition "%s": %w`, c.Name(), err),
			Value:  cty.NilVal,
		}
	}

	return core.ResourceOutput{
		Status: core.ResourceOutputSuccess,
		Error:  nil,
		Value:  c.Value,
	}

}

// String returns a human-friendly string ID for the condition resource
func (c *Condition) String() string {
	return fmt.Sprintf(
		"%s.%s.%s",
		c.APIVersion(), c.Kind(), c.Name(),
	)
}

// APIVersion returns the resource's API version
func (c *Condition) APIVersion() core.APIVersion {
	return core.APIVersion("v1")
}

// Kind returns the resource kind
func (c *Condition) Kind() core.ResourceKind {
	return core.ResourceKind("condition")
}

// Name returns the name of the condition
func (c *Condition) Name() string {
	return c.conditionBlockSpec.Name
}

// JSON returns a JSON representation of this condition block using the given
// ResourceContext.
// TODO
func (c *Condition) JSON(ctx *core.ResourceContext) ([]byte, error) {
	return nil, nil
}
