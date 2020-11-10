package core

import (
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/stretchr/testify/assert"

	"github.com/zclconf/go-cty/cty"
)

// tests ResourceOutput.Output
func TestOutput(t *testing.T) {
	tcs := []struct {
		desc     string
		input    ResourceOutput
		want     cty.Value
		expected cty.Value
	}{
		{
			desc: "basic Output",
			input: ResourceOutput{
				ID:     "example",
				Status: "Success",
				Value: cty.ObjectVal(
					map[string]cty.Value{
						"example": cty.ObjectVal(
							map[string]cty.Value{
								"nested_example": cty.StringVal("value"),
							},
						),
					},
				),
			},
			expected: cty.Value(
				cty.ObjectVal(
					map[string]cty.Value{
						"id":     cty.StringVal("example"),
						"status": cty.StringVal("Success"),
						"value": cty.ObjectVal(map[string]cty.Value{
							"example": cty.ObjectVal(map[string]cty.Value{
								"nested_example": cty.StringVal("value"),
							}),
						}),
					},
				),
			),
		},
	}

	for _, tc := range tcs {
		t.Run(tc.desc, func(t *testing.T) {
			actual := tc.input.Output()

			assert.NotNil(t, actual)

			assert.Equal(t, tc.expected, actual)
		})
	}
}

// tests Local.Reference
func TestReference(t *testing.T) {
	tcs := []struct {
		desc     string
		input    Local
		wantTrav hcl.Traversal
		wantCty  cty.Value
	}{
		{
			desc: "basic local block",
			input: Local{
				Name:  "api_version",
				Value: cty.StringVal("v1"),
			},
			wantTrav: hcl.Traversal{
				hcl.TraverseRoot{Name: "local"},
				hcl.TraverseAttr{Name: "api_version"},
			},
			wantCty: cty.StringVal("v1"),
		},
	}

	for _, tc := range tcs {
		t.Run(tc.desc, func(t *testing.T) {
			actualTrav, actualCty := tc.input.Reference()

			assert.NotNil(t, actualTrav)
			assert.NotNil(t, actualCty)

			assert.Equal(t, tc.wantTrav, actualTrav)
			assert.Equal(t, tc.wantCty, actualCty)
		})
	}
}
