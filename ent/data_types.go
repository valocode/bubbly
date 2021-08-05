// Code generated by entc, DO NOT EDIT.

package ent

import (
	"time"

	"github.com/valocode/bubbly/ent/artifact"
	"github.com/valocode/bubbly/ent/codeissue"
	"github.com/valocode/bubbly/ent/cve"
	"github.com/valocode/bubbly/ent/release"
	"github.com/valocode/bubbly/ent/releasecheck"
	"github.com/valocode/bubbly/ent/releaseentry"
	"github.com/zclconf/go-cty/cty"
)

func NewArtifactNode() *ArtifactNode {
	return &ArtifactNode{
		DataNode: &DataNode{
			Name: "artifact",
		},
	}
}

type ArtifactNode struct {
	*DataNode
}

func (a *ArtifactNode) Node() *DataNode {
	return a.DataNode
}
func (a *ArtifactNode) SetName(v string) *ArtifactNode {
	a.AddField("name", cty.StringVal(v))
	return a
}
func (a *ArtifactNode) SetSha256(v string) *ArtifactNode {
	a.AddField("sha256", cty.StringVal(v))
	return a
}
func (a *ArtifactNode) SetType(v artifact.Type) *ArtifactNode {
	a.AddField("type", cty.StringVal(v.String()))
	return a
}

func (a *ArtifactNode) SetRelease(input *ReleaseNode) *ArtifactNode {
	node := input.Node()

	// Add an edge to the input node
	a.AddEdge("release", node)
	// Create the invse edge also
	node.AddInverseEdge("artifacts", a.Node())
	return a
}

func (a *ArtifactNode) SetEntry(input *ReleaseEntryNode) *ArtifactNode {
	node := input.Node()

	// Add an edge to the input node
	a.AddEdge("entry", node)
	// Create the invse edge also
	node.AddInverseEdge("artifact", a.Node())
	return a
}

func (a *ArtifactNode) SetOperation(op NodeOperation) *ArtifactNode {
	a.Operation = op
	return a
}

func NewCVENode() *CVENode {
	return &CVENode{
		DataNode: &DataNode{
			Name: "cve",
		},
	}
}

type CVENode struct {
	*DataNode
}

func (c *CVENode) Node() *DataNode {
	return c.DataNode
}
func (c *CVENode) SetCveID(v string) *CVENode {
	c.AddField("cve_id", cty.StringVal(v))
	return c
}
func (c *CVENode) SetDescription(v string) *CVENode {
	c.AddField("description", cty.StringVal(v))
	return c
}
func (c *CVENode) SetSeverityScore(v float64) *CVENode {
	c.AddField("severity_score", cty.NumberFloatVal(v))
	return c
}
func (c *CVENode) SetSeverity(v cve.Severity) *CVENode {
	c.AddField("severity", cty.StringVal(v.String()))
	return c
}
func (c *CVENode) SetPublishedData(v time.Time) *CVENode {
	c.AddField("published_data", cty.StringVal(v.Format(time.RFC3339)))
	return c
}
func (c *CVENode) SetModifiedData(v time.Time) *CVENode {
	c.AddField("modified_data", cty.StringVal(v.Format(time.RFC3339)))
	return c
}

func (c *CVENode) AddFound(inputs ...*VulnerabilityNode) *CVENode {
	for _, input := range inputs {
		node := input.Node()

		// Add an edge to the input node
		c.AddEdge("found", node)
		// Create the invse edge also
		node.AddInverseEdge("cve", c.Node())
	}
	return c
}

func (c *CVENode) AddRules(inputs ...*CVERuleNode) *CVENode {
	for _, input := range inputs {
		node := input.Node()

		// Add an edge to the input node
		c.AddEdge("rules", node)
		// Create the invse edge also
		node.AddInverseEdge("cve", c.Node())
	}
	return c
}

func (c *CVENode) SetOperation(op NodeOperation) *CVENode {
	c.Operation = op
	return c
}

func NewCVERuleNode() *CVERuleNode {
	return &CVERuleNode{
		DataNode: &DataNode{
			Name: "cve_rule",
		},
	}
}

type CVERuleNode struct {
	*DataNode
}

func (cr *CVERuleNode) Node() *DataNode {
	return cr.DataNode
}
func (cr *CVERuleNode) SetName(v string) *CVERuleNode {
	cr.AddField("name", cty.StringVal(v))
	return cr
}

