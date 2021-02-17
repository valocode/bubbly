package server

import (
	"fmt"
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
func (*Server) Query(bCtx *env.BubblyContext, c echo.Context) error {
	var query queryReq
	if err := c.Bind(&query); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	nc := client.NewNATS(bCtx)

	if err := nc.Connect(bCtx); err != nil {
		return fmt.Errorf("failed to connect to NATS server: %w", err)
	}

	results, err := nc.Query(bCtx, query.Query)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSONBlob(http.StatusOK, results)
}
