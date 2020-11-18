package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/verifa/bubbly/env"
)

type queryReq struct {
	Query string `json:"query"`
}

// Query godoc
// @Summary Query performs graphql related tasks
// @ID graphql
// @Tags graphql
// @Param query body string true "Query String"
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /api/graphql [post]
func Query(bCtx *env.BubblyContext, c *gin.Context) {
	var query queryReq
	if bindErr := c.ShouldBindJSON(&query); bindErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": bindErr.Error()})
		return
	}

	results, queryErr := serverStore.Store.Query(query.Query)
	if queryErr != nil {
		bCtx.Logger.Error().Msg(queryErr.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": queryErr.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": results,
	})
}
