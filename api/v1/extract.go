package v1

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"time"

	"github.com/imdario/mergo"

	"github.com/clbanning/mxj"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/ext/typeexpr"

	"github.com/valocode/bubbly/api/common"
	"github.com/valocode/bubbly/api/core"
	"github.com/valocode/bubbly/env"
	"github.com/valocode/bubbly/events"

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

// Run returns the output from applying a resource
func (e *Extract) Run(bCtx *env.BubblyContext, ctx *core.ResourceContext) core.ResourceOutput {

	if err := e.decode(bCtx, ctx); err != nil {
		return core.ResourceOutput{
			ID:     e.String(),
			Status: events.ResourceRunFailure,
			Error:  fmt.Errorf("failed to decode resource %s: %w", e.String(), err),
		}
	}

	if e == nil {
		return core.ResourceOutput{
			Status: events.ResourceRunFailure,
			Error:  errors.New("cannot get output of a null extract"),
			Value:  cty.NilVal,
		}
	}

	if e.Spec.Source == nil {
		return core.ResourceOutput{
			Status: events.ResourceRunFailure,
			Error:  errors.New("cannot get output of an extract with null source"),
			Value:  cty.NilVal,
		}
	}

	if len(e.Spec.Source) == 0 {
		return core.ResourceOutput{
			Status: events.ResourceRunFailure,
			Error:  errors.New("cannot get output of an extract with no source"),
			Value:  cty.NilVal,
		}
	}

	var (
		vals = make([]cty.Value, 0, len(e.Spec.Source))
		errs = make([]core.ResourceOutput, 0)
	)
	for _, src := range e.Spec.Source {
		val, err := src.Resolve(bCtx)
		if err != nil {
			output := core.ResourceOutput{
				ID:     e.String(),
				Status: events.ResourceRunFailure,
				Error:  fmt.Errorf("error on extract: %w", err),
				Value:  cty.NilVal,
			}
			if e.Spec.ContinueOnError {
				errs = append(errs, output)
				// Append a value anyway, in case indexes are used to access
				// the results, this keeps the extract data "in sync"
				vals = append(vals, cty.NilVal)
				bCtx.Logger.Error().Err(err).Msg("error on extract but continue on error set")
				continue
			}
			return output
		}
		vals = append(vals, val)
	}

	if len(errs) == len(e.Spec.Source) {
		return core.ResourceOutput{
			ID:     e.ID(),
			Status: events.ResourceRunFailure,
			Error:  errors.New("error on all extracts"),
			Value:  cty.NilVal,
		}
	}

	// We support dynamic blocks (HCL extension), so there may be more than one
	// effective source block in the resource spec (`resource/spec/source[]`).
	// When there is only one, the data structure representing that part of
	// the spec is still a slice and in that case it has a length of one.
	var val cty.Value
	switch len(e.Spec.Source) {
	case 0:
		return core.ResourceOutput{
			ID:     e.String(),
			Status: events.ResourceRunFailure,
			Error:  fmt.Errorf("failed to resolve extract source: no sources defined"),
			Value:  cty.NilVal,
		}
	case 1:
		val = vals[0]
	default:
		val = cty.TupleVal(vals)
	}

	// There is the exception if multi source is set
	if e.Spec.MultiSource {
		val = cty.TupleVal(vals)
	}

	return core.ResourceOutput{
		ID:     e.String(),
		Status: events.ResourceRunSuccess,
		Error:  nil,
		Value:  val,
	}
}

// set{*}SourceDefaults are the initialisers for some types of Source(s),
// where Golang default values would not be sufficient. Their purpose is
// to simplify the Resolve(ers) logic by avoiding checks on null or default
// values where such checks would have been inconvenient or verbose.
//
// It was a deliberate design decision that the "correct" or "default" values for
// optional fields are set AFTER the HCL parser has created and populated
// the graphQLSource structure.
//
// TODO: why though? i don't remember anymore. it had something to do with HCL parser, maybe?
//

