package api

import (
	"github.com/valocode/bubbly/ent"
)

type (
	AdapterSaveRequest struct {
		Adapter *ent.AdapterModelCreate `json:"adapter,omitempty" validate:"required"`
	}

	AdapterGetRequest struct {
		Name *string `validate:"required"`
		Tag  *string
	}

	AdapterGetResponse struct {
		ent.AdapterModelRead `validate:"required"`
	}
)

type (
	CodeScanRequest struct {
		CodeScan *CodeScan `json:"code_scan,omitempty" validate:"required"`
		Commit   *string   `json:"commit,omitempty" validate:"required"`
	}

	CodeScan struct {
		*ent.CodeScanModelCreate `json:"code_scan" validate:"required"`
		Issues                   []*CodeScanIssue     `json:"issues,omitempty"`
		Components               []*CodeScanComponent `json:"components,omitempty"`
	}

	CodeScanIssue struct {
		ent.CodeIssueModelCreate `validate:"required" mapstructure:",squash"`
	}

	CodeScanComponent struct {
		ent.ComponentModelCreate `validate:"required"`
		Vulnerabilities          []*Vulnerability `json:"vulnerabilities"`
	}
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
