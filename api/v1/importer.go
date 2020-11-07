package v1

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"reflect"
	"strings"

	"github.com/rs/zerolog/log"

	"github.com/clbanning/mxj"
	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/hashicorp/hcl/v2"
	"github.com/verifa/bubbly/api/core"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/gocty"
)

// Compiler check to see that v1.Importer implements the Importer interface
var _ core.Importer = (*Importer)(nil)

// Importer represents an importer type
type Importer struct {
	*core.ResourceBlock

	Spec importerSpec `json:"spec"`
}

// NewImporter returns a new Importer
func NewImporter(resBlock *core.ResourceBlock) *Importer {
	return &Importer{
		ResourceBlock: resBlock,
	}
}

// Apply returns the output from applying a resource
func (i *Importer) Apply(ctx *core.ResourceContext) core.ResourceOutput {
	if err := i.decode(ctx.DecodeBody); err != nil {
		return core.ResourceOutput{
			Status: core.ResourceOutputFailure,
			Error:  fmt.Errorf("Failed to decode resource %s: %w", i.String(), err),
		}
	}

	if i == nil {
		return core.ResourceOutput{
			Status: core.ResourceOutputFailure,
			Error:  errors.New("Cannot get output of a null importer"),
			Value:  cty.NilVal,
		}
	}

	if i.Spec.Source == nil {
		return core.ResourceOutput{
			Status: core.ResourceOutputFailure,
			Error:  errors.New("Cannot get output of an importer with null source"),
			Value:  cty.NilVal,
		}
	}

	val, err := i.Spec.Source.Resolve()
	if err != nil {
		return core.ResourceOutput{
			Status: core.ResourceOutputFailure,
			Error:  fmt.Errorf("Failed to resolve importer source: %w", err),
			Value:  cty.NilVal,
		}
	}

	return core.ResourceOutput{
		Status: core.ResourceOutputSuccess,
		Error:  nil,
		Value:  val,
	}
}

// SpecValue method returns resource specification structure
func (i *Importer) SpecValue() core.ResourceSpec {
	return &i.Spec
}

// decode is responsible for decoding any necessary hcl.Body inside Importer
func (i *Importer) decode(decode core.DecodeBodyFn) error {
	// decode the resource spec into the importer's Spec
	if err := decode(i, i.SpecHCL.Body, &i.Spec); err != nil {
		return fmt.Errorf(`Failed to decode "%s" body spec: %w`, i.String(), err)
	}

	// based on the type of the importer, initiate the importer's Source
	switch i.Spec.Type {
	case jsonImporterType:
		i.Spec.Source = &jsonSource{}
	case xmlImporterType:
		i.Spec.Source = &xmlSource{}
	default:
		panic(fmt.Sprintf("Unsupported importer resource type %s", i.Spec.Type))
	}

	// decode the source HCL into the importer's Source
	if err := decode(i, i.Spec.SourceHCL.Body, i.Spec.Source); err != nil {
		return fmt.Errorf(`Failed to decode importer source: %w`, err)
	}

	return nil
}

var _ core.ResourceSpec = (*importerSpec)(nil)

// importerSpec defines the spec for an importer
type importerSpec struct {
	Inputs InputDeclarations `hcl:"input,block"`
	// the type is either json, xml, rest, etc.
	Type      importerType `hcl:"type,attr"`
	SourceHCL *struct {
		Body hcl.Body `hcl:",remain"`
	} `hcl:"source,block"`
	// Source stores the actual value for SourceHCL
	Source source
}

// importerType defines the type of an importer
type importerType string

const (
	jsonImporterType importerType = "json"
	xmlImporterType               = "xml"
	gitImporterType               = "git"
)

// Source is an interface for the different data sources that an Importer can have
type source interface {
	// returns an interface{} containing the parsed XML, JSON data, that should
	// be converted into the Output cty.Value
	Resolve() (cty.Value, error)
}

// Compiler check to see that the source interface is implemented
var _ source = (*gitSource)(nil)

// gitSource represents the importer type for a local git repository data
type gitSource struct {
	Directory string `hcl:"directory,attr"`
}