// setGraphQLSourceDefaults is the initialiser for the GraphQL Source type
func setGraphQLSourceDefaults(bCtx *env.BubblyContext, dst *graphqlSource) error {
	method := http.MethodPost
	timeout := uint(1)

	defaults := &graphqlSource{
		Method:  method,
		Params:  map[string]string{},
		Headers: map[string]string{},
		Timeout: timeout,
	}

	if err := mergo.Merge(dst, defaults); err != nil {
		return fmt.Errorf("failed to initialise the data structure: %w", err)
	}

	// Mergo does not set empty and nil values, and for purposes of
	// simplifying control flow logic it helps to set the following...

	// Params is optional in HCL, but it helps to have an empty map instead of nil.
	if dst.Params == nil {
		dst.Params = map[string]string{}
	}
	// Headers is optional in HCL
	if dst.Headers == nil {
		dst.Headers = map[string]string{}
	}

	// If HTTP Basic Authentication is requested in HCL, the username is compulsory field,
	// while the password and password_file are both optional since one or the other may be
	// provided by the user. If the username is given, initialise the other two fields, unless
	// already set.
	if dst.BasicAuth != nil {
		if dst.BasicAuth.Password == nil {
			dst.BasicAuth.Password = new(string)
		}
		if dst.BasicAuth.PasswordFile == nil {
			dst.BasicAuth.PasswordFile = new(string)
		}
	}

	return nil
}

// setRestSourceDefaults is the initialiser for the REST Source type
func setRestSourceDefaults(bCtx *env.BubblyContext, dst *restSource) error {

	method := http.MethodGet
	timeout := uint(1)

	defaults := &restSource{
		Method:  method,
		Params:  &map[string]string{},
		Headers: &map[string]string{},
		Timeout: &timeout,
	}

	if err := mergo.Merge(dst, defaults); err != nil {
		return fmt.Errorf("failed to initialise the data structure: %w", err)
	}

	// Mergo does not set empty and nil values, and for purposes of
	// simplifying control flow logic it helps to set the following...

	// Params is optional in HCL, but it helps to have an empty map instead of nil.
	if dst.Params == nil {
		dst.Params = &map[string]string{}
	}
	// Headers is optional in HCL
	if dst.Headers == nil {
		dst.Headers = &map[string]string{}
	}

	// If HTTP Basic Authentication is requested in HCL, the username is compulsory field,
	// while the password and password_file are both optional since one or the other may be
	// provided by the user. If the username is given, initialise the other two fields, unless
	// already set.
	if dst.BasicAuth != nil {
		if dst.BasicAuth.Password == nil {
			dst.BasicAuth.Password = new(string)
		}
		if dst.BasicAuth.PasswordFile == nil {
			dst.BasicAuth.PasswordFile = new(string)
		}
	}

	return nil
}

// decode is responsible for decoding any necessary hcl.Body inside Extract
func (e *Extract) decode(bCtx *env.BubblyContext, ctx *core.ResourceContext) error {
	if err := common.DecodeBodyWithInputs(bCtx, e.SpecHCL.Body, &e.Spec, ctx); err != nil {
		return err
	}

	e.Spec.Source = make(SourceBlocks, len(e.Spec.SourceHCL))

	for idx, source := range e.Spec.SourceHCL {

		// Initiate the Extract's Source structure
		switch e.Spec.Type {
		case jsonExtractType:
			e.Spec.Source[idx] = new(jsonSource)
		case xmlExtractType:
			e.Spec.Source[idx] = new(xmlSource)
		case gitExtractType:
			e.Spec.Source[idx] = new(gitSource)
		case restExtractType:
			e.Spec.Source[idx] = new(restSource)
		case graphQLExtractType:
			e.Spec.Source[idx] = new(graphqlSource)
		default:
			return fmt.Errorf("unsupported extract resource type: %s", e.Spec.Type)
		}

		// decode the source HCL into the extract's Source
		if err := common.DecodeBody(bCtx, source.Body, e.Spec.Source[idx], ctx); err != nil {
			return fmt.Errorf("failed to decode extract source: %w", err)
		}

		// Merge with default values for those source types where it's necessary
		switch dst := e.Spec.Source[idx].(type) {
		case *restSource:
			if err := setRestSourceDefaults(bCtx, dst); err != nil {
				return fmt.Errorf("failed to decode extract: %w", err)
			}
		case *graphqlSource:
			if err := setGraphQLSourceDefaults(bCtx, dst); err != nil {
				return fmt.Errorf("failed to decode extract: %w", err)
			}
		}
	}
	return nil
}

