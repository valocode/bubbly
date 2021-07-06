package store

import (
	"encoding/json"
	"fmt"

	"github.com/valocode/bubbly/api/core"
	"github.com/valocode/bubbly/bubbly/builtin"
	"github.com/zclconf/go-cty/cty"
)

//
// The Bubbly Schema is an abstraction describing the database tables
// containing the user data and Bubbly's internal data structures,
// as well as the connections between them.
//

// newBubblySchema makes an empty data structure, representing
// a new Bubbly Schema, and then populates it with the minimum
// set of tables necessary for Bubbly to work.
func newBubblySchema() *bubblySchema {

	tables := FlattenTables(builtin.BuiltinTables, nil)
	schema := &bubblySchema{
		Tables: make(map[string]core.Table, len(tables)),
	}

	for _, t := range tables {
		schema.Tables[t.Name] = t
	}
	return schema
}

func newBubblySchemaFromTables(tables []core.Table, internal bool) (*bubblySchema, error) {
	schemaTables := make(map[string]core.Table)
	// Create the base schema with the builtin tables.
	// First flatten them, so that we get implicit joins and can easily loop
	// through them
	builtinTables := FlattenTables(builtin.BuiltinTables, nil)
	for _, table := range builtinTables {
		schemaTables[table.Name] = table
	}
	// Compare the given tables with the builtin tables.
	// If a builtin table is trying to be changed, product an error
	tables = FlattenTables(tables, nil)
	for _, table := range tables {
		// If the request was not made internally, prevent end users from modifying
		// internal/builtin tables
		if !internal {
			_, ok := schemaTables[table.Name]
			if ok {
				return nil, fmt.Errorf("cannot modify builtin table %s", table.Name)
			}
		}
		schemaTables[table.Name] = table
	}
	schema := &bubblySchema{
		Tables: schemaTables,
	}
	return schema, nil
}

// bubblySchema contains the bubblySchema in a useable form, which is currently
// a map of the tables.
// This should be extended in the future to accommodate for schema diffing
// and other neat tricks.
type bubblySchema struct {
	Tables    map[string]core.Table
	changelog schemaUpdates
}

// Data returns a representation of the Bubbly Schema in a format,
// suitable for saving it in the Store, just like any other data.
func (b *bubblySchema) Data() (core.Data, error) {

	bTables, err := json.Marshal(b.Tables)

	if err != nil {
		return core.Data{}, fmt.Errorf("failed to convert bubblySchema into data blocks: %w", err)
	}

	return core.Data{
		TableName: core.SchemaTableName,
		Fields: &core.DataFields{Values: map[string]cty.Value{
			"tables": cty.StringVal(string(bTables)), // "Smuggle" the JSON as a string
		}},
	}, nil
}

// FlattenTables takes a list of tables flattens any nested tables, making sure
// the joins implied by the nesting are added.
// The table is already a flat list, this should return an identical list
func FlattenTables(tables []core.Table, parent *core.Table) []core.Table {
	var curTables []core.Table
	for _, t := range tables {
		if parent != nil {
			var hasParentID bool
			// Check if the join already exists, so that we don't add it twice
			for _, join := range t.Joins {
				if join.Table == parent.Name {
					hasParentID = true
				}
			}
			if !hasParentID {
				t.Joins = append(t.Joins, core.TableJoin{
					Table:  parent.Name,
					Single: t.Single,
					Unique: t.Unique,
				})
			}
		}
		childTables := FlattenTables(t.Tables, &t)
		// Clear the child tables
		t.Tables = nil
		curTables = append(curTables, t)
		curTables = append(curTables, childTables...)
	}
	return curTables
}
