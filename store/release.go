package store

import (
	"fmt"

	"github.com/valocode/bubbly/ent"
	"github.com/valocode/bubbly/ent/gitcommit"
	"github.com/valocode/bubbly/ent/release"
	"github.com/valocode/bubbly/ent/repository"
	"github.com/valocode/bubbly/store/api"
)

type ReleaseQuery struct {
	Where *ent.ReleaseWhereInput
	Order *api.Order

	WithLog        bool
	WithViolations bool
	WithPolicies   bool
}

func (h *Handler) GetReleases(query *ReleaseQuery) ([]*api.Release, error) {
	entQuery := h.client.Release.Query().
		WhereInput(query.Where).
		WithCommit(func(gcq *ent.GitCommitQuery) {
			gcq.WithRepository(func(rq *ent.RepositoryQuery) {
				rq.WithProject()
			})
		})
	if query.WithLog {
		entQuery.WithLog()
	}
	if query.WithViolations {
		entQuery.WithViolations()
	}
	if query.Order != nil {
		switch query.Order.Field {
		case "name":
			entQuery.Order(query.Order.Func("name"))
		case "version":
			entQuery.QueryCommit().Order(query.Order.Func("time")).QueryRelease()
		default:
			return nil, NewValidationError(nil, fmt.Sprintf("unknown order field: %s", query.Order.Field))
		}
	}
	dbReleases, err := entQuery.All(h.ctx)
	if err != nil {
		return nil, HandleEntError(err, "getting releases by commit")
	}
	var releases = make([]*api.Release, 0, len(dbReleases))
	for _, dbRelease := range dbReleases {
		var (
			commit     = dbRelease.Edges.Commit
			repository = commit.Edges.Repository
			project    = repository.Edges.Project
		)
		r := &api.Release{
			Project:    ent.NewProjectModelRead().FromEnt(project),
			Repository: ent.NewRepositoryModelRead().FromEnt(repository),
			Commit:     ent.NewGitCommitModelRead().FromEnt(commit),
			Release:    ent.NewReleaseModelRead().FromEnt(dbRelease),
		}
		// If the request said to get policies, then fetch the policies for the release.
		if query.WithPolicies {
			dbPolicies, err := h.policiesForRelease(h.client, dbRelease.ID)
			if err != nil {
				return nil, err
			}
			var policies = make([]*ent.ReleasePolicyModelRead, 0, len(dbPolicies))
			for _, p := range dbPolicies {
				policies = append(policies, ent.NewReleasePolicyModelRead().FromEnt(p))
			}
			r.Policies = policies
		}
		// Append the release violation
		for _, dbViolation := range dbRelease.Edges.Violations {
			r.Violations = append(r.Violations, ent.NewReleasePolicyViolationModelRead().FromEnt(dbViolation))
		}
		// Append the release log
		for _, dbEntry := range dbRelease.Edges.Log {
			r.Entries = append(r.Entries, ent.NewReleaseEntryModelRead().FromEnt(dbEntry))
		}
		releases = append(releases, r)
	}
	return releases, nil
}

func (h *Handler) CreateRelease(req *api.ReleaseCreateRequest) (*ent.Release, error) {
	if err := h.validator.Struct(req); err != nil {
		return nil, HandleValidatorError(err, "release create")
	}
	dbRelease, err := h.createRelease(req)
	if err != nil {
		return nil, err
	}
	// Once transaction is complete, evaluate the release.
	_, evalErr := h.EvaluateReleasePolicies(dbRelease.ID)
	if evalErr != nil {
		return nil, NewServerError(evalErr, "evaluating release policies")
	}
	return dbRelease, nil
}

func (h *Handler) createRelease(req *api.ReleaseCreateRequest) (*ent.Release, error) {
	var (
		commitTag string
	)
	// Set defaults / optional values
	if req.Commit.Tag != nil {
		commitTag = *req.Commit.Tag
	}
	// Check if the release exists already
	releaseExists, err := h.client.Release.Query().Where(
		release.Name(*req.Release.Name),
		release.Version(*req.Release.Version),
		release.HasCommitWith(
			gitcommit.Hash(*req.Commit.Hash),
		),
	).Exist(h.ctx)
	if err != nil {
		return nil, HandleEntError(err, "release")
	}
	if releaseExists {
		return nil, NewConflictError(nil, "release already exists")
	}

	// Release does not exist, so let's create it. First checking that the project,
	// repo and commit exist

	dbRepo, err := h.createRepository(h.client, &api.RepositoryCreateRequest{
		Project:    req.Project,
		Repository: req.Repository,
	})
	if err != nil {
		return nil, err
	}

	dbCommit, err := h.client.GitCommit.Query().Where(
		gitcommit.Hash(*req.Commit.Hash),
		gitcommit.HasRepositoryWith(repository.ID(dbRepo.ID)),
	).Only(h.ctx)
	if err != nil {
		if !ent.IsNotFound(err) {
			return nil, HandleEntError(err, "commit create")
		}
		if req.Commit.Tag != nil {
			commitTag = *req.Commit.Tag
		}
		dbCommit, err = h.client.GitCommit.Create().
			SetHash(*req.Commit.Hash).
			SetBranch(*req.Commit.Branch).
			SetTime(*req.Commit.Time).
			SetTag(commitTag).
			SetRepository(dbRepo).
			Save(h.ctx)
		if err != nil {
			return nil, HandleEntError(err, "commit create")
		}
	}

	dbRelease, err := h.client.Release.Create().
		SetModelCreate(req.Release).
		SetCommit(dbCommit).
		Save(h.ctx)
	if err != nil {
		return nil, HandleEntError(err, "release")
	}

	return dbRelease, nil
}
