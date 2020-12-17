package server

import (
	"github.com/labstack/echo/v4"
	"net/http"

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
	if bindErr := c.Bind(&query); bindErr != nil {
		bCtx.Logger.Debug().Err(bindErr).Msg("failed to bind request to queryReq")
		return c.JSON(http.StatusBadRequest, &Error{bindErr.Error()})
	}

	bCtx.Logger.Debug().Str("query", query.Query).Msg("querying the bubbly store")
	results, queryErr := serverStore.Store.Query(query.Query)
	if queryErr != nil {
		bCtx.Logger.Debug().Err(queryErr).Msg("failed while querying the bubbly store")
		return c.JSON(http.StatusBadRequest, &Error{queryErr.Error()})
	}

	return c.JSON(http.StatusOK, &Data{results})
}
