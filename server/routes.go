package server

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func InitializeRoutes(router *gin.Engine) {
	// Keep Alive Test
	router.GET("/healthz", func(c *gin.Context) {
		c.String(200, "pong")
	})

	// API level versioning
	// Establish grouping rules for versioning
	v1 := router.Group("/v1")
	{
		v1.GET("/test", test)
		v1.GET("/status", status)
		v1.GET("/version", version)
	}

	v2 := router.Group("/v2")
	{
		v2.GET("/test", test)
	}

	// Resource Level Versioning
	alpha1 := router.Group("/upload/alpha1")
	{
		alpha1.POST("/", upload)
	}
}

func test(c *gin.Context) {
	c.JSON(200, gin.H{
		"hey": "hey",
	})
}

func version(c *gin.Context) {
	fmt.Println(c.Request.Header["Accept-Version"])
	c.JSON(200, gin.H{
		"version": "0.0.0.1",
	})
}

func status(c *gin.Context) {
	c.JSON(200, gin.H{
		"status": "alive",
	})
}
