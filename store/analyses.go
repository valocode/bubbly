package store

import (
	"github.com/valocode/bubbly/ent"
	"github.com/valocode/bubbly/store/api"
)

func (h *Handler) CreateAnalysis(req *api.AnalysisCreateRequest) error {
	if err := h.validator.Struct(req); err != nil {
		return HandleValidatorError(err, "analysis create")
	}
	if err := h.createAnalysis(h.client, req); err != nil {
		return HandleEntError(err, "analysis create")
	}
	return nil
}

func (h *Handler) createAnalysis(client *ent.Client, req *api.AnalysisCreateRequest) error {
	for _, comp := range req.Components {
		_, err := h.createComponent(h.client, &api.ComponentCreateRequest{Component: comp})
		if err != nil {
			return err
		}
	}
	for _, vuln := range req.Vulnerabilities {
		_, err := h.createVulnerability(h.client, &api.VulnerabilityCreateRequest{Vulnerability: vuln})
		if err != nil {
			return err
		}
	}
	return nil
}
