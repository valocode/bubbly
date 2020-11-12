package v1

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/zclconf/go-cty/cty"
)

// tests InputDefinitions.Value
func TestValue(t *testing.T) {
	tcs := []struct {
		desc     string
		input    InputDefinitions
		expected cty.Value
	}{
		{
			desc: "basic Value",
			input: InputDefinitions{
				&InputDefinition{
					Name:  "api_version",
					Value: cty.StringVal("v1"),
				},
				&InputDefinition{
					Name:  "project",
					Value: cty.StringVal("project"),
				},
			},
			expected: cty.ObjectVal(map[string]cty.Value{
				"input": cty.ObjectVal(map[string]cty.Value{
					"api_version": cty.StringVal("v1"),
					"project":     cty.StringVal("project"),
				}),
			}),
		},
	}

	for _, tc := range tcs {
		t.Run(tc.desc, func(t *testing.T) {
			actual := tc.input.Value()

			assert.NotNil(t, actual)

			assert.Equal(t, tc.expected, actual)
		})
	}
}
