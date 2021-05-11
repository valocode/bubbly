package bubbly

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/graphql-go/graphql"
	"github.com/valocode/bubbly/bubbly/builtin"
	"github.com/valocode/bubbly/client"
	"github.com/valocode/bubbly/env"
	"github.com/valocode/bubbly/parser"
)

func GetRelease(bCtx *env.BubblyContext) (*builtin.Release, error) {
	releaseQuery := `
{
	release(name: "%s", version: "%s") {
		name
		version
		project(id: "%s") {
			id
			name
		}
		release_item(filter_on: true) {
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

	// Get the release in the current directory
	var fileParser BubblyFileParser
	err := parser.ParseConfig(bCtx, &fileParser)
	if err != nil {
		return nil, fmt.Errorf("error parsing bubbly configs: %w", err)
	}

	if fileParser.Release == nil {
		return nil, fmt.Errorf("no release definition found")
	}
	release := fileParser.Release

	// Ignore the output of Data, but this validates the release and makes sure
	// the default values are set
	if _, err := release.Data(); err != nil {
		return nil, fmt.Errorf("unable to process release definition: %w", err)
	}
	// Insert the necessary values into the GraphQL query
	releaseQuery = fmt.Sprintf(releaseQuery,
		release.Name, release.Version, release.Project,
	)

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

	if len(releases.Release) == 0 {
		return nil, fmt.Errorf("release does not exist")
	}

	return &releases.Release[0], nil
}
