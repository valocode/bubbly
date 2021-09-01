package api

import "github.com/valocode/bubbly/ent"

type ProjectCreateRequest struct {
	Project *ent.ProjectModelCreate `json:"project,omitempty" validate:"required"`
}
