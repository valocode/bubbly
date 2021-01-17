package server

import (
	"net/http"

	"github.com/labstack/echo/v4"

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
func Query(bCtx *env.BubblyContext, c echo.Context) error {
	var query queryReq
	if err := c.Bind(&query); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	bCtx.Logger.Info().Msgf("Received Graphql Query: %s", query)

	results, err := serverStore.Query(query.Query)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, &Data{results})
}
