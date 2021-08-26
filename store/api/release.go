package api

import (
	"github.com/valocode/bubbly/ent"
)

type (
	ReleaseCreateRequest struct {
		Project *ent.ProjectModelCreate   `json:"project,omitempty" validate:"required"`
		Repo    *ent.RepoModelCreate      `json:"repo,omitempty" validate:"required"`
		Commit  *ent.GitCommitModelCreate `json:"commit,omitempty" validate:"required"`
		Release *ent.ReleaseModelCreate   `json:"release,omitempty" validate:"required"`
	}

	ArtifactLogRequest struct {
		Artifact *ent.ArtifactModelCreate `json:"artifact,omitempty"`
		Commit   *string                  `json:"commit,omitempty"`
	}
)