func (cr *CVERuleNode) SetCve(input *CVENode) *CVERuleNode {
	node := input.Node()

	// Add an edge to the input node
	cr.AddInverseEdge("cve", node)
	// Create the invse edge also
	node.AddEdge("rules", cr.Node())
	return cr
}

func (cr *CVERuleNode) AddProject(inputs ...*ProjectNode) *CVERuleNode {
	for _, input := range inputs {
		node := input.Node()

		// Add an edge to the input node
		cr.AddEdge("project", node)
		// Create the invse edge also
		node.AddInverseEdge("cve_rules", cr.Node())
	}
	return cr
}

func (cr *CVERuleNode) AddRepo(inputs ...*RepoNode) *CVERuleNode {
	for _, input := range inputs {
		node := input.Node()

		// Add an edge to the input node
		cr.AddEdge("repo", node)
		// Create the invse edge also
		node.AddInverseEdge("cve_rules", cr.Node())
	}
	return cr
}

func (cr *CVERuleNode) SetOperation(op NodeOperation) *CVERuleNode {
	cr.Operation = op
	return cr
}

func NewCVEScanNode() *CVEScanNode {
	return &CVEScanNode{
		DataNode: &DataNode{
			Name: "cve_scan",
		},
	}
}

type CVEScanNode struct {
	*DataNode
}

func (cs *CVEScanNode) Node() *DataNode {
	return cs.DataNode
}
func (cs *CVEScanNode) SetTool(v string) *CVEScanNode {
	cs.AddField("tool", cty.StringVal(v))
	return cs
}

func (cs *CVEScanNode) SetRelease(input *ReleaseNode) *CVEScanNode {
	node := input.Node()

	// Add an edge to the input node
	cs.AddInverseEdge("release", node)
	// Create the invse edge also
	node.AddEdge("cve_scans", cs.Node())
	return cs
}

func (cs *CVEScanNode) SetEntry(input *ReleaseEntryNode) *CVEScanNode {
	node := input.Node()

	// Add an edge to the input node
	cs.AddEdge("entry", node)
	// Create the invse edge also
	node.AddInverseEdge("cve_scan", cs.Node())
	return cs
}

func (cs *CVEScanNode) AddVulnerabilities(inputs ...*VulnerabilityNode) *CVEScanNode {
	for _, input := range inputs {
		node := input.Node()

		// Add an edge to the input node
		cs.AddEdge("vulnerabilities", node)
		// Create the invse edge also
		node.AddInverseEdge("scan", cs.Node())
	}
	return cs
}

func (cs *CVEScanNode) SetOperation(op NodeOperation) *CVEScanNode {
	cs.Operation = op
	return cs
}

func NewCWENode() *CWENode {
	return &CWENode{
		DataNode: &DataNode{
			Name: "cwe",
		},
	}
}

type CWENode struct {
	*DataNode
}

func (c *CWENode) Node() *DataNode {
	return c.DataNode
}
func (c *CWENode) SetCweID(v string) *CWENode {
	c.AddField("cwe_id", cty.StringVal(v))
	return c
}
func (c *CWENode) SetDescription(v string) *CWENode {
	c.AddField("description", cty.StringVal(v))
	return c
}
func (c *CWENode) SetURL(v float64) *CWENode {
	c.AddField("url", cty.NumberFloatVal(v))
	return c
}

func (c *CWENode) AddIssues(inputs ...*CodeIssueNode) *CWENode {
	for _, input := range inputs {
		node := input.Node()

		// Add an edge to the input node
		c.AddEdge("issues", node)
		// Create the invse edge also
		node.AddInverseEdge("cwe", c.Node())
	}
	return c
}

func (c *CWENode) SetOperation(op NodeOperation) *CWENode {
	c.Operation = op
	return c
}

func NewCodeIssueNode() *CodeIssueNode {
	return &CodeIssueNode{
		DataNode: &DataNode{
			Name: "code_issue",
		},
	}
}

type CodeIssueNode struct {
	*DataNode
}

func (ci *CodeIssueNode) Node() *DataNode {
	return ci.DataNode
}
func (ci *CodeIssueNode) SetRuleID(v string) *CodeIssueNode {
	ci.AddField("rule_id", cty.StringVal(v))
	return ci
}
func (ci *CodeIssueNode) SetMessage(v string) *CodeIssueNode {
	ci.AddField("message", cty.StringVal(v))
	return ci
}
func (ci *CodeIssueNode) SetSeverity(v codeissue.Severity) *CodeIssueNode {
	ci.AddField("severity", cty.StringVal(v.String()))
	return ci
}
func (ci *CodeIssueNode) SetType(v codeissue.Type) *CodeIssueNode {
	ci.AddField("type", cty.StringVal(v.String()))
	return ci
}

