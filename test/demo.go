package test

import (
	"context"
	"crypto/sha1"
	"crypto/sha256"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/valocode/bubbly/ent"
	"github.com/valocode/bubbly/ent/artifact"
	"github.com/valocode/bubbly/ent/codeissue"
	"github.com/valocode/bubbly/ent/release"
	"github.com/valocode/bubbly/ent/releaseentry"
	"github.com/valocode/bubbly/integrations"
	"github.com/valocode/bubbly/store"
)

const (
	projectName = "dummy"
	repoName    = "github.com/valocode/dummy"
	branchName  = "main"

	numCommits    = 500
	numTestCases  = 500
	numCodeIssues = 10

	artifactTimeMin = 5
	artifactTimeMax = 10

	testRunTimeMin = 15
	testRunTimeMax = 30

	codeScanTimeMin = 15
	codeScanTimeMax = 30

	cveTimeMin = 25
	cveTimeMax = 35
)

func SaveSPDXData(client *ent.Client) error {
	list, err := integrations.FetchSPDXLicenses()
	if err != nil {
		return err
	}
	ctx := context.Background()
	for _, lic := range list.Licenses {
		l, err := client.License.Create().
			SetSpdxID(lic.LicenseID).
			SetName(lic.Name).
			SetDetailsURL(lic.DetailsURL).
			SetIsOsiApproved(lic.IsOSIApproved).SetReference(lic.Reference).
			Save(ctx)
		if err != nil {
			return err
		}
		fmt.Println("Created license: " + l.String())
	}
	return nil
}
func SaveCVEData(client *ent.Client) error {
	resp, err := integrations.FetchCVEs(integrations.DefaultCVEFetchOptions())
	if err != nil {
		return fmt.Errorf("error fetching CVEs: %w", err)
	}
	var graph ent.DataGraph
	for _, cve := range resp.Result.CVEItems {
		var description string
		// Get the english description
		for _, d := range cve.CVEItem.Description.DescriptionData {
			if d.Lang == "en" {
				description = d.Value
			}
		}

		graph.RootNodes = append(graph.RootNodes, ent.NewCVENode().
			SetCveID(cve.CVEItem.CVEDataMeta.ID).
			SetDescription(description).
			SetSeverityScore(float64(cve.Impact.BaseMetricV3.CVSSV3.BaseScore)).
			Node(),
		)
	}
	if err := graph.Save(client); err != nil {
		return fmt.Errorf("error saving CVEs: %w", err)
	}
	return nil
}

func FailSomeRandomReleases(db *store.Store) error {
	client := db.Client()
	releases, err := client.Release.Query().All(context.Background())
	if err != nil {
		log.Fatal("getting releases: ", err)
	}
	// Fail some release within the first 10
	nextEval := rand.Intn(10)
	for i := 0; i < len(releases); i++ {
		if i == nextEval {
			if err := db.EvaluateRelease(releases[i].ID); err != nil {
				return fmt.Errorf("evaluating release: %w", err)
			}
			// every ~5 releases should be failed
			nextEval += rand.Intn(10)
			continue
		}
		_, err := client.Release.UpdateOneID(releases[i].ID).
			SetStatus(release.StatusReady).
			Save(context.Background())
		if err != nil {
			return fmt.Errorf("setting release ready: %w", err)
		}
	}
	return nil
}

