package interim

import (
	"errors"
	"fmt"

	"github.com/hashicorp/go-memdb"
	"github.com/verifa/bubbly/api/core"
	"github.com/zclconf/go-cty/cty"
)

const (
	idFieldName = "ID"
	fkFieldName = "FK"
)

// newMemDB creates a new memdb for the given tables.
func newMemDB(tables []core.Table) (*memdb.MemDB, error) {
	if len(tables) == 0 {
		return nil, errors.New("at least one table is required")
	}

	schema := &memdb.DBSchema{
		Tables: make(map[string]*memdb.TableSchema, len(tables)),
	}

	for _, t := range tables {
		addTableToMemDB(t, schema)
	}

	db, err := memdb.NewMemDB(schema)
	if err != nil {
		return nil, fmt.Errorf("failed to create DB: %w", err)
	}

	return db, nil
}

// flattenTables takes a set of nested tables and flattens them
// into a single layer of *memdb.TableSchemas. The names of the
// table schemas are namespaces to represent where they were in
// the original nested hierarchy.
func addTableToMemDB(t core.Table, schema *memdb.DBSchema) {
	schema.Tables[t.Name] = tableSchema(t)
	for _, subT := range t.Tables {
		addTableToMemDB(subT, schema)
	}
}

// tableSchema returns a *memdb.TableSchema representation
// of the current table recursively.
func tableSchema(t core.Table) *memdb.TableSchema {
	schema := &memdb.TableSchema{
		Name:    t.Name,
		Indexes: make(map[string]*memdb.IndexSchema, len(t.Fields)),
	}

	for _, f := range t.Fields {
		schema.Indexes[f.Name] = fieldSchema(f)
	}

	// "ID" and "FK" are not part of the original
	// data. We will supply they as data is inserted
	// to allow us to query between tables. Note that
	// the name of the "id" index must be lowercase
	// for memdb.
	schema.Indexes["id"] = &memdb.IndexSchema{
		Name:    "id",
		Unique:  true,
		Indexer: &memdb.UUIDFieldIndex{Field: idFieldName},
	}
	schema.Indexes[fkFieldName] = &memdb.IndexSchema{
		Name:    fkFieldName,
		Unique:  false,
		Indexer: &memdb.UUIDFieldIndex{Field: fkFieldName},
	}

	return schema
}

// fieldSchema returns a *memdb.IndexSchema representation
// of the current field and the name that it be stored under.
func fieldSchema(f core.TableField) *memdb.IndexSchema {
	schema := &memdb.IndexSchema{
		Name:   f.Name,
		Unique: f.Unique,
	}

	switch f.Type {
	case cty.Bool:
		schema.Indexer = &memdb.BoolFieldIndex{Field: f.Name}
	case cty.Number:
		schema.Indexer = &memdb.IntFieldIndex{Field: f.Name}
	case cty.String:
		schema.Indexer = &memdb.StringFieldIndex{Field: f.Name}
	}

	return schema
}
