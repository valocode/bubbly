package store

import (
	"github.com/valocode/bubbly/ent"
	"github.com/valocode/bubbly/ent/vulnerability"
	"github.com/valocode/bubbly/store/api"
)

func (h *Handler) SaveComponentVulnerabilities(req *api.ComponentVulnerabilityRequest) error {
	if err := h.validator.Struct(req); err != nil {
		return HandleValidatorError(err, "component vulnerabilities")
	}
	txErr := WithTx(h.ctx, h.client, func(tx *ent.Tx) error {
		for _, comp := range req.Components {
			dbComp, err := h.GetComponentOrError(tx.Client(), *comp.ComponentID)
			if err != nil {
				return err
			}
			var dbVulns = make([]*ent.Vulnerability, 0, len(comp.Vulnerabilities))
			for _, vuln := range comp.Vulnerabilities {
				dbVuln, err := h.GetVulnerabilityOrCreate(tx.Client(), vuln)
				if err != nil {
					return err
				}
				dbVulns = append(dbVulns, dbVuln)
			}
			_, err = tx.Component.UpdateOne(dbComp).
				AddVulnerabilities(dbVulns...).
				Save(h.ctx)
			if err != nil {
				return HandleEntError(err, "component vulnerability")
			}
		}
		return nil
	})
	if txErr != nil {
		return txErr
	}
	return nil
}

func (h *Handler) GetComponentOrError(client *ent.Client, id int) (*ent.Component, error) {
	dbComp, err := client.Component.Get(h.ctx, id)
	if err != nil {
		if !ent.IsNotFound(err) {
			return nil, NewNotFoundError(nil, "component with ID %d", id)
		}
		return nil, HandleEntError(err, "component")
	}
	return dbComp, nil
}

func (h *Handler) GetVulnerabilityOrCreate(client *ent.Client, vuln *api.Vulnerability) (*ent.Vulnerability, error) {
	dbVuln, err := client.Vulnerability.Query().
		Where(vulnerability.Vid(*vuln.Vid)).
		Only(h.ctx)
	if err != nil {
		if !ent.IsNotFound(err) {
			return nil, NewNotFoundError(err, "vulnerability with ID %s", *vuln.Vid)
		}
		// If not found, create it
		dbVuln, err = client.Vulnerability.Create().
			SetModelCreate(&vuln.VulnerabilityModelCreate).
			Save(h.ctx)
		if err != nil {
			return nil, HandleEntError(err, "vulnerability")
		}
	}

	return dbVuln, nil
}
