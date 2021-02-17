package core

import (
	"encoding/json"
	"fmt"

	"github.com/verifa/bubbly/parser"
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
	TableName string     `hcl:",label" json:"data"`
	Fields    DataFields `hcl:"fields" json:"fields"`
	Joins     []string   `hcl:"joins,optional" json:"joins"`
	Data      DataBlocks `hcl:"data,block" json:"nested_data"`
}

// DataFields is a slice of DataField
type DataFields map[string]cty.Value

// UnmarshalJSON unmarshals json into a Data type.
// This is a bit hacky, but the problem was making sure Fields (which is a map)
// is initialized before it gets unmarshaled. Perhaps there is a cleaner way,
// but for now this works and is not that ugly.
func (d *Data) UnmarshalJSON(data []byte) error {
	v := struct {
		TableName string     `json:"data"`
		Fields    DataFields `json:"fields"`
		Joins     []string   `json:"joins"`
		Data      DataBlocks `json:"nested_data"`
	}{
		Fields: make(DataFields),
	}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	d.TableName = v.TableName
	d.Fields = v.Fields
	d.Joins = v.Joins
	d.Data = v.Data
	return nil
}

// IsValidResource verifies that any of Data's fields that are mandatory for a
// valid core.Resource are non-null
func (d *Data) IsValidResource() bool {
	for k, v := range d.Fields {
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
	var jsonFields = make([]DataFieldJSON, 0, len(d))
	for name, val := range d {
		jsonVal, dataRef := ctyToJSON(val)
		jsonFields = append(jsonFields, DataFieldJSON{
			Name:    name,
			Value:   jsonVal,
			DataRef: dataRef,
		})
	}
	return json.Marshal(jsonFields)
}

// UnmarshalJSON unmarshals json into DataFields
func (d DataFields) UnmarshalJSON(data []byte) error {
	var jsonFields []DataFieldJSON
	if err := json.Unmarshal(data, &jsonFields); err != nil {
		return fmt.Errorf("failed to unmarshal DataFields: %w", err)
	}
	for _, field := range jsonFields {
		d[field.Name] = jsonToCty(field.DataRef, field.Value)
	}

	return nil
}

// DataFieldJSON is a json friendly version of DataField
type DataFieldJSON struct {
	Name    string                  `json:"name"`
	Value   ctyjson.SimpleJSONValue `json:"value,omitempty"`
	DataRef *parser.DataRef         `json:"data_ref,omitempty"`
}

func ctyToJSON(value cty.Value) (ctyjson.SimpleJSONValue, *parser.DataRef) {
	var dataRef parser.DataRef
	switch {
	case value.Type() == parser.DataRefType:
		dataRef = *value.EncapsulatedValue().(*parser.DataRef)
	default:
		return ctyjson.SimpleJSONValue{Value: value}, nil
	}
	return ctyjson.SimpleJSONValue{Value: cty.NilVal}, &dataRef
}

func jsonToCty(dataRef *parser.DataRef, value ctyjson.SimpleJSONValue) cty.Value {
	switch {
	case dataRef != nil:
		return cty.CapsuleVal(parser.DataRefType, dataRef)
	default:
		return value.Value
	}
}
