package api

import (
	"github.com/valocode/bubbly/ent"
)

type (
	EventSaveRequest struct {
		Event     *ent.EventModelCreate `json:"event,omitempty" validate:"required"`
		ReleaseID *int                  `json:"release_id,omitempty"`
		RepoID    *int                  `json:"repo_id,omitempty"`
		ProjectID *int                  `json:"project_id,omitempty"`
	}

	EventGetRequest struct {
		Project        string `json:"project,omitempty" query:"project" mapstructure:"project,omitempty"`
		Repo           string `json:"repo,omitempty" query:"repo" mapstructure:"repo,omitempty"`
		ReleaseName    string `json:"release_name,omitempty" query:"release_name" mapstructure:"release_name,omitempty"`
		ReleaseVersion string `json:"release_version,omitempty" query:"release_version" mapstructure:"release_version,omitempty"`
		Commit         string `json:"commit,omitempty" query:"commit" mapstructure:"commit,omitempty"`
		Last           string `json:"last,omitempty" query:"last" mapstructure:"last,omitempty" validate:"number"`
	}

	EventGetResponse struct {
		Events []*ent.Event `json:"events,omitempty" query:"events"`
	}
)
