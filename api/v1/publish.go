package v1

import (
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/verifa/bubbly/api/core"
	"github.com/zclconf/go-cty/cty"
)

var _ core.Publish = (*Publish)(nil)

type Publish struct {
	*core.ResourceBlock
	Spec publishSpec
}

func NewPublish(resBlock *core.ResourceBlock) *Publish {
	return &Publish{
		ResourceBlock: resBlock,
	}
}

func (p *Publish) SpecValue() core.ResourceSpec {
	return &p.Spec
}

// Apply returns ...
func (p *Publish) Apply(ctx *core.ResourceContext) core.ResourceOutput {
	if err := ctx.DecodeBody(p, p.SpecHCL.Body, &p.Spec); err != nil {
		return core.ResourceOutput{
			Status: core.ResourceOutputFailure,
			Error:  fmt.Errorf(`Failed to decode "%s" body spec: %s`, p.String(), err.Error()),
			Value:  cty.NilVal,
		}
	}
	// TODO: need to add the actual publish
	log.Warn().Msgf("TODO: Should publish this JSON: %s", p.Spec.Data)

	return core.ResourceOutput{
		Status: core.ResourceOutputSuccess,
		Error:  nil,
		Value:  cty.NilVal,
	}
}

type publishSpec struct {
	Inputs InputDeclarations `hcl:"input,block"`
	Data   string            `hcl:"data,attr"`
}
