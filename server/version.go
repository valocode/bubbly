package server

import (
	"github.com/labstack/echo/v4"
)

var version string

// VersionMiddleware : add version on header.
func VersionMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	// Set out header value for each response
	return func(c echo.Context) error {
		c.Response().Header().Set("X-BUBBLY-VERSION", version)
		return next(c)
	}
}
