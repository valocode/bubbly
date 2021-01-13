package server

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/ziflex/lecho/v2"
	"golang.org/x/sync/errgroup"

	"github.com/verifa/bubbly/config"
	"github.com/verifa/bubbly/env"
	"github.com/verifa/bubbly/store"
)

var serverStore *store.Store

type Server struct {
	Config *config.ServerConfig
	Server *http.Server
}

func New(bCtx *env.BubblyContext) *Server {
	// create the http server
	a := &Server{
		Config: bCtx.ServerConfig,
		Server: &http.Server{
			// TODO: maybe we should use the bCtx Host here, unless it's localhost?
			Addr: fmt.Sprintf(":%s", bCtx.ServerConfig.Port),
		},
	}

	a.Server.Handler = a.setupRouter(bCtx)

	return a
}

// SetupRouter returns a pointer to an instance of our echo server with all the
// routes initialized
func (a *Server) setupRouter(bCtx *env.BubblyContext) *echo.Echo {
	// Initialize Router
	router := echo.New()
	router.Logger = lecho.From(*bCtx.Logger)
	router.Use(
		middleware.Recover(),
		middleware.RequestID(), // Generate a request IDs
		VersionMiddleware,
		// Setup CORS middleware to allow local docker-compose setup.
		// TODO: make this more restrictive (default "*") or come up with an
		// approach to avoid this
		middleware.CORS(),
	)
	// setup the error handler
	router.HTTPErrorHandler = func(err error, c echo.Context) {
		// TODO: send this to some telemetry/logging service...

		// Log the output for the user
		bCtx.Logger.Error().
			Str("Path", c.Path()).
			Strs("QueryParams", c.ParamValues()).
			Err(err).
			Msg("received an error")

		// Call the default handler to return the HTTP response
		router.DefaultHTTPErrorHandler(err, c)
	}

	// Initialize HTTP Routes
	a.initializeRoutes(bCtx, router)

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

func (a *Server) ListenAndServe(bCtx *env.BubblyContext) error {
	var g errgroup.Group
	g.Go(func() error {
		if err := a.Server.ListenAndServe(); err != nil && err != http.
			ErrServerClosed {
			return err
		}
		return nil
	})

	return g.Wait()
}
