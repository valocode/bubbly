package core

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/zclconf/go-cty/cty"
)

// ModuleBlocks is a wrapper type for a list of ModuleBlock instances
type ModuleBlocks []*ModuleBlock

// Ensure that ModuleBlock implements the Block interface
var _ Block = (*ModuleBlock)(nil)

// ModuleBlock represents a block of type module in HCL.
type ModuleBlock struct {
	Name      string     `hcl:",label"`
	Type      moduleType `hcl:"type,attr"`
	SourceHCL *struct {
		Body hcl.Body `hcl:",remain"`
	} `hcl:"source,block"`
	Inputs []moduleInput `hcl:"input,block"`

	Source moduleSource
}

// NewLocalModuleBlock returns a new ModuleBlock of type "local".
// This method is used when creating the "root" module.
func NewLocalModuleBlock(baseDir string) *ModuleBlock {
	return &ModuleBlock{
		Name:      "root",
		Type:      localmoduleType,
		SourceHCL: nil,
		Source: &ModuleLocalSource{
			BaseDir: baseDir,
			Path:    ".",
		},
	}
}

// Decode decodes any hcl.Bodies inside ModuleBlock
func (m *ModuleBlock) Decode(decode DecodeBodyFn) error {
	// Make sure there is an hcl.Body to decode. The root module does not have
	// a body, for example, because it is constructed manually and not through
	// HCL
	if m.SourceHCL == nil {
		return nil
	}
	// Behaviour based on the module type
	switch t := moduleType(m.Type); t {
	case localmoduleType:
		// If local type then we need to get the location of this ModuleBlock
		// definition so that we can retrieve the actual module based on this
		// relative path
		mFile := m.SourceHCL.Body.MissingItemRange().Filename
		if mFile == "" {
			panic(fmt.Sprintf(`Error: could not get location of HCL Body for module "%s"`, m.Name))
		}
		mDir := filepath.Dir(mFile)
		m.Source = &ModuleLocalSource{
			// set the base directory for the ModuleLocalSource to be the
			// directory containing the block of hcl that contains the
			// definition of this module... from which the "path" attribute
			// is the relative path from
			BaseDir: mDir,
		}
	case remotemoduleType:
		// TODO
		panic(fmt.Sprintf(`Module type "%s" not implemented yet...`, t))
	default:
		panic(fmt.Sprintf(`Invalid module type "%s"`, t))
	}

	// Finally decode the source{} block into the specific type of Source set
	// above
	return decode(m.SourceHCL.Body, m.Source)
}

// Body returns the Body of HCL that is referenced by the module source.
// If the module is a local module, then this is just parsing the HCL files
// and merging the hcl.Body together.
func (m *ModuleBlock) Body() (hcl.Body, error) {
	return m.Source.Body()
}

// moduleType represents the different types of modules
type moduleType string

const (
	localmoduleType  moduleType = "local"
	remotemoduleType            = "remote"
)

// moduleSource creates an interface for the different module source blocks
// to implement.
type moduleSource interface {
	// Body returns the Body of HCL referenced by the source block
	Body() (hcl.Body, error)
}

// Ensure that ModuleLocalSource implements the moduleSource interface
var _ moduleSource = &ModuleLocalSource{}

// ModuleLocalSource represents the source block of local module types
type ModuleLocalSource struct {
	BaseDir string
	Path    string `hcl:"path,attr"`
}

// Body returns the hcl.Body for local modules
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

// moduleInput represent the input{} block in a module
type moduleInput struct {
	Name  string    `hcl:",label"`
	Value cty.Value `hcl:"value,attr"`
}
