package server

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/ziflex/lecho/v2"
	"golang.org/x/sync/errgroup"

	"github.com/valocode/bubbly/client"
	"github.com/valocode/bubbly/config"
	"github.com/valocode/bubbly/env"
)

type Server struct {
	Config *config.ServerConfig
	Server *http.Server
	Client client.Client
	bCtx   *env.BubblyContext
}

func New(bCtx *env.BubblyContext) (*Server, error) {
	client, err := client.New(bCtx)
	if err != nil {
		return nil, fmt.Errorf("failed to create NATS client: %w", err)
	}
	// create the http server
	server := &Server{
		Config: bCtx.ServerConfig,
		Server: &http.Server{
			// TODO: maybe we should use the bCtx Host here, unless it's localhost?
			Addr: fmt.Sprintf(":%s", bCtx.ServerConfig.Port),
		},
		Client: client,
		bCtx:   bCtx,
	}

	server.Server.Handler = server.setupRouter()

	return server, nil
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
			Msg("API server error")

		// Call the default handler to return the HTTP response
		router.DefaultHTTPErrorHandler(err, c)
	}

	// Initialize HTTP Routes
	s.initializeRoutes(router)

	return router
}

func (s *Server) ListenAndServe() error {
	var g errgroup.Group
	g.Go(func() error {
		if err := s.Server.ListenAndServe(); err != nil && err != http.
			ErrServerClosed {
			return err
		}
		return nil
	})

	return g.Wait()
}

func (s *Server) Close() {
	s.Client.Close()
}
