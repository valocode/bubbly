package gql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/valocode/bubbly/ent"
)

func (r *cVEResolver) Found(ctx context.Context, obj *ent.CVE, first *int, last *int, where *ent.VulnerabilityWhereInput) ([]*ent.Vulnerability, error) {
	result, err := obj.Edges.FoundOrErr()
	if ent.IsNotLoaded(err) {
		result, err = obj.QueryFound().Filter(ctx, first, last, nil, where)
	}
	return result, err
}

func (r *cVEResolver) Rules(ctx context.Context, obj *ent.CVE, first *int, last *int, where *ent.CVERuleWhereInput, orderBy *ent.CVERuleOrder) ([]*ent.CVERule, error) {
	result, err := obj.Edges.RulesOrErr()
	if ent.IsNotLoaded(err) {
		result, err = obj.QueryRules().Filter(ctx, first, last, orderBy, where)
	}
	return result, err
}

func (r *cVERuleResolver) Project(ctx context.Context, obj *ent.CVERule, first *int, last *int, where *ent.ProjectWhereInput, orderBy *ent.ProjectOrder) ([]*ent.Project, error) {
	result, err := obj.Edges.ProjectOrErr()
	if ent.IsNotLoaded(err) {
		result, err = obj.QueryProject().Filter(ctx, first, last, orderBy, where)
	}
	return result, err
}

func (r *cVERuleResolver) Repo(ctx context.Context, obj *ent.CVERule, first *int, last *int, where *ent.RepoWhereInput, orderBy *ent.RepoOrder) ([]*ent.Repo, error) {
	result, err := obj.Edges.RepoOrErr()
	if ent.IsNotLoaded(err) {
		result, err = obj.QueryRepo().Filter(ctx, first, last, orderBy, where)
	}
	return result, err
}

func (r *cVEScanResolver) Vulnerabilities(ctx context.Context, obj *ent.CVEScan, first *int, last *int, where *ent.VulnerabilityWhereInput) ([]*ent.Vulnerability, error) {
	result, err := obj.Edges.VulnerabilitiesOrErr()
	if ent.IsNotLoaded(err) {
		result, err = obj.QueryVulnerabilities().Filter(ctx, first, last, nil, where)
	}
	return result, err
}

func (r *cWEResolver) Issues(ctx context.Context, obj *ent.CWE, first *int, last *int, where *ent.CodeIssueWhereInput, orderBy *ent.CodeIssueOrder) ([]*ent.CodeIssue, error) {
	result, err := obj.Edges.IssuesOrErr()
	if ent.IsNotLoaded(err) {
		result, err = obj.QueryIssues().Filter(ctx, first, last, orderBy, where)
	}
	return result, err
}

func (r *codeIssueResolver) Cwe(ctx context.Context, obj *ent.CodeIssue, first *int, last *int, where *ent.CWEWhereInput, orderBy *ent.CWEOrder) ([]*ent.CWE, error) {
	result, err := obj.Edges.CweOrErr()
	if ent.IsNotLoaded(err) {
		result, err = obj.QueryCwe().Filter(ctx, first, last, orderBy, where)
	}
	return result, err
}

func (r *codeScanResolver) Issues(ctx context.Context, obj *ent.CodeScan, first *int, last *int, where *ent.CodeIssueWhereInput, orderBy *ent.CodeIssueOrder) ([]*ent.CodeIssue, error) {
	result, err := obj.Edges.IssuesOrErr()
	if ent.IsNotLoaded(err) {
		result, err = obj.QueryIssues().Filter(ctx, first, last, orderBy, where)
	}
	return result, err
}

func (r *componentResolver) Vulnerabilities(ctx context.Context, obj *ent.Component, first *int, last *int, where *ent.VulnerabilityWhereInput) ([]*ent.Vulnerability, error) {
	result, err := obj.Edges.VulnerabilitiesOrErr()
	if ent.IsNotLoaded(err) {
		result, err = obj.QueryVulnerabilities().Filter(ctx, first, last, nil, where)
	}
	return result, err
}

func (r *componentResolver) Licenses(ctx context.Context, obj *ent.Component, first *int, last *int, where *ent.LicenseWhereInput, orderBy *ent.LicenseOrder) ([]*ent.License, error) {
	result, err := obj.Edges.LicensesOrErr()
	if ent.IsNotLoaded(err) {
		result, err = obj.QueryLicenses().Filter(ctx, first, last, orderBy, where)
	}
	return result, err
}

