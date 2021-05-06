package parser

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/valocode/bubbly/env"
	"github.com/zclconf/go-cty/cty"
)

func ParseFilename(bCtx *env.BubblyContext, filename string, val interface{}) error {
	mergedBody, err := MergedHCLBodies(bCtx, filename)
	if err != nil {
		return err
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

func MergedHCLBodies(bCtx *env.BubblyContext, filename string) (hcl.Body, error) {
	files, err := bubblyFiles(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to get bubbly files: %s", err.Error())
	}

	if len(files) == 0 {
		return nil, errors.New("no bubbly files found to parse")
	}

	parser := hclparse.NewParser()
	hclFiles := []*hcl.File{}
	for _, file := range files {
		hclFile, diags := parser.ParseHCLFile(file)
		if diags.HasErrors() {
			return nil, fmt.Errorf("failed to parse bubbly file: %s: %s", file, diags.Error())
		}
		hclFiles = append(hclFiles, hclFile)
	}
	mergedBody := hcl.MergeFiles(hclFiles)

	return mergedBody, nil
}

func bubblyFiles(filename string) ([]string, error) {
	fi, err := os.Stat(filename)
	if err != nil {
		return nil, fmt.Errorf(`cannot read from filename "%s"`, filename)
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
		return nil, fmt.Errorf(`unknown filename mode %s`, mode.String())
	}
}
