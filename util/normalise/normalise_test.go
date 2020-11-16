package normalise

import (
	"testing"

	"github.com/likexian/gokit/assert"
)

// tests normalise.LongDesc
func TestLongDesc(t *testing.T) {
	tcs := []struct {
		desc     string
		input    string
		expected string
	}{
		{
			desc: "basic",
			input: `
Show details of a specific resource or group of resources

Print a detailed description of the selected resources, including related resources such
as events or controllers. You may select all resources of a given
type by providing only the type, or additionally provide a name prefix to describe a single resource by name. 

For example:

	$ bubbly describe TYPE NAME_PREFIX`,
			expected: "Show details of a specific resource or group of resources\n\nPrint a detailed description of the selected resources, including related resources such\nas events or controllers. You may select all resources of a given\ntype by providing only the type, or additionally provide a name prefix to describe a single resource by name. \n\nFor example:\n\n\t$ bubbly describe TYPE NAME_PREFIX",
		},
		{
			desc:     "empty string",
			input:    ``,
			expected: ``,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.desc, func(t *testing.T) {
			sNorm := LongDesc(tc.input)
			assert.Equal(t, sNorm, tc.expected)
		})
	}

}

// tests normalise.Examples
func TestExamples(t *testing.T) {
	tcs := []struct {
		desc     string
		input    string
		expected string
	}{
		{
			desc: "basic",
			input: `
# Describe an extract with name 'default'
bubbly describe extract default`,
			expected: "  # Describe an extract with name 'default'\n  bubbly describe extract default",
		},
		{
			desc:     "empty string",
			input:    ``,
			expected: ``,
		},
		{
			desc: "basic with explicit trim",
			input: `
								# Describe an extract with name 'default'
bubbly describe extract default`,
			expected: "  # Describe an extract with name 'default'\n  bubbly describe extract default",
		},
	}

	for _, tc := range tcs {
		t.Run(tc.desc, func(t *testing.T) {
			sNorm := Examples(tc.input)
			assert.Equal(t, sNorm, tc.expected)
		})
	}

}
