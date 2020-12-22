package common

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/verifa/bubbly/api/core"
	"github.com/verifa/bubbly/client"
	"github.com/verifa/bubbly/env"
	"github.com/zclconf/go-cty/cty"
)

// RunResource takes a given resource id as string and ResourceContext, with
// the input values as cty.Value and runs the resource referenced by the id
func RunResource(bCtx *env.BubblyContext, ctx *core.ResourceContext, id string, inputs cty.Value) (core.Resource, core.ResourceOutput) {
	resource, err := GetResource(bCtx, ctx, id)
	if err != nil {
		return nil, core.ResourceOutput{
			Error: err,
		}
	}
	runCtx := core.NewResourceContext(resource.Namespace(), inputs, ctx.NewResource)
	return resource, resource.Apply(bCtx, runCtx)
}

// GetResource is used by resources that reference other resources (such as
// pipelines) and returns the referenced resource or an error.
// The bubbly client is used to access fetch the resource either via the REST
// API (external) or via NATS (internal, TODO)
func GetResource(bCtx *env.BubblyContext, ctx *core.ResourceContext, id string) (core.Resource, error) {
	resID, err := normalizeNamespace(id, ctx.Namespace)
	if err != nil {
		return nil, fmt.Errorf("could not normalize the resource ID: %w", err)
	}
	client, err := client.New(bCtx)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize the bubbly client: %w", err)
	}
	resByte, err := client.GetResource(bCtx, resID)
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

func normalizeNamespace(id string, namespaceCtx string) (string, error) {
	strVals := strings.Split(id, "/")
	switch len(strVals) {
	case 1:
		return "", fmt.Errorf(`resource ID "%s" is missing kind`, id)
	case 2:
		// concatenate the namespace to the string
		return fmt.Sprintf("%s/%s", namespaceCtx, id), nil
	case 3:
		return id, nil
	default:
		return "", fmt.Errorf(`incorrect resource ID format for "%s"`, id)
	}
}
