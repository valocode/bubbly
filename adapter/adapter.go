package adapter

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/mitchellh/mapstructure"
	"github.com/open-policy-agent/conftest/parser"
	"github.com/open-policy-agent/opa/rego"
	"github.com/valocode/bubbly/client"
	"github.com/valocode/bubbly/ent"
	"github.com/valocode/bubbly/env"
	"github.com/valocode/bubbly/store/api"
)

const (
	DefaultTag = "default"

	adapterQuery   = "data.adapter"
	adapterModule  = "adapter"
	adapterPackage = "data.adapter"

	codeScanResult  = "code_scan"
	codeIssueResult = "code_issue"
	componentResult = "component"

	testRunResult  = "test_run"
	testCaseResult = "test_case"
)

type (
	AdapterResult struct {
		CodeScan *api.CodeScan
		TestRun  *api.TestRun
		Traces   []string
	}

	runOptions struct {
		trace      bool
		inputFiles []string
	}
)

// RunFromID runs an adapter by the given id.
// This will get the relevant adapter (from filesystem or remotely) and perform
// the necessary rego queries and publish that data to the release.
func RunFromID(bCtx *env.BubblyContext, id string, opts ...func(r *runOptions)) (*AdapterResult, error) {
	//
	// Get the adapter module to run
	//
	module, err := adapterFromID(bCtx, id)
	if err != nil {
		return nil, err
	}

	return Run(module, opts...)
}

func RunFromFile(path string, opts ...func(r *runOptions)) (*AdapterResult, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return Run(string(b), opts...)
}

