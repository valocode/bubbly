package store

import (
	"github.com/hashicorp/go-multierror"
	"github.com/valocode/bubbly/ent"
	"github.com/valocode/bubbly/ent/vulnerability"
	"github.com/valocode/bubbly/store/api"
)

func (s *Store) SaveComponentVulnerabilities(req *api.ComponentVulnerabilityRequest) error {
	if err := s.validator.Struct(req); err != nil {
		return HandleValidatorError(err, "component vulnerabilities")
	}
	txErr := s.WithTx(func(tx *ent.Tx) error {
		for _, comp := range req.Components {
			dbComp, err := s.GetComponentOrError(tx, *comp.ComponentID)
			if err != nil {
				return err
			}
			var dbVulns = make([]*ent.Vulnerability, 0, len(comp.Vulnerabilities))
			for _, vuln := range comp.Vulnerabilities {
				dbVuln, err := s.GetVulnerabilityOrCreate(tx, vuln)
				if err != nil {
					return err
				}
				dbVulns = append(dbVulns, dbVuln)
			}
			_, err = tx.Component.UpdateOne(dbComp).
				AddVulnerabilities(dbVulns...).
				Save(s.ctx)
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

func (s *Store) GetComponentOrError(tx *ent.Tx, id int) (*ent.Component, error) {
	client := s.clientOrTx(tx)
	dbComp, err := client.Component.Get(s.ctx, id)
	if err != nil {
		if !ent.IsNotFound(err) {
			return nil, NewNotFoundError(nil, "component with ID %d", id)
		}
		return nil, HandleEntError(err, "component")
	}
	return dbComp, nil

	// Should we be able to get component by vendor:name:version?!
	// var (
	// 	vErr        multierror.Error
	// 	compName    = *comp.Name
	// 	compVendor  = *comp.Vendor
	// 	compVersion = *comp.Version
	// )

	// if vErr.ErrorOrNil() != nil {
	// 	return nil, HandleMultiVError(vErr)
	// }
	// dbComp, err := client.Component.Query().Where(
	// 	component.Name(compName),
	// 	component.Vendor(compVendor),
	// 	component.Version(compVersion),
	// ).Only(s.ctx)
	// if err != nil {
	// 	if ent.IsNotFound(err) {
	// 		return nil, NewNotFoundError(nil, "component %s:%s:%s", compVendor, compName, compVersion)
	// 	}
	// 	return nil, HandleEntError(err, "component")
	// }
	// return dbComp, nil
}

func (s *Store) GetVulnerabilityOrCreate(tx *ent.Tx, vuln *api.Vulnerability) (*ent.Vulnerability, error) {
	client := s.clientOrTx(tx)
	var (
		vErr multierror.Error
		vID  = *vuln.Vid
	)
	if vErr.ErrorOrNil() != nil {
		return nil, HandleMultiVError(vErr)
	}
	dbVuln, err := client.Vulnerability.Query().
		Where(vulnerability.Vid(vID)).
		Only(s.ctx)
	if err != nil {
		if !ent.IsNotFound(err) {
			return nil, NewNotFoundError(err, "vulnerability with ID %s", vID)
		}
		// If not found, create it
		dbVuln, err = client.Vulnerability.Create().
			SetModelCreate(&vuln.VulnerabilityModelCreate).
			Save(s.ctx)
		if err != nil {
			return nil, HandleEntError(err, "vulnerability")
		}
	}

	return dbVuln, nil
}
