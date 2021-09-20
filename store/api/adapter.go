package api

import (
	"github.com/valocode/bubbly/ent"
)

type (
	AdapterSaveRequest struct {
		Adapter *ent.AdapterModelCreate `json:"adapter,omitempty" validate:"required"`
	}

	AdapterGetRequest struct {
		Name *string `json:"name,omitempty" param:"name" validate:"required"`
		Tag  *string `json:"tag,omitempty" query:"tag"`
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
		ent.CodeScanModelCreate `json:"code_scan" validate:"required" mapstructure:",squash"`
		Issues                  []*CodeScanIssue     `json:"issues,omitempty" validate:"dive,required"`
		Components              []*CodeScanComponent `json:"components,omitempty" validate:"dive,required"`
	}

	CodeScanIssue struct {
		ent.CodeIssueModelCreate `validate:"required" mapstructure:",squash"`
	}

	CodeScanComponent struct {
		ent.ComponentModelCreate `validate:"required" mapstructure:",squash"`
		Vulnerabilities          []*Vulnerability `json:"vulnerabilities" validate:"dive,required" mapstructure:"vulnerabilities"`
		Licenses                 []*License       `json:"licenses" validate:"dive,required" mapstructure:"licenses"`
	}
)

type (
	TestRunRequest struct {
		TestRun *TestRun `json:"test_run,omitempty"`
		Commit  *string  `json:"commit,omitempty"`
	}

	TestRun struct {
		ent.TestRunModelCreate `mapstructure:",squash"`
		TestCases              []*TestRunCase `json:"test_cases,omitempty" validate:"dive,required"`
	}

	TestRunCase struct {
		ent.TestCaseModelCreate `validate:"required" alias:"test_case"`
	}
)
