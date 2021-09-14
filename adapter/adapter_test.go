package adapter

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRego(t *testing.T) {
	type test struct {
		adapter    string
		inputFiles []string
	}

	tests := []test{
		{adapter: "testdata/adapters/gosec.rego", inputFiles: []string{"testdata/adapters/gosec.json"}},
		{adapter: "testdata/snyk.rego", inputFiles: []string{"testdata/snyk.json"}},
	}
	for _, tc := range tests {
		name := filepath.Base(tc.adapter)
		t.Run(name, func(t *testing.T) {
			result, err := RunFromFile(
				tc.adapter,
				WithInputFileSlice(tc.inputFiles),
				WithTracing(true),
			)
			require.NoError(t, err)
			for _, trace := range result.Traces {
				fmt.Println(trace)
			}
			t.Logf("result: %#v", result.CodeScan)
		})
	}
}
