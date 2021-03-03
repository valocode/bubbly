package server

import (
	"net/http"

	"github.com/labstack/echo/v4"

	echoSwagger "github.com/swaggo/echo-swagger"
	"github.com/verifa/bubbly/env"
)

// initializeRoutes Builds the endpoints and grouping for a gin router
func (s *Server) initializeRoutes(bCtx *env.BubblyContext,
	router *echo.Echo) {
	// Keep Alive Test
	router.GET("/healthz", func(c echo.Context) error {
		return c.String(http.StatusOK, "pong")
	})

	api := router.Group("/api/v1")
	{
		api.POST("/resource", func(c echo.Context) error {
			return s.PostResource(bCtx, c)
		})

		api.GET("/resource/:kind/:name", func(c echo.Context) error {
			return s.GetResource(bCtx, c)
		})

		api.POST("/graphql", func(c echo.Context) error {
			return s.Query(bCtx, c)
		})

		api.POST("/schema", func(c echo.Context) error {
			return s.PostSchema(bCtx, c)
		})

		api.POST("/upload", func(c echo.Context) error {
			return s.upload(bCtx, c)
		})
	}

	// Serve Swagger files
	router.GET("/swagger/*", echoSwagger.WrapHandler)
}
