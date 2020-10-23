package v1

import (
	"encoding/json"
	"fmt"

	"github.com/verifa/bubbly/api/core"
	"github.com/zclconf/go-cty/cty"
	ctyjson "github.com/zclconf/go-cty/cty/json"
)

var _ core.Translator = (*Translator)(nil)

type Translator struct {
	*core.ResourceBlock
	Spec TranslatorSpec
}

func NewTranslator(resBlock *core.ResourceBlock) *Translator {
	return &Translator{
		ResourceBlock: resBlock,
	}
}

func (t *Translator) Decode(decode core.DecodeResourceFn) error {
	// decode the resource spec into the importer's Spec
	if err := decode(t, t.SpecHCL.Body, &t.Spec); err != nil {
		return fmt.Errorf(`Failed to decode translator body spec: %s`, err.Error())
	}
	return nil
}

func (t *Translator) JSON() ([]byte, error) {
	return json.Marshal(t.Spec.Data)
}

// Output returns ...
func (t *Translator) Output() core.ResourceOutput {
	return core.ResourceOutput{
		Status: core.ResourceOutputSuccess,
		Error:  nil,
		Value:  cty.NilVal,
	}
}

func (t *Translator) SpecValue() core.ResourceSpec {
	return &t.Spec
}

type TranslatorSpec struct {
	Inputs InputDeclarations `hcl:"input,block"`
	Data   DataBlocks        `hcl:"data,block"`
}

type DataBlocks []*Data

// Data will reference a Table name, and assign the Field values into the
// corresponding Field values in the Table
type Data struct {
	TableName string     `hcl:",label" json:"table"`
	RowName   string     `hcl:",label" json:"row"`
	Fields    Fields     `hcl:"field,block" json:"fields"`
	Data      DataBlocks `hcl:"data,block" json:"data"`
}

type Fields []*Field

type Field struct {
	Name  string    `hcl:",label" json:"name"`
	Value cty.Value `hcl:"value,attr" json:"-"`
}

// MashalJSON implements a JSON marshaller for Field
func (f *Field) MarshalJSON() ([]byte, error) {
	return json.Marshal(NewJSONField(f))
}

// UnmashalJSON implements a JSON unmarshaller for Field
func (f *Field) UnmarshalJSON(data []byte) error {
	var jf JSONField
	if err := json.Unmarshal(data, &jf); err != nil {
		return err
	}
	*f = jf.Field()
	return nil
}

// FieldAlias is an alias to avoid a recursive stack overflow with JSONField
type FieldAlias Field

// JSONField is a JSON-friendly version of Field, embedding the alias FieldAlias
type JSONField struct {
	FieldAlias
	Value ctyjson.SimpleJSONValue
}

// Field returns a Field equivalent of JSONField
func (f *JSONField) Field() Field {
	field := Field(f.FieldAlias)
	field.Value = f.Value.Value
	return field
}

// NewJSONField creates a new JSONField based on the given Field
func NewJSONField(f *Field) *JSONField {
	return &JSONField{
		FieldAlias: FieldAlias(*f),
		Value: ctyjson.SimpleJSONValue{
			Value: f.Value,
		},
	}
}
