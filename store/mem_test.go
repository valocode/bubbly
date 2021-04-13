package store

import (
	"os"
	"path/filepath"
	"testing"

	testData "github.com/valocode/bubbly/store/testdata"

	"github.com/graphql-go/graphql"

	"github.com/stretchr/testify/assert"
)

func TestSchemaMem(t *testing.T) {
	tenantID := "key"
	fields := graphql.Fields{
		"hello": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return "world", nil
			},
		},
	}
	rootQuery := graphql.ObjectConfig{Name: "RootQuery", Fields: fields}
	schemaConfig := graphql.SchemaConfig{Query: graphql.NewObject(rootQuery)}
	testSchema, _ := graphql.NewSchema(schemaConfig)
	err := saveSchema(tenantID, &testSchema)
	assert.NoError(t, err)
	_, err = fetchSchema(tenantID)
	assert.NoError(t, err)
	t.Cleanup(func() {
		_ = os.Remove(dbPath)
	})
}

func TestTriggerMem(t *testing.T) {
	tenantID := "1234"
	err := saveTrigger(tenantID, internalTriggers)
	assert.NoError(t, err)
	_, err = fetchTriggers(tenantID)
	assert.NoError(t, err)
	t.Cleanup(func() {
		_ = os.Remove(dbPath)
	})
}

func TestSchemaGraphMem(t *testing.T) {
	tenantID := "1234"
	tables := testData.Tables(t, filepath.FromSlash("testdata/tables.hcl"))
	graph, err := newSchemaGraph(tables)
	assert.NoError(t, err)
	err = saveSchemaGraph(tenantID, graph)
	assert.NoError(t, err)
	_, err = fetchSchemaGraph(tenantID)
	assert.NoError(t, err)
	t.Cleanup(func() {
		_ = os.Remove(dbPath)
	})
}
