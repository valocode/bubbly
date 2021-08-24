package adapter

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/hashicorp/hcl/v2/hclsimple"
	"github.com/valocode/bubbly/ent"
	entadapter "github.com/valocode/bubbly/ent/adapter"
)

const DefaultTag = "default"

type (
	// Adapter is the main definition of an adapter.
	//
	// An adapter is one of the main forms of input into bubbly and is created
	// to extract data from somewhere (e.g. JSON file or REST api) and map the
	// data to a bubbly data graph which can then be processed by the bubbly
	// server.
	Adapter struct {
		Name      string    `hcl:"name,attr"`
		Tag       *string   `hcl:"tag,optional"`
		Operation Operation `json:"operation,omitempty" hcl:"operation,block"`
		Results   Results   `json:"results,omitempty" hcl:"results,block"`
	}

	// RunOptions defines the options about the adapter to run, such as the
	// name, tag, or from file
	RunOptions struct {
		Name     string
		Tag      string
		Filename string
		BaseDir  string
	}
	// RunArgs defines the runtime arguments that are provided to the adapter
	// which come from user input
	RunArgs struct {
		Filename string
	}
	// adapterWrap is used for wrapping an Adapter to be decoded using gohcl
	adapterWrap struct {
		Adapter Adapter `hcl:"adapter,block"`
	}
)

func FromModel(model *ent.AdapterModelRead) (*Adapter, error) {
	a := &Adapter{
		Name: *model.Name,
		Tag:  model.Tag,
		Operation: Operation{
			Type: string(*model.Type),
		},
		Results: Results{
			Type: string(*model.ResultsType),
		},
	}
	if err := json.Unmarshal(*model.Operation, &a.Operation); err != nil {
		return nil, err
	}
	if err := a.Results.FromBytes(a.String(), *model.Results); err != nil {
		return nil, err
	}
	return a, nil
}

func FromFile(filename string) (*Adapter, error) {
	src, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("could not read file %s: %w", filename, err)
	}
	return Decode(filename, src)
}

func Decode(filename string, src []byte) (*Adapter, error) {
	var adaptWrap adapterWrap
	err := hclsimple.Decode(filename, src, nil, &adaptWrap)
	if err != nil {
		return nil, err
	}

	a := adaptWrap.Adapter
	if err := a.Operation.Decode(); err != nil {
		return nil, fmt.Errorf("error decoding operation: %w", err)
	}

	return &a, nil
}

func (a *Adapter) Run(args RunArgs) (*Output, error) {

	opValue, err := a.Operation.Spec.Run(a.Name, args)
	if err != nil {
		return nil, fmt.Errorf("error performing operation %s: %w", a.Operation.Type, err)
	}

	eCtx := NewEvalContext(nil)
	eCtx.Variables["data"] = opValue

	if err := a.Results.Decode(eCtx); err != nil {
		return nil, fmt.Errorf("error decoding results: %w", err)
	}

	return a.Results.Spec.Output(a.Name)
}

func (a *Adapter) Model() (*ent.AdapterModelCreate, error) {
	op, err := json.Marshal(a.Operation)
	if err != nil {
		return nil, fmt.Errorf("json marshalling adapter operation: %w", err)
	}
	results, err := a.Results.SpecBytes()
	if err != nil {
		return nil, fmt.Errorf("converting results body to bytes: %w", err)
	}
	return ent.NewAdapterModelCreate().
		SetName(a.Name).
		SetTag(a.TagOrDefault()).
		SetType(entadapter.Type(a.Operation.Type)).
		SetOperation(op).
		SetResultsType(entadapter.ResultsType(a.Results.Type)).
		SetResults(results), nil
}

func (a *Adapter) String() string {
	return a.Name + ":" + a.TagOrDefault()
}

func (a *Adapter) TagOrDefault() string {
	var tag = DefaultTag
	if a.Tag != nil {
		tag = *a.Tag
	}
	return tag
}

func ParseAdpaterID(id string) (string, string, error) {
	splitID := strings.Split(id, ":")
	switch len(splitID) {
	case 1:
		return id, "", nil
	case 2:
		return splitID[0], splitID[1], nil
	}
	return "", "", fmt.Errorf("adapter must be in the form \"name:tag\"")
}
