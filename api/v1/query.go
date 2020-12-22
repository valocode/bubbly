package v1

import (
	"fmt"

	"github.com/verifa/bubbly/client"
	"github.com/verifa/bubbly/parser"

	"github.com/verifa/bubbly/api/core"
	"github.com/verifa/bubbly/env"
	"github.com/zclconf/go-cty/cty"
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

// Apply returns a core.ResourceOutput, whose Value is a cty.Value
// representation of the bubbly server's response to the q.Spec.Query string
func (q *Query) Apply(bCtx *env.BubblyContext, ctx *core.ResourceContext) core.ResourceOutput {
	p := parser.WithInputs(bCtx, ctx.Inputs)
	if err := p.Scope.DecodeExpandBody(bCtx, q.SpecHCL.Body, &q.Spec); err != nil {
		return core.ResourceOutput{
			Status: core.ResourceOutputFailure,
			Error:  fmt.Errorf("failed to decode query body spec: %w", err),
			Value:  cty.NilVal,
		}
	}

	c, err := client.New(bCtx)

	if err != nil {
		return core.ResourceOutput{
			Status: core.ResourceOutputFailure,
			Error:  fmt.Errorf("failed to establish bubbly client: %w", err),

			Value: cty.NilVal,
		}
	}
	// run the query against the bubbly store
	byteRes, err := c.Query(bCtx, q.Spec.Query)

	if err != nil {
		return core.ResourceOutput{
			Status: core.ResourceOutputFailure,
			Error:  fmt.Errorf(`failed while running query "%s": %w`, q.Spec.Query, err),

			Value: cty.NilVal,
		}
	}

	bCtx.Logger.Debug().RawJSON("response", byteRes).Msg("received query response from bubbly server")

	// because the return from the bubbly server
	// is simply a []byte representation of the graphql-go
	// interface{} response, we can convert it directly to a
	// cty.StringVal for use as core.ResourceOutput.Value
	strNormalized := cty.NormalizeString(string(byteRes))

	cVal := cty.StringVal(strNormalized)

	// ctx.InsertValue(bCtx, cVal, []string{"self", "query", q.Name()})

	return core.ResourceOutput{
		Status: core.ResourceOutputSuccess,
		Error:  nil,
		Value:  cVal,
	}
}

func (q *Query) SpecValue() core.ResourceSpec {
	return &q.Spec
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
