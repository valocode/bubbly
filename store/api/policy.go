package api

import "github.com/valocode/bubbly/ent"

type (
	ReleasePolicyCreateRequest struct {
		Policy *ReleasePolicyCreate `json:"policy,omitempty" validate:"required"`
	}

	ReleasePolicyCreateResponse struct {
		Policy *ent.ReleasePolicyModelRead `json:"policy"`
	}

	ReleasePolicyUpdateRequest struct {
		ID     *int                 `param:"id" validate:"required"`
		Policy *ReleasePolicyUpdate `json:"policy,omitempty" validate:"required"`
	}

	ReleasePolicyUpdateResponse struct {
		Policy *ent.ReleasePolicyModelRead `json:"policy"`
	}

	ReleasePolicyGetRequest struct {
		Name string `json:"name" query:"name"`

		WithAffects bool `json:"with_affects,string" query:"with_affects"`
	}

	ReleasePolicyGetResponse struct {
		Policies []*ReleasePolicy `json:"policies"`
	}

	ReleasePolicyCreate struct {
		ent.ReleasePolicyModelCreate
		Affects *ReleasePolicyAffectsSet `json:"affects,omitempty"`
	}
	ReleasePolicyUpdate struct {
		ent.ReleasePolicyModelUpdate
		Affects *ReleasePolicyAffectsSet `json:"affects,omitempty"`
	}

	ReleasePolicy struct {
		ent.ReleasePolicyModelRead
		Affects *ReleasePolicyAffects `json:"affects,omitempty"`
	}

	// ReleasePolicyAffects
	ReleasePolicyAffects struct {
		Projects []string `json:"projects,omitempty"`
		Repos    []string `json:"repos,omitempty"`
	}

	// ReleasePolicyAffectsSet defines the entities that a policy should affect.
	// The Not will remove the entity from the policy relationship (if exists).
	ReleasePolicyAffectsSet struct {
		Projects    []string `json:"projects,omitempty"`
		NotProjects []string `json:"not_projects,omitempty"`

		Repos    []string `json:"repos,omitempty"`
		NotRepos []string `json:"not_repos,omitempty"`
	}
)
