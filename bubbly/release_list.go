package bubbly

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/graphql-go/graphql"
	"github.com/valocode/bubbly/bubbly/builtin"
	"github.com/valocode/bubbly/client"
	"github.com/valocode/bubbly/env"
)

func ListReleases(bCtx *env.BubblyContext) (*builtin.Release_Wrap, error) {
	releaseQuery := `
{
	release {
		name
		version
		project {
			id
			name
		}
		release_input(filter_on: true) {
			type
			commit {
				repo {
					name
				}
			}
		}
		release_stage(filter_on: true) {
			name
			release_criteria(filter_on: true) {
				entry_name
			}
		}
	}
}
	`

	var releases builtin.Release_Wrap

	client, err := client.New(bCtx)
	if err != nil {
		return nil, fmt.Errorf("error creating bubbly client: %w", err)
	}
	bytes, err := client.Query(bCtx, nil, releaseQuery)
	if err != nil {
		return nil, fmt.Errorf("error making GraphQL query: %w", err)
	}

	var results graphql.Result
	results.Data = &releases
	if err := json.Unmarshal(bytes, &results); err != nil {
		return nil, fmt.Errorf("error unmarshalling GraphQL results: %w", err)
	}
	if results.HasErrors() {
		var msgs []string
		for _, err := range results.Errors {
			msgs = append(msgs, err.Message)
		}
		return nil, fmt.Errorf("GraphQL query returned errors:\n%s", strings.Join(msgs, "\n\n"))
	}

	return &releases, nil
}
