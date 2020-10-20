package server

import (
	"os"

	"github.com/gin-contrib/logger"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var router *gin.Engine

func SetupRouter() *gin.Engine {
	// Enable Production Mode
	// gin.SetMode(gin.ReleaseMode)

	// Initialize Logger
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if gin.IsDebugging() {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	log.Logger = log.Output(
		zerolog.ConsoleWriter{
			Out:     os.Stderr,
			NoColor: false,
		},
	)

	// Initialize Router
	// router := gin.Default()  // Sets the Gin defaults
	router := gin.New() // Use a blank Gin server with no middleware loaded
	router.Use(logger.SetLogger())

	// Initialize HTTP Routes
	InitializeRoutes(router)

	return router
}
