package v1

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"path/filepath"
	"reflect"
	"strings"
	"time"

	"github.com/clbanning/mxj"
	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/hashicorp/hcl/v2"
	"github.com/verifa/bubbly/api/core"
	"github.com/verifa/bubbly/env"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/gocty"
)

// Compiler check to see that v1.Extract implements the Extract interface
var _ core.Extract = (*Extract)(nil)

// Extract represents an extract type
type Extract struct {
	*core.ResourceBlock

	Spec extractSpec `json:"spec"`
}

// NewExtract returns a new Extract
func NewExtract(resBlock *core.ResourceBlock) *Extract {
	return &Extract{
		ResourceBlock: resBlock,
	}
}

// Apply returns the output from applying a resource
func (i *Extract) Apply(bCtx *env.BubblyContext, ctx *core.ResourceContext) core.ResourceOutput {
	if err := i.decode(bCtx, ctx.DecodeBody); err != nil {
		return core.ResourceOutput{
			Status: core.ResourceOutputFailure,
			Error:  fmt.Errorf("Failed to decode resource %s: %w", i.String(), err),
		}
	}

	if i == nil {
		return core.ResourceOutput{
			Status: core.ResourceOutputFailure,
			Error:  errors.New("Cannot get output of a null extract"),
			Value:  cty.NilVal,
		}
	}

	if i.Spec.Source == nil {
		return core.ResourceOutput{
			Status: core.ResourceOutputFailure,
			Error:  errors.New("Cannot get output of an extract with null source"),
			Value:  cty.NilVal,
		}
	}

	val, err := i.Spec.Source.Resolve(bCtx)
	if err != nil {
		return core.ResourceOutput{
			Status: core.ResourceOutputFailure,
			Error:  fmt.Errorf("Failed to resolve extract source: %w", err),
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
func (i *Extract) SpecValue() core.ResourceSpec {
	return &i.Spec
}

// decode is responsible for decoding any necessary hcl.Body inside Extract
func (i *Extract) decode(bCtx *env.BubblyContext, decode core.DecodeBodyFn) error {
	// decode the resource spec into the extract's Spec
	if err := decode(bCtx, i, i.SpecHCL.Body, &i.Spec); err != nil {
		return fmt.Errorf(`Failed to decode "%s" body spec: %w`, i.String(), err)
	}

	// based on the type of the extract, initiate the extract's Source
	switch i.Spec.Type {
	case jsonExtractType:
		i.Spec.Source = &jsonSource{}
	case xmlExtractType:
		i.Spec.Source = &xmlSource{}
	case gitExtractType:
		i.Spec.Source = &gitSource{}
	case restExtractType:
		i.Spec.Source = &restSource{} // FIXME this is where mergo init happens
	default:
		panic(fmt.Sprintf("Unsupported extract resource type %s", i.Spec.Type))
	}

	// decode the source HCL into the extract's Source
	if err := decode(bCtx, i, i.Spec.SourceHCL.Body, i.Spec.Source); err != nil {
		return fmt.Errorf(`Failed to decode extract source: %w`, err)
	}

	return nil
}

var _ core.ResourceSpec = (*extractSpec)(nil)

// extractSpec defines the spec for an extract
type extractSpec struct {
	Inputs InputDeclarations `hcl:"input,block"`
	// the type is either json, xml, rest, etc.
	Type      extractType `hcl:"type,attr"`
	SourceHCL *struct {
		Body hcl.Body `hcl:",remain"`
	} `hcl:"source,block"`
	// Source stores the actual value for SourceHCL
	Source source
}

// extractType defines the type of an extract
type extractType string

const (
	jsonExtractType extractType = "json"
	xmlExtractType              = "xml"
	gitExtractType              = "git"
	restExtractType             = "rest"
)

// Source is an interface for the different data sources that an Extract can have
type source interface {
	// returns an interface{} containing the parsed XML, JSON data, that should
	// be converted into the Output cty.Value
	Resolve(*env.BubblyContext) (cty.Value, error)
}

// Compiler check to see that the source interface is implemented
var _ source = (*restSource)(nil)

type basicAuth struct {
	Username     string `hcl:"username,attr"`
	Password     string `hcl:"password,attr"`
	PasswordFile string `hcl:"password_file,attr"`
}

// restSource represents the extract type for a REST API query
type restSource struct {

	// Protocol is "http" or "https"
	Protocol string `hcl:"protocol,attr"`

	// Host is ?
	Host string `hcl:"host,attr"`

	// Port is TCP port number
	Port uint16 `hcl:"port,attr"`

	// Route is ?
	Route string `hcl:"path,attr"`

	// Method is "GET" or "POST" for protocols "http" and "https"
	Method string `hcl:"method,attr"`

	// Params are URL query params which get prepended at the end of the url
	// in the format ?key1=val1&key2=val2
	Params map[string]string `hcl:"params,attr"`

	// HTTP Headers to set, anything apart from Basic Authorization
	// and Bearer Token for which there are custom fields.
	Headers map[string]string `hcl:"headers,attr"`

	// BasicAuth sets the `Authorization` header on every
	// scrape request with the configured username and password.
	// The `password` field, if set, overrides the `password_file` field.
	// More information: https://swagger.io/docs/specification/authentication/basic-authentication/
	BasicAuth basicAuth `hcl:"basic_auth,attr"`

	// BearerToken sets the `Authorization` header on every
	// scrape request with the configured bearer token.
	// More information: https://swagger.io/docs/specification/authentication/bearer-authentication/
	// If set, this option overrides `bearer_token_file`.
	BearerToken string `hcl:"bearer_token,attr"`

	// BearerTokenFile sets the `Authorization` header on every
	// scrape request with the bearer token read from the configured file.
	// This option can be overridden by `bearer_token`.
	BearerTokenFile string `hcl:"bearer_token_file"`

	// Flavour is "json". This is the expected format of the response.
	Flavour string `hcl:"flavour,attr"`

	// Timeout in seconds is how long the extractor can wait before giving up
	// trying to extract the data from this resource.
	Timeout uint `hcl:"timeout,attr"`

	// Format is ?
	Format cty.Type `hcl:"format,attr"`
}

// Resolve returns a cty.Value representation of the data in response to a REST API query
func (s *restSource) Resolve(bCtx *env.BubblyContext) (cty.Value, error) {

	// HTTP Method:
	//   * "GET" or "POST", case-insensitive
	//   * the default value "" is interpreted as GET
	//   * any other value is invalid and raises error
	//
	method := strings.ToUpper(s.Method)
	switch method {
	case http.MethodGet:
		break
	case http.MethodPost:
		return cty.NilVal, fmt.Errorf(`not implemented: method: %s`, method)
	case "":
		method = http.MethodGet
	default:
		return cty.NilVal, fmt.Errorf(`unsupported method: %s`, method)
	}

	// Transport Protocol (aka scheme):
	//   * "http" or "https" only, case-insensitive
	//   * the default value "" is interpreted as "http"
	//   * any other value is invalid and raises error
	//
	scheme := strings.ToLower(s.Protocol)
	switch scheme {
	case "http":
		break
	case "https":
		break
	case "":
		scheme = "http"
	default:
		return cty.NilVal, fmt.Errorf(`unsupported protocol: %s`, s.Protocol)
	}

	// HTTP request timeout:
	//   * a positive integer is accepted and interpreted as value in seconds
	//   * the default value 0 is interpreted as 1 second
	//
	var timeout time.Duration
	switch {
	case s.Timeout == 0:
		timeout = time.Second
	default:
		timeout = time.Duration(s.Timeout) * time.Second
	}

	// Basic Authentication activates only if the username is provided
	//    * only on of `password` or `password_file` must be provided; the other must be empty string
	//    * no additional encoding is done on the values of `username` and `password`
	//
	var username, password string

	if s.BasicAuth.Username != "" {

		username = s.BasicAuth.Username

		switch {
		case s.BasicAuth.Password != "" && s.BasicAuth.PasswordFile != "":
			return cty.NilVal, fmt.Errorf("for basic auth, only one of password or password_file must be provided with the username")

		case s.BasicAuth.Password != "":
			password = s.BasicAuth.Password

		case s.BasicAuth.PasswordFile != "":
			byteArr, err := ioutil.ReadFile(filepath.FromSlash(s.BasicAuth.PasswordFile))
			if err != nil {
				return cty.NilVal, fmt.Errorf("failed to read the password file: %w", err)
			}
			password = string(byteArr)
		}
	}

	// Construct a URL string
	us := scheme + "://" + s.Host

	if s.Port != 0 {
		us += fmt.Sprint(":", s.Port)
	}

	if strings.HasPrefix(s.Route, "/") {
		us += s.Route
	} else {
		us += fmt.Sprint("/", s.Route)
	}

	// URL query string
	params := url.Values{}
	for k, v := range s.Params {
		params.Add(k, v)
	}
	if len(params) > 0 {
		us += fmt.Sprint("?", params.Encode())
	}
	// The URL has been built.

	// We don't **need** a value of type url.URL but it's a useful check
	u, err := url.Parse(us)
	if err != nil {
		return cty.NilVal, fmt.Errorf(`failed to parse endpoint url %s: %w`, us, err)
	}

	// Create an object representing a HTTP request
	var body io.Reader
	httpRequest, err := http.NewRequest(method, u.String(), body)
	if err != nil {
		return cty.NilVal, fmt.Errorf(`failed to craft HTTP request object: %w`, err)
	}

	// Authentication, if requested
	if username != "" {
		httpRequest.SetBasicAuth(s.BasicAuth.Username, password)
	}

	// Add a bearer token, if requests. Adds a Header.
	switch {
	case s.BearerToken != "":
		httpRequest.Header.Set("Authorization", "Bearer "+s.BearerToken)
	case s.BearerTokenFile != "":
		bearerToken, err := ioutil.ReadFile(filepath.FromSlash(s.BearerTokenFile))
		if err != nil {
			return cty.NilVal, fmt.Errorf("failed to read bearer token file: %w", err)
		}
		httpRequest.Header.Set("Authorization", "Bearer "+string(bearerToken))
	}

	// Any other headers, if reqested
	for k, v := range s.Headers {
		httpRequest.Header.Add(k, v)
	}

	// Initiate the HTTP client
	c := http.Client{Timeout: timeout}

	// Make REST API request
	httpResponse, err := c.Do(httpRequest)
	if err != nil {
		return cty.NilVal, fmt.Errorf(`failed to make HTTP request: %w`, err)
	}

	if httpResponse.StatusCode != http.StatusOK {
		return cty.NilVal, fmt.Errorf(`HTTP response status code: %d`, httpResponse.StatusCode)
	}

	defer httpResponse.Body.Close()

	// Parse the content of response body
	var data interface{}
	if err := json.NewDecoder(httpResponse.Body).Decode(&data); err != nil {
		return cty.NilVal, fmt.Errorf(`failed to decode JSON: %w`, err)
	}

	val, err := gocty.ToCtyValue(data, s.Format)
	if err != nil {
		return cty.NilVal, fmt.Errorf(`failed to convert to desired data format: %w`, err)
	}

	return val, nil
}

// Compiler check to see that the source interface is implemented
var _ source = (*gitSource)(nil)

// gitSource represents the extract type for a local git repository data
type gitSource struct {
	Directory string `hcl:"directory,attr"`
}

// Resolve returns a cty.Value representation of the data from a local Git repo
func (s *gitSource) Resolve(bCtx *env.BubblyContext) (cty.Value, error) {

	// The format of v1 Git extract output
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
		bCtx.Logger.Debug().Msgf(`Local branch: %v`, ref.Name().Short())
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
			bCtx.Logger.Debug().Str("ref.String()", ref.String()).Msg(`Reference not a remote:`)
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
			bCtx.Logger.Debug().Str("branch.Name", branch.Name).Str("branch.Remote", branch.Remote).Msg("Branch read from config")
		}
	*/

	// Tags
	var tag string

	tagrefs, err := repo.Tags()
	if err != nil {
		return cty.NilVal, fmt.Errorf(`failed to read tags from repo %s, error %w`, s.Directory, err)
	}
	err = tagrefs.ForEach(func(t *plumbing.Reference) error {
		bCtx.Logger.Debug().Str("short name", t.Name().Short()).Str("hash", t.Hash().String()).Msg(`Found tag:`)
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

// jsonSource represents the extract type for using a JSON file as the input
type jsonSource struct {
	File string `hcl:"file,attr"`
	// the format of the raw input data defined as a cty.Type
	Format cty.Type `hcl:"format,attr"`
}

// Resolve returns a cty.Value representation of the parsed JSON file
func (s *jsonSource) Resolve(bCtx *env.BubblyContext) (cty.Value, error) {

	var barr []byte
	var err error

	// FIXME GitHub issue #39
	barr, err = ioutil.ReadFile(s.File)
	if err != nil {
		return cty.NilVal, fmt.Errorf("failed to read file %s: %w", s.File, err)
	}

	// Attempt to unmarshall the data into an empty interface data type
	var data interface{}
	err = json.Unmarshal(barr, &data)
	if err != nil {
		return cty.NilVal, fmt.Errorf("failed to marshal JSON: %w", err)
	}

	val, err := gocty.ToCtyValue(data, s.Format)
	if err != nil {
		return cty.NilVal, err
	}

	return val, nil
}

// Compiler check to see that v1.XMLSource implements the Source interface
var _ source = (*xmlSource)(nil)

// xmlSource represents the extract type for using an XML file as the input
type xmlSource struct {
	File string `hcl:"file,attr"`
	// the format of the raw input data defined as a cty.Type
	Format cty.Type `hcl:"format,attr"`
}

// Resolve returns a cty.Value representation of the XML file
func (s *xmlSource) Resolve(bCtx *env.BubblyContext) (cty.Value, error) {

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
