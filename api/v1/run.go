package v1

import (
	"fmt"

	"github.com/valocode/bubbly/api/common"
	"github.com/valocode/bubbly/api/core"
	"github.com/valocode/bubbly/env"
	"github.com/valocode/bubbly/events"

	"github.com/zclconf/go-cty/cty"
)

var _ core.Run = (*Run)(nil)

type Run struct {
	*core.ResourceBlock
	Spec runSpec
}

func NewRun(resBlock *core.ResourceBlock) *Run {
	return &Run{
		ResourceBlock: resBlock,
	}
}

// Run returns ...
func (p *Run) Run(bCtx *env.BubblyContext, ctx *core.ResourceContext) core.ResourceOutput {
	if err := common.DecodeBody(bCtx, p.SpecHCL.Body, &p.Spec, ctx); err != nil {
		return core.ResourceOutput{
			ID:     p.String(),
			Status: events.ResourceRunFailure,
			Error:  fmt.Errorf(`failed to decode "%s" body spec: %s`, p.String(), err.Error()),
			Value:  cty.NilVal,
		}
	}

	_, output := common.RunResourceByID(bCtx, ctx, p.Spec.ResourceID, p.Spec.Inputs.Value())
	if output.Error != nil {
		return core.ResourceOutput{
			ID:     p.String(),
			Status: events.ResourceRunFailure,
			Error:  fmt.Errorf("run failed for resource %s: %w", p.Spec.ResourceID, output.Error),
			Value:  cty.NilVal,
		}
	}
	return core.ResourceOutput{
		ID:     p.String(),
		Status: events.ResourceRunSuccess,
		Value:  cty.NilVal,
	}
}

type runSpec struct {
	Inputs     core.InputDefinitions `hcl:"input,block"`
	ResourceID string                `hcl:"resource,attr"`
	Remote     *RemoteBlockSpec      `hcl:"remote,block"`
}

// remoteBlockSpec is the type representing any "remote {...}" definition block
// in a run resource's HCL
type RemoteBlockSpec struct {
	Interval string `hcl:"interval,optional"`
}
