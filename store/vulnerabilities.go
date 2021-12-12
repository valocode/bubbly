package store

import (
	"github.com/valocode/bubbly/ent"
	"github.com/valocode/bubbly/ent/vulnerability"
	"github.com/valocode/bubbly/store/api"
)

type VulnerabilityQuery struct {
	Where *ent.VulnerabilityWhereInput
}

func (h *Handler) GetVulnerabilities(query *VulnerabilityQuery) ([]*api.Vulnerability, error) {
	dbVulnerabilities, err := h.client.Vulnerability.Query().
		WhereInput(query.Where).
		// Where(func(s *sql.Selector) {
		// 	s.Where(
		// 		sqljson.ValueContains(vulnerability.FieldLabels, "windows", sqljson.Path("bubbly/os")),
		// 	)
		// }).
		All(h.ctx)
	if err != nil {
		return nil, HandleEntError(err, "getting vulnerabilities")
	}
	var vulnerabilities = make([]*api.Vulnerability, 0, len(dbVulnerabilities))
	for _, vuln := range dbVulnerabilities {
		vulnerabilities = append(vulnerabilities, &api.Vulnerability{
			VulnerabilityModelRead: ent.NewVulnerabilityModelRead().FromEnt(vuln),
		})
	}
	return vulnerabilities, nil
}

func (h *Handler) CreateVulnerability(req *api.VulnerabilityCreateRequest) (*api.VulnerabilityCreateResponse, error) {
	if err := h.validator.Struct(req); err != nil {
		return nil, HandleValidatorError(err, "vulnerability create")
	}
	dbVuln, err := h.createVulnerability(h.client, req)
	if err != nil {
		return nil, HandleEntError(err, "vulnerability create")
	}

	return &api.VulnerabilityCreateResponse{
		Vulnerability: ent.NewVulnerabilityModelRead().FromEnt(dbVuln),
	}, nil
}

func (h *Handler) createVulnerability(client *ent.Client, req *api.VulnerabilityCreateRequest) (*ent.Vulnerability, error) {
	dbVuln, err := client.Vulnerability.Query().Where(
		vulnerability.Vid(*req.Vulnerability.Vid),
	).Only(h.ctx)
	if err != nil {
		if !ent.IsNotFound(err) {
			return nil, HandleEntError(err, "vulnerability create")
		}
		dbVuln, err := client.Vulnerability.Create().
			SetModelCreate(req.Vulnerability.VulnerabilityModelCreate).
			Save(h.ctx)
		if err != nil {
			return nil, HandleEntError(err, "vulnerability create")
		}
		return dbVuln, nil
	}
	dbVuln, err = client.Vulnerability.UpdateOne(dbVuln).
		SetModelCreate(req.Vulnerability.VulnerabilityModelCreate).
		Save(h.ctx)
	if err != nil {
		return nil, HandleEntError(err, "vulnerability update")
	}
	return dbVuln, nil
}
