// This file provides support for the REST "resource"
package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path"

	"github.com/gin-gonic/gin"
	"github.com/tidwall/buntdb"
	"github.com/verifa/bubbly/api/core"
)

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

// PostResource Takes a POST request to upload a new resource to the in memory database
// ATM this will only accept one resource per request
func PostResource(c *gin.Context) {
	var resourceMap map[string]map[string]map[string]interface{}
	if err := c.ShouldBindJSON(&resourceMap); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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
	}

	// get the resource name
	var resourceName string
	for k := range resourceMap["resource"][resourceKind] {
		resourceName = k
	}
	if resourceName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Resource Name not defined"})
	}

	resource := core.ResourceJSON{
		Kind:     resourceKind,
		Name:     resourceName,
		Resource: string(request),
	}

	if err := uploadResouce(&resource); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "uploaded",
	})
}

// Uploads the resource to the in-mem db
func uploadResouce(resource *core.ResourceJSON) error {
	db, dbErr := buntdb.Open(DbPath())
	if dbErr != nil {
		fmt.Println(dbErr)
	}

	db.CreateIndex(dbIndexName(), "*", buntdb.IndexJSON("name"))
	db.Update(func(tx *buntdb.Tx) error {
		tx.Set(resource.Kind+"."+resource.Name, resource.Resource, nil)
		return nil
	})

	return nil
}

// GetResource Fetches a resource via GET
func GetResource(c *gin.Context) {
	db, dbErr := buntdb.Open(DbPath())
	if dbErr != nil {
		fmt.Println(dbErr)
	}

	db.CreateIndex(dbIndexName(), "*", buntdb.IndexJSON("name"))
	var resource string
	db.View(func(tx *buntdb.Tx) error {

		val, err := tx.Get(c.Param("id"))
		if err != nil {
			return err
		}
		resource = val
		return nil
	})

	c.Data(http.StatusOK, "application/json", []byte(resource))
}