// SourceBlocks stores the HCL for the `source` block in `extract` type resource
type SourceBlocks []source

// extractSpec defines the spec for an extract
type extractSpec struct {
	Inputs core.InputDeclarations `hcl:"input,block"`
	// the type is either json, xml, rest, etc.
	Type            extractType `hcl:"type,attr"`
	MultiSource     bool        `hcl:"multi_source,optional"`
	ContinueOnError bool        `hcl:"continue_on_error,optional"`
	SourceHCL       []struct {
		Body hcl.Body `hcl:",remain"`
	} `hcl:"source,block"`
	// Source stores the actual value for SourceHCL
	Source SourceBlocks
}

// extractType defines the type of an extract
type extractType string

const (
	jsonExtractType    extractType = "json"
	xmlExtractType     extractType = "xml"
	gitExtractType     extractType = "git"
	restExtractType    extractType = "rest"
	graphQLExtractType extractType = "graphql"
)

// Source is an interface for the different data sources that an Extract can have
type source interface {
	// Resolve requests data in a way specific to the dynamic type of the interface,
	// and handles the conversion of response to a dynamic value.
	Resolve(*env.BubblyContext) (cty.Value, error)
}

// Compiler check to see that the source interface is implemented
var _ source = (*graphqlSource)(nil)

// graphqlSource represents the extract type for a GraphQL API query
type graphqlSource struct {

	// The body of the GraphQL `query { ... }`
	Query string `hcl:"query"`

	// URL (technically a URI reference, as per RFC 3986) represents the unparsed URL string.
	// Note that URL Query Parameters (?key=val) must be provided separately as Params.
	URL string `hcl:"url"`

	// Method is "GET" or "POST" for protocols "http" and "https". The default is "POST".
	Method string `hcl:"method,optional"`

	// Params are URL query params which get prepended at the end of the url
	// in the format ?key1=val1&key2=val2
	Params map[string]string `hcl:"params,optional"`

	// HTTP Headers to set, anything apart from Basic Authorization
	// and Bearer Token for which there are custom fields.
	Headers map[string]string `hcl:"headers,optional"`

	// BasicAuth sets the `Authorization` header on every
	// scrape request with the configured username and password.
	// The `password` field, if set, overrides the `password_file` field.
	// More information: https://swagger.io/docs/specification/authentication/basic-authentication/
	BasicAuth *basicAuth `hcl:"basic_auth,block"`

	// BearerToken sets the `Authorization` header on every
	// scrape request with the configured bearer token.
	// More information: https://swagger.io/docs/specification/authentication/bearer-authentication/
	// If set, this option overrides `bearer_token_file`.
	BearerToken string `hcl:"bearer_token,optional"`

	// BearerTokenFile sets the `Authorization` header on every
	// scrape request with the bearer token read from the configured file.
	// This option can be overridden by `bearer_token`.
	BearerTokenFile string `hcl:"bearer_token_file,optional"`

	// Timeout in seconds is how long the extractor can wait before giving up
	// trying to extract the data from this resource.
	Timeout uint `hcl:"timeout,optional"`

	// Format is is a dynamic type, usually built from an HCL type expression.
	// It defines what is expected in response to the GraphQL API query.
	Format     cty.Type
	FormatExpr hcl.Expression `hcl:"format,attr"`
}

