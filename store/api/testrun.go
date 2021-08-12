package api

import "github.com/valocode/bubbly/ent/model"

type (
	TestRunRequest struct {
		TestRun *TestRun `json:"test_run,omitempty"`
		Commit  *string  `json:"commit,omitempty"`
	}

	TestRun struct {
		model.TestRunModel
		TestCases []*TestRunCase `json:"test_cases,omitempty"`
	}

	TestRunCase struct {
		model.TestCaseModel
	}
)
