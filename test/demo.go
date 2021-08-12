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
	"github.com/valocode/bubbly/ent/model"
	"github.com/valocode/bubbly/ent/release"
	"github.com/valocode/bubbly/ent/vulnerability"
	"github.com/valocode/bubbly/integrations"
	"github.com/valocode/bubbly/store"
	"github.com/valocode/bubbly/store/api"
)

type DemoRepoOptions struct {
	Name string
	// Project string
	Branch string

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
		Name: "demo",
		// Project:       "bubbly",
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
	ctx := context.Background()
	for _, vuln := range resp.Result.CVEItems {
		var description string
		// Get the english description
		for _, d := range vuln.CVEItem.Description.DescriptionData {
			if d.Lang == "en" {
				description = d.Value
			}
		}
		_, err := client.Vulnerability.Query().Where(
			vulnerability.Vid(vuln.CVEItem.CVEDataMeta.ID),
		).Only(ctx)
		if err != nil {
			if !ent.IsNotFound(err) {
				return fmt.Errorf("error checking vulnerability: %w", err)
			}
			_, err = client.Vulnerability.Create().
				SetVid(vuln.CVEItem.CVEDataMeta.ID).
				SetDescription(description).
				SetSeverityScore(float64(vuln.Impact.BaseMetricV3.CVSSV3.BaseScore)).
				Save(ctx)
			if err != nil {
				return fmt.Errorf("error creating vulnerability: %w", err)
			}
		}
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

func CreateDummyData(store *store.Store) error {

	var (
		client     = store.Client()
		vulns      []*ent.Vulnerability
		components []*ent.Component
		// licenses   []*ent.License
		// licenseIDs = []string{"MIT", "GPL-3.0", "MPL-2.0", "Apache-2.0"}
		ctx = context.Background()
	)
	{
		var err error
		_, err = client.License.Query().All(ctx)
		if err != nil {
			return err
		}
		vulns, err = client.Vulnerability.Query().All(ctx)
		if err != nil {
			return fmt.Errorf("error getting the Vulnerability list: %w", err)
		}
	}

	// Create some components
	{
		mods, err := integrations.ParseSBOM("adapter/testdata/spdx-sbom-generator.json")
		if err != nil {
			return err
		}

		var compVulns = make([]*api.ComponentVulnerability, 0, len(mods))
		for i := 0; i < len(mods); i++ {
			mod := mods[i]

			// This  *should* only be true for the first/root module
			if mod.Version == "" {
				continue
			}
			dbComp, err := client.Component.Create().
				SetName(mod.Name).
				SetVendor(mod.Supplier.Name).
				SetVersion(mod.Version).SetURL(mod.PackageURL).
				Save(ctx)
			if err != nil {
				return err
			}
			components = append(components, dbComp)
			compVuln := api.ComponentVulnerability{
				Component: api.Component{ComponentModel: *model.NewComponentModel().SetID(dbComp.ID)},
			}

			vulnID := fmt.Sprintf("CVE-%d", i)
			if len(vulns) > 0 {
				vulnID = vulns[i].Vid
			}

			compVuln.Vulnerabilities = append(compVuln.Vulnerabilities,
				&api.Vulnerability{
					VulnerabilityModel: *model.NewVulnerabilityModel().
						SetVid(vulnID),
				},
			)
			compVulns = append(compVulns, &compVuln)

			// Demo review - TODO...
			// project := ent.NewProjectNode().SetName("bubbly")
			// ent.NewVulnerabilityReviewNode().
			// 	SetName("My demo review...").
			// 	SetDecision(vulnerabilityreview.DecisionInProgress).
			// 	SetVulnerability(vuln).
			// 	AddProjects(project)
			// end TODO
			// spdxID := fmt.Sprintf("GPL-%d", i)
			// if len(licenses) > 0 {
			// 	spdxID = licenseIDs[i%len(licenseIDs)]
			// }
			// license := ent.NewLicenseNode().SetSpdxID(spdxID)

		}
		if err := store.SaveComponentVulnerabilities(&api.ComponentVulnerabilityRequest{
			Components: compVulns,
		}); err != nil {
			return fmt.Errorf("error saving component vulnerabilities: %w", err)
		}
	}

	for _, opt := range demoData {
		if err := createRepoData(store, opt, components); err != nil {
			return err
		}
	}

	return nil
}

func createRepoData(store *store.Store, opt DemoRepoOptions, components []*ent.Component) error {

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
		relReq := api.ReleaseCreateRequest{
			Repo: model.NewRepoModel().SetName(opt.Name),
			Commit: model.NewGitCommitModel().
				SetHash(hash).
				SetBranch(opt.Branch).
				SetTag(tag).
				SetTime(cTime),
		}
		_, err := store.CreateRelease(&relReq)
		if err != nil {
			return err
		}
		//
		// Artifact
		//
		{
			cTime = cTime.Add(time.Minute * time.Duration(
				rand.Intn(opt.ArtifactTimeMax-opt.ArtifactTimeMin+1)+opt.ArtifactTimeMin,
			))
			artReq := api.ArtifactLogRequest{
				Artifact: model.NewArtifactModel().
					SetName("dummy").
					SetSha256(fmt.Sprintf("%x", sha256.Sum256([]byte(hash)))).
					SetType(artifact.TypeDocker),
				Commit: &hash,
			}

			_, err := store.LogArtifact(&artReq)
			if err != nil {
				return err
			}
		}
		//
		// Test Run & Cases
		//
		{
			cTime = cTime.Add(time.Minute * time.Duration(
				rand.Intn(opt.TestRunTimeMax-opt.TestRunTimeMin+1)+opt.TestRunTimeMin,
			))

			run := api.TestRun{
				TestRunModel: *model.NewTestRunModel().SetTool("gotest"),
			}
			// run := ent.NewTestRunNode().
			// 	SetRelease(release).
			// 	SetEntry(
			// 		ent.NewReleaseEntryNode().
			// 			SetTime(cTime).
			// 			SetType(releaseentry.TypeTestRun),
			// 	)

			for j := 1; j <= opt.NumTests; j++ {
				// Calculate the result based on the test case number. The higher
				// the number, the more likely to fail
				// Create a failChance percentage
				failChance := (float64(j*100) / float64(opt.NumTests)) / 100.0
				result := (float64(rand.Intn(100)) * failChance) <= 80.0
				// if !result {
				// 	fmt.Printf("Test Case will fail: %d\n", j)
				// }
				run.TestCases = append(run.TestCases, &api.TestRunCase{
					TestCaseModel: *model.NewTestCaseModel().
						SetName(fmt.Sprintf("test_case_%d", j)).
						SetMessage("TODO").
						SetResult(result).
						SetElapsed(0),
				})
			}
			_, err := store.SaveTestRun(&api.TestRunRequest{
				Commit:  &hash,
				TestRun: &run,
			})
			if err != nil {
				return err
			}
		}

		//
		// Code Issues
		//
		{
			cTime = cTime.Add(time.Minute * time.Duration(
				rand.Intn(opt.IssueScanTimeMax-opt.IssueScanTimeMin+1)+opt.IssueScanTimeMin,
			))
			scan := api.CodeScan{
				CodeScanModel: *model.NewCodeScanModel().
					SetTool("gosec"),
			}
			// scan := ent.NewCodeScanNode().
			// 	SetRelease(release).
			// 	SetEntry(
			// 		ent.NewReleaseEntryNode().
			// 			SetTime(cTime).
			// 			SetType(releaseentry.TypeCodeScan),
			// 	)

			for j := 0; j < opt.NumCodeIssues; j++ {
				// cweID := fmt.Sprintf("CWE-%d", j)
				// tx.CWE.Query().Where(cwe.CweID()).Only(ctx)
				scan.Issues = append(scan.Issues, &api.CodeScanIssue{
					CodeIssueModel: *model.NewCodeIssueModel().
						SetRuleID(fmt.Sprintf("rule_%d", j)).
						SetMessage("TODO").
						SetSeverity(codeissue.SeverityHigh).
						SetType(codeissue.TypeSecurity),
				})
			}
			_, err := store.SaveCodeScan(&api.CodeScanRequest{
				Commit:   &hash,
				CodeScan: &scan,
			})
			if err != nil {
				return err
			}
		}

		//
		// Components
		//
		{
			cTime = cTime.Add(time.Minute * time.Duration(
				rand.Intn(opt.CveScanTimeMax-opt.CveScanTimeMin+1)+opt.CveScanTimeMin,
			))
			scan := api.CodeScan{
				CodeScanModel: *model.NewCodeScanModel().
					SetTool("spdx-sbom-generator"),
			}
			// scan := ent.NewCodeScanNode().
			// 	SetRelease(release).
			// 	SetTool("snyk").
			// 	SetEntry(
			// 		ent.NewReleaseEntryNode().
			// 			SetTime(cTime).
			// 			SetType(releaseentry.TypeCodeScan),
			// 	)

			for _, c := range components {
				scan.Components = append(scan.Components, &api.CodeScanComponent{
					ComponentModel: *model.NewComponentModel().
						SetName(c.Name).SetVendor(c.Vendor).SetVersion(c.Version),
				})
			}
			_, err := store.SaveCodeScan(&api.CodeScanRequest{
				Commit:   &hash,
				CodeScan: &scan,
			})
			if err != nil {
				return err
			}
		}
	}

	return nil
}
