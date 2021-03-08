package server

import (
	"net/http"

	"github.com/labstack/echo/v4"

	echoSwagger "github.com/swaggo/echo-swagger"
)

// initializeRoutes Builds the endpoints and grouping for a gin router
func (s *Server) initializeRoutes(router *echo.Echo) {
	// Keep Alive Test
	router.GET("/healthz", func(c echo.Context) error {
		return c.String(http.StatusOK, "pong")
	})

	api := router.Group("/api/v1")

	// If multitenancy is enabled, add it the URL path
	if s.bCtx.AuthConfig.MultiTenancy {
		api.GET("/organizations", s.getOrganizations)
		api.POST("/organization/new", s.createOrganization)
		api.GET("/organization/exists", s.existsOrganization)
		// Add base routes before adding organization level
		api = api.Group("/o/:organization")
		api.GET("/authorize", s.authorize)
		api.GET("/users", s.getUsersInOrganization)
		api.POST("/user/invite", s.inviteUserByEmail)
		api.POST("/user/delete", s.deleteUser)
		api.GET("/user/role", s.getUserRole)
		api.POST("/user/role", s.setUserRole)
		api.GET("/user/tokens", s.getUserTokens)
		api.POST("/user/token/new", s.createUserToken)
		api.POST("/user/token/delete", s.deleteUserToken)
	}

	api.POST("/run/:name", s.RunResource)
	api.POST("/resource", s.PostResource)
	api.GET("/resource/:kind/:name", s.GetResource)
	api.POST("/graphql", s.Query)
	api.POST("/schema", s.PostSchema)
	api.POST("/upload", s.upload)

	// Serve Swagger files
	router.GET("/swagger/*", echoSwagger.WrapHandler)
}
