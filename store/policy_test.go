package store

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/valocode/bubbly/config"
	"github.com/valocode/bubbly/ent"
	"github.com/valocode/bubbly/store/api"
	"github.com/valocode/bubbly/test"
)

func TestPolicyAffects(t *testing.T) {
	s, err := New(config.NewBubblyConfig())
	require.NoError(t, err)

	h, err := NewHandler(WithStore(s))
	require.NoError(t, err)

	dbProject, projErr := h.CreateProject(&api.ProjectCreateRequest{
		Project: ent.NewProjectModelCreate().SetName("test"),
	})
	require.NoError(t, projErr)

	dbRelease, relErr := h.CreateRelease(&api.ReleaseCreateRequest{
		Project:    dbProject.Name,
		Repository: ent.NewRepositoryModelCreate().SetName("test"),
		Commit:     ent.NewGitCommitModelCreate().SetHash("123").SetBranch("main").SetTime(time.Now()),
		Release:    ent.NewReleaseModelCreate().SetName("test").SetVersion("test"),
	})
	require.NoError(t, relErr)

	policies, err := test.ParsePolicies()
	require.NoError(t, err)
	for _, p := range policies {
		_, policyErr := h.CreateReleasePolicy(&api.ReleasePolicyCreateRequest{
			Policy: &api.ReleasePolicyCreate{
				ReleasePolicyModelCreate: *ent.NewReleasePolicyModelCreate().SetName(*p.Policy.Name).SetModule(*p.Policy.Module),
				Affects: &api.ReleasePolicyAffectsSet{
					Projects: []string{"test"},
					Repos:    []string{"test"},
				},
			},
		})
		require.NoError(t, policyErr)
	}

	dbPolicies, err := h.policiesForRelease(h.client, dbRelease.ID)
	require.NoError(t, err)

	assert.Lenf(t, dbPolicies, len(policies), "the number of policies for the release should be the same as the ones we just assigned")
}
