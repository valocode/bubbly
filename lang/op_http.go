package lang

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/ext/typeexpr"
	"github.com/zclconf/go-cty/cty"
)

type (
	httpOp struct {

		// The body of a POST request. Ignored for a GET request.
		Query string `hcl:"query,optional"`

		// URL (technically a URI reference, as per RFC 3986) represents the unparsed URL string.
		// Note that URL Query Parameters (?key=val) must be provided separately as Params.
		URL string `hcl:"url"`

		// Method is "GET" or "POST" for protocols "http" and "https". The default is "GET".
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

		// Decoder is one of either "json" and "xml". It will be applied to extract's input.
		Decoder string `hcl:"decoder,optional"`

		// Timeout is the http request timeout
		Timeout string `hcl:"timeout,optional"`

		// Format is a dynamic type, usually built from an HCL type expression.
		// It defines what is expected in response to the REST API query.
		Format     cty.Type
		FormatExpr hcl.Expression `hcl:"format,attr"`
	}

	basicAuth struct {
		Username     string  `hcl:"username"`
		Password     *string `hcl:"password"`
		PasswordFile *string `hcl:"password_file"`
	}
)

// Run performs an HTTP query, parses the response, and returns a corresponding dynamic value.
func (o *httpOp) Run() (cty.Value, error) {
	// Handle the format. First check if the current format is NilType (zero value)
	// and if so, get the type from FormatExpr
	if o.Format == cty.NilType {
		if o.FormatExpr == nil {
			return cty.NilVal, fmt.Errorf("no type available for format field")
		}
		var typeDiags hcl.Diagnostics
		o.Format, typeDiags = typeexpr.TypeConstraint(o.FormatExpr)
		if typeDiags.HasErrors() {
			return cty.NilVal, fmt.Errorf("invalid type for format field: %s", typeDiags.Error())
		}
	}

	// Format of the response
	kind := o.Decoder

	if kind == "" {
		kind = "json"
	}
	if !(kind == "json" || kind == "xml") {
		return cty.NilVal, fmt.Errorf("unsupported response kind: %s", kind)
	}

	// HTTP Method: "GET" or "POST", case-insensitive, default value "GET"
	method := strings.ToUpper(o.Method)

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
	case o.Timeout == "":
		// Set a default timeout
		timeout = time.Duration(time.Second * 10)
	default:
		var err error
		timeout, err = time.ParseDuration(o.Timeout)
		if err != nil {
			return cty.NilVal, fmt.Errorf("invalid timeout: %w", err)
		}
	}

	// The 'Basic' HTTP Authentication Scheme as defined in RFC 7617.
	//
	// Activated by non-nil `basicAuth` field. When active,
	// * the `basicAuth.username` must be a non-empty string
	// * either `basicAuth.password` or `basicAuth.passwordFile` must be a non-empty string,
	//   with the former overriding the latter.

	var username, password string

	if o.BasicAuth != nil {

		if o.BasicAuth.Username == "" {
			return cty.NilVal, fmt.Errorf("HTTP basic authentication requires a username")
		}

		username = o.BasicAuth.Username

		switch {
		case *o.BasicAuth.Password != "":
			password = *o.BasicAuth.Password

		case *o.BasicAuth.PasswordFile != "":
			byteArr, err := os.ReadFile(filepath.FromSlash(*o.BasicAuth.PasswordFile))
			if err != nil {
				return cty.NilVal, fmt.Errorf("failed to read the password for http basic authentication from a file: %w", err)
			}
			password = string(byteArr)
		}
	}

	// URL query string
	params := url.Values{}
	for k, v := range o.Params {
		params.Add(k, v)
	}

	// Construct a URL string
	us := o.URL
	if len(params) > 0 {
		us = fmt.Sprint(o.URL, "?", params.Encode())
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
	if o.BasicAuth != nil {
		httpRequest.SetBasicAuth(username, password)
	}

	// Add a bearer token, if requests. Adds a Header.
	if o.BearerToken != "" && o.BearerTokenFile != "" {
		return cty.NilVal, fmt.Errorf("cannot provide both bearer_token and bearer_token_file")
	}
	if o.BearerToken != "" {
		httpRequest.Header.Set("Authorization", fmt.Sprint("Bearer ", o.BearerToken))
	}
	if o.BearerTokenFile != "" {
		bt, err := os.ReadFile(filepath.FromSlash(o.BearerTokenFile))
		if err != nil {
			return cty.NilVal, fmt.Errorf("failed to read bearer token file: %w", err)
		}
		httpRequest.Header.Set("Authorization", fmt.Sprint("Bearer ", strings.TrimSpace(string(bt))))
	}

	// Any other headers, if reqested
	for k, v := range o.Headers {
		httpRequest.Header.Add(k, v)
	}

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
		return decodeJSON(httpResponse.Body, o.Format)
	case "xml":
		return decodeXML(httpResponse.Body, o.Format)
	}

	return cty.NilVal, fmt.Errorf("unsupported format: %s", kind)
}
