package release

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/valocode/bubbly/config"
	"github.com/valocode/bubbly/ent"
	"github.com/valocode/bubbly/store/api"
)

type ReleaseOptions struct {
	Filename string
	RepoDir  string
}

type ReleaseSpec struct {
	Project string `json:"project,omitempty" hcl:"project,optional" validate:"required"`
	// Repo is the name of the repository. If not provided, the name is created
	// from the git remote
	Repo    string  `json:"repo,omitempty" hcl:"repo,optional"`
	Name    string  `json:"name,omitempty" hcl:"name,optional"`
	Version string  `json:"version,omitempty" hcl:"version,optional"`
	GitSpec GitSpec `json:"git,omitempty" hcl:"git,block"`
}

type ReleaseSpecWrap struct {
	Release ReleaseSpec `json:"release" hcl:"release,block"`
}

func DefaultReleaseSpec() *ReleaseSpec {
	return &ReleaseSpec{
		Project: config.DefaultEnvStr("BUBBLY_PROJECT", "default"),
	}
}

func CreateRelease(opts ReleaseOptions) (*api.ReleaseCreateRequest, error) {
	var specFile = opts.Filename
	if opts.Filename == "" {
		file, err := defaultSpecFile()
		if err != nil {
			return nil, err
		}
		specFile = file
	}

	release, err := decodeReleaseSpec(specFile)
	if err != nil {
		return nil, fmt.Errorf("decoding release spec %s: %w", specFile, err)
	}
	// Set the path to the git repository relative to the release spec
	baseDir := filepath.Dir(opts.Filename)
	release.GitSpec.Path = filepath.Join(baseDir, release.GitSpec.Path)
	if err := release.Validate(); err != nil {
		return nil, err
	}

	commit, err := release.commit()
	if err != nil {
		return nil, err
	}

	return &api.ReleaseCreateRequest{
		Repo:    ent.NewRepoModelCreate().SetName(release.Repo),
		Release: ent.NewReleaseModelCreate().SetName(release.Name).SetVersion(release.Version),
		Commit:  commit,
	}, nil
}

func (r *ReleaseSpec) commit() (*ent.GitCommitModelCreate, error) {
	gs := r.GitSpec
	gitRepo, err := gs.openRepo(gs.Path)
	if err != nil {
		return nil, fmt.Errorf("failed to open git repository %s: %w", gs.Path, err)
	}
	ref, err := gitRepo.Head()
	if err != nil {
		return nil, fmt.Errorf("failed to get HEAD of repo %s: %w", gs.Path, err)
	}
	gitCommit, err := gitRepo.CommitObject(ref.Hash())
	if err != nil {
		return nil, fmt.Errorf("failed to get commit from HEAD %s: %w", ref.Hash().String(), err)
	}

	// Set the default repo name
	if r.Repo == "" {
		remoteName, err := gs.nameFromRemote()
		if err != nil {
			return nil, fmt.Errorf("error getting git repo name from remotes: %w", err)
		}
		r.Repo = remoteName
	}

	commit := ent.NewGitCommitModelCreate().
		SetHash(ref.Hash().String()).
		SetTime(gitCommit.Author.When)
	// If HEAD is not detached, then we can add the branch name to the git data
	// block
	if ref.Name().IsBranch() {
		commit.SetBranch(ref.Name().Short())
	}

	tag, _, err := gs.version()
	if err != nil {
		return nil, fmt.Errorf("error getting git tag: %w", err)
	}
	if tag != "" {
		commit.SetTag(tag)
	}

	return commit, nil
}

