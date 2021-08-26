package store

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/valocode/bubbly/ent/release"
	"github.com/valocode/bubbly/ent/releasepolicyviolation"
	"github.com/valocode/bubbly/env"
)

func TestPolicies(t *testing.T) {
	bCtx := env.NewBubblyContext()
	s, err := New(bCtx)
	require.NoError(t, err)
	{
		pErr := s.PopulateStoreWithPolicies("..")
		require.NoError(t, pErr)
		dErr := s.PopulateStoreWithDummyData()
		require.NoError(t, dErr)
	}
	ctx := context.Background()
	client := s.Client()

	dbRelease, err := client.Release.Query().
		Where(release.HasHeadOf()).
		WithHeadOf().
		Only(ctx)
	require.NoError(t, err)

	_, evalErr := s.EvaluateReleasePolicies(dbRelease.ID)
	require.NoError(t, evalErr)

	vs, err := client.ReleasePolicyViolation.Query().
		Where(releasepolicyviolation.HasReleaseWith(release.ID(dbRelease.ID))).
		All(ctx)
	require.NoError(t, err)

	for _, v := range vs {
		t.Log("Created Violation: ", v.String())
	}

}
