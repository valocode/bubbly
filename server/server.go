package server

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/ziflex/lecho/v2"
	"golang.org/x/sync/errgroup"

	"github.com/valocode/bubbly/config"
	"github.com/valocode/bubbly/env"
)

type Server struct {
	Config *config.ServerConfig
	Server *http.Server
	bCtx   *env.BubblyContext
}

func New(bCtx *env.BubblyContext) *Server {
	// create the http server
	a := &Server{
		Config: bCtx.ServerConfig,
		Server: &http.Server{
			// TODO: maybe we should use the bCtx Host here, unless it's localhost?
			Addr: fmt.Sprintf(":%s", bCtx.ServerConfig.Port),
		},
		bCtx: bCtx,
	}

	a.Server.Handler = a.setupRouter()

	return a
}

// SetupRouter returns a pointer to an instance of our echo server with all the
// routes initialized
func (s *Server) setupRouter() *echo.Echo {
	// Initialize Router
	router := echo.New()
	router.Logger = lecho.From(*s.bCtx.Logger)
	router.Use(
		middleware.Recover(),
		middleware.RequestID(), // Generate a request IDs
		VersionMiddleware,
		// Setup CORS middleware to allow local docker-compose setup.
		// TODO: make this more restrictive (default "*") or come up with an
		// approach to avoid this
		middleware.CORS(),
	)
	// If multitenancy should be enabled
	if s.bCtx.AuthConfig.Authentication {
		router.Use(s.authMiddleware)
	}

	// Setup the error handler
	router.HTTPErrorHandler = func(err error, c echo.Context) {
		// TODO: send this to some telemetry/logging service...

		// Log the output for the user
		s.bCtx.Logger.Error().
			Str("Path", c.Path()).
			Strs("QueryParams", c.ParamValues()).
			Err(err).
			Msg("received an error")

		// Call the default handler to return the HTTP response
		router.DefaultHTTPErrorHandler(err, c)
	}

	// Initialize HTTP Routes
	s.initializeRoutes(router)

	return router
}

func (a *Server) ListenAndServe() error {
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
