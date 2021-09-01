package store

import (
	"github.com/valocode/bubbly/ent"
	"github.com/valocode/bubbly/ent/gitcommit"
	"github.com/valocode/bubbly/ent/project"
	"github.com/valocode/bubbly/ent/release"
	"github.com/valocode/bubbly/ent/repo"
	"github.com/valocode/bubbly/store/api"
)

func (h *Handler) CreateRelease(req *api.ReleaseCreateRequest) (*ent.Release, error) {
	return h.createRelease(req)
}

func (h *Handler) GetReleases(req *api.ReleaseGetRequest) (*api.ReleaseGetResponse, error) {
	dbReleases, err := h.client.Release.Query().
		Where(release.HasCommitWith(gitcommit.Hash(*req.Commit))).
		// Get the commit, repo and project
		WithCommit(func(gcq *ent.GitCommitQuery) {
			gcq.WithRepo(func(rq *ent.RepoQuery) { rq.WithProject() })
		}).
		WithLog().
		WithViolations().
		All(h.ctx)
	if err != nil {
		return nil, HandleEntError(err, "getting releases by commit")
	}
	var resp api.ReleaseGetResponse
	for _, dbRelease := range dbReleases {
		var (
			commit  = dbRelease.Edges.Commit
			repo    = commit.Edges.Repo
			project = repo.Edges.Project
		)
		r := &api.ReleaseRead{
			Project: ent.NewProjectModelRead().FromEnt(project),
			Repo:    ent.NewRepoModelRead().FromEnt(repo),
			Commit:  ent.NewGitCommitModelRead().FromEnt(commit),
			Release: ent.NewReleaseModelRead().FromEnt(dbRelease),
		}
		if req.Policies {
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
		// // Getting the policies that apply to a release is a bit trickier, as they
		// // are actually applied to projects and repos that the release belongs to.
		// // TODO: this will produce duplicate policies if a policy is joined to both
		// // the project and repo
		// for _, dbPolicy := range repo.Edges.Policies {
		// 	r.Policies = append(r.Policies, ent.NewReleasePolicyModelRead().FromEnt(dbPolicy))
		// }
		// for _, dbPolicy := range project.Edges.Policies {
		// 	r.Policies = append(r.Policies, ent.NewReleasePolicyModelRead().FromEnt(dbPolicy))
		// }
		// Append the release violation
		for _, dbViolation := range dbRelease.Edges.Violations {
			r.Violations = append(r.Violations, ent.NewReleasePolicyViolationModelRead().FromEnt(dbViolation))
		}
		// Append the release log
		for _, dbEntry := range dbRelease.Edges.Log {
			r.Entries = append(r.Entries, ent.NewReleaseEntryModelRead().FromEnt(dbEntry))
		}
		resp.Releases = append(resp.Releases, r)
	}
	return &resp, nil
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
		SetCommit(dbCommit).
		SetName(*req.Release.Name).
		SetVersion(*req.Release.Version).
		Save(h.ctx)
	if err != nil {
		return nil, HandleEntError(err, "release")
	}

	return dbRelease, nil
}
