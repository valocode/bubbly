package server

import (
	"context"
	"fmt"
	"io/fs"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/valocode/bubbly/auth"
	"github.com/valocode/bubbly/config"
	"github.com/valocode/bubbly/ent"
	"github.com/valocode/bubbly/store"
	"github.com/valocode/bubbly/store/api"
	"github.com/ziflex/lecho/v2"
)

func New(bCtx *config.BubblyConfig) (*Server, error) {
	store, err := store.New(bCtx)
	if err != nil {
		return nil, fmt.Errorf("error initializing store: %w", err)
	}
	return NewWithStore(bCtx, store)
}

func NewWithStore(bCtx *config.BubblyConfig, store *store.Store) (*Server, error) {
	var (
		e = echo.New()
		s = Server{
			store:     store,
			bCtx:      bCtx,
			e:         e,
			validator: newValidator(),
		}
	)
	// Create an echo logger from our existing zerolog
	eLogger := lecho.From(bCtx.Logger,
		lecho.WithTimestamp(),
	)
	e.Logger = eLogger
	e.Use(
		middleware.Recover(),
		middleware.RequestID(), // Generate a request IDs
		lecho.Middleware(lecho.Config{
			Logger: eLogger,
		}),
		middleware.CORS(),
	)

	// Setup the error handler
	e.HTTPErrorHandler = func(err error, c echo.Context) {
		// TODO: send this to some telemetry/logging service...
		httpError := httpErrorHandler(s.bCtx, err, c)

		// Call the default handler to return the HTTP response
		e.DefaultHTTPErrorHandler(httpError, c)
	}

	authProvider, err := auth.NewProvider(context.TODO(), bCtx.AuthConfig)
	if err != nil {
		return &s, err
	}

	e.Use(authProvider.EchoMiddleware())

	e.GET("/auth/token", authProvider.EchoAuthorizeHandler())

	//
	// Setup the Bubbly UI
	//
	if bCtx.ServerConfig.UI {
		if bCtx.UI == nil {
			return nil, fmt.Errorf("cannot run the bubbly UI with a nil filesystem")
		}
		buildDir, err := fs.Sub(bCtx.UI, "build")
		if err != nil {
			log.Fatalf("creating sub filesystem for ui: %s", err.Error())
		}

		ui := e.Group("/ui")
		ui.Use(middleware.StaticWithConfig(middleware.StaticConfig{
			Filesystem: http.FS(buildDir),
			HTML5:      true,
		}))
	}

	// Keep Alive Test
	e.GET("/healthz", func(c echo.Context) error {
		return c.String(http.StatusOK, "still standin'")
	})

	v1 := e.Group("/api/v1")
	v1.GET("/events", s.getEvents)
	v1.GET("/projects", s.getProjects)
	v1.POST("/projects", s.postProject)
	v1.GET("/releases", s.getReleases)
	v1.GET("/releases/:id", s.getReleaseByID)
	v1.POST("/releases", s.postRelease)
	v1.POST("/artifacts", s.postArtifact)
	v1.POST("/codescans", s.postCodeScan)
	v1.POST("/testruns", s.postTestRun)
	v1.POST("/adapters", s.postAdapter)
	v1.GET("/adapters", s.getAdapters)
	v1.POST("/policies", s.postPolicy)
	v1.PUT("/policies/:id", s.putPolicy)
	v1.GET("/policies", s.getPolicies)

	v1.GET("/components", s.getComponents)

	v1.GET("/vulnerabilityreviews", s.getVulnerabilityReviews)
	v1.POST("/vulnerabilityreviews", s.postVulnerabilityReview)
	v1.PUT("/vulnerabilityreviews/:id", s.putVulnerabilityReview)

	// Authentication
	// /login endpoint that redirects to OIDC provider, which provides access token
	// v1.POST("/users", nil)
	// v1.PUT("/users/me", nil)
	// v1.PUT("/users/:id", nil)
	// v1.GET("/users/me", nil)
	// v1.GET("/users/:id", nil)
	// v1.DELETE("/users/:id", nil)

	// v1.POST("/groups", nil)
	// v1.POST("/groups/:id/users", nil) // Add users to a group
	// v1.POST("/groups/:id/roles", nil) // Add roles to a group
	// v1.PUT("/groups/:id", nil)
	// v1.GET("/groups/:id", nil)
	// v1.DELETE("/groups/:id", nil)

	// Initialise monitoring services
	if err := s.initMonitoring(); err != nil {
		return nil, fmt.Errorf("initializing monitoring: %w", err)
	}

	return &s, nil
}

type Server struct {
	store     *store.Store
	bCtx      *config.BubblyConfig
	e         *echo.Echo
	validator *validator.Validate
}

func (s *Server) Start() error {
	addr := fmt.Sprintf("%s:%s", s.bCtx.ServerConfig.Host, s.bCtx.ServerConfig.Port)
	if err := s.e.Start(addr); err != http.ErrServerClosed {
		return err
	}
	return nil
}

