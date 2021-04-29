package store

import (
	"fmt"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/valocode/bubbly/env"
	"github.com/valocode/bubbly/test"

	testData "github.com/valocode/bubbly/store/testdata"
)

func TestRelease(t *testing.T) {
	releaseQuery := `
{
	release(_id: "1") {
		release_stage {
			name
			release_criteria {
				entry_name
				release_item {
					release_entry {
						name
						result
					}
				}
			}
		}
	}
}`
	bCtx := env.NewBubblyContext()
	bCtx.UpdateLogLevel(zerolog.DebugLevel)
	resource := test.RunPostgresDocker(bCtx, t)
	bCtx.StoreConfig.PostgresAddr = fmt.Sprintf("localhost:%s", resource.GetPort("5432/tcp"))
	// bCtx.StoreConfig.PostgresAddr = "localhost:5432"

	// Parse the schema and data blocks
	tables := testData.Tables(t, bCtx, "./testdata/release/schema.hcl")
	data := testData.DataBlocks(t, bCtx, "./testdata/release/data.hcl")

	// Initialize a new bubbly store (connection to postgres)
	s, err := New(bCtx)
	require.NoErrorf(t, err, "failed to initialize store")
	err = s.Apply(DefaultTenantName, tables)
	require.NoErrorf(t, err, "failed to apply schema from tables")

	err = s.Save(DefaultTenantName, data)
	require.NoErrorf(t, err, "failed to save data blocks")

	// Query and get the result
	result, err := s.Query(DefaultTenantName, releaseQuery)
	assert.NoErrorf(t, err, "failed to run release query")
	assert.Empty(t, result.Errors)
	val, ok := result.Data.(map[string]interface{})
	require.True(t, ok)
	releases := val["release"].([]interface{})
	require.Len(t, releases, 1)

	// Quite hacky in golang, but traverse the received GraphQL response and verify
	// that the release criteria has been satisfied (by checking that there is
	// a release entry)
	stages := releases[0].(map[string]interface{})["release_stage"].([]interface{})
	for _, stage := range stages {
		stageMap := stage.(map[string]interface{})
		criterion := stageMap["release_criteria"].([]interface{})
		t.Logf("stage: %s", stageMap["name"])
		for _, criteria := range criterion {
			var entryOK bool
			criteriaMap := criteria.(map[string]interface{})
			entryName := criteriaMap["entry_name"]
			item := criteriaMap["release_item"]
			for _, entry := range item.(map[string]interface{})["release_entry"].([]interface{}) {
				entryMap := entry.(map[string]interface{})
				if entryMap["name"] == entryName && entryMap["result"].(bool) {
					entryOK = true
					break
				}
			}
			assert.Truef(t, entryOK, "entry is not satisfied")
			t.Logf("criteria: %#v", criteria)
		}
	}
}
