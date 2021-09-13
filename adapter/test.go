package adapter

import (
	"context"
	"fmt"

	"github.com/open-policy-agent/opa/loader"
	"github.com/open-policy-agent/opa/tester"
)

func RunTests(ctx context.Context, files []string) ([]*tester.Result, error) {
	result, err := loader.NewFileLoader().All(files)
	if err != nil {
		return nil, fmt.Errorf("loading rego files: %w", err)
	}
	runner := tester.NewRunner().
		SetModules(result.ParsedModules()).
		Filter(".*")

	ch, err := runner.RunTests(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("running rego tests: %w", err)
	}
	var rawResults []*tester.Result
	for result := range ch {
		rawResults = append(rawResults, result)
	}
	return rawResults, nil
}
