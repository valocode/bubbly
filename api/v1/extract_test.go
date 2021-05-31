package v1

import (
	"fmt"
	"os"

	"net/http"
	"path/filepath"
	"testing"

	"github.com/rs/zerolog"
	"github.com/valocode/bubbly/env"

	"github.com/stretchr/testify/assert"
	"github.com/zclconf/go-cty/cty"

	"github.com/stretchr/testify/require"

	"gopkg.in/h2non/gock.v1"

	fixtureJSON "github.com/valocode/bubbly/api/v1/testdata/extract/json"
	restGitHub0 "github.com/valocode/bubbly/api/v1/testdata/extract/rest/github"
	fixtureXML "github.com/valocode/bubbly/api/v1/testdata/extract/xml"
)

// newBasicAuth returns a basicAuth struct, with optional fields initialised
// to "correct" values. They are "correct" in a sense that the code in this module
// does not need to do unnecessary checks on empty fields which simplifies the logic.
func newBasicAuth(username, password, passwordFile string) *basicAuth {
	return &basicAuth{username, &password, &passwordFile}
}

func TestExtractJSON(t *testing.T) {
	bCtx := env.NewBubblyContext()

	// Helper function that runs the test defined by its arguments
	run := func(t *testing.T, jsonFile string, ctyType cty.Type, expected cty.Value) {

		t.Helper()

		source := jsonSource{
			File:   jsonFile,
			Format: ctyType,
		}

		val, err := source.Resolve(bCtx)

		assert.Nil(t, err, "failed to Resolve() the extract")
		require.False(t, val.IsNull(), "the extract returned null type value")
		assert.Equal(t, cty.BoolVal(true), val.Equals(expected), "the extract returned unexpected value")
	}

	t.Run("sonarqube-example", func(t *testing.T) {
		run(t,
			filepath.FromSlash("testdata/extract/json/sonarqube-example.json"),
			fixtureJSON.ExpectedType(),
			fixtureJSON.ExpectedValue(),
		)
	})
}

// The XML format is different from JSON in a way that it
// does not have syntax for lists. So the XML parser does not
// know whether an element is by itself, or it's in a list of length one.
// This information is available only in cty.Type data structure we build from HCL
func TestExtractXML(t *testing.T) {

	// Helper function that runs the test defined by its arguments
	run := func(bCtx *env.BubblyContext, t *testing.T, xmlFile string, ctyType cty.Type, expected cty.Value) {

		t.Helper()

		source := xmlSource{
			File:   xmlFile,
			Format: ctyType,
		}

		val, err := source.Resolve(bCtx)

		assert.Nil(t, err, "failed to Resolve() the extract")
		require.False(t, val.IsNull(), "the extract returned null type value")
		assert.Equal(t, cty.BoolVal(true), val.Equals(expected), "the extract returned unexpected value")
	}

	// Easiest. Baseline. No "short" lists.
	t.Run("junit0", func(t *testing.T) {
		bCtx := env.NewBubblyContext()
		run(bCtx,
			t,
			filepath.FromSlash(`testdata/extract/xml/junit0.xml`),
			fixtureXML.ExpectedType(),
			fixtureXML.ExpectedValue0(),
		)
	})

	// Harder. A single "testsuite" element but multiple "testcase" elements therein.
	t.Run("junit1", func(t *testing.T) {
		bCtx := env.NewBubblyContext()
		run(bCtx,
			t,
			filepath.FromSlash(`testdata/extract/xml/junit1.xml`),
			fixtureXML.ExpectedType(),
			fixtureXML.ExpectedValue1(),
		)
	})

	// Hardest. A single "testsuite" element with a single "testcase" elements within.
	t.Run("junit2", func(t *testing.T) {
		bCtx := env.NewBubblyContext()
		run(bCtx,
			t,
			filepath.FromSlash(`testdata/extract/xml/junit2.xml`),
			fixtureXML.ExpectedType(),
			fixtureXML.ExpectedValue2(),
		)
	})
}

