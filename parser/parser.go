package parser

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/verifa/bubbly/env"
	"github.com/zclconf/go-cty/cty"
)

func ParseFilename(bCtx *env.BubblyContext, filename string, val interface{}) error {
	files, err := bubblyFiles(filename)
	if err != nil {
		return fmt.Errorf("failed to get bubbly files: %s", err.Error())
	}

	if len(files) == 0 {
		return errors.New("no bubbly files found to parse")
	}

	parser := hclparse.NewParser()
	hclFiles := []*hcl.File{}
	for _, file := range files {
		hclFile, diags := parser.ParseHCLFile(file)
		if diags.HasErrors() {
			return errors.New(diags.Error())
		}
		hclFiles = append(hclFiles, hclFile)
	}

	mergedBody := hcl.MergeFiles(hclFiles)

	if mergedBody == nil {
		return fmt.Errorf("HCL body is nil")
	}

	if err := DecodeBody(bCtx, mergedBody, val, cty.NilVal); err != nil {
		return fmt.Errorf(`failed to decode body: %s`, err.Error())
	}
	return nil
}

func ParseResource(bCtx *env.BubblyContext, src []byte, value interface{}) error {
	hclParser := hclparse.NewParser()
	file, diags := hclParser.ParseHCL(src, "TODO_FROM_JSON")
	if diags.HasErrors() {
		return fmt.Errorf("failed to parse resource: %s", diags.Error())
	}
	if err := DecodeBody(bCtx, file.Body, value, cty.NilVal); err != nil {
		return fmt.Errorf("failed to decode resource: %w", err)
	}
	return nil
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
