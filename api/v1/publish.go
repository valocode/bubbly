package v1

import (
	"github.com/verifa/bubbly/api/core"
	"github.com/zclconf/go-cty/cty"
)

var _ core.Publish = (*Publish)(nil)

type Publish struct {
	*core.ResourceBlock
	Spec PublishSpec
}

func NewPublish(resBlock *core.ResourceBlock) *Publish {
	return &Publish{
		ResourceBlock: resBlock,
	}
}

func (p *Publish) Decode(decode core.DecodeResourceFn) error {
	return nil
}

func (p *Publish) SpecValue() core.ResourceSpec {
	return &p.Spec
}

// Output returns ...
func (p *Publish) Output() core.ResourceOutput {
	return core.ResourceOutput{
		Status: core.ResourceOutputSuccess,
		Error:  nil,
		Value:  cty.NilVal,
	}
}

type PublishSpec struct {
}
