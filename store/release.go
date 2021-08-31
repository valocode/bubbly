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

func (h *Handler) LogArtifact(req *api.ArtifactLogRequest) (*ent.Artifact, error) {
	return h.logArtifact(req)
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

func (h *Handler) logArtifact(req *api.ArtifactLogRequest) (*ent.Artifact, error) {
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
	return dbArtifact, nil
}
