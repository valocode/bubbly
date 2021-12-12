package api

import "github.com/valocode/bubbly/ent"

type (
	RepositoryGetRequest struct {
		Name string `query:"name"`
	}

	RepositoryGetResponse struct {
		Repositories []*Repository `json:"repositories"`
	}
	Repository struct {
		*ent.RepositoryModelRead
	}

	RepositoryCreateRequest struct {
		Project    *string                    `json:"project,omitempty" validate:"required"`
		Repository *ent.RepositoryModelCreate `json:"repository,omitempty" validate:"required"`
	}
	RepositoryCreateResponse struct {
		Repository *ent.RepositoryModelRead `json:"repository"`
	}
)
