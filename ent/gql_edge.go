// Code generated by entc, DO NOT EDIT.

package ent

import "context"

func (a *Adapter) Owner(ctx context.Context) (*Organization, error) {
	result, err := a.Edges.OwnerOrErr()
	if IsNotLoaded(err) {
		result, err = a.QueryOwner().Only(ctx)
	}
	return result, err
}

func (a *Artifact) Release(ctx context.Context) (*Release, error) {
	result, err := a.Edges.ReleaseOrErr()
	if IsNotLoaded(err) {
		result, err = a.QueryRelease().Only(ctx)
	}
	return result, MaskNotFound(err)
}

func (a *Artifact) Entry(ctx context.Context) (*ReleaseEntry, error) {
	result, err := a.Edges.EntryOrErr()
	if IsNotLoaded(err) {
		result, err = a.QueryEntry().Only(ctx)
	}
	return result, MaskNotFound(err)
}

func (ci *CodeIssue) Scan(ctx context.Context) (*CodeScan, error) {
	result, err := ci.Edges.ScanOrErr()
	if IsNotLoaded(err) {
		result, err = ci.QueryScan().Only(ctx)
	}
	return result, err
}

func (cs *CodeScan) Release(ctx context.Context) (*Release, error) {
	result, err := cs.Edges.ReleaseOrErr()
	if IsNotLoaded(err) {
		result, err = cs.QueryRelease().Only(ctx)
	}
	return result, err
}

func (cs *CodeScan) Entry(ctx context.Context) (*ReleaseEntry, error) {
	result, err := cs.Edges.EntryOrErr()
	if IsNotLoaded(err) {
		result, err = cs.QueryEntry().Only(ctx)
	}
	return result, MaskNotFound(err)
}

func (cs *CodeScan) Issues(ctx context.Context) ([]*CodeIssue, error) {
	result, err := cs.Edges.IssuesOrErr()
	if IsNotLoaded(err) {
		result, err = cs.QueryIssues().All(ctx)
	}
	return result, err
}

func (cs *CodeScan) Vulnerabilities(ctx context.Context) ([]*ReleaseVulnerability, error) {
	result, err := cs.Edges.VulnerabilitiesOrErr()
	if IsNotLoaded(err) {
		result, err = cs.QueryVulnerabilities().All(ctx)
	}
	return result, err
}

func (cs *CodeScan) Components(ctx context.Context) ([]*ReleaseComponent, error) {
	result, err := cs.Edges.ComponentsOrErr()
	if IsNotLoaded(err) {
		result, err = cs.QueryComponents().All(ctx)
	}
	return result, err
}

func (c *Component) Owner(ctx context.Context) (*Organization, error) {
	result, err := c.Edges.OwnerOrErr()
	if IsNotLoaded(err) {
		result, err = c.QueryOwner().Only(ctx)
	}
	return result, err
}

func (c *Component) Vulnerabilities(ctx context.Context) ([]*Vulnerability, error) {
	result, err := c.Edges.VulnerabilitiesOrErr()
	if IsNotLoaded(err) {
		result, err = c.QueryVulnerabilities().All(ctx)
	}
	return result, err
}

func (c *Component) Licenses(ctx context.Context) ([]*License, error) {
	result, err := c.Edges.LicensesOrErr()
	if IsNotLoaded(err) {
		result, err = c.QueryLicenses().All(ctx)
	}
	return result, err
}

func (c *Component) Uses(ctx context.Context) ([]*ReleaseComponent, error) {
	result, err := c.Edges.UsesOrErr()
	if IsNotLoaded(err) {
		result, err = c.QueryUses().All(ctx)
	}
	return result, err
}

func (e *Event) Release(ctx context.Context) (*Release, error) {
	result, err := e.Edges.ReleaseOrErr()
	if IsNotLoaded(err) {
		result, err = e.QueryRelease().Only(ctx)
	}
	return result, MaskNotFound(err)
}

func (e *Event) Repo(ctx context.Context) (*Repo, error) {
	result, err := e.Edges.RepoOrErr()
	if IsNotLoaded(err) {
		result, err = e.QueryRepo().Only(ctx)
	}
	return result, MaskNotFound(err)
}

func (e *Event) Project(ctx context.Context) (*Project, error) {
	result, err := e.Edges.ProjectOrErr()
	if IsNotLoaded(err) {
		result, err = e.QueryProject().Only(ctx)
	}
	return result, MaskNotFound(err)
}

