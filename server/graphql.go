package server

import (
	// "encoding/json"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/verifa/bubbly/client"
	"github.com/verifa/bubbly/env"
)

type queryReq struct {
	Query string `json:"query"`
}

// Query godoc
// @Summary Query performs graphql related tasks
// @ID graphql
// @Tags graphql
// @Param query body string true "Query String"
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /api/graphql [post]
func (a *Server) Query(bCtx *env.BubblyContext, c echo.Context) error {
	var query queryReq
	if err := c.Bind(&query); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	nc := client.NewNATS(bCtx)

	nc.Connect(bCtx)

	results, err := nc.Query(bCtx, query.Query)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSONBlob(200, results)
}
