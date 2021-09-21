package store

import (
	"github.com/valocode/bubbly/ent"
	"github.com/valocode/bubbly/store/api"
)

type ProjectQuery struct {
	Where *ent.ProjectWhereInput
}

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

func (h *Handler) GetProjects(query *ProjectQuery) ([]*api.Project, error) {
	dbProjects, err := h.client.Project.Query().
		WhereInput(query.Where).
		All(h.ctx)
	if err != nil {
		return nil, HandleEntError(err, "getting projects")
	}
	var projects = make([]*api.Project, 0, len(dbProjects))
	for _, proj := range dbProjects {
		projects = append(projects, &api.Project{
			ProjectModelRead: *ent.NewProjectModelRead().FromEnt(proj),
		})
	}
	return projects, nil
}
