package server

import (
	"net/http"

	"github.com/rs/zerolog/log"

	"github.com/gin-gonic/gin"
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
func Query(c *gin.Context) {
	var query queryReq
	if bindErr := c.ShouldBindJSON(&query); bindErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": bindErr.Error()})
		return
	}

	results, queryErr := serverStore.Store.Query(query.Query)
	if queryErr != nil {
		log.Error().Msg(queryErr.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": queryErr.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": results,
	})
}
