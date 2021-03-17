package server

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/valocode/bubbly/client"
	"github.com/valocode/bubbly/env"
)

type queryReq struct {
	Query string `json:"query"`
}

type apiResponse struct {
	Code    int    `json:"code" example:"200"`
	Message string `json:"message" example:"data"`
}

// Query godoc
// @Summary Query performs graphql related tasks
// @ID graphql
// @Tags graphql
// @Param query body queryReq true "Query String"
// @Accept json
// @Produce json
// @Success 200 {object} apiResponse
// @Failure 400 {object} apiResponse
// @Failure 404 {object} apiResponse
// @Failure 500 {object} apiResponse
// @Router /graphql [post]
func (*Server) Query(bCtx *env.BubblyContext, c echo.Context) error {
	var query queryReq
	if err := c.Bind(&query); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	nc := client.NewNATS(bCtx)

	if err := nc.Connect(bCtx); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to connect to NATS server: %w", err))
	}

	results, err := nc.Query(bCtx, query.Query)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSONBlob(http.StatusOK, results)
}
