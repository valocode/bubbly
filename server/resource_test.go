package server

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/rs/zerolog/log"

	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/verifa/bubbly/api/core"

	"github.com/appleboy/gofight"
	"github.com/stretchr/testify/assert"
)

var tearDown = true

func TestPostResource(t *testing.T) {
	// Setup Test Environment
	router := SetupRouter()

	jsonFile, readErr := os.Open("testdata/resource.json")
	if readErr != nil {
		log.Error().Msg(readErr.Error())
	}
	byteValue, _ := ioutil.ReadAll(jsonFile)

	var resourceMap map[string]map[string]map[string]interface{}

	err := json.Unmarshal(byteValue, &resourceMap)
	if err != nil {
		log.Error().Msg(err.Error())
	}

	body, _ := json.Marshal(resourceMap)
	r := httptest.NewRecorder()

	// Test
	req, _ := http.NewRequest("POST", "/api/resource", bytes.NewBuffer(body))
	router.ServeHTTP(r, req)

	assert.Equal(t, 200, r.Code)
	assert.Equal(t, "{\"status\":\"uploaded\"}", r.Body.String())

	// Cleanup
	if tearDown {
		tearDownDb()
	}
}

func TestGetResource(t *testing.T) {
	r := gofight.New()

	// Adds a resource to the db so that it can be fetched
	// Set tearDown to false so that there will be data to GET from the DB
	tearDown = false
	TestPostResource(t)
	tearDown = true
	r.GET("/api/resource/default/translator/junit").
		Run(SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			data := []byte(r.Body.String())
			assert.Equal(t, http.StatusOK, r.Code)
			// Creating a resource based off of the response
			file, diags := hclparse.NewParser().ParseJSON(data, "test.json")
			if diags.HasErrors() {
				log.Error().Msg(diags.Error())
				t.Errorf(diags.Error())
			}
			resWrap := &core.ResourceBlockHCLWrapper{}
			diags = gohcl.DecodeBody(file.Body, nil, resWrap)
			if diags.HasErrors() {
				log.Error().Msg(diags.Error())
				t.Errorf(diags.Error())
			}
			assert.NotNil(t, resWrap)
		})
	tearDownDb()
}

// Wipes the test DB
func tearDownDb() {
	err := os.Remove(DbPath())
	if err != nil {
		log.Error().Msg(err.Error())
	}
}
