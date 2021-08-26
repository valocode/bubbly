package store

import (
	"github.com/valocode/bubbly/ent"
	"github.com/valocode/bubbly/ent/gitcommit"
	"github.com/valocode/bubbly/ent/project"
	"github.com/valocode/bubbly/ent/release"
	"github.com/valocode/bubbly/ent/repo"
	"github.com/valocode/bubbly/store/api"
)

func (s *Store) CreateRelease(req *api.ReleaseCreateRequest) (*ent.Release, error) {
	return s.createRelease(req)
}

func (s *Store) LogArtifact(req *api.ArtifactLogRequest) (*ent.Artifact, error) {
	return s.logArtifact(req)
}

func (s *Store) createRelease(req *api.ReleaseCreateRequest) (*ent.Release, error) {
	if err := s.validator.Struct(req); err != nil {
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
	dbRelease, err := s.client.Release.Query().Where(
		release.Name(*req.Release.Name), release.Version(*req.Release.Version),
		release.HasCommitWith(
			gitcommit.Hash(*req.Commit.Hash),
		),
	).Only(s.ctx)
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
	dbProject, err := s.client.Project.Query().Where(
		project.Name(*req.Project.Name),
	).Only(s.ctx)
	if err != nil {
		if !ent.IsNotFound(err) {
			return nil, HandleEntError(err, "repo")
		}
		dbProject, err = s.client.Project.Create().
			SetModelCreate(req.Project).
			Save(s.ctx)
		if err != nil {
			return nil, HandleEntError(err, "repo")
		}
	}

	dbRepo, err := s.client.Repo.Query().Where(
		repo.Name(*req.Repo.Name),
	).Only(s.ctx)
	if err != nil {
		if !ent.IsNotFound(err) {
			return nil, HandleEntError(err, "repo")
		}
		dbRepo, err = s.client.Repo.Create().
			SetModelCreate(req.Repo).
			SetProject(dbProject).
			Save(s.ctx)
		if err != nil {
			return nil, HandleEntError(err, "repo")
		}
	}

	dbCommit, err := s.client.GitCommit.Query().Where(
		gitcommit.Hash(*req.Commit.Hash),
		gitcommit.HasRepoWith(repo.ID(dbRepo.ID)),
	).Only(s.ctx)
	if err != nil {
		if !ent.IsNotFound(err) {
			return nil, HandleEntError(err, "commit")
		}
		if req.Commit.Tag != nil {
			commitTag = *req.Commit.Tag
		}
		dbCommit, err = s.client.GitCommit.Create().
			SetHash(*req.Commit.Hash).
			SetBranch(*req.Commit.Branch).
			SetTime(*req.Commit.Time).
			SetTag(commitTag).
			SetRepo(dbRepo).
			Save(s.ctx)
		if err != nil {
			return nil, HandleEntError(err, "commit")
		}
	}

	dbRelease, err = s.client.Release.Create().
		SetCommit(dbCommit).
		SetName(*req.Release.Name).
		SetVersion(*req.Release.Version).
		Save(s.ctx)
	if err != nil {
		return nil, HandleEntError(err, "release")
	}

	return dbRelease, nil
}

func (s *Store) logArtifact(req *api.ArtifactLogRequest) (*ent.Artifact, error) {
	if err := s.validator.Struct(req); err != nil {
		return nil, HandleValidatorError(err, "artifact")
	}
	dbRelease, err := s.releaseFromCommit(req.Commit)
	if err != nil {
		return nil, err
	}
	dbArtifact, err := s.client.Artifact.Create().
		SetModelCreate(req.Artifact).
		SetRelease(dbRelease).
		Save(s.ctx)
	if err != nil {
		return nil, HandleEntError(err, "artifact")
	}
	return dbArtifact, nil
}
