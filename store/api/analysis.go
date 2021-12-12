package api

type (
	AnalysisCreateRequest struct {
		Tool            *string                `json:"tool,omitempty" validate:"required"`
		Components      []*ComponentCreate     `json:"components,omitempty" validate:"dive,required"`
		Vulnerabilities []*VulnerabilityCreate `json:"vulnerability,omitempty" validate:"dive,required"`
	}
)
