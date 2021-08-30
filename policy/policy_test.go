package policy

import (
	"fmt"
	"os"
	"testing"

	"github.com/open-policy-agent/opa/ast"
	"github.com/open-policy-agent/opa/rego"
	"github.com/open-policy-agent/opa/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var _ Resolver = (*fakeResolver)(nil)

type fakeResolver struct {
	data map[string][]map[string]interface{}
}

func (r *fakeResolver) Functions() []func(*rego.Rego) {
	var funcs = make([]func(*rego.Rego), 0, len(r.data))
	for name, value := range r.data {
		funcs = append(funcs,
			rego.Function1(&rego.Function{
				Name: name,
				Decl: types.NewFunction(
					nil,
					types.NewArray(nil, types.NewObject(nil, types.NewDynamicProperty(types.A, types.A))),
				),
			}, func(bctx rego.BuiltinContext, op1 *ast.Term) (*ast.Term, error) {
				v, err := ast.InterfaceToValue(value)
				if err != nil {
					return nil, err
				}
				return ast.NewTerm(v), nil
			}),
		)
	}
	return funcs
}

func TestPolicy(t *testing.T) {
	type test struct {
		name   string
		input  map[string][]map[string]interface{}
		policy string
		want   int
	}
	tests := []test{
		{
			name:   "test_case_policy_failing",
			input:  map[string][]map[string]interface{}{"test_cases": {{"result": false}}},
			policy: "./testdata/test_case_fail.rego",
			want:   1,
		},
		{
			name:   "test_case_policy_passing",
			input:  map[string][]map[string]interface{}{"test_cases": {{"result": true}}},
			policy: "./testdata/test_case_fail.rego",
			want:   0,
		},
		{
			name:   "code_issues_high_severity_failing",
			input:  map[string][]map[string]interface{}{"code_issues": {{"severity": "high"}}},
			policy: "./testdata/code_issue_high_severity.rego",
			want:   1,
		},
		{
			name:   "code_issues_high_severity_passing",
			input:  map[string][]map[string]interface{}{"code_issues": {{"severity": "low"}}},
			policy: "./testdata/code_issue_high_severity.rego",
			want:   0,
		},
		{
			name:   "release_checklist_failing",
			input:  map[string][]map[string]interface{}{"release_entries": {{"type": "code_scan"}}},
			policy: "./testdata/release_checklist.rego",
			want:   1,
		},
		// {
		// 	name:   "code_issues_high_severity_passing",
		// 	input:  map[string][]map[string]interface{}{"code_issues": {{"severity": "low"}}},
		// 	policy: "./testdata/code_issue_high_severity.rego",
		// 	want:   0,
		// },
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			b, err := os.ReadFile(tc.policy)
			require.NoError(t, err)
			result, err := EvaluatePolicy(string(b),
				WithResolver(&fakeResolver{data: tc.input}),
				WithTracing(true),
			)
			require.NoError(t, err)
			assert.Len(t, result.Violations, tc.want)

			for _, t := range result.Traces {
				fmt.Println(t)
			}
		})
	}
}
