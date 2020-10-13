package v1

import (
	"fmt"

	"github.com/hashicorp/hcl/v2"
	"github.com/verifa/bubbly/api/core"
	"github.com/zclconf/go-cty/cty"
)

var _ core.Importer = (*Importer)(nil)

// Importer represents an importer type
type Importer struct {
	ResourceBlock *core.ResourceBlock

	Spec   ImporterSpec
	Source Source
}

// NewImporter returns a new Importer
func NewImporter(resBlock *core.ResourceBlock) *Importer {
	return &Importer{
		ResourceBlock: resBlock,
	}
}

// Decode implements the Block interfaces Decode() method and is responsible
// for decoding any necessary hcl.Body inside Importer
func (i *Importer) Decode(decode core.DecodeBodyFn) error {
	// decode the resource spec into the importer's Spec
	if err := decode(i.ResourceBlock.Spec.Body, &i.Spec); err != nil {
		return fmt.Errorf(`Failed to decode importer body spec: %s`, err.Error())
	}

	// based on the type of the importer, initiate the importer's Source
	switch i.Spec.Type {
	case jsonImporterType:
		i.Source = &JSONSource{}
	case xmlImporterType:
		i.Source = &XMLSource{}
	default:
		panic(fmt.Sprintf("Unsupported importer resource type %s", i.Spec.Type))
	}

	// decode the source HCL into the importer's Source
	if err := decode(i.Spec.SourceHCL.Body, i.Source); err != nil {
		return fmt.Errorf(`Failed to decode importer source: %s`, err.Error())
	}

	return nil
}

// String returns a string representation of the resource
func (i *Importer) String() string {
	return fmt.Sprintf(
		"%s.%s.%s",
		i.ResourceBlock.APIVersion, i.ResourceBlock.Kind, i.ResourceBlock.Name,
	)
}

// Resolve returns the cty.Value of the importer, or an error
func (i *Importer) Resolve() (cty.Value, error) {
	return cty.NilVal, nil
}

var _ core.ResourceSpec = (*ImporterSpec)(nil)

type ImporterSpec struct {
	// the type is either json, xml, rest, etc.
	Type      ImporterType `hcl:"type,attr"`
	SourceHCL struct {
		Body hcl.Body `hcl:",remain"`
	} `hcl:"source,block"`
	// the format of the raw input data defined as a cty.Type
	Format hcl.Expression `hcl:"format,attr"`
}

type ImporterType string

const (
	jsonImporterType ImporterType = "json"
	xmlImporterType               = "xml"
)

type Source interface {
	// returns an interface{} containing the parsed XML, JSON data, that should
	// be converted into the Output cty.Value
	Resolve() interface{}
}

var _ Source = (*JSONSource)(nil)

type JSONSource struct {
	File string `hcl:"file,attr"`
}

func (s *JSONSource) Resolve() interface{} {
	return nil
}

var _ Source = (*XMLSource)(nil)

type XMLSource struct {
	File string `hcl:"file,attr"`
}

func (s *XMLSource) Resolve() interface{} {
	return nil
}