func (r *componentResolver) Release(ctx context.Context, obj *ent.Component, first *int, last *int, where *ent.ReleaseWhereInput, orderBy *ent.ReleaseOrder) ([]*ent.Release, error) {
	result, err := obj.Edges.ReleaseOrErr()
	if ent.IsNotLoaded(err) {
		result, err = obj.QueryRelease().Filter(ctx, first, last, orderBy, where)
	}
	return result, err
}

func (r *licenseResolver) Components(ctx context.Context, obj *ent.License, first *int, last *int, where *ent.ComponentWhereInput, orderBy *ent.ComponentOrder) ([]*ent.Component, error) {
	result, err := obj.Edges.ComponentsOrErr()
	if ent.IsNotLoaded(err) {
		result, err = obj.QueryComponents().Filter(ctx, first, last, orderBy, where)
	}
	return result, err
}

func (r *licenseResolver) Usages(ctx context.Context, obj *ent.License, first *int, last *int, where *ent.LicenseUsageWhereInput) ([]*ent.LicenseUsage, error) {
	result, err := obj.Edges.UsagesOrErr()
	if ent.IsNotLoaded(err) {
		result, err = obj.QueryUsages().Filter(ctx, first, last, nil, where)
	}
	return result, err
}

func (r *licenseScanResolver) Licenses(ctx context.Context, obj *ent.LicenseScan, first *int, last *int, where *ent.LicenseUsageWhereInput) ([]*ent.LicenseUsage, error) {
	result, err := obj.Edges.LicensesOrErr()
	if ent.IsNotLoaded(err) {
		result, err = obj.QueryLicenses().Filter(ctx, first, last, nil, where)
	}
	return result, err
}

func (r *projectResolver) Repos(ctx context.Context, obj *ent.Project, first *int, last *int, where *ent.RepoWhereInput, orderBy *ent.RepoOrder) ([]*ent.Repo, error) {
	result, err := obj.Edges.ReposOrErr()
	if ent.IsNotLoaded(err) {
		result, err = obj.QueryRepos().Filter(ctx, first, last, orderBy, where)
	}
	return result, err
}

func (r *projectResolver) Releases(ctx context.Context, obj *ent.Project, first *int, last *int, where *ent.ReleaseWhereInput, orderBy *ent.ReleaseOrder) ([]*ent.Release, error) {
	result, err := obj.Edges.ReleasesOrErr()
	if ent.IsNotLoaded(err) {
		result, err = obj.QueryReleases().Filter(ctx, first, last, orderBy, where)
	}
	return result, err
}

func (r *projectResolver) CveRules(ctx context.Context, obj *ent.Project, first *int, last *int, where *ent.CVERuleWhereInput, orderBy *ent.CVERuleOrder) ([]*ent.CVERule, error) {
	result, err := obj.Edges.CveRulesOrErr()
	if ent.IsNotLoaded(err) {
		result, err = obj.QueryCveRules().Filter(ctx, first, last, orderBy, where)
	}
	return result, err
}

func (r *queryResolver) TestRunConnection(ctx context.Context, first *int, last *int, before *ent.Cursor, after *ent.Cursor, orderBy *ent.TestRunOrder, where *ent.TestRunWhereInput) (*ent.TestRunConnection, error) {
	return r.client.TestRun.Query().Paginate(ctx, after, first, before, last,
		ent.WithTestRunOrder(orderBy),
		ent.WithTestRunFilter(where.Filter),
	)
}

func (r *queryResolver) Vulnerability(ctx context.Context, first *int, last *int, where *ent.VulnerabilityWhereInput) ([]*ent.Vulnerability, error) {
	return r.client.Vulnerability.Query().Filter(ctx, first, last, nil, where)
}

func (r *queryResolver) Commit(ctx context.Context, first *int, last *int, orderBy *ent.GitCommitOrder, where *ent.GitCommitWhereInput) ([]*ent.GitCommit, error) {
	return r.client.GitCommit.Query().Filter(ctx, first, last, orderBy, where)
}

func (r *queryResolver) License(ctx context.Context, first *int, last *int, orderBy *ent.LicenseOrder, where *ent.LicenseWhereInput) ([]*ent.License, error) {
	return r.client.License.Query().Filter(ctx, first, last, orderBy, where)
}

