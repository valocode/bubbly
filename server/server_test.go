package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/valocode/bubbly/config"
	"github.com/valocode/bubbly/ent"
	"github.com/valocode/bubbly/store/api"
	"github.com/valocode/bubbly/test"
)

func testGetRequest(t *testing.T, e *echo.Echo, h echo.HandlerFunc, path string, pathParams map[string]string) (*httptest.ResponseRecorder, error) {
	req := httptest.NewRequest(http.MethodGet, path, nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath(path)
	for name, value := range pathParams {
		c.SetParamNames(name)
		c.SetParamValues(value)
	}
	return rec, h(c)
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

func testPutRequest(t *testing.T, e *echo.Echo, h echo.HandlerFunc, path string, pathParams map[string]string, data interface{}) error {
	b := new(bytes.Buffer)
	{
		err := json.NewEncoder(b).Encode(data)
		require.NoError(t, err)
	}
	req := httptest.NewRequest(http.MethodPut, path, b)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	q := req.URL.Query()
	req.URL.RawQuery = q.Encode()
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	for name, value := range pathParams {
		c.SetParamNames(name)
		c.SetParamValues(value)
	}
	return h(c)
}

func TestRelease(t *testing.T) {
	s, err := New(config.NewBubblyConfig())
	require.NoError(t, err)

	{
		err := testPostRequest(t, s.e, s.postRelease, "/releases", &api.ReleaseCreateRequest{
			Project: ent.NewProjectModelCreate().SetName("test-project"),
			Repo:    ent.NewRepoModelCreate().SetName("test-repo"),
			Commit:  ent.NewGitCommitModelCreate().SetHash("123").SetBranch("main").SetTime(time.Now()),
			Release: ent.NewReleaseModelCreate().SetName("test-repo").SetVersion("123"),
		})
		require.NoError(t, err)
	}
	{
		err := testPostRequest(t, s.e, s.postRelease, "/releases", &api.ReleaseCreateRequest{
			Project: ent.NewProjectModelCreate().SetName("another-project"),
			Repo:    ent.NewRepoModelCreate().SetName("another-repo"),
			Commit:  ent.NewGitCommitModelCreate().SetHash("456").SetBranch("main").SetTime(time.Now()),
			Release: ent.NewReleaseModelCreate().SetName("another-repo").SetVersion("456"),
		})
		require.NoError(t, err)
	}
	t.Run("get all releases", func(t *testing.T) {
		rec, err := testGetRequest(t, s.e, s.getReleases, "/releases", nil)
		assert.NoError(t, err)
		var resp api.ReleaseGetResponse
		json.NewDecoder(rec.Body).Decode(&resp)
		assert.Len(t, resp.Releases, 2)
	})
	t.Run("get project releases", func(t *testing.T) {
		rec, err := testGetRequest(t, s.e, s.getReleases, "/releases?projects=test-project", nil)
		assert.NoError(t, err)
		var resp api.ReleaseGetResponse
		json.NewDecoder(rec.Body).Decode(&resp)
		assert.Len(t, resp.Releases, 1)
	})
}

func TestAdapter(t *testing.T) {
	s, err := New(config.NewBubblyConfig())
	require.NoError(t, err)

	{
		err := testPostRequest(t, s.e, s.postAdapter, "/adapters", &api.AdapterSaveRequest{
			Adapter: ent.NewAdapterModelCreate().SetName("test").SetTag("tag").SetModule("test"),
		})
		require.NoError(t, err)
	}
	{
		err := testPostRequest(t, s.e, s.postAdapter, "/adapters", &api.AdapterSaveRequest{
			Adapter: ent.NewAdapterModelCreate().SetName("another").SetTag("default").SetModule("test"),
		})
		require.NoError(t, err)
	}
	t.Run("get all adapters", func(t *testing.T) {
		rec, err := testGetRequest(t, s.e, s.getAdapters, "/adapters", nil)
		assert.NoError(t, err)
		var resp api.AdapterGetResponse
		json.NewDecoder(rec.Body).Decode(&resp)
		assert.Len(t, resp.Adapters, 2)
	})
	t.Run("get existing adapter", func(t *testing.T) {
		rec, err := testGetRequest(t, s.e, s.getAdapters, "/adapters?name=test&tag=tag", nil)
		assert.NoError(t, err)
		var resp api.AdapterGetResponse
		json.NewDecoder(rec.Body).Decode(&resp)
		assert.Len(t, resp.Adapters, 1)
	})
	t.Run("get no adapters", func(t *testing.T) {
		rec, err := testGetRequest(t, s.e, s.getAdapters, "/adapters?name=noname&tag=notag", nil)
		assert.NoError(t, err)
		var resp api.AdapterGetResponse
		json.NewDecoder(rec.Body).Decode(&resp)
		assert.Len(t, resp.Adapters, 0)
	})
}

func TestPolicy(t *testing.T) {
	s, err := New(config.NewBubblyConfig())
	require.NoError(t, err)

	{
		err := testPostRequest(t, s.e, s.postPolicy, "/policies", &api.ReleasePolicySaveRequest{
			Policy: &api.ReleasePolicyCreate{
				ReleasePolicyModelCreate: *ent.NewReleasePolicyModelCreate().
					SetName("test").SetModule("package test"),
			},
		})
		require.NoError(t, err)
	}
	{
		err := testPostRequest(t, s.e, s.postPolicy, "/policies", &api.ReleasePolicySaveRequest{
			Policy: &api.ReleasePolicyCreate{
				ReleasePolicyModelCreate: *ent.NewReleasePolicyModelCreate().
					SetName("another").SetModule("package test"),
			},
		})
		require.NoError(t, err)
	}
	t.Run("get all policies", func(t *testing.T) {
		rec, err := testGetRequest(t, s.e, s.getPolicies, "/policies", nil)
		assert.NoError(t, err)
		var resp api.ReleasePolicyGetResponse
		json.NewDecoder(rec.Body).Decode(&resp)
		assert.Len(t, resp.Policies, 2)
	})
	t.Run("get one policy", func(t *testing.T) {
		rec, err := testGetRequest(t, s.e, s.getPolicies, "/policies?name=test", nil)
		assert.NoError(t, err)
		var resp api.ReleasePolicyGetResponse
		json.NewDecoder(rec.Body).Decode(&resp)
		assert.Len(t, resp.Policies, 1)
	})
	t.Run("update policy", func(t *testing.T) {
		rec, err := testGetRequest(t, s.e, s.getPolicies, "/policies?name=test", nil)
		assert.NoError(t, err)
		var resp api.ReleasePolicyGetResponse
		json.NewDecoder(rec.Body).Decode(&resp)
		require.Len(t, resp.Policies, 1)
		policy := resp.Policies[0]
		putErr := testPutRequest(t, s.e, s.putPolicy, "/policies/:id", map[string]string{"id": strconv.Itoa(*policy.ID)}, api.ReleasePolicyUpdateRequest{
			Policy: &api.ReleasePolicyUpdate{
				ReleasePolicyModelUpdate: *ent.NewReleasePolicyModelUpdate().SetName("somethingelse"),
			},
		})
		assert.NoError(t, putErr)
	})
}

func TestServer(t *testing.T) {
	s, err := New(config.NewBubblyConfig())
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
