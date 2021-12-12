package server

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/valocode/bubbly/store"
	"github.com/valocode/bubbly/store/api"
)

func (s *Server) postAnalysis(c echo.Context) error {
	var req api.AnalysisCreateRequest
	if err := (&echo.DefaultBinder{}).Bind(&req, c); err != nil {
		return err
	}
	h, err := store.NewHandler(store.WithStore(s.store))
	if err != nil {
		return err
	}
	if err := h.CreateAnalysis(&req); err != nil {
		return err
	}
	return c.NoContent(http.StatusOK)
}
