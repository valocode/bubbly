package parser

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/verifa/bubbly/env"
	"github.com/zclconf/go-cty/cty"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclparse"
)

// Parser is the main type in the parser package and represents a type for
// parsing HCL code, and maintaining the known resources.
//
// The parser can be used for both HCL and a JSON representations of HCL,
// however the Body member of Parser must be from an HCL file.
type Parser struct {
	Body  hcl.Body
	Scope *Scope

	HCLParser *hclparse.Parser
}

// NewParserFromFilename creates a new parser from the given filename.
func NewParserFromFilename(filename string) (*Parser, error) {
	files, err := bubblyFiles(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to get bubbly files: %s", err.Error())
	}

	if len(files) == 0 {
		return nil, errors.New("No bubbly files found to parse")
	}

	parser := hclparse.NewParser()
	hclFiles := []*hcl.File{}
	for _, file := range files {
		hclFile, diags := parser.ParseHCLFile(file)
		if diags.HasErrors() {
			return nil, errors.New(diags.Error())
		}
		hclFiles = append(hclFiles, hclFile)
	}

	mergedBody := hcl.MergeFiles(hclFiles)

	return newParser(mergedBody, parser), nil
}

// WithInputs returns a parser with the provided inputs added to the
// EvalContext
func WithInputs(bCtx *env.BubblyContext, inputs cty.Value) *Parser {
	// bCtx.Logger.Debug().Str("inputs", inputs.GoString()).Msg("Creating parser with inputs")
	p := newParser(nil, nil)
	p.Scope.InsertValue(bCtx, inputs, []string{"self"})
	return p
}

// EmptyParser returns an empty parser that can be used, e.g. to turn a JSON
// string into a resource without parsing any local files
func EmptyParser() *Parser {
	return newParser(nil, nil)
}

func ParseResource(bCtx *env.BubblyContext, src []byte, value interface{}) error {
	p := newParser(nil, nil)
	hclParser := hclparse.NewParser()
	file, diags := hclParser.ParseHCL(src, "TODO")
	if diags.HasErrors() {
		return fmt.Errorf("failed to parse resource: %s", diags.Error())
	}
	diags = p.Scope.decodeBody(bCtx, file.Body, value)
	if diags.HasErrors() {
		return fmt.Errorf("failed to decode resource: %s", diags.Error())
	}
	return nil
}

func newParser(body hcl.Body, hclparser *hclparse.Parser) *Parser {

	if hclparser == nil {
		hclparser = hclparse.NewParser()
	}

	return &Parser{
		Body:      body,
		Scope:     NewScope(),
		HCLParser: hclparser,
	}
}

// Parse performs the actual parsing for the parser
func (p *Parser) Parse(bCtx *env.BubblyContext, value interface{}) error {
	if p.Body == nil {
		return errors.New("Trying to run the parser on an empty body")
	}

	if err := p.Scope.decodeBody(bCtx, p.Body, value); err != nil {
		return fmt.Errorf(`failed to decode main body: %s`, err.Error())
	}
	return nil
}

func (p *Parser) ParseBytes(bCtx *env.BubblyContext, src []byte, value interface{}) error {
	file, diags := p.HCLParser.ParseHCL(src, "TODO")
	if diags.HasErrors() {
		return fmt.Errorf("failed to parse bytes: %s", diags.Error())
	}
	return p.Scope.DecodeExpandBody(bCtx, file.Body, value)
}

func bubblyFiles(filename string) ([]string, error) {
	fi, err := os.Stat(filename)
	if err != nil {
		return nil, fmt.Errorf(`Cannot read from filename "%s"`, filename)
	}
	switch mode := fi.Mode(); {
	case mode.IsRegular():
		return []string{filename}, nil
	case mode.IsDir():
		// walk the directory and get .bubbly files
		files := []string{}
		filepath.Walk(filename, func(path string, file os.FileInfo, err error) error {
			if filepath.Ext(path) == ".bubbly" {
				files = append(files, path)
			}
			return nil
		})
		return files, nil
	default:
		return nil, fmt.Errorf(`Unknown filename mode %s`, mode.String())
	}
}
