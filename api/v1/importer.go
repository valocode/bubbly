package v1

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"reflect"
	"strings"

	"github.com/clbanning/mxj"
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

	Spec importerSpec `json:"spec"`
}

// NewImporter returns a new Importer
func NewImporter(resBlock *core.ResourceBlock) *Importer {
	return &Importer{
		ResourceBlock: resBlock,
	}
}

// Apply returns the output from applying a resource
func (i *Importer) Apply(ctx *core.ResourceContext) core.ResourceOutput {
	if err := i.decode(ctx.DecodeBody); err != nil {
		return core.ResourceOutput{
			Status: core.ResourceOutputFailure,
			Error:  fmt.Errorf("Failed to decode resource %s: %s", i.String(), err.Error()),
		}
	}

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

	return core.ResourceOutput{
		Status: core.ResourceOutputSuccess,
		Error:  nil,
		Value:  val,
	}
}

// SpecValue method returns resource specification structure
func (i *Importer) SpecValue() core.ResourceSpec {
	return &i.Spec
}

// decode is responsible for decoding any necessary hcl.Body inside Importer
func (i *Importer) decode(decode core.DecodeBodyFn) error {
	// decode the resource spec into the importer's Spec
	if err := decode(i, i.SpecHCL.Body, &i.Spec); err != nil {
		return fmt.Errorf(`Failed to decode "%s" body spec: %s`, i.String(), err.Error())
	}

	// based on the type of the importer, initiate the importer's Source
	switch i.Spec.Type {
	case jsonImporterType:
		i.Spec.Source = &jsonSource{}
	case xmlImporterType:
		i.Spec.Source = &xmlSource{}
	default:
		panic(fmt.Sprintf("Unsupported importer resource type %s", i.Spec.Type))
	}

	// decode the source HCL into the importer's Source
	if err := decode(i, i.Spec.SourceHCL.Body, i.Spec.Source); err != nil {
		return fmt.Errorf(`Failed to decode importer source: %s`, err.Error())
	}

	return nil
}

var _ core.ResourceSpec = (*importerSpec)(nil)

// importerSpec defines the spec for an importer
type importerSpec struct {
	Inputs InputDeclarations `hcl:"input,block"`
	// the type is either json, xml, rest, etc.
	Type      importerType `hcl:"type,attr"`
	SourceHCL *struct {
		Body hcl.Body `hcl:",remain"`
	} `hcl:"source,block"`
	// Source stores the actual value for SourceHCL
	Source source
}

// importerType defines the type of an importer
type importerType string

const (
	jsonImporterType importerType = "json"
	xmlImporterType               = "xml"
)

// Source is an interface for the different data sources that an Importer can have
type source interface {
	// returns an interface{} containing the parsed XML, JSON data, that should
	// be converted into the Output cty.Value
	Resolve() (cty.Value, error)
}

// Compiler check to see that v1.JSONSource implements the Source interface
var _ source = (*jsonSource)(nil)

// jsonSource represents the importer type for using a JSON file as the input
type jsonSource struct {
	File string `hcl:"file,attr"`
	// the format of the raw input data defined as a cty.Type
	Format cty.Type `hcl:"format,attr"`
}

// Resolve returns a cty.Value representation of the parsed JSON file
func (s *jsonSource) Resolve() (cty.Value, error) {

	var barr []byte
	var err error

	// FIXME GitHub issue #39
	barr, err = ioutil.ReadFile(s.File)
	if err != nil {
		return cty.NilVal, err
	}

	// Attempt to unmarshall the data into an empty interface data type
	var data interface{}
	err = json.Unmarshal(barr, &data)
	if err != nil {
		return cty.NilVal, err
	}

	val, err := gocty.ToCtyValue(data, s.Format)
	if err != nil {
		return cty.NilVal, nil
	}

	return val, nil
}

// Compiler check to see that v1.XMLSource implements the Source interface
var _ source = (*xmlSource)(nil)

// xmlSource represents the importer type for using an XML file as the input
type xmlSource struct {
	File string `hcl:"file,attr"`
	// the format of the raw input data defined as a cty.Type
	Format cty.Type `hcl:"format,attr"`
}

// Resolve returns a cty.Value representation of the XML file
func (s *xmlSource) Resolve() (cty.Value, error) {

	var barr []byte
	var err error

	// FIXME GitHub issue #39
	barr, err = ioutil.ReadFile(s.File)
	if err != nil {
		return cty.NilVal, err
	}

	mxj.PrependAttrWithHyphen(false) // no "-" prefix on attributes
	mxj.CastNanInf(true)             // use float64, not string for extremes

	// Unmarshall the XML data into a Go object
	data, err := mxj.NewMapXml(barr, true)
	if err != nil {
		return cty.NilVal, err
	}

	if err := walkTypeTransformData(&data, s.Format); err != nil {
		return cty.NilVal, err
	}

	val, err := gocty.ToCtyValue(data, s.Format)
	if err != nil {
		return cty.NilVal, err
	}

	return val, nil
}

func walkTypeTransformData(data *mxj.Map, ty cty.Type) error {
	path := make([]string, 0)
	return walk(data, ty, path, 0)
}

func walk(data *mxj.Map, ty cty.Type, path []string, idx int) error {

	pathStr := strings.Join(path, ".")

	if idx > 0 {
		pathStr += fmt.Sprint("[", idx, "]")
	}

	if ty.IsObjectType() {
		for x := range ty.AttributeTypes() {
			path = append(path, x)
			pathIdx := len(path) - 1

			walk(data, ty.AttributeType(x), path, 0)
			path = path[0:pathIdx]
		}
	}

	if ty.IsListType() {

		vs, err := data.ValuesForPath(pathStr)
		if err != nil {
			return fmt.Errorf("wrong path (%s) in xml structure: %w", pathStr, err)
		}

		n := len(vs)
		//t.Logf("ValuesForPath(%s): %d", pathStr, n)

		switch n {
		case 0:
			return fmt.Errorf("xml data structure inconsistent state, ValuesForPath are zero at %s", pathStr)
		case 1:
			v := vs[0]

			if reflect.TypeOf(v).Kind() == reflect.Map {
				vv := make([]interface{}, 0)
				vv = append(vv, v)
				if err := data.SetValueForPath(vv, pathStr); err != nil {
					return fmt.Errorf("cannot convert at path %s, error %w", pathStr, err)
				}
			}
			fallthrough
		default:
			for i := range vs {
				return walk(data, ty.ElementType(), path, i)
			}
		}
	}

	return nil
}
