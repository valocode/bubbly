package store

import (
	"fmt"

	"github.com/valocode/bubbly/ent"
	"github.com/valocode/bubbly/ent/gitcommit"
	"github.com/valocode/bubbly/ent/project"
	"github.com/valocode/bubbly/ent/release"
	"github.com/valocode/bubbly/ent/releasepolicy"
	"github.com/valocode/bubbly/ent/repo"
	"github.com/valocode/bubbly/policy"
	"github.com/valocode/bubbly/store/api"
)

func (s *Store) SaveReleasePolicy(req *api.ReleasePolicySaveRequest) (*ent.ReleasePolicy, error) {
	if err := s.validator.Struct(req); err != nil {
		return nil, HandleValidatorError(err, "release policy create")
	}
	var dbPolicy *ent.ReleasePolicy
	txErr := s.WithTx(func(tx *ent.Tx) error {
		var err error
		dbPolicy, err = tx.ReleasePolicy.Query().
			Where(releasepolicy.Name(*req.Policy.Name)).
			Only(s.ctx)
		if err != nil {
			if !ent.IsNotFound(err) {
				return HandleEntError(err, "release policy query")
			}
			// If not found, then create the policy
			dbPolicy, err = tx.ReleasePolicy.Create().
				SetModelCreate(req.Policy).
				Save(s.ctx)
			if err != nil {
				return HandleEntError(err, "release policy create")
			}
			// Finish here
			return nil
		}
		// If the policy did exist, then we should update it
		dbPolicy, err = tx.ReleasePolicy.UpdateOne(dbPolicy).
			SetModelCreate(req.Policy).
			Save(s.ctx)
		if err != nil {
			return HandleEntError(err, "release policy update")
		}
		return nil
	})
	if txErr != nil {
		return nil, txErr
	}

	return dbPolicy, nil
}
func (s *Store) EvaluateReleasePolicies(releaseID int) ([]*ent.ReleasePolicyViolation, error) {
	var dbViolations []*ent.ReleasePolicyViolation
	txErr := s.WithTx(func(tx *ent.Tx) error {
		policies, err := tx.ReleasePolicy.Query().
			Where(
				releasepolicy.Or(
					releasepolicy.HasProjectsWith(
						project.HasReposWith(
							repo.HasCommitsWith(
								gitcommit.HasReleaseWith(
									release.ID(releaseID),
								),
							),
						),
					),
					releasepolicy.HasReposWith(
						repo.HasCommitsWith(
							gitcommit.HasReleaseWith(
								release.ID(releaseID),
							),
						),
					),
				),
			).All(s.ctx)
		if err != nil {
			return fmt.Errorf("error getting policies for release: %w", err)
		}
		for _, p := range policies {
			fmt.Println("Evaluate Policy: ", p.String())
			result, err := policy.EvaluatePolicy(p.Module,
				policy.WithResolver(&policy.EntResolver{
					Ctx:       s.ctx,
					Client:    tx.Client(),
					ReleaseID: releaseID,
				}),
			)
			if err != nil {
				return fmt.Errorf("error evaluating policy %s: %w", p.Name, err)
			}
			for _, v := range result.Violations {
				dbViolation, err := tx.ReleasePolicyViolation.Create().
					SetModelCreate(*v).
					Save(s.ctx)
				if err != nil {
					return err
				}
				fmt.Println("Created Policy Violation: ", dbViolation.String())
				dbViolations = append(dbViolations, dbViolation)
			}
		}
		return nil
	})
	if txErr != nil {
		return nil, txErr
	}
	return dbViolations, nil
}
