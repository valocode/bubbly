package release

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/valocode/bubbly/store/api"
)

type ReleaseSpec struct {
	Project      string           `hcl:"project,optional"`
	Name         string           `hcl:"name,optional"`
	Version      string           `hcl:"version,optional"`
	GitSpec      *GitSpec         `hcl:"git,block"`
	Dependencies []DependencySpec `hcl:"dependency,block"`
	// BaseDir is set at runtime so that relative paths can be resolved
	BaseDir string
}

type DependencySpec struct {
	Project string `hcl:"project,optional"`
	Name    string `hcl:"name,optional"`
	Version string `hcl:"version,optional"`
}

func CreateReleaseNode(file string) (*api.ReleaseCreateRequest, error) {
	spec, err := CreateReleaseFromSpec(file)
	if err != nil {
		return nil, err
	}
	fmt.Printf("TODODODOD: %#v\n", spec)
	return nil, nil
}

// func CreateReleaseNodeQuery(file string) (*ent.ReleaseNode, error) {
// 	spec, err := CreateReleaseFromSpec(file)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return spec.Node(true)
// }

func CreateReleaseFromSpec(file string) (*ReleaseSpec, error) {

	hclFile, diags := hclparse.NewParser().ParseHCLFile(file)
	if diags.HasErrors() {
		return nil, fmt.Errorf("error parsing bubbly file: %s: %s", file, diags.Error())
	}

	var releaseWrap struct {
		Release ReleaseSpec `hcl:"release,block"`
	}
	if diags := gohcl.DecodeBody(hclFile.Body, nil, &releaseWrap); diags.HasErrors() {
		return nil, fmt.Errorf("error decoding bubbly release from %s: %s", file, diags.Error())
	}

	release := releaseWrap.Release

	// Set the basedir of the release spec
	release.BaseDir = filepath.Dir(file)

	err := release.Validate()
	if err != nil {
		return nil, err
	}
	return &release, nil
	// return createRelease(&release)
}

// func createRelease(spec *ReleaseSpec) (*ReleaseSpec, error) {
// 	node, err := spec.Node()
// 	if err != nil {
// 		return nil, err
// 	}
// 	// TODO: post!
// 	fmt.Printf("Should send this: %#v\n", node.Graph().RootNodes)

// 	return spec, nil
// }

// func (r *ReleaseSpec) Node(query bool) (*ent.ReleaseNode, error) {
// 	err := r.Validate()
// 	if err != nil {
// 		return nil, err
// 	}
// 	//
// 	// Create the release
// 	//
// 	release := ent.NewReleaseNode().
// 		SetName(r.Name).
// 		SetVersion(r.Version).
// 		SetProject(
// 			ent.NewProjectNode().
// 				SetName(r.Project).
// 				SetOperation(ent.NodeOperationQuery),
// 		).
// 		SetOperation(ent.NodeOperationCreate)
// 	if query {
// 		release.SetOperation(ent.NodeOperationQuery)
// 	}

// 	//
// 	// Create the commit, unless we are querying the release in which case
// 	// we don't need it
// 	//
// 	if !query {
// 		commit, err := r.GitSpec.Commit(r.BaseDir)
// 		if err != nil {
// 			return nil, fmt.Errorf("error getting git data for repo %s: %w", r.GitSpec.Repo, err)
// 		}
// 		release.SetCommit(commit)
// 	}

// 	return release, nil
// }

func (r *ReleaseSpec) Validate() error {
	// Project is required
	if r.Project == "" {
		return fmt.Errorf("project is required")
	}
	// If Name was not provided, try to get a default value
	if r.Name == "" {
		if r.GitSpec == nil {
			r.GitSpec = &GitSpec{}
		}
		gitId, err := r.GitSpec.nameFromRemote(r.BaseDir)
		if err != nil {
			return fmt.Errorf("error getting git repo id: %w", err)
		}
		r.Name = gitId
	}
	// If Version was not provided, we can set a default based on git data
	if r.Version == "" {
		if r.GitSpec != nil {
			tag, commit, err := r.GitSpec.version(r.BaseDir)
			if err != nil {
				return fmt.Errorf("error getting git tag: %w", err)
			}
			r.Version = commit
			if tag != "" {
				r.Version = tag
			}
		}
	}
	return nil
}

const (
	gitRemoteGitPrefix   = "git@"
	gitRemoteHTTPSPrefix = "https://"

	gitRemoteSuffix = ".git"
)

type GitSpec struct {
	Name   string `hcl:"name,optional"`
	Repo   string `hcl:"repo,optional"`
	Remote string `hcl:"remote,optional"`
}

// func (g *GitSpec) Commit(baseDir string) (*ent.GitCommitNode, error) {
// 	gitRepo, err := g.openRepo(baseDir)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to open git repository %s: %w", g.Repo, err)
// 	}
// 	ref, err := gitRepo.Head()
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to get HEAD of repo %s: %w", g.Repo, err)
// 	}
// 	gitCommit, err := gitRepo.CommitObject(ref.Hash())
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to get commit from HEAD %s: %w", ref.Hash().String(), err)
// 	}

