package v1

import (
	"encoding/json"
	"fmt"

	"github.com/valocode/bubbly/api/common"
	"github.com/valocode/bubbly/api/core"
	"github.com/valocode/bubbly/env"
	"github.com/valocode/bubbly/events"

	"github.com/zclconf/go-cty/cty"
)

var _ core.Transform = (*Transform)(nil)

type Transform struct {
	*core.ResourceBlock
	Spec transformSpec
}

func NewTransform(resBlock *core.ResourceBlock) *Transform {
	return &Transform{
		ResourceBlock: resBlock,
	}
}

// Run returns ...
func (t *Transform) Run(bCtx *env.BubblyContext, ctx *core.ResourceContext) core.ResourceOutput {
	if err := common.DecodeBodyWithInputs(bCtx, t.SpecHCL.Body, &t.Spec, ctx); err != nil {
		return core.ResourceOutput{
			ID:     t.String(),
			Status: events.ResourceRunFailure,
			Error:  fmt.Errorf(`failed to decode "%s" body spec: %s`, t.String(), err.Error()),
			Value:  cty.NilVal,
		}
	}

	json, err := t.toJSON()
	if err != nil {
		return core.ResourceOutput{
			ID:     t.String(),
			Status: events.ResourceRunFailure,
			Error:  err,
			Value:  cty.NilVal,
		}
	}
	return core.ResourceOutput{
		ID:     t.String(),
		Status: events.ResourceRunSuccess,
		Error:  nil,
		Value:  cty.StringVal(string(json)),
	}
}

func (t *Transform) toJSON() ([]byte, error) {
	if t.Spec.Data == nil {
		return nil, fmt.Errorf("transform %s has not output data", t.String())
	}
	return json.Marshal(t.Spec.Data)
}

type transformSpec struct {
	Inputs core.InputDeclarations `hcl:"input,block"`
	Locals *core.LocalsDefinition `hcl:"locals,block"`
	Data   core.DataBlocks        `hcl:"data,block"`
}
