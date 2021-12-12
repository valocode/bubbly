package server

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/valocode/bubbly/store"
	"github.com/valocode/bubbly/store/api"
)

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
