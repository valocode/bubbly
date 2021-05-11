package server

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/valocode/bubbly/agent/component"
)

func (s *Server) authMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		var (
			authHeader = c.Request().Header.Get(echo.HeaderAuthorization)
			client     = &http.Client{}
		)
		// Make an exception that /healthz does not require authentication
		if c.Path() == "/healthz" {
			return next(c)
		}
		// Ignore requests to swagger documentation
		if strings.HasPrefix(c.Path(), "/swagger") {
			return next(c)
		}
		if authHeader == "" {
			return echo.NewHTTPError(http.StatusForbidden, "missing authentication token")
		}
		// Create a new request to get the authorization information
		req, err := http.NewRequest(http.MethodGet, s.bCtx.AuthConfig.AuthAddr+"/authorize", nil)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "failed to create authorization request")
		}
		// Copy the received header into the request
		req.Header.Add(echo.HeaderAuthorization, authHeader)

		// If the URL path contains the organization param in it we need to
		// provide that to the authorization service to check that the user
		// has access to the organization
		if org := c.Param("organization"); org != "" {
			s.bCtx.Logger.Debug().Str("org", org).Msgf("Adding org to authz request")
			// Add parameters to the query
			q := req.URL.Query()
			q.Add("organization", org)
			req.URL.RawQuery = q.Encode()
		}

		resp, err := client.Do(req)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "authorization request failed: "+err.Error())
		}
		if resp.StatusCode != http.StatusOK {
			return echo.NewHTTPError(resp.StatusCode, "authorization request failed")
		}

		var auth component.MessageAuth
		err = json.NewDecoder(resp.Body).Decode(&auth)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "failed getting request authentication: "+err.Error())
		}
		// Set the request auth into the echo request so that the handlers can
		// use it later
		c.Set(authContextKey, auth)
		return next(c)
	}
}
