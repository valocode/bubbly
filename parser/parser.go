package parser

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/verifa/bubbly/config"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/verifa/bubbly/api"
	"github.com/verifa/bubbly/api/core"
	"github.com/zclconf/go-cty/cty"
)

// Parser is the main type in the parser package and represents a type for
// parsing HCL code, and maintaining the known resources.
//
// The parser can be used for both HCL and a JSON representations of HCL,
// however the Body member of Parser must be from an HCL file.
type Parser struct {
	Body      hcl.Body
	Scope     *Scope
	Value     core.HCLMainType
	Resources api.Resources

	HCLParser *hclparse.Parser
}

// NewParserFromFilename creates a new parser from the given filename.
func NewParserFromFilename(filename string) (*Parser, error) {
	files, err := bubblyFiles(filename)
	if err != nil {
		return nil, fmt.Errorf("Failed to get bubbly files: %s", err.Error())
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

func newParser(body hcl.Body, hclparser *hclparse.Parser) *Parser {

	if hclparser == nil {
		hclparser = hclparse.NewParser()
	}

	return &Parser{
		Body:      body,
		Scope:     NewScope(),
		Resources: *api.NewResources(),
		HCLParser: hclparser,
	}
}

// Parse performs the actual parsing for the parser
func (p *Parser) Parse() error {
	if p.Body == nil {
		return errors.New("Trying to run the parser on an empty body")
	}

	if err := p.Scope.decodeBody(p.Body, &p.Value); err != nil {
		return fmt.Errorf(`Failed to decode main body: %s`, err.Error())
	}

	// populate the parser's EvalContext with variables
	p.populateEvalContext()

	// populate the parser's resources
	p.populateResources()

	return nil
}

// Context produces a context for resources to decode themselves and get
// resources, as well as managing some of the EvalContext scope and creating
// a nested context.
func (p *Parser) Context(inputs cty.Value) *core.ResourceContext {
	// create the nested scope to assign the "self" traversal
	nestedScope := p.Scope.NestedScope(inputs)
	return &core.ResourceContext{
		GetResource:     p.GetResource,
		DecodeBody:      nestedScope.DecodeExpandBody,
		NewContext:      p.Context,
		InsertValue:     nestedScope.InsertValue,
		BodyToJSON:      p.BodyToJSON,
		GetServerConfig: p.GetServerConfig,
		Debug: func() *hcl.EvalContext {
			return nestedScope.EvalContext
		},
	}
}

// GetServerConfig returns the server configuration defined by viper bindings
func (p *Parser) GetServerConfig() (*config.ServerConfig, error) {
	c, err := config.SetupConfigs()

	if err != nil {
		return nil, fmt.Errorf("failed to create config.ServerConfig: %w", err)
	}
	return c.ServerConfig, nil
}

// GetResource takes a given resource kind and name and returns the resource
// if it exists, either locally in parsed HCL files or on the bubbly server.
func (p *Parser) GetResource(kind core.ResourceKind, name string) (core.Resource, error) {
	if resource := p.Resources.Get(kind, name); resource != nil {
		return resource, nil
	}
	// TODO: we should fetch a resource using the bubbly server REST API
	return nil, fmt.Errorf(`Could not obtain resource "%s" of kind %s`, name, string(kind))
}

func (p *Parser) populateEvalContext() {
	for _, local := range p.Value.Locals {
		traversal, value := local.Reference()
		p.Scope.insert(value, traversal)
	}
}

func (p *Parser) populateResources() {
	for _, resBlock := range p.Value.ResourceBlocks {
		p.Resources.NewResource(resBlock)
	}
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
