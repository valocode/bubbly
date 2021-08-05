package integrations

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFetchCVEs(t *testing.T) {
	resp, err := FetchCVEs(DefaultCVEFetchOptions())
	require.NoError(t, err)
	t.Logf("Received %d CVEs", len(resp.Result.CVEItems))
}
