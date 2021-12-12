package server

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/valocode/bubbly/config"
	"github.com/valocode/bubbly/ent"
	"github.com/valocode/bubbly/store/api"
)

func TestRelease(t *testing.T) {
	s, err := New(config.NewBubblyConfig())
	require.NoError(t, err)

	{
		project := "default"
		err := testPostRequest(t, s.e, s.postRelease, "/releases", &api.ReleaseCreateRequest{
			Project:    &project,
			Repository: ent.NewRepositoryModelCreate().SetName("test-repo"),
			Commit:     ent.NewGitCommitModelCreate().SetHash("123").SetBranch("main").SetTime(time.Now()),
			Release:    ent.NewReleaseModelCreate().SetName("test-repo").SetVersion("123"),
		})
		require.NoError(t, err)
	}
	{
		project := "default"
		err := testPostRequest(t, s.e, s.postRelease, "/releases", &api.ReleaseCreateRequest{
			Project:    &project,
			Repository: ent.NewRepositoryModelCreate().SetName("another-repo"),
			Commit:     ent.NewGitCommitModelCreate().SetHash("456").SetBranch("main").SetTime(time.Now()),
			Release:    ent.NewReleaseModelCreate().SetName("another-repo").SetVersion("456"),
		})
		require.NoError(t, err)
	}
	t.Run("get all releases", func(t *testing.T) {
		rec, err := testGetRequest(t, s.e, s.getReleases, "/releases", nil)
		assert.NoError(t, err)
		var resp api.ReleaseGetResponse
		json.NewDecoder(rec.Body).Decode(&resp)
		assert.Len(t, resp.Releases, 2)
	})
	t.Run("get project releases", func(t *testing.T) {
		rec, err := testGetRequest(t, s.e, s.getReleases, "/releases?projects=default", nil)
		assert.NoError(t, err)
		var resp api.ReleaseGetResponse
		json.NewDecoder(rec.Body).Decode(&resp)
		assert.Len(t, resp.Releases, 2)
	})
	t.Run("get project releases - unknown project", func(t *testing.T) {
		rec, err := testGetRequest(t, s.e, s.getReleases, "/releases?projects=unknown", nil)
		assert.NoError(t, err)
		var resp api.ReleaseGetResponse
		json.NewDecoder(rec.Body).Decode(&resp)
		assert.Len(t, resp.Releases, 0)
	})
}
