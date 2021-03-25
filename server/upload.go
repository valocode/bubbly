package server

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/valocode/bubbly/api/core"
	"github.com/valocode/bubbly/client"
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
func (s *Server) upload(c echo.Context) error {
	var data core.DataBlocks

	binder := &echo.DefaultBinder{}
	if err := binder.BindBody(c, &data); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	nc, err := client.New(s.bCtx)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to connect to the NATS server: %w", err))
	}

	if err := nc.Load(s.bCtx, data); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, &Status{"uploaded"})
}
