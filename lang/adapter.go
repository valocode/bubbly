package lang

import (
	"fmt"
	"os"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclsimple"
	"github.com/valocode/bubbly/ent"
	"github.com/zclconf/go-cty/cty"
)

type (
	// Adapter is the main definition of an adapter.
	//
	// An adapter is one of the main forms of input into bubbly and is created
	// to extract data from somewhere (e.g. JSON file or REST api) and map the
	// data to a bubbly data graph which can then be processed by the bubbly
	// server.
	Adapter struct {
		Name string `hcl:"name,attr"`
		Tag  string `hcl:"tag,optional"`
		Type string `hcl:"type,attr"`

		Operation *AdapterOperation `hcl:"operation,block"`
		Mapping   AdapterMapping    `hcl:"mapping,block"`
	}
	// AdapterOperation contains the HCL body that will be decoded to the
	// specific operation based on the type of the adapter
	AdapterOperation struct {
		Body hcl.Body `hcl:",remain"`
	}
	// AdapterMapping contains the HCL body for creating the bubbly data graph
	// that will be sent to the bubbly server to store the data
	AdapterMapping struct {
		Body hcl.Body `hcl:",remain"`
	}
	// AdapterOptions defines the runtime options that are provided to the adapter
	// which can come from command line flags or environment variables
	AdapterOptions struct {
		Filename string
	}
	// adapterWrap is used for wrapping an Adapter to be decoded using gohcl
	adapterWrap struct {
		Adapter Adapter `hcl:"adapter,block"`
	}
	adapterOp interface {
		Run(adapter *Adapter, opts AdapterOptions) (cty.Value, error)
	}
)

func NewAdapterFromFile(filename string) (*Adapter, error) {

	src, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("could not read file %s: %w", filename, err)
	}
	return DecodeAdapter(filename, src)

}

func DecodeAdapter(filename string, src []byte) (*Adapter, error) {
	var adaptWrap adapterWrap
	err := hclsimple.Decode(filename, src, nil, &adaptWrap)
	if err != nil {
		return nil, err
	}

	adapter := adaptWrap.Adapter
	if err := adapter.Validate(); err != nil {
		return nil, err
	}
	return &adapter, nil
}

func (a *Adapter) Run(opts AdapterOptions) (*ent.DataGraph, error) {
	var (
		eCtx = NewEvalContext(nil)
		op   adapterOp
	)
	switch AdapterSource(a.Type) {
	case JSONSource:
		op = &jsonOp{}
		diags := gohcl.DecodeBody(a.Operation.Body, eCtx, op)
		if diags.HasErrors() {
			return nil, fmt.Errorf("error decoding operation for adapter %s: %s", a.Name, diags.Error())
		}
	default:
		return nil, fmt.Errorf("unsupported adapter type \"%s\"", a.Type)
	}
	opValue, err := op.Run(a, opts)
	if err != nil {
		return nil, fmt.Errorf("error performing operation %s: %w", a.Type, err)
	}

	eCtx.Variables["data"] = opValue

	graph, err := ent.DecodeDataGraph(a.Mapping.Body, eCtx)
	if err != nil {
		return nil, fmt.Errorf("error decoding data graph: %w", err)
	}
	return graph, nil
}

func (a *Adapter) Validate() error {
	// TODO: options, preprocess, type, format, etc
	// if err := a.SetFormat(); err != nil {
	// 	return err
	// }
	return nil
}

type AdapterSource string

const (
	HTTPSource AdapterSource = "http"
	JSONSource AdapterSource = "json"
	XMLSource  AdapterSource = "xml"
)

func (a AdapterSource) String() string {
	return string(a)
}

func NewEvalContext(variables map[string]cty.Value) *hcl.EvalContext {
	if variables == nil {
		variables = make(map[string]cty.Value)
	}
	return &hcl.EvalContext{
		Variables: variables,
		Functions: stdfunctions(),
	}
}