func (gc *GitCommit) Repo(ctx context.Context) (*Repo, error) {
	result, err := gc.Edges.RepoOrErr()
	if IsNotLoaded(err) {
		result, err = gc.QueryRepo().Only(ctx)
	}
	return result, err
}

func (gc *GitCommit) Release(ctx context.Context) (*Release, error) {
	result, err := gc.Edges.ReleaseOrErr()
	if IsNotLoaded(err) {
		result, err = gc.QueryRelease().Only(ctx)
	}
	return result, MaskNotFound(err)
}

func (l *License) Components(ctx context.Context) ([]*Component, error) {
	result, err := l.Edges.ComponentsOrErr()
	if IsNotLoaded(err) {
		result, err = l.QueryComponents().All(ctx)
	}
	return result, err
}

func (l *License) Uses(ctx context.Context) ([]*LicenseUse, error) {
	result, err := l.Edges.UsesOrErr()
	if IsNotLoaded(err) {
		result, err = l.QueryUses().All(ctx)
	}
	return result, err
}

func (lu *LicenseUse) License(ctx context.Context) (*License, error) {
	result, err := lu.Edges.LicenseOrErr()
	if IsNotLoaded(err) {
		result, err = lu.QueryLicense().Only(ctx)
	}
	return result, err
}

func (o *Organization) Projects(ctx context.Context) ([]*Project, error) {
	result, err := o.Edges.ProjectsOrErr()
	if IsNotLoaded(err) {
		result, err = o.QueryProjects().All(ctx)
	}
	return result, err
}

func (o *Organization) Repos(ctx context.Context) ([]*Repo, error) {
	result, err := o.Edges.ReposOrErr()
	if IsNotLoaded(err) {
		result, err = o.QueryRepos().All(ctx)
	}
	return result, err
}

func (pr *Project) Owner(ctx context.Context) (*Organization, error) {
	result, err := pr.Edges.OwnerOrErr()
	if IsNotLoaded(err) {
		result, err = pr.QueryOwner().Only(ctx)
	}
	return result, err
}

func (pr *Project) Repos(ctx context.Context) ([]*Repo, error) {
	result, err := pr.Edges.ReposOrErr()
	if IsNotLoaded(err) {
		result, err = pr.QueryRepos().All(ctx)
	}
	return result, err
}

func (pr *Project) VulnerabilityReviews(ctx context.Context) ([]*VulnerabilityReview, error) {
	result, err := pr.Edges.VulnerabilityReviewsOrErr()
	if IsNotLoaded(err) {
		result, err = pr.QueryVulnerabilityReviews().All(ctx)
	}
	return result, err
}

func (pr *Project) Policies(ctx context.Context) ([]*ReleasePolicy, error) {
	result, err := pr.Edges.PoliciesOrErr()
	if IsNotLoaded(err) {
		result, err = pr.QueryPolicies().All(ctx)
	}
	return result, err
}

func (r *Release) Subreleases(ctx context.Context) ([]*Release, error) {
	result, err := r.Edges.SubreleasesOrErr()
	if IsNotLoaded(err) {
		result, err = r.QuerySubreleases().All(ctx)
	}
	return result, err
}

func (r *Release) Dependencies(ctx context.Context) ([]*Release, error) {
	result, err := r.Edges.DependenciesOrErr()
	if IsNotLoaded(err) {
		result, err = r.QueryDependencies().All(ctx)
	}
	return result, err
}

func (r *Release) Commit(ctx context.Context) (*GitCommit, error) {
	result, err := r.Edges.CommitOrErr()
	if IsNotLoaded(err) {
		result, err = r.QueryCommit().Only(ctx)
	}
	return result, err
}

func (r *Release) HeadOf(ctx context.Context) (*Repo, error) {
	result, err := r.Edges.HeadOfOrErr()
	if IsNotLoaded(err) {
		result, err = r.QueryHeadOf().Only(ctx)
	}
	return result, MaskNotFound(err)
}

func (r *Release) Log(ctx context.Context) ([]*ReleaseEntry, error) {
	result, err := r.Edges.LogOrErr()
	if IsNotLoaded(err) {
		result, err = r.QueryLog().All(ctx)
	}
	return result, err
}

func (r *Release) Violations(ctx context.Context) ([]*ReleasePolicyViolation, error) {
	result, err := r.Edges.ViolationsOrErr()
	if IsNotLoaded(err) {
		result, err = r.QueryViolations().All(ctx)
	}
	return result, err
}

