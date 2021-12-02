package api

import "github.com/valocode/bubbly/ent"

type ProjectCreateRequest struct {
	Project *ent.ProjectModelCreate `json:"project,omitempty" validate:"required"`
}

type ProjectGetResponse struct {
	Projects []*Project `json:"projects"`
}

type ProjectGetRequest struct {
	Name string `query:"name"`
}

type Project struct {
	ent.ProjectModelRead
}
