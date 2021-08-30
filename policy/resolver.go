package policy

import (
	"context"

	"github.com/open-policy-agent/opa/ast"
	"github.com/open-policy-agent/opa/rego"
	"github.com/open-policy-agent/opa/types"
	"github.com/valocode/bubbly/ent"
	"github.com/valocode/bubbly/ent/codeissue"
	"github.com/valocode/bubbly/ent/codescan"
	"github.com/valocode/bubbly/ent/release"
	"github.com/valocode/bubbly/ent/testcase"
	"github.com/valocode/bubbly/ent/testrun"
)

type Resolver interface {
	Functions() []func(*rego.Rego)
}

var _ Resolver = (*EntResolver)(nil)

type EntResolver struct {
	Ctx       context.Context
	Client    *ent.Client
	ReleaseID int
}

func (r *EntResolver) Functions() []func(*rego.Rego) {
	return []func(*rego.Rego){
		r.CodeIssues(),
		r.TestCases(),
	}
}
func (r *EntResolver) CodeIssues() func(*rego.Rego) {
	return rego.Function1(&rego.Function{
		Name: "code_issues",
		Decl: types.NewFunction(
			nil,
			types.NewArray(nil, types.NewObject(nil, types.NewDynamicProperty(types.A, types.A))),
		),
	}, func(bctx rego.BuiltinContext, op1 *ast.Term) (*ast.Term, error) {
		issues, err := r.Client.CodeIssue.Query().Where(
			codeissue.HasScanWith(codescan.HasReleaseWith(release.ID(r.ReleaseID))),
		).All(r.Ctx)
		if err != nil {
			return nil, err
		}
		v, err := ast.InterfaceToValue(issues)
		if err != nil {
			return nil, err
		}
		return ast.NewTerm(v), nil
	})
}

func (r *EntResolver) TestCases() func(*rego.Rego) {
	return rego.Function1(&rego.Function{
		Name: "test_cases",
		Decl: types.NewFunction(
			nil,
			types.NewArray(nil, types.NewObject(nil, types.NewDynamicProperty(types.A, types.A))),
		),
	}, func(bctx rego.BuiltinContext, op1 *ast.Term) (*ast.Term, error) {
		tests, err := r.Client.TestCase.Query().Where(
			testcase.HasRunWith(testrun.HasReleaseWith(release.ID(r.ReleaseID))),
		).All(r.Ctx)
		if err != nil {
			return nil, err
		}
		v, err := ast.InterfaceToValue(tests)
		if err != nil {
			return nil, err
		}
		return ast.NewTerm(v), nil
	})
}
