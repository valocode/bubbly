package api

import (
	"github.com/valocode/bubbly/ent"
)

type (
	TestRunRequest struct {
		TestRun *TestRun `json:"test_run,omitempty"`
		Commit  *string  `json:"commit,omitempty"`
	}

	TestRun struct {
		ent.TestRunModelCreate
		TestCases []*TestRunCase `json:"test_cases,omitempty" validate:"dive,required"`
	}

	TestRunCase struct {
		ent.TestCaseModelCreate `validate:"required" alias:"test_case"`
	}
)
