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

func (p *Run) SpecValue() core.ResourceSpec {
	return &p.Spec
}

// Apply returns ...
func (p *Run) Apply(bCtx *env.BubblyContext, ctx *core.ResourceContext) core.ResourceOutput {
	if err := common.DecodeBody(bCtx, p.SpecHCL.Body, &p.Spec, ctx); err != nil {
		return core.ResourceOutput{
			ID:     p.String(),
			Status: events.ResourceApplyFailure,
			Error:  fmt.Errorf(`failed to decode "%s" body spec: %s`, p.String(), err.Error()),
			Value:  cty.NilVal,
		}
	}

	_, output := common.RunResource(bCtx, ctx, p.Spec.ResourceID, p.Spec.Inputs.Value())
	return output
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
