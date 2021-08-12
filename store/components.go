package store

import (
	"github.com/hashicorp/go-multierror"
	"github.com/valocode/bubbly/ent"
	"github.com/valocode/bubbly/ent/component"
	"github.com/valocode/bubbly/ent/vulnerability"
	"github.com/valocode/bubbly/store/api"
)

func (s *Store) SaveComponentVulnerabilities(components *api.ComponentVulnerabilityRequest) error {
	txErr := s.WithTx(func(tx *ent.Tx) error {
		for _, comp := range components.Components {
			dbComp, err := s.GetComponentOrError(tx, &comp.Component)
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

func (s *Store) GetComponentOrError(tx *ent.Tx, comp *api.Component) (*ent.Component, error) {
	client := s.clientOrTx(tx)
	if comp.ID != nil {
		dbComp, err := client.Component.Get(s.ctx, *comp.ID)
		if err != nil {
			if !ent.IsNotFound(err) {
				return nil, NewNotFoundError(nil, "component with ID %d", comp.ID)
			}
			return nil, HandleEntError(err, "component")
		}
		return dbComp, nil
	}
	var (
		vErr        multierror.Error
		compName    = comp.GetNameOrErr(&vErr)
		compVendor  = comp.GetVendorOrErr(&vErr)
		compVersion = comp.GetVersionOrErr(&vErr)
	)

	if vErr.ErrorOrNil() != nil {
		return nil, HandleMultiVError(vErr)
	}
	dbComp, err := client.Component.Query().Where(
		component.Name(*compName),
		component.Vendor(*compVendor),
		component.Version(*compVersion),
	).Only(s.ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, NewNotFoundError(nil, "component %s:%s:%s", *compVendor, *compName, *compVersion)
		}
		return nil, HandleEntError(err, "component")
	}
	return dbComp, nil
}

func (s *Store) GetVulnerabilityOrCreate(tx *ent.Tx, vuln *api.Vulnerability) (*ent.Vulnerability, error) {
	client := s.clientOrTx(tx)
	var (
		vErr multierror.Error
		vID  = vuln.GetVidOrErr(&vErr)
	)
	if vErr.ErrorOrNil() != nil {
		return nil, HandleMultiVError(vErr)
	}
	dbVuln, err := client.Vulnerability.Query().
		Where(vulnerability.Vid(*vID)).
		Only(s.ctx)
	if err != nil {
		if !ent.IsNotFound(err) {
			return nil, NewNotFoundError(err, "vulnerability with ID %s", *vID)
		}
		// If not found, create it
		vulnCreate := client.Vulnerability.Create()
		vuln.SetMutatorFields(vulnCreate.Mutation())
		dbVuln, err = vulnCreate.Save(s.ctx)
		if err != nil {
			return nil, HandleEntError(err, "vulnerability")
		}
	}

	return dbVuln, nil
}