// Resolve performs a GraphQL query, parses the response, and returns a corresponding cty.Value
func (s *graphqlSource) Resolve(bCtx *env.BubblyContext) (cty.Value, error) {
	// Handle the format. First check if the current format is NilType (zero value)
	// and if so, get the type from FormatExpr
	if s.Format == cty.NilType {
		if s.FormatExpr == nil {
			return cty.NilVal, fmt.Errorf("no type available for format field")
		}
		var typeDiags hcl.Diagnostics
		s.Format, typeDiags = typeexpr.TypeConstraint(s.FormatExpr)
		if typeDiags.HasErrors() {
			return cty.NilVal, fmt.Errorf("invalid type for format field: %s", typeDiags.Error())
		}
	}

	// Wrap the query into proper JSON as this is what the GraphQL end-point expects
	query, err := json.Marshal(
		struct {
			Query string `json:"query"`
		}{
			Query: s.Query,
		},
	)

	if err != nil {
		return cty.NilVal, fmt.Errorf(`failed to marshal GraphQL query string into JSON: "%s"`, s.Query)
	}

	// HTTP method: "GET" or "POST", case-insensitive, default value "POST"
	method := strings.ToUpper(s.Method)

	if method != http.MethodPost && method != http.MethodGet {
		return cty.NilVal, fmt.Errorf("unsupported method: %s", method)
	}

	// HTTP request timeout:
	//   * a positive integer is accepted and interpreted as value in seconds
	//   * the default value 0 is interpreted as 1 second
	//
	var timeout time.Duration

	switch {
	case s.Timeout == 0:
		return cty.NilVal, fmt.Errorf("invalid timeout value: %d", s.Timeout)
	default:
		timeout = time.Duration(s.Timeout) * time.Second
	}

	// The 'Basic' HTTP Authentication Scheme as defined in RFC 7617.
	//
	// Activated by non-nil `basicAuth` field. When active,
	// * the `basicAuth.username` must be a non-empty string
	// * either `basicAuth.password` or `basicAuth.passwordFile` must be a non-empty string,
	//   with the former overriding the latter.

	var username, password string

	if s.BasicAuth != nil {

		if s.BasicAuth.Username == "" {
			return cty.NilVal, fmt.Errorf("HTTP basic authentication requires a username")
		}

		username = s.BasicAuth.Username

		switch {
		case *s.BasicAuth.Password != "":
			password = *s.BasicAuth.Password

		case *s.BasicAuth.PasswordFile != "":
			byteArr, err := os.ReadFile(filepath.FromSlash(*s.BasicAuth.PasswordFile))
			if err != nil {
				return cty.NilVal, fmt.Errorf("failed to read the password for http basic authentication from a file: %w", err)
			}
			password = string(byteArr)
		}
	}

	// URL query string
	params := url.Values{}
	for k, v := range s.Params {
		params.Add(k, v)
	}

	// Construct a URL string
	us := s.URL
	if len(params) > 0 {
		us = fmt.Sprint(s.URL, "?", params.Encode())
	}

	// Validate the URL string using the `url` standard library module
	u, err := url.Parse(us)
	if err != nil {
		return cty.NilVal, fmt.Errorf("failed to parse endpoint url %s: %w", us, err)
	}

	// Create an object representing a HTTP request
	body := bytes.NewReader(query)
	httpRequest, err := http.NewRequest(method, u.String(), body)
	if err != nil {
		return cty.NilVal, fmt.Errorf("failed to craft HTTP request object: %w", err)
	}

	// Authentication, if requested
	if s.BasicAuth != nil {
		httpRequest.SetBasicAuth(username, password)
	}

	// Add a bearer token, if requests. Adds a Header.
	var bearerToken string
	switch {
	case s.BearerToken != "":
		bearerToken = s.BearerToken
	case s.BearerTokenFile != "":
		bt, err := os.ReadFile(filepath.FromSlash(s.BearerTokenFile))
		if err != nil {
			return cty.NilVal, fmt.Errorf("failed to read bearer token file: %w", err)
		}
		bearerToken = strings.TrimSpace(string(bt))
	}
	if bearerToken != "" {
		httpRequest.Header.Set("Authorization", fmt.Sprint("Bearer ", bearerToken))
	}

	// Any other headers, if reqested
	for k, v := range s.Headers {
		httpRequest.Header.Add(k, v)
	}

	// Initiate the HTTP client
	c := http.Client{Timeout: timeout}

	// Make a request to GraphQL API end-point
	httpResponse, err := c.Do(httpRequest)
	if err != nil {
		return cty.NilVal, fmt.Errorf("failed to make HTTP request: %w", err)
	}

	if httpResponse.StatusCode != http.StatusOK {
		return cty.NilVal, fmt.Errorf("HTTP response status code: %d", httpResponse.StatusCode)
	}

	defer httpResponse.Body.Close()

	// Parse the content of response body into `interface{}` for further processing later
	var graphQLresponse interface{}
	if err := json.NewDecoder(httpResponse.Body).Decode(&graphQLresponse); err != nil {
		return cty.NilVal, fmt.Errorf("failed to decode GraphQL response: %w", err)
	}

	// As per GraphQL June 2018 spec, section 7.1:
	// A response to a GraphQL operation must be a map.
	// If the operation encountered any errors, the response map must contain an entry with key `errors`.
	// If the operation completed without encountering any errors, this entry must not be present.
	// https://spec.graphql.org/June2018/#sec-Response-Format

	// So, the strategy is to attempt to the `errors` part of the response first,
	// and if this fails, proceed to decoding the data.

	// Describes the GraphQL response `errors` field format.
	errfmt := cty.Object(map[string]cty.Type{
		"errors": cty.List(cty.Object(map[string]cty.Type{
			"message": cty.String,
			"locations": cty.List(cty.Object(map[string]cty.Type{
				"line":   cty.Number,
				"column": cty.Number,
			})),
			"path": cty.List(cty.String),
		})),
	})
	// TODO: not sure about `path` as it is a list of strings and numbers (mixed type)

	// TODO: the above format assumes that the GraphQL server returns a well-formed error message
	//       and that we don't fail to parse it for whatever reason. Maybe a more robust way would
	//       be first to check if `errors` exist, and if so, attempt to decode using a more detailed
	//       format. As it stands now, it's hard to distinguish between the cases when there is
	//       no `errors` field from the case when there is but we fail to parse it.

	// If there is an `error` field in response, the following operation succeeds,
	// and this function must return error, and should not process the `data` block
	errval, err := gocty.ToCtyValue(graphQLresponse, errfmt)
	// TODO: there got to be a better way to check this...
	if err == nil &&
		!errval.IsNull() &&
		!errval.GetAttr("errors").IsNull() &&
		errval.GetAttr("errors").Length().GreaterThan(cty.NumberIntVal(0)) == cty.True {
		// TODO: perhaps, if there is a locations/path in `error`, then `data` should be dumped as well?
		return cty.NilVal, fmt.Errorf("error in GraphQL response: %v", errval)
	}

	// If there was no error returned by GraphQL, the `data` part of the response can be processed
	// The GraphQL response is wrapped in `data` block but we don't require this in .bubbly `format` directive
	// to avoid being verbose. So, wrap the expected type of response with the `data` block.
	valfmt := cty.Object(map[string]cty.Type{
		"data": s.Format,
	})

	val, err := gocty.ToCtyValue(graphQLresponse, valfmt)
	if err != nil {
		return cty.NilVal, fmt.Errorf("failed to convert to desired data format: %w", err)
	}

	// Don't return the top-level element as it's part of GraphQL response format,
	// return only the actual data within it.
	return val.GetAttr("data"), nil
}

