package server

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclparse"
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

	var resourceMap map[string]map[string]map[string]interface{}

	err = json.Unmarshal(byteValue, &resourceMap)
	if err != nil {
		bCtx.Logger.Error().Msg(err.Error())
	}

	body, _ := json.Marshal(resourceMap)
	r := httptest.NewRecorder()

	// Test
	req, _ := http.NewRequest(http.MethodPost, "/api/resource", bytes.NewBuffer(body))
	router.ServeHTTP(r, req)

	assert.Equal(t, http.StatusOK, r.Code)
	assert.Equal(t, "{\"status\":\"uploaded\"}", r.Body.String())

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
	r.GET("/api/resource/default/transform/junit").
		Run(setupRouter(bCtx), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			data := []byte(r.Body.String())
			assert.Equal(t, http.StatusOK, r.Code)
			// Creating a resource based off of the response
			file, diags := hclparse.NewParser().ParseJSON(data, "test.json")
			if diags.HasErrors() {
				bCtx.Logger.Error().Msg(diags.Error())
				t.Errorf(diags.Error())
			}
			resWrap := &core.ResourceBlockHCLWrapper{}
			diags = gohcl.DecodeBody(file.Body, nil, resWrap)
			if diags.HasErrors() {
				bCtx.Logger.Error().Msg(diags.Error())
				t.Errorf(diags.Error())
			}
			assert.NotNil(t, resWrap)
		})
	tearDownDb()
}

// Wipes the test DB
func tearDownDb() {
	_ = os.Remove(resource.MustBuntDBPath())
}