func (ci *CodeIssueNode) AddCwe(inputs ...*CWENode) *CodeIssueNode {
	for _, input := range inputs {
		node := input.Node()

		// Add an edge to the input node
		ci.AddEdge("cwe", node)
		// Create the invse edge also
		node.AddInverseEdge("issues", ci.Node())
	}
	return ci
}

func (ci *CodeIssueNode) SetScan(input *CodeScanNode) *CodeIssueNode {
	node := input.Node()

	// Add an edge to the input node
	ci.AddInverseEdge("scan", node)
	// Create the invse edge also
	node.AddEdge("issues", ci.Node())
	return ci
}

func (ci *CodeIssueNode) SetOperation(op NodeOperation) *CodeIssueNode {
	ci.Operation = op
	return ci
}

func NewCodeScanNode() *CodeScanNode {
	return &CodeScanNode{
		DataNode: &DataNode{
			Name: "code_scan",
		},
	}
}

type CodeScanNode struct {
	*DataNode
}

func (cs *CodeScanNode) Node() *DataNode {
	return cs.DataNode
}
func (cs *CodeScanNode) SetTool(v string) *CodeScanNode {
	cs.AddField("tool", cty.StringVal(v))
	return cs
}

func (cs *CodeScanNode) SetRelease(input *ReleaseNode) *CodeScanNode {
	node := input.Node()

	// Add an edge to the input node
	cs.AddInverseEdge("release", node)
	// Create the invse edge also
	node.AddEdge("code_scans", cs.Node())
	return cs
}

func (cs *CodeScanNode) AddIssues(inputs ...*CodeIssueNode) *CodeScanNode {
	for _, input := range inputs {
		node := input.Node()

		// Add an edge to the input node
		cs.AddEdge("issues", node)
		// Create the invse edge also
		node.AddInverseEdge("scan", cs.Node())
	}
	return cs
}

func (cs *CodeScanNode) SetEntry(input *ReleaseEntryNode) *CodeScanNode {
	node := input.Node()

	// Add an edge to the input node
	cs.AddEdge("entry", node)
	// Create the invse edge also
	node.AddInverseEdge("code_scan", cs.Node())
	return cs
}

func (cs *CodeScanNode) SetOperation(op NodeOperation) *CodeScanNode {
	cs.Operation = op
	return cs
}

func NewComponentNode() *ComponentNode {
	return &ComponentNode{
		DataNode: &DataNode{
			Name: "component",
		},
	}
}

type ComponentNode struct {
	*DataNode
}

func (c *ComponentNode) Node() *DataNode {
	return c.DataNode
}
func (c *ComponentNode) SetName(v string) *ComponentNode {
	c.AddField("name", cty.StringVal(v))
	return c
}
func (c *ComponentNode) SetVendor(v string) *ComponentNode {
	c.AddField("vendor", cty.StringVal(v))
	return c
}
func (c *ComponentNode) SetVersion(v string) *ComponentNode {
	c.AddField("version", cty.StringVal(v))
	return c
}
func (c *ComponentNode) SetDescription(v string) *ComponentNode {
	c.AddField("description", cty.StringVal(v))
	return c
}
func (c *ComponentNode) SetURL(v string) *ComponentNode {
	c.AddField("url", cty.StringVal(v))
	return c
}

func (c *ComponentNode) AddVulnerabilities(inputs ...*VulnerabilityNode) *ComponentNode {
	for _, input := range inputs {
		node := input.Node()

		// Add an edge to the input node
		c.AddEdge("vulnerabilities", node)
		// Create the invse edge also
		node.AddInverseEdge("component", c.Node())
	}
	return c
}

func (c *ComponentNode) AddLicenses(inputs ...*LicenseNode) *ComponentNode {
	for _, input := range inputs {
		node := input.Node()

		// Add an edge to the input node
		c.AddEdge("licenses", node)
		// Create the invse edge also
		node.AddInverseEdge("components", c.Node())
	}
	return c
}

