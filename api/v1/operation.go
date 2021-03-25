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

var _ core.Operation = (*Operation)(nil)

type operationBlockSpec struct {
	Name string   `hcl:",label"`
	Body hcl.Body `hcl:",remain"`
}

type Operation struct {
	*operationBlockSpec
	Value cty.Value `hcl:"value,attr"`
}

func NewOperation(operationBlock *operationBlockSpec) *Operation {
	return &Operation{
		operationBlockSpec: operationBlock,
	}
}

// Apply returns the output from the decoding of the operation's hcl.Body into
// an Operation struct. Namely, the o.Value, which represents the final
// criteria's return value
func (o *Operation) Apply(bCtx *env.BubblyContext, ctx *core.ResourceContext) core.ResourceOutput {
	if err := common.DecodeBody(bCtx, o.operationBlockSpec.Body, o, ctx); err != nil {
		return core.ResourceOutput{
			Status: events.ResourceApplyFailure,
			Error:  fmt.Errorf(`failed to decode operation "%s" body spec: %s`, o.Name(), err.Error()),
			Value:  cty.NilVal,
		}
	}

	if o.Value == cty.BoolVal(false) {
		return core.ResourceOutput{
			Status: events.ResourceApplyFailure,
			Error:  nil,
			Value:  o.Value,
		}
	}

	return core.ResourceOutput{
		Status: events.ResourceApplySuccess,
		Error:  nil,
		Value:  o.Value,
	}
}

// Name returns the name of the operation
func (o *Operation) Name() string {
	return o.operationBlockSpec.Name
}
