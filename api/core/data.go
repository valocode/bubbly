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
	TableName   string     `hcl:",label" json:"data"`
	Fields      DataFields `hcl:"field,block" json:"fields"`
	Data        DataBlocks `hcl:"data,block" json:"nested_data"`
	ParentTable string     `json:"-"`
}

// DataFields is a slice of DataField
type DataFields []DataField

// DataField is a field within a Data block
type DataField struct {
	Name  string    `hcl:",label" json:"name"`
	Value cty.Value `hcl:"value,attr" json:"-"`
}

// MarshalJSON implements a JSON marshaller for DataField
func (f *DataField) MarshalJSON() ([]byte, error) {
	d, err := NewJSONField(f)
	if err != nil {
		return nil, fmt.Errorf("failed to convert field into json: %w", err)
	}
	return json.Marshal(d)
}

// UnmarshalJSON implements a JSON unmarshaller for DataField
func (f *DataField) UnmarshalJSON(data []byte) error {
	var jf JSONDataField
	if err := json.Unmarshal(data, &jf); err != nil {
		return fmt.Errorf("failed to unmarshal DataField: %w", err)
	}
	*f = jf.DataField()
	return nil
}

// DataFieldAlias is an alias to avoid a recursive stack overflow with JSONDataField
type DataFieldAlias DataField

// JSONDataField is a JSON-friendly version of Field
type JSONDataField struct {
	DataFieldAlias
	Value   ctyjson.SimpleJSONValue `json:"value"`
	DataRef *parser.DataRef         `json:"data_ref"`
}

// DataField returns a DataField equivalent of JSONDataField
func (f *JSONDataField) DataField() DataField {
	field := DataField(f.DataFieldAlias)
	// if JSONDataField stored a DataRef instead of a native cty.Value, then
	// use that value. Else, use the native cty.Value
	if f.DataRef != nil {
		field.Value = cty.CapsuleVal(parser.DataRefType, f.DataRef)
	} else {
		field.Value = f.Value.Value
	}
	return field
}

// NewJSONField creates a new JSONDataField based on the given DataField.
// Noteworthy here is that it handles Capsule Values separately from the
// SimpleJSONValue.
func NewJSONField(f *DataField) (*JSONDataField, error) {
	var retData *JSONDataField
	switch {
	case f.Value.Type().IsCapsuleType():
		dataRef, ok := f.Value.EncapsulatedValue().(*parser.DataRef)
		if !ok {
			return nil, fmt.Errorf("invalid DataRef: %s", f.Value.GoString())
		}
		retData = &JSONDataField{
			DataFieldAlias: DataFieldAlias(*f),
			// assign NilVal to Value as we are using CapsuleVal
			Value: ctyjson.SimpleJSONValue{
				Value: cty.NilVal,
			},
			DataRef: dataRef,
		}
	default:
		retData = &JSONDataField{
			DataFieldAlias: DataFieldAlias(*f),
			Value: ctyjson.SimpleJSONValue{
				Value: f.Value,
			},
		}
	}
	return retData, nil
}
