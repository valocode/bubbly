package adapter

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRego(t *testing.T) {
	result, err := RunFromFile(
		"./testdata/adapters/gosec.rego",
		WithInputFiles("./testdata/adapters/gosec.json"),
		// WithTracing(true),
	)
	require.NoError(t, err)
	t.Logf("result: %#v", result.CodeScan)
	for _, trace := range result.Traces {
		fmt.Println(trace)
	}
}
