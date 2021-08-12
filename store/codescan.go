package store

import (
	"github.com/valocode/bubbly/ent"
	"github.com/valocode/bubbly/ent/component"
	"github.com/valocode/bubbly/ent/gitcommit"
	"github.com/valocode/bubbly/ent/vulnerability"
	"github.com/valocode/bubbly/store/api"
)

func (s *Store) SaveCodeScan(req *api.CodeScanRequest) (*ent.CodeScan, error) {
	release, err := s.releaseFromCommit(req.Commit)
	if err != nil {
		return nil, err
	}
	return s.saveCodeScan(release, req.CodeScan)
}

func (s *Store) saveCodeScan(release *ent.Release, scan *api.CodeScan) (*ent.CodeScan, error) {
	var codeScan *ent.CodeScan
	txErr := s.WithTx(func(tx *ent.Tx) error {

		scanCreate := tx.CodeScan.Create()
		if scan.Tool == nil {
			return NewValidationError(nil, "scan tool is required")
		}
		scanCreate.SetTool(*scan.Tool).
			SetRelease(release)
		var err error
		codeScan, err = scanCreate.Save(s.ctx)
		if err != nil {
			return HandleEntError(err, "code scan")
		}

		for _, issue := range scan.Issues {
			issueCreate := tx.CodeIssue.Create().
				SetScan(codeScan)
			issue.SetMutatorFields(issueCreate.Mutation())
			_, err := issueCreate.Save(s.ctx)
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
			existingComp, err = tx.Component.Query().Where(
				component.Vendor(*comp.Vendor),
				component.Name(*comp.Name),
				component.Version(*comp.Version),
			).Only(s.ctx)
			if err != nil {
				if !ent.IsNotFound(err) {
					return HandleEntError(err, "component")
				}
				// It is not found, so create the component...
				compCreate := tx.Component.Create()
				comp.SetMutatorFields(compCreate.Mutation())
				existingComp, err = compCreate.Save(s.ctx)
				if err != nil {
					return HandleEntError(err, "component")
				}
			}
			relComp, err := tx.ReleaseComponent.Create().
				SetComponent(existingComp).
				AddScans(codeScan).
				SetRelease(release).
				Save(s.ctx)
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
				existingVuln, err = tx.Vulnerability.Query().Where(
					vulnerability.Vid(*vuln.Vid),
				).Only(s.ctx)
				if err != nil {
					if !ent.IsNotFound(err) {
						return HandleEntError(err, "vulnerability")
					}
					// If not found, create!
					vulnCreate := tx.Vulnerability.Create()
					vuln.SetMutatorFields(vulnCreate.Mutation())
					existingVuln, err = vulnCreate.Save(s.ctx)
					if err != nil {
						return HandleEntError(err, "vulnerability")
					}
				}
				_, err = tx.ReleaseVulnerability.Create().
					SetRelease(release).
					SetVulnerability(existingVuln).
					AddScans(codeScan).
					AddComponents(relComp).
					Save(s.ctx)
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

func (s *Store) releaseFromCommit(commitHash *string) (*ent.Release, error) {
	if commitHash == nil {
		return nil, NewValidationError(nil, "commit is required")
	}
	commit, err := s.client.GitCommit.Query().
		Where(gitcommit.Hash(*commitHash)).
		WithRelease().
		Only(s.ctx)
	if err != nil {
		return nil, HandleEntError(err, "commit")
	}

	if commit.Edges.Release == nil {
		return nil, NewNotFoundError(nil, "release")
	}
	return commit.Edges.Release, nil
}
