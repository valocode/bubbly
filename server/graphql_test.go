package server

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/appleboy/gofight"
	"github.com/stretchr/testify/assert"
	testData "github.com/verifa/bubbly/server/testdata/upload"
)

func IntegrationTestQuery(t *testing.T) {
	r := gofight.New()
	router := SetupRouter()

	// First, insert data into MemDb using the Upload functionality
	r.POST("/alpha1/upload").
		SetJSON(gofight.D{"data": testData.DataStruct()}).
		Run(router, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, http.StatusOK, r.Code)
			assert.Equal(t, "{\"status\":\"uploaded\"}", r.Body.String())

		})
	// Finally, Test the query function of graphql
	r.POST("/api/graphql").
		SetJSON(gofight.D{
			"query": `{product(Name:"1234"){Name}}`,
		}).
		Run(router, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, http.StatusOK, r.Code)
			fmt.Println("")
		})
}

func IntegrationTestQueryFail(t *testing.T) {
	r := gofight.New()

	r.POST("/api/graphql").
		SetJSON(gofight.D{
			"query": `{tablez(Name:"Boise Zoo"){Name}}`,
		}).
		Run(SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, http.StatusBadRequest, r.Code)
		})
}