func (r *Release) Artifacts(ctx context.Context) ([]*Artifact, error) {
	result, err := r.Edges.ArtifactsOrErr()
	if IsNotLoaded(err) {
		result, err = r.QueryArtifacts().All(ctx)
	}
	return result, err
}

func (r *Release) Components(ctx context.Context) ([]*ReleaseComponent, error) {
	result, err := r.Edges.ComponentsOrErr()
	if IsNotLoaded(err) {
		result, err = r.QueryComponents().All(ctx)
	}
	return result, err
}

func (r *Release) Vulnerabilities(ctx context.Context) ([]*ReleaseVulnerability, error) {
	result, err := r.Edges.VulnerabilitiesOrErr()
	if IsNotLoaded(err) {
		result, err = r.QueryVulnerabilities().All(ctx)
	}
	return result, err
}

func (r *Release) CodeScans(ctx context.Context) ([]*CodeScan, error) {
	result, err := r.Edges.CodeScansOrErr()
	if IsNotLoaded(err) {
		result, err = r.QueryCodeScans().All(ctx)
	}
	return result, err
}

func (r *Release) TestRuns(ctx context.Context) ([]*TestRun, error) {
	result, err := r.Edges.TestRunsOrErr()
	if IsNotLoaded(err) {
		result, err = r.QueryTestRuns().All(ctx)
	}
	return result, err
}

func (r *Release) VulnerabilityReviews(ctx context.Context) ([]*VulnerabilityReview, error) {
	result, err := r.Edges.VulnerabilityReviewsOrErr()
	if IsNotLoaded(err) {
		result, err = r.QueryVulnerabilityReviews().All(ctx)
	}
	return result, err
}

func (rc *ReleaseComponent) Release(ctx context.Context) (*Release, error) {
	result, err := rc.Edges.ReleaseOrErr()
	if IsNotLoaded(err) {
		result, err = rc.QueryRelease().Only(ctx)
	}
	return result, err
}

func (rc *ReleaseComponent) Scans(ctx context.Context) ([]*CodeScan, error) {
	result, err := rc.Edges.ScansOrErr()
	if IsNotLoaded(err) {
		result, err = rc.QueryScans().All(ctx)
	}
	return result, err
}

func (rc *ReleaseComponent) Component(ctx context.Context) (*Component, error) {
	result, err := rc.Edges.ComponentOrErr()
	if IsNotLoaded(err) {
		result, err = rc.QueryComponent().Only(ctx)
	}
	return result, err
}

func (rc *ReleaseComponent) Vulnerabilities(ctx context.Context) ([]*ReleaseVulnerability, error) {
	result, err := rc.Edges.VulnerabilitiesOrErr()
	if IsNotLoaded(err) {
		result, err = rc.QueryVulnerabilities().All(ctx)
	}
	return result, err
}

func (re *ReleaseEntry) Artifact(ctx context.Context) (*Artifact, error) {
	result, err := re.Edges.ArtifactOrErr()
	if IsNotLoaded(err) {
		result, err = re.QueryArtifact().Only(ctx)
	}
	return result, MaskNotFound(err)
}

func (re *ReleaseEntry) CodeScan(ctx context.Context) (*CodeScan, error) {
	result, err := re.Edges.CodeScanOrErr()
	if IsNotLoaded(err) {
		result, err = re.QueryCodeScan().Only(ctx)
	}
	return result, MaskNotFound(err)
}

func (re *ReleaseEntry) TestRun(ctx context.Context) (*TestRun, error) {
	result, err := re.Edges.TestRunOrErr()
	if IsNotLoaded(err) {
		result, err = re.QueryTestRun().Only(ctx)
	}
	return result, MaskNotFound(err)
}

func (re *ReleaseEntry) Release(ctx context.Context) (*Release, error) {
	result, err := re.Edges.ReleaseOrErr()
	if IsNotLoaded(err) {
		result, err = re.QueryRelease().Only(ctx)
	}
	return result, err
}

func (rl *ReleaseLicense) License(ctx context.Context) (*License, error) {
	result, err := rl.Edges.LicenseOrErr()
	if IsNotLoaded(err) {
		result, err = rl.QueryLicense().Only(ctx)
	}
	return result, err
}

func (rl *ReleaseLicense) Component(ctx context.Context) (*ReleaseComponent, error) {
	result, err := rl.Edges.ComponentOrErr()
	if IsNotLoaded(err) {
		result, err = rl.QueryComponent().Only(ctx)
	}
	return result, MaskNotFound(err)
}

