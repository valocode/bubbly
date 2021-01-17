package core

import "github.com/zclconf/go-cty/cty"

// Tables holds a slice of table
type Tables []Table

// Table is a schema table. It may
// contain fields, tables, or any
// combination of the two.
type Table struct {
	Name   string       `hcl:",label" json:"name"`
	Fields []TableField `hcl:"field,block" json:"fields"`
	Joins  []TableJoin  `hcl:"join,block" json:"joins,omitempty"`
	Unique bool         `hcl:"unique,optional" json:"unique,omitempty"`
	Tables []Table      `hcl:"table,block" json:"tables,omitempty"`
}

// TableField is a schema field.
type TableField struct {
	Name   string   `hcl:",label" json:"name"`
	Unique bool     `hcl:"unique,optional" json:"unique,omitempty"`
	Type   cty.Type `hcl:"type,attr" json:"type"`
}

type TableJoin struct {
	Table  string `hcl:",label" json:"name"`
	Unique bool   `hcl:"unique,optional" json:"unique,omitempty"`
}
