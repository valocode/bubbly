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

// TestUniqueConstraints takes some tables and data blocks which should all be
// unique, applies those data blocks multiple times and then checks that only
// one instance of the data blocks exists
func TestUniqueConstraints(t *testing.T) {
	bCtx := env.NewBubblyContext()
	resource := test.RunPostgresDocker(bCtx, t)
	bCtx.StoreConfig.PostgresAddr = fmt.Sprintf("localhost:%s", resource.GetPort("5432/tcp"))
	// bCtx.StoreConfig.PostgresAddr = "localhost:5432"

	// Parse the schema and data blocks
	tables := testData.Tables(t, bCtx, "./testdata/unique/tables.hcl")
	data := testData.DataBlocks(t, bCtx, "./testdata/unique/data.hcl")
	// Initialize a new bubbly store (connection to postgres)
	s, err := New(bCtx)
	require.NoErrorf(t, err, "failed to initialize store")
	err = s.Apply(DefaultTenantName, tables, true)
	require.NoErrorf(t, err, "failed to apply schema from tables")
	// Save it more than once to test the unique constraints
	for i := 0; i < 10; i++ {
		err = s.Save(DefaultTenantName, data)
		require.NoErrorf(t, err, "failed to save data for data blocks")
	}

	for _, d := range data {
		t.Run("Data block "+d.TableName, func(t *testing.T) {
			query := fmt.Sprintf("{ %s { _id } }", d.TableName)
			result, err := s.Query(DefaultTenantName, query)
			require.NoError(t, err)
			require.Empty(t, result.Errors)
			assert.Len(t, result.Data.(map[string]interface{})[d.TableName], 1)
		})
	}
}

func TestNotUniqueUpdate(t *testing.T) {
	bCtx := env.NewBubblyContext()
	resource := test.RunPostgresDocker(bCtx, t)
	bCtx.StoreConfig.PostgresAddr = fmt.Sprintf("localhost:%s", resource.GetPort("5432/tcp"))
	// bCtx.StoreConfig.PostgresAddr = "localhost:5432"

	// Parse the schema and data blocks
	tables := testData.Tables(t, bCtx, "./testdata/unique/tables_update.hcl")
	data := testData.DataBlocks(t, bCtx, "./testdata/unique/data_update.hcl")
	// Initialize a new bubbly store (connection to postgres)
	s, err := New(bCtx)
	require.NoError(t, err)
	err = s.Apply(DefaultTenantName, tables, true)
	require.NoError(t, err)
	// Save it more than once to test the unique constraints
	for i := 0; i < 1; i++ {
		err = s.Save(DefaultTenantName, data)
		require.NoErrorf(t, err, "failed to save data for data blocks")
	}
	for _, d := range data {
		t.Run("Data block "+d.TableName, func(t *testing.T) {
			query := fmt.Sprintf("{ %s { _id } }", d.TableName)
			result, err := s.Query(DefaultTenantName, query)
			require.NoError(t, err)
			require.Empty(t, result.Errors)
			assert.Len(t, result.Data.(map[string]interface{})[d.TableName], 1)
		})
	}
}
