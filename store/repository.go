package store

import (
	"github.com/valocode/bubbly/ent"
	"github.com/valocode/bubbly/ent/project"
	"github.com/valocode/bubbly/ent/repository"
	"github.com/valocode/bubbly/store/api"
)

type RepositoryQuery struct {
	Where *ent.RepositoryWhereInput
}

func (h *Handler) GetRepositories(query *RepositoryQuery) (*api.RepositoryGetResponse, error) {
	dbRepositorys, err := h.client.Repository.Query().
		WhereInput(query.Where).
		All(h.ctx)
	if err != nil {
		return nil, HandleEntError(err, "repository get")
	}
	var repositories = make([]*api.Repository, 0, len(dbRepositorys))
	for _, repo := range dbRepositorys {
		repositories = append(repositories, &api.Repository{
			RepositoryModelRead: ent.NewRepositoryModelRead().FromEnt(repo),
		})
	}
	return &api.RepositoryGetResponse{
		Repositories: repositories,
	}, nil
}

func (h *Handler) CreateRepository(req *api.RepositoryCreateRequest) (*api.RepositoryCreateResponse, error) {
	if err := h.validator.Struct(req); err != nil {
		return nil, HandleValidatorError(err, "repository create")
	}
	dbRepository, err := h.createRepository(h.client, req)
	if err != nil {
		return nil, HandleEntError(err, "repository create")
	}
	return &api.RepositoryCreateResponse{
		Repository: ent.NewRepositoryModelRead().FromEnt(dbRepository),
	}, nil
}

func (h *Handler) createRepository(client *ent.Client, req *api.RepositoryCreateRequest) (*ent.Repository, error) {
	dbRepository, err := client.Repository.Query().
		Where(
			repository.HasProjectWith(project.Name(*req.Project)),
			repository.Name(*req.Repository.Name),
		).Only(h.ctx)
	if err != nil {
		if !ent.IsNotFound(err) {
			return nil, HandleEntError(err, "repository create")
		}
		projectID, err := h.GetProjectIDByName(req.Project)
		if err != nil {
			return nil, HandleEntError(err, "repository create")
		}
		dbRepository, err := client.Repository.Create().
			SetModelCreate(req.Repository).
			SetProjectID(projectID).
			SetOwnerID(h.orgID).
			Save(h.ctx)
		if err != nil {
			return nil, HandleEntError(err, "repository create")
		}
		return dbRepository, nil
	}
	dbRepository, err = client.Repository.UpdateOne(dbRepository).
		SetModelCreate(req.Repository).
		Save(h.ctx)
	if err != nil {
		return nil, HandleEntError(err, "repository create")
	}
	return dbRepository, nil
}
