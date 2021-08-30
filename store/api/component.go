package api

import (
	"github.com/valocode/bubbly/ent"
)

type (
	ComponentVulnerabilityRequest struct {
		DataSource *string                   `json:"data_source,omitempty" validate:"required"`
		Components []*ComponentVulnerability `json:"components,omitempty" validate:"required,dive,required"`
	}

	ComponentVulnerability struct {
		ComponentID     *int             `json:"component_id" validate:"required"`
		Vulnerabilities []*Vulnerability `json:"vulnerabilities,omitempty" validate:"dive,required"`
	}

	Vulnerability struct {
		ent.VulnerabilityModelCreate `validate:"required" mapstructure:",squash"`
	}

	ComponentRead struct {
		ent.ComponentModelRead `validate:"required"`
	}
)