func (r *queryResolver) LicenseScan(ctx context.Context, first *int, last *int, orderBy *ent.LicenseScanOrder, where *ent.LicenseScanWhereInput) ([]*ent.LicenseScan, error) {
	return r.client.LicenseScan.Query().Filter(ctx, first, last, orderBy, where)
}

func (r *queryResolver) Release(ctx context.Context, first *int, last *int, orderBy *ent.ReleaseOrder, where *ent.ReleaseWhereInput) ([]*ent.Release, error) {
	return r.client.Release.Query().Filter(ctx, first, last, orderBy, where)
}

func (r *queryResolver) ReleaseEntry(ctx context.Context, first *int, last *int, orderBy *ent.ReleaseEntryOrder, where *ent.ReleaseEntryWhereInput) ([]*ent.ReleaseEntry, error) {
	return r.client.ReleaseEntry.Query().Filter(ctx, first, last, orderBy, where)
}

func (r *queryResolver) ReleaseCheck(ctx context.Context, first *int, last *int, orderBy *ent.ReleaseCheckOrder, where *ent.ReleaseCheckWhereInput) ([]*ent.ReleaseCheck, error) {
	return r.client.ReleaseCheck.Query().Filter(ctx, first, last, orderBy, where)
}

func (r *queryResolver) ReleaseEntryConnection(ctx context.Context, first *int, last *int, before *ent.Cursor, after *ent.Cursor, orderBy *ent.ReleaseEntryOrder, where *ent.ReleaseEntryWhereInput) (*ent.ReleaseEntryConnection, error) {
	return r.client.ReleaseEntry.Query().Paginate(ctx, after, first, before, last,
		ent.WithReleaseEntryOrder(orderBy),
		ent.WithReleaseEntryFilter(where.Filter),
	)
}

func (r *queryResolver) ArtifactConnection(ctx context.Context, first *int, last *int, before *ent.Cursor, after *ent.Cursor, orderBy *ent.ArtifactOrder, where *ent.ArtifactWhereInput) (*ent.ArtifactConnection, error) {
	return r.client.Artifact.Query().Paginate(ctx, after, first, before, last,
		ent.WithArtifactOrder(orderBy),
		ent.WithArtifactFilter(where.Filter),
	)
}

func (r *queryResolver) CveRuleConnection(ctx context.Context, first *int, last *int, before *ent.Cursor, after *ent.Cursor, orderBy *ent.CVERuleOrder, where *ent.CVERuleWhereInput) (*ent.CVERuleConnection, error) {
	return r.client.CVERule.Query().Paginate(ctx, after, first, before, last,
		ent.WithCVERuleOrder(orderBy),
		ent.WithCVERuleFilter(where.Filter),
	)
}

func (r *queryResolver) Cwe(ctx context.Context, first *int, last *int, orderBy *ent.CWEOrder, where *ent.CWEWhereInput) ([]*ent.CWE, error) {
	return r.client.CWE.Query().Filter(ctx, first, last, orderBy, where)
}

func (r *queryResolver) CodeIssueConnection(ctx context.Context, first *int, last *int, before *ent.Cursor, after *ent.Cursor, orderBy *ent.CodeIssueOrder, where *ent.CodeIssueWhereInput) (*ent.CodeIssueConnection, error) {
	return r.client.CodeIssue.Query().Paginate(ctx, after, first, before, last,
		ent.WithCodeIssueOrder(orderBy),
		ent.WithCodeIssueFilter(where.Filter),
	)
}

func (r *queryResolver) CommitConnection(ctx context.Context, first *int, last *int, before *ent.Cursor, after *ent.Cursor, orderBy *ent.GitCommitOrder, where *ent.GitCommitWhereInput) (*ent.GitCommitConnection, error) {
	return r.client.GitCommit.Query().Paginate(ctx, after, first, before, last,
		ent.WithGitCommitOrder(orderBy),
		ent.WithGitCommitFilter(where.Filter),
	)
}

func (r *queryResolver) CodeScanConnection(ctx context.Context, first *int, last *int, before *ent.Cursor, after *ent.Cursor, orderBy *ent.CodeScanOrder, where *ent.CodeScanWhereInput) (*ent.CodeScanConnection, error) {
	return r.client.CodeScan.Query().Paginate(ctx, after, first, before, last,
		ent.WithCodeScanOrder(orderBy),
		ent.WithCodeScanFilter(where.Filter),
	)
}

