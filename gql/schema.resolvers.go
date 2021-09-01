package gql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/valocode/bubbly/ent"
)

func (r *artifactResolver) Metadata(ctx context.Context, obj *ent.Artifact) (map[string]interface{}, error) {
	return obj.Metadata, nil
}

func (r *codeIssueResolver) Metadata(ctx context.Context, obj *ent.CodeIssue) (map[string]interface{}, error) {
	return obj.Metadata, nil
}

func (r *codeScanResolver) Metadata(ctx context.Context, obj *ent.CodeScan) (map[string]interface{}, error) {
	return obj.Metadata, nil
}

func (r *codeScanResolver) Issues(ctx context.Context, obj *ent.CodeScan, first *int, last *int, where *ent.CodeIssueWhereInput, orderBy *ent.CodeIssueOrder) ([]*ent.CodeIssue, error) {
	result, err := obj.Edges.IssuesOrErr()
	if ent.IsNotLoaded(err) {
		result, err = obj.QueryIssues().Filter(ctx, first, last, orderBy, where)
	}
	return result, err
}

func (r *codeScanResolver) Vulnerabilities(ctx context.Context, obj *ent.CodeScan, first *int, last *int, where *ent.ReleaseVulnerabilityWhereInput) ([]*ent.ReleaseVulnerability, error) {
	result, err := obj.Edges.VulnerabilitiesOrErr()
	if ent.IsNotLoaded(err) {
		result, err = obj.QueryVulnerabilities().Filter(ctx, first, last, nil, where)
	}
	return result, err
}

func (r *codeScanResolver) Components(ctx context.Context, obj *ent.CodeScan, first *int, last *int, where *ent.ReleaseComponentWhereInput) ([]*ent.ReleaseComponent, error) {
	result, err := obj.Edges.ComponentsOrErr()
	if ent.IsNotLoaded(err) {
		result, err = obj.QueryComponents().Filter(ctx, first, last, nil, where)
	}
	return result, err
}

func (r *componentResolver) Metadata(ctx context.Context, obj *ent.Component) (map[string]interface{}, error) {
	return obj.Metadata, nil
}

