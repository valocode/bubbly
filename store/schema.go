package store

import (
	"encoding/json"
	"fmt"

	"github.com/valocode/bubbly/api/core"
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

	schema := &bubblySchema{
		Tables: make(map[string]core.Table, len(internalTables)),
	}

	for _, t := range internalTables {
		schema.Tables[t.Name] = t
	}

	return schema
}

func newBubblySchemaFromTables(tables core.Tables) *bubblySchema {
	schemaTables := make(map[string]core.Table)
	// Append the internal tables containing definition of the schema and
	// resource tables.
	tables = append(tables, internalTables...)
	for _, table := range tables {
		schemaTables[table.Name] = table
	}
	schema := &bubblySchema{
		Tables: schemaTables,
	}
	addImplicitJoins(schema, tables, nil)
	return schema
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
		Fields: map[string]cty.Value{
			"tables": cty.StringVal(string(bTables)), // "Smuggle" the JSON as a string
		},
	}, nil
}

func addImplicitJoins(schema *bubblySchema, tables core.Tables, parent *core.Table) {
	for index := range tables {
		// Get a reference to the table so that we can modify it
		t := tables[index]
		if parent != nil {
			var hasParentID bool
			// Check if the parent was already added to the schema
			for _, f := range t.Joins {
				if f.Table == parent.Name {
					hasParentID = true
				}
			}
			if !hasParentID {
				t.Joins = append(t.Joins, core.TableJoin{
					Table:  parent.Name,
					Single: t.Single,
				})
			}
		}

		addImplicitJoins(schema, t.Tables, &t)
		// Clear the child tables
		t.Tables = nil
		schema.Tables[t.Name] = t
	}
}