func (c *ComponentNode) AddRelease(inputs ...*ReleaseNode) *ComponentNode {
	for _, input := range inputs {
		node := input.Node()

		// Add an edge to the input node
		c.AddEdge("release", node)
		// Create the invse edge also
		node.AddInverseEdge("components", c.Node())
	}
	return c
}

func (c *ComponentNode) SetOperation(op NodeOperation) *ComponentNode {
	c.Operation = op
	return c
}

func NewGitCommitNode() *GitCommitNode {
	return &GitCommitNode{
		DataNode: &DataNode{
			Name: "commit",
		},
	}
}

type GitCommitNode struct {
	*DataNode
}

func (gc *GitCommitNode) Node() *DataNode {
	return gc.DataNode
}
func (gc *GitCommitNode) SetHash(v string) *GitCommitNode {
	gc.AddField("hash", cty.StringVal(v))
	return gc
}
func (gc *GitCommitNode) SetBranch(v string) *GitCommitNode {
	gc.AddField("branch", cty.StringVal(v))
	return gc
}
func (gc *GitCommitNode) SetTag(v string) *GitCommitNode {
	gc.AddField("tag", cty.StringVal(v))
	return gc
}
func (gc *GitCommitNode) SetTime(v time.Time) *GitCommitNode {
	gc.AddField("time", cty.StringVal(v.Format(time.RFC3339)))
	return gc
}

func (gc *GitCommitNode) SetRepo(input *RepoNode) *GitCommitNode {
	node := input.Node()

	// Add an edge to the input node
	gc.AddInverseEdge("repo", node)
	// Create the invse edge also
	node.AddEdge("commits", gc.Node())
	return gc
}

func (gc *GitCommitNode) SetRelease(input *ReleaseNode) *GitCommitNode {
	node := input.Node()

	// Add an edge to the input node
	gc.AddEdge("release", node)
	// Create the invse edge also
	node.AddInverseEdge("commit", gc.Node())
	return gc
}

func (gc *GitCommitNode) SetOperation(op NodeOperation) *GitCommitNode {
	gc.Operation = op
	return gc
}

func NewLicenseNode() *LicenseNode {
	return &LicenseNode{
		DataNode: &DataNode{
			Name: "license",
		},
	}
}

type LicenseNode struct {
	*DataNode
}

func (l *LicenseNode) Node() *DataNode {
	return l.DataNode
}
func (l *LicenseNode) SetSpdxID(v string) *LicenseNode {
	l.AddField("spdx_id", cty.StringVal(v))
	return l
}
func (l *LicenseNode) SetName(v string) *LicenseNode {
	l.AddField("name", cty.StringVal(v))
	return l
}
func (l *LicenseNode) SetReference(v string) *LicenseNode {
	l.AddField("reference", cty.StringVal(v))
	return l
}
func (l *LicenseNode) SetDetailsURL(v string) *LicenseNode {
	l.AddField("details_url", cty.StringVal(v))
	return l
}
func (l *LicenseNode) SetIsOsiApproved(v bool) *LicenseNode {
	l.AddField("is_osi_approved", cty.BoolVal(v))
	return l
}

func (l *LicenseNode) AddComponents(inputs ...*ComponentNode) *LicenseNode {
	for _, input := range inputs {
		node := input.Node()

		// Add an edge to the input node
		l.AddEdge("components", node)
		// Create the invse edge also
		node.AddInverseEdge("licenses", l.Node())
	}
	return l
}

func (l *LicenseNode) AddUsages(inputs ...*LicenseUsageNode) *LicenseNode {
	for _, input := range inputs {
		node := input.Node()

		// Add an edge to the input node
		l.AddEdge("usages", node)
		// Create the invse edge also
		node.AddInverseEdge("license", l.Node())
	}
	return l
}

func (l *LicenseNode) SetOperation(op NodeOperation) *LicenseNode {
	l.Operation = op
	return l
}

func NewLicenseScanNode() *LicenseScanNode {
	return &LicenseScanNode{
		DataNode: &DataNode{
			Name: "license_scan",
		},
	}
}

type LicenseScanNode struct {
	*DataNode
}

func (ls *LicenseScanNode) Node() *DataNode {
	return ls.DataNode
}
func (ls *LicenseScanNode) SetTool(v string) *LicenseScanNode {
	ls.AddField("tool", cty.StringVal(v))
	return ls
}