func TestExtractGit(t *testing.T) {
	bCtx := env.NewBubblyContext()

	source := gitSource{
		Directory: filepath.FromSlash(`../../integration/testdata/git/repo1.git`),
	}

	expected := cty.ObjectVal(map[string]cty.Value{
		"active_branch": cty.StringVal("master"),
		"branches": cty.ObjectVal(map[string]cty.Value{
			"local":  cty.ListVal([]cty.Value{cty.StringVal("dev"), cty.StringVal("master")}),
			"remote": cty.NullVal(cty.List(cty.String)),
		}),
		"commit_id": cty.StringVal("81411ea85f68f64f727f140400d7107786d93ba4"),
		"is_bare":   cty.True,
		"remotes": cty.ListValEmpty(cty.Object(map[string]cty.Type{
			"name": cty.String,
			"url":  cty.String,
		})),
		"tag": cty.StringVal("kawabunga"),
	})

	val, err := source.Resolve(bCtx)

	assert.Nil(t, err, "failed to Resolve() the extract")
	require.False(t, val.IsNull(), "the extract returned null type value")
	assert.Equal(t, cty.BoolVal(true), val.Equals(expected), "the extract returned unexpected value")

}

func TestExtractGraphQL(t *testing.T) {

	// Setting up Bubbly
	bCtx := env.NewBubblyContext()
	bCtx.UpdateLogLevel(zerolog.DebugLevel)

	// Using GitHub GraphQL API as it is rich, mature, and interesting!
	url := "https://api.github.com/graphql"

	// The GraphQL query to send:
	// Read public repo information.
	query := `query { 
		repository(owner:"octocat", name:"Hello-World") {
			issues(first:3, states:CLOSED) {
		  		edges {
					node {
			  			title
			  			url
			  		}
		  		}
			}
	  	}
	}`

	// To communicate with the GraphQL server, you'll need an OAuth token with the right scopes.
	// To read public repo information, `public_repo` scope (a subscope within `repo`) is required.
	// For private repo information, the full `repo` scope is required.
	// More info: https://docs.github.com/en/graphql/guides/forming-calls-with-graphql#authenticating-with-graphql
	var bearerToken string
	if value, ok := os.LookupEnv("GH_TOKEN"); ok {
		bearerToken = value
	} else {
		if value, ok := os.LookupEnv("GITHUB_TOKEN"); ok {
			bearerToken = value
		} else {
			t.Skip("GitHub authentication token not set. Must set either $GH_TOKEN or $GITHUB_TOKEN to run this test.")
		}
	}

	// Set up `extract/graphql` instance
	source := graphqlSource{
		URL:         url,
		Query:       query,
		BearerToken: bearerToken,
		Format: cty.Object(map[string]cty.Type{
			"repository": cty.Object(map[string]cty.Type{
				"issues": cty.Object(map[string]cty.Type{
					"edges": cty.List(cty.Object(map[string]cty.Type{
						"node": cty.Object(map[string]cty.Type{
							"title": cty.String,
							"url":   cty.String,
						}),
					})),
				}),
			}),
		}),
	}
	setGraphQLSourceDefaults(bCtx, &source)

	expected := cty.ObjectVal(map[string]cty.Value{
		"repository": cty.ObjectVal(map[string]cty.Value{
			"issues": cty.ObjectVal(map[string]cty.Value{
				"edges": cty.ListVal([]cty.Value{
					cty.ObjectVal(map[string]cty.Value{
						"node": cty.ObjectVal(map[string]cty.Value{
							"title": cty.StringVal("Hello World in all programming languages"),
							"url":   cty.StringVal("https://github.com/octocat/Hello-World/issues/7"),
						}),
					}),
					cty.ObjectVal(map[string]cty.Value{
						"node": cty.ObjectVal(map[string]cty.Value{
							"title": cty.StringVal("test100"),
							"url":   cty.StringVal("https://github.com/octocat/Hello-World/issues/10"),
						}),
					}),
					cty.ObjectVal(map[string]cty.Value{
						"node": cty.ObjectVal(map[string]cty.Value{
							"title": cty.StringVal("test100"),
							"url":   cty.StringVal("https://github.com/octocat/Hello-World/issues/11"),
						}),
					}),
				}),
			}),
		}),
	})

	// Send the request to GitHub API endpoint and read the result
	val, err := source.Resolve(bCtx)

	// Validate the response
	require.NoError(t, err, "error in resolving GraphQL extract")
	require.NotNil(t, val, "nil value returned by GraphQL extract resolver")
	require.False(t, val.IsNull(), "cty.NullVal returned by GraphQL extract resolver")
	require.Equal(t, cty.BoolVal(true), val.Equals(expected))
}

