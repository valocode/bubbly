package server

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/valocode/bubbly/agent/component"
)

// authorize godoc
// @Summary Authorizes the request against the bubbly auth service
// @ID authorize
// @Tags authorize
// @Accept json
// @Produce json
// @Success 200 {object} apiResponse
// @Failure 400 {object} apiResponse
// @Router /authorize [post]
func (s *Server) authorize(c echo.Context) error {
	auth := s.getAuthFromContext(c)
	// Create a new request to get the authorization information
	return c.JSON(http.StatusOK, auth)
}

func (s *Server) getOrganizations(c echo.Context) error {
	return s.handleAuthServiceRequest(c, http.MethodGet, "/organizations")
}

func (s *Server) createOrganization(c echo.Context) error {
	var (
		auth = s.getAuthFromContext(c)
		org  = make(map[string]string)
	)
	// We cannot read the body from the request io.Reader twice, so let's read
	// it once and then copy it back for the auth service handler to also read it
	bodyBytes, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "failed to read body of request")
	}
	// Have to close the body immediately so that we can re-write to it
	c.Request().Body.Close()
	c.Request().Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
	if err := json.Unmarshal(bodyBytes, &org); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "error decoding organisation: "+err.Error())
	}

	if err := s.handleAuthServiceRequest(c, http.MethodPost, "/organization/new"); err != nil {
		// If the request returned an error, the organization/tenant was not created
		// so just pass through the error
		return err
	}
	// Else, the organization was successfully created so we need to tell the
	// store to create the tenant!
	if err := s.Client.CreateTenant(s.bCtx, auth, org["name"]); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "error creating tenant: "+err.Error())
	}
	return nil
}

func (s *Server) existsOrganization(c echo.Context) error {
	return s.handleAuthServiceRequest(c, http.MethodGet, "/organization/exists")
}

func (s *Server) getUsersInOrganization(c echo.Context) error {
	return s.handleAuthServiceRequest(c, http.MethodGet, "/users")
}

func (s *Server) inviteUserByEmail(c echo.Context) error {
	return s.handleAuthServiceRequest(c, http.MethodPost, "/user/invite")
}

func (s *Server) deleteUser(c echo.Context) error {
	return s.handleAuthServiceRequest(c, http.MethodPost, "/user/delete")
}

func (s *Server) setUserRole(c echo.Context) error {
	return s.handleAuthServiceRequest(c, http.MethodPost, "/user/role")
}

func (s *Server) getUserRole(c echo.Context) error {
	return s.handleAuthServiceRequest(c, http.MethodGet, "/user/role")
}

func (s *Server) getUserTokens(c echo.Context) error {
	return s.handleAuthServiceRequest(c, http.MethodGet, "/user/tokens")
}

func (s *Server) createUserToken(c echo.Context) error {
	return s.handleAuthServiceRequest(c, http.MethodPost, "/user/token/new")
}

func (s *Server) deleteUserToken(c echo.Context) error {
	return s.handleAuthServiceRequest(c, http.MethodPost, "/user/token/delete")
}

func (s *Server) handleAuthServiceRequest(c echo.Context, method string, path string) error {
	var (
		auth   = s.getAuthFromContext(c)
		client = &http.Client{}
	)
	// Make a new request simply passing in the body of the original request as the
	// body for the new request
	req, err := http.NewRequest(method, s.bCtx.AuthConfig.AuthAddr+path, c.Request().Body)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to create auth service request: "+err.Error())
	}
	q := req.URL.Query()
	q.Add("user_id", auth.UserID)
	if auth.Organization != "" {
		q.Add("organization", auth.Organization)
	}
	req.URL.RawQuery = q.Encode()

	// Add the Content-Type header to the forwarded request
	if contentType := c.Request().Header.Get(echo.HeaderContentType); contentType != "" {
		req.Header.Add(echo.HeaderContentType, contentType)
	}

	respBytes, err := doOrganizationQuery(client, req)
	if err != nil {
		return err
	}
	return c.JSONBlob(http.StatusOK, respBytes)
}

func doOrganizationQuery(client *http.Client, req *http.Request) ([]byte, error) {
	resp, err := client.Do(req)
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusInternalServerError, "failed to get user organizations: "+err.Error())
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, echo.NewHTTPError(resp.StatusCode, "received non-200 HTTP status code from auth service")
	}

	// Read and forward the bytes
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusInternalServerError, "failed to read response from auth service: "+err.Error())
	}
	return respBytes, nil
}

const authContextKey = "auth"

func (s *Server) getAuthFromContext(c echo.Context) *component.MessageAuth {
	if !s.bCtx.AuthConfig.Authentication {
		return nil
	}
	auth := c.Get(authContextKey).(component.MessageAuth)
	return &auth
}
