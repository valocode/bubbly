// Code generated by entc, DO NOT EDIT.

package ent

import (
	"fmt"
	"strings"

	"entgo.io/ent/dialect/sql"
	"github.com/valocode/bubbly/ent/gitcommit"
	"github.com/valocode/bubbly/ent/release"
)

// Release is the model entity for the Release schema.
type Release struct {
	config `json:"-"`
	// ID of the ent.
	ID int `json:"id,omitempty"`
	// Name holds the value of the "name" field.
	Name string `json:"name,omitempty"`
	// Version holds the value of the "version" field.
	Version string `json:"version,omitempty"`
	// Status holds the value of the "status" field.
	Status release.Status `json:"status,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the ReleaseQuery when eager-loading is set.
	Edges              ReleaseEdges `json:"edges"`
	git_commit_release *int
}

// ReleaseEdges holds the relations/edges for other nodes in the graph.
type ReleaseEdges struct {
	// Subreleases holds the value of the subreleases edge.
	Subreleases []*Release `json:"subreleases,omitempty"`
	// Dependencies holds the value of the dependencies edge.
	Dependencies []*Release `json:"dependencies,omitempty"`
	// Commit holds the value of the commit edge.
	Commit *GitCommit `json:"commit,omitempty"`
	// Log holds the value of the log edge.
	Log []*ReleaseEntry `json:"log,omitempty"`
	// Artifacts holds the value of the artifacts edge.
	Artifacts []*Artifact `json:"artifacts,omitempty"`
	// Components holds the value of the components edge.
	Components []*ReleaseComponent `json:"components,omitempty"`
	// Vulnerabilities holds the value of the vulnerabilities edge.
	Vulnerabilities []*ReleaseVulnerability `json:"vulnerabilities,omitempty"`
	// CodeScans holds the value of the code_scans edge.
	CodeScans []*CodeScan `json:"code_scans,omitempty"`
	// TestRuns holds the value of the test_runs edge.
	TestRuns []*TestRun `json:"test_runs,omitempty"`
	// VulnerabilityReviews holds the value of the vulnerability_reviews edge.
	VulnerabilityReviews []*VulnerabilityReview `json:"vulnerability_reviews,omitempty"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [10]bool
}

// SubreleasesOrErr returns the Subreleases value or an error if the edge
// was not loaded in eager-loading.
func (e ReleaseEdges) SubreleasesOrErr() ([]*Release, error) {
	if e.loadedTypes[0] {
		return e.Subreleases, nil
	}
	return nil, &NotLoadedError{edge: "subreleases"}
}

// DependenciesOrErr returns the Dependencies value or an error if the edge
// was not loaded in eager-loading.
func (e ReleaseEdges) DependenciesOrErr() ([]*Release, error) {
	if e.loadedTypes[1] {
		return e.Dependencies, nil
	}
	return nil, &NotLoadedError{edge: "dependencies"}
}

// CommitOrErr returns the Commit value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e ReleaseEdges) CommitOrErr() (*GitCommit, error) {
	if e.loadedTypes[2] {
		if e.Commit == nil {
			// The edge commit was loaded in eager-loading,
			// but was not found.
			return nil, &NotFoundError{label: gitcommit.Label}
		}
		return e.Commit, nil
	}
	return nil, &NotLoadedError{edge: "commit"}
}

// LogOrErr returns the Log value or an error if the edge
// was not loaded in eager-loading.
func (e ReleaseEdges) LogOrErr() ([]*ReleaseEntry, error) {
	if e.loadedTypes[3] {
		return e.Log, nil
	}
	return nil, &NotLoadedError{edge: "log"}
}

// ArtifactsOrErr returns the Artifacts value or an error if the edge
// was not loaded in eager-loading.
func (e ReleaseEdges) ArtifactsOrErr() ([]*Artifact, error) {
	if e.loadedTypes[4] {
		return e.Artifacts, nil
	}
	return nil, &NotLoadedError{edge: "artifacts"}
}

// ComponentsOrErr returns the Components value or an error if the edge
// was not loaded in eager-loading.
func (e ReleaseEdges) ComponentsOrErr() ([]*ReleaseComponent, error) {
	if e.loadedTypes[5] {
		return e.Components, nil
	}
	return nil, &NotLoadedError{edge: "components"}
}

// VulnerabilitiesOrErr returns the Vulnerabilities value or an error if the edge
// was not loaded in eager-loading.
func (e ReleaseEdges) VulnerabilitiesOrErr() ([]*ReleaseVulnerability, error) {
	if e.loadedTypes[6] {
		return e.Vulnerabilities, nil
	}
	return nil, &NotLoadedError{edge: "vulnerabilities"}
}

// CodeScansOrErr returns the CodeScans value or an error if the edge
// was not loaded in eager-loading.
func (e ReleaseEdges) CodeScansOrErr() ([]*CodeScan, error) {
	if e.loadedTypes[7] {
		return e.CodeScans, nil
	}
	return nil, &NotLoadedError{edge: "code_scans"}
}

// TestRunsOrErr returns the TestRuns value or an error if the edge
// was not loaded in eager-loading.
func (e ReleaseEdges) TestRunsOrErr() ([]*TestRun, error) {
	if e.loadedTypes[8] {
		return e.TestRuns, nil
	}
	return nil, &NotLoadedError{edge: "test_runs"}
}

