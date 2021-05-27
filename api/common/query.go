package common

import (
	"encoding/json"
	"fmt"

	"github.com/graphql-go/graphql"
	"github.com/valocode/bubbly/api/core"
	"github.com/valocode/bubbly/client"
	"github.com/valocode/bubbly/env"
	"github.com/zclconf/go-cty/cty"
	ctyjson "github.com/zclconf/go-cty/cty/json"
)

// QueryToCtyValue executes the criteria query and creates a cty.Value containing
// the results, which can then be used to evaluate the criteria conditions
func QueryToCtyValue(bCtx *env.BubblyContext, ctx *core.ResourceContext, query string) (cty.Value, error) {
	client, err := client.New(bCtx)
	if err != nil {
		return cty.NilVal, fmt.Errorf("error creating bubbly client: %w", err)
	}
	defer client.Close()

	bytes, err := client.Query(bCtx, ctx.Auth, query)
	if err != nil {
		return cty.NilVal, fmt.Errorf("error executing query: %w", err)
	}

	var result graphql.Result
	if err := json.Unmarshal(bytes, &result); err != nil {
		return cty.NilVal, fmt.Errorf("error unmarshalling query result: %w", err)
	}
	if result.HasErrors() {
		return cty.NilVal, fmt.Errorf("received errors from query: %v", result.Errors)
	}

	// Operation: we are given a result from the GraphQL query, which comes as
	// an interface{} but is a map[string]interface{}. The structure of the map
	// and slices therein, depends on the GraphQL query and also the bubbly
	// schema that has been applied.
	// This data needs to be converted into a cty.Value so that the conditions
	// for this criteria can be evaluated.
	// Right now the easiest way seems to be to get the implied type using JSON,
	// but this has limitations/implications... such as missing values in JSON
	// may lead to the wrong type. Probably we need to write some of our own
	// logic here in the future, or have the ability to create a cty.Type from
	// the GraphQL query and Bubbly schema. Anyway, until such an issue arises...
	dBytes, err := json.Marshal(result.Data)
	if err != nil {
		return cty.NilVal, fmt.Errorf("error marshalling query result data: %w", err)
	}
	dType, err := ctyjson.ImpliedType(dBytes)
	if err != nil {
		return cty.NilVal, fmt.Errorf("could not imply type from query result: %w", err)
	}
	queryVal, err := ctyjson.Unmarshal(dBytes, dType)
	if err != nil {
		return cty.NilVal, fmt.Errorf("could not unmarshal result data into cty value: %w", err)
	}
	return queryVal, nil
}