func (ls *LicenseScanNode) SetRelease(input *ReleaseNode) *LicenseScanNode {
	node := input.Node()

	// Add an edge to the input node
	ls.AddInverseEdge("release", node)
	// Create the invse edge also
	node.AddEdge("license_scans", ls.Node())
	return ls
}

func (ls *LicenseScanNode) SetEntry(input *ReleaseEntryNode) *LicenseScanNode {
	node := input.Node()

	// Add an edge to the input node
	ls.AddEdge("entry", node)
	// Create the invse edge also
	node.AddInverseEdge("license_scan", ls.Node())
	return ls
}

func (ls *LicenseScanNode) AddLicenses(inputs ...*LicenseUsageNode) *LicenseScanNode {
	for _, input := range inputs {
		node := input.Node()

		// Add an edge to the input node
		ls.AddEdge("licenses", node)
		// Create the invse edge also
		node.AddInverseEdge("scan", ls.Node())
	}
	return ls
}

func (ls *LicenseScanNode) SetOperation(op NodeOperation) *LicenseScanNode {
	ls.Operation = op
	return ls
}

func NewLicenseUsageNode() *LicenseUsageNode {
	return &LicenseUsageNode{
		DataNode: &DataNode{
			Name: "license_usage",
		},
	}
}

type LicenseUsageNode struct {
	*DataNode
}

func (lu *LicenseUsageNode) Node() *DataNode {
	return lu.DataNode
}

func (lu *LicenseUsageNode) SetLicense(input *LicenseNode) *LicenseUsageNode {
	node := input.Node()

	// Add an edge to the input node
	lu.AddInverseEdge("license", node)
	// Create the invse edge also
	node.AddEdge("usages", lu.Node())
	return lu
}

func (lu *LicenseUsageNode) SetScan(input *LicenseScanNode) *LicenseUsageNode {
	node := input.Node()

	// Add an edge to the input node
	lu.AddInverseEdge("scan", node)
	// Create the invse edge also
	node.AddEdge("licenses", lu.Node())
	return lu
}

func (lu *LicenseUsageNode) SetOperation(op NodeOperation) *LicenseUsageNode {
	lu.Operation = op
	return lu
}

func NewProjectNode() *ProjectNode {
	return &ProjectNode{
		DataNode: &DataNode{
			Name: "project",
		},
	}
}

type ProjectNode struct {
	*DataNode
}

func (pr *ProjectNode) Node() *DataNode {
	return pr.DataNode
}
func (pr *ProjectNode) SetName(v string) *ProjectNode {
	pr.AddField("name", cty.StringVal(v))
	return pr
}

func (pr *ProjectNode) AddRepos(inputs ...*RepoNode) *ProjectNode {
	for _, input := range inputs {
		node := input.Node()

		// Add an edge to the input node
		pr.AddEdge("repos", node)
		// Create the invse edge also
		node.AddInverseEdge("project", pr.Node())
	}
	return pr
}

func (pr *ProjectNode) AddReleases(inputs ...*ReleaseNode) *ProjectNode {
	for _, input := range inputs {
		node := input.Node()

		// Add an edge to the input node
		pr.AddEdge("releases", node)
		// Create the invse edge also
		node.AddInverseEdge("project", pr.Node())
	}
	return pr
}

func (pr *ProjectNode) AddCveRules(inputs ...*CVERuleNode) *ProjectNode {
	for _, input := range inputs {
		node := input.Node()

		// Add an edge to the input node
		pr.AddEdge("cve_rules", node)
		// Create the invse edge also
		node.AddInverseEdge("project", pr.Node())
	}
	return pr
}

func (pr *ProjectNode) SetOperation(op NodeOperation) *ProjectNode {
	pr.Operation = op
	return pr
}

func NewReleaseNode() *ReleaseNode {
	return &ReleaseNode{
		DataNode: &DataNode{
			Name: "release",
		},
	}
}

type ReleaseNode struct {
	*DataNode
}

func (r *ReleaseNode) Node() *DataNode {
	return r.DataNode
}
func (r *ReleaseNode) SetName(v string) *ReleaseNode {
	r.AddField("name", cty.StringVal(v))
	return r
}
func (r *ReleaseNode) SetVersion(v string) *ReleaseNode {
	r.AddField("version", cty.StringVal(v))
	return r
}
func (r *ReleaseNode) SetStatus(v release.Status) *ReleaseNode {
	r.AddField("status", cty.StringVal(v.String()))
	return r
}

