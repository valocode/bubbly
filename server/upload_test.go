package server

import (
	"net/http"
	"testing"

	"github.com/appleboy/gofight/v2"
	"github.com/stretchr/testify/assert"
	testData "github.com/verifa/bubbly/server/testdata/upload"
)

// This creates a passing test for the Upload route
func TestUploadPassing(t *testing.T) {
	r := gofight.New()
	r.POST("/upload/alpha1/").
		SetJSON(testData.DataStruct()).
		Run(SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, http.StatusOK, r.Code)
			assert.Equal(t, "{\"status\":\"uploaded\"}", r.Body.String())
		})
}

// This creates a failing test case. The struct binding should fail on account of improperly named fields
func TestUploadFailing(t *testing.T) {
	r := gofight.New()

	r.POST("/upload/alpha1/").
		SetJSON(testData.DataStructFailing()).
		Run(SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, http.StatusBadRequest, r.Code)
		})
}
