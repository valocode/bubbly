package api

import (
	"github.com/valocode/bubbly/ent"
)

type (
	CodeScanRequest struct {
		CodeScan *CodeScan `json:"code_scan,omitempty" validate:"required"`
		Commit   *string   `json:"commit,omitempty" validate:"required"`
	}

	CodeScan struct {
		*ent.CodeScanModelCreate `json:"code_scan" validate:"required"`
		Issues                   []*CodeScanIssue     `json:"issues,omitempty" hcl:"issue,block"`
		Components               []*CodeScanComponent `json:"components,omitempty" hcl:"component,block"`
	}

	CodeScanIssue struct {
		ent.CodeIssueModelCreate `validate:"required" hcl:",remain"`
	}

	CodeScanComponent struct {
		ent.ComponentModelCreate `validate:"required" hcl:",remain"`
		Vulnerabilities          []*Vulnerability `json:"vulnerabilities" hcl:"vulnerability,block"`
	}
)
