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
		return
	}
	if upload.Data == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "malformed request"})
		return
	}

	log.Debug().Interface("data", upload).Msg("loading data into intermediary database")

	importErr := db.Import(upload.Data)
	if importErr != nil {
		log.Error().Msg(importErr.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": importErr.Error()})
		return
	}

	// TODO when Uploader is ready, send data here

	c.JSON(http.StatusOK, gin.H{
		"status": "uploaded",
	})
}
