package bubbly

import (
	"encoding/json"
	"fmt"
	"path/filepath"

	"github.com/valocode/bubbly/api/core"
	"github.com/valocode/bubbly/bubbly/builtin"
	"github.com/valocode/bubbly/client"
	"github.com/valocode/bubbly/env"
	"github.com/valocode/bubbly/parser"
)

// Schema is the Go-native struct representation of a bubbly
// schema file

// ApplySchema parses a .bubbly schema file into a Schema, then posts
// the []core.Table of the Schema to the bubbly store
func ApplySchema(bCtx *env.BubblyContext, file string) error {
	var schema builtin.SchemaWrapper

	err := parser.ParseFilename(bCtx, file, &schema)
	if err != nil {
		return fmt.Errorf(
			`failed to parse schema file at "%s": %w`,
			filepath.ToSlash(file),
			err)
	}
	tables, err := core.TablesFromHCL(schema.Tables)
	if err != nil {
		return fmt.Errorf("error creating schema tables: %w", err)
	}

	tableBytes, err := json.Marshal(tables)
	if err != nil {
		return fmt.Errorf("failed to json marshal schema tables: %w", err)
	}

	c, err := client.New(bCtx)
	if err != nil {
		return fmt.Errorf("failed to create bubbly HTTP client: %w", err)
	}
	defer c.Close()

	if err := c.PostSchema(bCtx, nil, tableBytes); err != nil {
		return fmt.Errorf("failed to post schema to bubbly server: %w", err)
	}

	return nil
}
