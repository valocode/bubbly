// Package server provides support for the REST "resource"
package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/zclconf/go-cty/cty"

	"github.com/verifa/bubbly/api/core"
	"github.com/verifa/bubbly/env"
	"github.com/verifa/bubbly/events"
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
	_, err := resJSON.ResourceBlock()
	if err != nil {
		return fmt.Errorf("failed to unmarshal JSON resource: %w", err)
	}

	d, err := resJSON.Data()
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	resource_event := core.Data{
		TableName: core.EventTableName,
		Fields: map[string]cty.Value{
			"status": cty.StringVal(events.ResourceCreated.String()),
			"time":   cty.StringVal(events.TimeNow()),
		},
		// this join means the _id pulled of the _resource row entry will be
		// mapped to the
		// _resource_id column of this row entry in _event
		Joins: []string{core.ResourceTableName},
	}

	res_db := core.DataBlocks{d, resource_event}

	if err := serverStore.Save(res_db); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, &Status{"uploaded"})
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

	resQuery := fmt.Sprintf(`
		{
			%s(id: "%s") {
				name
				kind
				api_version
				metadata
				spec
			}
		}
	`, core.ResourceTableName, resBlock.String())

	result, err := serverStore.Query(resQuery)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	var (
		resJSON  core.ResourceBlockJSON
		inputMap = result.(map[string]interface{})[core.ResourceTableName].([]interface{})
	)
	b, err := json.Marshal(inputMap[0])
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to marshal resource: ", err.Error())
	}
	err = json.Unmarshal(b, &resJSON)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to unmarshal resource: ", err.Error())
	}

	return c.JSON(http.StatusOK, resJSON)
}
