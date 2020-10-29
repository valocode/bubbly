package v1

import (
	"encoding/json"
	"fmt"

	"github.com/verifa/bubbly/api/core"
	"github.com/zclconf/go-cty/cty"
)

var _ core.Translator = (*Translator)(nil)

type Translator struct {
	*core.ResourceBlock
	Spec translatorSpec
}

func NewTranslator(resBlock *core.ResourceBlock) *Translator {
	return &Translator{
		ResourceBlock: resBlock,
	}
}

// Apply returns ...
func (t *Translator) Apply(ctx *core.ResourceContext) core.ResourceOutput {
	if err := ctx.DecodeBody(t, t.SpecHCL.Body, &t.Spec); err != nil {
		return core.ResourceOutput{
			Status: core.ResourceOutputFailure,
			Error:  fmt.Errorf(`Failed to decode translator body spec: %s`, err.Error()),
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

func (t *Translator) toJSON() ([]byte, error) {
	if t.Spec.Data == nil {
		return nil, fmt.Errorf("Translator %s has not output data", t.String())
	}
	return json.Marshal(t.Spec.Data)
}

func (t *Translator) SpecValue() core.ResourceSpec {
	return &t.Spec
}

type translatorSpec struct {
	Inputs InputDeclarations `hcl:"input,block"`
	Data   core.DataBlocks   `hcl:"data,block"`
}
