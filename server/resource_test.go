package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/verifa/bubbly/api/core"
	"github.com/verifa/bubbly/env"
	"github.com/verifa/bubbly/resource"

	"github.com/appleboy/gofight"
	"github.com/stretchr/testify/assert"
)

var tearDown = true

func TestPostResource(t *testing.T) {
	bCtx := env.NewBubblyContext()
	// Setup Test Environment
	router := setupRouter(bCtx)

	byteValue, err := ioutil.ReadFile("testdata/resource.json")
	if err != nil {
		t.Error(err)
	}

	r := httptest.NewRecorder()

	// Test
	req, _ := http.NewRequest(http.MethodPost, "/api/resource", bytes.NewBuffer(byteValue))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(r, req)

	assert.Equal(t, http.StatusOK, r.Code)

	// Cleanup
	if tearDown {
		tearDownDb()
	}
}

func TestGetResource(t *testing.T) {
	bCtx := env.NewBubblyContext()
	r := gofight.New()

	// Adds a resource to the db so that it can be fetched
	// Set tearDown to false so that there will be data to GET from the DB
	tearDown = false
	TestPostResource(t)
	tearDown = true
	resourceID := "qa/extract/sonarqube"
	r.GET(fmt.Sprintf("/api/resource/%s", resourceID)).
		Run(setupRouter(bCtx), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			data := []byte(r.Body.String())
			assert.Equal(t, http.StatusOK, r.Code)
			var resJSON core.ResourceBlockJSON
			// Creating a resource based off of the response
			err := json.Unmarshal(data, &resJSON)
			assert.NoError(t, err, "unmarshal json resource")
			resBlock, err := resJSON.ResourceBlock()
			assert.NoError(t, err, "convert to ResourceBlock")
			assert.Equal(t, resBlock.String(), resourceID)
		})
	tearDownDb()
}

// Wipes the test DB
func tearDownDb() {
	_ = os.Remove(resource.MustBuntDBPath())
}
