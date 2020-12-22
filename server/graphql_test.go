package server

import (
	"net/http"
	"testing"

	"github.com/appleboy/gofight"
	"github.com/stretchr/testify/assert"
	"github.com/verifa/bubbly/env"
	testData "github.com/verifa/bubbly/server/testdata/upload"
)

func IntegrationTestQuery(t *testing.T) {
	bCtx := env.NewBubblyContext()
	r := gofight.New()
	router := setupRouter(bCtx)

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
		})
}

func IntegrationTestQueryFail(t *testing.T) {
	bCtx := env.NewBubblyContext()
	r := gofight.New()

	r.POST("/api/graphql").
		SetJSON(gofight.D{
			"query": `{tablez(Name:"Boise Zoo"){Name}}`,
		}).
		Run(setupRouter(bCtx), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, http.StatusBadRequest, r.Code)
		})
}
