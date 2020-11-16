package client

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/verifa/bubbly/server"
	"gopkg.in/h2non/gock.v1"

	"github.com/appleboy/gofight/v2"
	"github.com/stretchr/testify/assert"
	"github.com/verifa/bubbly/api/core"
)

// TestLoad validates that a valid core.DataBlocks, with corresponding valid memdb schema, is imported correctly into the memdb database
func TestLoad(t *testing.T) {

	var loadData core.DataBlocks

	filename := "./testdata/load/load_output_sonarqube.json"
	loadJSONData, err := ioutil.ReadFile(filename)

	if err != nil {
		t.Errorf(err.Error())
	}

	// verify that it unmarshals to a valid core.DataBlocks struct instance
	err = json.Unmarshal(loadJSONData, &loadData)

	if err != nil {
		t.Errorf(err.Error())
	}

	t.Log("Unmarshalled to a core.DataBlocks successfully.")

	r := gofight.New()
	r.POST("/alpha1/upload").
		SetJSON(gofight.D{"data": loadData}).
		Run(server.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			t.Logf("Server response: %+v", r)
			assert.Equal(t, http.StatusOK, r.Code)
			assert.Equal(t, "{\"status\":\"uploaded\"}", r.Body.String())
		})
}

// TestClientLoad verifies that some valid core.DataBlocks is POSTed to the correct server route.
func TestClientLoad(t *testing.T) {
	for _, c := range loadDataCases {
		t.Run(c.desc, func(t *testing.T) {
			hostURL := c.sc.Host + ":" + c.sc.Port
			// Create a new server route for mocking a Bubbly server response
			gock.New(hostURL).
				Post(c.route).
				Reply(c.responseCode).
				JSON(c.response)

			var loadData core.DataBlocks
			loadJSONData, err := ioutil.ReadFile(c.inputFile)

			if err != nil {
				t.Errorf(err.Error())
			}

			t.Log("Read file with test data successfully.")

			// verify that it unmarshals to a valid core.DataBlocks
			err = json.Unmarshal(loadJSONData, &loadData)

			if err != nil {
				t.Errorf(err.Error())
			}

			t.Log("Unmarshalled to a core.DataBlocks successfully.")

			c, err := NewClient(c.sc)

			if err != nil {
				t.Errorf(err.Error())
			}

			err = c.Load(loadData)

			if err != nil {
				t.Errorf(err.Error())
			}

		})
	}
}
