package store

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/valocode/bubbly/ent/release"
	"github.com/valocode/bubbly/ent/releasepolicyviolation"
	"github.com/valocode/bubbly/ent/repo"
	"github.com/valocode/bubbly/env"
)

func TestStore(t *testing.T) {
	bCtx := env.NewBubblyContext()
	s, err := New(bCtx)
	require.NoError(t, err)
	// Create the "demo" project
	_, projectErr := s.client.Project.Create().SetName("demo").Save(s.ctx)
	require.NoError(t, projectErr)
	{
		pErr := s.PopulateStoreWithPolicies("..")
		require.NoError(t, pErr)
		dErr := s.PopulateStoreWithDummyData()
		require.NoError(t, dErr)
	}
	ctx := context.Background()
	client := s.Client()

	//
	// Get the head of the demo repo
	//
	dbRelease, err := client.Release.Query().
		Where(release.HasHeadOfWith(repo.Name("demo"))).
		WithHeadOf().
		Only(ctx)
	require.NoError(t, err)

	vs, err := client.ReleasePolicyViolation.Query().
		Where(releasepolicyviolation.HasReleaseWith(release.ID(dbRelease.ID))).
		All(ctx)
	require.NoError(t, err)

	t.Log("Release Violations for head of demo:")
	for _, v := range vs {
		t.Log(v.String())
	}

}
