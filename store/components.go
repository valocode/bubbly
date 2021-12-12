package store

import (
	"github.com/valocode/bubbly/ent"
	"github.com/valocode/bubbly/ent/component"
	"github.com/valocode/bubbly/store/api"
)

type ComponentQuery struct {
	Where *ent.ComponentWhereInput
}

func (h *Handler) GetComponents(query *ComponentQuery) ([]*api.Component, error) {
	dbComonents, err := h.client.Component.Query().
		WhereInput(query.Where).
		All(h.ctx)
	if err != nil {
		return nil, HandleEntError(err, "component get")
	}
	var components = make([]*api.Component, 0, len(dbComonents))
	for _, comp := range dbComonents {
		components = append(components, &api.Component{
			ComponentModelRead: *ent.NewComponentModelRead().FromEnt(comp),
		})
	}
	return components, nil
}

func (h *Handler) CreateComponent(req *api.ComponentCreateRequest) (*ent.Component, error) {
	if err := h.validator.Struct(req); err != nil {
		return nil, HandleValidatorError(err, "component create")
	}
	return h.createComponent(h.client, req)
}
func (h *Handler) createComponent(client *ent.Client, req *api.ComponentCreateRequest) (*ent.Component, error) {

	var dbComponent *ent.Component
	txErr := WithTx(h.ctx, h.client, func(tx *ent.Tx) error {
		var err error
		// TODO: check this is complete
		dbComponent, err = tx.Component.Query().Where(
			component.Scheme(*req.Component.Scheme),
			component.Namespace(*req.Component.Namespace),
			component.Name(*req.Component.Name),
			component.Version(*req.Component.Version),
		).Only(h.ctx)
		if err != nil {
			if !ent.IsNotFound(err) {
				return HandleEntError(err, "component query")
			}
			// If not found, then create it
			dbComponent, err = tx.Component.Create().
				SetModelCreate(req.Component.ComponentModelCreate).
				Save(h.ctx)
			if err != nil {
				return HandleEntError(err, "component create")
			}
		}
		return nil
	})
	if txErr != nil {
		return nil, txErr
	}
	return dbComponent, nil
}

// func (h *Handler) SaveComponentVulnerabilities(req *api.ComponentVulnerabilityRequest) error {
// 	if err := h.validator.Struct(req); err != nil {
// 		return HandleValidatorError(err, "component vulnerabilities")
// 	}
// 	txErr := WithTx(h.ctx, h.client, func(tx *ent.Tx) error {
// 		for _, comp := range req.Components {
// 			dbComp, err := h.GetComponentOrError(tx.Client(), *comp.ComponentID)
// 			if err != nil {
// 				return err
// 			}
// 			dbVulns := make([]*ent.Vulnerability, 0, len(comp.Vulnerabilities))
// 			for _, vuln := range comp.Vulnerabilities {
// 				dbVuln, err := h.GetVulnerabilityOrCreate(tx.Client(), vuln)
// 				if err != nil {
// 					return err
// 				}
// 				dbVulns = append(dbVulns, dbVuln)
// 			}
// 			_, err = tx.Component.UpdateOne(dbComp).
// 				AddVulnerabilities(dbVulns...).
// 				Save(h.ctx)
// 			if err != nil {
// 				return HandleEntError(err, "component vulnerability")
// 			}
// 		}
// 		return nil
// 	})
// 	if txErr != nil {
// 		return txErr
// 	}
// 	return nil
// }

// func (h *Handler) GetComponentOrError(client *ent.Client, id int) (*ent.Component, error) {
// 	dbComp, err := client.Component.Get(h.ctx, id)
// 	if err != nil {
// 		if !ent.IsNotFound(err) {
// 			return nil, NewNotFoundError(nil, "component with ID %d", id)
// 		}
// 		return nil, HandleEntError(err, "component")
// 	}
// 	return dbComp, nil
// }

// func (h *Handler) GetVulnerabilityOrCreate(client *ent.Client, vuln *api.VulnerabilityCreatePatch) (*ent.Vulnerability, error) {
// 	dbVuln, err := client.Vulnerability.Query().
// 		Where(vulnerability.Vid(*vuln.Vid)).
// 		Only(h.ctx)
// 	if err != nil {
// 		if !ent.IsNotFound(err) {
// 			return nil, NewNotFoundError(err, "vulnerability with ID %s", *vuln.Vid)
// 		}
// 		// If not found, create it
// 		dbVuln, err = client.Vulnerability.Create().
// 			SetModelCreate(&vuln.VulnerabilityModelCreate).
// 			SetOwnerID(h.orgID).
// 			Save(h.ctx)
// 		if err != nil {
// 			return nil, HandleEntError(err, "vulnerability")
// 		}
// 	}

// 	return dbVuln, nil
// }
