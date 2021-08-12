package api

import "github.com/valocode/bubbly/ent/model"

type (
	ReleaseCreateRequest struct {
		Repo    *model.RepoModel
		Commit  *model.GitCommitModel
		Release *model.ReleaseModel
	}

	ArtifactLogRequest struct {
		Artifact *model.ArtifactModel `json:"artifact,omitempty"`
		Commit   *string              `json:"commit,omitempty"`
	}
)
