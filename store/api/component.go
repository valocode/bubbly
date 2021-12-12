package api

import (
	"github.com/valocode/bubbly/ent"
)

type (
	ComponentGetRequest struct {
	}

	ComponentGetResponse struct {
		Components []*Component `json:"components"`
	}

	ComponentCreateRequest struct {
		Component *ComponentCreate `json:"component,omitempty" validate:"required"`
	}

	ComponentCreate struct {
		*ent.ComponentModelCreate
		Vulnerabilities []*ent.VulnerabilityModelCreate `json:"vulnerabilities,omitempty" validate:"dive,required"`
	}

	License struct {
		ent.LicenseModelCreate `validate:"required" mapstructure:",squash"`
	}

	Component struct {
		ent.ComponentModelRead `validate:"required"`
	}
)
