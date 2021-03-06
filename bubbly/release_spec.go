package bubbly

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	dockerclient "github.com/docker/docker/client"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/valocode/bubbly/api"
	"github.com/valocode/bubbly/api/common"
	"github.com/valocode/bubbly/api/core"
	"github.com/valocode/bubbly/bubbly/builtin"
	"github.com/valocode/bubbly/env"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/gocty"
)

type ReleaseSpec struct {
	Name         string         `hcl:"name,optional"`
	Version      string         `hcl:"version,optional"`
	Project      string         `hcl:"project,optional"`
	GitItem      *gitItem       `hcl:"git,block"`
	ArtifactItem *artifactItem  `hcl:"artifact,block"`
	Stages       []releaseStage `hcl:"stage,block"`
	// BaseDir is set at runtime so that relative paths can be resolved
	BaseDir string
}

const (
	releaseItemGitType      = "git"
	releaseItemArtifactType = "artifact"
	releaseItemReleaseType  = "release"
)

func (r *ReleaseSpec) Data() (core.DataBlocks, error) {
	err := r.Validate()
	if err != nil {
		return nil, err
	}

	//
	// Create the release and release items data blocks
	//
	release := builtin.Release{
		Name:    r.Name,
		Version: r.Version,
		Project: &builtin.Project{
			Name: r.Project,
		},
		// Do not allow update, as that causes all kinds of complications with
		// releases. If you want to change a specific release, first delete it
		// and recreate it
		DBlock_Policy: core.CreatePolicy,
	}

	//
	// Create the release_stage and release_criteria data blocks
	//
	for _, stage := range r.Stages {
		relStage := builtin.ReleaseStage{
			Name: stage.Name,
		}
		for _, criteria := range stage.Criterion {
			relCriteria := builtin.ReleaseCriteria{
				EntryName:    criteria.Name,
				DBlock_Joins: []string{"release"},
			}
			relStage.ReleaseCriteria = append(relStage.ReleaseCriteria, relCriteria)
		}
		release.ReleaseStage = append(release.ReleaseStage, relStage)
	}

	//
	// Create the release item and criteria data blocks
	//
	// For each release item, create the item type and the release item.
	// E.g. for git we create the commit and also the release_item that joins
	// to that specific commit.
	if r.GitItem != nil {
		commit, err := r.GitItem.Commit(r.BaseDir)
		if err != nil {
			return nil, fmt.Errorf("error getting git data for repo %s: %w", r.GitItem.Repo, err)
		}
		relItem := builtin.ReleaseItem{
			Type:   releaseItemGitType,
			Commit: commit,
		}
		release.ReleaseItem = append(release.ReleaseItem, relItem)
	}
	if r.ArtifactItem != nil {
		// TODO: Remove ReleaseItems...
	}

	return builtin.ToDataBlocks(release), nil
}

func (r *ReleaseSpec) DataRef() (core.DataBlocks, error) {
	err := r.Validate()
	if err != nil {
		return nil, err
	}

	//
	// Create the release and release items data blocks
	//
	release := builtin.Release{
		Name:    r.Name,
		Version: r.Version,
		Project: &builtin.Project{
			Name:          r.Project,
			DBlock_Policy: core.ReferencePolicy,
		},
		DBlock_Policy: core.ReferencePolicy,
	}
	return builtin.ToDataBlocks(release), nil
}

func (r *ReleaseSpec) String() string {
	var rType string
	switch {
	case r.ArtifactItem != nil:
		rType = "artifact"
	case r.GitItem != nil:
		rType = "git"
	default:
		rType = "unknown"
	}
	return "Project: " + r.Project + "\nName: " + r.Name + "\nVersion: " + r.Version +
		"\nType: " + rType
}

