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
	// Single describes a one-to-one relationship for an implicit join
	// If a Table is nested within another table it is an implicit join from the
	// nested Table to the parent Table. Just like joins have properties (like
	// single & unique), these implicit joins also need to have them
	Single bool `hcl:"single,optional" json:"single,omitempty"`
	// Unique makes an implicit join part of the unique constraint
	Unique bool    `hcl:"unique,optional" json:"unique,omitempty"`
	Tables []Table `hcl:"table,block" json:"tables,omitempty"`
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
	Single bool   `hcl:"single,optional" json:"single,omitempty"`
}
