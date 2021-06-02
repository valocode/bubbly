package builtin

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/valocode/bubbly/api/core"
	"github.com/zclconf/go-cty/cty"
)

func TestGen(t *testing.T) {
	tests := []struct {
		name  string
		input interface{}
		exp   core.DataBlocks
	}{
		{
			name: "simple release",
			input: Release{
				Name:    "release",
				Version: "version",
				Project: &Project{
					Name:          "project",
					DBlock_Policy: core.ReferencePolicy,
				},
				ReleaseStage: []ReleaseStage{
					{
						Name: "stage",
						ReleaseCriteria: []ReleaseCriteria{
							{
								EntryName:    "entry",
								DBlock_Joins: []string{"release"},
							},
						},
					},
				},
			},
			exp: core.DataBlocks{
				core.Data{
					TableName:     "project",
					Fields:        &core.DataFields{Values: map[string]cty.Value{"name": cty.StringVal("project")}},
					Policy:        "reference",
					IgnoreNesting: true,
				},
				core.Data{
					TableName: "release",
					Fields: &core.DataFields{Values: map[string]cty.Value{
						"name":    cty.StringVal("release"),
						"version": cty.StringVal("version"),
					}},
					Joins: []string{"project"},
					Data: core.DataBlocks{
						core.Data{
							TableName: "release_stage",
							Fields:    &core.DataFields{Values: map[string]cty.Value{"name": cty.StringVal("stage")}},
							Data: core.DataBlocks{
								core.Data{
									TableName: "release_criteria",
									Fields:    &core.DataFields{Values: map[string]cty.Value{"entry_name": cty.StringVal("entry")}},
									Joins:     []string{"release"},
								},
							},
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dblocks := ToDataBlocks(tt.input)
			assert.Equal(t, tt.exp, dblocks)
		})
	}
}