func (r *ReleaseNode) AddSubreleases(inputs ...*ReleaseNode) *ReleaseNode {
	for _, input := range inputs {
		node := input.Node()

		// Add an edge to the input node
		r.AddEdge("subreleases", node)
		// Create the invse edge also
		node.AddInverseEdge("dependencies", r.Node())
	}
	return r
}

func (r *ReleaseNode) AddDependencies(inputs ...*ReleaseNode) *ReleaseNode {
	for _, input := range inputs {
		node := input.Node()

		// Add an edge to the input node
		r.AddEdge("dependencies", node)
		// Create the invse edge also
		node.AddInverseEdge("subreleases", r.Node())
	}
	return r
}

func (r *ReleaseNode) SetProject(input *ProjectNode) *ReleaseNode {
	node := input.Node()

	// Add an edge to the input node
	r.AddInverseEdge("project", node)
	// Create the invse edge also
	node.AddEdge("releases", r.Node())
	return r
}

func (r *ReleaseNode) SetCommit(input *GitCommitNode) *ReleaseNode {
	node := input.Node()

	// Add an edge to the input node
	r.AddInverseEdge("commit", node)
	// Create the invse edge also
	node.AddEdge("release", r.Node())
	return r
}

func (r *ReleaseNode) AddArtifacts(inputs ...*ArtifactNode) *ReleaseNode {
	for _, input := range inputs {
		node := input.Node()

		// Add an edge to the input node
		r.AddEdge("artifacts", node)
		// Create the invse edge also
		node.AddInverseEdge("release", r.Node())
	}
	return r
}

func (r *ReleaseNode) AddChecks(inputs ...*ReleaseCheckNode) *ReleaseNode {
	for _, input := range inputs {
		node := input.Node()

		// Add an edge to the input node
		r.AddEdge("checks", node)
		// Create the invse edge also
		node.AddInverseEdge("release", r.Node())
	}
	return r
}

func (r *ReleaseNode) AddLog(inputs ...*ReleaseEntryNode) *ReleaseNode {
	for _, input := range inputs {
		node := input.Node()

		// Add an edge to the input node
		r.AddEdge("log", node)
		// Create the invse edge also
		node.AddInverseEdge("release", r.Node())
	}
	return r
}

func (r *ReleaseNode) AddCodeScans(inputs ...*CodeScanNode) *ReleaseNode {
	for _, input := range inputs {
		node := input.Node()

		// Add an edge to the input node
		r.AddEdge("code_scans", node)
		// Create the invse edge also
		node.AddInverseEdge("release", r.Node())
	}
	return r
}

func (r *ReleaseNode) AddCveScans(inputs ...*CVEScanNode) *ReleaseNode {
	for _, input := range inputs {
		node := input.Node()

		// Add an edge to the input node
		r.AddEdge("cve_scans", node)
		// Create the invse edge also
		node.AddInverseEdge("release", r.Node())
	}
	return r
}

func (r *ReleaseNode) AddLicenseScans(inputs ...*LicenseScanNode) *ReleaseNode {
	for _, input := range inputs {
		node := input.Node()

		// Add an edge to the input node
		r.AddEdge("license_scans", node)
		// Create the invse edge also
		node.AddInverseEdge("release", r.Node())
	}
	return r
}

func (r *ReleaseNode) AddTestRuns(inputs ...*TestRunNode) *ReleaseNode {
	for _, input := range inputs {
		node := input.Node()

		// Add an edge to the input node
		r.AddEdge("test_runs", node)
		// Create the invse edge also
		node.AddInverseEdge("release", r.Node())
	}
	return r
}

func (r *ReleaseNode) AddComponents(inputs ...*ComponentNode) *ReleaseNode {
	for _, input := range inputs {
		node := input.Node()

		// Add an edge to the input node
		r.AddEdge("components", node)
		// Create the invse edge also
		node.AddInverseEdge("release", r.Node())
	}
	return r
}

func (r *ReleaseNode) SetOperation(op NodeOperation) *ReleaseNode {
	r.Operation = op
	return r
}

func NewReleaseCheckNode() *ReleaseCheckNode {
	return &ReleaseCheckNode{
		DataNode: &DataNode{
			Name: "release_check",
		},
	}
}

type ReleaseCheckNode struct {
	*DataNode
}

func (rc *ReleaseCheckNode) Node() *DataNode {
	return rc.DataNode
}
func (rc *ReleaseCheckNode) SetType(v releasecheck.Type) *ReleaseCheckNode {
	rc.AddField("type", cty.StringVal(v.String()))
	return rc
}

