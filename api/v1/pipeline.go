package v1

import (
	"fmt"

	"github.com/verifa/bubbly/api/core"
	"github.com/zclconf/go-cty/cty"
)

var _ core.Pipeline = (*Pipeline)(nil)

type Pipeline struct {
	*core.ResourceBlock
	Spec  PipelineSpec
	Tasks Tasks
}

func NewPipeline(resBlock *core.ResourceBlock) *Pipeline {
	return &Pipeline{
		ResourceBlock: resBlock,
		Tasks:         Tasks{},
	}
}

func (p *Pipeline) Decode(decode core.DecodeResourceFn) error {
	// decode the resource spec into the importer's Spec
	if err := decode(p, p.SpecHCL.Body, &p.Spec); err != nil {
		return fmt.Errorf(`Failed to decode "%s" body spec: %s`, p.String(), err.Error())
	}

	for _, taskSpec := range p.Spec.TaskBlocks {
		println(taskSpec.Name)
		task := &Task{
			Name: taskSpec.Name,
		}
		if err := decode(p, taskSpec.Body, task); err != nil {
			return fmt.Errorf(`Failed to decode task "%s" in pipeline "%s": %s`, taskSpec.Name, p.String(), err.Error())
		}
		p.Tasks[task.Name] = task
	}
	return nil
}

func (p *Pipeline) SpecValue() core.ResourceSpec {
	return &p.Spec
}

// Output returns ...
func (p *Pipeline) Output() core.ResourceOutput {
	return core.ResourceOutput{
		Status: core.ResourceOutputSuccess,
		Error:  nil,
		Value:  cty.NilVal,
	}
}

func (p *Pipeline) Task(name string) core.Task {
	return p.Tasks[name]
}

type PipelineSpec struct {
	Inputs     InputDeclarations `hcl:"input,block"`
	TaskBlocks []TaskBlockSpec   `hcl:"task,block"`
}
