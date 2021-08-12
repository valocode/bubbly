package api

import "github.com/valocode/bubbly/ent/model"

type (
	CodeScanRequest struct {
		CodeScan *CodeScan `json:"code_scan,omitempty"`
		Commit   *string   `json:"commit,omitempty" validate:"required"`
	}

	CodeScan struct {
		model.CodeScanModel
		Issues     []*CodeScanIssue     `json:"issues,omitempty" hcl:"issue,block"`
		Components []*CodeScanComponent `json:"components,omitempty" hcl:"component,block"`
	}

	CodeScanIssue struct {
		model.CodeIssueModel `hcl:",remain"`
	}

	CodeScanComponent struct {
		model.ComponentModel
		Vulnerabilities []*Vulnerability `json:"vulnerabilities" hcl:"vulnerability,block"`
	}
)
