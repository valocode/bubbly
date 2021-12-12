package server

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/valocode/bubbly/ent"
	"github.com/valocode/bubbly/store"
	"github.com/valocode/bubbly/store/api"
)

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
	resp, err := h.GetProjects(&store.ProjectQuery{Where: &where})
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, resp)
}

func (s *Server) postProject(c echo.Context) error {
	var req api.ProjectCreateRequest
	if err := (&echo.DefaultBinder{}).Bind(&req, c); err != nil {
		return err
	}
	h, err := store.NewHandler(store.WithStore(s.store))
	if err != nil {
		return err
	}
	resp, err := h.CreateProject(&req)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, resp)
}
