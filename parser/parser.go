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

func ParseResource(bCtx *env.BubblyContext, id string, src []byte, value interface{}) error {
	hclParser := hclparse.NewParser()
	file, diags := hclParser.ParseHCL(src, id)
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
	var (
		files []string
	)
	fi, err := os.Stat(filename)
	if err != nil {
		return nil, fmt.Errorf("cannot read from filename %s", filename)
	}
	// Handle the file or directory that was provided
	switch mode := fi.Mode(); {
	case mode.IsRegular():
		files = append(files, filename)
		// Get the directory that the file is in
		filename = filepath.Dir(filename)
	case mode.IsDir():
		dirFiles, err := bubblyFilesInDir(filename)
		if err != nil {
			return nil, fmt.Errorf("error getting bubbly files in directory %s: %w", filename, err)
		}
		files = append(files, dirFiles...)
	default:
		return nil, fmt.Errorf("unknown filename mode %s", mode.String())
	}

	dir, err := filepath.Abs(filepath.Dir(filename))
	if err != nil {
		return nil, fmt.Errorf("error getting absolute path of %s: %w", filename, err)
	}
	// TODO: this will only work when the root filesystem is "/"
	for dir != "/" {
		dirFiles, err := bubblyFilesInDir(dir)
		if err != nil {
			return nil, fmt.Errorf("error getting bubbly files in directory %s: %w", filename, err)
		}
		files = append(files, dirFiles...)
		dir = filepath.Dir(dir)
	}
	return files, nil
}

func bubblyFilesInDir(dir string) ([]string, error) {
	var files []string
	// walk the directory and get .bubbly files
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("error opening directory %s: %w", dir, err)
	}
	for _, e := range entries {
		if filepath.Ext(e.Name()) == ".bubbly" && !e.IsDir() {
			files = append(files, filepath.Join(dir, e.Name()))
		}
	}
	return files, err
}