func (rc *ReleaseCheckNode) SetRelease(input *ReleaseNode) *ReleaseCheckNode {
	node := input.Node()

	// Add an edge to the input node
	rc.AddInverseEdge("release", node)
	// Create the invse edge also
	node.AddEdge("checks", rc.Node())
	return rc
}

func (rc *ReleaseCheckNode) SetOperation(op NodeOperation) *ReleaseCheckNode {
	rc.Operation = op
	return rc
}

func NewReleaseEntryNode() *ReleaseEntryNode {
	return &ReleaseEntryNode{
		DataNode: &DataNode{
			Name: "release_entry",
		},
	}
}

type ReleaseEntryNode struct {
	*DataNode
}

func (re *ReleaseEntryNode) Node() *DataNode {
	return re.DataNode
}
func (re *ReleaseEntryNode) SetType(v releaseentry.Type) *ReleaseEntryNode {
	re.AddField("type", cty.StringVal(v.String()))
	return re
}
func (re *ReleaseEntryNode) SetTime(v time.Time) *ReleaseEntryNode {
	re.AddField("time", cty.StringVal(v.Format(time.RFC3339)))
	return re
}

func (re *ReleaseEntryNode) SetArtifact(input *ArtifactNode) *ReleaseEntryNode {
	node := input.Node()

	// Add an edge to the input node
	re.AddEdge("artifact", node)
	// Create the invse edge also
	node.AddInverseEdge("entry", re.Node())
	return re
}

func (re *ReleaseEntryNode) SetCodeScan(input *CodeScanNode) *ReleaseEntryNode {
	node := input.Node()

	// Add an edge to the input node
	re.AddEdge("code_scan", node)
	// Create the invse edge also
	node.AddInverseEdge("entry", re.Node())
	return re
}

func (re *ReleaseEntryNode) SetTestRun(input *TestRunNode) *ReleaseEntryNode {
	node := input.Node()

	// Add an edge to the input node
	re.AddEdge("test_run", node)
	// Create the invse edge also
	node.AddInverseEdge("entry", re.Node())
	return re
}

func (re *ReleaseEntryNode) SetCveScan(input *CVEScanNode) *ReleaseEntryNode {
	node := input.Node()

	// Add an edge to the input node
	re.AddEdge("cve_scan", node)
	// Create the invse edge also
	node.AddInverseEdge("entry", re.Node())
	return re
}

func (re *ReleaseEntryNode) SetLicenseScan(input *LicenseScanNode) *ReleaseEntryNode {
	node := input.Node()

	// Add an edge to the input node
	re.AddEdge("license_scan", node)
	// Create the invse edge also
	node.AddInverseEdge("entry", re.Node())
	return re
}

func (re *ReleaseEntryNode) SetRelease(input *ReleaseNode) *ReleaseEntryNode {
	node := input.Node()

	// Add an edge to the input node
	re.AddInverseEdge("release", node)
	// Create the invse edge also
	node.AddEdge("log", re.Node())
	return re
}

func (re *ReleaseEntryNode) SetOperation(op NodeOperation) *ReleaseEntryNode {
	re.Operation = op
	return re
}

func NewRepoNode() *RepoNode {
	return &RepoNode{
		DataNode: &DataNode{
			Name: "repo",
		},
	}
}

type RepoNode struct {
	*DataNode
}

func (r *RepoNode) Node() *DataNode {
	return r.DataNode
}
func (r *RepoNode) SetName(v string) *RepoNode {
	r.AddField("name", cty.StringVal(v))
	return r
}

func (r *RepoNode) SetProject(input *ProjectNode) *RepoNode {
	node := input.Node()

	// Add an edge to the input node
	r.AddEdge("project", node)
	// Create the invse edge also
	node.AddInverseEdge("repos", r.Node())
	return r
}

func (r *RepoNode) AddCommits(inputs ...*GitCommitNode) *RepoNode {
	for _, input := range inputs {
		node := input.Node()

		// Add an edge to the input node
		r.AddEdge("commits", node)
		// Create the invse edge also
		node.AddInverseEdge("repo", r.Node())
	}
	return r
}

func (r *RepoNode) AddCveRules(inputs ...*CVERuleNode) *RepoNode {
	for _, input := range inputs {
		node := input.Node()

		// Add an edge to the input node
		r.AddEdge("cve_rules", node)
		// Create the invse edge also
		node.AddInverseEdge("repo", r.Node())
	}
	return r
}

