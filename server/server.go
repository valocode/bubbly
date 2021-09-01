package server

import (
	"fmt"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/valocode/bubbly/env"
	"github.com/valocode/bubbly/gql"
	"github.com/valocode/bubbly/store"
	"github.com/valocode/bubbly/store/api"
	"github.com/ziflex/lecho/v2"
)

func New(bCtx *env.BubblyContext) (*Server, error) {
	store, err := store.New(bCtx)
	if err != nil {
		return nil, fmt.Errorf("error initializing store: %w", err)
	}
	return NewWithStore(bCtx, store), nil
}

func NewWithStore(bCtx *env.BubblyContext, store *store.Store) *Server {
	var (
		e = echo.New()
		s = Server{
			store: store,
			bCtx:  bCtx,
			e:     e,
		}
	)
	var logLevel = log.INFO
	if s.bCtx.CLIConfig.Debug {
		logLevel = log.DEBUG
	}
	e.Logger = lecho.New(
		os.Stdout,
		lecho.WithLevel(logLevel),
	)
	e.Use(
		middleware.Recover(),
		middleware.RequestID(), // Generate a request IDs
		middleware.Logger(),
	)

	// Setup the error handler
	e.HTTPErrorHandler = func(err error, c echo.Context) {
		// TODO: send this to some telemetry/logging service...
		httpError := httpErrorHandler(s.bCtx, err, c)

		// Call the default handler to return the HTTP response
		e.DefaultHTTPErrorHandler(httpError, c)
	}

	if true {
		e.Use(middleware.CORS())
	}

	// Keep Alive Test
	e.GET("/healthz", func(c echo.Context) error {
		return c.String(http.StatusOK, "still standin'")
	})

	// Setup graphql query and playground
	srv := handler.NewDefaultServer(gql.NewSchema(s.store.Client()))
	playgroundHandler := playground.Handler("Bubbly", "/graphql")
	e.POST("/graphql", func(c echo.Context) error {
		srv.ServeHTTP(c.Response(), c.Request())
		return nil
	})
	e.GET("/playground", func(c echo.Context) error {
		playgroundHandler.ServeHTTP(c.Response(), c.Request())
		return nil
	})

	v1 := e.Group("/api/v1")
	v1.POST("/projects", s.postProject)
	v1.POST("/releases", s.postRelease)
	v1.GET("/releases", s.getRelease)
	v1.POST("/artifacts", s.postArtifact)
	v1.POST("/codescans", s.postCodeScan)
	v1.POST("/testruns", s.postTestRun)
	v1.POST("/adapters", s.postAdapter)
	v1.GET("/adapters/:name", s.getAdapter)

	// TODO: Miika - integrate Casbin as middleware or similar
	// https://echo.labstack.com/middleware/casbin-auth/

	// Authentication
	// /login endpoint that redirects to OIDC provider, which provides access token

	// Send X-Bubbly-Flavour=oss then auto redirect to default org

	// If SaaS, create a group "/o/:organization" and mount these under there
	v1.POST("/users", nil)
	v1.PUT("/users/me", nil)
	v1.PUT("/users/:id", nil)
	v1.GET("/users/me", nil)
	v1.GET("/users/:id", nil)
	v1.DELETE("/users/:id", nil)

	v1.POST("/groups", nil)
	v1.POST("/groups/:id/users", nil) // Add users to a group
	v1.POST("/groups/:id/roles", nil) // Add roles to a group
	v1.PUT("/groups/:id", nil)
	v1.GET("/groups/:id", nil)
	v1.DELETE("/groups/:id", nil)

	// What are the predifined roles?
	// admins: can add/remove users/groups from their organization
	// users: ??
	// Nice to have:
	// - access control by project (that they are assigned to, either by user or by group)

	// SaaS specific things:
	// - organizations/ PUT,GET, etc

	return &s
}

type Server struct {
	store *store.Store
	bCtx  *env.BubblyContext
	e     *echo.Echo
}

func (s *Server) Start() error {

	addr := fmt.Sprintf("%s:%s", s.bCtx.ServerConfig.Host, s.bCtx.ServerConfig.Port)
	if err := s.e.Start(addr); err != http.ErrServerClosed {
		return err
	}
	return nil
}

func (s *Server) postProject(c echo.Context) error {
	var req api.ProjectCreateRequest
	binder := &echo.DefaultBinder{}
	if err := binder.BindBody(c, &req); err != nil {
		return err
	}
	h, err := store.NewHandler(store.WithStore(s.store))
	if err != nil {
		return err
	}
	if _, err := h.CreateProject(&req); err != nil {
		return err
	}
	return nil
}

func (s *Server) postRelease(c echo.Context) error {
	var req api.ReleaseCreateRequest
	binder := &echo.DefaultBinder{}
	if err := binder.BindBody(c, &req); err != nil {
		return err
	}
	h, err := store.NewHandler(store.WithStore(s.store))
	if err != nil {
		return err
	}
	if _, err := h.CreateRelease(&req); err != nil {
		return err
	}
	return nil
}

func (s *Server) getRelease(c echo.Context) error {
	var req api.ReleaseGetRequest
	if err := (&echo.DefaultBinder{}).BindQueryParams(c, &req); err != nil {
		return err
	}

	h, err := store.NewHandler(store.WithStore(s.store))
	if err != nil {
		return err
	}
	dbRelease, err := h.GetReleases(&req)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, dbRelease)
}

func (s *Server) postArtifact(c echo.Context) error {
	var req api.ArtifactLogRequest
	binder := &echo.DefaultBinder{}
	if err := binder.BindBody(c, &req); err != nil {
		return err
	}
	h, err := store.NewHandler(store.WithStore(s.store))
	if err != nil {
		return err
	}
	if _, err := h.LogArtifact(&req); err != nil {
		return err
	}
	return nil
}

func (s *Server) postCodeScan(c echo.Context) error {
	var req api.CodeScanRequest
	binder := &echo.DefaultBinder{}
	if err := binder.BindBody(c, &req); err != nil {
		return err
	}
	h, err := store.NewHandler(store.WithStore(s.store))
	if err != nil {
		return err
	}
	if _, err := h.SaveCodeScan(&req); err != nil {
		return err
	}
	return nil
}

func (s *Server) postTestRun(c echo.Context) error {
	var req api.TestRunRequest
	binder := &echo.DefaultBinder{}
	if err := binder.BindBody(c, &req); err != nil {
		return err
	}
	h, err := store.NewHandler(store.WithStore(s.store))
	if err != nil {
		return err
	}
	if _, err := h.SaveTestRun(&req); err != nil {
		return err
	}
	return nil
}

func (s *Server) getAdapter(c echo.Context) error {
	var req api.AdapterGetRequest
	if err := (&echo.DefaultBinder{}).BindQueryParams(c, &req); err != nil {
		return err
	}
	h, err := store.NewHandler(store.WithStore(s.store))
	if err != nil {
		return err
	}
	adapter, err := h.GetAdapter(&req)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, adapter)
}

func (s *Server) postAdapter(c echo.Context) error {
	var req api.AdapterSaveRequest
	binder := &echo.DefaultBinder{}
	if err := binder.BindBody(c, &req); err != nil {
		return err
	}
	h, err := store.NewHandler(store.WithStore(s.store))
	if err != nil {
		return err
	}
	if _, err := h.SaveAdapter(&req); err != nil {
		return err
	}
	return nil
}
