package store

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/valocode/bubbly/env"
	"github.com/valocode/bubbly/test"
)

func TestCreateTenant(t *testing.T) {
	const tenant string = "acme"

	bCtx := env.NewBubblyContext()
	res := test.RunPostgresDocker(bCtx, t)
	bCtx.StoreConfig.PostgresAddr = fmt.Sprintf("localhost:%s", res.GetPort("5432/tcp"))
	bCtx.AuthConfig.MultiTenancy = true

	s, err := New(bCtx)
	require.NoError(t, err)

	err = s.CreateTenant(tenant)
	require.NoError(t, err)

	// Run a dummy query
	result, err := s.Query(tenant, "{ release { name } }")
	require.NoError(t, err)
	assert.Empty(t, result.Errors)
}
