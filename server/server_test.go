package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/valocode/bubbly/ent"
	"github.com/valocode/bubbly/env"
	"github.com/valocode/bubbly/store"
	"github.com/valocode/bubbly/store/api"
	"github.com/valocode/bubbly/test"
)

func testGetRequest(t *testing.T, e *echo.Echo, h echo.HandlerFunc, path string, pathParams map[string]string) error {
	req := httptest.NewRequest(http.MethodGet, path, nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	q := req.URL.Query()
	for name, value := range pathParams {
		q.Add(name, value)
	}
	req.URL.RawQuery = q.Encode()
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath(path)
	for name, value := range pathParams {
		c.SetParamNames(name)
		c.SetParamValues(value)
	}
	return h(c)
}

func testPostRequest(t *testing.T, e *echo.Echo, h echo.HandlerFunc, path string, data interface{}) error {
	b := new(bytes.Buffer)
	{
		err := json.NewEncoder(b).Encode(data)
		require.NoError(t, err)
	}
	req := httptest.NewRequest(http.MethodPost, path, b)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	return h(c)
}

func testPutRequest(t *testing.T, e *echo.Echo, h echo.HandlerFunc, path string, data interface{}) error {
	b := new(bytes.Buffer)
	{
		err := json.NewEncoder(b).Encode(data)
		require.NoError(t, err)
	}
	req := httptest.NewRequest(http.MethodPut, path, b)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	return h(c)
}

func TestAdapter(t *testing.T) {
	s, err := New(env.NewBubblyConfig())
	require.NoError(t, err)

	{
		err := testPostRequest(t, s.e, s.postAdapter, "/adapters", &api.AdapterSaveRequest{
			Adapter: ent.NewAdapterModelCreate().SetName("test").SetTag("tag").SetModule("test"),
		})
		require.NoError(t, err)
	}
	{
		err := testGetRequest(t, s.e, s.getAdapter, "/adapters/:name", map[string]string{
			"name": "test",
		})
		var expErr *store.NotFoundError
		assert.ErrorAs(t, err, &expErr)
	}
	{
		err := testGetRequest(t, s.e, s.getAdapter, "/adapters/:name?tag=tag", map[string]string{
			"name": "test",
		})
		assert.NoError(t, err)
	}
}

func TestPolicy(t *testing.T) {
	s, err := New(env.NewBubblyConfig())
	require.NoError(t, err)

	{
		err := testPostRequest(t, s.e, s.postPolicy, "/policies", &api.ReleasePolicySaveRequest{
			Policy: ent.NewReleasePolicyModelCreate().
				SetName("test").SetModule("package test"),
		})
		require.NoError(t, err)
	}
	{
		err := testGetRequest(t, s.e, s.getPolicy, "/policies/:name", map[string]string{
			"name": "not exist",
		})
		var expErr *store.NotFoundError
		assert.ErrorAs(t, err, &expErr)
	}
	{
		err := testGetRequest(t, s.e, s.getPolicy, "/policies/:name", map[string]string{
			"name": "test",
		})
		assert.NoError(t, err)
	}
}

func TestServer(t *testing.T) {
	s, err := New(env.NewBubblyConfig())
	require.NoError(t, err)

	data := test.CreateDummyData()
	for _, repo := range data {
		for _, release := range repo.Releases {
			// t.Log("Creating Release: ", release.Release.Commit.Hash)
			testPostRequest(t, s.e, s.postRelease, "/releases", release.Release)
			for _, art := range release.Artifacts {
				testPostRequest(t, s.e, s.postArtifact, "/artifacts", art)
			}
			for _, scan := range release.CodeScans {
				testPostRequest(t, s.e, s.postCodeScan, "/codescans", scan)
			}
			for _, run := range release.TestRuns {
				testPostRequest(t, s.e, s.postTestRun, "/testruns", run)
			}
		}
	}
}