// Resolve returns a cty.Value representation of the data from a local Git repo
func (s *gitSource) Resolve() (cty.Value, error) {

	//zerolog.SetGlobalLevel(zerolog.DebugLevel)
	//log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, NoColor: true})

	// The format of v1 Git importer output
	format := cty.Object(map[string]cty.Type{
		"is_bare":       cty.Bool,
		"commit_id":     cty.String,
		"tag":           cty.String,
		"active_branch": cty.String,
		"branches": cty.Object(map[string]cty.Type{
			"local":  cty.List(cty.String),
			"remote": cty.List(cty.String),
		}),
		"remotes": cty.List(cty.Object(map[string]cty.Type{
			"name": cty.String,
			"url":  cty.String,
		})),
	})

	// Find and open the repo
	repo, err := git.PlainOpen(s.Directory)

	if err != nil {
		return cty.NilVal, fmt.Errorf(`cannot open repository %s, error %w`, s.Directory, err)
	}

	// Is the repo bare or not
	cfg, err := repo.Config()
	if err != nil {
		return cty.NilVal, fmt.Errorf(`failed to check repo status (bare or not) for repo %s, error %w`, s.Directory, err)
	}
	isBare := cfg.Core.IsBare

	// Find HEAD and establish whether it's pointing to a proper branch
	headRef, err := repo.Head()
	if err != nil {
		return cty.NilVal, fmt.Errorf(`failed to read the repo (%s) HEAD, error %w`, s.Directory, err)
	}

	var headBranch string

	if headRef.Name().IsBranch() {
		headBranch = headRef.Name().Short()
	} else {
		headBranch = `Detached HEAD`
	}

	// Local branches: iterate to extract short names
	var localBranchNames []string
	branches, err := repo.Branches()
	if err != nil {
		return cty.NilVal, fmt.Errorf(`failed to get a list of local branches for repo %s, error %w`, s.Directory, err)
	}

	err = branches.ForEach(func(ref *plumbing.Reference) error {
		log.Debug().Msgf(`Local branch: %v`, ref.Name().Short())
		localBranchNames = append(localBranchNames, ref.Name().Short())
		return nil
	})

	if err != nil {
		return cty.NilVal, fmt.Errorf(`failed to iterate over a list of local branches for repo %s, error %w`, s.Directory, err)
	}

	// Remotes
	remotesList, err := repo.Remotes()
	if err != nil {
		return cty.NilVal, fmt.Errorf(`failed to get a list of remotes for repo %s, error %w`, s.Directory, err)
	}

	var remotes = make([]map[string]string, len(remotesList))

	for i, remote := range remotesList {
		remotes[i] = map[string]string{
			"name": remote.Config().Name,
			"url":  remote.Config().URLs[0], // always non-empty; first elem is for `git fetch`
		}
	}

	// Remote branches
	var remoteBranchNames []string
	refs, _ := repo.References()
	err = refs.ForEach(func(ref *plumbing.Reference) error {

		if ref.Type() == plumbing.HashReference && ref.Name().IsRemote() {
			remoteBranchNames = append(remoteBranchNames, ref.Name().Short())
		} else {
			log.Debug().Str("ref.String()", ref.String()).Msg(`Reference not a remote:`)
		}
		return nil
	})
	if err != nil {
		return cty.NilVal, fmt.Errorf(`failed to compile a list of known remote branches for repo %s, error %w`, s.Directory, err)
	}

	/*
		// The config file and the data structure representing it would only have those branches which have
		// upstream tracking set up.
		cfg, _ := repo.Config()
		for _, branch := range cfg.Branches {
			log.Debug().Str("branch.Name", branch.Name).Str("branch.Remote", branch.Remote).Msg("Branch read from config")
		}
	*/

	// Tags
	var tag string

	tagrefs, err := repo.Tags()
	if err != nil {
		return cty.NilVal, fmt.Errorf(`failed to read tags from repo %s, error %w`, s.Directory, err)
	}
	err = tagrefs.ForEach(func(t *plumbing.Reference) error {
		log.Debug().Str("short name", t.Name().Short()).Str("hash", t.Hash().String()).Msg(`Found tag:`)
		if t.Hash() == headRef.Hash() {
			tag = t.Name().Short()
		}
		return nil
	})
	if err != nil {
		return cty.NilVal, fmt.Errorf(`failed to iterate over tags from repo %s, error %w`, s.Directory, err)
	}

	// Construct Go data structure for conversion
	// to cty.Value using well-defined cty.Type
	data := map[string]interface{}{
		"is_bare":       isBare,
		"commit_id":     headRef.Hash().String(),
		"tag":           tag,
		"active_branch": headBranch,
		"branches": map[string][]string{
			"local":  localBranchNames,
			"remote": remoteBranchNames,
		},
		"remotes": remotes,
	}

	val, err := gocty.ToCtyValue(data, format)
	if err != nil {
		return cty.NilVal, fmt.Errorf(`failed to tranform the data for output, repo %s, error %w`, s.Directory, err)
	}

	return val, nil
}

