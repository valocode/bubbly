package api

import "github.com/valocode/bubbly/ent/model"

type (
	ComponentVulnerabilityRequest struct {
		DataSource *string                   `json:"data_source,omitempty"`
		Components []*ComponentVulnerability `json:"components,omitempty"`
	}

	ComponentVulnerability struct {
		Component
		Vulnerabilities []*Vulnerability `json:"vulnerabilities,omitempty"`
	}

	Vulnerability struct {
		model.VulnerabilityModel
	}

	Component struct {
		model.ComponentModel
	}
)
