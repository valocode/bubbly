package policy

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/mitchellh/mapstructure"
	"github.com/open-policy-agent/opa/rego"
	"github.com/valocode/bubbly/ent"
)

const (
	violationQuery = "data.bubbly.violation"
)

// Violation is just an alias for a very long struct name
type Violation *ent.ReleasePolicyViolationModelCreate

func EvaluatePolicy(module string, resolver Resolver) ([]Violation, error) {
	ctx := context.Background()
	// Create a simple query over the input.
	rego := rego.New(
		rego.Query(violationQuery),
		rego.Module("module", module),
		resolver.CodeIssues(),
		resolver.TestCases(),
	)

	// Run evaluation
	rs, err := rego.Eval(ctx)
	if err != nil {
		return nil, fmt.Errorf("error evaluating policy: %w", err)
	}
	// ResultSet should never be empty, but just to double check
	if len(rs) == 0 {
		return nil, errors.New("policy result set is empty; check the policy query")
	}
	vResult := rs[0]
	if len(vResult.Bindings) != 0 {
		return nil, fmt.Errorf("violation result has variable bindings; check the policy query: %v", vResult.Bindings)
	}
	if len(vResult.Expressions) == 0 {
		return nil, fmt.Errorf("violation result has no expressions; check the policy query")
	}
	expr := vResult.Expressions[0]
	rawViolations, ok := expr.Value.([]interface{})
	if !ok {
		return nil, errors.New("violation result should be a list of map values, e.g. [{msg: \"policy failed\", severity: \"error\"}, ...]")
	}
	var violations []Violation
	if err := mapstructure.Decode(rawViolations, &violations); err != nil {
		return nil, fmt.Errorf("error decoding policy violations: %w", err)
	}
	if err := validator.New().Var(violations, "required,dive"); err != nil {
		return nil, err
	}

	return violations, nil
}
