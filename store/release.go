package store

import (
	"errors"

	"github.com/hashicorp/go-multierror"
	"github.com/valocode/bubbly/ent"
	"github.com/valocode/bubbly/ent/gitcommit"
	"github.com/valocode/bubbly/ent/release"
	"github.com/valocode/bubbly/ent/repo"
	"github.com/valocode/bubbly/store/api"
)

func (s *Store) CreateRelease(rel *api.ReleaseCreateRequest) (*ent.Release, error) {
	return s.createRelease(rel)
}

func (s *Store) LogArtifact(req *api.ArtifactLogRequest) (*ent.Artifact, error) {
	return s.logArtifact(req)
}

func (s *Store) createRelease(rel *api.ReleaseCreateRequest) (*ent.Release, error) {

	var vErr multierror.Error
	if rel == nil {
		return nil, NewValidationError(nil, "request is nil")
	}
	if rel.Repo == nil {
		vErr = *multierror.Append(&vErr, errors.New("repo is required"))
	}
	if rel.Commit == nil {
		vErr = *multierror.Append(&vErr, errors.New("commit is required"))
	}
	if err := vErr.ErrorOrNil(); err != nil {
		return nil, HandleMultiVError(vErr)
	}
	var (
		repoName     = rel.Repo.GetNameOrErr(&vErr)
		commitHash   = rel.Commit.GetHashOrErr(&vErr)
		commitBranch = rel.Commit.GetBranchOrErr(&vErr)
		commitTime   = rel.Commit.GetTimeOrErr(&vErr)
		commitTag    string
		relName      string
		relVersion   string
	)
	// Validate that all the required values exist
	if err := vErr.ErrorOrNil(); err != nil {
		return nil, HandleMultiVError(vErr)
	}

	// Set defaults / optional values
	if rel.Commit.Tag != nil {
		commitTag = *rel.Commit.Tag
	}
	if rel.Release != nil {
		if rel.Release.Name != nil {
			relName = *rel.Release.Name
		}
		if rel.Release.Version != nil {
			relVersion = *rel.Release.Version
		}
	}
	if relName == "" {
		relName = *repoName
	}
	if relVersion == "" {
		relVersion = *commitHash
		if commitTag != "" {
			relVersion = commitTag
		}
	}

	// Check if the release exists already
	dbRelease, err := s.client.Release.Query().Where(
		release.Name(relName), release.Version(relVersion),
		release.HasCommitWith(
			gitcommit.Hash(*commitHash),
			gitcommit.HasRepoWith(repo.Name(*repoName)),
		),
	).Only(s.ctx)
	if err != nil {
		if !ent.IsNotFound(err) {
			return nil, HandleEntError(err, "release")
		}
		// Otherwise continue and create the release
	}
	if dbRelease != nil {
		return dbRelease, nil
	}

	// First we need to get the project, repo and commit
	// dbProject, err := s.client.Project.Query().Where(
	// 	project.Name(*projectName),
	// ).Only(s.ctx)
	// if err != nil {
	// 	if !ent.IsNotFound(err) {
	// 		return nil, HandleEntError(err, "project")
	// 	}
	// 	dbProject, err = s.client.Project.Create().SetName(*projectName).Save(s.ctx)
	// 	if err != nil {
	// 		return nil, HandleEntError(err, "project")
	// 	}
	// }

	dbRepo, err := s.client.Repo.Query().Where(
		repo.Name(*repoName),
	).Only(s.ctx)
	if err != nil {
		if !ent.IsNotFound(err) {
			return nil, HandleEntError(err, "repo")
		}
		dbRepo, err = s.client.Repo.Create().
			SetName(*repoName).
			Save(s.ctx)
		if err != nil {
			return nil, HandleEntError(err, "repo")
		}
	}

	dbCommit, err := s.client.GitCommit.Query().Where(
		gitcommit.Hash(*commitHash),
		gitcommit.HasRepoWith(repo.ID(dbRepo.ID)),
	).Only(s.ctx)
	if err != nil {
		if !ent.IsNotFound(err) {
			return nil, HandleEntError(err, "commit")
		}
		if rel.Commit.Tag != nil {
			commitTag = *rel.Commit.Tag
		}
		dbCommit, err = s.client.GitCommit.Create().
			SetHash(*commitHash).
			SetBranch(*commitBranch).
			SetTime(*commitTime).
			SetTag(commitTag).
			SetRepo(dbRepo).
			Save(s.ctx)
		if err != nil {
			return nil, HandleEntError(err, "commit")
		}
	}

	dbRelease, err = s.client.Release.Create().
		SetCommit(dbCommit).
		SetName(relName).
		SetVersion(relVersion).
		Save(s.ctx)
	if err != nil {
		return nil, HandleEntError(err, "release")
	}

	return dbRelease, nil
}

func (s *Store) logArtifact(req *api.ArtifactLogRequest) (*ent.Artifact, error) {
	var vErr error
	if req.Commit == nil {
		vErr = multierror.Append(vErr, NewValidationError(nil, "commit is required"))
	}
	if req.Artifact == nil {
		vErr = multierror.Append(vErr, NewValidationError(nil, "release is required"))
	}
	if vErr != nil {
		return nil, vErr
	}
	var (
		artName   = req.Artifact.GetNameOrErr(vErr)
		artSha256 = req.Artifact.GetSha256OrErr(vErr)
		artType   = req.Artifact.GetTypeOrErr(vErr)
	)
	if vErr != nil {
		return nil, vErr
	}
	rel, err := s.releaseFromCommit(req.Commit)
	if err != nil {
		return nil, err
	}
	dbArtifact, err := s.client.Artifact.Create().
		SetName(*artName).
		SetSha256(*artSha256).
		SetType(*artType).
		SetRelease(rel).
		Save(s.ctx)
	if err != nil {
		return nil, HandleEntError(err, "artifact")
	}
	return dbArtifact, nil
}