func CreateDummyData(client *ent.Client) error {
	// Create a silly CVE rule... this doesn't work yet with queries
	// if err := ent.NewCVERuleNode().
	// 	SetCve(ent.NewCVENode().SetCveID("CVE-1")).
	// 	Graph().
	// 	Save(client); err != nil {
	// 	return err
	// }

	var cves []*ent.CVE
	{
		var err error
		cves, err = client.CVE.Query().All(context.Background())
		if err != nil {
			return fmt.Errorf("error getting the CVE list: %w", err)
		}
	}

	project := ent.NewProjectNode().SetName(projectName)
	repo := ent.NewRepoNode().
		SetName(repoName).
		SetProject(project)

	cTime := time.Now()

	for i := 1; i <= numCommits; i++ {
		var (
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
		release := ent.NewReleaseNode().
			SetCommit(
				ent.NewGitCommitNode().
					SetHash(hash).
					SetBranch(branchName).
					SetTag(tag).
					SetTime(cTime).
					SetRepo(repo),
			).
			SetProject(project).
			SetName(repoName).
			SetVersion(version)

		//
		// Artifact
		//
		{
			cTime = cTime.Add(time.Minute * time.Duration(
				rand.Intn(artifactTimeMax-artifactTimeMin+1)+artifactTimeMin,
			))
			artifact := ent.NewArtifactNode().
				SetName("dummy").
				SetSha256(fmt.Sprintf("%x", sha256.Sum256([]byte(hash)))).
				SetType(artifact.TypeDocker).
				// Fabriate the entry as well to override the one that is created
				SetEntry(
					ent.NewReleaseEntryNode().
						SetTime(cTime).
						SetType(releaseentry.TypeArtifact),
				)
			release.AddArtifacts(artifact)

		}
		//
		// Test Run & Cases
		//
		{
			cTime = cTime.Add(time.Minute * time.Duration(
				rand.Intn(testRunTimeMax-testRunTimeMin+1)+testRunTimeMin,
			))
			run := ent.NewTestRunNode().
				SetRelease(release).
				SetTool("gotest").
				SetEntry(
					ent.NewReleaseEntryNode().
						SetTime(cTime).
						SetType(releaseentry.TypeTestRun),
				)

			for j := 1; j <= numTestCases; j++ {
				// Calculate the result based on the test case number. The higher
				// the number, the more likely to fail
				// Create a failChance percentage
				failChance := (float64(j*100) / float64(numTestCases)) / 100.0
				result := (float64(rand.Intn(100)) * failChance) <= 80.0
				// if !result {
				// 	fmt.Printf("Test Case will fail: %d\n", j)
				// }
				ent.NewTestCaseNode().
					SetRun(run).
					SetName(fmt.Sprintf("test_case_%d", j)).
					SetMessage("TODO").
					SetResult(result).
					SetElapsed(0).
					SetRun(run).
					// Optimization to tell bubbly to only create (not update)
					SetOperation(ent.NodeOperationCreate)
			}
			release.AddTestRuns(run)
		}

		//
		// Code Issues
		//
		{
			cTime = cTime.Add(time.Minute * time.Duration(
				rand.Intn(codeScanTimeMax-codeScanTimeMin+1)+codeScanTimeMin,
			))
			scan := ent.NewCodeScanNode().
				SetRelease(release).
				SetTool("gosec").
				SetEntry(
					ent.NewReleaseEntryNode().
						SetTime(cTime).
						SetType(releaseentry.TypeCodeScan),
				)

			for j := 0; j < numCodeIssues; j++ {
				// cweID := fmt.Sprintf("CWE-%d", j)
				// tx.CWE.Query().Where(cwe.CweID()).Only(ctx)
				ent.NewCodeIssueNode().
					SetScan(scan).
					SetRuleID(fmt.Sprintf("rule_%d", j)).
					SetMessage("TODO").
					SetSeverity(codeissue.SeverityHigh).
					SetType(codeissue.TypeSecurity)
			}
		}

		//
		// Components
		//
		{
			cTime = cTime.Add(time.Minute * time.Duration(
				rand.Intn(cveTimeMax-cveTimeMin+1)+cveTimeMin,
			))
			scan := ent.NewCodeScanNode().
				SetRelease(release).
				SetTool("snyk").
				SetEntry(
					ent.NewReleaseEntryNode().
						SetTime(cTime).
						SetType(releaseentry.TypeCodeScan),
				)

			for j := 0; j < 10; j++ {
				comp := ent.NewComponentNode().
					SetName("comp1").
					SetDescription("TODO").
					SetVendor("TODO").
					SetURL("TODO").
					SetVersion(fmt.Sprintf("0.0.%d", j))
				compUse := ent.NewComponentUseNode().
					// ent.NewComponentUseNode().
					AddScans(scan).
					SetComponent(comp)

				release.AddComponents(compUse)

				//
				// Vulnerabilities
				//
				{
					// HACK: do this once only
					if i == 1 {
						// Set a default CVE ID in case we don't have any CVEs
						cveID := fmt.Sprintf("CVE-%d", j)
						if len(cves) > 0 {
							cveID = cves[j].CveID
						}
						cve := ent.NewCVENode().
							SetCveID(cveID).AddComponents(comp)

						if j == 9 {
							ent.NewCVERuleNode().SetCve(cve).SetName("ignore").AddProject(project)
						}

					}
					// ent.NewVulnerabilityNode().
					// 	SetComponent(compUse).
					// 	SetRelease(release).
					// 	SetCve(
					// 		ent.NewCVENode().
					// 			SetCveID(cveID),
					// 	).AddScans(scan)
				}

			}
		}

		if err := release.Graph().Save(client); err != nil {
			return err
		}
	}

	return nil
}
