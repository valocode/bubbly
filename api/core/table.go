package core

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/ext/typeexpr"
	"github.com/valocode/bubbly/parser"
	"github.com/zclconf/go-cty/cty"
)

// Tables holds a slice of table
type Tables []Table

func (t Tables) Resolve() error {
	for _, tbl := range t {
		if err := tbl.Resolve(); err != nil {
			return err
		}
	}
	return nil
}

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

func (t *Table) Resolve() error {
	var diags hcl.Diagnostics
	for idx := range t.Fields {
		var (
			typeDiags hcl.Diagnostics
			// Get a pointer to the field so that it can be updated
			field = &t.Fields[idx]
		)
		exprWord := hcl.ExprAsKeyword(field.TypeExpr)
		if exprWord == "time" {
			field.Type = parser.TimeType
			continue
		}
		field.Type, typeDiags = typeexpr.TypeConstraint(field.TypeExpr)
		diags = append(diags, typeDiags...)
	}
	if diags.HasErrors() {
		return parser.NewParserError(nil, diags)
	}
	// Recursively go through child tables
	for _, tbl := range t.Tables {
		if err := tbl.Resolve(); err != nil {
			return err
		}
	}
	return nil
}

// TableField is a schema field.
type TableField struct {
	Name     string         `hcl:",label" json:"name"`
	Unique   bool           `hcl:"unique,optional" json:"unique,omitempty"`
	TypeExpr hcl.Expression `hcl:"type,attr" json:"-"`
	Type     cty.Type       `json:"type"`
}

func (f TableField) MarshalJSON() ([]byte, error) {
	alias := TableFieldAlias(f)
	if f.Type.IsCapsuleType() {
		if f.Type.Equals(parser.TimeType) {
			js := TableFieldJSON{
				TableFieldAlias: alias,
			}
			js.TypeJSON = TypeJSON(f.Type)
			return json.Marshal(js)
		} else {
			return nil, fmt.Errorf("unsupported capsule type for field %s: %s", f.Name, f.Type.GoString())
		}
	}
	return json.Marshal(alias)
}

func (f *TableField) UnmarshalJSON(data []byte) error {
	var js TableFieldJSON
	if err := json.Unmarshal(data, &js); err != nil {
		return fmt.Errorf("error unmarshalling table field: %w", err)
	}
	*f = TableField(js.TableFieldAlias)
	f.Type = cty.Type(js.TypeJSON)
	return nil
}

type TableFieldAlias TableField

type TypeJSON cty.Type

func (t TypeJSON) MarshalJSON() ([]byte, error) {
	ty := cty.Type(t)
	if ty.IsCapsuleType() {
		if ty.Equals(parser.TimeType) {
			return []byte{'"', 't', 'i', 'm', 'e', '"'}, nil
		}
	}
	return json.Marshal(ty)
}
func (t *TypeJSON) UnmarshalJSON(b []byte) error {
	r := bytes.NewReader(b)
	tok, err := json.NewDecoder(r).Token()
	if err != nil {
		return err
	}
	if str, ok := tok.(string); ok {
		if str == "time" {
			*t = TypeJSON(parser.TimeType)
			return nil
		}
	}
	var ty cty.Type
	err = json.Unmarshal(b, &ty)
	if err != nil {
		return err
	}
	*t = TypeJSON(ty)
	return nil
}

type TableFieldJSON struct {
	TableFieldAlias
	TypeJSON TypeJSON `json:"type"`
}

type TableJoin struct {
	Table  string `hcl:",label" json:"name"`
	Unique bool   `hcl:"unique,optional" json:"unique,omitempty"`
	Single bool   `hcl:"single,optional" json:"single,omitempty"`
}
