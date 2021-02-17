package util

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

// tests util.UsageErrorf
func TestUsageErrorf(t *testing.T) {
	tcs := []struct {
		desc     string
		expected string
	}{
		{
			desc:     "basic",
			expected: "Unexpected args: [example wrong arguments]\nSee 'parent -h' for help and examples",
		},
	}

	for _, tc := range tcs {
		t.Run(tc.desc, func(t *testing.T) {
			cmd := &cobra.Command{Use: "parent"}
			actual := UsageErrorf(cmd, "Unexpected args: %v", []string{"example", "wrong", "arguments"})

			assert.Equal(t, tc.expected, actual.Error())

		})
	}
}