func TestExtractRESTfulJSON(t *testing.T) {

	defer gock.Off()
	bCtx := env.NewBubblyContext()

	// Expected format of the end-point response
	rFormat := cty.Object(map[string]cty.Type{
		"status": cty.String,
		"data": cty.Object(map[string]cty.Type{
			"verified": cty.Bool,
			"items":    cty.List(cty.Number),
		}),
	})

	// The value returned by the JSON end-point
	rJSON := map[string]interface{}{
		"status": "ok",
		"data": map[string]interface{}{
			"verified": true,
			"items":    []int{0, 1, 2},
		},
	}

	// Expected result of unmarshaling the value received
	expected := cty.ObjectVal(map[string]cty.Value{
		"status": cty.StringVal("ok"),
		"data": cty.ObjectVal(map[string]cty.Value{
			"verified": cty.BoolVal(true),
			"items": cty.ListVal([]cty.Value{
				cty.NumberIntVal(0),
				cty.NumberIntVal(1),
				cty.NumberIntVal(2),
			}),
		}),
	})

	// Describe the request
	scheme := "https"
	host := "localhost"
	port := uint16(8080)
	route := "get/sequence"

	url := fmt.Sprint(scheme, "://", host, ":", port, "/", route)

	source := restSource{
		URL:     url,
		Decoder: "json",
		Format:  rFormat,
	}
	setRestSourceDefaults(bCtx, &source)

	//
	// Subtests...
	//

	t.Run("json/GET", func(t *testing.T) {

		// Mock the HTTP server and the REST API endpoint
		s := source
		method := http.MethodGet
		s.Method = method

		gockResponse := gock.New(s.URL).
			Get(route).
			Reply(http.StatusOK).
			JSON(rJSON)

		// Make a REST API request and decode the response
		val, err := s.Resolve(bCtx)

		require.True(t, gockResponse.Done(), "mock is not done")
		assert.Equal(t, http.StatusOK, gockResponse.StatusCode, "HTTP error")

		assert.Nil(t, err, "failed to Resolve() the extract")
		require.False(t, val.IsNull(), "null value unmarshaled")
		assert.Equal(t, cty.BoolVal(true), val.Equals(expected), "unexpected value unmarshaled")
	})

	t.Run("json/POST", func(t *testing.T) {

		// Mock the HTTP server and the REST API endpoint
		s := source
		method := http.MethodPost
		s.Method = method

		gockResponse := gock.New(s.URL).
			Post(route).
			Reply(http.StatusOK).
			JSON(rJSON)

		// Make a REST API request and decode the response
		val, err := s.Resolve(bCtx)

		require.True(t, gockResponse.Done(), "mock is not done")
		assert.Equal(t, http.StatusOK, gockResponse.StatusCode, "HTTP error")

		assert.Nil(t, err, "failed to Resolve() the extract")
		require.False(t, val.IsNull(), "null value unmarshaled")
		assert.Equal(t, cty.BoolVal(true), val.Equals(expected), "unexpected value unmarshaled")
	})
}

func TestExtractRESTfulXML(t *testing.T) {

	defer gock.Off()
	bCtx := env.NewBubblyContext()

	// Expected format of the end-point response
	rFormat := cty.Object(map[string]cty.Type{
		"root": cty.Object(map[string]cty.Type{
			"status": cty.String,
			"data": cty.Object(map[string]cty.Type{
				"verified": cty.Bool,
				"count":    cty.Number,
			}),
		}),
	})

	// The value returned by the XML end-point
	sXML := `<root><status>ok</status><data><verified>true</verified><count>101</count></data></root>`

	// Expected result of unmarshaling the value received
	expected := cty.ObjectVal(map[string]cty.Value{
		"root": cty.ObjectVal(map[string]cty.Value{
			"status": cty.StringVal("ok"),
			"data": cty.ObjectVal(map[string]cty.Value{
				"verified": cty.BoolVal(true),
				"count":    cty.NumberIntVal(101),
			}),
		}),
	})

	// Describe the request
	scheme := "https"
	host := "localhost"
	port := uint16(8080)
	route := "get/sequence"

	url := fmt.Sprint(scheme, "://", host, ":", port, "/", route)

	source := restSource{
		URL:     url,
		Decoder: "xml",
		Format:  rFormat,
	}
	setRestSourceDefaults(bCtx, &source)

	t.Run("xml/GET", func(t *testing.T) {

		// Mock the HTTP server and the REST API endpoint
		s := source
		method := http.MethodGet
		s.Method = method

		gockResponse := gock.New(s.URL).
			Get(route).
			Reply(http.StatusOK).
			BodyString(sXML)

		// Make a REST API request and decode the response
		val, err := s.Resolve(bCtx)

		require.True(t, gockResponse.Done(), "mock is not done")
		assert.Equal(t, http.StatusOK, gockResponse.StatusCode, "HTTP error")

		assert.Nil(t, err, "failed to Resolve() the extract")
		require.False(t, val.IsNull(), "null value unmarshaled")
		assert.Equal(t, cty.BoolVal(true), val.Equals(expected), "unexpected value unmarshaled")
	})

	// TODO: xml/POST with a Param
}