func Run(module string, opts ...func(r *runOptions)) (*AdapterResult, error) {
	ctx := context.Background()
	regoInstance, err := newRego(module, opts...)
	if err != nil {
		return nil, err
	}
	rs, err := regoInstance.Eval(ctx)
	if err != nil {
		return nil, err
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
		return nil, errors.New("adapter result set is empty; check the adapter query")
	}
	queryResult := rs[0]
	if len(queryResult.Bindings) != 0 {
		return nil, fmt.Errorf("result has variable bindings; check the adapter query: %v", queryResult.Bindings)
	}
	if len(queryResult.Expressions) == 0 {
		return nil, fmt.Errorf("result has no expressions; check the adapter query")
	}

	expr := queryResult.Expressions[0]
	obj, ok := expr.Value.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("internal error: result is not a map[string]interface{}")
	}
	var (
		codeScans  []*api.CodeScan
		codeIssues []*api.CodeScanIssue
		components []*api.CodeScanComponent

		testRuns  []*ent.TestRunModelCreate
		testCases []*api.TestRun
	)

	for name, value := range obj {
		switch name {
		case codeScanResult:
			if err := mapstructure.Decode(value, &codeScans); err != nil {
				return nil, fmt.Errorf("error decoding code scan: %w", err)
			}
			if len(codeScans) > 1 {
				return nil, fmt.Errorf("multiple code scans detected: only one is allowed per adapter")
			}
		case codeIssueResult:
			if err := mapstructure.Decode(value, &codeIssues); err != nil {
				return nil, fmt.Errorf("error decoding code issues: %w", err)
			}
		case componentResult:
			if err := mapstructure.Decode(value, &components); err != nil {
				return nil, fmt.Errorf("error decoding components: %w", err)
			}
		case testRunResult:
			if err := mapstructure.Decode(value, &testRuns); err != nil {
				return nil, fmt.Errorf("error decoding test run: %w", err)
			}
			if len(testRuns) > 1 {
				return nil, fmt.Errorf("multiple test runs detected: only one is allowed per adapter")
			}
		case testCaseResult:
			if err := mapstructure.Decode(value, &testCases); err != nil {
				return nil, fmt.Errorf("error decoding test cases: %w", err)
			}
		}
	}

	result := AdapterResult{
		Traces: traces,
	}
	validate := validator.New()
	// Check that the data seems appropriate
	if codeScans != nil {
		result.CodeScan = codeScans[0]
		result.CodeScan.Issues = codeIssues
		if err := validate.Struct(result.CodeScan); err != nil {
			return nil, err
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
	mod := query.Modules()[adapterModule]
	pkg := mod.Package.Path.String()
	if pkg != adapterPackage {
		return fmt.Errorf("package should be \"%s\", received \"%s\"", adapterPackage, pkg)
	}
	return nil
}

// WithTracing enables tracing for the run
func WithTracing(trace bool) func(r *runOptions) {
	return func(r *runOptions) {
		r.trace = trace
	}
}

// WithInputFiles adds the given file list (comma-separated) to the input files
func WithInputFiles(files string) func(r *runOptions) {
	return func(r *runOptions) {
		r.inputFiles = append(r.inputFiles, strings.Split(files, ",")...)
	}
}

// WithInputFileSlice adds the given file list to the input files
func WithInputFileSlice(files []string) func(r *runOptions) {
	return func(r *runOptions) {
		r.inputFiles = append(r.inputFiles, files...)
	}
}

func newRego(module string, opts ...func(*runOptions)) (*rego.Rego, error) {
	r := runOptions{}
	for _, opt := range opts {
		opt(&r)
	}
	//
	// Get the inputs by parsing the inputs
	//
	configurations, err := parser.ParseConfigurations(r.inputFiles)
	if err != nil {
		return nil, fmt.Errorf("error parsing input files")
	}

	regoOptions := []func(*rego.Rego){
		rego.Query(adapterQuery),
		rego.Module(adapterModule, module),
		rego.Trace(r.trace),
		rego.Input(configurations),
	}
	return rego.New(regoOptions...), nil
}

func ParseAdpaterID(id string) (string, string, error) {
	splitID := strings.Split(id, ":")
	var (
		adapterName = splitID[0]
		adapterTag  = DefaultTag
	)
	if len(splitID) == 2 {
		adapterTag = splitID[1]
	}
	if len(splitID) > 2 {
		return "", "", fmt.Errorf("invalid adapter format, should be \"name:tag\": %s", id)
	}
	return adapterName, adapterTag, nil
}

// adapterFromID takes the (user provided) id of an adapter which can either be
// a path to a local file, the name of an adapter (which could lead to a local
// an adapter file in the bubbly directory), or a name:tag for which the adapter
// is fetched remotely.
func adapterFromID(bCtx *env.BubblyContext, id string) (string, error) {
	switch strings.Count(id, ":") {
	case 0:
		// It could be a local adapter, or remote with a default tag.
		// Create a list of possible paths and check them
		possiblePaths := []string{
			id, filepath.Join(bCtx.ReleaseConfig.BubblyDir, "adapter", id+".rego"),
		}
		for _, path := range possiblePaths {
			if _, err := os.Stat(path); err == nil {
				// We have a local file so read and return it
				b, err := os.ReadFile(path)
				if err != nil {
					return "", fmt.Errorf("error reading adapter file: %w", err)
				}
				return string(b), nil
			} else if os.IsNotExist(err) {
				// If it doesn't exist, don't worry, just carry on
				continue
			} else {
				return "", fmt.Errorf("error checking for adapter existence: %s: %w", path, err)
			}
		}
		// If not found, and no errors, continue to fetch remotely
	case 1:
		// It is an adapter with a name:tag so fetch remotely
	default:
		return "", fmt.Errorf("invalid format of adapter: %s", id)
	}

	splitID := strings.Split(id, ":")
	var (
		adapterName = splitID[0]
		adapterTag  = DefaultTag
	)
	if len(splitID) == 2 {
		adapterTag = splitID[1]
	}
	resp, err := client.GetAdapter(bCtx, &api.AdapterGetRequest{
		Name: &adapterName,
		Tag:  &adapterTag,
	})
	if err != nil {
		return "", fmt.Errorf("error getting adapter remotely: %s: %w", id, err)
	}
	return *resp.Module, nil
}
