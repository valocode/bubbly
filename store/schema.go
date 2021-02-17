package store

import (
	"encoding/json"
	"fmt"

	"github.com/verifa/bubbly/api/core"
	"github.com/zclconf/go-cty/cty"
)

func newBubblySchema() *bubblySchema {
	schema := &bubblySchema{
		Tables: make(map[string]core.Table),
	}
	// Populate the fresh schema with the internal tables
	for _, t := range internalTables {
		schema.Tables[t.Name] = t
	}
	return schema
}

// bubblySchema contains the bubblySchema in a useable form, which is currently
// a map of the tables.
// This should be extended in the future to accommodate for schema diffing
// and other neat tricks.
type bubblySchema struct {
	Tables    map[string]core.Table
	Changelog Changelog
}

// Data returns a core.Data value of the schema so that it can be saved in the
// store just like any other data.
func (b *bubblySchema) Data() (core.Data, error) {
	bTables, err := json.Marshal(b.Tables)
	if err != nil {
		return core.Data{}, fmt.Errorf("failed to convert bubblySchema into data blocks: %w", err)
	}
	return core.Data{
		TableName: core.SchemaTableName,
		Fields: map[string]cty.Value{
			// "Smuggle" the json as a string
			"tables": cty.StringVal(string(bTables)),
		},
	}, nil
}
