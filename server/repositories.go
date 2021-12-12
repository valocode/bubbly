package server

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/valocode/bubbly/ent"
	"github.com/valocode/bubbly/store"
	"github.com/valocode/bubbly/store/api"
)

func (s *Server) getRepositories(c echo.Context) error {
	var req api.RepositoryGetRequest
	if err := (&echo.DefaultBinder{}).Bind(&req, c); err != nil {
		return err
	}
	h, err := store.NewHandler(store.WithStore(s.store))
	if err != nil {
		return err
	}
	var where ent.RepositoryWhereInput
	if req.Name != "" {
		where.Name = &req.Name
	}
	resp, err := h.GetRepositories(&store.RepositoryQuery{Where: &where})
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, resp)
}

func (s *Server) postRepository(c echo.Context) error {
	var req api.RepositoryCreateRequest
	if err := (&echo.DefaultBinder{}).Bind(&req, c); err != nil {
		return err
	}
	h, err := store.NewHandler(store.WithStore(s.store))
	if err != nil {
		return err
	}
	resp, err := h.CreateRepository(&req)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, resp)
}