func (rl *ReleaseLicense) Release(ctx context.Context) (*Release, error) {
	result, err := rl.Edges.ReleaseOrErr()
	if IsNotLoaded(err) {
		result, err = rl.QueryRelease().Only(ctx)
	}
	return result, err
}

func (rl *ReleaseLicense) Scans(ctx context.Context) ([]*CodeScan, error) {
	result, err := rl.Edges.ScansOrErr()
	if IsNotLoaded(err) {
		result, err = rl.QueryScans().All(ctx)
	}
	return result, err
}

func (rp *ReleasePolicy) Owner(ctx context.Context) (*Organization, error) {
	result, err := rp.Edges.OwnerOrErr()
	if IsNotLoaded(err) {
		result, err = rp.QueryOwner().Only(ctx)
	}
	return result, err
}

func (rp *ReleasePolicy) Projects(ctx context.Context) ([]*Project, error) {
	result, err := rp.Edges.ProjectsOrErr()
	if IsNotLoaded(err) {
		result, err = rp.QueryProjects().All(ctx)
	}
	return result, err
}

func (rp *ReleasePolicy) Repos(ctx context.Context) ([]*Repo, error) {
	result, err := rp.Edges.ReposOrErr()
	if IsNotLoaded(err) {
		result, err = rp.QueryRepos().All(ctx)
	}
	return result, err
}

func (rp *ReleasePolicy) Violations(ctx context.Context) ([]*ReleasePolicyViolation, error) {
	result, err := rp.Edges.ViolationsOrErr()
	if IsNotLoaded(err) {
		result, err = rp.QueryViolations().All(ctx)
	}
	return result, err
}

func (rpv *ReleasePolicyViolation) Policy(ctx context.Context) (*ReleasePolicy, error) {
	result, err := rpv.Edges.PolicyOrErr()
	if IsNotLoaded(err) {
		result, err = rpv.QueryPolicy().Only(ctx)
	}
	return result, err
}

func (rpv *ReleasePolicyViolation) Release(ctx context.Context) (*Release, error) {
	result, err := rpv.Edges.ReleaseOrErr()
	if IsNotLoaded(err) {
		result, err = rpv.QueryRelease().Only(ctx)
	}
	return result, err
}

func (rv *ReleaseVulnerability) Vulnerability(ctx context.Context) (*Vulnerability, error) {
	result, err := rv.Edges.VulnerabilityOrErr()
	if IsNotLoaded(err) {
		result, err = rv.QueryVulnerability().Only(ctx)
	}
	return result, err
}

func (rv *ReleaseVulnerability) Component(ctx context.Context) (*ReleaseComponent, error) {
	result, err := rv.Edges.ComponentOrErr()
	if IsNotLoaded(err) {
		result, err = rv.QueryComponent().Only(ctx)
	}
	return result, MaskNotFound(err)
}

func (rv *ReleaseVulnerability) Release(ctx context.Context) (*Release, error) {
	result, err := rv.Edges.ReleaseOrErr()
	if IsNotLoaded(err) {
		result, err = rv.QueryRelease().Only(ctx)
	}
	return result, err
}

func (rv *ReleaseVulnerability) Reviews(ctx context.Context) ([]*VulnerabilityReview, error) {
	result, err := rv.Edges.ReviewsOrErr()
	if IsNotLoaded(err) {
		result, err = rv.QueryReviews().All(ctx)
	}
	return result, err
}

func (rv *ReleaseVulnerability) Scan(ctx context.Context) (*CodeScan, error) {
	result, err := rv.Edges.ScanOrErr()
	if IsNotLoaded(err) {
		result, err = rv.QueryScan().Only(ctx)
	}
	return result, MaskNotFound(err)
}

func (r *Repo) Owner(ctx context.Context) (*Organization, error) {
	result, err := r.Edges.OwnerOrErr()
	if IsNotLoaded(err) {
		result, err = r.QueryOwner().Only(ctx)
	}
	return result, err
}

func (r *Repo) Project(ctx context.Context) (*Project, error) {
	result, err := r.Edges.ProjectOrErr()
	if IsNotLoaded(err) {
		result, err = r.QueryProject().Only(ctx)
	}
	return result, err
}

func (r *Repo) Head(ctx context.Context) (*Release, error) {
	result, err := r.Edges.HeadOrErr()
	if IsNotLoaded(err) {
		result, err = r.QueryHead().Only(ctx)
	}
	return result, MaskNotFound(err)
}

