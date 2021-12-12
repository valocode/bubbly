package server

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/valocode/bubbly/ent"
	"github.com/valocode/bubbly/store"
	"github.com/valocode/bubbly/store/api"
)

func (s *Server) getArtifacts(c echo.Context) error {
	var req api.ArtifactGetRequest
	if err := (&echo.DefaultBinder{}).Bind(&req, c); err != nil {
		return err
	}
	h, err := store.NewHandler(store.WithStore(s.store))
	if err != nil {
		return err
	}
	var where ent.ArtifactWhereInput
	if req.Name != "" {
		where.Name = &req.Name
	}
	resp, err := h.GetArtifacts(&store.ArtifactQuery{Where: &where})
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, resp)
}

func (s *Server) postArtifact(c echo.Context) error {
	var req api.ArtifactCreateRequest
	if err := (&echo.DefaultBinder{}).Bind(&req, c); err != nil {
		return err
	}
	h, err := store.NewHandler(store.WithStore(s.store))
	if err != nil {
		return err
	}
	resp, err := h.CreateArtifact(&req)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, resp)
}
