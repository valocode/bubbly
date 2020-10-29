package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zclconf/go-cty/cty/json"
)

type Data struct {
	Name   string      `json:"name" binding:"required"`
	Fields []DataField `json:"fields"`
	Table  []Data      `json:"table"`
}

type DataField struct {
	Name  string               `json:"name"`
	Value json.SimpleJSONValue `json:"value"`
}

func upload(c *gin.Context) {
	var data Data
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	// TODO when Uploader is ready, send data here

	c.JSON(http.StatusOK, gin.H{
		"status": "uploaded",
	})
}
