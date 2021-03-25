package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/valocode/bubbly/api/core"
	"github.com/valocode/bubbly/client"
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
	var schema core.Tables
	binder := &echo.DefaultBinder{}
	if err := binder.BindBody(c, &schema); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	nc, err := client.New(s.bCtx)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to connect to the NATS server: %w", err))
	}

	sBytes, err := json.Marshal(schema)

	if err != nil {
		echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := nc.PostSchema(s.bCtx, sBytes); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, &Status{"schema created!"})
}
