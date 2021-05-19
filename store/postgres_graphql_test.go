package store

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestScanTableColumns tests the unpacking of SQL row results (flat list) into
// a nested structure following the hierarchy of the GraphQL query that created
// the results.
// The hierarchy is described by tableColumns. If you just spend a second looking
// at a tableColumn it maps quite nicely to the structure of GraphQL, so there
// is no "magic" here. The tricky part is creating the ending golang
// map[string]interface{} but that's handled by the function under test...
func TestScanTableColumns(t *testing.T) {

	tests := []struct {
		name   string
		tc     tableColumns
		values [][]interface{}
		exp    map[string]interface{}
	}{
		{
			name: "simple nested no scalar",
			tc: tableColumns{
				table:   "a",
				columns: []string{"a1", "a2", "a3"},
				children: []*tableColumns{
					{
						table:   "b",
						columns: []string{"b1", "b2"},
					},
				},
			},
			values: [][]interface{}{
				{
					"a1", "a2", "a3", "b1", "b2",
				},
			},
			exp: map[string]interface{}{
				"a": []map[string]interface{}{
					{
						"a1": "a1", "a2": "a2", "a3": "a3",
						"b": []map[string]interface{}{
							{
								"b1": "b1", "b2": "b2",
							},
						},
					},
				},
			},
		},
		{
			name: "simple nested scalar",
			tc: tableColumns{
				table:   "a",
				columns: []string{"a1", "a2", "a3"},
				children: []*tableColumns{
					{
						table:   "b",
						columns: []string{"b1", "b2"},
						scalar:  true,
					},
				},
			},
			values: [][]interface{}{
				{
					"a1", "a2", "a3", "b1", "b2",
				},
			},
			exp: map[string]interface{}{
				"a": []map[string]interface{}{
					{
						"a1": "a1", "a2": "a2", "a3": "a3",
						"b": map[string]interface{}{
							"b1": "b1", "b2": "b2",
						},
					},
				},
			},
		},
		{
			name: "nil node with children",
			tc: tableColumns{
				table:   "a",
				columns: []string{"a1", "a2", "a3"},
				children: []*tableColumns{
					{
						table:   "b",
						columns: []string{"b1", "b2"},
						children: []*tableColumns{
							{
								table:   "c",
								columns: []string{"c1", "c2"},
							},
						},
					},
				},
			},
			values: [][]interface{}{
				{
					"a1", "a2", "a3", "b1", "b2", "c1", "c2",
				},
				{
					"a1", "a2", "a3", nil, nil, "c1", nil,
				},
			},
			exp: map[string]interface{}{
				"a": []map[string]interface{}{
					{
						"a1": "a1", "a2": "a2", "a3": "a3",
						"b": []map[string]interface{}{
							{
								"b1": "b1", "b2": "b2",
								"c": []map[string]interface{}{
									{
										"c1": "c1", "c2": "c2",
									},
								},
							},
						},
					},
					{
						"a1": "a1", "a2": "a2", "a3": "a3",
						"b": []map[string]interface{}{
							{
								"b1": nil, "b2": nil,
								"c": []map[string]interface{}{
									{
										"c1": "c1", "c2": nil,
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "nil node with children same id",
			tc: tableColumns{
				table: "a",
				// NOTE: a1 is now tableIDField so that two records with
				// the same ID are the "same"
				columns: []string{tableIDField, "a2", "a3"},
				children: []*tableColumns{
					{
						table: "b",
						// NOTE: b1 is now tableIDField so that two records with
						// the same ID are the "same"
						columns: []string{tableIDField, "b2"},
						children: []*tableColumns{
							{
								table:   "c",
								columns: []string{"c1", "c2"},
							},
						},
					},
				},
			},
			values: [][]interface{}{
				{
					"a1", "a2", "a3", "b1", "b2", "c1", "c2",
				},
				{
					"a1", "a2", "a3", "b1", nil, "c1", nil,
				},
			},
			exp: map[string]interface{}{
				"a": []map[string]interface{}{
					{
						tableIDField: "a1", "a2": "a2", "a3": "a3",
						"b": []map[string]interface{}{
							{
								tableIDField: "b1", "b2": "b2",
								"c": []map[string]interface{}{
									{
										"c1": "c1", "c2": "c2",
									},
									{
										"c1": "c1", "c2": nil,
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "example release view",
			tc: tableColumns{
				table:   "release",
				columns: []string{tableIDField, "name", "version"},
				children: []*tableColumns{
					{
						table:   "project",
						columns: []string{tableIDField, "name"},
						scalar:  true,
					},
					{
						table:   "release_item",
						columns: []string{tableIDField, "type"},
						children: []*tableColumns{
							{
								table:   "commit",
								columns: []string{tableIDField},
								scalar:  true,
								children: []*tableColumns{
									{
										table:   "repo",
										columns: []string{tableIDField, "name"},
										scalar:  true,
									},
								},
							},
						},
					},
					{
						table:   "release_stage",
						columns: []string{tableIDField, "name"},
						children: []*tableColumns{
							{
								table:   "release_criteria",
								columns: []string{tableIDField, "entry_name"},
								children: []*tableColumns{
									{
										table:   "release_entry",
										columns: []string{tableIDField, "result", "reason"},
									},
								},
							},
						},
					},
				},
			},
			values: [][]interface{}{
				{
					1, "github.com/valocode/bubbly", "bcabc689bbae4f69a04d1be56a22e1f27ae797f2", 1, "bubbly", 1, "git", 1, 1, nil, 2, "Code Analysis", 2, "gosec", nil, nil, nil,
				},
				{
					1, "github.com/valocode/bubbly", "bcabc689bbae4f69a04d1be56a22e1f27ae797f2", 1, "bubbly", 1, "git", 1, 1, nil, 3, "Testing", 4, "integration_test", nil, nil, nil,
				},
				{
					1, "github.com/valocode/bubbly", "bcabc689bbae4f69a04d1be56a22e1f27ae797f2", 1, "bubbly", 1, "git", 1, 1, nil, 1, "Artifact", 1, "artifact", nil, nil, nil,
				},
				{
					1, "github.com/valocode/bubbly", "bcabc689bbae4f69a04d1be56a22e1f27ae797f2", 1, "bubbly", 1, "git", 1, 1, nil, 3, "Testing", 3, "unit_test", 1, false, "no test run found",
				},
			},
			exp: map[string]interface{}{
				"release": []map[string]interface{}{
					{
						tableIDField: 1, "name": "github.com/valocode/bubbly", "version": "bcabc689bbae4f69a04d1be56a22e1f27ae797f2",
						"project": map[string]interface{}{
							tableIDField: 1, "name": "bubbly",
						},
						"release_item": []map[string]interface{}{
							{
								tableIDField: 1, "type": "git",
								"commit": map[string]interface{}{
									tableIDField: 1,
									"repo": map[string]interface{}{
										tableIDField: 1, "name": nil,
									},
								},
							},
						},
						"release_stage": []map[string]interface{}{
							{
								tableIDField: 2, "name": "Code Analysis",
								"release_criteria": []map[string]interface{}{
									{
										tableIDField: 2, "entry_name": "gosec",
										"release_entry": []map[string]interface{}{},
									},
								},
							},
							{
								tableIDField: 3, "name": "Testing",
								"release_criteria": []map[string]interface{}{
									{
										tableIDField: 4, "entry_name": "integration_test",
										"release_entry": []map[string]interface{}{},
									},
									{
										tableIDField: 3, "entry_name": "unit_test",
										"release_entry": []map[string]interface{}{
											{
												tableIDField: 1,
												"result":     false,
												"reason":     "no test run found",
											},
										},
									},
								},
							},
							{
								tableIDField: 1, "name": "Artifact",
								"release_criteria": []map[string]interface{}{
									{
										tableIDField: 1, "entry_name": "artifact",
										"release_entry": []map[string]interface{}{},
									},
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
			var (
				result = make(map[string]interface{})
				index  int
			)
			for _, val := range tt.values {
				// Reset the index each time
				index = 0
				psqlScanTableColumns(result, tt.tc, val, &index)
			}
			assert.Equal(t, tt.exp, result)
		})
	}
}
