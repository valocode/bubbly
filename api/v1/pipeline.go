package v1

import (
	"fmt"

	"github.com/verifa/bubbly/api/core"
	"github.com/verifa/bubbly/env"
	"github.com/zclconf/go-cty/cty"
)

var _ core.Pipeline = (*Pipeline)(nil)

type Pipeline struct {
	*core.ResourceBlock
	Spec  pipelineSpec
	Tasks core.Tasks
}

func NewPipeline(resBlock *core.ResourceBlock) *Pipeline {
	return &Pipeline{
		ResourceBlock: resBlock,
		Tasks:         core.Tasks{},
	}
}

func (p *Pipeline) SpecValue() core.ResourceSpec {
	return &p.Spec
}

// Apply returns ...
func (p *Pipeline) Apply(bCtx *env.BubblyContext, ctx *core.ResourceContext) core.ResourceOutput {

	if err := ctx.DecodeBody(bCtx, p, p.SpecHCL.Body, &p.Spec); err != nil {
		return core.ResourceOutput{
			Status: core.ResourceOutputFailure,
			Error:  fmt.Errorf(`Failed to decode "%s" body spec: %s`, p.String(), err.Error()),
			Value:  cty.NilVal,
		}
	}

	for idx, taskSpec := range p.Spec.TaskBlocks {
		bCtx.Logger.Debug().Msgf("Applying task: %s", taskSpec.Name)
		t := NewTask(taskSpec)

		err := t.Apply(bCtx, ctx)
		if err != nil {
			return core.ResourceOutput{
				Status: core.ResourceOutputFailure,
				Error:  fmt.Errorf(`Failed to apply task "%s" with index %d in pipeline "%s": %s"`, taskSpec.Name, idx, p.String(), err.Error()),
				Value:  cty.NilVal,
			}
		}
		p.Tasks[t.Name] = t

	}

	return core.ResourceOutput{
		Status: core.ResourceOutputSuccess,
		Error:  nil,
		Value:  cty.NilVal,
	}
}

type pipelineSpec struct {
	Inputs     InputDeclarations `hcl:"input,block"`
	TaskBlocks []*taskBlockSpec  `hcl:"task,block"`
}