func decodeReleaseSpec(filename string) (*ReleaseSpec, error) {
	relWrap := ReleaseSpecWrap{
		Release: *DefaultReleaseSpec(),
	}
	// If no spec file was found and no spec file was provided as an argument,
	// return the default config
	if filename == "" {
		return &relWrap.Release, nil
	}
	switch {
	case strings.HasSuffix(filename, ".json"):
		file, err := os.Open(filename)
		if err != nil {
			return nil, err
		}
		if err := json.NewDecoder(file).Decode(&relWrap); err != nil {
			return nil, err
		}
	case strings.HasSuffix(filename, ".hcl"):
		hclFile, diags := hclparse.NewParser().ParseHCLFile(filename)
		if diags.HasErrors() {
			return nil, fmt.Errorf("parsing bubbly file: %s: %s", filename, diags.Error())
		}

		if diags := gohcl.DecodeBody(hclFile.Body, nil, &relWrap); diags.HasErrors() {
			return nil, fmt.Errorf("decoding bubbly release from %s: %s", filename, diags.Error())
		}
	default:
		return nil, fmt.Errorf("unknown bubbly release specification file type: %s", filename)
	}
	return &relWrap.Release, nil
}

func defaultSpecFile() (string, error) {
	var exts = []string{
		".json", ".hcl",
	}
	for _, ext := range exts {
		f, err := os.Open(".bubbly" + ext)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return "", err
		}
		defer f.Close()
		return ".bubbly" + ext, nil
	}
	return "", nil
}

func (r *ReleaseSpec) Validate() error {
	// Project is required
	if r.Project == "" {
		return fmt.Errorf("project is required")
	}
	// If Name was not provided, try to get a default value
	if r.Name == "" {
		gitId, err := r.GitSpec.nameFromRemote()
		if err != nil {
			return fmt.Errorf("error getting git repo id: %w", err)
		}
		r.Name = gitId
	}
	// If Version was not provided, we can set a default based on git data
	if r.Version == "" {
		tag, commit, err := r.GitSpec.version()
		if err != nil {
			return fmt.Errorf("error checking for git tag: %w", err)
		}
		r.Version = commit
		if tag != "" {
			r.Version = tag
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
	// Path points to the git repository, relative from the release spec file
	Path   string `json:"path,omitempty" hcl:"path,optional"`
	Remote string `json:"remote,omitempty" hcl:"remote,optional"`
}

func (g *GitSpec) openRepo(baseDir string) (*git.Repository, error) {
	if _, err := os.Stat(g.Path); err != nil {
		if errors.Is(os.ErrNotExist, err) {
			return nil, fmt.Errorf("git repository does not exist: %s", g.Path)
		}
		return nil, err
	}
	return git.PlainOpen(g.Path)
}

// nameFromRemote creates the git repo ID from the remote.
// Default remote is origin, unless specified otherwise, and ID should be the
// URL "normalized": no https:// or git@ prefix, no .git suffix, and only
// forward slashes (no colons)
func (g *GitSpec) nameFromRemote() (string, error) {
	repo, err := g.openRepo(g.Path)
	if err != nil {
		return "", fmt.Errorf("failed to open git repository %s: %w", g.Path, err)
	}
	if g.Remote == "" {
		g.Remote = "origin"
	}
	remote, err := repo.Remote(g.Remote)
	if err != nil {
		return "", fmt.Errorf("git repo does not have a remote called %s: %w", g.Remote, err)
	}
	if len(remote.Config().URLs) == 0 {
		return "", fmt.Errorf("git repo at %s does not have any URLs for remote %s", g.Path, g.Remote)
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

// version returns <tag>, <commit>, <error>
func (g *GitSpec) version() (string, string, error) {
	repo, err := g.openRepo(g.Path)
	if err != nil {
		return "", "", fmt.Errorf("failed to open git repository %s: %w", g.Path, err)
	}
	ref, err := repo.Head()
	if err != nil {
		return "", "", fmt.Errorf("failed to get HEAD of repo %s: %w", g.Path, err)
	}
	// Add the tag, if there is one for this commit, to the commit data
	tagrefs, err := repo.Tags()
	if err != nil {
		return "", "", fmt.Errorf("failed to read tags from repo %s, error %w", g.Path, err)
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