func (s *Server) getEvents(c echo.Context) error {
	var req api.EventGetRequest
	if err := (&echo.DefaultBinder{}).BindQueryParams(c, &req); err != nil {
		return err
	}

	h, err := store.NewHandler(store.WithStore(s.store))
	if err != nil {
		return err
	}
	dbEvents, err := h.GetEvents(&req)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, api.EventGetResponse{
		Events: dbEvents,
	})
}

func (s *Server) getProjects(c echo.Context) error {
	var req api.ProjectGetRequest
	if err := (&echo.DefaultBinder{}).Bind(&req, c); err != nil {
		return err
	}
	h, err := store.NewHandler(store.WithStore(s.store))
	if err != nil {
		return err
	}
	var where ent.ProjectWhereInput
	if req.Name != "" {
		where.Name = &req.Name
	}
	projects, err := h.GetProjects(&store.ProjectQuery{Where: &where})
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, api.ProjectGetResponse{
		Projects: projects,
	})
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
	return c.NoContent(http.StatusOK)
}

func (s *Server) getReleases(c echo.Context) error {
	var req api.ReleaseGetRequest
	if err := (&echo.DefaultBinder{}).Bind(&req, c); err != nil {
		return err
	}

	h, err := store.NewHandler(store.WithStore(s.store))
	if err != nil {
		return err
	}
	var where ent.ReleaseWhereInput
	if req.Commit != "" {
		where.HasCommitWith = append(where.HasCommitWith, &ent.GitCommitWhereInput{
			Hash: &req.Commit,
		})
	}
	if req.Repos != "" {
		where.HasCommitWith = append(where.HasCommitWith, &ent.GitCommitWhereInput{
			HasRepoWith: []*ent.RepoWhereInput{{NameIn: strings.Split(req.Repos, ",")}},
		})
	}
	if req.Projects != "" {
		where.HasCommitWith = append(where.HasCommitWith, &ent.GitCommitWhereInput{
			HasRepoWith: []*ent.RepoWhereInput{{HasProjectWith: []*ent.ProjectWhereInput{{NameIn: strings.Split(req.Projects, ",")}}}},
		})
	}
	if req.HeadOnly {
		where.HasHeadOf = &req.HeadOnly
	}

	var order *api.Order
	if req.SortBy != "" {
		order, err = api.OrderFromSortBy(req.SortBy)
		if err != nil {
			return err
		}
	}

	dbReleases, err := h.GetReleases(&store.ReleaseQuery{
		Where:          &where,
		Order:          order,
		WithLog:        req.Log,
		WithPolicies:   req.Policies,
		WithViolations: true,
	})
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, api.ReleaseGetResponse{
		Releases: dbReleases,
	})
}

func (s *Server) getReleaseByID(c echo.Context) error {
	var req api.ReleaseGetByIDRequest
	if err := (&echo.DefaultBinder{}).Bind(&req, c); err != nil {
		return err
	}

	h, err := store.NewHandler(store.WithStore(s.store))
	if err != nil {
		return err
	}

	dbReleases, err := h.GetReleases(&store.ReleaseQuery{
		Where: &ent.ReleaseWhereInput{
			ID: &req.ID,
		},
		WithLog:        true,
		WithViolations: true,
		WithPolicies:   true,
	})
	if err != nil {
		return err
	}
	if len(dbReleases) == 0 {
		return store.NewNotFoundError(nil, "release with id %d not found", req.ID)
	}
	return c.JSON(http.StatusOK, api.ReleaseGetByIDResponse{
		Release: dbReleases[0],
	})
}

func (s *Server) postRelease(c echo.Context) error {
	var req api.ReleaseCreateRequest
	binder := &echo.DefaultBinder{}
	if err := binder.BindBody(c, &req); err != nil {
		return err
	}
	h, err := store.NewHandler(
		store.WithStore(s.store),
		// TODO(miika): here is how to set user id for handler
		// store.WithUserID(""),
	)
	if err != nil {
		return err
	}
	if _, err := h.CreateRelease(&req); err != nil {
		return err
	}
	return nil
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
	return c.NoContent(http.StatusOK)
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
	return c.NoContent(http.StatusOK)
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
	return c.NoContent(http.StatusOK)
}

func (s *Server) getAdapters(c echo.Context) error {
	var req api.AdapterGetRequest
	if err := (&echo.DefaultBinder{}).Bind(&req, c); err != nil {
		return err
	}
	h, err := store.NewHandler(store.WithStore(s.store))
	if err != nil {
		return err
	}
	var where ent.AdapterWhereInput
	if req.Name != "" {
		where.Name = &req.Name
	}
	if req.Tag != "" {
		where.Tag = &req.Tag
	}
	adapters, err := h.GetAdapters(&store.AdapterQuery{
		Where: &where,
	})
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, api.AdapterGetResponse{
		Adapters: adapters,
	})
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
	return c.NoContent(http.StatusOK)
}

func (s *Server) getPolicies(c echo.Context) error {
	var req api.ReleasePolicyGetRequest
	if err := (&echo.DefaultBinder{}).Bind(&req, c); err != nil {
		return err
	}
	h, err := store.NewHandler(store.WithStore(s.store))
	if err != nil {
		return err
	}
	var where ent.ReleasePolicyWhereInput
	if req.Name != "" {
		where.Name = &req.Name
	}
	policies, err := h.GetReleasePolicies(&store.ReleasePolicyQuery{
		Where:       &where,
		WithAffects: req.WithAffects,
	})
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, api.ReleasePolicyGetResponse{
		Policies: policies,
	})
}

