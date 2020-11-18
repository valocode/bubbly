package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"github.com/verifa/bubbly/env"
)

// InitializeRoutes Builds the endpoints and grouping for a gin router
func InitializeRoutes(bCtx *env.BubblyContext, router *gin.Engine) {
	// Keep Alive Test
	router.GET("/healthz", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	api := router.Group("/api")
	{
		api.POST("/resource", func(c *gin.Context) {
			PostResource(bCtx, c)
		})

		api.GET("/resource/:namespace/:kind/:name", func(c *gin.Context) {
			GetResource(bCtx, c)
		})

		api.POST("/graphql", func(c *gin.Context) {
			Query(bCtx, c)
		})
	}

	// API level versioning
	// Establish grouping rules for versioning
	v1 := router.Group("/v1")
	{
		v1.GET("/status", status)
		v1.GET("/version", versionHandler)
	}

	// Resource Level Versioning
	alpha1 := router.Group("/alpha1")
	{
		alpha1.POST("/upload", func(c *gin.Context) {

			upload(bCtx, c)
		})
	}

	// Serve Swagger files
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}

func status(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "alive",
	})
}