// Compiler check to see that the source interface is implemented
var _ source = (*restSource)(nil)

type basicAuth struct {
	Username     string  `hcl:"username"`
	Password     *string `hcl:"password"`
	PasswordFile *string `hcl:"password_file"`
}

// restSource represents the extract type for a REST API query
type restSource struct {

	// The body of a POST request. Ignored for a GET request.
	Query *string `hcl:"query"`

	// URL (technically a URI reference, as per RFC 3986) represents the unparsed URL string.
	// Note that URL Query Parameters (?key=val) must be provided separately as Params.
	URL string `hcl:"url"`

	// Method is "GET" or "POST" for protocols "http" and "https". The default is "GET".
	Method string `hcl:"method,optional"`

	// Params are URL query params which get prepended at the end of the url
	// in the format ?key1=val1&key2=val2
	Params *map[string]string `hcl:"params"`

	// HTTP Headers to set, anything apart from Basic Authorization
	// and Bearer Token for which there are custom fields.
	Headers *map[string]string `hcl:"headers"`

	// BasicAuth sets the `Authorization` header on every
	// scrape request with the configured username and password.
	// The `password` field, if set, overrides the `password_file` field.
	// More information: https://swagger.io/docs/specification/authentication/basic-authentication/
	BasicAuth *basicAuth `hcl:"basic_auth,block"`

	// BearerToken sets the `Authorization` header on every
	// scrape request with the configured bearer token.
	// More information: https://swagger.io/docs/specification/authentication/bearer-authentication/
	// If set, this option overrides `bearer_token_file`.
	BearerToken *string `hcl:"bearer_token"`

	// BearerTokenFile sets the `Authorization` header on every
	// scrape request with the bearer token read from the configured file.
	// This option can be overridden by `bearer_token`.
	BearerTokenFile *string `hcl:"bearer_token_file"`

	// Decoder is one of either "json" and "xml". It will be applied to extract's input.
	Decoder string `hcl:"decoder,optional"`

	// Timeout in seconds is how long the extractor can wait before giving up
	// trying to extract the data from this resource.
	Timeout *uint `hcl:"timeout"`

	// Format is a dynamic type, usually built from an HCL type expression.
	// It defines what is expected in response to the REST API query.
	Format     cty.Type
	FormatExpr hcl.Expression `hcl:"format,attr"`
}

