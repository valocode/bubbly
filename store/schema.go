package store

import (
	"encoding/json"
	"fmt"

	"github.com/verifa/bubbly/api/core"
	"github.com/zclconf/go-cty/cty"
)

func newBubblySchema() *bubblySchema {
	return &bubblySchema{
		Tables: make(map[string]core.Table),
	}
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