func (r *ReleaseSpec) Validate() error {
	// Project is required
	if r.Project == "" {
		return fmt.Errorf("project is required")
	}
	// If Name was not provided, try to get a default value
	if r.Name == "" {
		if r.GitItem != nil {
			gitId, err := r.GitItem.idFromRemote(r.BaseDir)
			if err != nil {
				return fmt.Errorf("error getting git repo id: %w", err)
			}
			r.Name = gitId
		}
	}
	// If Version was not provided, we can set a default based on git data
	if r.Version == "" {
		if r.GitItem != nil {
			tag, commit, err := r.GitItem.version(r.BaseDir)
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

type gitItem struct {
	ID     string `hcl:"id,optional"`
	Repo   string `hcl:"repo,optional"`
	Remote string `hcl:"remote,optional"`
}

func (g *gitItem) Commit(baseDir string) (*builtin.Commit, error) {
	gitRepo, err := g.openRepo(baseDir)
	if err != nil {
		return nil, fmt.Errorf("failed to open git repository %s: %w", g.Repo, err)
	}
	ref, err := gitRepo.Head()
	if err != nil {
		return nil, fmt.Errorf("failed to get HEAD of repo %s: %w", g.Repo, err)
	}
	gitCommit, err := gitRepo.CommitObject(ref.Hash())
	if err != nil {
		return nil, fmt.Errorf("failed to get commit from HEAD %s: %w", ref.Hash().String(), err)
	}

	// If the ID is not given, or empty, then automatically fetch it from the
	// remote
	if g.ID == "" {
		var err error
		g.ID, err = g.idFromRemote(baseDir)
		if err != nil {
			return nil, fmt.Errorf("error getting git repo ID from remotes: %w", err)
		}
	}

	commit := builtin.Commit{
		Id:   ref.Hash().String(),
		Time: gitCommit.Author.When.String(),
		Repo: &builtin.Repo{
			Name: g.ID,
		},
	}
	// If HEAD is not detached, then we can add the branch name to the git data
	// block
	if ref.Name().IsBranch() {
		branch := builtin.Branch{Name: ref.Name().Short()}
		commit.Repo.Branch = append(commit.Repo.Branch, branch)
		commit.Branch = &branch
	}

	tag, _, err := g.version(baseDir)
	if err != nil {
		return nil, fmt.Errorf("error getting git tag: %w", err)
	}
	if tag != "" {
		commit.Tag = tag
	}

	return &commit, nil
}

func (g *gitItem) openRepo(baseDir string) (*git.Repository, error) {
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

// idFromRemote creates the git repo ID from the remote.
// Default remote is origin, unless specified otherwise, and ID should be the
// URL "normalized": no https:// or git@ prefix, no .git suffix, and only
// forward slashes (no colons)
func (g *gitItem) idFromRemote(baseDir string) (string, error) {
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

func (g *gitItem) version(baseDir string) (string, string, error) {
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

const (
	artifactTypeDiv      = "://"
	artifactFilePrefix   = "file" + artifactTypeDiv
	artifactDockerPrefix = "docker" + artifactTypeDiv
)

type artifactItem struct {
	Name     string `hcl:"name,attr"`
	Location string `hcl:"location,attr"`
}

func (a *artifactItem) Artifact(baseDir string) (*builtin.Artifact, error) {
	var (
		sha256 string
		err    error
	)
	switch {
	case strings.HasPrefix(a.Location, artifactFilePrefix):
		sha256, err = a.sha256SumFile(baseDir)
	case strings.HasPrefix(a.Location, artifactDockerPrefix):
		sha256, err = a.sha256SumDocker()
	default:
		// Check if there was supposed to be a type given before the ://
		// If so, it's an unknown type and we should error
		if strings.Contains(a.Location, artifactTypeDiv) {
			// Get the substring up to and including ://
			typeStr := a.Location[0 : strings.Index(a.Location, artifactTypeDiv)+3]
			return nil, fmt.Errorf("unkown artifact type: %s", typeStr)
		}
		// Treat as a file by default
		sha256, err = a.sha256SumFile(baseDir)
	}
	if err != nil {
		return nil, fmt.Errorf("error calculating sha256 sum of %s: %w", a.Location, err)
	}
	return &builtin.Artifact{
		Name:     a.Name,
		Sha256:   sha256,
		Location: a.Location,
	}, nil
}

func (a *artifactItem) sha256SumFile(baseDir string) (string, error) {
	loc := strings.TrimPrefix(a.Location, artifactFilePrefix)
	loc = filepath.Join(baseDir, loc)
	f, err := os.Open(loc)
	if err != nil {
		return "", fmt.Errorf("error opening artifact file %s: %w", loc, err)
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", fmt.Errorf("error reading artifact file %s: %w", loc, err)
	}
	// Sum comes back as hex, so need to encode that
	return hex.EncodeToString(h.Sum(nil)), nil
}

func (a *artifactItem) sha256SumDocker() (string, error) {
	loc := strings.TrimPrefix(a.Location, artifactDockerPrefix)
	cli, err := dockerclient.NewClientWithOpts(dockerclient.FromEnv)
	if err != nil {
		return "", fmt.Errorf("error creating docker client: %w", err)
	}
	t, _, err := cli.ImageInspectWithRaw(context.TODO(), loc)
	if err != nil {
		return "", fmt.Errorf("error inspecting image %s: %w", loc, err)
	}

	return strings.TrimPrefix(t.ID, "sha256:"), nil
}

type releaseStage struct {
	Name      string            `hcl:",label"`
	Criterion []releaseCriteria `hcl:"criteria,block"`
}

type releaseCriteria struct {
	Name     string        `hcl:",label"`
	Artifact *artifactItem `hcl:"artifact,block"`
	Run      []resourceRun `hcl:"run,block"`
}

// Evaluate evaluates a release criteria and produces data blocks to log a release entry
func (c *releaseCriteria) Evaluate(bCtx *env.BubblyContext, release *ReleaseSpec, dataCtx core.DataBlocks, baseDir string) (core.DataBlocks, error) {
	// Validation of the release criteria
	if c.Artifact != nil && len(c.Run) > 0 {
		return nil, errors.New("release criteria cannot specify both artifact and resource runs")
	}
	if c.Artifact == nil && len(c.Run) == 0 {
		return nil, errors.New("release criteria has no criteria")
	}

	var (
		data        core.DataBlocks
		entryReason string
		entryResult = true
	)

	if c.Artifact != nil {
		// Get the data block for the artifact
		// TODO: handle if the artifact doesn't exist. We do not want to error,
		// but it should mark the release_entry as false, e.g.
		// entry.Result = <ARTIFACT_RESULT>
		artifact, err := c.Artifact.Artifact(baseDir)
		if err != nil {
			return nil, fmt.Errorf("error with artifact %s: %w", c.Artifact.Name, err)
		}
		data = append(data, builtin.ToDataBlocks(*artifact)...)
	}
	// Iterate through the run blocks. As soon as one fails, or one is a criteria
	// kind with a failing result, break the loop
	for _, run := range c.Run {
		// Add the release name and version as inputs to the run, in case the
		// resources use them as input
		run.Inputs = append(run.Inputs, &core.InputDefinition{
			Name: "release",
			Value: cty.MapVal(map[string]cty.Value{
				"name":    cty.StringVal(release.Name),
				"version": cty.StringVal(release.Version),
			}),
		})
		resource, output := run.Run(bCtx, dataCtx)
		if output.Error != nil {
			return nil, fmt.Errorf("resource run failed for %s: %w", run.Resource, output.Error)
		}
		// If the resource was of criteria kind, then we care about the resource
		// output value
		if resource.Kind() == core.CriteriaResourceKind {
			var result core.CriteriaResult
			if err := gocty.FromCtyValue(output.Value, &result); err != nil {
				return nil, fmt.Errorf("error getting criteria result from resource output: %w", err)
			}
			entryResult = result.Result
			entryReason = result.Reason
		}
	}
	// Create the data block for the release entry
	relCriteria := builtin.ReleaseCriteria{
		EntryName:     c.Name,
		DBlock_Joins:  []string{"release"},
		DBlock_Policy: core.ReferencePolicy,
		ReleaseEntry: []builtin.ReleaseEntry{
			{
				Name:   c.Name,
				Result: entryResult,
				Reason: entryReason,
			},
		},
	}

	data = append(data, builtin.ToDataBlocks(relCriteria)...)
	return data, nil
}

type resourceRun struct {
	Resource string                `hcl:",label"`
	Inputs   core.InputDefinitions `hcl:"input,block"`
}

func (r *resourceRun) Run(bCtx *env.BubblyContext, dataCtx core.DataBlocks) (core.Resource, core.ResourceOutput) {
	ctx := core.NewResourceContext(cty.NilVal, api.NewResource, nil)
	// Add the data block containing the release into the context
	ctx.DataBlocks = dataCtx
	return common.RunResourceByID(bCtx, ctx, r.Resource, r.Inputs.Value())
}
