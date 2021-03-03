package bubbly

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/graphql-go/graphql"
	"github.com/hashicorp/go-multierror"

	"github.com/verifa/bubbly/api/core"
	"github.com/verifa/bubbly/client"
	"github.com/verifa/bubbly/env"
	"github.com/verifa/bubbly/events"
)

// errors
var ErrNoResourcesFound = errors.New("no resources found")

// Resource is a struct representing a bubbly resource extracted from a query
// against the store.
type Resource struct {
	Id     string         // the unique identifier for the resource: kind/name
	Events []events.Event `json:"_event"` // a slice of Events
}

// QueryResources takes a graphql query string, queries the bubbly store and
// returns resource and event information matching the query
func QueryResources(bCtx *env.BubblyContext, query string) ([]Resource, error) {
	c, err := client.NewHTTP(bCtx)
	if err != nil {
		return nil, fmt.Errorf("failed to create bubbly client: %w", err)
	}

	res, err := c.Query(bCtx, query)

	if err != nil {
		return nil, fmt.Errorf("failed to query: %w", err)
	}

	var result graphql.Result

	if err := json.Unmarshal(res, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal get resp: %w", err)
	}

	if result.HasErrors() {
		var graphqlErrors error

		for _, qlError := range result.Errors {
			graphqlErrors = multierror.Append(graphqlErrors, qlError)
		}

		return nil, fmt.Errorf("failed to get resources: %w", graphqlErrors)
	}

	r, err := convertResponseToResources(result.Data)

	if err != nil {
		switch err {
		case ErrNoResourcesFound:
			// return known errors upstream to the cmd package unwrapped
			return nil, err
		default:
			// if unknown, wrap with generic message
			return nil, fmt.Errorf("failed to convert response to resources")
		}
	}

	return r, nil
}

// convertResponseToResources takes the graphql.Response.Data and
// produces a []Resource slice.
// It does this by marshalling the relevant subset of the
// response and then unmarshalling directly into a []Resource struct
func convertResponseToResources(response interface{}) ([]Resource, error) {
	// pull out the []interface{} representing the list of returned resources
	resources := response.(map[string]interface{})[core.ResourceTableName]

	if resources == nil {
		return nil, ErrNoResourcesFound
	}

	resourcesBytes, err := json.Marshal(resources.([]interface{}))

	if err != nil {
		return nil, fmt.Errorf("failed to marshal graphql query response: %w", err)
	}

	var resSlice []Resource

	err = json.Unmarshal(resourcesBytes, &resSlice)

	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal []Resource from graphql query response: %w", err)
	}

	return resSlice, nil
}
