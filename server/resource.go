// Package server provides support for the REST "resource"
package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/valocode/bubbly/api/core"
)

// PostResource godoc
// @Summary Takes a POST request to upload a new resource to the in memory database
// @Description ATM this will only accept one resource per request
// @ID Post-resource
// @Tags resource
// @Param resource body core.ResourceBlockJSON true "Resource Body"
// @Accept  json
// @Produce  json
// @Success 200 {object} apiResponse
// @Failure 400 {object} apiResponse
// @Router /resource [post]
func (s *Server) PostResource(c echo.Context) error {
	// read the resource into a ResourceBlockJSON which keeps the spec{} block
	// as bytes
	var resJSON core.ResourceBlockJSON
	binder := &echo.DefaultBinder{}
	if err := binder.BindBody(c, &resJSON); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// get the ResourceBlock from the JSON representation. We don't actually
	// need the ResourceBlock right now but this is just to validate that the
	// received resource is correctly formatted and to get the resource ID
	// If it fails, return an error code to show it
	_, err := resJSON.ResourceBlock()
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("invalid resource: %w", err))
	}

	dBytes, err := json.Marshal(core.DataBlocks{resJSON.Data()})
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("failed to marshal: %w", err))
	}

	auth := s.getAuthFromContext(c)
	if err := s.Client.PostResource(s.bCtx, auth, dBytes); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, &Status{"uploaded"})
}

// RunResource godoc
// @Summary Takes a POST request to run a named `run` resource, using content
// provided by a multipart form in the run if provided
// @Description Will run the `run` resource specified by the provided name parameter.
// Any inputs required by the resource should be provided within the POST request.
// @ID Run-resource
// @Tags resource,run
// @Param name path string true "Run Resource Name"
// @Accept  mpfd
// @Produce  json
// @Success 200 {object} apiResponse
// @Failure 400 {object} apiResponse
// @Failure 415 {object} apiResponse
// @Failure 500 {object} apiResponse
// @Router /run/{name} [post]
func (s *Server) RunResource(c echo.Context) error {

	workerRun, err := ProcessRunData(s.bCtx, c)
	if err != nil {
		if err == http.ErrNotMultipart {
			return echo.NewHTTPError(http.StatusUnsupportedMediaType, fmt.Errorf("content must be of type %s", echo.MIMEMultipartForm))
		}
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("error while processing content: %w", err))
	}
	workerRun.Name = c.Param("name")

	workerRunBytes, err := json.Marshal(workerRun)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("error marshalling workerRun: %w", err))
	}

	auth := s.getAuthFromContext(c)
	if err := s.Client.PostResourceToWorker(s.bCtx, auth, workerRunBytes); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("error publishing run content to worker: %w", err))
	}

	return c.JSON(http.StatusOK, "run request sent to the Bubbly Worker")
}

// GetResource godoc
// @Summary GetResource Fetches a resource via GET
// @Description Will fetch a resource based on the given ID
// @ID Get-resource
// @Tags resource
// @Param id path string true "Resource ID"
// @Accept  json
// @Produce  json
// @Success 200 {object} apiResponse
// @Failure 400 {object} apiResponse
// @Router /resource/{id} [get]
func (s *Server) GetResource(c echo.Context) error {
	resBlock := core.ResourceBlock{
		ResourceName: c.Param("name"),
		Metadata:     &core.Metadata{},
		ResourceKind: c.Param("kind"),
	}

	auth := s.getAuthFromContext(c)
	resultBytes, err := s.Client.GetResource(s.bCtx, auth, resBlock.String())
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("error getting resource: %w", err))
	}

	return c.JSONBlob(http.StatusOK, resultBytes)
}
