package server

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/valocode/bubbly/ent"
	"github.com/valocode/bubbly/store"
	"github.com/valocode/bubbly/store/api"
)

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
	if req.Repositories != "" {
		where.HasCommitWith = append(where.HasCommitWith, &ent.GitCommitWhereInput{
			HasRepositoryWith: []*ent.RepositoryWhereInput{{NameIn: strings.Split(req.Repositories, ",")}},
		})
	}
	if req.Projects != "" {
		where.HasCommitWith = append(where.HasCommitWith, &ent.GitCommitWhereInput{
			HasRepositoryWith: []*ent.RepositoryWhereInput{{HasProjectWith: []*ent.ProjectWhereInput{{NameIn: strings.Split(req.Projects, ",")}}}},
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

	releases, err := h.GetReleases(&store.ReleaseQuery{
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
		Releases: releases,
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

	releases, err := h.GetReleases(&store.ReleaseQuery{
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
	if len(releases) == 0 {
		return store.NewNotFoundError(nil, "release with id %d not found", req.ID)
	}
	return c.JSON(http.StatusOK, api.ReleaseGetByIDResponse{
		Release: releases[0],
	})
}

func (s *Server) postRelease(c echo.Context) error {
	var req api.ReleaseCreateRequest
	if err := (&echo.DefaultBinder{}).Bind(&req, c); err != nil {
		return err
	}
	// binder := &echo.DefaultBinder{}
	// if err := binder.BindBody(c, &req); err != nil {
	// 	return err
	// }
	h, err := store.NewHandler(
		store.WithStore(s.store),
	)
	if err != nil {
		return err
	}
	if _, err := h.CreateRelease(&req); err != nil {
		return err
	}
	return c.NoContent(http.StatusOK)
}
