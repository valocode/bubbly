package store

import (
	"github.com/valocode/bubbly/ent"
	"github.com/valocode/bubbly/ent/component"
	"github.com/valocode/bubbly/ent/gitcommit"
	"github.com/valocode/bubbly/ent/release"
	"github.com/valocode/bubbly/ent/releasecomponent"
	"github.com/valocode/bubbly/ent/releasevulnerability"
	"github.com/valocode/bubbly/ent/vulnerability"
	"github.com/valocode/bubbly/ent/vulnerabilityreview"
	"github.com/valocode/bubbly/store/api"
)

func (h *Handler) SaveCodeScan(req *api.CodeScanRequest) (*ent.CodeScan, error) {
	if err := h.validator.Struct(req); err != nil {
		return nil, HandleValidatorError(err, "code scan create")
	}
	dbRelease, err := h.releaseFromCommit(req.Commit)
	if err != nil {
		return nil, err
	}
	scan, err := h.saveCodeScan(dbRelease, req.CodeScan)
	if err != nil {
		return nil, err
	}
	// Once transaction is complete, evaluate the release.
	_, evalErr := h.EvaluateReleasePolicies(dbRelease.ID)
	if evalErr != nil {
		return nil, NewServerError(evalErr, "evaluating release policies")
	}

	return scan, nil
}

func (h *Handler) saveCodeScan(dbRelease *ent.Release, scan *api.CodeScan) (*ent.CodeScan, error) {
	var codeScan *ent.CodeScan
	txErr := WithTx(h.ctx, h.client, func(tx *ent.Tx) error {

		var err error
		codeScan, err = tx.CodeScan.Create().
			SetModelCreate(&scan.CodeScanModelCreate).
			SetRelease(dbRelease).
			Save(h.ctx)
		if err != nil {
			return HandleEntError(err, "code scan")
		}

		for _, issue := range scan.Issues {
			_, err := tx.CodeIssue.Create().
				SetModelCreate(&issue.CodeIssueModelCreate).
				SetScan(codeScan).
				Save(h.ctx)
			if err != nil {
				return HandleEntError(err, "code issue")
			}
		}
		//
		// Save code scan components
		//
		for _, comp := range scan.Components {
			var (
				existingComp *ent.Component
				err          error
			)
			if err := h.validator.Struct(comp); err != nil {
				return HandleValidatorError(err, "release component create")
			}
			existingComp, err = tx.Component.Query().
				WhereName(comp.Name).
				WhereVendor(comp.Vendor).
				WhereVersion(comp.Version).
				Only(h.ctx)
			if err != nil {
				if !ent.IsNotFound(err) {
					return HandleEntError(err, "component")
				}
				// It is not found, so create the component...
				existingComp, err = tx.Component.Create().
					SetModelCreate(&comp.ComponentModelCreate).
					SetOwnerID(h.orgID).
					Save(h.ctx)
				if err != nil {
					return HandleEntError(err, "component")
				}
			}
			// Get the release component, which might already exist if the
			// component already exists
			relComp, err := tx.ReleaseComponent.Query().Where(
				releasecomponent.HasComponentWith(component.ID(existingComp.ID)),
				releasecomponent.HasReleaseWith(release.ID(dbRelease.ID)),
			).Only(h.ctx)
			if err != nil {
				if !ent.IsNotFound(err) {
					return HandleEntError(err, "query release component")
				}
				relComp, err = tx.ReleaseComponent.Create().
					SetComponent(existingComp).
					AddScans(codeScan).
					SetRelease(dbRelease).
					Save(h.ctx)
				if err != nil {
					return HandleEntError(err, "create release component")
				}
			}
			//
			// Save component vulnerabilities
			//
			for _, vuln := range comp.Vulnerabilities {
				var (
					existingVuln *ent.Vulnerability
					err          error
				)
				if err := h.validator.Struct(comp); err != nil {
					return HandleValidatorError(err, "release vulnerability create")
				}
				existingVuln, err = tx.Vulnerability.Query().Where(
					vulnerability.Vid(*vuln.Vid),
				).Only(h.ctx)
				if err != nil {
					if !ent.IsNotFound(err) {
						return HandleEntError(err, "vulnerability")
					}
					// If not found, create!
					existingVuln, err = tx.Vulnerability.Create().
						SetModelCreate(&vuln.VulnerabilityModelCreate).
						SetOwnerID(h.orgID).
						Save(h.ctx)
					if err != nil {
						return HandleEntError(err, "vulnerability")
					}
				}
				// Check if the release vulnerability already exists, which is
				// the combination of release ID and vulnerability ID
				dbRelVuln, err := tx.ReleaseVulnerability.Query().Where(
					releasevulnerability.HasReleaseWith(release.ID(dbRelease.ID)),
					releasevulnerability.HasVulnerabilityWith(vulnerability.ID(existingVuln.ID)),
				).Only(h.ctx)
				if err != nil {
					if !ent.IsNotFound(err) {
						return HandleEntError(err, "query release vulnerability")
					}
					dbRelVuln, err = tx.ReleaseVulnerability.Create().
						SetRelease(dbRelease).
						SetVulnerability(existingVuln).
						SetScan(codeScan).
						SetComponent(relComp).
						Save(h.ctx)
					if err != nil {
						return HandleEntError(err, "create release vulnerability")
					}
				}
				if vuln.Patch != nil {
					_, err := tx.VulnerabilityReview.Create().
						SetName(*vuln.Patch.Message).
						SetDecision(vulnerabilityreview.DecisionPatched).
						SetVulnerability(existingVuln).
						AddInstanceIDs(dbRelVuln.ID).
						Save(h.ctx)
					if err != nil {
						return HandleEntError(err, "create vulnerability patch")
					}
				}
			}
		}
		return nil
	})
	if txErr != nil {
		return nil, txErr
	}

	return codeScan, nil
}

func (h *Handler) releaseFromCommit(commitHash *string) (*ent.Release, error) {
	if commitHash == nil {
		return nil, NewValidationError(nil, "commit is required")
	}
	commit, err := h.client.GitCommit.Query().
		Where(gitcommit.Hash(*commitHash)).
		WithRelease().
		Only(h.ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, NewNotFoundError(nil, "no release for commit %s. Please create one.", *commitHash)
		}
		return nil, HandleEntError(err, "commit")
	}

	if commit.Edges.Release == nil {
		return nil, NewNotFoundError(nil, "release")
	}
	return commit.Edges.Release, nil
}
