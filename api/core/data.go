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
	Joins       DataJoins  `hcl:"join,block" json:"joins"`
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

// DataFieldJSON is a json friendly version of DataField
type DataFieldJSON struct {
	Name    string                  `json:"name"`
	Value   ctyjson.SimpleJSONValue `json:"value,omitempty"`
	DataRef *parser.DataRef         `json:"data_ref,omitempty"`
}

// MarshalJSON implements a JSON marshaller for DataField
func (d *DataField) MarshalJSON() ([]byte, error) {
	jsonVal, dataRef := ctyToJSON(d.Value)

	return json.Marshal(&DataFieldJSON{
		Name:    d.Name,
		Value:   jsonVal,
		DataRef: dataRef,
	})
}

// UnmarshalJSON implements a JSON unmarshaller for DataField
func (d *DataField) UnmarshalJSON(data []byte) error {
	var j DataFieldJSON
	if err := json.Unmarshal(data, &j); err != nil {
		return fmt.Errorf("failed to unmarshal DataFieldJSON: %w", err)
	}

	*d = DataField{
		Name:  j.Name,
		Value: jsonToCty(j.DataRef, j.Value),
	}
	return nil
}

// DataJoins is a slice of DataJoin
type DataJoins []DataJoin

// DataJoin is a join within a Data block
type DataJoin struct {
	Table string    `hcl:",label" json:"table"`
	Value cty.Value `hcl:"value,attr" json:"-"`
}

// DataJoinJSON is a json friendly version of DataJoin
type DataJoinJSON struct {
	Table   string                  `json:"table"`
	Value   ctyjson.SimpleJSONValue `json:"value,omitempty"`
	DataRef *parser.DataRef         `json:"data_ref,omitempty"`
}

// MarshalJSON implements a JSON marshaller for DataJoin
func (d *DataJoin) MarshalJSON() ([]byte, error) {
	jsonVal, dataRef := ctyToJSON(d.Value)

	return json.Marshal(&DataJoinJSON{
		Table:   d.Table,
		Value:   jsonVal,
		DataRef: dataRef,
	})
}

// UnmarshalJSON implements a JSON unmarshaller for DataJoin
func (d *DataJoin) UnmarshalJSON(data []byte) error {
	var j DataJoinJSON
	if err := json.Unmarshal(data, &j); err != nil {
		return fmt.Errorf("failed to unmarshal DataJoinJSON: %w", err)
	}

	*d = DataJoin{
		Table: j.Table,
		Value: jsonToCty(j.DataRef, j.Value),
	}
	return nil
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
