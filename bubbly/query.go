package bubbly

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/graphql-go/graphql"
	"github.com/hashicorp/go-multierror"

	"github.com/valocode/bubbly/bubbly/builtin"
	"github.com/valocode/bubbly/client"
	"github.com/valocode/bubbly/env"
)

var ErrNoResourcesFound = errors.New("no resources found")

// QueryResources takes a graphql query string, queries the bubbly store and
// returns resource and event information matching the query
func QueryResources(bCtx *env.BubblyContext, query string) ([]builtin.Resource, error) {
	c, err := client.New(bCtx)
	if err != nil {
		return nil, fmt.Errorf("failed to create bubbly client: %w", err)
	}

	res, err := c.Query(bCtx, nil, query)

	if err != nil {
		return nil, fmt.Errorf("failed to query: %w", err)
	}
	var (
		resWrap builtin.Resource_Wrap
		result  graphql.Result
	)
	// Assign the resource pointer to data so that we unmarshall directly
	// into resWrap
	result.Data = &resWrap
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

	if len(resWrap.Resource) == 0 {
		return nil, ErrNoResourcesFound
	}
	return resWrap.Resource, nil
}
