package policy

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/mitchellh/mapstructure"
	"github.com/open-policy-agent/opa/rego"
	"github.com/valocode/bubbly/ent"
)

const (
	policyQuery   = "data.policy.violation"
	policyModule  = "policy"
	policyPackage = "data.policy"
)

type runOptions struct {
	trace    bool
	resolver Resolver
}

// Violation is just an alias for a very long struct name
type Violation *ent.ReleasePolicyViolationModelCreate

type PolicyResult struct {
	Violations []*Violation
	Traces     []string
}

func EvaluatePolicy(module string, opts ...func(*runOptions)) (*PolicyResult, error) {
	ctx := context.Background()
	regoInstance, err := newRego(module, opts...)
	if err != nil {
		return nil, err
	}

	// Run evaluation
	rs, err := regoInstance.Eval(ctx)
	if err != nil {
		return nil, fmt.Errorf("error evaluating policy: %w", err)
	}

	traceBuf := new(bytes.Buffer)
	rego.PrintTrace(traceBuf, regoInstance)
	var traces []string
	for _, line := range strings.Split(traceBuf.String(), "\n") {
		if len(line) > 0 {
			traces = append(traces, line)
		}
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
	var violations []*Violation
	if err := mapstructure.Decode(rawViolations, &violations); err != nil {
		return nil, fmt.Errorf("error decoding policy violations: %w", err)
	}
	if err := validator.New().Var(violations, "required,dive"); err != nil {
		return nil, err
	}

	return &PolicyResult{
		Violations: violations,
		Traces:     traces,
	}, nil
}

func Validate(module string, opts ...func(r *runOptions)) error {
	ctx := context.Background()
	regoInstance, err := newRego(module, opts...)
	if err != nil {
		return err
	}
	query, pErr := regoInstance.PrepareForEval(ctx)
	if pErr != nil {
		return pErr
	}
	// Check that the module and such is correct
	mod := query.Modules()[policyModule]
	pkg := mod.Package.Path.String()
	if pkg != policyPackage {
		return fmt.Errorf("package should be \"%s\", received \"%s\"", policyPackage, pkg)
	}
	return pErr
}

// WithTracing enables tracing for the run
func WithTracing(trace bool) func(r *runOptions) {
	return func(r *runOptions) {
		r.trace = trace
	}
}

// WithResolver adds a resolver to the policy
func WithResolver(resolver Resolver) func(r *runOptions) {
	return func(r *runOptions) {
		r.resolver = resolver
	}
}

func newRego(module string, opts ...func(*runOptions)) (*rego.Rego, error) {
	r := runOptions{}
	for _, opt := range opts {
		opt(&r)
	}
	regoOptions := []func(*rego.Rego){
		rego.Query(policyQuery),
		rego.Module(policyModule, module),
		rego.Trace(r.trace),
	}
	regoOptions = append(regoOptions, r.resolver.Functions()...)

	return rego.New(regoOptions...), nil
}