func (r *queryResolver) LicenseUsage(ctx context.Context, first *int, last *int, where *ent.LicenseUsageWhereInput) ([]*ent.LicenseUsage, error) {
	return r.client.LicenseUsage.Query().Filter(ctx, first, last, nil, where)
}

func (r *queryResolver) ReleaseCheckConnection(ctx context.Context, first *int, last *int, before *ent.Cursor, after *ent.Cursor, orderBy *ent.ReleaseCheckOrder, where *ent.ReleaseCheckWhereInput) (*ent.ReleaseCheckConnection, error) {
	return r.client.ReleaseCheck.Query().Paginate(ctx, after, first, before, last,
		ent.WithReleaseCheckOrder(orderBy),
		ent.WithReleaseCheckFilter(where.Filter),
	)
}

func (r *queryResolver) TestCase(ctx context.Context, first *int, last *int, orderBy *ent.TestCaseOrder, where *ent.TestCaseWhereInput) ([]*ent.TestCase, error) {
	return r.client.TestCase.Query().Filter(ctx, first, last, orderBy, where)
}

func (r *queryResolver) TestRun(ctx context.Context, first *int, last *int, orderBy *ent.TestRunOrder, where *ent.TestRunWhereInput) ([]*ent.TestRun, error) {
	return r.client.TestRun.Query().Filter(ctx, first, last, orderBy, where)
}

func (r *queryResolver) CveScan(ctx context.Context, first *int, last *int, orderBy *ent.CVEScanOrder, where *ent.CVEScanWhereInput) ([]*ent.CVEScan, error) {
	return r.client.CVEScan.Query().Filter(ctx, first, last, orderBy, where)
}

func (r *queryResolver) CweConnection(ctx context.Context, first *int, last *int, before *ent.Cursor, after *ent.Cursor, orderBy *ent.CWEOrder, where *ent.CWEWhereInput) (*ent.CWEConnection, error) {
	return r.client.CWE.Query().Paginate(ctx, after, first, before, last,
		ent.WithCWEOrder(orderBy),
		ent.WithCWEFilter(where.Filter),
	)
}

func (r *queryResolver) TestCaseConnection(ctx context.Context, first *int, last *int, before *ent.Cursor, after *ent.Cursor, orderBy *ent.TestCaseOrder, where *ent.TestCaseWhereInput) (*ent.TestCaseConnection, error) {
	return r.client.TestCase.Query().Paginate(ctx, after, first, before, last,
		ent.WithTestCaseOrder(orderBy),
		ent.WithTestCaseFilter(where.Filter),
	)
}

func (r *queryResolver) Repo(ctx context.Context, first *int, last *int, orderBy *ent.RepoOrder, where *ent.RepoWhereInput) ([]*ent.Repo, error) {
	return r.client.Repo.Query().Filter(ctx, first, last, orderBy, where)
}

func (r *queryResolver) CveConnection(ctx context.Context, first *int, last *int, before *ent.Cursor, after *ent.Cursor, orderBy *ent.CVEOrder, where *ent.CVEWhereInput) (*ent.CVEConnection, error) {
	return r.client.CVE.Query().Paginate(ctx, after, first, before, last,
		ent.WithCVEOrder(orderBy),
		ent.WithCVEFilter(where.Filter),
	)
}

func (r *queryResolver) Cve(ctx context.Context, first *int, last *int, orderBy *ent.CVEOrder, where *ent.CVEWhereInput) ([]*ent.CVE, error) {
	return r.client.CVE.Query().Filter(ctx, first, last, orderBy, where)
}

func (r *queryResolver) CodeIssue(ctx context.Context, first *int, last *int, orderBy *ent.CodeIssueOrder, where *ent.CodeIssueWhereInput) ([]*ent.CodeIssue, error) {
	return r.client.CodeIssue.Query().Filter(ctx, first, last, orderBy, where)
}

func (r *queryResolver) CodeScan(ctx context.Context, first *int, last *int, orderBy *ent.CodeScanOrder, where *ent.CodeScanWhereInput) ([]*ent.CodeScan, error) {
	return r.client.CodeScan.Query().Filter(ctx, first, last, orderBy, where)
}

func (r *queryResolver) LicenseScanConnection(ctx context.Context, first *int, last *int, before *ent.Cursor, after *ent.Cursor, orderBy *ent.LicenseScanOrder, where *ent.LicenseScanWhereInput) (*ent.LicenseScanConnection, error) {
	return r.client.LicenseScan.Query().Paginate(ctx, after, first, before, last,
		ent.WithLicenseScanOrder(orderBy),
		ent.WithLicenseScanFilter(where.Filter),
	)
}