func TestExtractRestBasicAuth(t *testing.T) {

	defer gock.Off()
	bCtx := env.NewBubblyContext()

	// Describe a REST API Extract Resource with HTTP Basic Authorization
	scheme := "https"
	host := "api.cloud84.dev"
	port := uint16(9090)
	route := "gruffalo/hello"

	url := fmt.Sprint(scheme, "://", host, ":", port, "/", route)

	source := restSource{
		URL:     url,
		Decoder: "json",
		Format:  cty.EmptyObject,
	}
	setRestSourceDefaults(bCtx, &source)

	// JSON returned upon successful authorization
	responseBody := "{}"

	// The result of successful conversion of server response
	expected := cty.EmptyObjectVal

	// Subtest
	t.Run("username only", func(t *testing.T) {

		s := source
		s.BasicAuth = newBasicAuth("mouse", "", "")

		// The answer we expect: error

		// The API request that we expect and the response that we send in that case
		gockResponse := gock.New(s.URL).
			Get("/").
			Reply(http.StatusUnauthorized)

		// Make API request
		val, err := s.Resolve(bCtx)

		// Do checks suitable for the testing scenario
		assert.True(t, gockResponse.Done(), "server did not understand the request")
		assert.Equal(t, http.StatusUnauthorized, gockResponse.StatusCode, "HTTP status code")
		assert.NotNil(t, err, "expected error but no error was returned")
		require.True(t, val.IsNull(), "expected error: the extract returned null type value")
	})

	// Subtest
	t.Run("username and password", func(t *testing.T) {

		s := source
		s.BasicAuth = newBasicAuth("mouse", "correct horse battery staple", "")

		gockResponse := gock.New(s.URL).
			Get("/").
			BasicAuth(s.BasicAuth.Username, *s.BasicAuth.Password).
			Reply(http.StatusOK).
			BodyString(responseBody)

		val, err := s.Resolve(bCtx)

		assert.True(t, gockResponse.Done(), "server did not understand the request")
		assert.Equal(t, http.StatusOK, gockResponse.StatusCode, "HTTP status code")
		assert.Nil(t, err, "failed to Resolve() the extract")
		require.False(t, val.IsNull(), "the extract returned null type value")
		assert.Equal(t, cty.BoolVal(true), val.Equals(expected), "the extract returned unexpected value")
	})

	// Subtest
	t.Run("username and password file", func(t *testing.T) {

		s := source
		s.BasicAuth = newBasicAuth(
			"mouse",
			"",
			filepath.FromSlash("./testdata/extract/rest/secret"),
		)

		// Read password value from a test fixture file
		password, err := os.ReadFile(*s.BasicAuth.PasswordFile)
		require.Nil(t, err, "test fixture: password containing file")

		gockResponse := gock.New(s.URL).
			Get("/").
			BasicAuth(s.BasicAuth.Username, string(password)).
			Reply(http.StatusOK).
			BodyString(responseBody)

		val, err := s.Resolve(bCtx)

		assert.True(t, gockResponse.Done(), "server did not understand the request")
		assert.Equal(t, http.StatusOK, gockResponse.StatusCode, "HTTP status code")
		assert.Nil(t, err, "failed to Resolve() the extract")
		require.False(t, val.IsNull(), "the extract returned null type value")
		assert.Equal(t, cty.BoolVal(true), val.Equals(expected), "the extract returned unexpected value")
	})

	// Subtest
	t.Run("username and password and password file", func(t *testing.T) {

		s := source
		s.BasicAuth = newBasicAuth(
			"mouse",
			"correct horse battery staple",
			filepath.FromSlash("./testdata/extract/rest/secret"),
		)

		gockResponse := gock.New(s.URL).
			Get("/").
			BasicAuth(s.BasicAuth.Username, *s.BasicAuth.Password).
			Reply(http.StatusOK).
			BodyString(responseBody)

		val, err := s.Resolve(bCtx)

		assert.True(t, gockResponse.Done(), "server did not understand the request")
		assert.Equal(t, http.StatusOK, gockResponse.StatusCode, "HTTP status code")
		assert.Nil(t, err, "failed to Resolve() the extract")
		require.False(t, val.IsNull(), "the extract returned null type value")
		assert.Equal(t, cty.BoolVal(true), val.Equals(expected), "the extract returned unexpected value")
	})
}

