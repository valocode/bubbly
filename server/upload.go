package server

import (
	"encoding/json"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/valocode/bubbly/api/core"
	"github.com/valocode/bubbly/client"
	"github.com/valocode/bubbly/env"
)

// upload godoc
// @Summary This function will upload core.DataBlocks
// @ID upload data
// @Tags datablocks
// @Param data body Data true "Datablocks"
// @Accept json
// @Produce json
// @Success 200 {object} apiResponse
// @Failure 400 {object} apiResponse
// @Router /upload [post]
func (a *Server) upload(bCtx *env.BubblyContext, c echo.Context) error {
	var data core.DataBlocks

	binder := &echo.DefaultBinder{}
	if err := binder.BindBody(c, &data); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	nc := client.NewNATS(bCtx)

	nc.Connect(bCtx)

	sBytes, err := json.Marshal(data)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := nc.Upload(bCtx, sBytes); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, &Status{"uploaded"})
}
