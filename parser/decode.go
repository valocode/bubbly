package parser

import (
	"fmt"
	"reflect"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/ext/dynblock"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/verifa/bubbly/env"
	"github.com/zclconf/go-cty/cty"
)

func DecodeBody(bCtx *env.BubblyContext, body hcl.Body, val interface{}, inputs cty.Value) error {
	if diags := gohcl.DecodeBody(body, newEvalContext(inputs), val); diags.HasErrors() {
		return fmt.Errorf(`failed to decode body using type "%s": %s`, reflect.TypeOf(val).String(), diags.Error())
	}
	return nil
}

func DecodeExpandBody(bCtx *env.BubblyContext, body hcl.Body, val interface{}, inputs cty.Value) error {
	eCtx := newEvalContext(inputs)
	// we have to expand before we resolve variables, otherwise the variables
	// will not exist
	body = dynblock.Expand(body, eCtx)
	if diags := gohcl.DecodeBody(body, eCtx, val); diags.HasErrors() {
		return fmt.Errorf(`failed to decode body using type "%s": %s`, reflect.TypeOf(val).String(), diags.Error())
	}

	return nil
}

func newEvalContext(inputs cty.Value) *hcl.EvalContext {
	return &hcl.EvalContext{
		Variables: map[string]cty.Value{
			"self": inputs,
		},
		Functions: stdfunctions(),
	}
}
