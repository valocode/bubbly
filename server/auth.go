package server

import (
	"io/ioutil"
	"net/http"

	"github.com/labstack/echo/v4"
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
	var reqAuth = c.Get(requestAuthKey).(requestAuth)
	// Create a new request to get the authorization information
	return c.JSON(http.StatusOK, reqAuth)
}

func (s *Server) getOrganizations(c echo.Context) error {
	return s.handleAuthServiceRequest(c, http.MethodGet, "/organizations")
}

func (s *Server) createOrganization(c echo.Context) error {
	return s.handleAuthServiceRequest(c, http.MethodPost, "/organization/new")
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
		reqAuth = c.Get(requestAuthKey).(requestAuth)
		client  = &http.Client{}
	)
	req, err := http.NewRequest(method, s.bCtx.AuthConfig.AuthAddr+path, c.Request().Body)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to create auth service request: "+err.Error())
	}
	q := req.URL.Query()
	q.Add("user_id", reqAuth.UserID)
	if reqAuth.Organization != "" {
		q.Add("organization", reqAuth.Organization)
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
