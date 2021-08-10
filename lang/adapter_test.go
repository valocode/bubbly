package lang

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAdapter(t *testing.T) {
	adapt, err := NewAdapterFromFile("testdata/adapters/gosec.adapt.hcl")
	require.NoError(t, err)

	opts := AdapterOptions{
		Filename: "testdata/adapters/gosec.json",
	}
	result, err := adapt.Run(opts)
	require.NoError(t, err)

	t.Logf("result: %#v", result)
}

// func TestAdapter(t *testing.T) {
// 	adapt, err := ParseAdapter("../bubbly/gosec.adapt.hcl")
// 	require.NoError(t, err)

// 	inputs := map[string]cty.Value{
// 		"json": cty.StringVal("../gosec.json"),
// 	}
// 	_, err = adapt.Run(inputs)
// 	require.NoError(t, err)
// }