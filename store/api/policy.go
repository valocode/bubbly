package api

import "github.com/valocode/bubbly/ent"

type (
	ReleasePolicySaveRequest struct {
		Policy *ent.ReleasePolicyModelCreate `json:"policy,omitempty" validate:"required"`
	}

	ReleasePolicyGetRequest struct {
		Name *string `validate:"required"`
	}

	ReleasePolicyGetResponse struct {
		ent.ReleasePolicyModelRead `validate:"required"`
	}
)
