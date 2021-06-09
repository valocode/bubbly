package bubbly

import (
	"errors"
	"fmt"

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

	var resWrap builtin.Resource_Wrap
	if err := c.QueryType(bCtx, nil, query, &resWrap); err != nil {
		return nil, err
	}

	if len(resWrap.Resource) == 0 {
		return nil, ErrNoResourcesFound
	}
	return resWrap.Resource, nil
}
