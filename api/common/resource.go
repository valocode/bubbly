package common

import (
	"encoding/json"
	"fmt"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/zclconf/go-cty/cty"

	"github.com/valocode/bubbly/api/core"
	"github.com/valocode/bubbly/client"
	"github.com/valocode/bubbly/env"
	"github.com/valocode/bubbly/events"
	"github.com/valocode/bubbly/parser"
)

// RunResource takes a given resource id as string and ResourceContext, with
// the input values as cty.Value and runs the resource referenced by the id
func RunResource(bCtx *env.BubblyContext, ctx *core.ResourceContext, id string, inputs cty.Value) (core.Resource, core.ResourceOutput) {
	resource, err := GetResource(bCtx, ctx, id)
	if err != nil {
		return nil, core.ResourceOutput{
			ID:     id,
			Status: events.ResourceApplyFailure,
			Error:  err,
		}
	}
	runCtx := core.NewResourceContext(inputs, ctx.NewResource, ctx.Auth)
	// TODO: handle this better... it get's copied around in multiple places...
	runCtx.DataCtx = ctx.DataCtx
	return resource, resource.Apply(bCtx, runCtx)
}

// GetResource is used by resources that reference other resources (such as
// pipelines) and returns the referenced resource or an error.
// The bubbly client is used to access fetch the resource either via the REST
// API (external) or via NATS (internal, TODO)
func GetResource(bCtx *env.BubblyContext, ctx *core.ResourceContext, resID string) (core.Resource, error) {
	client, err := client.New(bCtx)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize the bubbly client: %w", err)
	}
	defer client.Close()

	resByte, err := client.GetResource(bCtx, ctx.Auth, resID)
	if err != nil {
		return nil, fmt.Errorf("failed to get resource using bubbly client: %w", err)
	}
	var resBlock core.ResourceBlock
	err = json.Unmarshal(resByte, &resBlock)
	if err != nil {
		return nil, fmt.Errorf(`failed to unmarshal resource "%s": %w`, resID, err)
	}
	resource, err := ctx.NewResource(&resBlock)
	if err != nil {
		return nil, fmt.Errorf(`failed to create new resource: %w`, err)
	}

	return resource, nil
}

// DecodeBodyWithInputs decodes a body that is expected to have input
// declarations and validates that all the inputs are provided in the
// ResourceContext. It then updates the inputs in the ResourceContext with any
// default values from the input declarations.
func DecodeBodyWithInputs(bCtx *env.BubblyContext, body hcl.Body, val interface{}, ctx *core.ResourceContext) error {
	retInputs, err := ValidateResourceInputs(bCtx, body, ctx.Inputs)
	if err != nil {
		return fmt.Errorf("resource input validation failed: %w", err)
	}
	// override the current resource context inputs
	ctx.Inputs = retInputs

	if err := parser.DecodeExpandBody(body, val, ctx.Inputs); err != nil {
		return fmt.Errorf("failed to decode resource: %w", err)
	}
	return nil
}

// DecodeBody is the same as DecodeBodyWithInputs but it does not validate any
// inputs before decoding, which means it is simply a wrapper around the parser
// but is here for readability... Could also be removed.
func DecodeBody(bCtx *env.BubblyContext, body hcl.Body, val interface{}, ctx *core.ResourceContext) error {
	if err := parser.DecodeExpandBody(body, val, ctx.Inputs); err != nil {
		return fmt.Errorf("failed to decode resource: %w", err)
	}
	return nil
}

// ValidateResourceInputs takes the body of a resource and the given inputs and
// validates whether all the provided inputs have been given, and returns the
// inputs with any default values from the input declaration added.
// If all inputs are given where no defaults are given, no error is returned.
// If there are missing inputs and no default values, an error is returned
// suggesting which inputs are missing.
func ValidateResourceInputs(bCtx *env.BubblyContext, body hcl.Body, inputs cty.Value) (cty.Value, error) {
	var inputDeclsWrap core.InputDeclarationHCLWrapper
	if diags := gohcl.DecodeBody(body, nil, &inputDeclsWrap); diags.HasErrors() {
		return cty.NilVal, fmt.Errorf("failed to get input declarations: %v", diags.Errs())
	}
	return compareInputsWithDecls(inputDeclsWrap.InputDeclarations, inputs)
}

func compareInputsWithDecls(decls core.InputDeclarations, inputs cty.Value) (cty.Value, error) {
	if !inputs.Type().IsObjectType() {
		return cty.NilVal, fmt.Errorf(`inputs should be of cty.Object type, not "%s"`, inputs.Type().GoString())
	}
	// store the resulting inputs to use, which includes the default values
	// where necessary
	retInputs := make(map[string]cty.Value)

	// init an empty object in case there are no inputs
	inputVals := cty.EmptyObjectVal
	if inputs.Type().HasAttribute("input") {
		inputVals = inputs.GetAttr("input")
	}

	// check the format is correct
	if !inputVals.CanIterateElements() {
		return cty.NilVal, fmt.Errorf(
			"value of inputs to resource is invalid. Should be an %s not %s",
			cty.EmptyObject.GoString(), inputVals.Type().GoString(),
		)
	}

	inputMap := inputVals.AsValueMap()
	for _, decl := range decls {
		val, exists := inputMap[decl.Name]
		if exists {
			// delete the key from the map so that we can check leftovers later
			delete(inputMap, decl.Name)
			retInputs[decl.Name] = val
		} else {
			// if the input was not provided and no default is given, then we
			// have an error
			if decl.Default.IsNull() {
				return cty.NilVal, fmt.Errorf(`input "%s" not provided and does not have a default value`, decl.Name)
			}
			// else, add the default value to the return input values
			retInputs[decl.Name] = decl.Default
		}
	}

	// check if there are any given inputs that should not be given
	if len(inputMap) > 0 {
		var errorInputs []string
		for input := range inputMap {
			errorInputs = append(errorInputs, input)
		}
		return cty.NilVal, fmt.Errorf("%d unexpected inputs provided to resource: %v", len(errorInputs), errorInputs)
	}

	return cty.ObjectVal(map[string]cty.Value{
		"input": cty.ObjectVal(retInputs),
	}), nil
}
