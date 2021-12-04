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

	ComponentSaveRequest struct {
		Component *ComponentCreate `json:"component,omitempty" validate:"required"`
	}

	ComponentCreate struct {
		ent.ComponentModelCreate
	}

	ComponentVulnerabilityRequest struct {
		DataSource *string                   `json:"data_source,omitempty" validate:"required"`
		Components []*ComponentVulnerability `json:"components,omitempty" validate:"required,dive,required"`
	}

	ComponentVulnerability struct {
		ComponentID     *int                   `json:"component_id" validate:"required"`
		Vulnerabilities []*VulnerabilityCreate `json:"vulnerabilities,omitempty" validate:"dive,required"`
	}

	License struct {
		ent.LicenseModelCreate `validate:"required" mapstructure:",squash"`
	}

	Component struct {
		ent.ComponentModelRead `validate:"required"`
	}
)
