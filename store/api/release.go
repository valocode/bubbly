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

	ReleaseGetRequest struct {
		Commit   *string `json:"commit,omitempty" query:"commit" validate:"required"`
		Repo     *string `json:"repo,omitempty" query:"repo"`
		Policies bool    `json:"policies,string" query:"policies"`
		Log      bool    `json:"log,string" query:"log"`
	}

	ReleaseGetResponse struct {
		Releases []*ReleaseRead `json:"releases,omitempty" validate:"dive,required"`
	}

	ReleaseRead struct {
		Project    *ent.ProjectModelRead                  `json:"project,omitempty" validate:"required"`
		Repo       *ent.RepoModelRead                     `json:"repo,omitempty" validate:"required"`
		Commit     *ent.GitCommitModelRead                `json:"commit,omitempty" validate:"required"`
		Release    *ent.ReleaseModelRead                  `json:"release,omitempty" validate:"required"`
		Policies   []*ent.ReleasePolicyModelRead          `json:"policies,omitempty" validate:"dive,required"`
		Violations []*ent.ReleasePolicyViolationModelRead `json:"violations,omitempty" validate:"dive,required"`
		Entries    []*ent.ReleaseEntryModelRead           `json:"entries,omitempty" validate:"dive,required"`
	}
	ArtifactLogRequest struct {
		Artifact *ent.ArtifactModelCreate `json:"artifact,omitempty"`
		Commit   *string                  `json:"commit,omitempty"`
	}
)
