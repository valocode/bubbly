package main

import (
	"github.com/verifa/bubbly/server"

	"github.com/gin-gonic/gin"
)

var router *gin.Engine
var BUBBLY_API_PORT = ":8080"

// The starting point for launching Bubbly's web API
func main() {
	// initialize the router's endpoints
	router = server.SetupRouter()

	// Execute Bubbly server (default port is 8080)
	router.Run(BUBBLY_API_PORT)
}
