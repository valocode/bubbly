package adapter

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRego(t *testing.T) {
	_, err := RunFromFile("./testdata/adapters/gosec.rego")
	require.NoError(t, err)
}
