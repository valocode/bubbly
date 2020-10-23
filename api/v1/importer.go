package v1

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/hashicorp/hcl/v2"
	"github.com/verifa/bubbly/api/core"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/gocty"
)

// Compiler check to see that v1.Importer implements the Importer interface
var _ core.Importer = (*Importer)(nil)

// Importer represents an importer type
type Importer struct {
	*core.ResourceBlock

	Spec ImporterSpec `json:"spec"`
}

// NewImporter returns a new Importer
func NewImporter(resBlock *core.ResourceBlock) *Importer {
	return &Importer{
		ResourceBlock: resBlock,
	}
}

// Decode implements the Block interfaces Decode() method and is responsible
// for decoding any necessary hcl.Body inside Importer
func (i *Importer) Decode(decode core.DecodeResourceFn) error {
	// decode the resource spec into the importer's Spec
	if err := decode(i, i.SpecHCL.Body, &i.Spec); err != nil {
		return fmt.Errorf(`Failed to decode "%s" body spec: %s`, i.String(), err.Error())
	}

	// based on the type of the importer, initiate the importer's Source
	switch i.Spec.Type {
	case jsonImporterType:
		i.Spec.Source = &JSONSource{}
	case xmlImporterType:
		i.Spec.Source = &XMLSource{}
	default:
		panic(fmt.Sprintf("Unsupported importer resource type %s", i.Spec.Type))
	}

	// decode the source HCL into the importer's Source
	if err := decode(i, i.Spec.SourceHCL.Body, i.Spec.Source); err != nil {
		return fmt.Errorf(`Failed to decode importer source: %s`, err.Error())
	}

	return nil
}

func (i *Importer) SpecValue() core.ResourceSpec {
	return &i.Spec
}

// Output returns the output from applying a resource
func (i *Importer) Output() core.ResourceOutput {
	if i == nil {
		return core.ResourceOutput{
			Status: core.ResourceOutputFailure,
			Error:  errors.New("Cannot get output of a null importer"),
			Value:  cty.NilVal,
		}
	}

	if i.Spec.Source == nil {
		return core.ResourceOutput{
			Status: core.ResourceOutputFailure,
			Error:  errors.New("Cannot get output of an importer with null source"),
			Value:  cty.NilVal,
		}
	}

	val, err := i.Spec.Source.Resolve()
	if err != nil {
		return core.ResourceOutput{
			Status: core.ResourceOutputFailure,
			Error:  fmt.Errorf("Failed to resolve importer source: %s", err.Error()),
			Value:  cty.NilVal,
		}
	}

	ctyVal, err := gocty.ToCtyValue(val, i.Spec.Format)
	if err != nil {
		return core.ResourceOutput{
			Status: core.ResourceOutputFailure,
			Error:  fmt.Errorf("Failed to convert import source to format: %s", err.Error()),
			Value:  cty.NilVal,
		}
	}

	return core.ResourceOutput{
		Status: core.ResourceOutputSuccess,
		Error:  nil,
		Value:  ctyVal,
	}
}

var _ core.ResourceSpec = (*ImporterSpec)(nil)

// ImporterSpec defines the spec for an importer
type ImporterSpec struct {
	Inputs InputDeclarations `hcl:"input,block"`
	// the type is either json, xml, rest, etc.
	Type      ImporterType `hcl:"type,attr"`
	SourceHCL struct {
		Body hcl.Body `hcl:",remain"`
	} `hcl:"source,block"`
	// Source stores the actual value for SourceHCL
	Source Source
	// the format of the raw input data defined as a cty.Type
	Format cty.Type `hcl:"format,attr"`
}

// ImporterType defines the type of an importer
type ImporterType string

const (
	jsonImporterType ImporterType = "json"
	xmlImporterType               = "xml"
)

// Source is an interface for the different data sources that an Importer
// can have
type Source interface {
	// returns an interface{} containing the parsed XML, JSON data, that should
	// be converted into the Output cty.Value
	Resolve() (interface{}, error)
}

var _ Source = (*JSONSource)(nil)

type JSONSource struct {
	File string `hcl:"file,attr"`
}

func (s *JSONSource) Resolve() (interface{}, error) {

	var barr []byte
	var err error

	// FIXME reading the whole file at once may be too much
	barr, err = ioutil.ReadFile(s.File)
	if err != nil {
		return nil, err
	}

	// Attempt to unmarshall the data into an empty interface data type
	var data interface{}
	err = json.Unmarshal(barr, &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

var _ Source = (*XMLSource)(nil)

type XMLSource struct {
	File string `hcl:"file,attr"`
}

func (s *XMLSource) Resolve() (interface{}, error) {
	return nil, errors.New("not implemented")
}