func (r *RepoNode) SetOperation(op NodeOperation) *RepoNode {
	r.Operation = op
	return r
}

func NewTestCaseNode() *TestCaseNode {
	return &TestCaseNode{
		DataNode: &DataNode{
			Name: "test_case",
		},
	}
}

type TestCaseNode struct {
	*DataNode
}

func (tc *TestCaseNode) Node() *DataNode {
	return tc.DataNode
}
func (tc *TestCaseNode) SetName(v string) *TestCaseNode {
	tc.AddField("name", cty.StringVal(v))
	return tc
}
func (tc *TestCaseNode) SetResult(v bool) *TestCaseNode {
	tc.AddField("result", cty.BoolVal(v))
	return tc
}
func (tc *TestCaseNode) SetMessage(v string) *TestCaseNode {
	tc.AddField("message", cty.StringVal(v))
	return tc
}
func (tc *TestCaseNode) SetElapsed(v float64) *TestCaseNode {
	tc.AddField("elapsed", cty.NumberFloatVal(v))
	return tc
}

func (tc *TestCaseNode) SetRun(input *TestRunNode) *TestCaseNode {
	node := input.Node()

	// Add an edge to the input node
	tc.AddInverseEdge("run", node)
	// Create the invse edge also
	node.AddEdge("tests", tc.Node())
	return tc
}

func (tc *TestCaseNode) SetOperation(op NodeOperation) *TestCaseNode {
	tc.Operation = op
	return tc
}

func NewTestRunNode() *TestRunNode {
	return &TestRunNode{
		DataNode: &DataNode{
			Name: "test_run",
		},
	}
}

type TestRunNode struct {
	*DataNode
}

func (tr *TestRunNode) Node() *DataNode {
	return tr.DataNode
}
func (tr *TestRunNode) SetTool(v string) *TestRunNode {
	tr.AddField("tool", cty.StringVal(v))
	return tr
}

func (tr *TestRunNode) SetRelease(input *ReleaseNode) *TestRunNode {
	node := input.Node()

	// Add an edge to the input node
	tr.AddInverseEdge("release", node)
	// Create the invse edge also
	node.AddEdge("test_runs", tr.Node())
	return tr
}

func (tr *TestRunNode) SetEntry(input *ReleaseEntryNode) *TestRunNode {
	node := input.Node()

	// Add an edge to the input node
	tr.AddEdge("entry", node)
	// Create the invse edge also
	node.AddInverseEdge("test_run", tr.Node())
	return tr
}

func (tr *TestRunNode) AddTests(inputs ...*TestCaseNode) *TestRunNode {
	for _, input := range inputs {
		node := input.Node()

		// Add an edge to the input node
		tr.AddEdge("tests", node)
		// Create the invse edge also
		node.AddInverseEdge("run", tr.Node())
	}
	return tr
}

func (tr *TestRunNode) SetOperation(op NodeOperation) *TestRunNode {
	tr.Operation = op
	return tr
}

func NewVulnerabilityNode() *VulnerabilityNode {
	return &VulnerabilityNode{
		DataNode: &DataNode{
			Name: "vulnerability",
		},
	}
}

type VulnerabilityNode struct {
	*DataNode
}

func (v *VulnerabilityNode) Node() *DataNode {
	return v.DataNode
}

func (v *VulnerabilityNode) SetCve(input *CVENode) *VulnerabilityNode {
	node := input.Node()

	// Add an edge to the input node
	v.AddInverseEdge("cve", node)
	// Create the invse edge also
	node.AddEdge("found", v.Node())
	return v
}

func (v *VulnerabilityNode) SetScan(input *CVEScanNode) *VulnerabilityNode {
	node := input.Node()

	// Add an edge to the input node
	v.AddInverseEdge("scan", node)
	// Create the invse edge also
	node.AddEdge("vulnerabilities", v.Node())
	return v
}

func (v *VulnerabilityNode) SetComponent(input *ComponentNode) *VulnerabilityNode {
	node := input.Node()

	// Add an edge to the input node
	v.AddInverseEdge("component", node)
	// Create the invse edge also
	node.AddEdge("vulnerabilities", v.Node())
	return v
}

func (v *VulnerabilityNode) SetOperation(op NodeOperation) *VulnerabilityNode {
	v.Operation = op
	return v
}
