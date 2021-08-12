package client

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/valocode/bubbly/env"
	"github.com/valocode/bubbly/store/api"
)

func TestClient(t *testing.T) {
	{
		err := CreateRelease(env.NewBubblyContext(), &api.ReleaseCreateRequest{})
		require.NoError(t, err)
	}
	{
		err := SaveCodeScan(env.NewBubblyContext(), &api.CodeScanRequest{})
		require.NoError(t, err)
	}
}
