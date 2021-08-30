package api

import "github.com/valocode/bubbly/ent"

type (
	ReleasePolicySaveRequest struct {
		Policy  *ent.ReleasePolicyModelCreate `json:"policy,omitempty" validate:"required"`
		Affects *ReleasePolicyAffects         `json:"affects,omitempty"`
	}

	ReleasePolicySetRequest struct {
		Policy  *string               `json:"policy,omitempty" validate:"required"`
		Affects *ReleasePolicyAffects `json:"affects,omitempty"`
	}

	ReleasePolicyGetRequest struct {
		Name *string `validate:"required"`
	}

	ReleasePolicyGetResponse struct {
		ent.ReleasePolicyModelRead `validate:"required"`
	}

	// ReleasePolicyAffects defines the entities that a policy should affect.
	// The Not will remove the entity from the policy relationship (if exists).
	ReleasePolicyAffects struct {
		Projects    []string
		NotProjects []string

		Repos    []string
		NotRepos []string
	}
)
