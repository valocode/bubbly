package server

import (
	"net/http"

	"github.com/rs/zerolog/log"

	"github.com/gin-gonic/gin"
	"github.com/verifa/bubbly/api/core"
)

type uploadStruct struct {
	Data core.DataBlocks `json:"data"`
}

func upload(c *gin.Context) {
	var upload uploadStruct
	if err := c.ShouldBindJSON(&upload); err != nil {
		log.Error().Msg(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	if upload.Data == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "malformed request"})
	}

	importErr := db.Import(upload.Data)
	if importErr != nil {
		log.Error().Msg(importErr.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": importErr.Error()})
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "uploaded",
	})
}
