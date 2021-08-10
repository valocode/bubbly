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

var demoData = []DemoRepoOptions{
	{
		Name:          "demo",
		Project:       "bubbly",
		Branch:        "main",
		NumCommits:    250,
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

	var (
		licenses   []*ent.License
		licenseIDs = []string{"MIT", "GPL-3.0", "MPL-2.0", "Apache-2.0"}
		cves       []*ent.CVE
		components []*ent.Component
		ctx        = context.Background()
	)
	{
		var err error
		licenses, err = client.License.Query().All(ctx)
		if err != nil {
			return err
		}
		cves, err = client.CVE.Query().All(ctx)
		if err != nil {
			return fmt.Errorf("error getting the CVE list: %w", err)
		}
	}

	// Create some components
	{
		mods, err := integrations.ParseSBOM("lang/testdata/spdx-sbom-generator.json")
		if err != nil {
			return err
		}
		for i := 0; i < len(mods); i++ {
			mod := mods[i]
			// This  *should* only be true for the first/root module
			if mod.Version == "" {
				continue
			}
			cveID := fmt.Sprintf("CVE-%d", i)
			if len(cves) > 0 {
				cveID = cves[i].CveID
			}
			cve := ent.NewCVENode().
				SetCveID(cveID)
			spdxID := fmt.Sprintf("GPL-%d", i)
			if len(licenses) > 0 {
				spdxID = licenseIDs[i%len(licenseIDs)]
				// if i >= len(licenseIDs) {
				// }
			}
			license := ent.NewLicenseNode().SetSpdxID(spdxID)

			comp := ent.NewComponentNode().
				SetName(mod.Name).
				SetVendor(mod.Supplier.Name).
				SetVersion(mod.Version).SetURL(mod.PackageURL).
				AddCves(cve).
				AddLicenses(license)
			if err := comp.Graph().Save(client); err != nil {
				return err
			}

			components = append(components, comp.Value.(*ent.Component))

		}
	}

	for _, opt := range demoData {
		if err := createRepoData(client, opt, components); err != nil {
			return err
		}
	}

	return nil
}

func createRepoData(client *ent.Client, opt DemoRepoOptions, components []*ent.Component) error {
	project := ent.NewProjectNode().SetName(opt.Project)
	repo := ent.NewRepoNode().
		SetName(opt.Name).
		SetProject(project)

	cTime := time.Now()

	for i := 1; i <= opt.NumCommits; i++ {
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
					SetBranch(opt.Branch).
					SetTag(tag).
					SetTime(cTime).
					SetRepo(repo),
			).
			SetProject(project).
			SetName(opt.Name).
			SetVersion(version)

		//
		// Artifact
		//
		{
			cTime = cTime.Add(time.Minute * time.Duration(
				rand.Intn(opt.ArtifactTimeMax-opt.ArtifactTimeMin+1)+opt.ArtifactTimeMin,
			))
			art := ent.NewArtifactNode().
				SetName("dummy").
				SetSha256(fmt.Sprintf("%x", sha256.Sum256([]byte(hash)))).
				SetType(artifact.TypeDocker).
				// Fabriate the entry as well to override the one that is created
				SetEntry(
					ent.NewReleaseEntryNode().
						SetTime(cTime).
						SetType(releaseentry.TypeArtifact),
				)
			release.AddArtifacts(art)

		}
		//
		// Test Run & Cases
		//
		{
			cTime = cTime.Add(time.Minute * time.Duration(
				rand.Intn(opt.TestRunTimeMax-opt.TestRunTimeMin+1)+opt.TestRunTimeMin,
			))
			run := ent.NewTestRunNode().
				SetRelease(release).
				SetTool("gotest").
				SetEntry(
					ent.NewReleaseEntryNode().
						SetTime(cTime).
						SetType(releaseentry.TypeTestRun),
				)

			for j := 1; j <= opt.NumTests; j++ {
				// Calculate the result based on the test case number. The higher
				// the number, the more likely to fail
				// Create a failChance percentage
				failChance := (float64(j*100) / float64(opt.NumTests)) / 100.0
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
				rand.Intn(opt.IssueScanTimeMax-opt.IssueScanTimeMin+1)+opt.IssueScanTimeMin,
			))
			scan := ent.NewCodeScanNode().
				SetRelease(release).
				SetTool("spdx-sbom-generator").
				SetEntry(
					ent.NewReleaseEntryNode().
						SetTime(cTime).
						SetType(releaseentry.TypeCodeScan),
				)

			for j := 0; j < opt.NumCodeIssues; j++ {
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
				rand.Intn(opt.CveScanTimeMax-opt.CveScanTimeMin+1)+opt.CveScanTimeMin,
			))
			scan := ent.NewCodeScanNode().
				SetRelease(release).
				SetTool("snyk").
				SetEntry(
					ent.NewReleaseEntryNode().
						SetTime(cTime).
						SetType(releaseentry.TypeCodeScan),
				)

			for _, c := range components {
				comp := ent.NewComponentNode().
					SetName(c.Name)
				compUse := ent.NewComponentUseNode().
					// ent.NewComponentUseNode().
					AddScans(scan).
					SetComponent(comp)

				release.AddComponents(compUse)

			}
		}

		if err := release.Graph().Save(client); err != nil {
			return err
		}
	}

	return nil
}