package core

import "github.com/zclconf/go-cty/cty"

// Table is a schema table. It may
// contain fields, tables, or any
// combination of the two.
type Table struct {
	Name   string
	Fields []TableField
	Tables []Table
}

// TableField is a schema field.
type TableField struct {
	Name   string
	Unique bool
	Type   cty.Type
}
