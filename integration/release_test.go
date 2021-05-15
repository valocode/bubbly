package integration

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/valocode/bubbly/bubbly"
	"github.com/valocode/bubbly/env"
)

func TestRelease(t *testing.T) {

	tests := []struct {
		name    string
		release string
	}{
		{
			name:    "release git commit",
			release: "./testdata/release/release-git-commit.bubbly",
		},
		{
			name:    "release git tag",
			release: "./testdata/release/release-git-tag.bubbly",
		},
	}
	bCtx := env.NewBubblyContext()
	// Apply all resources in the release directory
	err := bubbly.Apply(bCtx, "./testdata/release/resources")
	require.NoError(t, err)

	t.Run("list releases", func(t *testing.T) {
		// There should be no releases right now
		rel, err := bubbly.ListReleases(bCtx)
		require.NoError(t, err)
		t.Logf("Found %d releases", len(rel.Release))
	})
	for _, tt := range tests {
		t.Logf("TEST: %#v", tt)
		t.Run("create release - "+tt.name, func(t *testing.T) {
			rel, err := bubbly.CreateRelease(bCtx, tt.release)
			// TODO: currently no way to know if error is acceptable or not.
			// E.g. if the release already exists, it's ok. But if the release
			// creation failed for some other reason, we want to know!
			if err == nil {
				t.Logf("Created release: %#v", rel.Name)
			}
			if err != nil {
				t.Logf("Release creation failed: %s", err.Error())
			}
		})
		t.Run("get release - "+tt.name, func(t *testing.T) {
			rel, err := bubbly.GetRelease(bCtx, tt.release)
			// We may get an error, but it should be ErrReleaseNotExist
			if err == nil {
				t.Logf("Current release: %s", rel.Name)
			}
			if err != nil {
				assert.True(t, errors.Is(err, bubbly.ErrReleaseNotExist))
			}
		})
		// TODO: this is not idempotent... test requires the release to be created
		rel, err := bubbly.GetRelease(bCtx, tt.release)
		require.NoError(t, err)
		assert.NotEmpty(t, rel.ReleaseStage)
		for _, stage := range rel.ReleaseStage {
			assert.NotEmpty(t, stage.ReleaseCriteria)
			for _, criteria := range stage.ReleaseCriteria {
				t.Run("eval criteria - "+tt.name+" - "+criteria.EntryName, func(t *testing.T) {
					_, err := bubbly.EvalReleaseCriteria(bCtx, tt.release, criteria.EntryName)
					assert.NoError(t, err)
				})
			}
		}
	}

	t.Run("list releases", func(t *testing.T) {
		// There should be no releases right now
		rel, err := bubbly.ListReleases(bCtx)
		require.NoError(t, err)
		t.Logf("Found %d releases", len(rel.Release))
	})
}