// Resolve performs a REST query, parses the response, and returns a corresponding dynamic value.
func (s *restSource) Resolve(bCtx *env.BubblyContext) (cty.Value, error) {
	// Handle the format. First check if the current format is NilType (zero value)
	// and if so, get the type from FormatExpr
	if s.Format == cty.NilType {
		if s.FormatExpr == nil {
			return cty.NilVal, fmt.Errorf("no type available for format field")
		}
		var typeDiags hcl.Diagnostics
		s.Format, typeDiags = typeexpr.TypeConstraint(s.FormatExpr)
		if typeDiags.HasErrors() {
			return cty.NilVal, fmt.Errorf("invalid type for format field: %s", typeDiags.Error())
		}
	}

	// Format of the response
	kind := s.Decoder

	if kind == "" {
		kind = "json"
	}
	if !(kind == "json" || kind == "xml") {
		return cty.NilVal, fmt.Errorf("unsupported response kind: %s", kind)
	}

	// HTTP Method: "GET" or "POST", case-insensitive, default value "GET"
	method := strings.ToUpper(s.Method)

	switch method {
	case "":
		method = http.MethodGet
	case http.MethodGet:
		break
	case http.MethodPost:
		break
	default:
		return cty.NilVal, fmt.Errorf("unsupported method: %s", method)
	}

	// HTTP request timeout:
	//   * a positive integer is accepted and interpreted as value in seconds
	//   * the default value 0 is interpreted as 1 second
	//
	var timeout time.Duration

	switch {
	case *s.Timeout == 0:
		return cty.NilVal, fmt.Errorf("invalid timeout value: %d", *s.Timeout)
	default:
		timeout = time.Duration(*s.Timeout) * time.Second
	}

	// The 'Basic' HTTP Authentication Scheme as defined in RFC 7617.
	//
	// Activated by non-nil `basicAuth` field. When active,
	// * the `basicAuth.username` must be a non-empty string
	// * either `basicAuth.password` or `basicAuth.passwordFile` must be a non-empty string,
	//   with the former overriding the latter.

	var username, password string

	if s.BasicAuth != nil {

		if s.BasicAuth.Username == "" {
			return cty.NilVal, fmt.Errorf("HTTP basic authentication requires a username")
		}

		username = s.BasicAuth.Username

		switch {
		case *s.BasicAuth.Password != "":
			password = *s.BasicAuth.Password

		case *s.BasicAuth.PasswordFile != "":
			byteArr, err := os.ReadFile(filepath.FromSlash(*s.BasicAuth.PasswordFile))
			if err != nil {
				return cty.NilVal, fmt.Errorf("failed to read the password for http basic authentication from a file: %w", err)
			}
			password = string(byteArr)
		}
	}

	// URL query string
	params := url.Values{}
	for k, v := range *s.Params {
		params.Add(k, v)
	}

	// Construct a URL string
	us := s.URL
	if len(params) > 0 {
		us = fmt.Sprint(s.URL, "?", params.Encode())
	}

	// Validate the URL string using the `url` standard library module
	u, err := url.Parse(us)
	if err != nil {
		return cty.NilVal, fmt.Errorf("failed to parse endpoint url %s: %w", us, err)
	}

	// Create an object representing a HTTP request
	var body io.Reader
	httpRequest, err := http.NewRequest(method, u.String(), body)
	if err != nil {
		return cty.NilVal, fmt.Errorf("failed to craft HTTP request object: %w", err)
	}

	// Authentication, if requested
	if s.BasicAuth != nil {
		httpRequest.SetBasicAuth(username, password)
	}

	// Add a bearer token, if requests. Adds a Header.
	var bearerToken string
	switch {
	case s.BearerToken != nil:
		bearerToken = *s.BearerToken
	case s.BearerTokenFile != nil:
		bt, err := os.ReadFile(filepath.FromSlash(*s.BearerTokenFile))
		if err != nil {
			return cty.NilVal, fmt.Errorf("failed to read bearer token file: %w", err)
		}
		bearerToken = strings.TrimSpace(string(bt))
	}
	if bearerToken != "" {
		httpRequest.Header.Set("Authorization", fmt.Sprint("Bearer ", bearerToken))
	}

	// Any other headers, if reqested
	for k, v := range *s.Headers {
		httpRequest.Header.Add(k, v)
	}

	bCtx.Logger.Debug().Str("url", httpRequest.URL.String()).Str("timeout", timeout.String()).Msg("Making HTTP request")
	// Initiate the HTTP client
	c := http.Client{Timeout: timeout}
	// Make REST API request
	httpResponse, err := c.Do(httpRequest)
	if err != nil {
		return cty.NilVal, fmt.Errorf("failed to make HTTP request: %w", err)
	}

	if httpResponse.StatusCode != http.StatusOK {
		body, err := io.ReadAll(httpResponse.Body)
		if err != nil {
			return cty.NilVal, fmt.Errorf("error getting body of response: %w", err)
		}
		return cty.NilVal, fmt.Errorf("HTTP response status code: %d: %s", httpResponse.StatusCode, body)
	}
	defer httpResponse.Body.Close()

	// Decode the body
	switch kind {
	case "json":
		return readJSON(httpResponse.Body, s.Format)
	case "xml":
		return readXML(httpResponse.Body, s.Format)
	}

	return cty.NilVal, fmt.Errorf("unsupported format: %s", kind)
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
	// to cty.Value using a well-defined cty.Type
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
	File     string `hcl:"file,optional"`
	Contents string `hcl:"contents,optional"`
	// Format is the structure of the raw input data defined as a cty.Type
	Format     cty.Type
	FormatExpr hcl.Expression `hcl:"format,attr"`
}

