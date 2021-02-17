package v1

import (
	"fmt"

	"github.com/verifa/bubbly/api/common"
	"github.com/verifa/bubbly/api/core"
	"github.com/verifa/bubbly/env"
	"github.com/verifa/bubbly/events"

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
func (p *PipelineRun) Apply(bCtx *env.BubblyContext, ctx *core.ResourceContext) core.ResourceOutput {
	if err := common.DecodeBody(bCtx, p.SpecHCL.Body, &p.Spec, ctx); err != nil {
		return core.ResourceOutput{
			ID:     p.String(),
			Status: events.ResourceApplyFailure,
			Error:  fmt.Errorf(`failed to decode "%s" body spec: %s`, p.String(), err.Error()),
			Value:  cty.NilVal,
		}
	}

	_, output := common.RunResource(bCtx, ctx, p.Spec.PipelineID, p.Spec.Inputs.Value())
	return output
}

type pipelineRunSpec struct {
	Inputs     core.InputDefinitions `hcl:"input,block"`
	PipelineID string                `hcl:"pipeline,attr"`
	Interval   string                `hcl:"interval,optional"`
}