func (r *queryResolver) CveScanConnection(ctx context.Context, first *int, last *int, before *ent.Cursor, after *ent.Cursor, orderBy *ent.CVEScanOrder, where *ent.CVEScanWhereInput) (*ent.CVEScanConnection, error) {
	return r.client.CVEScan.Query().Paginate(ctx, after, first, before, last,
		ent.WithCVEScanOrder(orderBy),
		ent.WithCVEScanFilter(where.Filter),
	)
}

func (r *queryResolver) ReleaseConnection(ctx context.Context, first *int, last *int, before *ent.Cursor, after *ent.Cursor, orderBy *ent.ReleaseOrder, where *ent.ReleaseWhereInput) (*ent.ReleaseConnection, error) {
	return r.client.Release.Query().Paginate(ctx, after, first, before, last,
		ent.WithReleaseOrder(orderBy),
		ent.WithReleaseFilter(where.Filter),
	)
}

func (r *queryResolver) RepoConnection(ctx context.Context, first *int, last *int, before *ent.Cursor, after *ent.Cursor, orderBy *ent.RepoOrder, where *ent.RepoWhereInput) (*ent.RepoConnection, error) {
	return r.client.Repo.Query().Paginate(ctx, after, first, before, last,
		ent.WithRepoOrder(orderBy),
		ent.WithRepoFilter(where.Filter),
	)
}

func (r *queryResolver) CveRule(ctx context.Context, first *int, last *int, orderBy *ent.CVERuleOrder, where *ent.CVERuleWhereInput) ([]*ent.CVERule, error) {
	return r.client.CVERule.Query().Filter(ctx, first, last, orderBy, where)
}

func (r *queryResolver) ComponentConnection(ctx context.Context, first *int, last *int, before *ent.Cursor, after *ent.Cursor, orderBy *ent.ComponentOrder, where *ent.ComponentWhereInput) (*ent.ComponentConnection, error) {
	return r.client.Component.Query().Paginate(ctx, after, first, before, last,
		ent.WithComponentOrder(orderBy),
		ent.WithComponentFilter(where.Filter),
	)
}

func (r *queryResolver) LicenseConnection(ctx context.Context, first *int, last *int, before *ent.Cursor, after *ent.Cursor, orderBy *ent.LicenseOrder, where *ent.LicenseWhereInput) (*ent.LicenseConnection, error) {
	return r.client.License.Query().Paginate(ctx, after, first, before, last,
		ent.WithLicenseOrder(orderBy),
		ent.WithLicenseFilter(where.Filter),
	)
}

func (r *queryResolver) ProjectConnection(ctx context.Context, first *int, last *int, before *ent.Cursor, after *ent.Cursor, orderBy *ent.ProjectOrder, where *ent.ProjectWhereInput) (*ent.ProjectConnection, error) {
	return r.client.Project.Query().Paginate(ctx, after, first, before, last,
		ent.WithProjectOrder(orderBy),
		ent.WithProjectFilter(where.Filter),
	)
}

func (r *queryResolver) VulnerabilityConnection(ctx context.Context, first *int, last *int, before *ent.Cursor, after *ent.Cursor, where *ent.VulnerabilityWhereInput) (*ent.VulnerabilityConnection, error) {
	return r.client.Vulnerability.Query().Paginate(ctx, after, first, before, last,
		ent.WithVulnerabilityFilter(where.Filter),
	)
}

func (r *queryResolver) Artifact(ctx context.Context, first *int, last *int, orderBy *ent.ArtifactOrder, where *ent.ArtifactWhereInput) ([]*ent.Artifact, error) {
	return r.client.Artifact.Query().Filter(ctx, first, last, orderBy, where)
}

func (r *queryResolver) Component(ctx context.Context, first *int, last *int, orderBy *ent.ComponentOrder, where *ent.ComponentWhereInput) ([]*ent.Component, error) {
	return r.client.Component.Query().Filter(ctx, first, last, orderBy, where)
}

