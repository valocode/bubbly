// Package server provides support for the REST "resource"
package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/verifa/bubbly/api/core"
	"github.com/verifa/bubbly/env"
)

// PostResource godoc
// @Summary Takes a POST request to upload a new resource to the in memory database
// @Description ATM this will only accept one resource per request
// @ID Post-resource
// @Tags resource
// @Param resource body resourceMap true "Resource Body"
// @Accept  json
// @Produce  json
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /api/resource [post]
func PostResource(bCtx *env.BubblyContext, c echo.Context) error {
	// read the resource into a ResourceBlockJSON which keeps the spec{} block
	// as bytes
	resJSON := core.ResourceBlockJSON{}
	if err := c.Bind(&resJSON); err != nil {
		return err
	}

	// get the ResourceBlock from the JSON representation. We don't actually
	// need the ResourceBlock right now but this is just to validate that the
	// received resource is correctly formatted and to get the resource ID
	// If it fails, return an error code to show it
	resBlock, err := resJSON.ResourceBlock()
	if err != nil {
		return fmt.Errorf("failed to unmarshal JSON resource: %w", err)
	}

	if err := uploadResource(bCtx, resBlock.String(), resJSON); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, &Status{"uploaded"})
}

// Uploads the resource to the suitable provider
func uploadResource(bCtx *env.BubblyContext, id string, resJSON core.ResourceBlockJSON) error {
	resBytes, err := json.Marshal(resJSON)
	if err != nil {
		return fmt.Errorf("failed to marshal resource: %w", err)
	}

	return serverStore.PutResource(id, string(resBytes))
}

// GetResource godoc
// @Summary GetResource Fetches a resource via GET
// @Description Will fetch a resource based on the given ID
// @ID Get-resource
// @Tags resource
// @Param id path string true "Resource ID"
// @Accept  json
// @Produce  json
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @x-examples 12345
// @Router /api/resource/{id} [get]
func GetResource(bCtx *env.BubblyContext, c echo.Context) error {
	resBlock := core.ResourceBlock{
		ResourceName: c.Param("name"),
		Metadata:     &core.Metadata{Namespace: c.Param("namespace")},
		ResourceKind: c.Param("kind"),
	}

	val, err := serverStore.GetResource(resBlock.String())
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	resJSON := core.ResourceBlockJSON{}
	if err := json.NewDecoder(val).Decode(&resJSON); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, resJSON)
}
