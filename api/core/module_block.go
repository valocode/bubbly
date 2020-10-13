package core

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/zclconf/go-cty/cty"
)

type ModuleBlocks []*ModuleBlock

var _ Block = (*ModuleBlock)(nil)

type ModuleBlock struct {
	Name      string     `hcl:",label"`
	Type      ModuleType `hcl:"type,attr"`
	SourceHCL *struct {
		Body hcl.Body `hcl:",remain"`
	} `hcl:"source,block"`
	Inputs []ModuleInput `hcl:"input,block"`

	Source ModuleSource
}

func NewLocalModuleBlock(baseDir string) *ModuleBlock {
	return &ModuleBlock{
		Name:      "root",
		Type:      LocalModuleType,
		SourceHCL: nil,
		Source: &ModuleLocalSource{
			BaseDir: baseDir,
			Path:    ".",
		},
	}
}

func (m *ModuleBlock) Decode(decFn DecodeBodyFn) error {

	if m.SourceHCL == nil {
		return nil
	}
	// if module is root then the Source is already configured
	switch t := ModuleType(m.Type); t {
	case LocalModuleType:
		// if local then the parent module must be local as well
		// localSource, ok := m.Source.(*ModuleLocalSource)
		mFile := m.SourceHCL.Body.MissingItemRange().Filename
		if mFile == "" {
			panic(fmt.Sprintf(`Error: trying to create local module from non-local parent module "%s"`, mFile))
		}
		mDir := filepath.Dir(mFile)
		m.Source = &ModuleLocalSource{
			// set the base directory for the ModuleLocalSource to be the
			// directory containing the block of hcl that contains the
			// definition of this module... from which the "path" attribute
			// is the relative path from
			BaseDir: mDir,
		}
	case RemoteModuleType:
		// TODO
		panic("ModuleType remote not implemented yet...")
	default:
		panic(fmt.Sprintf(`Invalid module type "%s"`, t))
	}

	decFn(m.SourceHCL.Body, m.Source)
	return nil
}

// func (m *ModuleBlock) Variables() []hcl.Traversal {
// 	traversals := []hcl.Traversal{}
// 	// for _, attr := range m.Params.Attributes {
// 	// 	traversals = append(traversals, attr.Expr.Variables()...)
// 	// }
// 	return traversals
// }

// func (m *ModuleBlock) DecodeTasks() DecodeTasks {
// 	if m.SourceHCL == nil {
// 		return nil
// 	}
// 	// if module is root then the Source is already configured
// 	switch t := ModuleType(m.Type); t {
// 	case LocalModuleType:
// 		// if local then the parent module must be local as well
// 		// localSource, ok := m.Source.(*ModuleLocalSource)
// 		mFile := m.SourceHCL.Body.MissingItemRange().Filename
// 		if mFile == "" {
// 			panic(fmt.Sprintf(`Error: trying to create local module from non-local parent module "%s"`, mFile))
// 		}
// 		mDir := filepath.Dir(
// 			m.SourceHCL.Body.MissingItemRange().Filename,
// 		)
// 		// println(localSource.BaseDir)
// 		m.Source = &ModuleLocalSource{
// 			// set the base directory for the ModuleLocalSource to be the
// 			// directory containing the block of hcl that contains the
// 			// definition of this module... from which the "path" attribute
// 			// is the relative path from
// 			BaseDir: mDir,
// 		}
// 	case RemoteModuleType:
// 		// TODO
// 		panic("ModuleType remote not implemented yet...")
// 	default:
// 		panic(fmt.Sprintf(`Invalid module type "%s"`, t))
// 	}

// 	// decode the module source
// 	return DecodeTasks{NewDecodeTask(m.SourceHCL.Body, m.Source)}
// }

func (m *ModuleBlock) Body() (hcl.Body, error) {
	return m.Source.Body()
}

type ModuleType string

const (
	LocalModuleType  ModuleType = "local"
	RemoteModuleType            = "remote"
)

type ModuleSource interface {
	Body() (hcl.Body, error)
}

var _ ModuleSource = &ModuleLocalSource{}

type ModuleLocalSource struct {
	BaseDir string
	Path    string `hcl:"path,attr"`
}

func (m *ModuleLocalSource) Body() (hcl.Body, error) {
	modulePath := filepath.Join(m.BaseDir, m.Path)
	files, err := filepath.Glob(fmt.Sprintf("%s/**.bubbly", modulePath))

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

	return hcl.MergeFiles(hclFiles), nil
}

type ModuleInput struct {
	Name  string    `hcl:",label"`
	Value cty.Value `hcl:"value,attr"`
}
