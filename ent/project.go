// Code generated by entc, DO NOT EDIT.

package ent

import (
	"fmt"
	"strings"

	"entgo.io/ent/dialect/sql"
	"github.com/valocode/bubbly/ent/project"
)

// Project is the model entity for the Project schema.
type Project struct {
	config `json:"-"`
	// ID of the ent.
	ID int `json:"id,omitempty"`
	// Name holds the value of the "name" field.
	Name string `json:"name,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the ProjectQuery when eager-loading is set.
	Edges ProjectEdges `json:"edges"`
}

// ProjectEdges holds the relations/edges for other nodes in the graph.
type ProjectEdges struct {
	// Repos holds the value of the repos edge.
	Repos []*Repo `json:"repos,omitempty"`
	// Releases holds the value of the releases edge.
	Releases []*Release `json:"releases,omitempty"`
	// CveRules holds the value of the cve_rules edge.
	CveRules []*CVERule `json:"cve_rules,omitempty"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [3]bool
}

// ReposOrErr returns the Repos value or an error if the edge
// was not loaded in eager-loading.
func (e ProjectEdges) ReposOrErr() ([]*Repo, error) {
	if e.loadedTypes[0] {
		return e.Repos, nil
	}
	return nil, &NotLoadedError{edge: "repos"}
}

// ReleasesOrErr returns the Releases value or an error if the edge
// was not loaded in eager-loading.
func (e ProjectEdges) ReleasesOrErr() ([]*Release, error) {
	if e.loadedTypes[1] {
		return e.Releases, nil
	}
	return nil, &NotLoadedError{edge: "releases"}
}

// CveRulesOrErr returns the CveRules value or an error if the edge
// was not loaded in eager-loading.
func (e ProjectEdges) CveRulesOrErr() ([]*CVERule, error) {
	if e.loadedTypes[2] {
		return e.CveRules, nil
	}
	return nil, &NotLoadedError{edge: "cve_rules"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*Project) scanValues(columns []string) ([]interface{}, error) {
	values := make([]interface{}, len(columns))
	for i := range columns {
		switch columns[i] {
		case project.FieldID:
			values[i] = new(sql.NullInt64)
		case project.FieldName:
			values[i] = new(sql.NullString)
		default:
			return nil, fmt.Errorf("unexpected column %q for type Project", columns[i])
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the Project fields.
func (pr *Project) assignValues(columns []string, values []interface{}) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case project.FieldID:
			value, ok := values[i].(*sql.NullInt64)
			if !ok {
				return fmt.Errorf("unexpected type %T for field id", value)
			}
			pr.ID = int(value.Int64)
		case project.FieldName:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field name", values[i])
			} else if value.Valid {
				pr.Name = value.String
			}
		}
	}
	return nil
}

// QueryRepos queries the "repos" edge of the Project entity.
func (pr *Project) QueryRepos() *RepoQuery {
	return (&ProjectClient{config: pr.config}).QueryRepos(pr)
}

// QueryReleases queries the "releases" edge of the Project entity.
func (pr *Project) QueryReleases() *ReleaseQuery {
	return (&ProjectClient{config: pr.config}).QueryReleases(pr)
}

// QueryCveRules queries the "cve_rules" edge of the Project entity.
func (pr *Project) QueryCveRules() *CVERuleQuery {
	return (&ProjectClient{config: pr.config}).QueryCveRules(pr)
}

// Update returns a builder for updating this Project.
// Note that you need to call Project.Unwrap() before calling this method if this Project
// was returned from a transaction, and the transaction was committed or rolled back.
func (pr *Project) Update() *ProjectUpdateOne {
	return (&ProjectClient{config: pr.config}).UpdateOne(pr)
}

// Unwrap unwraps the Project entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (pr *Project) Unwrap() *Project {
	tx, ok := pr.config.driver.(*txDriver)
	if !ok {
		panic("ent: Project is not a transactional entity")
	}
	pr.config.driver = tx.drv
	return pr
}

// String implements the fmt.Stringer.
func (pr *Project) String() string {
	var builder strings.Builder
	builder.WriteString("Project(")
	builder.WriteString(fmt.Sprintf("id=%v", pr.ID))
	builder.WriteString(", name=")
	builder.WriteString(pr.Name)
	builder.WriteByte(')')
	return builder.String()
}

// Projects is a parsable slice of Project.
type Projects []*Project

func (pr Projects) config(cfg config) {
	for _i := range pr {
		pr[_i].config = cfg
	}
}