package server

import (
	"net/http"

	"github.com/labstack/echo/v4"

	echoSwagger "github.com/swaggo/echo-swagger"
	"github.com/verifa/bubbly/env"
)

// initializeRoutes Builds the endpoints and grouping for a gin router
func (a *Server) initializeRoutes(bCtx *env.BubblyContext,
	router *echo.Echo) {
	// Keep Alive Test
	router.GET("/healthz", func(c echo.Context) error {
		return c.String(http.StatusOK, "pong")
	})

	api := router.Group("/api")
	{
		api.POST("/resource", func(c echo.Context) error {
			return a.PostResource(bCtx, c)
		})

		api.GET("/resource/:namespace/:kind/:name", func(c echo.Context) error {
			return a.GetResource(bCtx, c)
		})

		api.POST("/graphql", func(c echo.Context) error {
			return a.Query(bCtx, c)
		})

		api.POST("/schema", func(c echo.Context) error {
			return a.PostSchema(bCtx, c)
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
		alpha1.POST("/upload", func(c echo.Context) error {
			return a.upload(bCtx, c)
		})
	}

	// Serve Swagger files
	router.GET("/swagger/*", echoSwagger.WrapHandler)
}

func status(c echo.Context) error {
	return c.JSON(http.StatusOK, &Status{"alive"})
}
