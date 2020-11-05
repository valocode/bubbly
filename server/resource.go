// Package server provides support for the REST "resource"
package server

import (
	"encoding/json"
	"net/http"
	"os"
	"path"

	"github.com/rs/zerolog/log"

	"github.com/gin-gonic/gin"
	"github.com/tidwall/buntdb"
	"github.com/verifa/bubbly/api/core"
)

const defaultNamespace = "default"

// returns an index for ensuring unique names
func dbIndexName() string {
	return "unique_name"
}

// DbPath returns the path to the DB
func DbPath() string {
	url := "test.db"
	e, err := os.Executable()
	if err != nil {
		panic(err)
	}
	dbpath := path.Join(path.Dir(e), url)
	return dbpath
}

type resourceMap map[string]map[string]map[string]interface{}

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
func PostResource(c *gin.Context) {
	var resourceMap map[string]map[string]map[string]map[string]interface{}
	if err := c.ShouldBindJSON(&resourceMap); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "Resource not defined"})
		return
	}

	// get the resource name
	var resourceName string
	for k := range resourceMap["resource"][resourceKind] {
		resourceName = k
	}
	if resourceName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Resource Name not defined"})
		return
	}

	// If the namespace is not specified, it will default as defaultNamespace
	rawMetadata, _ := json.Marshal(resourceMap["resource"][resourceKind][resourceName]["metadata"])
	var metadata core.Metadata
	err := json.Unmarshal(rawMetadata, &metadata)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "metadata not present"})
		return
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

	if err := uploadResource(&resource); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "uploaded",
	})
}

// Uploads the resource to the in-mem db
func uploadResource(resource *core.ResourceJSON) error {
	db, dbErr := buntdb.Open(DbPath())
	if dbErr != nil {
		log.Error().Msg(dbErr.Error())
		return dbErr
	}

	db.CreateIndex(dbIndexName(), "*", buntdb.IndexJSON("name"))
	db.Update(func(tx *buntdb.Tx) error {
		tx.Set(resource.GetID(), resource.Resource, nil)
		return nil
	})

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
func GetResource(c *gin.Context) {
	resource := core.ResourceJSON{
		Name:      c.Param("name"),
		Namespace: c.Param("namespace"),
		Kind:      c.Param("kind"),
	}

	db, dbErr := buntdb.Open(DbPath())
	if dbErr != nil {
		log.Error().Msg(dbErr.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": dbErr.Error()})
		return
	}

	db.CreateIndex(dbIndexName(), "*", buntdb.IndexJSON("name"))
	var resourceString string
	db.View(func(tx *buntdb.Tx) error {
		val, err := tx.Get(resource.GetID())
		if err != nil {
			return err
		}
		resourceString = val
		return nil
	})

	c.Data(http.StatusOK, "application/json", []byte(resourceString))
}
