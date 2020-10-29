package v1

import (
	"fmt"

	"github.com/verifa/bubbly/api/core"
	"github.com/zclconf/go-cty/cty"
)

var _ core.PipelineRun = (*PipelineRun)(nil)

type PipelineRun struct {
	*core.ResourceBlock
	Spec pipelineRunSpec
}

func NewPipelineRun(resBlock *core.ResourceBlock) *PipelineRun {
	return &PipelineRun{
		ResourceBlock: resBlock,
	}
}

func (p *PipelineRun) SpecValue() core.ResourceSpec {
	return &p.Spec
}

// Apply returns ...
func (p *PipelineRun) Apply(ctx *core.ResourceContext) core.ResourceOutput {
	// decode the resource spec into the pipeline runs's Spec
	if err := ctx.DecodeBody(p, p.SpecHCL.Body, &p.Spec); err != nil {
		return core.ResourceOutput{
			Status: core.ResourceOutputFailure,
			Error:  fmt.Errorf(`Failed to decode "%s" body spec: %s`, p.String(), err.Error()),
			Value:  cty.NilVal,
		}
	}

	pipeline, err := ctx.GetResource(core.PipelineResourceKind, p.Spec.PipelineID)
	if err != nil {
		return core.ResourceOutput{
			Status: core.ResourceOutputFailure,
			Error:  fmt.Errorf(`Pipeline "%s" does not exist: %s`, p.Spec.PipelineID, err.Error()),
			Value:  cty.NilVal,
		}
	}

	out := pipeline.Apply(ctx.NewContext(p.Spec.Inputs.Value()))

	return core.ResourceOutput{
		Status: out.Status,
		Error:  out.Error,
		Value:  cty.NilVal,
	}
}

type pipelineRunSpec struct {
	Inputs     InputDefinitions `hcl:"input,block"`
	PipelineID string           `hcl:"pipeline,attr"`
}
