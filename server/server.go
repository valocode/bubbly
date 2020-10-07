package server

import (
	"github.com/gin-gonic/gin"
)

var router *gin.Engine

func SetupRouter() *gin.Engine {
	// Enable Production Mode
	// gin.SetMode(gin.ReleaseMode)

	// Initialize Router
	router := gin.Default()

	// Initialize HTTP Routes
	InitializeRoutes(router)

	return router
}
