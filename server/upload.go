package server

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/verifa/bubbly/api/core"
	"github.com/verifa/bubbly/env"
)

// upload godoc
// @Summary This function will upload core.DataBlocks
// @ID upload data
// @Tags datablocks
// @Param data body uploadStruct true "Datablocks"
// @Accept json
// @Produce json
// @Router /alpha1/upload [post]
func upload(bCtx *env.BubblyContext, c echo.Context) error {
	var data core.DataBlocks
	if err := c.Bind(&data); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if serverStore == nil {
		return &echo.HTTPError{
			Code:    http.StatusBadRequest,
			Message: "server store is not initialized",
		}
	}

	bCtx.Logger.Debug().
		Interface("data", upload).
		Interface("store", serverStore).
		Interface("store.Schema()", serverStore.Schema()).
		Msg("loading data into intermediary database")

	if err := serverStore.Save(data); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, &Status{"uploaded"})
}
