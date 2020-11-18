package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/verifa/bubbly/api/core"
	"github.com/verifa/bubbly/env"
)

type uploadStruct struct {
	Data core.DataBlocks `json:"data"`
}

// upload godoc
// @Summary This function will upload core.DataBlocks
// @ID upload data
// @Tags datablocks
// @Param data body uploadStruct true "Datablocks"
// @Accept json
// @Produce json
// @Router /alpha1/upload [post]
func upload(bCtx *env.BubblyContext, c *gin.Context) {
	var upload uploadStruct
	if err := c.ShouldBindJSON(&upload); err != nil {
		bCtx.Logger.Error().Msg(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if upload.Data == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "malformed request"})
		return
	}

	bCtx.Logger.Debug().Interface("data", upload).Msg("loading data into intermediary database")

	importErr := serverStore.Store.Save(upload.Data)
	if importErr != nil {
		bCtx.Logger.Error().Msg(importErr.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": importErr.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "uploaded",
	})
}
