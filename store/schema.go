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

// bubblySchema contains the bubblySchema in a useable form, which is
// currently a map of the tables. It also contains a list of expected
// changes that will be applied by the migration.
type bubblySchema struct {
	Tables    map[string]core.Table
	changelog schemaUpdates
}

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
