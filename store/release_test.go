package store

import (
	"fmt"
	"testing"

	"github.com/davecgh/go-spew/spew"
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

	spew.Dump(releases)
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

func TestDataConstraints(t *testing.T) {
	bCtx := env.NewBubblyContext()

	tests := []struct {
		name   string
		schema string
		data   string
		query  string
		want   map[string]interface{}
	}{
		{
			name:   "create/update data with missing unique field",
			schema: "schema1.hcl",
			data:   "data1.hcl",
			query:  "{ t1 { f1 } }",
			want: map[string]interface{}{
				"t1": []interface{}{map[string]interface{}{
					"f1": "f1val",
				}},
			},
		},
		{
			name:   "create/update data with unique join",
			schema: "schema2.hcl",
			data:   "data2.hcl",
			query:  "{ t2 { t1 { f1 } f1 } }",
			want: map[string]interface{}{
				"t2": []interface{}{
					map[string]interface{}{
						"t1": map[string]interface{}{"f1": "f1val"},
						"f1": "",
					},
					map[string]interface{}{
						"t1": map[string]interface{}{"f1": nil},
						"f1": "f1val",
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resource := test.RunPostgresDocker(bCtx, t)
			bCtx.StoreConfig.PostgresAddr = fmt.Sprintf("localhost:%s", resource.GetPort("5432/tcp"))
			// bCtx.StoreConfig.PostgresAddr = "localhost:5432"

			// Parse the schema and data blocks
			tables := testData.Tables(t, bCtx, "./testdata/release/"+tt.schema)
			data := testData.DataBlocks(t, bCtx, "./testdata/release/"+tt.data)

			// Initialize a new bubbly store (connection to postgres)
			s, err := New(bCtx)
			require.NoErrorf(t, err, "failed to initialize store")
			err = s.Apply(DefaultTenantName, tables)
			require.NoErrorf(t, err, "failed to apply schema from tables")
			err = s.Save(DefaultTenantName, data)
			require.NoErrorf(t, err, "failed to save data")
			result, err := s.Query(DefaultTenantName, tt.query)
			require.NoError(t, err)
			require.Empty(t, result.Errors)
			t.Logf("%s:\n%#v", tt.name, result.Data)
			assert.Equal(t, tt.want, result.Data)
		})
	}

}