// 	// If the name is not given, or empty, then automatically fetch it from the
// 	// remote
// 	if g.Name == "" {
// 		var err error
// 		g.Name, err = g.nameFromRemote(baseDir)
// 		if err != nil {
// 			return nil, fmt.Errorf("error getting git repo ID from remotes: %w", err)
// 		}
// 	}

// 	commit := ent.NewGitCommitNode().
// 		SetHash(ref.Hash().String()).
// 		SetTime(gitCommit.Author.When).
// 		SetRepo(
// 			ent.NewRepoNode().SetName(g.Name),
// 		)
// 	// If HEAD is not detached, then we can add the branch name to the git data
// 	// block
// 	if ref.Name().IsBranch() {
// 		commit.SetBranch(ref.Name().Short())
// 	}

// 	tag, _, err := g.version(baseDir)
// 	if err != nil {
// 		return nil, fmt.Errorf("error getting git tag: %w", err)
// 	}
// 	if tag != "" {
// 		commit.SetTag(tag)
// 	}

// 	return commit, nil
// }

func (g *GitSpec) openRepo(baseDir string) (*git.Repository, error) {
	// Make sure there is a default Repo set
	if g.Repo == "" {
		g.Repo = "."
	}
	repoPath, err := filepath.Abs(filepath.Join(baseDir, g.Repo))
	if err != nil {
		return nil, fmt.Errorf("cannot get path to git repo %s: %w", repoPath, err)
	}
	if _, err = os.Stat(repoPath); err != nil {
		if errors.Is(os.ErrNotExist, err) {
			return nil, errors.New("git repository does not exist")
		}
		return nil, err
	}
	return git.PlainOpen(repoPath)
}

// nameFromRemote creates the git repo ID from the remote.
// Default remote is origin, unless specified otherwise, and ID should be the
// URL "normalized": no https:// or git@ prefix, no .git suffix, and only
// forward slashes (no colons)
func (g *GitSpec) nameFromRemote(baseDir string) (string, error) {
	repo, err := g.openRepo(baseDir)
	if err != nil {
		return "", fmt.Errorf("failed to open git repository %s: %w", g.Repo, err)
	}
	if g.Remote == "" {
		g.Remote = "origin"
	}
	remote, err := repo.Remote(g.Remote)
	if err != nil {
		return "", fmt.Errorf("git repo does not have a remote called %s: %w", g.Remote, err)
	}
	if len(remote.Config().URLs) == 0 {
		return "", fmt.Errorf("git repo at %s does not have any URLs for remote %s", g.Repo, g.Remote)
	}
	// TODO: maybe a warning if there are multiple URLs, but stil select the first?
	// if len(remote.Config().URLs) > 1 {
	// }
	var (
		id  string
		url = remote.Config().URLs[0]
	)
	switch {
	case strings.HasPrefix(url, gitRemoteGitPrefix):
		// Example: git@github.com:valocode/bubbly
		// Want: github.com/valocode/bubbly
		id = strings.TrimPrefix(url, gitRemoteGitPrefix)
		id = strings.ReplaceAll(id, ":", "/")
	case strings.HasPrefix(url, gitRemoteHTTPSPrefix):
		// Example: https://github.com/valocode/bubbly
		// Want: github.com/valocode/bubbly
		id = strings.TrimPrefix(url, gitRemoteHTTPSPrefix)
	default:
		return "", fmt.Errorf("supported git remote protocol for remote %s: %s", g.Remote, url)
	}

	// Remove unwanted .git suffix (if exists)
	id = strings.TrimSuffix(id, gitRemoteSuffix)

	return id, nil
}

func (g *GitSpec) version(baseDir string) (string, string, error) {
	repo, err := g.openRepo(baseDir)
	if err != nil {
		return "", "", fmt.Errorf("failed to open git repository %s: %w", g.Repo, err)
	}
	ref, err := repo.Head()
	if err != nil {
		return "", "", fmt.Errorf("failed to get HEAD of repo %s: %w", g.Repo, err)
	}
	// Add the tag, if there is one for this commit, to the commit data
	tagrefs, err := repo.Tags()
	if err != nil {
		return "", "", fmt.Errorf("failed to read tags from repo %s, error %w", g.Repo, err)
	}
	var (
		tag    string
		commit = ref.Hash().String()
	)
	// Ignore the returned error, as it shouldn't be triggered
	err = tagrefs.ForEach(func(t *plumbing.Reference) error {
		// Noteworthy: annotated git tags are actually individual objects in git
		// with their own Hash. Thus, we cannot just compare the tag hash with
		// the commit hash. For lightweight tags, then yes we could simply do that.
		// Solution: get the tag object for the tag hash, from which we can get
		// the commit object, which we can compare to the commit we have from HEAD
		tagObj, err := repo.TagObject(t.Hash())
		if err != nil {
			return fmt.Errorf("error getting tag object from list of tag refs: %w", err)
		}
		commitObj, err := tagObj.Commit()
		if err != nil {
			return fmt.Errorf("error getting commit from tag: %w", err)
		}

		// If the tag ref is for the same commit, we have a match and so add the
		// tag to the commit data
		if ref.Hash() == commitObj.Hash {
			tag = t.Name().Short()
		}
		return nil
	})
	if err != nil {
		return "", "", err
	}
	return tag, commit, nil
}
