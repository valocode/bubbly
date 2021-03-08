package bubbly

import (
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"

	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/zclconf/go-cty/cty"

	"github.com/valocode/bubbly/api/core"
	"github.com/valocode/bubbly/client"
	"github.com/valocode/bubbly/env"
	"github.com/valocode/bubbly/parser"
)

// Schema is the Go-native struct representation of a bubbly
// schema file
type Schema struct {
	Tables core.Tables `hcl:"table,block"`
}

// ApplySchema parses a .bubbly schema file into a Schema, then posts
// the core.Tables of the Schema to the bubbly store
func ApplySchema(bCtx *env.BubblyContext, file string) error {
	var schema Schema

	err := parseSchemaFile(bCtx, file, &schema)

	if err != nil {
		return fmt.Errorf(
			`failed to parse schema file at "%s": %w`,
			filepath.ToSlash(file),
			err)
	}

	tableBytes, err := json.Marshal(schema.Tables)
	if err != nil {
		return fmt.Errorf("failed to json marshal schema tables: %w", err)
	}

	c, err := client.New(bCtx)
	if err != nil {
		return fmt.Errorf("failed to create bubbly HTTP client: %w", err)
	}

	if err := c.PostSchema(bCtx, tableBytes); err != nil {
		return fmt.Errorf("failed to post schema to bubbly server: %w", err)
	}

	return nil
}

// parseSchemaFile reads and parses a .bubbly schema file
// and decodes the schema into the provided interface,
// which is typically of type Schema
func parseSchemaFile(bCtx *env.BubblyContext, file string, val interface{}) error {
	hclFile, diags := hclparse.NewParser().ParseHCLFile(file)
	if diags != nil {
		return errors.New(diags.Error())
	}

	return parser.DecodeExpandBody(bCtx, hclFile.Body, val, cty.NilVal)
}