// VulnerabilityReviewsOrErr returns the VulnerabilityReviews value or an error if the edge
// was not loaded in eager-loading.
func (e ReleaseEdges) VulnerabilityReviewsOrErr() ([]*VulnerabilityReview, error) {
	if e.loadedTypes[9] {
		return e.VulnerabilityReviews, nil
	}
	return nil, &NotLoadedError{edge: "vulnerability_reviews"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*Release) scanValues(columns []string) ([]interface{}, error) {
	values := make([]interface{}, len(columns))
	for i := range columns {
		switch columns[i] {
		case release.FieldID:
			values[i] = new(sql.NullInt64)
		case release.FieldName, release.FieldVersion, release.FieldStatus:
			values[i] = new(sql.NullString)
		case release.ForeignKeys[0]: // git_commit_release
			values[i] = new(sql.NullInt64)
		default:
			return nil, fmt.Errorf("unexpected column %q for type Release", columns[i])
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the Release fields.
func (r *Release) assignValues(columns []string, values []interface{}) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case release.FieldID:
			value, ok := values[i].(*sql.NullInt64)
			if !ok {
				return fmt.Errorf("unexpected type %T for field id", value)
			}
			r.ID = int(value.Int64)
		case release.FieldName:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field name", values[i])
			} else if value.Valid {
				r.Name = value.String
			}
		case release.FieldVersion:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field version", values[i])
			} else if value.Valid {
				r.Version = value.String
			}
		case release.FieldStatus:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field status", values[i])
			} else if value.Valid {
				r.Status = release.Status(value.String)
			}
		case release.ForeignKeys[0]:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for edge-field git_commit_release", value)
			} else if value.Valid {
				r.git_commit_release = new(int)
				*r.git_commit_release = int(value.Int64)
			}
		}
	}
	return nil
}

// QuerySubreleases queries the "subreleases" edge of the Release entity.
func (r *Release) QuerySubreleases() *ReleaseQuery {
	return (&ReleaseClient{config: r.config}).QuerySubreleases(r)
}

// QueryDependencies queries the "dependencies" edge of the Release entity.
func (r *Release) QueryDependencies() *ReleaseQuery {
	return (&ReleaseClient{config: r.config}).QueryDependencies(r)
}

// QueryCommit queries the "commit" edge of the Release entity.
func (r *Release) QueryCommit() *GitCommitQuery {
	return (&ReleaseClient{config: r.config}).QueryCommit(r)
}

// QueryLog queries the "log" edge of the Release entity.
func (r *Release) QueryLog() *ReleaseEntryQuery {
	return (&ReleaseClient{config: r.config}).QueryLog(r)
}

// QueryArtifacts queries the "artifacts" edge of the Release entity.
func (r *Release) QueryArtifacts() *ArtifactQuery {
	return (&ReleaseClient{config: r.config}).QueryArtifacts(r)
}

// QueryComponents queries the "components" edge of the Release entity.
func (r *Release) QueryComponents() *ReleaseComponentQuery {
	return (&ReleaseClient{config: r.config}).QueryComponents(r)
}

// QueryVulnerabilities queries the "vulnerabilities" edge of the Release entity.
func (r *Release) QueryVulnerabilities() *ReleaseVulnerabilityQuery {
	return (&ReleaseClient{config: r.config}).QueryVulnerabilities(r)
}

// QueryCodeScans queries the "code_scans" edge of the Release entity.
func (r *Release) QueryCodeScans() *CodeScanQuery {
	return (&ReleaseClient{config: r.config}).QueryCodeScans(r)
}

// QueryTestRuns queries the "test_runs" edge of the Release entity.
func (r *Release) QueryTestRuns() *TestRunQuery {
	return (&ReleaseClient{config: r.config}).QueryTestRuns(r)
}

// QueryVulnerabilityReviews queries the "vulnerability_reviews" edge of the Release entity.
func (r *Release) QueryVulnerabilityReviews() *VulnerabilityReviewQuery {
	return (&ReleaseClient{config: r.config}).QueryVulnerabilityReviews(r)
}

// Update returns a builder for updating this Release.
// Note that you need to call Release.Unwrap() before calling this method if this Release
// was returned from a transaction, and the transaction was committed or rolled back.
func (r *Release) Update() *ReleaseUpdateOne {
	return (&ReleaseClient{config: r.config}).UpdateOne(r)
}

// Unwrap unwraps the Release entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (r *Release) Unwrap() *Release {
	tx, ok := r.config.driver.(*txDriver)
	if !ok {
		panic("ent: Release is not a transactional entity")
	}
	r.config.driver = tx.drv
	return r
}

// String implements the fmt.Stringer.
func (r *Release) String() string {
	var builder strings.Builder
	builder.WriteString("Release(")
	builder.WriteString(fmt.Sprintf("id=%v", r.ID))
	builder.WriteString(", name=")
	builder.WriteString(r.Name)
	builder.WriteString(", version=")
	builder.WriteString(r.Version)
	builder.WriteString(", status=")
	builder.WriteString(fmt.Sprintf("%v", r.Status))
	builder.WriteByte(')')
	return builder.String()
}

// Releases is a parsable slice of Release.
type Releases []*Release

func (r Releases) config(cfg config) {
	for _i := range r {
		r[_i].config = cfg
	}
}
