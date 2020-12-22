package server

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/verifa/bubbly/env"
	"github.com/verifa/bubbly/store"
	"github.com/ziflex/lecho/v2"
	"golang.org/x/sync/errgroup"
)

var serverStore *store.Store

// SetupRouter returns a pointer to a gin engine after setting up middleware
// and initializing routes
func setupRouter(bCtx *env.BubblyContext) *echo.Echo {
	// Initialize Router
	router := echo.New()
	router.Logger = lecho.From(*bCtx.Logger)
	router.Use(
		middleware.Recover(),
		middleware.RequestID(), // Generate a request IDs
		VersionMiddleware,
	)
	// setup the error handler
	router.HTTPErrorHandler = func(err error, c echo.Context) {
		// Should send this to some telemetry/logging service...
		bCtx.Logger.Error().
			Str("Path", c.Path()).
			Strs("QueryParams", c.ParamValues()).
			Err(err).
			Msg("Received an error")

		// Call the default handler to return the HTTP response
		router.DefaultHTTPErrorHandler(err, c)
	}

	// Initialize HTTP Routes
	InitializeRoutes(bCtx, router)

	return router
}

func InitStore(bCtx *env.BubblyContext) error {
	var err error
	serverStore, err = store.New(bCtx)
	if err != nil {
		return fmt.Errorf("failed to create store: %w", err)
	}
	return nil
}

// GetStore returns a pointer to the DB
func GetStore() *store.Store {
	return serverStore
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
