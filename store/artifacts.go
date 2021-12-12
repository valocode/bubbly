package store

import (
	"github.com/valocode/bubbly/ent"
	"github.com/valocode/bubbly/store/api"
)

type ArtifactQuery struct {
	Where *ent.ArtifactWhereInput
}

func (h *Handler) GetArtifacts(query *ArtifactQuery) ([]*api.Artifact, error) {
	dbArtifacts, err := h.client.Artifact.Query().
		WhereInput(query.Where).
		All(h.ctx)
	if err != nil {
		return nil, HandleEntError(err, "artifact get")
	}
	var artifacts = make([]*api.Artifact, 0, len(dbArtifacts))
	for _, repo := range dbArtifacts {
		artifacts = append(artifacts, &api.Artifact{
			ArtifactModelRead: ent.NewArtifactModelRead().FromEnt(repo),
		})
	}
	return artifacts, nil
}

func (h *Handler) CreateArtifact(req *api.ArtifactCreateRequest) (*ent.ArtifactModelRead, error) {
	if err := h.validator.Struct(req); err != nil {
		return nil, HandleValidatorError(err, "artifact create")
	}
	dbArtifact, err := h.createArtifact(h.client, req)
	if err != nil {
		return nil, HandleEntError(err, "artifact create")
	}
	return ent.NewArtifactModelRead().FromEnt(dbArtifact), nil
}

func (h *Handler) createArtifact(client *ent.Client, req *api.ArtifactCreateRequest) (*ent.Artifact, error) {
	dbArtifact, err := h.client.Artifact.Create().
		SetModelCreate(req.Artifact).
		// SetRelease()
		// SetProjectID(dbProject.ID).
		// SetOwnerID(h.orgID).
		Save(h.ctx)
	if err != nil {
		return nil, HandleEntError(err, "artifact create")
	}
	return dbArtifact, nil
}
