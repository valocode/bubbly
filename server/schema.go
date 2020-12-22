package server

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/verifa/bubbly/api/core"
	"github.com/verifa/bubbly/env"
)

// Schema godoc
// @Summary Schema uploads the schema for bubbly
// @ID schema
// @Tags schema
// @Param TODO
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /api/schema [post]
func Schema(bCtx *env.BubblyContext, c echo.Context) error {
	var schema core.Tables
	if err := c.Bind(&schema); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if serverStore != nil {
		// TODO: this is a horrible hack for now... the store should support
		// updates, but currently it only supports create, which does not
		// work a second time. hence, if serverStore is not nil, we cannot
		// update or re-create right now, so just return
		return c.JSON(http.StatusOK, &Status{"schema uploaded!"})
	}

	if err := InitStore(bCtx); err != nil {
		return echo.NewHTTPError(
			http.StatusInternalServerError,
			fmt.Sprintf("failed to initialize the store: %s", err.Error()),
		)
	}

	if err := serverStore.Create(schema); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, &Status{"schema created!"})
}
