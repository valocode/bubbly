package server

import (
	"fmt"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/valocode/bubbly/env"
	"github.com/valocode/bubbly/gql"
	"github.com/valocode/bubbly/store"
	"github.com/valocode/bubbly/store/api"
)

func ListenAndServe(bCtx *env.BubblyContext, store *store.Store) error {
	e := echo.New()
	// e.Logger = lecho.From(bCtx.Logger)
	e.Use(
		middleware.Recover(),
		middleware.RequestID(), // Generate a request IDs
	)

	// Setup the error handler
	e.HTTPErrorHandler = func(err error, c echo.Context) {
		// TODO: send this to some telemetry/logging service...
		httpError := httpErrorHandler(bCtx, err, c)

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
	srv := handler.NewDefaultServer(gql.NewSchema(store.Client()))
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
	v1.POST("/codescans", func(c echo.Context) error {
		return postCodeScan(c, store)
	})
	v1.POST("/testruns", func(c echo.Context) error {
		return postTestRun(c, store)
	})
	v1.POST("/releases", func(c echo.Context) error {
		return postRelease(c, store)
	})
	v1.POST("/adapters", func(c echo.Context) error {
		return postAdapter(c, store)
	})
	v1.GET("/adapters/:name", func(c echo.Context) error {
		return getAdapter(c, store)
	})

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

	addr := fmt.Sprintf("%s:%s", bCtx.ServerConfig.Host, bCtx.ServerConfig.Port)
	if err := e.Start(addr); err != http.ErrServerClosed {
		return err
	}
	return nil
}

func postCodeScan(c echo.Context, store *store.Store) error {
	var req api.CodeScanRequest
	binder := &echo.DefaultBinder{}
	if err := binder.BindBody(c, &req); err != nil {
		return err
	}
	_, err := store.SaveCodeScan(&req)
	if err != nil {
		return err
	}
	return nil
}

func postTestRun(c echo.Context, store *store.Store) error {
	var req api.TestRunRequest
	binder := &echo.DefaultBinder{}
	if err := binder.BindBody(c, &req); err != nil {
		return err
	}
	_, err := store.SaveTestRun(&req)
	if err != nil {
		return err
	}
	return nil
}

func postRelease(c echo.Context, store *store.Store) error {
	var req api.ReleaseCreateRequest
	binder := &echo.DefaultBinder{}
	if err := binder.BindBody(c, &req); err != nil {
		return err
	}
	_, err := store.CreateRelease(&req)
	if err != nil {
		return err
	}
	return nil
}

func getAdapter(c echo.Context, store *store.Store) error {
	var (
		req  api.AdapterGetRequest
		name = c.Param("name")
		tag  = c.QueryParam("tag")
		ty   = c.QueryParam("type")
	)
	if name != "" {
		req.Name = &name
	}
	if tag != "" {
		req.Tag = &tag
	}
	if ty != "" {
		req.Type = &ty
	}
	adapter, err := store.GetAdapter(&req)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, adapter)
}

func postAdapter(c echo.Context, store *store.Store) error {
	var req api.AdapterSaveRequest
	binder := &echo.DefaultBinder{}
	if err := binder.BindBody(c, &req); err != nil {
		return err
	}
	_, err := store.SaveAdapter(&req)
	if err != nil {
		return err
	}
	return nil
}