func (r *queryResolver) LicenseUsageConnection(ctx context.Context, first *int, last *int, before *ent.Cursor, after *ent.Cursor, where *ent.LicenseUsageWhereInput) (*ent.LicenseUsageConnection, error) {
	return r.client.LicenseUsage.Query().Paginate(ctx, after, first, before, last,
		// ent.WithLicenseUsageOrder(orderBy),
		ent.WithLicenseUsageFilter(where.Filter),
	)
}

func (r *queryResolver) Project(ctx context.Context, first *int, last *int, orderBy *ent.ProjectOrder, where *ent.ProjectWhereInput) ([]*ent.Project, error) {
	return r.client.Project.Query().Filter(ctx, first, last, orderBy, where)
}

func (r *releaseResolver) Subreleases(ctx context.Context, obj *ent.Release, first *int, last *int, where *ent.ReleaseWhereInput, orderBy *ent.ReleaseOrder) ([]*ent.Release, error) {
	result, err := obj.Edges.SubreleasesOrErr()
	if ent.IsNotLoaded(err) {
		result, err = obj.QuerySubreleases().Filter(ctx, first, last, orderBy, where)
	}
	return result, err
}

func (r *releaseResolver) Dependencies(ctx context.Context, obj *ent.Release, first *int, last *int, where *ent.ReleaseWhereInput, orderBy *ent.ReleaseOrder) ([]*ent.Release, error) {
	result, err := obj.Edges.DependenciesOrErr()
	if ent.IsNotLoaded(err) {
		result, err = obj.QueryDependencies().Filter(ctx, first, last, orderBy, where)
	}
	return result, err
}

func (r *releaseResolver) Artifacts(ctx context.Context, obj *ent.Release, first *int, last *int, where *ent.ArtifactWhereInput, orderBy *ent.ArtifactOrder) ([]*ent.Artifact, error) {
	result, err := obj.Edges.ArtifactsOrErr()
	if ent.IsNotLoaded(err) {
		result, err = obj.QueryArtifacts().Filter(ctx, first, last, orderBy, where)
	}
	return result, err
}

func (r *releaseResolver) Checks(ctx context.Context, obj *ent.Release, first *int, last *int, where *ent.ReleaseCheckWhereInput, orderBy *ent.ReleaseCheckOrder) ([]*ent.ReleaseCheck, error) {
	result, err := obj.Edges.ChecksOrErr()
	if ent.IsNotLoaded(err) {
		result, err = obj.QueryChecks().Filter(ctx, first, last, orderBy, where)
	}
	return result, err
}

func (r *releaseResolver) Log(ctx context.Context, obj *ent.Release, first *int, last *int, where *ent.ReleaseEntryWhereInput, orderBy *ent.ReleaseEntryOrder) ([]*ent.ReleaseEntry, error) {
	result, err := obj.Edges.LogOrErr()
	if ent.IsNotLoaded(err) {
		result, err = obj.QueryLog().Filter(ctx, first, last, orderBy, where)
	}
	return result, err
}

func (r *releaseResolver) CodeScans(ctx context.Context, obj *ent.Release, first *int, last *int, where *ent.CodeScanWhereInput, orderBy *ent.CodeScanOrder) ([]*ent.CodeScan, error) {
	result, err := obj.Edges.CodeScansOrErr()
	if ent.IsNotLoaded(err) {
		result, err = obj.QueryCodeScans().Filter(ctx, first, last, orderBy, where)
	}
	return result, err
}

func (r *releaseResolver) CveScans(ctx context.Context, obj *ent.Release, first *int, last *int, where *ent.CVEScanWhereInput, orderBy *ent.CVEScanOrder) ([]*ent.CVEScan, error) {
	result, err := obj.Edges.CveScansOrErr()
	if ent.IsNotLoaded(err) {
		result, err = obj.QueryCveScans().Filter(ctx, first, last, orderBy, where)
	}
	return result, err
}

func (r *releaseResolver) LicenseScans(ctx context.Context, obj *ent.Release, first *int, last *int, where *ent.LicenseScanWhereInput, orderBy *ent.LicenseScanOrder) ([]*ent.LicenseScan, error) {
	result, err := obj.Edges.LicenseScansOrErr()
	if ent.IsNotLoaded(err) {
		result, err = obj.QueryLicenseScans().Filter(ctx, first, last, orderBy, where)
	}
	return result, err
}

