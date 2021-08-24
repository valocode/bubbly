package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
	"github.com/valocode/bubbly/env"
	"github.com/valocode/bubbly/test"
)

func testPostRequest(t *testing.T, e *echo.Echo, h echo.HandlerFunc, path string, data interface{}) {
	b := new(bytes.Buffer)
	{
		err := json.NewEncoder(b).Encode(data)
		require.NoError(t, err)
	}
	req := httptest.NewRequest(http.MethodPost, path, b)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	{
		err := h(c)
		require.NoError(t, err)
	}
}

func TestServer(t *testing.T) {
	s, err := New(env.NewBubblyContext())
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