func TestExtractRestBearerToken(t *testing.T) {

	defer gock.Off()
	bCtx := env.NewBubblyContext()

	scheme := "https"
	host := "api.cloud84.dev"
	port := uint16(9090)
	route := "private"

	url := fmt.Sprint(scheme, "://", host, ":", port, "/", route)

	source := restSource{
		URL:     url,
		Decoder: "json",
		Format:  cty.EmptyObject,
	}
	setRestSourceDefaults(bCtx, &source)

	responseBody := "{}"
	expected := cty.EmptyObjectVal

	// Subtest
	t.Run("bearer token only", func(t *testing.T) {

		bearerToken := "3048d70dc7c4e4ccf47916e809ef2019eaef41d68e46ff100560807bbe1572f9"
		s := source
		s.BearerToken = &bearerToken

		gockResponse := gock.New(s.URL).
			Get("/").
			MatchHeader("Authorization", fmt.Sprint("Bearer ", bearerToken)).
			Reply(http.StatusOK).
			BodyString(responseBody)

		val, err := s.Resolve(bCtx)

		assert.True(t, gockResponse.Done(), "server did not understand the request")
		assert.Equal(t, http.StatusOK, gockResponse.StatusCode, "HTTP status code")
		assert.Nil(t, err, "failed to Resolve() the extract")
		require.False(t, val.IsNull(), "the extract returned null type value")
		assert.Equal(t, cty.BoolVal(true), val.Equals(expected), "the extract returned unexpected value")
	})

	// Subtest
	t.Run("bearer token file only", func(t *testing.T) {

		bearerTokenFile := filepath.FromSlash("./testdata/extract/rest/bearer_token_secret")

		s := source
		s.BearerTokenFile = &bearerTokenFile

		// Read bearer token value from a test fixture file
		bearerToken, err := os.ReadFile(*s.BearerTokenFile)
		require.Nil(t, err, "test fixture: bearer token file")

		gockResponse := gock.New(s.URL).
			Get("/").
			MatchHeader("Authorization", fmt.Sprint("Bearer ", string(bearerToken))).
			Reply(http.StatusOK).
			BodyString(responseBody)

		val, err := s.Resolve(bCtx)

		assert.True(t, gockResponse.Done(), "server did not understand the request")
		assert.Equal(t, http.StatusOK, gockResponse.StatusCode, "HTTP status code")
		assert.Nil(t, err, "failed to Resolve() the extract")
		require.False(t, val.IsNull(), "the extract returned null type value")
		assert.Equal(t, cty.BoolVal(true), val.Equals(expected), "the extract returned unexpected value")
	})

	// Subtest
	t.Run("bearer token and bearer token file both", func(t *testing.T) {

		bearerToken := "3048d70dc7c4e4ccf47916e809ef2019eaef41d68e46ff100560807bbe1572f9"
		bearerTokenFile := filepath.FromSlash("./testdata/extract/rest/bearer_token_secret")

		s := source
		s.BearerToken = &bearerToken
		s.BearerTokenFile = &bearerTokenFile

		gockResponse := gock.New(s.URL).
			Get("/").
			MatchHeader("Authorization", fmt.Sprint("Bearer ", *s.BearerToken)).
			Reply(http.StatusOK).
			BodyString(responseBody)

		val, err := s.Resolve(bCtx)

		assert.True(t, gockResponse.Done(), "server did not understand the request")
		assert.Equal(t, http.StatusOK, gockResponse.StatusCode, "HTTP status code")
		assert.Nil(t, err, "failed to Resolve() the extract")
		require.False(t, val.IsNull(), "the extract returned null type value")
		assert.Equal(t, cty.BoolVal(true), val.Equals(expected), "the extract returned unexpected value")
	})
}

