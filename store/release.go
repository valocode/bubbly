package store

import (
	"fmt"

	"github.com/valocode/bubbly/ent"
	"github.com/valocode/bubbly/ent/gitcommit"
	"github.com/valocode/bubbly/ent/project"
	"github.com/valocode/bubbly/ent/release"
	"github.com/valocode/bubbly/ent/repo"
	"github.com/valocode/bubbly/store/api"
)

type ReleaseQuery struct {
	Where *ent.ReleaseWhereInput
	Order *api.Order

	WithLog        bool
	WithViolations bool
	WithPolicies   bool
}

func (h *Handler) CreateRelease(req *api.ReleaseCreateRequest) (*ent.Release, error) {
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

func (h *Handler) GetReleases(query *ReleaseQuery) ([]*api.Release, error) {
	entQuery := h.client.Release.Query().
		WhereInput(query.Where).
		WithCommit(func(gcq *ent.GitCommitQuery) {
			gcq.WithRepo(func(rq *ent.RepoQuery) {
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
			commit  = dbRelease.Edges.Commit
			repo    = commit.Edges.Repo
			project = repo.Edges.Project
		)
		r := &api.Release{
			Project: ent.NewProjectModelRead().FromEnt(project),
			Repo:    ent.NewRepoModelRead().FromEnt(repo),
			Commit:  ent.NewGitCommitModelRead().FromEnt(commit),
			Release: ent.NewReleaseModelRead().FromEnt(dbRelease),
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

func (h *Handler) LogArtifact(req *api.ArtifactLogRequest) (*ent.Artifact, error) {
	if err := h.validator.Struct(req); err != nil {
		return nil, HandleValidatorError(err, "artifact")
	}
	dbRelease, err := h.releaseFromCommit(req.Commit)
	if err != nil {
		return nil, err
	}
	dbArtifact, err := h.client.Artifact.Create().
		SetModelCreate(req.Artifact).
		SetRelease(dbRelease).
		Save(h.ctx)
	if err != nil {
		return nil, HandleEntError(err, "artifact")
	}
	if err != nil {
		return nil, err
	}
	// Once transaction is complete, evaluate the release.
	_, evalErr := h.EvaluateReleasePolicies(dbRelease.ID)
	if evalErr != nil {
		return nil, NewServerError(evalErr, "evaluating release policies")
	}
	return dbArtifact, nil
}

func (h *Handler) createRelease(req *api.ReleaseCreateRequest) (*ent.Release, error) {
	if err := h.validator.Struct(req); err != nil {
		return nil, HandleValidatorError(err, "release")
	}

	var (
		commitTag string
	)
	// Set defaults / optional values
	if req.Commit.Tag != nil {
		commitTag = *req.Commit.Tag
	}
	// Check if the release exists already
	dbRelease, err := h.client.Release.Query().Where(
		release.Name(*req.Release.Name), release.Version(*req.Release.Version),
		release.HasCommitWith(
			gitcommit.Hash(*req.Commit.Hash),
		),
	).Only(h.ctx)
	if err != nil {
		if !ent.IsNotFound(err) {
			return nil, HandleEntError(err, "release")
		}
		// Otherwise continue and create the release
	}
	if dbRelease != nil {
		return dbRelease, NewConflictError(nil, "release already exists")
	}

	//
	// Create the release, first the project, then repo, then commit
	//
	dbProject, err := h.client.Project.Query().Where(
		project.Name(*req.Project.Name),
	).Only(h.ctx)
	if err != nil {
		if !ent.IsNotFound(err) {
			return nil, HandleEntError(err, "repo")
		}
		dbProject, err = h.client.Project.Create().
			SetModelCreate(req.Project).
			SetOwnerID(h.orgID).
			Save(h.ctx)
		if err != nil {
			return nil, HandleEntError(err, "repo")
		}
	}

	dbRepo, err := h.client.Repo.Query().Where(
		repo.Name(*req.Repo.Name),
	).Only(h.ctx)
	if err != nil {
		if !ent.IsNotFound(err) {
			return nil, HandleEntError(err, "repo")
		}
		dbRepo, err = h.client.Repo.Create().
			SetModelCreate(req.Repo).
			SetOwnerID(h.orgID).
			SetProject(dbProject).
			Save(h.ctx)
		if err != nil {
			return nil, HandleEntError(err, "repo")
		}
	}

	dbCommit, err := h.client.GitCommit.Query().Where(
		gitcommit.Hash(*req.Commit.Hash),
		gitcommit.HasRepoWith(repo.ID(dbRepo.ID)),
	).Only(h.ctx)
	if err != nil {
		if !ent.IsNotFound(err) {
			return nil, HandleEntError(err, "commit")
		}
		if req.Commit.Tag != nil {
			commitTag = *req.Commit.Tag
		}
		dbCommit, err = h.client.GitCommit.Create().
			SetHash(*req.Commit.Hash).
			SetBranch(*req.Commit.Branch).
			SetTime(*req.Commit.Time).
			SetTag(commitTag).
			SetRepo(dbRepo).
			Save(h.ctx)
		if err != nil {
			return nil, HandleEntError(err, "commit")
		}
	}

	dbRelease, err = h.client.Release.Create().
		SetModelCreate(req.Release).
		SetCommit(dbCommit).
		Save(h.ctx)
	if err != nil {
		return nil, HandleEntError(err, "release")
	}

	return dbRelease, nil
}
