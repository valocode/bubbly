package test

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/valocode/bubbly/ent"
	"github.com/valocode/bubbly/store/api"
)

var policyFiles = []string{
	"policy/testdata/code_issue_high_severity.rego",
	"policy/testdata/test_case_fail.rego",
	"policy/testdata/release_checklist.rego",
}

func ParsePolicies(basedir string) ([]*api.ReleasePolicySaveRequest, error) {
	var reqs []*api.ReleasePolicySaveRequest
	for _, f := range policyFiles {
		b, err := os.ReadFile(filepath.Join(basedir, f))
		if err != nil {
			return nil, err
		}

		fname := filepath.Base(f)
		fname = strings.TrimSuffix(fname, filepath.Ext(fname))

		reqs = append(reqs, &api.ReleasePolicySaveRequest{
			Policy: ent.NewReleasePolicyModelCreate().
				SetName(fname).SetModule(string(b)),
		})
	}

	return reqs, nil
}
