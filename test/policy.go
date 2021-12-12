package test

import (
	"fmt"

	"github.com/valocode/bubbly/policy"
	"github.com/valocode/bubbly/store/api"
)

func ParsePolicies() ([]*api.ReleasePolicyCreateRequest, error) {
	var reqs []*api.ReleasePolicyCreateRequest
	policies, err := policy.BuiltinPolicies()
	if err != nil {
		return nil, err
	}
	if len(policies) == 0 {
		return nil, fmt.Errorf("no builtin policies found")
	}
	for _, policy := range policies {
		reqs = append(reqs, &api.ReleasePolicyCreateRequest{
			Policy: &api.ReleasePolicyCreate{
				ReleasePolicyModelCreate: *policy,
			},
		})
	}

	return reqs, nil
}
