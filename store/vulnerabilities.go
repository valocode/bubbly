package store

import (
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqljson"
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
		Where(func(s *sql.Selector) {
			s.Where(
				sqljson.ValueContains(vulnerability.FieldLabels, "windows", sqljson.Path("bubbly/os")),
			)
		}).
		All(h.ctx)
	if err != nil {
		return nil, HandleEntError(err, "getting vulnerabilities")
	}
	var vulnerabilities = make([]*api.Vulnerability, 0, len(dbVulnerabilities))
	for _, vuln := range dbVulnerabilities {
		vulnerabilities = append(vulnerabilities, &api.Vulnerability{
			VulnerabilityModelRead: *ent.NewVulnerabilityModelRead().FromEnt(vuln),
		})
	}
	return vulnerabilities, nil
}
