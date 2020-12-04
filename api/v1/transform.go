package v1

import (
	"encoding/json"
	"fmt"

	"github.com/verifa/bubbly/api/core"
	"github.com/verifa/bubbly/env"
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

// Apply returns ...
func (t *Transform) Apply(bCtx *env.BubblyContext, ctx *core.ResourceContext) core.ResourceOutput {
	if err := ctx.DecodeBody(bCtx, t.SpecHCL.Body, &t.Spec); err != nil {
		return core.ResourceOutput{
			Status: core.ResourceOutputFailure,
			Error:  fmt.Errorf(`Failed to decode transform body spec: %s`, err.Error()),
			Value:  cty.NilVal,
		}
	}

	json, err := t.toJSON()
	if err != nil {
		return core.ResourceOutput{
			Status: core.ResourceOutputFailure,
			Error:  err,
			Value:  cty.NilVal,
		}
	}
	return core.ResourceOutput{
		Status: core.ResourceOutputSuccess,
		Error:  nil,
		Value:  cty.StringVal(string(json)),
	}
}

func (t *Transform) toJSON() ([]byte, error) {
	if t.Spec.Data == nil {
		return nil, fmt.Errorf("Transform %s has not output data", t.String())
	}
	return json.Marshal(t.Spec.Data)
}

func (t *Transform) SpecValue() core.ResourceSpec {
	return &t.Spec
}

type transformSpec struct {
	Inputs InputDeclarations `hcl:"input,block"`
	Data   core.DataBlocks   `hcl:"data,block"`
}
