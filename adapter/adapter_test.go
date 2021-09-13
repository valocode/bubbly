package adapter

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/open-policy-agent/conftest/output"
	"github.com/open-policy-agent/opa/loader"
	"github.com/open-policy-agent/opa/tester"
	"github.com/open-policy-agent/opa/topdown"
	"github.com/stretchr/testify/require"
)

func TestRego(t *testing.T) {
	result, err := RunFromFile(
		"./testdata/adapters/gosec.rego",
		WithInputFiles("./testdata/adapters/gosec.json"),
		// WithTracing(true),
	)
	require.NoError(t, err)
	t.Logf("result: %#v", result.CodeScan)
	for _, trace := range result.Traces {
		fmt.Println(trace)
	}
}

func TestUnit(t *testing.T) {
	// result, err := loader.NewFileLoader().All([]string{"./testdata/adapters/gosec.rego", "./testdata/adapters/gosec_test.rego"})
	result, err := loader.NewFileLoader().All([]string{"./testdata/adapters/gosec.rego", "./testdata/adapters/gosec_test.rego"})
	require.NoError(t, err)

	// r, err := newRego("package adapter")
	require.NoError(t, err)
	// pq, err := r.PrepareForEval(context.Background())
	require.NoError(t, err)
	// runner := tester.NewRunner().SetModules(pq.Modules())
	runner := tester.NewRunner().
		SetModules(result.ParsedModules()).
		Filter(".*")

	ch, err := runner.RunTests(context.Background(), nil)
	require.NoError(t, err)

	var results []output.CheckResult
	var rawResults []*tester.Result
	for result := range ch {
		require.NoError(t, err)
		rawResults = append(rawResults, result)
		buf := new(bytes.Buffer)
		topdown.PrettyTrace(buf, result.Trace)
		var traces []string
		for _, line := range strings.Split(buf.String(), "\n") {
			if len(line) > 0 {
				traces = append(traces, line)
			}
		}

		var outputResult output.Result
		if result.Fail || result.Skip {
			outputResult.Message = result.Package + "." + result.Name
		}

		queryResult := output.QueryResult{
			Query:   result.Name,
			Results: []output.Result{outputResult},
			Traces:  traces,
		}

		checkResult := output.CheckResult{
			FileName: result.Location.File,
			Queries:  []output.QueryResult{queryResult},
		}
		if result.Fail {
			checkResult.Failures = []output.Result{outputResult}
		} else if result.Skip {
			checkResult.Skipped = []output.Result{outputResult}
		} else {
			checkResult.Successes++
		}

		results = append(results, checkResult)
		fmt.Printf("Result: %#v\n", checkResult)
	}
}
