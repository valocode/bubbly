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

// Table is a schema table and defines the relationships (joins), fields,
// unique constraints, etc
type Table struct {
	Name   string       `json:"name"`
	Fields []TableField `json:"fields"`
	Joins  []TableJoin  `json:"joins,omitempty"`
	// Single describes a one-to-one relationship for an implicit join
	// If a Table is nested within another table it is an implicit join from the
	// nested Table to the parent Table. Just like joins have properties (like
	// single & unique), these implicit joins also need to have them
	Single bool `json:"single,omitempty"`
	// Unique makes an implicit join part of the unique constraint
	Unique bool    `json:"unique,omitempty"`
	Tables []Table `json:"tables,omitempty"`
}

// TableHCL is almost a replica of Table, but is used for parsing tables
// from HCL files, and should be converted into a Table type which resolves
// fields like the cty.Type.
// The reason for separating the types is to encourage compile/parse errors,
// rather than difficult/hard to catch errors by combining the two types
type TableHCL struct {
	Name   string          `hcl:",label"`
	Fields []TableFieldHCL `hcl:"field,block"`
	Joins  []TableJoin     `hcl:"join,block"`
	Single bool            `hcl:"single,optional"`
	Unique bool            `hcl:"unique,optional"`
	Tables []TableHCL      `hcl:"table,block"`
}

// TablesFromHCL takes tables parsed from HCL and returns a slice of Table
func TablesFromHCL(hclTables []TableHCL) ([]Table, error) {
	var tables []Table
	for _, ht := range hclTables {
		table, err := TableFromHCL(ht)
		if err != nil {
			return nil, err
		}
		tables = append(tables, table)
	}

	return tables, nil
}

// TableFromHCL takes a table parsed from HCL and returns a Table
func TableFromHCL(hclTable TableHCL) (Table, error) {
	var (
		table Table
		err   error
	)
	table.Name = hclTable.Name
	table.Joins = hclTable.Joins
	table.Single = hclTable.Single
	table.Unique = hclTable.Unique
	table.Fields, err = FieldsFromHCL(hclTable.Fields)
	if err != nil {
		return Table{}, err
	}
	table.Tables, err = TablesFromHCL(hclTable.Tables)
	if err != nil {
		return Table{}, err
	}
	return table, nil
}

// FieldsFromHCL takes fields parsed from HCL and returns a slice of TableField
func FieldsFromHCL(hclFields []TableFieldHCL) ([]TableField, error) {
	var fields []TableField
	for _, f := range hclFields {
		field, err := FieldFromHCL(f)
		if err != nil {
			return nil, err
		}
		fields = append(fields, field)
	}
	return fields, nil
}

// FieldFromHCL takes a field parsed from HCL and returns TableField
func FieldFromHCL(hclField TableFieldHCL) (TableField, error) {
	var (
		field TableField
		err   error
	)

	field.Name = hclField.Name
	field.Unique = hclField.Unique
	field.Type, err = typeFromTypeExpr(hclField.TypeExpr)
	if err != nil {
		return TableField{}, err
	}
	return field, nil
}

// typeFromTypeExpr takes an hcl.Expression which is expected to be a type expr
// (i.e. representing a cty.Type), and returns the cty.Type or error if
// the type was invalid
func typeFromTypeExpr(expr hcl.Expression) (cty.Type, error) {
	// Some special cases, like time, are handled internally to Bubbly
	exprWord := hcl.ExprAsKeyword(expr)
	if exprWord == "time" {
		return parser.TimeType, nil
	}
	// If no special case, delegate to typeexpr
	ty, diags := typeexpr.TypeConstraint(expr)
	if diags.HasErrors() {
		return cty.NilType, fmt.Errorf("invalid type: %s", diags.Error())
	}
	return ty, nil
}

// TableField is a schema field.
type TableField struct {
	Name   string   `json:"name"`
	Unique bool     `json:"unique,omitempty"`
	Type   cty.Type `json:"type"`
}

type TableFieldHCL struct {
	Name     string         `hcl:",label"`
	Unique   bool           `hcl:"unique,optional"`
	TypeExpr hcl.Expression `hcl:"type,attr"`
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

// An alias for cty.Type so that we can override JSON marshal/unmarshal
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

// TableJoin can be the same type used for HCL parsing because there are no
// differences (yet)
type TableJoin struct {
	Table  string `hcl:",label" json:"name"`
	Unique bool   `hcl:"unique,optional" json:"unique,omitempty"`
	Single bool   `hcl:"single,optional" json:"single,omitempty"`
}
