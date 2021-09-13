package server

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/valocode/bubbly/env"
	"github.com/valocode/bubbly/store"
)

func httpErrorHandler(bCtx *env.BubblyConfig, err error, c echo.Context) error {
	bCtx.Logger.Error().
		Str("Path", c.Path()).
		Strs("QueryParams", c.ParamValues()).
		Err(err).
		Msg("API server error")

	// If the error is already an echo HTTP error, just pass it
	if _, ok := err.(*echo.HTTPError); ok {
		return err
	}
	switch {
	case store.IsConflictError(err):
		return echo.NewHTTPError(http.StatusConflict, err.Error())
	case store.IsNotFound(err):
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	case store.IsServerError(err):
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	case store.IsValidationError(err):
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	default:
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
}
