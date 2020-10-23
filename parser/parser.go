package parser

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/verifa/bubbly/api"
	"github.com/verifa/bubbly/api/core"
)

type Parser struct {
	Body      hcl.Body
	Scope     *Scope
	Value     core.HCLMainType
	Resources api.Resources

	HCLParser *hclparse.Parser
}

func NewParser(baseDir string) (*Parser, error) {
	files, err := filepath.Glob(fmt.Sprintf("%s/**.bubbly", baseDir))

	if err != nil {
		return nil, err
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

func (p *Parser) Decode() error {
	if err := p.Scope.decodeBody(p.Body, &p.Value); err != nil {
		return fmt.Errorf(`Failed to decode main body: %s`, err.Error())
	}

	// populate the parser's EvalContext with variables
	p.populateEvalContext()

	// populate the parser's resources
	p.populateResources()

	return nil
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
