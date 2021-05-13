package v1

import (
	"fmt"

	"github.com/valocode/bubbly/api/common"
	"github.com/valocode/bubbly/api/core"
	"github.com/valocode/bubbly/env"
	"github.com/valocode/bubbly/events"

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

	if err := common.DecodeBodyWithInputs(bCtx, p.SpecHCL.Body, &p.Spec, ctx); err != nil {
		return core.ResourceOutput{
			ID:     p.String(),
			Status: events.ResourceApplyFailure,
			Error:  fmt.Errorf(`failed to decode "%s" body spec: %s`, p.String(), err.Error()),
			Value:  cty.NilVal,
		}
	}

	for idx, taskSpec := range p.Spec.TaskBlocks {
		bCtx.Logger.Debug().Msgf("Applying task: %s", taskSpec.Name)
		t := NewTask(taskSpec)

		// create the run ResourceContext for the SubResource to apply
		inputs := core.AppendInputObjects(ctx.State.ValueWithPath([]string{"query"}), ctx.Inputs)
		runCtx := core.SubResourceContext(inputs, ctx)
		output := t.Apply(bCtx, runCtx)
		if output.Error != nil {
			return core.ResourceOutput{
				ID:     p.String(),
				Status: events.ResourceApplyFailure,
				Error:  fmt.Errorf(`failed to apply task "%s" with index %d in pipeline "%s": %w"`, taskSpec.Name, idx, p.String(), output.Error),
				Value:  cty.NilVal,
			}
		}

		// add the output of the task to the parser
		ctx.State.Insert(t.Name(), output.Value)

		p.Tasks[t.Name()] = t
	}

	return core.ResourceOutput{
		ID:     p.String(),
		Status: events.ResourceApplySuccess,
		Error:  nil,
		Value:  cty.NilVal,
	}
}

type pipelineSpec struct {
	Inputs     core.InputDeclarations `hcl:"input,block"`
	TaskBlocks []*taskBlockSpec       `hcl:"task,block"`
}