func TestExtractRestHeaders(t *testing.T) {

	defer gock.Off()
	bCtx := env.NewBubblyContext()

	// Some APIs request that certain HTTP headers
	// "must" or "highly recommended" to be set.
	scheme := "http"
	host := "api.github.com"
	port := uint16(9090)
	route := "users/olliefr"

	url := fmt.Sprint(scheme, "://", host, ":", port, "/", route)

	source := restSource{
		URL:     url,
		Decoder: "json",
		Format:  cty.EmptyObject,
	}
	setRestSourceDefaults(bCtx, &source)

	responseBody := "{}"
	expected := cty.EmptyObjectVal

	// Subtest
	t.Run("github content type", func(t *testing.T) {
		s := source

		s.Headers = &map[string]string{
			"Accept":     "application/vnd.github.v3+json",
			"User-Agent": "Bubbly REST API Extract",
		}

		gockResponse := gock.New(s.URL).
			Get("/").
			MatchHeaders(*s.Headers).
			Reply(http.StatusOK).
			BodyString(responseBody)

		val, err := s.Resolve(bCtx)

		assert.True(t, gockResponse.Done(), "server did not understand the request")
		assert.Equal(t, http.StatusOK, gockResponse.StatusCode, "HTTP status code")
		assert.Nil(t, err, "failed to Resolve() the extract")
		require.False(t, val.IsNull(), "the extract returned null type value")
		assert.Equal(t, cty.BoolVal(true), val.Equals(expected), "the extract returned unexpected value")
	})
}

func TestExtractRestParams(t *testing.T) {

	defer gock.Off()
	bCtx := env.NewBubblyContext()

	// This is a more advanced REST API request which encodes certain
	// parameters in its query URL. The "full" query returns three entities,
	// while a following query with a parameter "per_page=2" set - returns only two.

	// The type and value of the result of this mocked REST API request
	// is taken from Prometheus' "HTTP API" *Rules* section. As the /rules
	// endpoint is fairly new, it does not have the same stability guarantees
	// as the overarching API v1.
	//
	// The API endpoint is:
	//   GET http://localhost:9090/api/v1/rules

	scheme := "https"
	host := "api.github.com"
	port := uint16(9090)
	route := "repos/octocat/hello-world/branches"

	url := fmt.Sprint(scheme, "://", host, ":", port, "/", route)

	source := restSource{
		URL:     url,
		Decoder: "json",
		Format:  restGitHub0.ExpectedType(),
	}
	setRestSourceDefaults(bCtx, &source)

	// Subtest
	t.Run("no url query string", func(t *testing.T) {

		s := source

		// No parameters for request means all entries are returned
		expected := restGitHub0.ExpectedValue()

		gockResponse := gock.New(s.URL).
			Get("/").
			Reply(http.StatusOK).
			File(filepath.FromSlash("./testdata/extract/rest/github/github0.json"))

		val, err := s.Resolve(bCtx)

		assert.True(t, gockResponse.Done(), "server did not understand the request")
		assert.Equal(t, http.StatusOK, gockResponse.StatusCode, "HTTP status code")
		assert.Empty(t, gockResponse.Mock.Request().URLStruct.RawQuery, "unexpected URL query key-value pairs in the request")
		assert.Nil(t, err, "failed to Resolve() the extract")
		require.False(t, val.IsNull(), "the extract returned null type value")
		assert.Equal(t, cty.BoolVal(true), val.Equals(expected), "the extract returned unexpected value")
	})

	// Subtest
	t.Run("set per_page=2 in URL query string", func(t *testing.T) {

		s := source

		// Add a key-value pair to the URL query string
		s.Params = &map[string]string{
			"per_page": "2",
		}

		// The page_limit parameter limits the response to the first two entries
		expected := cty.ListVal(restGitHub0.ExpectedValue().AsValueSlice()[:2])

		gockResponse := gock.New(s.URL).
			Get("/").
			MatchParam("per_page", "2").
			Reply(http.StatusOK).
			File(filepath.FromSlash("./testdata/extract/rest/github/github1.json"))

		val, err := s.Resolve(bCtx)

		require.True(t, gockResponse.Done(), "request did not match what server expected")
		assert.Equal(t, http.StatusOK, gockResponse.StatusCode, "HTTP status code")
		assert.NotEmpty(t, gockResponse.Mock.Request().URLStruct.RawQuery, "URL query string")
		assert.Nil(t, err, "failed to Resolve() the extract")
		require.False(t, val.IsNull(), "the extract returned null type value")
		assert.Equal(t, cty.BoolVal(true), val.Equals(expected), "the extract returned unexpected value")
	})
}
