package api

import "github.com/valocode/bubbly/ent"

type (
	ProjectGetRequest struct {
		Name string `query:"name"`
	}
	ProjectGetResponse struct {
		Projects []*Project `json:"projects"`
	}
	Project struct {
		*ent.ProjectModelRead
	}

	ProjectCreateRequest struct {
		Project *ent.ProjectModelCreate `json:"project,omitempty" validate:"required"`
	}
	ProjectCreateResponse struct {
		Project *ent.ProjectModelRead `json:"project"`
	}
)
