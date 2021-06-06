package core

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/valocode/bubbly/parser"
	"github.com/zclconf/go-cty/cty"
	ctyjson "github.com/zclconf/go-cty/cty/json"
)

// DataRefs is a slice of DataRef.
type DataRefs []parser.DataRef

// DataBlocks is a slice of type Data
type DataBlocks []Data

// Data will reference a Table name, and assign the Field values into the
// corresponding Field values in the Table
type Data struct {
	TableName     string          `hcl:",label" json:"table"`
	Fields        *DataFields     `hcl:"fields,block" json:"fields,omitempty"`
	Lifecycle     *Lifecycle      `hcl:"lifecycle,block" json:"lifecycle,omitempty"`
	Joins         []string        `hcl:"joins,optional" json:"joins,omitempty"`
	Policy        DataBlockPolicy `hcl:"policy,optional" json:"policy,omitempty"`
	IgnoreNesting bool            `hcl:"ignore_nesting,optional" json:"ignore_nesting,omitempty"`
	Data          DataBlocks      `hcl:"data,block" json:"data,omitempty"`
}

// DataBlockPolicy defines the policy for how the data block shall be handled.
// When the bubbly store goes to save a data block, it should consider whether
// it should create and/or update the data block (default behaviour), only
// create the data block (fail on conflict), or only reference an existing data
// block with matching field values so that another data block can join to it
type DataBlockPolicy string

const (
	EmptyPolicy DataBlockPolicy = ""
	// DefaultPolicy is to create or update
	DefaultPolicy DataBlockPolicy = CreateUpdatePolicy
	// CreateUpdatePolicy means either create or update an existing data block
	// based on the unique constraints applied to the schema table that this data
	// block refers to
	CreateUpdatePolicy DataBlockPolicy = "create_update"
	// CreatePolicy means only create. If a conflict occurs on unique constraints
	// on the corresponding schema table, then error
	CreatePolicy DataBlockPolicy = "create"
	// UpdatePolicy means only update. If the entity to update does not exist
	// an error is returned
	UpdatePolicy DataBlockPolicy = "update"
	// ReferencePolicy means do not create or update, but only retrieve a reference
	// to an already saved data block, with the matching field values
	ReferencePolicy DataBlockPolicy = "reference"
	// ReferenceIfExistsPolicy is the same as ReferencePolicy but it does not
	// error in case a reference does not exist
	ReferenceIfExistsPolicy DataBlockPolicy = "reference_if_exists"
)

// DataFields contains a map of values that can be assigned to, e.g.
// fields {
// 	  my_val = "abc"
// 	  other_val = 123
// }
type DataFields struct {
	Values map[string]cty.Value `hcl:",remain"`
}

// UnmarshalJSON unmarshals json into a Data type.
// This is a bit hacky, but the problem was making sure Fields (which is a map)
// is initialized before it gets unmarshaled. Perhaps there is a cleaner way,
// but for now this works and is not that ugly.
func (d *Data) UnmarshalJSON(data []byte) error {
	v := struct {
		TableName     string          `json:"table"`
		Fields        *DataFields     `json:"fields"`
		Joins         []string        `json:"joins"`
		Policy        DataBlockPolicy `json:"policy"`
		IgnoreNesting bool            `json:"ignore_nesting"`
		Data          DataBlocks      `json:"data"`
	}{}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	d.TableName = v.TableName
	d.Fields = v.Fields
	d.Joins = v.Joins
	d.Policy = v.Policy
	d.IgnoreNesting = v.IgnoreNesting
	d.Data = v.Data
	return nil
}

// IsValidResource verifies that any of Data's fields that are mandatory for a
// valid core.Resource are non-null
func (d *Data) IsValidResource() bool {
	for k, v := range d.Fields.Values {
		if k == "metadata" {
			continue
		}
		if v.IsNull() {
			return false
		}
	}
	return true
}

// MarshalJSON marshals DataFields into json
func (d DataFields) MarshalJSON() ([]byte, error) {
	var jsonFields = make([]DataFieldJSON, 0, len(d.Values))
	for name, val := range d.Values {
		if val.Type().IsCapsuleType() {
			switch val.Type() {
			case parser.DataRefType:
				jsonFields = append(jsonFields, DataFieldJSON{
					Name:    name,
					DataRef: val.EncapsulatedValue().(*parser.DataRef),
				})
			case parser.TimeType:
				jsonFields = append(jsonFields, DataFieldJSON{
					Name: name,
					Time: val.EncapsulatedValue().(*time.Time),
				})
			}
			continue
		}
		jsonFields = append(jsonFields, DataFieldJSON{
			Name:  name,
			Value: ctyjson.SimpleJSONValue{Value: val},
		})
	}
	return json.Marshal(jsonFields)
}

// UnmarshalJSON unmarshals json into DataFields
func (d *DataFields) UnmarshalJSON(data []byte) error {
	var jsonFields []DataFieldJSON
	if err := json.Unmarshal(data, &jsonFields); err != nil {
		return fmt.Errorf("failed to unmarshal DataFields: %w", err)
	}
	if d.Values == nil {
		d.Values = make(map[string]cty.Value)
	}
	for _, field := range jsonFields {
		switch {
		case !field.Value.IsNull():
			d.Values[field.Name] = field.Value.Value
		case field.DataRef != nil:
			d.Values[field.Name] = cty.CapsuleVal(parser.DataRefType, field.DataRef)
		case field.Time != nil:
			d.Values[field.Name] = cty.CapsuleVal(parser.TimeType, field.Time)
		}
	}

	return nil
}

type Lifecycle struct {
	Status  string           `hcl:"status,optional" json:"status,omitempty"`
	Entries []LifecycleEntry `hcl:"entry,block" json:"entries,omitempty"`
}

type LifecycleEntry struct {
	Message string `hcl:"message" json:"message"`
}

// DataFieldJSON is a json friendly version of DataField
type DataFieldJSON struct {
	Name    string                  `json:"name"`
	Value   ctyjson.SimpleJSONValue `json:"value,omitempty"`
	DataRef *parser.DataRef         `json:"data_ref,omitempty"`
	Time    *time.Time              `json:"time,omitempty"`
}
