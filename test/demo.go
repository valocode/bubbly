package test

import (
	"crypto/sha1"
	"crypto/sha256"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/valocode/bubbly/ent"
	"github.com/valocode/bubbly/ent/artifact"
	"github.com/valocode/bubbly/ent/codeissue"
	schema "github.com/valocode/bubbly/ent/schema/types"
	"github.com/valocode/bubbly/store/api"
)

type DemoRepoOptions struct {
	Name    string
	Project string
	Branch  string

	NumCommits    int
	NumTests      int
	NumCodeIssues int

	ArtifactTimeMin int
	ArtifactTimeMax int

	TestRunTimeMin int
	TestRunTimeMax int

	IssueScanTimeMin int
	IssueScanTimeMax int

	CveScanTimeMin int
	CveScanTimeMax int
}

type RepoData struct {
	Releases []ReleaseData
}

type ReleaseData struct {
	Release   *api.ReleaseCreateRequest
	Artifacts []*api.ArtifactLogRequest
	CodeScans []*api.CodeScanRequest
	TestRuns  []*api.TestRunRequest
}

var demoData = []DemoRepoOptions{
	{
		Name:          "demo",
		Project:       "demo",
		Branch:        "main",
		NumCommits:    30,
		NumTests:      500,
		NumCodeIssues: 50,

		ArtifactTimeMin: 5,
		ArtifactTimeMax: 10,

		TestRunTimeMin: 15,
		TestRunTimeMax: 30,

		IssueScanTimeMin: 10,
		IssueScanTimeMax: 20,

		CveScanTimeMin: 15,
		CveScanTimeMax: 20,
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
			Project: ent.NewProjectModelCreate().SetName(opt.Project),
			Repo:    ent.NewRepoModelCreate().SetName(opt.Name),
			Release: ent.NewReleaseModelCreate().SetName(opt.Name).SetVersion(version),
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
			cTime = cTime.Add(time.Minute * time.Duration(
				rand.Intn(opt.ArtifactTimeMax-opt.ArtifactTimeMin+1)+opt.ArtifactTimeMin,
			))
			artReq := api.ArtifactLogRequest{
				Artifact: ent.NewArtifactModelCreate().
					SetName("dummy").
					SetSha256(fmt.Sprintf("%x", sha256.Sum256([]byte(hash)))).
					SetType(artifact.TypeDocker).
					SetTime(cTime),
				Commit: &hash,
			}

			data.Artifacts = append(data.Artifacts, &artReq)
		}
		//
		// Test Run & Cases
		//
		{
			cTime = cTime.Add(time.Minute * time.Duration(
				rand.Intn(opt.TestRunTimeMax-opt.TestRunTimeMin+1)+opt.TestRunTimeMin,
			))
			hostname, _ := os.Hostname()
			run := api.TestRun{
				TestRunModelCreate: *ent.NewTestRunModelCreate().
					SetTool("gotest").
					SetMetadata(schema.Metadata{
						"env": map[string]interface{}{
							"hostname": hostname,
						},
					}).
					SetTime(cTime),
			}

			var passAll = false
			if rand.Intn(10) > 3 {
				passAll = true
			}
			for j := 1; j <= opt.NumTests; j++ {
				// Calculate the result based on the test case number. The higher
				// the number, the more likely to fail
				// Create a failChance percentage
				failChance := (float64(j*100) / float64(opt.NumTests)) / 100.0
				result := passAll || (float64(rand.Intn(100))*failChance) <= 80.0
				run.TestCases = append(run.TestCases, &api.TestRunCase{
					TestCaseModelCreate: *ent.NewTestCaseModelCreate().
						SetName(fmt.Sprintf("test_case_%d", j)).
						SetMessage("TODO").
						SetResult(result).
						SetElapsed(0),
				})
			}
			data.TestRuns = append(data.TestRuns, &api.TestRunRequest{
				Commit:  &hash,
				TestRun: &run,
			})
		}

		//
		// Code Issues
		//
		{
			cTime = cTime.Add(time.Minute * time.Duration(
				rand.Intn(opt.IssueScanTimeMax-opt.IssueScanTimeMin+1)+opt.IssueScanTimeMin,
			))
			scan := api.CodeScan{
				CodeScanModelCreate: *ent.NewCodeScanModelCreate().
					SetTool("gosec").
					SetTime(cTime),
			}

			for j := 0; j < opt.NumCodeIssues; j++ {
				// cweID := fmt.Sprintf("CWE-%d", j)
				// tx.CWE.Query().Where(cwe.CweID()).Only(ctx)
				scan.Issues = append(scan.Issues, &api.CodeScanIssue{
					CodeIssueModelCreate: *ent.NewCodeIssueModelCreate().
						SetRuleID(fmt.Sprintf("rule_%d", j)).
						SetMessage("TODO").
						SetSeverity(codeissue.SeverityHigh).
						SetType(codeissue.TypeSecurity),
				})
			}
			data.CodeScans = append(data.CodeScans, &api.CodeScanRequest{
				Commit:   &hash,
				CodeScan: &scan,
			})
		}

		//
		// Components
		//
		{
			cTime = cTime.Add(time.Minute * time.Duration(
				rand.Intn(opt.CveScanTimeMax-opt.CveScanTimeMin+1)+opt.CveScanTimeMin,
			))
			scan := api.CodeScan{
				CodeScanModelCreate: *ent.NewCodeScanModelCreate().
					SetTool("spdx-sbom-generator").
					SetTime(cTime),
			}

			//
			// TODO:: set components with "real" data...
			//
			var components []*ent.Component
			for _, c := range components {
				scan.Components = append(scan.Components, &api.CodeScanComponent{
					ComponentModelCreate: *ent.NewComponentModelCreate().
						SetName(c.Name).SetVendor(c.Vendor).SetVersion(c.Version),
				})
			}
			data.CodeScans = append(data.CodeScans, &api.CodeScanRequest{
				Commit:   &hash,
				CodeScan: &scan,
			})
			// _, err := store.SaveCodeScan(&api.CodeScanRequest{
			// 	Commit:   &hash,
			// 	CodeScan: &scan,
			// })
			// if err != nil {
			// 	return err
			// }
		}

		// Prepend the release so that the dates of the releases go in ascending order
		releases = append([]ReleaseData{data}, releases...)
	}

	return releases
}
