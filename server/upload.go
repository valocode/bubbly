package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/verifa/bubbly/api/core"
)

func upload(c *gin.Context) {
	var data core.Data

	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	// TODO when Uploader is ready, send data here

	c.JSON(http.StatusOK, gin.H{
		"status": "uploaded",
	})
}