func (r *componentResolver) Vulnerabilities(ctx context.Context, obj *ent.Component, first *int, last *int, where *ent.VulnerabilityWhereInput, orderBy *ent.VulnerabilityOrder) ([]*ent.Vulnerability, error) {
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

func (r *componentResolver) Uses(ctx context.Context, obj *ent.Component, first *int, last *int, where *ent.ReleaseComponentWhereInput) ([]*ent.ReleaseComponent, error) {
	result, err := obj.Edges.UsesOrErr()
	if ent.IsNotLoaded(err) {
		result, err = obj.QueryUses().Filter(ctx, first, last, nil, where)
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

func (r *licenseResolver) Uses(ctx context.Context, obj *ent.License, first *int, last *int, where *ent.LicenseUseWhereInput) ([]*ent.LicenseUse, error) {
	return obj.QueryUses().Filter(ctx, first, last, nil, where)
}

func (r *organizationResolver) Projects(ctx context.Context, obj *ent.Organization, first *int, last *int, where *ent.ProjectWhereInput, orderBy *ent.ProjectOrder) ([]*ent.Project, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *organizationResolver) Repos(ctx context.Context, obj *ent.Organization, first *int, last *int, where *ent.RepoWhereInput, orderBy *ent.RepoOrder) ([]*ent.Repo, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *projectResolver) Repos(ctx context.Context, obj *ent.Project, first *int, last *int, where *ent.RepoWhereInput, orderBy *ent.RepoOrder) ([]*ent.Repo, error) {
	result, err := obj.Edges.ReposOrErr()
	if ent.IsNotLoaded(err) {
		result, err = obj.QueryRepos().Filter(ctx, first, last, orderBy, where)
	}
	return result, err
}

func (r *projectResolver) VulnerabilityReviews(ctx context.Context, obj *ent.Project, first *int, last *int, where *ent.VulnerabilityReviewWhereInput, orderBy *ent.VulnerabilityReviewOrder) ([]*ent.VulnerabilityReview, error) {
	result, err := obj.Edges.VulnerabilityReviewsOrErr()
	if ent.IsNotLoaded(err) {
		result, err = obj.QueryVulnerabilityReviews().Filter(ctx, first, last, orderBy, where)
	}
	return result, err
}

func (r *projectResolver) Policies(ctx context.Context, obj *ent.Project, first *int, last *int, where *ent.ReleasePolicyWhereInput, orderBy *ent.ReleasePolicyOrder) ([]*ent.ReleasePolicy, error) {
	result, err := obj.Edges.PoliciesOrErr()
	if ent.IsNotLoaded(err) {
		result, err = obj.QueryPolicies().Filter(ctx, first, last, orderBy, where)
	}
	return result, err
}

func (r *queryResolver) ArtifactConnection(ctx context.Context, first *int, last *int, before *ent.Cursor, after *ent.Cursor, orderBy *ent.ArtifactOrder, where *ent.ArtifactWhereInput) (*ent.ArtifactConnection, error) {
	return r.client.Artifact.Query().Paginate(ctx, after, first, before, last,
		ent.WithArtifactOrder(orderBy),
		ent.WithArtifactFilter(where.Filter),
	)
}

func (r *queryResolver) Policy(ctx context.Context, first *int, last *int, orderBy *ent.ReleasePolicyOrder, where *ent.ReleasePolicyWhereInput) ([]*ent.ReleasePolicy, error) {
	return r.client.ReleasePolicy.Query().Filter(ctx, first, last, orderBy, where)
}

func (r *queryResolver) Artifact(ctx context.Context, first *int, last *int, orderBy *ent.ArtifactOrder, where *ent.ArtifactWhereInput) ([]*ent.Artifact, error) {
	return r.client.Artifact.Query().Filter(ctx, first, last, orderBy, where)
}

func (r *queryResolver) CodeIssueConnection(ctx context.Context, first *int, last *int, before *ent.Cursor, after *ent.Cursor, orderBy *ent.CodeIssueOrder, where *ent.CodeIssueWhereInput) (*ent.CodeIssueConnection, error) {
	return r.client.CodeIssue.Query().Paginate(ctx, after, first, before, last,
		ent.WithCodeIssueOrder(orderBy),
		ent.WithCodeIssueFilter(where.Filter),
	)
}

func (r *queryResolver) ComponentConnection(ctx context.Context, first *int, last *int, before *ent.Cursor, after *ent.Cursor, orderBy *ent.ComponentOrder, where *ent.ComponentWhereInput) (*ent.ComponentConnection, error) {
	return r.client.Component.Query().Paginate(ctx, after, first, before, last,
		ent.WithComponentOrder(orderBy),
		ent.WithComponentFilter(where.Filter),
	)
}

func (r *queryResolver) TestCaseConnection(ctx context.Context, first *int, last *int, before *ent.Cursor, after *ent.Cursor, orderBy *ent.TestCaseOrder, where *ent.TestCaseWhereInput) (*ent.TestCaseConnection, error) {
	return r.client.TestCase.Query().Paginate(ctx, after, first, before, last,
		ent.WithTestCaseOrder(orderBy),
		ent.WithTestCaseFilter(where.Filter),
	)
}

func (r *queryResolver) VulnerabilityConnection(ctx context.Context, first *int, last *int, before *ent.Cursor, after *ent.Cursor, orderBy *ent.VulnerabilityOrder, where *ent.VulnerabilityWhereInput) (*ent.VulnerabilityConnection, error) {
	return r.client.Vulnerability.Query().Paginate(ctx, after, first, before, last,
		ent.WithVulnerabilityFilter(where.Filter),
	)
}

func (r *queryResolver) Component(ctx context.Context, first *int, last *int, orderBy *ent.ComponentOrder, where *ent.ComponentWhereInput) ([]*ent.Component, error) {
	return r.client.Component.Query().Filter(ctx, first, last, orderBy, where)
}

func (r *queryResolver) CommitConnection(ctx context.Context, first *int, last *int, before *ent.Cursor, after *ent.Cursor, orderBy *ent.GitCommitOrder, where *ent.GitCommitWhereInput) (*ent.GitCommitConnection, error) {
	return r.client.GitCommit.Query().Paginate(ctx, after, first, before, last,
		ent.WithGitCommitOrder(orderBy),
		ent.WithGitCommitFilter(where.Filter),
	)
}

func (r *queryResolver) TestCase(ctx context.Context, first *int, last *int, orderBy *ent.TestCaseOrder, where *ent.TestCaseWhereInput) ([]*ent.TestCase, error) {
	return r.client.TestCase.Query().Filter(ctx, first, last, orderBy, where)
}

func (r *queryResolver) LicenseConnection(ctx context.Context, first *int, last *int, before *ent.Cursor, after *ent.Cursor, orderBy *ent.LicenseOrder, where *ent.LicenseWhereInput) (*ent.LicenseConnection, error) {
	return r.client.License.Query().Paginate(ctx, after, first, before, last,
		ent.WithLicenseOrder(orderBy),
		ent.WithLicenseFilter(where.Filter),
	)
}

func (r *queryResolver) ReleaseEntryConnection(ctx context.Context, first *int, last *int, before *ent.Cursor, after *ent.Cursor, orderBy *ent.ReleaseEntryOrder, where *ent.ReleaseEntryWhereInput) (*ent.ReleaseEntryConnection, error) {
	return r.client.ReleaseEntry.Query().Paginate(ctx, after, first, before, last,
		ent.WithReleaseEntryOrder(orderBy),
		ent.WithReleaseEntryFilter(where.Filter),
	)
}

func (r *queryResolver) TestRun(ctx context.Context, first *int, last *int, orderBy *ent.TestRunOrder, where *ent.TestRunWhereInput) ([]*ent.TestRun, error) {
	return r.client.TestRun.Query().Filter(ctx, first, last, orderBy, where)
}

func (r *queryResolver) VulnerabilityReview(ctx context.Context, first *int, last *int, orderBy *ent.VulnerabilityReviewOrder, where *ent.VulnerabilityReviewWhereInput) ([]*ent.VulnerabilityReview, error) {
	return r.client.VulnerabilityReview.Query().Filter(ctx, first, last, orderBy, where)
}

func (r *queryResolver) License(ctx context.Context, first *int, last *int, orderBy *ent.LicenseOrder, where *ent.LicenseWhereInput) ([]*ent.License, error) {
	return r.client.License.Query().Filter(ctx, first, last, orderBy, where)
}

func (r *queryResolver) ReleaseVulnerability(ctx context.Context, first *int, last *int, where *ent.ReleaseVulnerabilityWhereInput) ([]*ent.ReleaseVulnerability, error) {
	return r.client.ReleaseVulnerability.Query().Filter(ctx, first, last, nil, where)
}

func (r *queryResolver) Repo(ctx context.Context, first *int, last *int, orderBy *ent.RepoOrder, where *ent.RepoWhereInput) ([]*ent.Repo, error) {
	return r.client.Repo.Query().Filter(ctx, first, last, orderBy, where)
}

func (r *queryResolver) TestRunConnection(ctx context.Context, first *int, last *int, before *ent.Cursor, after *ent.Cursor, orderBy *ent.TestRunOrder, where *ent.TestRunWhereInput) (*ent.TestRunConnection, error) {
	return r.client.TestRun.Query().Paginate(ctx, after, first, before, last,
		ent.WithTestRunOrder(orderBy),
		ent.WithTestRunFilter(where.Filter),
	)
}

func (r *queryResolver) LicenseUseConnection(ctx context.Context, first *int, last *int, before *ent.Cursor, after *ent.Cursor, where *ent.LicenseUseWhereInput) (*ent.LicenseUseConnection, error) {
	return r.client.LicenseUse.Query().Paginate(ctx, after, first, before, last,
		// ent.WithLicenseUseOrder(orderBy),
		ent.WithLicenseUseFilter(where.Filter),
	)
}

func (r *queryResolver) LicenseUse(ctx context.Context, first *int, last *int, where *ent.LicenseUseWhereInput) ([]*ent.LicenseUse, error) {
	return r.client.LicenseUse.Query().Filter(ctx, first, last, nil, where)
}

func (r *queryResolver) Project(ctx context.Context, first *int, last *int, orderBy *ent.ProjectOrder, where *ent.ProjectWhereInput) ([]*ent.Project, error) {
	return r.client.Project.Query().Filter(ctx, first, last, orderBy, where)
}

func (r *queryResolver) Release(ctx context.Context, first *int, last *int, orderBy *ent.ReleaseOrder, where *ent.ReleaseWhereInput) ([]*ent.Release, error) {
	return r.client.Release.Query().Filter(ctx, first, last, orderBy, where)
}

func (r *queryResolver) RepoConnection(ctx context.Context, first *int, last *int, before *ent.Cursor, after *ent.Cursor, orderBy *ent.RepoOrder, where *ent.RepoWhereInput) (*ent.RepoConnection, error) {
	return r.client.Repo.Query().Paginate(ctx, after, first, before, last,
		ent.WithRepoOrder(orderBy),
		ent.WithRepoFilter(where.Filter),
	)
}

func (r *queryResolver) VulnerabilityReviewConnection(ctx context.Context, first *int, last *int, before *ent.Cursor, after *ent.Cursor, orderBy *ent.VulnerabilityReviewOrder, where *ent.VulnerabilityReviewWhereInput) (*ent.VulnerabilityReviewConnection, error) {
	return r.client.VulnerabilityReview.Query().Paginate(ctx, after, first, before, last,
		ent.WithVulnerabilityReviewOrder(orderBy),
		ent.WithVulnerabilityReviewFilter(where.Filter),
	)
}

func (r *queryResolver) CodeScanConnection(ctx context.Context, first *int, last *int, before *ent.Cursor, after *ent.Cursor, orderBy *ent.CodeScanOrder, where *ent.CodeScanWhereInput) (*ent.CodeScanConnection, error) {
	return r.client.CodeScan.Query().Paginate(ctx, after, first, before, last,
		ent.WithCodeScanOrder(orderBy),
		ent.WithCodeScanFilter(where.Filter),
	)
}

func (r *queryResolver) Commit(ctx context.Context, first *int, last *int, orderBy *ent.GitCommitOrder, where *ent.GitCommitWhereInput) ([]*ent.GitCommit, error) {
	return r.client.GitCommit.Query().Filter(ctx, first, last, orderBy, where)
}

func (r *queryResolver) CodeScan(ctx context.Context, first *int, last *int, orderBy *ent.CodeScanOrder, where *ent.CodeScanWhereInput) ([]*ent.CodeScan, error) {
	return r.client.CodeScan.Query().Filter(ctx, first, last, orderBy, where)
}

func (r *queryResolver) ReleaseComponentConnection(ctx context.Context, first *int, last *int, before *ent.Cursor, after *ent.Cursor, where *ent.ReleaseComponentWhereInput) (*ent.ReleaseComponentConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) ReleaseEntry(ctx context.Context, first *int, last *int, orderBy *ent.ReleaseEntryOrder, where *ent.ReleaseEntryWhereInput) ([]*ent.ReleaseEntry, error) {
	return r.client.ReleaseEntry.Query().Filter(ctx, first, last, orderBy, where)
}

func (r *queryResolver) CodeIssue(ctx context.Context, first *int, last *int, orderBy *ent.CodeIssueOrder, where *ent.CodeIssueWhereInput) ([]*ent.CodeIssue, error) {
	return r.client.CodeIssue.Query().Filter(ctx, first, last, orderBy, where)
}

func (r *queryResolver) ProjectConnection(ctx context.Context, first *int, last *int, before *ent.Cursor, after *ent.Cursor, orderBy *ent.ProjectOrder, where *ent.ProjectWhereInput) (*ent.ProjectConnection, error) {
	return r.client.Project.Query().Paginate(ctx, after, first, before, last,
		ent.WithProjectOrder(orderBy),
		ent.WithProjectFilter(where.Filter),
	)
}

func (r *queryResolver) ReleaseConnection(ctx context.Context, first *int, last *int, before *ent.Cursor, after *ent.Cursor, orderBy *ent.ReleaseOrder, where *ent.ReleaseWhereInput) (*ent.ReleaseConnection, error) {
	return r.client.Release.Query().Paginate(ctx, after, first, before, last,
		ent.WithReleaseOrder(orderBy),
		ent.WithReleaseFilter(where.Filter),
	)
}

func (r *queryResolver) ReleaseComponent(ctx context.Context, first *int, last *int, where *ent.ReleaseComponentWhereInput) ([]*ent.ReleaseComponent, error) {
	return r.client.ReleaseComponent.Query().Filter(ctx, first, last, nil, where)
}

func (r *queryResolver) ReleaseVulnerabilityConnection(ctx context.Context, first *int, last *int, before *ent.Cursor, after *ent.Cursor, where *ent.ReleaseVulnerabilityWhereInput) (*ent.ReleaseVulnerabilityConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Vulnerability(ctx context.Context, first *int, last *int, orderBy *ent.VulnerabilityOrder, where *ent.VulnerabilityWhereInput) ([]*ent.Vulnerability, error) {
	return r.client.Vulnerability.Query().Filter(ctx, first, last, nil, where)
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

func (r *releaseResolver) Log(ctx context.Context, obj *ent.Release, first *int, last *int, where *ent.ReleaseEntryWhereInput, orderBy *ent.ReleaseEntryOrder) ([]*ent.ReleaseEntry, error) {
	result, err := obj.Edges.LogOrErr()
	if ent.IsNotLoaded(err) {
		result, err = obj.QueryLog().Filter(ctx, first, last, orderBy, where)
	}
	return result, err
}

func (r *releaseResolver) Violations(ctx context.Context, obj *ent.Release, first *int, last *int, where *ent.ReleasePolicyViolationWhereInput) ([]*ent.ReleasePolicyViolation, error) {
	result, err := obj.Edges.ViolationsOrErr()
	if ent.IsNotLoaded(err) {
		result, err = obj.QueryViolations().Filter(ctx, first, last, nil, where)
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

func (r *releaseResolver) Components(ctx context.Context, obj *ent.Release, first *int, last *int, where *ent.ReleaseComponentWhereInput) ([]*ent.ReleaseComponent, error) {
	result, err := obj.Edges.ComponentsOrErr()
	if ent.IsNotLoaded(err) {
		result, err = obj.QueryComponents().Filter(ctx, first, last, nil, where)
	}
	return result, err
}

func (r *releaseResolver) Vulnerabilities(ctx context.Context, obj *ent.Release, first *int, last *int, where *ent.ReleaseVulnerabilityWhereInput) ([]*ent.ReleaseVulnerability, error) {
	result, err := obj.Edges.VulnerabilitiesOrErr()
	if ent.IsNotLoaded(err) {
		result, err = obj.QueryVulnerabilities().Filter(ctx, first, last, nil, where)
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

func (r *releaseResolver) TestRuns(ctx context.Context, obj *ent.Release, first *int, last *int, where *ent.TestRunWhereInput, orderBy *ent.TestRunOrder) ([]*ent.TestRun, error) {
	result, err := obj.Edges.TestRunsOrErr()
	if ent.IsNotLoaded(err) {
		result, err = obj.QueryTestRuns().Filter(ctx, first, last, orderBy, where)
	}
	return result, err
}

func (r *releaseResolver) VulnerabilityReviews(ctx context.Context, obj *ent.Release, first *int, last *int, where *ent.VulnerabilityReviewWhereInput, orderBy *ent.VulnerabilityReviewOrder) ([]*ent.VulnerabilityReview, error) {
	result, err := obj.Edges.VulnerabilityReviewsOrErr()
	if ent.IsNotLoaded(err) {
		result, err = obj.QueryVulnerabilityReviews().Filter(ctx, first, last, orderBy, where)
	}
	return result, err
}

func (r *releaseComponentResolver) Scans(ctx context.Context, obj *ent.ReleaseComponent, first *int, last *int, where *ent.CodeScanWhereInput, orderBy *ent.CodeScanOrder) ([]*ent.CodeScan, error) {
	result, err := obj.Edges.ScansOrErr()
	if ent.IsNotLoaded(err) {
		result, err = obj.QueryScans().Filter(ctx, first, last, orderBy, where)
	}
	return result, err
}

func (r *releaseComponentResolver) Vulnerabilities(ctx context.Context, obj *ent.ReleaseComponent, first *int, last *int, where *ent.ReleaseVulnerabilityWhereInput) ([]*ent.ReleaseVulnerability, error) {
	result, err := obj.Edges.VulnerabilitiesOrErr()
	if ent.IsNotLoaded(err) {
		result, err = obj.QueryVulnerabilities().Filter(ctx, first, last, nil, where)
	}
	return result, err
}

func (r *releaseLicenseResolver) Scans(ctx context.Context, obj *ent.ReleaseLicense, first *int, last *int, where *ent.CodeScanWhereInput, orderBy *ent.CodeScanOrder) ([]*ent.CodeScan, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *releasePolicyResolver) Projects(ctx context.Context, obj *ent.ReleasePolicy, first *int, last *int, where *ent.ProjectWhereInput, orderBy *ent.ProjectOrder) ([]*ent.Project, error) {
	result, err := obj.Edges.ProjectsOrErr()
	if ent.IsNotLoaded(err) {
		result, err = obj.QueryProjects().Filter(ctx, first, last, orderBy, where)
	}
	return result, err
}

func (r *releasePolicyResolver) Repos(ctx context.Context, obj *ent.ReleasePolicy, first *int, last *int, where *ent.RepoWhereInput, orderBy *ent.RepoOrder) ([]*ent.Repo, error) {
	result, err := obj.Edges.ReposOrErr()
	if ent.IsNotLoaded(err) {
		result, err = obj.QueryRepos().Filter(ctx, first, last, orderBy, where)
	}
	return result, err
}

func (r *releasePolicyResolver) Violations(ctx context.Context, obj *ent.ReleasePolicy, first *int, last *int, where *ent.ReleasePolicyViolationWhereInput) ([]*ent.ReleasePolicyViolation, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *releaseVulnerabilityResolver) Reviews(ctx context.Context, obj *ent.ReleaseVulnerability, first *int, last *int, where *ent.VulnerabilityReviewWhereInput, orderBy *ent.VulnerabilityReviewOrder) ([]*ent.VulnerabilityReview, error) {
	result, err := obj.Edges.ReviewsOrErr()
	if ent.IsNotLoaded(err) {
		result, err = obj.QueryReviews().Filter(ctx, first, last, orderBy, where)
	}
	return result, err
}

func (r *repoResolver) Commits(ctx context.Context, obj *ent.Repo, first *int, last *int, where *ent.GitCommitWhereInput, orderBy *ent.GitCommitOrder) ([]*ent.GitCommit, error) {
	result, err := obj.Edges.CommitsOrErr()
	if ent.IsNotLoaded(err) {
		result, err = obj.QueryCommits().Filter(ctx, first, last, orderBy, where)
	}
	return result, err
}

func (r *repoResolver) VulnerabilityReviews(ctx context.Context, obj *ent.Repo, first *int, last *int, where *ent.VulnerabilityReviewWhereInput, orderBy *ent.VulnerabilityReviewOrder) ([]*ent.VulnerabilityReview, error) {
	result, err := obj.Edges.VulnerabilityReviewsOrErr()
	if ent.IsNotLoaded(err) {
		result, err = obj.QueryVulnerabilityReviews().Filter(ctx, first, last, orderBy, where)
	}
	return result, err
}

func (r *repoResolver) Policies(ctx context.Context, obj *ent.Repo, first *int, last *int, where *ent.ReleasePolicyWhereInput, orderBy *ent.ReleasePolicyOrder) ([]*ent.ReleasePolicy, error) {
	result, err := obj.Edges.PoliciesOrErr()
	if ent.IsNotLoaded(err) {
		result, err = obj.QueryPolicies().Filter(ctx, first, last, orderBy, where)
	}
	return result, err
}

func (r *testCaseResolver) Metadata(ctx context.Context, obj *ent.TestCase) (map[string]interface{}, error) {
	return obj.Metadata, nil
}

func (r *testRunResolver) Metadata(ctx context.Context, obj *ent.TestRun) (map[string]interface{}, error) {
	return obj.Metadata, nil
}

func (r *testRunResolver) Tests(ctx context.Context, obj *ent.TestRun, first *int, last *int, where *ent.TestCaseWhereInput, orderBy *ent.TestCaseOrder) ([]*ent.TestCase, error) {
	result, err := obj.Edges.TestsOrErr()
	if ent.IsNotLoaded(err) {
		result, err = obj.QueryTests().Filter(ctx, first, last, orderBy, where)
	}
	return result, err
}

func (r *vulnerabilityResolver) Metadata(ctx context.Context, obj *ent.Vulnerability) (map[string]interface{}, error) {
	return obj.Metadata, nil
}

func (r *vulnerabilityResolver) Components(ctx context.Context, obj *ent.Vulnerability, first *int, last *int, where *ent.ComponentWhereInput, orderBy *ent.ComponentOrder) ([]*ent.Component, error) {
	result, err := obj.Edges.ComponentsOrErr()
	if ent.IsNotLoaded(err) {
		result, err = obj.QueryComponents().Filter(ctx, first, last, nil, where)
	}
	return result, err
}

func (r *vulnerabilityResolver) Reviews(ctx context.Context, obj *ent.Vulnerability, first *int, last *int, where *ent.VulnerabilityReviewWhereInput, orderBy *ent.VulnerabilityReviewOrder) ([]*ent.VulnerabilityReview, error) {
	result, err := obj.Edges.ReviewsOrErr()
	if ent.IsNotLoaded(err) {
		result, err = obj.QueryReviews().Filter(ctx, first, last, orderBy, where)
	}
	return result, err
}

func (r *vulnerabilityResolver) Instances(ctx context.Context, obj *ent.Vulnerability, first *int, last *int, where *ent.ReleaseVulnerabilityWhereInput) ([]*ent.ReleaseVulnerability, error) {
	result, err := obj.Edges.InstancesOrErr()
	if ent.IsNotLoaded(err) {
		result, err = obj.QueryInstances().Filter(ctx, first, last, nil, where)
	}
	return result, err
}

func (r *vulnerabilityReviewResolver) Projects(ctx context.Context, obj *ent.VulnerabilityReview, first *int, last *int, where *ent.ProjectWhereInput, orderBy *ent.ProjectOrder) ([]*ent.Project, error) {
	result, err := obj.Edges.ProjectsOrErr()
	if ent.IsNotLoaded(err) {
		result, err = obj.QueryProjects().Filter(ctx, first, last, nil, where)
	}
	return result, err
}

func (r *vulnerabilityReviewResolver) Repos(ctx context.Context, obj *ent.VulnerabilityReview, first *int, last *int, where *ent.RepoWhereInput, orderBy *ent.RepoOrder) ([]*ent.Repo, error) {
	result, err := obj.Edges.ReposOrErr()
	if ent.IsNotLoaded(err) {
		result, err = obj.QueryRepos().Filter(ctx, first, last, nil, where)
	}
	return result, err
}

func (r *vulnerabilityReviewResolver) Releases(ctx context.Context, obj *ent.VulnerabilityReview, first *int, last *int, where *ent.ReleaseWhereInput, orderBy *ent.ReleaseOrder) ([]*ent.Release, error) {
	result, err := obj.Edges.ReleasesOrErr()
	if ent.IsNotLoaded(err) {
		result, err = obj.QueryReleases().Filter(ctx, first, last, nil, where)
	}
	return result, err
}

func (r *vulnerabilityReviewResolver) Instances(ctx context.Context, obj *ent.VulnerabilityReview, first *int, last *int, where *ent.ReleaseVulnerabilityWhereInput) ([]*ent.ReleaseVulnerability, error) {
	result, err := obj.Edges.InstancesOrErr()
	if ent.IsNotLoaded(err) {
		result, err = obj.QueryInstances().Filter(ctx, first, last, nil, where)
	}
	return result, err
}

// Artifact returns ArtifactResolver implementation.
func (r *Resolver) Artifact() ArtifactResolver { return &artifactResolver{r} }

// CodeIssue returns CodeIssueResolver implementation.
func (r *Resolver) CodeIssue() CodeIssueResolver { return &codeIssueResolver{r} }

// CodeScan returns CodeScanResolver implementation.
func (r *Resolver) CodeScan() CodeScanResolver { return &codeScanResolver{r} }

// Component returns ComponentResolver implementation.
func (r *Resolver) Component() ComponentResolver { return &componentResolver{r} }

// License returns LicenseResolver implementation.
func (r *Resolver) License() LicenseResolver { return &licenseResolver{r} }

// Organization returns OrganizationResolver implementation.
func (r *Resolver) Organization() OrganizationResolver { return &organizationResolver{r} }

// Project returns ProjectResolver implementation.
func (r *Resolver) Project() ProjectResolver { return &projectResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

// Release returns ReleaseResolver implementation.
func (r *Resolver) Release() ReleaseResolver { return &releaseResolver{r} }

// ReleaseComponent returns ReleaseComponentResolver implementation.
func (r *Resolver) ReleaseComponent() ReleaseComponentResolver { return &releaseComponentResolver{r} }

// ReleaseLicense returns ReleaseLicenseResolver implementation.
func (r *Resolver) ReleaseLicense() ReleaseLicenseResolver { return &releaseLicenseResolver{r} }

// ReleasePolicy returns ReleasePolicyResolver implementation.
func (r *Resolver) ReleasePolicy() ReleasePolicyResolver { return &releasePolicyResolver{r} }

// ReleaseVulnerability returns ReleaseVulnerabilityResolver implementation.
func (r *Resolver) ReleaseVulnerability() ReleaseVulnerabilityResolver {
	return &releaseVulnerabilityResolver{r}
}

// Repo returns RepoResolver implementation.
func (r *Resolver) Repo() RepoResolver { return &repoResolver{r} }

// TestCase returns TestCaseResolver implementation.
func (r *Resolver) TestCase() TestCaseResolver { return &testCaseResolver{r} }

// TestRun returns TestRunResolver implementation.
func (r *Resolver) TestRun() TestRunResolver { return &testRunResolver{r} }

// Vulnerability returns VulnerabilityResolver implementation.
func (r *Resolver) Vulnerability() VulnerabilityResolver { return &vulnerabilityResolver{r} }

// VulnerabilityReview returns VulnerabilityReviewResolver implementation.
func (r *Resolver) VulnerabilityReview() VulnerabilityReviewResolver {
	return &vulnerabilityReviewResolver{r}
}

type artifactResolver struct{ *Resolver }
type codeIssueResolver struct{ *Resolver }
type codeScanResolver struct{ *Resolver }
type componentResolver struct{ *Resolver }
type licenseResolver struct{ *Resolver }
type organizationResolver struct{ *Resolver }
type projectResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type releaseResolver struct{ *Resolver }
type releaseComponentResolver struct{ *Resolver }
type releaseLicenseResolver struct{ *Resolver }
type releasePolicyResolver struct{ *Resolver }
type releaseVulnerabilityResolver struct{ *Resolver }
type repoResolver struct{ *Resolver }
type testCaseResolver struct{ *Resolver }
type testRunResolver struct{ *Resolver }
type vulnerabilityResolver struct{ *Resolver }
type vulnerabilityReviewResolver struct{ *Resolver }
