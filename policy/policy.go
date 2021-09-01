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
	"github.com/valocode/bubbly/ent/releasepolicyviolation"
)

const (
	// policyQuery = `
	// {
	// 	"deny"   : data.policy.deny,
	// 	"require": data.policy.require,
	// }`
	policyQuery   = `data.policy`
	policyModule  = "policy"
	policyPackage = "data.policy"

	denyResult    = "deny"
	requireResult = "require"
)

type runOptions struct {
	trace    bool
	resolver Resolver
}

type PolicyResult struct {
	Violations []*ent.ReleasePolicyViolationModelCreate
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
	queryResult := rs[0]
	if len(queryResult.Bindings) != 0 {
		return nil, fmt.Errorf("result has variable bindings; check the policy query: %v", queryResult.Bindings)
	}
	if len(queryResult.Expressions) == 0 {
		return nil, fmt.Errorf("result has no expressions; check the policy query")
	}
	expr := queryResult.Expressions[0]
	obj, ok := expr.Value.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("internal error: result is not a map[string]interface{}")
	}
	var violations []*ent.ReleasePolicyViolationModelCreate
	for name, value := range obj {
		switch name {
		case denyResult, requireResult:
			rawViolations, ok := value.([]interface{})
			if !ok {
				return nil, fmt.Errorf("expected expression result to be a list")
			}
			// Iterate over each violation as we need to set the type based on the
			// rule. E.g. For a deny rule, a deny violation.
			for _, v := range rawViolations {
				var violation ent.ReleasePolicyViolationModelCreate
				if err := mapstructure.Decode(v, &violation); err != nil {
					return nil, fmt.Errorf("error decoding policy violations: %w", err)
				}
				if violation.Type == nil {
					violation.SetType(releasepolicyviolation.Type(name))
				}
				if err := validator.New().Struct(violation); err != nil {
					return nil, err
				}
				violations = append(violations, &violation)
			}
		}
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
