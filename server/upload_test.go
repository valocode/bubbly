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

	jsonReq, _ := json.Marshal(data)
	r := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/alpha1/upload", bytes.NewBuffer(jsonReq))
	router.ServeHTTP(r, req)
	fmt.Println(r.Body)

	r.POST("/upload/alpha1/").
		SetJSON(testData.DataStructFailing()).
		Run(SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, http.StatusBadRequest, r.Code)
		})
}
