package store

import (
	"fmt"

	"github.com/valocode/bubbly/config"
	"github.com/valocode/bubbly/ent"
	"github.com/valocode/bubbly/ent/organization"
	"github.com/valocode/bubbly/ent/project"
	"github.com/valocode/bubbly/store/api"
)

type ProjectQuery struct {
	Where *ent.ProjectWhereInput
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
			ProjectModelRead: ent.NewProjectModelRead().FromEnt(proj),
		})
	}
	return projects, nil
}

func (h *Handler) DefaultProjectID() (int, error) {
	defaultProject := config.DefaultProject
	return h.GetProjectIDByName(&defaultProject)
}

func (h *Handler) GetProjectIDByName(name *string) (int, error) {
	if name == nil {
		return -1, fmt.Errorf("project cannot be nil")
	}
	projectID, err := h.client.Project.Query().Where(
		project.HasOwnerWith(organization.ID(h.orgID)),
		project.Name(*name),
	).OnlyID(h.ctx)
	if err != nil {
		return -1, fmt.Errorf("could not get project %s: %w", *name, err)
	}
	return projectID, nil
}

func (h *Handler) CreateProject(req *api.ProjectCreateRequest) (*ent.ProjectModelRead, error) {
	if err := h.validator.Struct(req); err != nil {
		return nil, HandleValidatorError(err, "project create")
	}
	dbProject, err := h.createProject(h.client, req)
	if err != nil {
		return nil, HandleEntError(err, "project create")
	}
	return ent.NewProjectModelRead().FromEnt(dbProject), nil
}

func (h *Handler) createProject(client *ent.Client, req *api.ProjectCreateRequest) (*ent.Project, error) {
	dbProject, err := client.Project.Create().
		SetModelCreate(req.Project).
		SetOwnerID(h.orgID).
		Save(h.ctx)
	if err != nil {
		if ent.IsConstraintError(err) {
			return nil, NewConflictError(err, "project already exists")
		}
		return nil, HandleEntError(err, "project create")
	}
	return dbProject, nil
}
