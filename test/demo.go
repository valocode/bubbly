package test

import (
	"crypto/sha1"
	"crypto/sha256"
	"fmt"
	"math/rand"
	"time"

	"github.com/valocode/bubbly/ent"
	"github.com/valocode/bubbly/ent/artifact"
	"github.com/valocode/bubbly/store/api"
)

type DemoRepoOptions struct {
	Name    string
	Project string
	Branch  string

	NumCommits int
}

type RepoData struct {
	Releases []ReleaseData
}

type ReleaseData struct {
	Release    *api.ReleaseCreateRequest
	Artifacts  []*api.ArtifactCreateRequest
	Components []*api.ComponentCreateRequest
}

var demoData = []DemoRepoOptions{
	{
		Name:       "demo",
		Project:    "demo",
		Branch:     "main",
		NumCommits: 30,
	},
}

func CreateDummyData() []RepoData {
	var repos []RepoData
	for _, opt := range demoData {
		repos = append(repos, RepoData{
			Releases: createRepoData(opt),
		})
	}

	return repos
}

func createRepoData(opt DemoRepoOptions) []ReleaseData {
	var releases []ReleaseData
	cTime := time.Now()
	for i := 1; i <= opt.NumCommits; i++ {
		var (
			data    ReleaseData
			tag     string
			version string
		)
		// Calculate a random time diff for the commits
		cTime = cTime.Add(time.Hour * -time.Duration(rand.Intn(100)))

		// Calculate a random commit hash
		hash := fmt.Sprintf("%x",
			sha1.Sum([]byte(fmt.Sprintf("%d", i))),
		)
		if i%20 == 0 {
			tag = fmt.Sprintf("%d.%d.%d", i/200, i/50, i/20)
			version = tag
		}
		if tag == "" {
			version = hash
		}

		fmt.Printf("Creating release with version %s at %s\n", version, cTime.Format(time.RFC3339))
		relReq := api.ReleaseCreateRequest{
			Project:    &opt.Project,
			Repository: ent.NewRepositoryModelCreate().SetName(opt.Name),
			Release:    ent.NewReleaseModelCreate().SetName(opt.Name).SetVersion(version),
			Commit: ent.NewGitCommitModelCreate().
				SetHash(hash).
				SetBranch(opt.Branch).
				SetTag(tag).
				SetTime(cTime),
		}
		data.Release = &relReq
		//
		// Artifact
		//
		{
			artReq := api.ArtifactCreateRequest{
				Artifact: ent.NewArtifactModelCreate().
					SetName("dummy").
					SetSha256(fmt.Sprintf("%x", sha256.Sum256([]byte(hash)))).
					SetType(artifact.TypeDocker),
				Commit: &hash,
			}

			data.Artifacts = append(data.Artifacts, &artReq)
		}

		//
		// Components
		//
		{
			//
			// TODO:: set components with "real" data...
			//
			// var components []*ent.Component
			// for _, c := range components {
			// &api.ComponentCreate{
			// 	ComponentModelCreate: *ent.NewComponentModelCreate().
			// 		SetName("comp-1").SetScheme("pkg").SetVersion("123"),
			// 	Vulnerabilities: []*ent.VulnerabilityModelCreate{
			// 		ent.NewVulnerabilityModelCreate().SetVid("demo-123").SetLabels(schema.Labels{
			// 			"bubbly/os": "windows",
			// 		}),
			// 		ent.NewVulnerabilityModelCreate().SetVid("demo-456").SetLabels(schema.Labels{
			// 			"bubbly/os": "windows",
			// 		}),
			// 	},
			// }
			data.Components = append(data.Components, &api.ComponentCreateRequest{
				Component: &api.ComponentCreate{
					ComponentModelCreate: ent.NewComponentModelCreate().SetName("walal"),
				},
			})
		}

		// Prepend the release so that the dates of the releases go in ascending order
		releases = append([]ReleaseData{data}, releases...)
	}

	return releases
}
