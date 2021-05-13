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
	files, err := bubblyFilesByFilename(filename)
	if err != nil {
		return fmt.Errorf("failed to get bubbly files: %s", err.Error())
	}
	mergedBody, err := MergedHCLBodies(bCtx, files)
	if err != nil {
		return err
	}
	if err := DecodeBody(mergedBody, val, cty.NilVal); err != nil {
		return fmt.Errorf(`failed to decode body: %s`, err.Error())
	}
	return nil
}

func ParseConfig(bCtx *env.BubblyContext, val interface{}) error {
	files, err := bubblyFilesWithConfig()
	if err != nil {
		return fmt.Errorf("failed to get bubbly files: %s", err.Error())
	}
	mergedBody, err := MergedHCLBodies(bCtx, files)
	if err != nil {
		return err
	}
	if err := DecodeBody(mergedBody, val, cty.NilVal); err != nil {
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
	if err := DecodeBody(file.Body, value, cty.NilVal); err != nil {
		return fmt.Errorf("failed to decode resource: %w", err)
	}
	return nil
}

func MergedHCLBodies(bCtx *env.BubblyContext, files []string) (hcl.Body, error) {

	if len(files) == 0 {
		return nil, errors.New("no bubbly files found")
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

func bubblyFilesByFilename(filename string) ([]string, error) {
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
	case mode.IsDir():
		// walk the directory and get .bubbly files
		entries, err := os.ReadDir(filename)
		if err != nil {
			return nil, fmt.Errorf("error opening directory %s: %w", filename, err)
		}
		for _, e := range entries {
			if filepath.Ext(e.Name()) == ".bubbly" && !e.IsDir() {
				files = append(files, filepath.Join(filename, e.Name()))
			}
		}
	default:
		return nil, fmt.Errorf("unknown filename mode %s", mode.String())
	}

	return files, nil
}

func bubblyFilesWithConfig() ([]string, error) {
	var files []string
	cwd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("error getting working directory: %w", err)
	}

	dir, err := filepath.Abs(filepath.Join(cwd, ".bubbly"))
	if err != nil {
		return nil, fmt.Errorf("error gettings absolute path: %w", err)
	}

	for {
		bfs, err := bubblyFilesByFilename(dir)
		if err != nil {
			return nil, fmt.Errorf("error getting bubbly files with %s: %w", dir, err)
		}
		files = append(files, bfs...)
		// Get the parent directory
		pDir := filepath.Dir(dir)
		if pDir == dir {
			// Then we are at the root directory so exit
			break
		}
		dir = pDir
	}
	return files, nil
}
