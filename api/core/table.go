package core

import "github.com/zclconf/go-cty/cty"

// Table is a schema table. It may
// contain fields, tables, or any
// combination of the two.
type Table struct {
	Name   string       `hcl:",label"`
	Fields []TableField `hcl:"field,block"`
	Tables []Table      `hcl:"table,block"`
}

// TableField is a schema field.
type TableField struct {
	Name   string   `hcl:",label"`
	Unique bool     `hcl:"unique,attr"`
	Type   cty.Type `hcl:"type,attr"`
}