// readJSON reads in, decodes, and validates the format of data
func readJSON(r io.Reader, ty cty.Type) (cty.Value, error) {

	var data interface{}
	if err := json.NewDecoder(r).Decode(&data); err != nil {
		return cty.NilVal, fmt.Errorf("failed to decode JSON: %w", err)
	}
	val, err := gocty.ToCtyValue(data, ty)
	if err != nil {
		return cty.NilVal, err
	}

	return val, nil
}

// Resolve returns a cty.Value representation of the parsed JSON file
func (s *jsonSource) Resolve(bCtx *env.BubblyContext) (cty.Value, error) {
	// Handle the format. First check if the current format is NilType (zero value)
	// and if so, get the type from FormatExpr
	if s.Format == cty.NilType {
		if s.FormatExpr == nil {
			return cty.NilVal, fmt.Errorf("no type available for format field")
		}
		var typeDiags hcl.Diagnostics
		s.Format, typeDiags = typeexpr.TypeConstraint(s.FormatExpr)
		if typeDiags.HasErrors() {
			return cty.NilVal, fmt.Errorf("invalid type for format field: %s", typeDiags.Error())
		}
	}

	if s.File != "" && s.Contents != "" {
		return cty.NilVal, errors.New("cannot provide both file and contents")
	}
	if s.File == "" && s.Contents == "" {
		return cty.NilVal, errors.New("must provide one of file and contents")
	}

	var r io.Reader
	if s.File != "" {
		var err error
		r, err = os.Open(s.File)
		if err != nil {
			return cty.NilVal, fmt.Errorf("error opening file %s: %w", s.File, err)
		}
	} else {
		r = strings.NewReader(s.Contents)
	}

	return readJSON(r, s.Format)
}

