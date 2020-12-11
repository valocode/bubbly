package v1

import (
	"fmt"

	"github.com/hashicorp/hcl/v2"
	"github.com/verifa/bubbly/api/core"
	"github.com/verifa/bubbly/env"
	"github.com/zclconf/go-cty/cty"
)

var _ core.Operation = (*Operation)(nil)

type operationBlockSpec struct {
	Name string   `hcl:",label"`
	Body hcl.Body `hcl:",remain"`
}

type Operation struct {
	*operationBlockSpec
	Value cty.Value `hcl:"value,attr"`
}

func NewOperation(operationBlock *operationBlockSpec) *Operation {
	return &Operation{
		operationBlockSpec: operationBlock,
	}
}

// Apply returns the output from the decoding of the operation's hcl.Body into
// an Operation struct. Namely, the o.Value, which represents the final
// criteria's return value
func (o *Operation) Apply(bCtx *env.BubblyContext, ctx *core.ResourceContext) core.ResourceOutput {
	if err := ctx.DecodeBody(bCtx, o.operationBlockSpec.Body, o); err != nil {
		return core.ResourceOutput{
			Status: core.ResourceOutputFailure,
			Error:  fmt.Errorf(`Failed to decode Operation "%s": %w`, o.Name(), err),
			Value:  cty.NilVal,
		}
	}

	return core.ResourceOutput{
		Status: core.ResourceOutputSuccess,
		Error:  nil,
		Value:  o.Value,
	}

}

// String returns a human-friendly string ID for the operation resource
func (o *Operation) String() string {
	return fmt.Sprintf(
		"%s.%s.%s",
		o.APIVersion(), o.Kind(), o.Name(),
	)
}

// APIVersion returns the resource's API version
func (o *Operation) APIVersion() core.APIVersion {
	return core.APIVersion("v1")
}

// Kind returns the resource kind
func (o *Operation) Kind() core.ResourceKind {
	return core.ResourceKind("operation")
}

// Name returns the name of the operation
func (o *Operation) Name() string {
	return o.operationBlockSpec.Name
}

// JSON returns a JSON representation of this operation block using the given
// ResourceContext.
// TODO
func (o *Operation) JSON(ctx *core.ResourceContext) ([]byte, error) {
	return nil, nil
}
