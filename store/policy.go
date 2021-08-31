package store

import (
	"fmt"

	"github.com/valocode/bubbly/ent"
	"github.com/valocode/bubbly/ent/gitcommit"
	"github.com/valocode/bubbly/ent/predicate"
	"github.com/valocode/bubbly/ent/project"
	"github.com/valocode/bubbly/ent/release"
	"github.com/valocode/bubbly/ent/releasepolicy"
	"github.com/valocode/bubbly/ent/releasepolicyviolation"
	"github.com/valocode/bubbly/ent/repo"
	"github.com/valocode/bubbly/policy"
	"github.com/valocode/bubbly/store/api"
)

func (h *Handler) SaveReleasePolicy(req *api.ReleasePolicySaveRequest) (*ent.ReleasePolicy, error) {
	if err := h.validator.Struct(req); err != nil {
		return nil, HandleValidatorError(err, "release policy create")
	}
	var dbPolicy *ent.ReleasePolicy
	txErr := WithTx(h.ctx, h.client, func(tx *ent.Tx) error {
		var err error
		dbPolicy, err = tx.ReleasePolicy.Query().
			Where(releasepolicy.Name(*req.Policy.Name)).
			Only(h.ctx)
		if err != nil {
			if !ent.IsNotFound(err) {
				return HandleEntError(err, "release policy query")
			}
			// If not found, then create the policy
			dbPolicy, err = tx.ReleasePolicy.Create().
				SetModelCreate(req.Policy).
				Save(h.ctx)
			if err != nil {
				return HandleEntError(err, "release policy create")
			}
		} else {
			// If the policy did exist, then we should update it
			dbPolicy, err = tx.ReleasePolicy.UpdateOne(dbPolicy).
				SetModelCreate(req.Policy).
				Save(h.ctx)
			if err != nil {
				return HandleEntError(err, "release policy update")
			}
		}

		dbPolicy, err = h.updatePolicyAffects(tx.Client(), dbPolicy, req.Affects)
		if err != nil {
			return err
		}

		return nil
	})
	if txErr != nil {
		return nil, txErr
	}

	return dbPolicy, nil
}

func (h *Handler) SetReleasePolicyAffects(req *api.ReleasePolicySetRequest) (*ent.ReleasePolicy, error) {
	if err := h.validator.Struct(req); err != nil {
		return nil, HandleValidatorError(err, "release policy set affects")
	}
	var dbPolicy *ent.ReleasePolicy
	var err error
	dbPolicy, err = h.client.ReleasePolicy.Query().
		Where(releasepolicy.Name(*req.Policy)).
		Only(h.ctx)
	if err != nil {
		return nil, HandleEntError(err, "release policy query")
	}

	dbPolicy, err = h.updatePolicyAffects(h.client, dbPolicy, req.Affects)
	if err != nil {
		return nil, err
	}

	return dbPolicy, nil
}

func (h *Handler) EvaluateReleasePolicies(releaseID int) ([]*ent.ReleasePolicyViolation, error) {
	var dbViolations []*ent.ReleasePolicyViolation
	txErr := WithTx(h.ctx, h.client, func(tx *ent.Tx) error {
		dbPolicies, err := tx.ReleasePolicy.Query().
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
			).All(h.ctx)
		if err != nil {
			return fmt.Errorf("error getting policies for release: %w", err)
		}
		for _, dbPolicy := range dbPolicies {
			result, err := policy.EvaluatePolicy(dbPolicy.Module,
				policy.WithResolver(&policy.EntResolver{
					Ctx:       h.ctx,
					Client:    tx.Client(),
					ReleaseID: releaseID,
				}),
			)
			if err != nil {
				return fmt.Errorf("error evaluating policy %s: %w", dbPolicy.Name, err)
			}
			// First we need to clear the existing violations so that we don't
			// create duplicates
			_, deleteErr := tx.ReleasePolicyViolation.Delete().Where(
				releasepolicyviolation.HasPolicyWith(releasepolicy.ID(dbPolicy.ID)),
				releasepolicyviolation.HasReleaseWith(release.ID(releaseID)),
			).Exec(h.ctx)
			if deleteErr != nil {
				return NewServerError(deleteErr, "deleteing release violations")
			}
			for _, v := range result.Violations {
				dbViolation, err := tx.ReleasePolicyViolation.Create().
					SetModelCreate(*v).
					SetPolicy(dbPolicy).
					SetReleaseID(releaseID).
					Save(h.ctx)
				if err != nil {
					return err
				}
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

func (h *Handler) updatePolicyAffects(client *ent.Client, dbPolicy *ent.ReleasePolicy, affects *api.ReleasePolicyAffects) (*ent.ReleasePolicy, error) {
	if affects == nil {
		return dbPolicy, nil
	}

	projectSetIDs, err := h.policyAffectsProjectIDs(client, affects.Projects)
	if err != nil {
		return nil, err
	}
	projectSetNotIDs, err := h.policyAffectsProjectIDs(client, affects.NotProjects)
	if err != nil {
		return nil, err
	}
	repoSetIDs, err := h.policyAffectsReposIDs(client, affects.Repos)
	if err != nil {
		return nil, err
	}
	repoSetNotIDs, err := h.policyAffectsReposIDs(client, affects.NotRepos)
	if err != nil {
		return nil, err
	}

	// Update the policy with the entities that it affects
	dbPolicy, err = client.ReleasePolicy.UpdateOne(dbPolicy).
		AddProjectIDs(projectSetIDs...).
		RemoveProjectIDs(projectSetNotIDs...).
		AddRepoIDs(repoSetIDs...).
		RemoveRepoIDs(repoSetNotIDs...).
		Save(h.ctx)
	if err != nil {
		return nil, HandleEntError(err, "set policy affects")
	}
	return dbPolicy, nil
}

func (h *Handler) policyAffectsProjectIDs(client *ent.Client, projects []string) ([]int, error) {
	var (
		projectSetIDs   []int
		projectSetNames []predicate.Project
		err             error
	)
	// Check if there are any projects, otherwise we have nothing to do
	if len(projects) == 0 {
		return projectSetIDs, nil
	}

	for _, setProject := range projects {
		projectSetNames = append(projectSetNames, project.Name(setProject))
	}

	projectSetIDs, err = client.Project.Query().Where(
		project.Or(projectSetNames...),
	).IDs(h.ctx)
	if err != nil {
		return nil, HandleEntError(err, "policy project IDs")
	}
	return projectSetIDs, nil
}

func (h *Handler) policyAffectsReposIDs(client *ent.Client, repos []string) ([]int, error) {
	var (
		repoSetIDs   []int
		repoSetNames []predicate.Repo
		err          error
	)
	// Check if there are any repos, otherwise we have nothing to do
	if len(repos) == 0 {
		return repoSetIDs, nil
	}

	for _, setRepo := range repos {
		repoSetNames = append(repoSetNames, repo.Name(setRepo))
	}

	repoSetIDs, err = client.Repo.Query().Where(
		repo.Or(repoSetNames...),
	).IDs(h.ctx)
	if err != nil {
		return nil, HandleEntError(err, "policy repo IDs")
	}
	return repoSetIDs, nil
}
