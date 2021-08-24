package adapter

import (
	"fmt"
	"os"
	"reflect"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/ext/dynblock"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/valocode/bubbly/ent"
	"github.com/valocode/bubbly/store/api"
)

type (
	Results struct {
		Type string   `json:"type,omitempty" hcl:",label"`
		Body hcl.Body `json:"-" hcl:",remain"`

		Spec ResultsSpec `json:"-"`
	}
	ResultsAlias Results

	Output struct {
		CodeScan *api.CodeScan `json:"code_scan,omitempty"`
		TestRun  *api.TestRun  `json:"test_run,omitempty"`
	}

	ResultsSpec interface {
		Output(name string) (*Output, error)
	}

	CodeScanResults struct {
		Tool       string                   `json:"tool,omitempty" hcl:"tool,optional"`
		CodeIssues []*api.CodeScanIssue     `json:"issues,omitempty" hcl:"issue,block"`
		Components []*api.CodeScanComponent `json:"components,omitempty"  hcl:"component,block"`
	}

	TestRunResults struct {
		Tool string `hcl:"tool,optional"`
	}
)

func (r *Results) Decode(eCtx *hcl.EvalContext) error {
	switch ResultsType(r.Type) {
	case ResultsCodeScan:
		r.Spec = &CodeScanResults{}
	case ResultsTestRun:
		r.Spec = &TestRunResults{}
	default:
		return fmt.Errorf("unsupported results type %s", r.Type)
	}
	dynBody := dynblock.Expand(r.Body, eCtx)
	diags := gohcl.DecodeBody(dynBody, eCtx, r.Spec)
	if diags.HasErrors() {
		return fmt.Errorf("error decoding code scan: %w", diags)
	}
	return nil
}

func (r *Results) FromBytes(id string, b []byte) error {
	file, diags := hclparse.NewParser().ParseHCL(b, id)
	if diags.HasErrors() {
		return fmt.Errorf("decoding results: %w", diags)
	}
	r.Body = file.Body
	return nil
}

func (r Results) SpecBytes() ([]byte, error) {
	if r.Body == nil {
		return nil, fmt.Errorf("cannot get src from nil results body")
	}
	// get the source range of the hcl body, so that we can extract it as raw text
	var srcRange hcl.Range
	switch body := r.Body.(type) {
	case *hclsyntax.Body:
		srcRange = body.SrcRange
	default:
		return nil, fmt.Errorf("cannot get src range for unknown hcl.Body type %s", reflect.TypeOf(body).String())
	}
	// read the bubbly file containing the HCL
	src, err := os.ReadFile(srcRange.Filename)
	if err != nil {
		return nil, fmt.Errorf("reading adapter file: %w", err)
	}
	if !srcRange.CanSliceBytes(src) {
		return nil, fmt.Errorf("cannot slice bytes for adapter results in filename %s", srcRange.Filename)
	}
	specBytes := srcRange.SliceBytes(src)
	// specBytes contains the block paranthesis "{" and "}". Remove them
	specBytes = specBytes[1 : len(specBytes)-1]
	// Format the bytes so that they become "pretty"
	fmtBytes := hclwrite.Format(specBytes)

	return fmtBytes, nil
}

func (r *CodeScanResults) Output(name string) (*Output, error) {
	var tool = name
	if r.Tool != "" {
		tool = r.Tool
	}
	return &Output{
		CodeScan: &api.CodeScan{
			CodeScanModelCreate: ent.NewCodeScanModelCreate().SetTool(tool),
			Issues:              r.CodeIssues,
			Components:          r.Components,
		},
	}, nil
}

func (r *TestRunResults) Output(name string) (*Output, error) {
	return nil, fmt.Errorf("TestRunResults Output TODO")
}

func (o *Output) HasCodeScan() bool {
	return o.CodeScan != nil
}

func (o *Output) HasTestRun() bool {
	return o.TestRun != nil
}

type ResultsType string

const (
	ResultsCodeScan ResultsType = "code_scan"
	ResultsTestRun  ResultsType = "test_run"
)

func (r ResultsType) String() string {
	return string(r)
}
