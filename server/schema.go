package server

import (
	"fmt"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
)

// PostSchema godoc
// @Summary PostSchema uploads the schema for bubbly
// @ID schema
// @Tag schema
// @Accept json
// @Produce json
// @Success 200 {object} apiResponse
// @Failure 400 {object} apiResponse
// @Router /schema [post]
func (s *Server) PostSchema(c echo.Context) error {

	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("failed to read body of request: %w", err))
	}

	if err := s.Client.PostSchema(s.bCtx, body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, &Status{"schema created!"})
}