func (s *Server) postPolicy(c echo.Context) error {
	var req api.ReleasePolicySaveRequest
	binder := &echo.DefaultBinder{}
	if err := binder.BindBody(c, &req); err != nil {
		return err
	}
	h, err := store.NewHandler(store.WithStore(s.store))
	if err != nil {
		return err
	}
	if _, err := h.SaveReleasePolicy(&req); err != nil {
		return err
	}
	return c.NoContent(http.StatusOK)
}

func (s *Server) putPolicy(c echo.Context) error {
	var req api.ReleasePolicyUpdateRequest
	if err := (&echo.DefaultBinder{}).Bind(&req, c); err != nil {
		return err
	}
	if err := (&echo.DefaultBinder{}).BindPathParams(c, &req); err != nil {
		return err
	}
	fmt.Println("req: ", req.ID)
	fmt.Println("req: ", c.Param("id"))
	h, err := store.NewHandler(store.WithStore(s.store))
	if err != nil {
		return err
	}
	if _, err := h.UpdateReleasePolicy(&req); err != nil {
		return err
	}
	return c.NoContent(http.StatusOK)
}

func (s *Server) getComponents(c echo.Context) error {
	var req api.ComponentGetRequest
	if err := (&echo.DefaultBinder{}).Bind(&req, c); err != nil {
		return err
	}
	if err := s.validator.Struct(req); err != nil {
		return store.HandleValidatorError(err, "query components")
	}
	h, err := store.NewHandler(store.WithStore(s.store))
	if err != nil {
		return err
	}

	components, err := h.GetComponents(&store.ComponentQuery{})
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, api.ComponentGetResponse{
		Components: components,
	})
}
func (s *Server) getVulnerabilityReviews(c echo.Context) error {
	var req api.VulnerabilityReviewGetRequest
	if err := (&echo.DefaultBinder{}).Bind(&req, c); err != nil {
		return err
	}
	if err := s.validator.Struct(req); err != nil {
		return store.HandleValidatorError(err, "query vulnerability reviews")
	}
	h, err := store.NewHandler(store.WithStore(s.store))
	if err != nil {
		return err
	}
	var (
		where         ent.VulnerabilityReviewWhereInput
		whereVulns    = strings.Split(req.Vulnerabilities, ",")
		whereProjects = strings.Split(req.Projects, ",")
		whereRepos    = strings.Split(req.Repos, ",")
	)

	for _, vulnID := range whereVulns {
		where.HasVulnerabilityWith = append(where.HasVulnerabilityWith,
			&ent.VulnerabilityWhereInput{
				Vid: &vulnID,
			},
		)
	}
	for _, project := range whereProjects {
		where.HasProjectsWith = append(where.HasProjectsWith,
			&ent.ProjectWhereInput{
				Name: &project,
			},
		)
	}
	for _, repo := range whereRepos {
		where.HasReposWith = append(where.HasReposWith,
			&ent.RepoWhereInput{
				Name: &repo,
			},
		)
	}
	if req.Commit != "" {
		// Add where condition for releases with that Hash (should only be one, or none)
		where.HasReleasesWith = append(where.HasReleasesWith, &ent.ReleaseWhereInput{
			HasCommitWith: []*ent.GitCommitWhereInput{{Hash: &req.Commit}},
		})
	}

	vulnReviews, err := h.GetVulnerabilityReviews(&store.VulnerabilityReviewQuery{
		Where:        &where,
		WithProjects: req.WithProjects,
		WithRepos:    req.WithRepos,
	})
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, api.VulnerabilityReviewGetResponse{
		Reviews: vulnReviews,
	})
}

func (s *Server) postVulnerabilityReview(c echo.Context) error {
	var req api.VulnerabilityReviewSaveRequest
	binder := &echo.DefaultBinder{}
	if err := binder.BindBody(c, &req); err != nil {
		return err
	}
	if err := s.validator.Struct(req); err != nil {
		return store.HandleValidatorError(err, "save vulnerability review")
	}
	h, err := store.NewHandler(store.WithStore(s.store))
	if err != nil {
		return err
	}
	if _, err := h.SaveVulnerabilityReview(&req); err != nil {
		return err
	}
	return c.NoContent(http.StatusOK)
}

func (s *Server) putVulnerabilityReview(c echo.Context) error {
	var req api.VulnerabilityReviewUpdateRequest
	if err := (&echo.DefaultBinder{}).Bind(&req, c); err != nil {
		return err
	}
	if err := s.validator.Struct(req); err != nil {
		return store.HandleValidatorError(err, "update vulnerability review")
	}
	h, err := store.NewHandler(store.WithStore(s.store))
	if err != nil {
		return err
	}
	if _, err := h.UpdateVulnerabilityReview(&req); err != nil {
		return err
	}
	return c.NoContent(http.StatusOK)
}
