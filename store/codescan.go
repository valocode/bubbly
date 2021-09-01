package store

import (
	"github.com/valocode/bubbly/ent"
	"github.com/valocode/bubbly/ent/component"
	"github.com/valocode/bubbly/ent/gitcommit"
	"github.com/valocode/bubbly/ent/vulnerability"
	"github.com/valocode/bubbly/store/api"
)

func (h *Handler) SaveCodeScan(req *api.CodeScanRequest) (*ent.CodeScan, error) {
	if err := h.validator.Struct(req); err != nil {
		return nil, HandleValidatorError(err, "code scan create")
	}
	release, err := h.releaseFromCommit(req.Commit)
	if err != nil {
		return nil, err
	}
	scan, err := h.saveCodeScan(release, req.CodeScan)
	if err != nil {
		return nil, err
	}
	// Once transaction is complete, evaluate the release.
	_, evalErr := h.EvaluateReleasePolicies(release.ID)
	if evalErr != nil {
		return nil, NewServerError(evalErr, "evaluating release policies")
	}

	return scan, nil
}

func (h *Handler) saveCodeScan(release *ent.Release, scan *api.CodeScan) (*ent.CodeScan, error) {
	var codeScan *ent.CodeScan
	txErr := WithTx(h.ctx, h.client, func(tx *ent.Tx) error {

		var err error
		codeScan, err = tx.CodeScan.Create().
			SetModelCreate(&scan.CodeScanModelCreate).
			SetRelease(release).
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
			existingComp, err = tx.Component.Query().Where(
				component.Vendor(*comp.Vendor),
				component.Name(*comp.Name),
				component.Version(*comp.Version),
			).Only(h.ctx)
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
			relComp, err := tx.ReleaseComponent.Create().
				SetComponent(existingComp).
				AddScans(codeScan).
				SetRelease(release).
				Save(h.ctx)
			if err != nil {
				return HandleEntError(err, "release component")
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
						SetModelCreate(&vuln.VulnerabilityModelCreate).Save(h.ctx)
					if err != nil {
						return HandleEntError(err, "vulnerability")
					}
				}
				_, err = tx.ReleaseVulnerability.Create().
					SetRelease(release).
					SetVulnerability(existingVuln).
					SetScan(codeScan).
					SetComponent(relComp).
					Save(h.ctx)
				if err != nil {
					return HandleEntError(err, "release vulnerability")
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
