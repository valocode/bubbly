package adapter

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/mitchellh/mapstructure"
	"github.com/open-policy-agent/opa/ast"
	"github.com/open-policy-agent/opa/rego"
	"github.com/open-policy-agent/opa/types"
	"github.com/valocode/bubbly/ent"
	"github.com/valocode/bubbly/store/api"
)

const (
	DefaultTag = "default"

	adapterQuery   = "data.adapter"
	adapterModule  = "adapter"
	adapterPackage = "data.adapter"

	codeScanResult  = "code_scan"
	codeIssueResult = "code_issue"
)

type (
	AdapterResult struct {
		CodeScan *api.CodeScan
		TestRun  *api.TestRun
	}
)

func RunFromFile(path string) (*AdapterResult, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return Run(string(b))
}

func Run(module string) (*AdapterResult, error) {
	ctx := context.Background()
	r := rego.New(
		rego.Query(adapterQuery),
		rego.Module(adapterModule, module),
		rego.Function1(&rego.Function{
			Name: "readfile",
			Decl: types.NewFunction(
				types.Args(types.S),
				types.S,
			),
		}, func(bctx rego.BuiltinContext, op1 *ast.Term) (*ast.Term, error) {
			path := op1.Value.(ast.String)
			b, err := os.ReadFile(string(path))
			if err != nil {
				return nil, err
			}
			return ast.StringTerm(string(b)), nil
		}),
		rego.Function1(&rego.Function{
			Name: "json",
			Decl: types.NewFunction(
				types.Args(types.S),
				types.NewObject(nil, types.NewDynamicProperty(types.A, types.A)),
			),
		}, func(bctx rego.BuiltinContext, op1 *ast.Term) (*ast.Term, error) {
			path := op1.Value.(ast.String)
			file, err := os.Open(string(path))
			if err != nil {
				return nil, err
			}
			var result map[string]interface{}
			if err := json.NewDecoder(file).Decode(&result); err != nil {
				return nil, err
			}
			v, err := ast.InterfaceToValue(result)
			if err != nil {
				return nil, err
			}
			return ast.NewTerm(v), nil
		}),
	)
	query, err := r.PrepareForEval(ctx)
	if err != nil {
		return nil, fmt.Errorf("error preparing adapter for evaluation: %w", err)
	}

	// Check that the module and such is correct
	mod := query.Modules()[adapterModule]
	pkg := mod.Package.Path.String()
	if pkg != adapterPackage {
		return nil, fmt.Errorf("package should be \"%s\", received \"%s\"", adapterPackage, pkg)
	}

	// Run evaluation
	// TODO: add some inputs
	rs, err := query.Eval(ctx)
	// rs, err := r.Eval(ctx)
	if err != nil {
		return nil, fmt.Errorf("error evaluating adapter: %w", err)
	}
	// ResultSet should never be empty, but just to double check
	if len(rs) == 0 {
		return nil, errors.New("adapter result set is empty; check the adapter query")
	}
	queryResult := rs[0]
	if len(queryResult.Bindings) != 0 {
		return nil, fmt.Errorf("result has variable bindings; check the adapter query: %v", queryResult.Bindings)
	}
	if len(queryResult.Expressions) == 0 {
		return nil, fmt.Errorf("result has no expressions; check the policy query")
	}
	expr := queryResult.Expressions[0]
	obj, ok := expr.Value.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("internal error: result is not a map[string]interface{}")
	}
	var (
		codeScans  []*ent.CodeScanModelCreate
		codeIssues []*api.CodeScanIssue
	)
	validate := validator.New()
	for name, value := range obj {
		switch name {
		case codeScanResult:
			if err := mapstructure.Decode(value, &codeScans); err != nil {
				return nil, fmt.Errorf("error decoding code scan: %w", err)
			}
			if len(codeScans) > 1 {
				return nil, fmt.Errorf("multiple code scans detected: only one is allowed per adapter")
			}
			if err := validate.Var(codeScans, "required,dive"); err != nil {
				return nil, err
			}
		case codeIssueResult:
			if err := mapstructure.Decode(value, &codeIssues); err != nil {
				return nil, fmt.Errorf("error decoding code issues: %w", err)
			}
			if err := validate.Var(codeIssues, "required,dive"); err != nil {
				return nil, err
			}
		}
	}

	var result AdapterResult
	// Check that the data seems appropriate
	if codeScans != nil {
		result.CodeScan = &api.CodeScan{
			CodeScanModelCreate: codeScans[0],
			Issues:              codeIssues,
		}
	} else {
		// Check if there was no code_scan that there were no results to associate
		// the missing code_scan
		if codeIssues != nil {
			return nil, fmt.Errorf("cannot provide code_issue without code_scan: please provide a code_scan")
		}
	}

	return &result, nil
}

func ParseAdpaterID(id string) (string, string, error) {
	splitID := strings.Split(id, ":")
	switch len(splitID) {
	case 1:
		return id, "", nil
	case 2:
		return splitID[0], splitID[1], nil
	}
	return "", "", fmt.Errorf("adapter must be in the form \"name:tag\"")
}
