package bubbly

import (
	"fmt"

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
			name
		}
		release_item(not_null: true) {
			type
			commit {
				repo {
					name
				}
			}
		}
		release_stage(not_null: true) {
			name
			release_criteria(not_null: true) {
				release_entry {
					result
				}
			}
		}
	}
}
	`

	client, err := client.New(bCtx)
	if err != nil {
		return nil, fmt.Errorf("error creating bubbly client: %w", err)
	}
	var releases builtin.Release_Wrap
	err = client.QueryType(bCtx, nil, releaseQuery, &releases)
	if err != nil {
		return nil, fmt.Errorf("error executing GraphQL query: %w", err)
	}
	return &releases, nil
}
