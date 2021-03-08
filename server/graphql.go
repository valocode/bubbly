package server

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/valocode/bubbly/client"
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
func (s *Server) Query(c echo.Context) error {
	var query queryReq

	binder := &echo.DefaultBinder{}
	if err := binder.BindBody(c, &query); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	nc, err := client.New(s.bCtx)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to connect to the NATS server: %w", err))
	}

	results, err := nc.Query(s.bCtx, query.Query)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSONBlob(http.StatusOK, results)
}
