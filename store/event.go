package store

import (
	"fmt"
	"strconv"

	"github.com/valocode/bubbly/ent"
	"github.com/valocode/bubbly/ent/event"
	"github.com/valocode/bubbly/ent/gitcommit"
	"github.com/valocode/bubbly/ent/project"
	"github.com/valocode/bubbly/ent/release"
	"github.com/valocode/bubbly/ent/repository"
	"github.com/valocode/bubbly/store/api"
)

func (h *Handler) GetEvents(req *api.EventGetRequest) ([]*ent.Event, error) {
	eQuery := h.client.Event.Query()
	if req.Project != "" {
		eQuery.Where(event.HasProjectWith(project.Name(req.Project)))
	}
	if req.Repo != "" {
		eQuery.Where(event.HasRepositoryWith(repository.Name(req.Repo)))
	}
	if req.ReleaseName != "" {
		eQuery.Where(event.HasReleaseWith(release.Name(req.ReleaseName)))
	}
	if req.ReleaseVersion != "" {
		eQuery.Where(event.HasReleaseWith(release.Name(req.ReleaseVersion)))
	}
	if req.Commit != "" {
		eQuery.Where(event.HasReleaseWith(release.HasCommitWith(gitcommit.Hash(req.Commit))))
	}
	// Set a default limit if none given
	if req.Last == "" {
		req.Last = "20"
	}
	// Limit the query and order by time
	limit, err := strconv.Atoi(req.Last)
	if err != nil {
		return nil, fmt.Errorf("\"last\" paramter is not a number: %w", err)
	}
	eQuery.Limit(limit).Order(ent.Desc(event.FieldTime))

	dbEvents, err := eQuery.All(h.ctx)
	if err != nil {
		return nil, HandleEntError(err, "query events")
	}
	return dbEvents, nil
}

func (h *Handler) CreateEvent(req *api.EventSaveRequest) (*ent.Event, error) {
	if err := h.validator.Struct(req); err != nil {
		return nil, HandleValidatorError(err, "create event")
	}

	eventCreate := h.client.Event.Create().SetModelCreate(req.Event)
	if req.ReleaseID != nil {
		eventCreate.SetReleaseID(*req.ReleaseID)
	}
	if req.RepoID != nil {
		eventCreate.SetRepositoryID(*req.RepoID)
	}
	if req.ProjectID != nil {
		eventCreate.SetProjectID(*req.ProjectID)
	}
	dbEvent, err := eventCreate.Save(h.ctx)
	if err != nil {
		return nil, HandleEntError(err, "create event")
	}
	return dbEvent, nil
}
