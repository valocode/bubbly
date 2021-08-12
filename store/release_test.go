package store

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/valocode/bubbly/ent/model"
	"github.com/valocode/bubbly/env"
	"github.com/valocode/bubbly/store/api"

	// Required since we use schema hooks!
	_ "github.com/valocode/bubbly/ent/runtime"
)

func TestCreateRelease(t *testing.T) {
	store, err := New(env.NewBubblyContext())
	require.NoError(t, err)

	{
		_, err := store.CreateRelease(&api.ReleaseCreateRequest{
			// Project: model.NewProjectModel().SetName("project"),
			Repo:    model.NewRepoModel().SetName("repo"),
			Commit:  model.NewGitCommitModel().SetBranch("main").SetHash("12345").SetTime(time.Now()),
			Release: model.NewReleaseModel(),
		})
		require.NoError(t, err)
	}
}
