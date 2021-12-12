package server

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/valocode/bubbly/store"
	"github.com/valocode/bubbly/store/api"
)

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

func (s *Server) postComponent(c echo.Context) error {
	var req api.ComponentCreateRequest
	binder := &echo.DefaultBinder{}
	if err := binder.BindBody(c, &req); err != nil {
		return err
	}
	h, err := store.NewHandler(store.WithStore(s.store))
	if err != nil {
		return err
	}
	if _, err := h.CreateComponent(&req); err != nil {
		return err
	}
	return c.NoContent(http.StatusOK)
}
