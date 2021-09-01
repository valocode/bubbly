package store

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/valocode/bubbly/ent"
	"github.com/valocode/bubbly/env"
	"github.com/valocode/bubbly/store/api"
	"github.com/valocode/bubbly/test"
)

func TestPolicyAffects(t *testing.T) {
	s, err := New(env.NewBubblyContext())
	require.NoError(t, err)

	h, err := NewHandler(WithStore(s))
	require.NoError(t, err)

	dbRelease, relErr := h.CreateRelease(&api.ReleaseCreateRequest{
		Project: ent.NewProjectModelCreate().SetName("test"),
		Repo:    ent.NewRepoModelCreate().SetName("test"),
		Commit:  ent.NewGitCommitModelCreate().SetHash("123").SetBranch("main").SetTime(time.Now()),
		Release: ent.NewReleaseModelCreate().SetName("test").SetVersion("test"),
	})
	require.NoError(t, relErr)

	policies, err := test.ParsePolicies("..")
	require.NoError(t, err)
	for _, p := range policies {
		_, policyErr := h.SaveReleasePolicy(&api.ReleasePolicySaveRequest{
			Policy: ent.NewReleasePolicyModelCreate().SetName(*p.Policy.Name).SetModule(*p.Policy.Module),
			Affects: &api.ReleasePolicyAffects{
				Projects: []string{"test"},
				Repos:    []string{"test"},
			},
		})
		require.NoError(t, policyErr)
	}

	dbPolicies, err := h.policiesForRelease(h.client, dbRelease.ID)
	require.NoError(t, err)

	assert.Lenf(t, dbPolicies, len(policies), "the number of policies for the release should be the same as the ones we just assigned")
}