func (r *Repo) Commits(ctx context.Context) ([]*GitCommit, error) {
	result, err := r.Edges.CommitsOrErr()
	if IsNotLoaded(err) {
		result, err = r.QueryCommits().All(ctx)
	}
	return result, err
}

func (r *Repo) VulnerabilityReviews(ctx context.Context) ([]*VulnerabilityReview, error) {
	result, err := r.Edges.VulnerabilityReviewsOrErr()
	if IsNotLoaded(err) {
		result, err = r.QueryVulnerabilityReviews().All(ctx)
	}
	return result, err
}

func (r *Repo) Policies(ctx context.Context) ([]*ReleasePolicy, error) {
	result, err := r.Edges.PoliciesOrErr()
	if IsNotLoaded(err) {
		result, err = r.QueryPolicies().All(ctx)
	}
	return result, err
}

func (tc *TestCase) Run(ctx context.Context) (*TestRun, error) {
	result, err := tc.Edges.RunOrErr()
	if IsNotLoaded(err) {
		result, err = tc.QueryRun().Only(ctx)
	}
	return result, err
}

func (tr *TestRun) Release(ctx context.Context) (*Release, error) {
	result, err := tr.Edges.ReleaseOrErr()
	if IsNotLoaded(err) {
		result, err = tr.QueryRelease().Only(ctx)
	}
	return result, err
}

func (tr *TestRun) Entry(ctx context.Context) (*ReleaseEntry, error) {
	result, err := tr.Edges.EntryOrErr()
	if IsNotLoaded(err) {
		result, err = tr.QueryEntry().Only(ctx)
	}
	return result, MaskNotFound(err)
}

func (tr *TestRun) Tests(ctx context.Context) ([]*TestCase, error) {
	result, err := tr.Edges.TestsOrErr()
	if IsNotLoaded(err) {
		result, err = tr.QueryTests().All(ctx)
	}
	return result, err
}

func (v *Vulnerability) Owner(ctx context.Context) (*Organization, error) {
	result, err := v.Edges.OwnerOrErr()
	if IsNotLoaded(err) {
		result, err = v.QueryOwner().Only(ctx)
	}
	return result, err
}

func (v *Vulnerability) Components(ctx context.Context) ([]*Component, error) {
	result, err := v.Edges.ComponentsOrErr()
	if IsNotLoaded(err) {
		result, err = v.QueryComponents().All(ctx)
	}
	return result, err
}

func (v *Vulnerability) Reviews(ctx context.Context) ([]*VulnerabilityReview, error) {
	result, err := v.Edges.ReviewsOrErr()
	if IsNotLoaded(err) {
		result, err = v.QueryReviews().All(ctx)
	}
	return result, err
}

func (v *Vulnerability) Instances(ctx context.Context) ([]*ReleaseVulnerability, error) {
	result, err := v.Edges.InstancesOrErr()
	if IsNotLoaded(err) {
		result, err = v.QueryInstances().All(ctx)
	}
	return result, err
}

func (vr *VulnerabilityReview) Vulnerability(ctx context.Context) (*Vulnerability, error) {
	result, err := vr.Edges.VulnerabilityOrErr()
	if IsNotLoaded(err) {
		result, err = vr.QueryVulnerability().Only(ctx)
	}
	return result, err
}

func (vr *VulnerabilityReview) Projects(ctx context.Context) ([]*Project, error) {
	result, err := vr.Edges.ProjectsOrErr()
	if IsNotLoaded(err) {
		result, err = vr.QueryProjects().All(ctx)
	}
	return result, err
}

func (vr *VulnerabilityReview) Repos(ctx context.Context) ([]*Repo, error) {
	result, err := vr.Edges.ReposOrErr()
	if IsNotLoaded(err) {
		result, err = vr.QueryRepos().All(ctx)
	}
	return result, err
}

func (vr *VulnerabilityReview) Releases(ctx context.Context) ([]*Release, error) {
	result, err := vr.Edges.ReleasesOrErr()
	if IsNotLoaded(err) {
		result, err = vr.QueryReleases().All(ctx)
	}
	return result, err
}

func (vr *VulnerabilityReview) Instances(ctx context.Context) ([]*ReleaseVulnerability, error) {
	result, err := vr.Edges.InstancesOrErr()
	if IsNotLoaded(err) {
		result, err = vr.QueryInstances().All(ctx)
	}
	return result, err
}