// Compiler check to see that v1.JSONSource implements the Source interface
var _ source = (*jsonSource)(nil)

// jsonSource represents the importer type for using a JSON file as the input
type jsonSource struct {
	File string `hcl:"file,attr"`
	// the format of the raw input data defined as a cty.Type
	Format cty.Type `hcl:"format,attr"`
}

// Resolve returns a cty.Value representation of the parsed JSON file
func (s *jsonSource) Resolve() (cty.Value, error) {

	var barr []byte
	var err error

	// FIXME GitHub issue #39
	barr, err = ioutil.ReadFile(s.File)
	if err != nil {
		return cty.NilVal, err
	}

	// Attempt to unmarshall the data into an empty interface data type
	var data interface{}
	err = json.Unmarshal(barr, &data)
	if err != nil {
		return cty.NilVal, err
	}

	val, err := gocty.ToCtyValue(data, s.Format)
	if err != nil {
		return cty.NilVal, err
	}

	return val, nil
}

// Compiler check to see that v1.XMLSource implements the Source interface
var _ source = (*xmlSource)(nil)

// xmlSource represents the importer type for using an XML file as the input
type xmlSource struct {
	File string `hcl:"file,attr"`
	// the format of the raw input data defined as a cty.Type
	Format cty.Type `hcl:"format,attr"`
}

// Resolve returns a cty.Value representation of the XML file
func (s *xmlSource) Resolve() (cty.Value, error) {

	var barr []byte
	var err error

	// FIXME GitHub issue #39
	barr, err = ioutil.ReadFile(s.File)
	if err != nil {
		return cty.NilVal, err
	}

	mxj.PrependAttrWithHyphen(false) // no "-" prefix on attributes
	mxj.CastNanInf(true)             // use float64, not string for extremes

	// Unmarshall the XML data into a Go object
	data, err := mxj.NewMapXml(barr, true)
	if err != nil {
		return cty.NilVal, err
	}

	if err := walkTypeTransformData(&data, s.Format); err != nil {
		return cty.NilVal, err
	}

	val, err := gocty.ToCtyValue(data, s.Format)
	if err != nil {
		return cty.NilVal, err
	}

	return val, nil
}

func walkTypeTransformData(data *mxj.Map, ty cty.Type) error {
	path := make([]string, 0)
	return walk(data, ty, path, 0)
}

func walk(data *mxj.Map, ty cty.Type, path []string, idx int) error {

	pathStr := strings.Join(path, ".")

	if idx > 0 {
		pathStr += fmt.Sprint("[", idx, "]")
	}

	if ty.IsObjectType() {
		for x := range ty.AttributeTypes() {
			path = append(path, x)
			pathIdx := len(path) - 1

			walk(data, ty.AttributeType(x), path, 0)
			path = path[0:pathIdx]
		}
	}

	if ty.IsListType() {

		vs, err := data.ValuesForPath(pathStr)
		if err != nil {
			return fmt.Errorf("wrong path (%s) in xml structure: %w", pathStr, err)
		}

		n := len(vs)
		//t.Logf("ValuesForPath(%s): %d", pathStr, n)

		switch n {
		case 0:
			return fmt.Errorf("xml data structure inconsistent state, ValuesForPath are zero at %s", pathStr)
		case 1:
			v := vs[0]

			if reflect.TypeOf(v).Kind() == reflect.Map {
				vv := make([]interface{}, 0)
				vv = append(vv, v)
				if err := data.SetValueForPath(vv, pathStr); err != nil {
					return fmt.Errorf("cannot convert at path %s, error %w", pathStr, err)
				}
			}
			fallthrough
		default:
			for i := range vs {
				return walk(data, ty.ElementType(), path, i)
			}
		}
	}

	return nil
}
