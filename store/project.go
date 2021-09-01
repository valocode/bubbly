package store

import (
	"github.com/valocode/bubbly/ent"
	"github.com/valocode/bubbly/store/api"
)

func (h *Handler) CreateProject(req *api.ProjectCreateRequest) (*ent.Project, error) {
	if err := h.validator.Struct(req); err != nil {
		return nil, HandleValidatorError(err, "create project")
	}
	dbProject, err := h.client.Project.Create().
		SetModelCreate(req.Project).
		SetOwnerID(h.orgID).
		Save(h.ctx)
	if err != nil {
		return nil, HandleEntError(err, "creating project")
	}
	return dbProject, nil
}
