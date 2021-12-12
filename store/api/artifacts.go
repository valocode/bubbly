package api

import "github.com/valocode/bubbly/ent"

type (
	ArtifactGetRequest struct {
		Name string `query:"name"`
	}

	ArtifactGetResponse struct {
		Artifacts []*Artifact `json:"artifacts"`
	}
	Artifact struct {
		*ent.ArtifactModelRead
	}

	ArtifactCreateRequest struct {
		Commit   *string                  `json:"commit,omitempty" validate:"required"`
		Artifact *ent.ArtifactModelCreate `json:"artifact,omitempty" validate:"required"`
	}
	ArtifactCreateResponse struct {
		Artifact *ent.ArtifactModelRead `json:"artifact"`
	}
)
