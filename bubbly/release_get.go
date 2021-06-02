package bubbly

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/graphql-go/graphql"
	"github.com/valocode/bubbly/bubbly/builtin"
	"github.com/valocode/bubbly/client"
	"github.com/valocode/bubbly/env"
)

var ErrReleaseNotExist = errors.New("release does not exist")

func GetRelease(bCtx *env.BubblyContext, filename string) (*builtin.Release, error) {
	// A query for a specific release (by name and version).
	// The only tricky part is the order_by the release_entry to only get the
	// latest for the release_criteria... as there may be multiple
	// release_entry for one release_criteria
	releaseQuery := `
{
	release(name: "%s", version: "%s") {
		name
		version
		project(name: "%s") {
			name
		}
		release_item {
			type
			commit {
				repo {
					name
				}
			}
		}
		release_stage {
			name
			release_criteria {
				entry_name
				release_entry(last: 1) {
					result
					reason
				}
			}
		}
	}
}
	`

	release, err := createReleaseSpec(bCtx, filename)
	if err != nil {
		return nil, fmt.Errorf("error creating release spec: %w", err)
	}

	// Validate the release data and
	if err := release.Validate(); err != nil {
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
		return nil, ErrReleaseNotExist
	}

	return &releases.Release[0], nil
}
