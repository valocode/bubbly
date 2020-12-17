package server

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

var version string

// SetVersion for setup version string.
func SetVersion(ver string) {
	version = ver
}

// GetVersion for get current version.
func GetVersion() string {
	return version
}

func versionHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, &VersionHeaders{
		Source:  "https://github.com/verifa/bubbly",
		Version: GetVersion(),
	})
}

// VersionMiddleware : add version on header.
func VersionMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	// Set out header value for each response
	return func(c echo.Context) error {
		c.Response().Header().Set("X-BUBBLY-VERSION", version)
		return next(c)
	}
}
