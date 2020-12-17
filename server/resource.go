// Package server provides support for the REST "resource"
package server

import (
	"encoding/json"
	"github.com/labstack/echo/v4"
	"net/http"

	"github.com/verifa/bubbly/api/core"
	"github.com/verifa/bubbly/env"
	"github.com/verifa/bubbly/resource"
)

const defaultNamespace = "default"

// PostResource godoc
// @Summary Takes a POST request to upload a new resource to the in memory database
// @Description ATM this will only accept one resource per request
// @ID Post-resource
// @Tags resource
// @Param resource body resourceMap true "Resource Body"
// @Accept  json
// @Produce  json
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /api/resource [post]
func PostResource(bCtx *env.BubblyContext, c echo.Context) error {
	var resourceMap map[string]map[string]map[string]map[string]interface{}
	if err := c.Bind(&resourceMap); err != nil {
		return c.JSON(http.StatusBadRequest, &Error{err.Error()})
	}

	// The json body that will be stored as core.ResourceJSON.Resource
	request, _ := json.Marshal(resourceMap)

	// get resource kind
	var resourceKind string
	for k := range resourceMap["resource"] {
		for _, item := range core.ResourceKindPriority() {
			if string(item) == k {
				resourceKind = k
			}
		}
	}
	if resourceKind == "" {
		return c.JSON(http.StatusBadRequest, &Error{"Resource not defined"})
	}

	// get the resource name
	var resourceName string
	for k := range resourceMap["resource"][resourceKind] {
		resourceName = k
	}
	if resourceName == "" {
		return c.JSON(http.StatusBadRequest, &Error{"Resource Name not defined"})
	}

	// If the namespace is not specified, it will default as defaultNamespace
	rawMetadata, _ := json.Marshal(resourceMap["resource"][resourceKind][resourceName]["metadata"])
	var metadata core.Metadata
	err := json.Unmarshal(rawMetadata, &metadata)

	if err != nil {
		return c.JSON(http.StatusBadRequest, &Error{"metadata not present"})
	}
	namespace := metadata.Namespace
	if namespace == "" {
		namespace = defaultNamespace
	}

	resource := core.ResourceJSON{
		Kind:      resourceKind,
		Name:      resourceName,
		Namespace: namespace,
		Resource:  string(request),
	}

	err = uploadResource(bCtx, &resource)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &Error{err.Error()})
	}

	return c.JSON(http.StatusOK, &Status{"uploaded"})
}

// Uploads the resource to the in-mem db
func uploadResource(bCtx *env.BubblyContext, r *core.ResourceJSON) error {
	db, err := resource.New(resource.Config{
		// This is hardcoded as "buntdb" for the time so that way bubbly can still be run without extra containers
		// and command line options
		Provider: "buntdb",
	})
	if err != nil {
		return err
	}

	err = db.P.Save(r.GetID(), r.Resource)
	if err != nil {
		return err
	}

	return nil
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
	r := core.ResourceJSON{
		Name:      c.Param("name"),
		Namespace: c.Param("namespace"),
		Kind:      c.Param("kind"),
	}
	db, err := resource.New(resource.Config{
		// This is hardcoded as "buntdb" for the time so that way bubbly can still be run without extra containers
		// and command line options
		Provider: "buntdb",
	})
	if err != nil {
		bCtx.Logger.Error().Msg(err.Error())
		return c.JSON(http.StatusBadRequest, &Error{err.Error()})
	}

	resourceString, err := db.P.Query(r.GetID())
	if err != nil {
		bCtx.Logger.Error().Msg(err.Error())
		return c.JSON(http.StatusBadRequest, &Error{err.Error()})
	}

	return c.String(http.StatusOK, resourceString)
}
