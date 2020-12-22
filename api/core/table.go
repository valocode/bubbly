package core

import "github.com/zclconf/go-cty/cty"

// Tables holds a slice of table
type Tables []Table

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
	Name   string   `hcl:",label" json:"table"`
	Unique bool     `hcl:"unique,attr" json:"unique"`
	Type   cty.Type `hcl:"type,attr" json:"type"`
}
