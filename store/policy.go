package store

import (
	"fmt"
	"strings"

	"github.com/valocode/bubbly/ent"
	"github.com/valocode/bubbly/ent/event"
	"github.com/valocode/bubbly/ent/gitcommit"
	"github.com/valocode/bubbly/ent/organization"
	"github.com/valocode/bubbly/ent/project"
	"github.com/valocode/bubbly/ent/release"
	"github.com/valocode/bubbly/ent/releasepolicy"
	"github.com/valocode/bubbly/ent/releasepolicyviolation"
	"github.com/valocode/bubbly/ent/repo"
	"github.com/valocode/bubbly/policy"
	"github.com/valocode/bubbly/store/api"
)

type ReleasePolicyQuery struct {
	Where *ent.ReleasePolicyWhereInput

	WithAffects bool
}

func (h *Handler) SaveReleasePolicy(req *api.ReleasePolicySaveRequest) (*ent.ReleasePolicy, error) {
	if err := h.validator.Struct(req); err != nil {
		return nil, HandleValidatorError(err, "release policy create")
	}
	var dbPolicy *ent.ReleasePolicy
	txErr := WithTx(h.ctx, h.client, func(tx *ent.Tx) error {
		var err error
		dbPolicy, err = tx.ReleasePolicy.Query().
			Where(releasepolicy.Name(*req.Policy.Name), releasepolicy.HasOwnerWith(organization.ID(h.orgID))).
			Only(h.ctx)
		if err != nil {
			if !ent.IsNotFound(err) {
				return HandleEntError(err, "release policy query")
			}
			// If not found, then create the policy
			dbPolicy, err = tx.ReleasePolicy.Create().
				SetModelCreate(&req.Policy.ReleasePolicyModelCreate).
				SetOwnerID(h.orgID).
				Save(h.ctx)
			if err != nil {
				return HandleEntError(err, "release policy create")
			}
		} else {
			// If the policy did exist, then we should update it
			dbPolicy, err = tx.ReleasePolicy.UpdateOne(dbPolicy).
				SetModelCreate(&req.Policy.ReleasePolicyModelCreate).
				Save(h.ctx)
			if err != nil {
				return HandleEntError(err, "release policy update")
			}
		}

		dbPolicy, err = h.updatePolicyAffects(tx.Client(), dbPolicy, req.Policy.Affects)
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

func (h *Handler) GetReleasePolicies(query *ReleasePolicyQuery) ([]*api.ReleasePolicy, error) {
	entQuery := h.client.ReleasePolicy.Query().WhereInput(query.Where)
	if query.WithAffects {
		entQuery.WithProjects().WithRepos()
	}
	dbPolicies, err := entQuery.All(h.ctx)
	if err != nil {
		return nil, HandleEntError(err, "get release policy")
	}
	var policies []*api.ReleasePolicy
	for _, p := range dbPolicies {
		var affects api.ReleasePolicyAffects
		for _, p := range p.Edges.Projects {
			affects.Projects = append(affects.Projects, p.Name)
		}
		for _, r := range p.Edges.Repos {
			affects.Repos = append(affects.Repos, r.Name)
		}
		policies = append(policies, &api.ReleasePolicy{
			ReleasePolicyModelRead: *ent.NewReleasePolicyModelRead().FromEnt(p),
			Affects:                &affects,
		})
	}
	return policies, nil
}

func (h *Handler) UpdateReleasePolicy(req *api.ReleasePolicyUpdateRequest) (*ent.ReleasePolicy, error) {
	if err := h.validator.Struct(req); err != nil {
		return nil, HandleValidatorError(err, "release policy set affects")
	}
	var (
		dbPolicy *ent.ReleasePolicy
		err      error
	)
	dbPolicy, err = h.client.ReleasePolicy.UpdateOneID(*req.ID).
		SetModelUpdate(&req.Policy.ReleasePolicyModelUpdate).
		Save(h.ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, NewNotFoundError(nil, "release policy not found")
		}
		return nil, HandleEntError(err, "updating policy")
	}

	dbPolicy, err = h.updatePolicyAffects(h.client, dbPolicy, req.Policy.Affects)
	if err != nil {
		return nil, err
	}

	return dbPolicy, nil
}

func (h *Handler) EvaluateReleasePolicies(releaseID int) ([]*ent.ReleasePolicyViolation, error) {
	var (
		dbViolations []*ent.ReleasePolicyViolation
		dbPolicies   []*ent.ReleasePolicy
	)
	txErr := WithTx(h.ctx, h.client, func(tx *ent.Tx) error {
		var err error
		dbPolicies, err = h.policiesForRelease(tx.Client(), releaseID)
		if err != nil {
			return err
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
					SetModelCreate(v).
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
	// Save the event that policies were evaluated
	var (
		policies  []string
		policyStr string
	)

	for _, p := range dbPolicies {
		policies = append(policies, p.Name)
	}
	if len(policies) == 0 {
		policyStr = "None"
	} else {
		policyStr = strings.Join(policies, ",")
	}
	if _, err := h.SaveEvent(&api.EventSaveRequest{
		ReleaseID: &releaseID,
		Event: ent.NewEventModelCreate().
			SetMessage(fmt.Sprintf("Policies: %s\nViolations: %d", policyStr, len(dbViolations))).
			SetType(event.TypeEvaluateRelease),
	}); err != nil {
		return nil, err
	}

	return dbViolations, nil
}

func (h *Handler) policiesForRelease(client *ent.Client, releaseID int) ([]*ent.ReleasePolicy, error) {
	dbPolicies, err := client.ReleasePolicy.Query().
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
		return nil, HandleEntError(err, "getting policies for release")
	}
	return dbPolicies, nil
}

func (h *Handler) updatePolicyAffects(client *ent.Client, dbPolicy *ent.ReleasePolicy, affects *api.ReleasePolicyAffectsSet) (*ent.ReleasePolicy, error) {
	if affects == nil {
		return dbPolicy, nil
	}

	projectSetIDs, err := h.projectIDs(client, affects.Projects)
	if err != nil {
		return nil, err
	}
	projectSetNotIDs, err := h.projectIDs(client, affects.NotProjects)
	if err != nil {
		return nil, err
	}
	repoSetIDs, err := h.repoIDs(client, affects.Repos)
	if err != nil {
		return nil, err
	}
	repoSetNotIDs, err := h.repoIDs(client, affects.NotRepos)
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
