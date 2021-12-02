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
		Commit   string `json:"commit,omitempty" query:"commit"`
		Repos    string `json:"repos,omitempty" query:"repos"`
		Projects string `json:"projects,omitempty" query:"projects"`
		Policies bool   `json:"policies,string" query:"policies"`
		Log      bool   `json:"log,string" query:"log"`
		HeadOnly bool   `json:"head_only,string" query:"head_only"`

		SortBy string `json:"sort_by,omitempty" query:"sort_by"`
	}

	ReleaseGetResponse struct {
		Releases []*Release `json:"releases" validate:"dive,required"`
	}

	ReleaseGetByIDRequest struct {
		ID int `param:"id"`
	}

	ReleaseGetByIDResponse struct {
		Release *Release `json:"release,omitempty" validate:"required"`
	}
	Release struct {
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
