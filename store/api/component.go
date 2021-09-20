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
		Patch                        *VulnerabilityPatch `json:"patch,omitempty" mapstructure:"patch"`
	}

	VulnerabilityPatch struct {
		Note *string `json:"note,omitempty" validate:"required" mapstructure:"note"`
	}

	License struct {
		ent.LicenseModelCreate `validate:"required" mapstructure:",squash"`
	}

	ComponentRead struct {
		ent.ComponentModelRead `validate:"required"`
	}
)
