package store

import (
	"github.com/valocode/bubbly/ent"
	"github.com/valocode/bubbly/ent/gitcommit"
	"github.com/valocode/bubbly/ent/project"
	"github.com/valocode/bubbly/ent/release"
	"github.com/valocode/bubbly/ent/repo"
)

func (h *Handler) projectIDs(client *ent.Client, projects []string) ([]int, error) {
	var projectIDs []int
	if len(projects) == 0 {
		return projectIDs, nil
	}
	// Check that the projects exist
	for _, p := range projects {
		ok, err := client.Project.Query().Where(project.Name(p)).Exist(h.ctx)
		if err != nil {
			return nil, HandleEntError(err, "checking projects exist")
		}
		if !ok {
			return nil, NewNotFoundError(nil, "project %s does not exist", p)
		}
	}
	projectIDs, err := client.Project.Query().Where(
		project.NameIn(projects...),
	).IDs(h.ctx)
	if err != nil {
		return nil, HandleEntError(err, "getting project IDs")
	}
	return projectIDs, nil
}

func (h *Handler) repoIDs(client *ent.Client, repos []string) ([]int, error) {
	var repoIDs []int
	if len(repos) == 0 {
		return repoIDs, nil
	}
	// Check that the repos exist
	for _, r := range repos {
		ok, err := client.Repo.Query().Where(repo.Name(r)).Exist(h.ctx)
		if err != nil {
			return nil, HandleEntError(err, "checking repo exist")
		}
		if !ok {
			return nil, NewNotFoundError(nil, "repo %s does not exist", r)
		}
	}
	repoIDs, err := client.Repo.Query().Where(
		repo.NameIn(repos...),
	).IDs(h.ctx)
	if err != nil {
		return nil, HandleEntError(err, "getting repo IDs")
	}
	return repoIDs, nil
}

func (h *Handler) releaseIDs(client *ent.Client, commits []string) ([]int, error) {
	var releaseIDs []int
	if len(commits) == 0 {
		return releaseIDs, nil
	}
	// Check that the repos exist
	for _, c := range commits {
		ok, err := client.Release.Query().Where(release.HasCommitWith(gitcommit.Hash(c))).Exist(h.ctx)
		if err != nil {
			return nil, HandleEntError(err, "checking release exist")
		}
		if !ok {
			return nil, NewNotFoundError(nil, "release with commit %s does not exist", c)
		}
	}
	releaseIDs, err := client.Release.Query().Where(
		release.HasCommitWith(gitcommit.HashIn(commits...)),
	).IDs(h.ctx)
	if err != nil {
		return nil, HandleEntError(err, "getting repo IDs")
	}
	return releaseIDs, nil
}
