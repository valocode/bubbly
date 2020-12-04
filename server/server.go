package server

import (
	"fmt"
	"os"

	"net/http"

	"github.com/gin-contrib/logger"
	"github.com/gin-gonic/gin"
	"github.com/verifa/bubbly/env"
	"github.com/verifa/bubbly/store"
	"golang.org/x/sync/errgroup"
)

var router *gin.Engine

var serverStore struct {
	*store.Store
}

// SetupRouter returns a pointer to a gin engine after setting up middleware
// and initializing routes
func setupRouter(bCtx *env.BubblyContext) *gin.Engine {
	// Initialize Router
	// router := gin.Default()  // Sets the Gin defaults
	router := gin.New() // Use a blank Gin server with no middleware loaded
	router.Use(logger.SetLogger())
	router.Use(gin.Recovery())
	router.Use(VersionMiddleware())

	// Initialize HTTP Routes
	InitializeRoutes(bCtx, router)

	return router
}

func InitStore() error {
	var err error
	serverStore.Store, err = store.New(store.Config{
		Provider:         store.ProviderType(os.Getenv("PROVIDER")),
		PostgresAddr:     os.Getenv("POSTGRES_ADDR"),
		PostgresUser:     os.Getenv("POSTGRES_USER"),
		PostgresPassword: os.Getenv("POSTGRES_PASSWORD"),
		PostgresDatabase: os.Getenv("POSTGRES_DATABASE"),
	})
	if err != nil {
		return fmt.Errorf("failed to create store: %w", err)
	}
	return nil
}

// GetStore returns a pointer to the DB
func GetStore() *store.Store {
	return serverStore.Store
}

func ListenAndServe(bCtx *env.BubblyContext) error {

	// initialize the router's endpoints
	router := setupRouter(bCtx)

	// create the http server
	serv := &http.Server{
		// TODO: maybe we should use the bCtx Host here, unless it's localhost?
		Addr:    fmt.Sprintf(":%s", bCtx.ServerConfig.Port),
		Handler: router,
	}

	var g errgroup.Group
	g.Go(func() error {
		if err := serv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			return err
		}
		return nil
	})
	return g.Wait()
}