// Compiler check to see that v1.XMLSource implements the Source interface
var _ source = (*xmlSource)(nil)

// xmlSource represents the extract type for using an XML file as the input
type xmlSource struct {
	File string `hcl:"file,attr"`
	// the format of the raw input data defined as a cty.Type
	Format     cty.Type
	FormatExpr hcl.Expression `hcl:"format,attr"`
}

// readXML reads in, decodes, and validates the format of data
func readXML(r io.Reader, ty cty.Type) (cty.Value, error) {

	data, err := mxj.NewMapXmlReader(r, true)
	if err != nil {
		return cty.NilVal, fmt.Errorf("failed to decode XML: %w", err)
	}

	if err := fixListsInXML(&data, ty); err != nil {
		return cty.NilVal, err
	}

	val, err := gocty.ToCtyValue(data, ty)
	if err != nil {
		return cty.NilVal, err
	}

	return val, nil
}

// Resolve returns a cty.Value representation of the XML file
func (s *xmlSource) Resolve(bCtx *env.BubblyContext) (cty.Value, error) {
	// Handle the format. First check if the current format is NilType (zero value)
	// and if so, get the type from FormatExpr
	if s.Format == cty.NilType {
		if s.FormatExpr == nil {
			return cty.NilVal, fmt.Errorf("no type available for format field")
		}
		var typeDiags hcl.Diagnostics
		s.Format, typeDiags = typeexpr.TypeConstraint(s.FormatExpr)
		if typeDiags.HasErrors() {
			return cty.NilVal, fmt.Errorf("invalid type for format field: %s", typeDiags.Error())
		}
	}

	mxj.PrependAttrWithHyphen(false) // no "-" prefix on attributes
	mxj.CastNanInf(true)             // use float64, not string for extremes

	f, err := os.Open(s.File)
	if err != nil {
		return cty.NilVal, fmt.Errorf("failed to open file %s: %w", s.File, err)
	}
	defer f.Close()

	return readXML(f, s.Format)
}

// TODO: fixListsInXML could do with extensive unit testing of edge cases and better documentation

// fixListsInXML updates those elements in XML tree who should have been
// the only member of a list of length one. In the absence of a schema or
// a document type definition, the XML parser cannot tell which elements
// must be placed into a list of length one. This is because XML does not
// have syntax for lists, unlike JSON. But the format of data is known to _us_
// via the resource definition. As there is no easy way to communicate this
// information to the XML parser, the second best approach is to add a post-
// processing step. It traverses the data format definition recursively, identifying
// the lists in it, and validates selected branches of the XML tree, updating
// those elements which should have been places in a list.
func fixListsInXML(data *mxj.Map, ty cty.Type) error {

	// Forward declaration: recursive inner function over the format spec
	var f func(*mxj.Map, cty.Type, []string, int) error

	// Implementation: recursive inner function over the format spec
	f = func(data *mxj.Map, ty cty.Type, path []string, idx int) error {

		// The full path to the current node
		// in a format that XML manipulation functions understand.
		pathStr := strings.Join(path, ".")

		if idx > 0 {
			pathStr += fmt.Sprint("[", idx, "]")
		}

		// Traverse all fields of the objects recursively,
		// as they may contain lists as well.
		if ty.IsObjectType() {
			for x := range ty.AttributeTypes() {
				path = append(path, x)
				pathIdx := len(path) - 1

				f(data, ty.AttributeType(x), path, 0)
				path = path[0:pathIdx]
			}
		}

		// For any list declarated in the format spec,
		// the XML tree has to be checked and updated,
		// if necessary...
		if ty.IsListType() {

			// Get the XML tree elements at that level
			vs, err := data.ValuesForPath(pathStr)
			if err != nil {
				return fmt.Errorf("wrong path (%s) in XML syntax tree: %w", pathStr, err)
			}

			n := len(vs)

			// If there is only one conforming element present in the XML tree, it means
			// the parser would not know that it had to make a list for it. To fix that,
			// create a list of length one for that element and update the XML tree.
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
					return f(data, ty.ElementType(), path, i)
				}
			}
		}

		return nil
	}

	path := make([]string, 0)
	return f(data, ty, path, 0)
}
