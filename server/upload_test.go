package server

import (
	"net/http"
	"testing"

	"github.com/appleboy/gofight/v2"
	"github.com/stretchr/testify/assert"
	testData "github.com/verifa/bubbly/server/testdata/upload"
)

// This creates a passing test for the upload route
func IntegrationTestUploadPassing(t *testing.T) {
	r := gofight.New()

	r.POST("/alpha1/upload").
		SetJSON(gofight.D{"data": testData.DataStruct()}).
		Run(SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, http.StatusOK, r.Code)
			assert.Equal(t, "{\"status\":\"uploaded\"}", r.Body.String())
		})

}

func TestUploadFailing(t *testing.T) {
	r := gofight.New()
	r.POST("/alpha1/upload").
		SetJSON(testData.DataStructFailing()).
		Run(SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, http.StatusBadRequest, r.Code)
		})
}