func (r *releaseResolver) TestRuns(ctx context.Context, obj *ent.Release, first *int, last *int, where *ent.TestRunWhereInput, orderBy *ent.TestRunOrder) ([]*ent.TestRun, error) {
	result, err := obj.Edges.TestRunsOrErr()
	if ent.IsNotLoaded(err) {
		result, err = obj.QueryTestRuns().Filter(ctx, first, last, orderBy, where)
	}
	return result, err
}

func (r *releaseResolver) Components(ctx context.Context, obj *ent.Release, first *int, last *int, where *ent.ComponentWhereInput, orderBy *ent.ComponentOrder) ([]*ent.Component, error) {
	result, err := obj.Edges.ComponentsOrErr()
	if ent.IsNotLoaded(err) {
		result, err = obj.QueryComponents().Filter(ctx, first, last, orderBy, where)
	}
	return result, err
}

func (r *repoResolver) Commits(ctx context.Context, obj *ent.Repo, first *int, last *int, where *ent.GitCommitWhereInput, orderBy *ent.GitCommitOrder) ([]*ent.GitCommit, error) {
	result, err := obj.Edges.CommitsOrErr()
	if ent.IsNotLoaded(err) {
		result, err = obj.QueryCommits().Filter(ctx, first, last, orderBy, where)
		// result, err = obj.QueryCommits().All(ctx)
	}
	return result, err
}

func (r *repoResolver) CveRules(ctx context.Context, obj *ent.Repo, first *int, last *int, where *ent.CVERuleWhereInput, orderBy *ent.CVERuleOrder) ([]*ent.CVERule, error) {
	result, err := obj.Edges.CveRulesOrErr()
	if ent.IsNotLoaded(err) {
		result, err = obj.QueryCveRules().Filter(ctx, first, last, orderBy, where)
	}
	return result, err
}

func (r *testRunResolver) Tests(ctx context.Context, obj *ent.TestRun, first *int, last *int, where *ent.TestCaseWhereInput, orderBy *ent.TestCaseOrder) ([]*ent.TestCase, error) {
	result, err := obj.Edges.TestsOrErr()
	if ent.IsNotLoaded(err) {
		result, err = obj.QueryTests().Filter(ctx, first, last, orderBy, where)
	}
	return result, err
}

// CVE returns CVEResolver implementation.
func (r *Resolver) CVE() CVEResolver { return &cVEResolver{r} }

// CVERule returns CVERuleResolver implementation.
func (r *Resolver) CVERule() CVERuleResolver { return &cVERuleResolver{r} }

// CVEScan returns CVEScanResolver implementation.
func (r *Resolver) CVEScan() CVEScanResolver { return &cVEScanResolver{r} }

// CWE returns CWEResolver implementation.
func (r *Resolver) CWE() CWEResolver { return &cWEResolver{r} }

// CodeIssue returns CodeIssueResolver implementation.
func (r *Resolver) CodeIssue() CodeIssueResolver { return &codeIssueResolver{r} }

// CodeScan returns CodeScanResolver implementation.
func (r *Resolver) CodeScan() CodeScanResolver { return &codeScanResolver{r} }

// Component returns ComponentResolver implementation.
func (r *Resolver) Component() ComponentResolver { return &componentResolver{r} }

// License returns LicenseResolver implementation.
func (r *Resolver) License() LicenseResolver { return &licenseResolver{r} }

// LicenseScan returns LicenseScanResolver implementation.
func (r *Resolver) LicenseScan() LicenseScanResolver { return &licenseScanResolver{r} }

// Project returns ProjectResolver implementation.
func (r *Resolver) Project() ProjectResolver { return &projectResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

// Release returns ReleaseResolver implementation.
func (r *Resolver) Release() ReleaseResolver { return &releaseResolver{r} }

// Repo returns RepoResolver implementation.
func (r *Resolver) Repo() RepoResolver { return &repoResolver{r} }

// TestRun returns TestRunResolver implementation.
func (r *Resolver) TestRun() TestRunResolver { return &testRunResolver{r} }

type cVEResolver struct{ *Resolver }
type cVERuleResolver struct{ *Resolver }
type cVEScanResolver struct{ *Resolver }
type cWEResolver struct{ *Resolver }
type codeIssueResolver struct{ *Resolver }
type codeScanResolver struct{ *Resolver }
type componentResolver struct{ *Resolver }
type licenseResolver struct{ *Resolver }
type licenseScanResolver struct{ *Resolver }
type projectResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type releaseResolver struct{ *Resolver }
type repoResolver struct{ *Resolver }
type testRunResolver struct{ *Resolver }
