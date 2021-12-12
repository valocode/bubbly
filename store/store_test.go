package store

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/valocode/bubbly/config"
	"github.com/valocode/bubbly/ent/release"
	"github.com/valocode/bubbly/ent/releasepolicyviolation"
	"github.com/valocode/bubbly/ent/repository"
)

func TestStore(t *testing.T) {
	bCtx := config.NewBubblyConfig()
	s, err := New(bCtx)
	require.NoError(t, err)
	h, err := NewHandler(WithStore(s))
	require.NoError(t, err)
	// Create the "demo" project
	_, projectErr := h.client.Project.Create().SetName("demo").SetOwnerID(h.orgID).Save(h.ctx)
	require.NoError(t, projectErr)
	{
		pErr := h.PopulateStoreWithPolicies()
		require.NoError(t, pErr)
		dErr := h.PopulateStoreWithDummyData()
		require.NoError(t, dErr)
	}
	ctx := context.Background()
	client := h.Client()

	//
	// Get the head of the demo repo
	//
	dbRelease, err := client.Release.Query().
		Where(release.HasHeadOfWith(repository.Name("demo"))).
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
