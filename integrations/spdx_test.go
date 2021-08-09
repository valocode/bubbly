package integrations

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSpdx(t *testing.T) {
	list, err := FetchSPDXLicenses()
	require.NoError(t, err)
	t.Logf("list: %#v", list)
}
