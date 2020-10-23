package v1

import (
	"fmt"

	"github.com/verifa/bubbly/api/core"
	"github.com/zclconf/go-cty/cty"
)

var _ core.PipelineRun = (*PipelineRun)(nil)

type PipelineRun struct {
	*core.ResourceBlock
	Spec PipelineRunSpec
}

func NewPipelineRun(resBlock *core.ResourceBlock) *PipelineRun {
	return &PipelineRun{
		ResourceBlock: resBlock,
	}
}

func (p *PipelineRun) Decode(decode core.DecodeResourceFn) error {
	// decode the resource spec into the pipeline runs's Spec
	if err := decode(p, p.SpecHCL.Body, &p.Spec); err != nil {
		return fmt.Errorf(`Failed to decode "%s" body spec: %s`, p.String(), err.Error())
	}
	return nil
}

func (p *PipelineRun) SpecValue() core.ResourceSpec {
	return &p.Spec
}

// Output returns ...
func (p *PipelineRun) Output() core.ResourceOutput {
	return core.ResourceOutput{
		Status: core.ResourceOutputSuccess,
		Error:  nil,
		Value:  cty.NilVal,
	}
}

func (p *PipelineRun) Pipeline() string {
	return p.Spec.PipelineID
}

func (p *PipelineRun) Inputs() cty.Value {
	inputs := map[string]cty.Value{}
	for _, input := range p.Spec.Inputs {
		inputs[input.Name] = input.Value
	}
	return cty.ObjectVal(
		map[string]cty.Value{
			"input": cty.ObjectVal(inputs),
		},
	)
}

type PipelineRunSpec struct {
	Inputs     InputDefinitions `hcl:"input,block"`
	PipelineID string           `hcl:"pipeline,attr"`
}
