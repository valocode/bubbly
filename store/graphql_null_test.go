package store

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/valocode/bubbly/env"
	"github.com/valocode/bubbly/test"

	testData "github.com/valocode/bubbly/store/testdata"
)

// TestNullFilters ...
func TestNullFilters(t *testing.T) {
	bCtx := env.NewBubblyContext()
	resource := test.RunPostgresDocker(bCtx, t)
	bCtx.StoreConfig.PostgresAddr = fmt.Sprintf("localhost:%s", resource.GetPort("5432/tcp"))
	// bCtx.StoreConfig.PostgresAddr = "localhost:5432"

	// Parse the schema and data blocks
	tables := testData.Tables(t, bCtx, "./testdata/null/tables.hcl")
	data := testData.DataBlocks(t, bCtx, "./testdata/null/data.hcl")
	// Initialize a new bubbly store (connection to postgres)
	s, err := New(bCtx)
	require.NoErrorf(t, err, "failed to initialize store")
	err = s.Apply(DefaultTenantName, tables, true)
	require.NoErrorf(t, err, "failed to apply schema from tables")
	err = s.Save(DefaultTenantName, data)
	require.NoErrorf(t, err, "failed to save data")

	tests := []struct {
		name  string
		query string
		exp   map[string]interface{}
	}{
		{
			name:  "test not_null",
			query: "{ t1 { t2 { t3(not_null: true) { f1 } } } }",
			exp: map[string]interface{}{
				"t1": []interface{}{map[string]interface{}{
					"t2": []interface{}{map[string]interface{}{
						"t3": []interface{}{map[string]interface{}{
							"f1": "abc",
						}},
					}},
				}},
			},
		},
		{
			name:  "test is_null",
			query: "{ t1 { t2(filter_is_null: t3) { f1 } } }",
			exp: map[string]interface{}{
				"t1": []interface{}{map[string]interface{}{
					"t2": []interface{}{map[string]interface{}{
						"f1": "def",
					}},
				}},
			},
		},
		{
			name:  "test is_null with list",
			query: "{ t1 { t2(filter_is_null: [t3]) { f1 } } }",
			exp: map[string]interface{}{
				"t1": []interface{}{map[string]interface{}{
					"t2": []interface{}{map[string]interface{}{
						"f1": "def",
					}},
				}},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := s.Query(DefaultTenantName, tt.query)
			require.NoError(t, err)
			require.Empty(t, result.Errors)
			assert.Equal(t, tt.exp, result.Data)
		})
	}
}
