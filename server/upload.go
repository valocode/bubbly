package server

import (
	"fmt"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
)

// upload godoc
// @Summary This function will upload core.DataBlocks
// @ID upload data
// @Tags datablocks
// @Param data body []byte true "Datablocks"
// @Accept json
// @Produce json
// @Success 200 {object} apiResponse
// @Failure 400 {object} apiResponse
// @Router /upload [post]
func (s *Server) upload(c echo.Context) error {

	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("failed to read body of request: %w", err))
	}

	if err := s.Client.Load(s.bCtx, body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, &Status{"uploaded"})
}
