package builtin

import (
	"fmt"

	"github.com/valocode/bubbly/api/core"
	"github.com/valocode/bubbly/env"
	"github.com/valocode/bubbly/parser"
	"github.com/zclconf/go-cty/cty"
)

//go:generate go run gen/gen.go

const schemaFile = "schema.bubbly"

func BuiltinSchema() (core.Tables, error) {
	var schema SchemaWrapper
	bCtx := env.NewBubblyContext()

	err := parser.ParseFilename(bCtx, schemaFile, &schema)
	if err != nil {
		return nil, fmt.Errorf("error parsing bubbly schema: %w", err)
	}

	return schema.Tables, nil
}

func table(name string, fields []core.TableField, Joins []core.TableJoin) core.Table {
	return core.Table{
		Name:   name,
		Fields: fields,
		Joins:  Joins,
	}
}

func fields(fields ...core.TableField) []core.TableField {
	return fields
}

func field(name string, ty cty.Type, unique bool) core.TableField {
	return core.TableField{
		Name:   name,
		Unique: unique,
		Type:   ty,
	}
}

func joins(joins ...core.TableJoin) []core.TableJoin {
	return joins
}

func join(name string, single bool, unique bool) core.TableJoin {
	return core.TableJoin{
		Table:  name,
		Single: single,
		Unique: unique,
	}
}
