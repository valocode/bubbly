package server

import (
	"github.com/labstack/echo/v4"
	"net/http"

	"github.com/verifa/bubbly/api/core"
	"github.com/verifa/bubbly/env"
)

type uploadStruct struct {
	Data core.DataBlocks `json:"data"`
}

// upload godoc
// @Summary This function will upload core.DataBlocks
// @ID upload data
// @Tags datablocks
// @Param data body uploadStruct true "Datablocks"
// @Accept json
// @Produce json
// @Router /alpha1/upload [post]
func upload(bCtx *env.BubblyContext, c echo.Context) error {
	var upload uploadStruct
	if err := c.Bind(&upload); err != nil {
		bCtx.Logger.Error().Msg(err.Error())
		return c.JSON(http.StatusBadRequest, &Error{err.Error()})
	}
	if upload.Data == nil {
		return c.JSON(http.StatusBadRequest, &Error{"malformed request"})
	}

	bCtx.Logger.Debug().
		Interface("data", upload).
		Interface("store", serverStore).
		Interface("store.Schema()", serverStore.Schema()).
		Msg("loading data into intermediary database")

	importErr := serverStore.Store.Save(upload.Data)
	if importErr != nil {
		bCtx.Logger.Error().Msg(importErr.Error())
		return c.JSON(http.StatusBadRequest, &Error{importErr.Error()})
	}

	return c.JSON(http.StatusOK, &Status{"uploaded"})
}
