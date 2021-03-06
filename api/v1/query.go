package v1

import (
	"fmt"

	"github.com/valocode/bubbly/events"

	"github.com/zclconf/go-cty/cty"

	"github.com/valocode/bubbly/api/common"
	"github.com/valocode/bubbly/api/core"
	"github.com/valocode/bubbly/env"
)

var _ core.Query = (*Query)(nil)

type Query struct {
	*core.ResourceBlock
	Spec querySpec
}

func NewQuery(resBlock *core.ResourceBlock) *Query {
	return &Query{
		ResourceBlock: resBlock,
	}
}

// Run returns a core.ResourceOutput, whose Value is a cty.Value
// representation of the bubbly server's response to the q.Spec.Query string
func (q *Query) Run(bCtx *env.BubblyContext, ctx *core.ResourceContext) core.ResourceOutput {
	if err := common.DecodeBody(bCtx, q.SpecHCL.Body, &q.Spec, ctx); err != nil {
		return core.ResourceOutput{
			ID:     q.String(),
			Status: events.ResourceRunFailure,
			Error:  fmt.Errorf(`failed to decode "%s" body spec: %s`, q.String(), err.Error()),
			Value:  cty.NilVal,
		}
	}

	queryVal, err := common.QueryToCtyValue(bCtx, ctx, q.Spec.Query)
	if err != nil {
		return core.ResourceOutput{
			ID:     q.ID(),
			Status: events.ResourceRunFailure,
			Error:  err,
			Value:  cty.NilVal,
		}
	}

	return core.ResourceOutput{
		ID:     q.String(),
		Status: events.ResourceRunSuccess,
		Error:  nil,
		Value:  queryVal,
	}
}

type querySpec struct {
	Query string `hcl:"query,attr"`
}

// QueryDeclarations is a wrapper for a slice of QueryDeclaration
type QueryDeclarations []*QueryDeclaration

// QueryDeclaration is the type representing any "query "<name>" {}" declaration
// within a criteria's HCL spec block
type QueryDeclaration struct {
	Name string `hcl:",label"`
}
