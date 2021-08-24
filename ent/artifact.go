// Code generated by entc, DO NOT EDIT.

package ent

import (
	"fmt"
	"strings"
	"time"

	"entgo.io/ent/dialect/sql"
	"github.com/valocode/bubbly/ent/artifact"
	"github.com/valocode/bubbly/ent/release"
	"github.com/valocode/bubbly/ent/releaseentry"
)

// Artifact is the model entity for the Artifact schema.
type Artifact struct {
	config `json:"-"`
	// ID of the ent.
	ID int `json:"id,omitempty"`
	// Name holds the value of the "name" field.
	Name string `json:"name,omitempty"`
	// Sha256 holds the value of the "sha256" field.
	Sha256 string `json:"sha256,omitempty"`
	// Type holds the value of the "type" field.
	Type artifact.Type `json:"type,omitempty"`
	// Time holds the value of the "time" field.
	Time time.Time `json:"time,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the ArtifactQuery when eager-loading is set.
	Edges                  ArtifactEdges `json:"edges"`
	artifact_release       *int
	release_entry_artifact *int
}

// ArtifactEdges holds the relations/edges for other nodes in the graph.
type ArtifactEdges struct {
	// Release holds the value of the release edge.
	Release *Release `json:"release,omitempty"`
	// Entry holds the value of the entry edge.
	Entry *ReleaseEntry `json:"entry,omitempty"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [2]bool
}

// ReleaseOrErr returns the Release value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e ArtifactEdges) ReleaseOrErr() (*Release, error) {
	if e.loadedTypes[0] {
		if e.Release == nil {
			// The edge release was loaded in eager-loading,
			// but was not found.
			return nil, &NotFoundError{label: release.Label}
		}
		return e.Release, nil
	}
	return nil, &NotLoadedError{edge: "release"}
}

// EntryOrErr returns the Entry value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e ArtifactEdges) EntryOrErr() (*ReleaseEntry, error) {
	if e.loadedTypes[1] {
		if e.Entry == nil {
			// The edge entry was loaded in eager-loading,
			// but was not found.
			return nil, &NotFoundError{label: releaseentry.Label}
		}
		return e.Entry, nil
	}
	return nil, &NotLoadedError{edge: "entry"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*Artifact) scanValues(columns []string) ([]interface{}, error) {
	values := make([]interface{}, len(columns))
	for i := range columns {
		switch columns[i] {
		case artifact.FieldID:
			values[i] = new(sql.NullInt64)
		case artifact.FieldName, artifact.FieldSha256, artifact.FieldType:
			values[i] = new(sql.NullString)
		case artifact.FieldTime:
			values[i] = new(sql.NullTime)
		case artifact.ForeignKeys[0]: // artifact_release
			values[i] = new(sql.NullInt64)
		case artifact.ForeignKeys[1]: // release_entry_artifact
			values[i] = new(sql.NullInt64)
		default:
			return nil, fmt.Errorf("unexpected column %q for type Artifact", columns[i])
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the Artifact fields.
func (a *Artifact) assignValues(columns []string, values []interface{}) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case artifact.FieldID:
			value, ok := values[i].(*sql.NullInt64)
			if !ok {
				return fmt.Errorf("unexpected type %T for field id", value)
			}
			a.ID = int(value.Int64)
		case artifact.FieldName:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field name", values[i])
			} else if value.Valid {
				a.Name = value.String
			}
		case artifact.FieldSha256:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field sha256", values[i])
			} else if value.Valid {
				a.Sha256 = value.String
			}
		case artifact.FieldType:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field type", values[i])
			} else if value.Valid {
				a.Type = artifact.Type(value.String)
			}
		case artifact.FieldTime:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field time", values[i])
			} else if value.Valid {
				a.Time = value.Time
			}
		case artifact.ForeignKeys[0]:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for edge-field artifact_release", value)
			} else if value.Valid {
				a.artifact_release = new(int)
				*a.artifact_release = int(value.Int64)
			}
		case artifact.ForeignKeys[1]:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for edge-field release_entry_artifact", value)
			} else if value.Valid {
				a.release_entry_artifact = new(int)
				*a.release_entry_artifact = int(value.Int64)
			}
		}
	}
	return nil
}

// QueryRelease queries the "release" edge of the Artifact entity.
func (a *Artifact) QueryRelease() *ReleaseQuery {
	return (&ArtifactClient{config: a.config}).QueryRelease(a)
}

// QueryEntry queries the "entry" edge of the Artifact entity.
func (a *Artifact) QueryEntry() *ReleaseEntryQuery {
	return (&ArtifactClient{config: a.config}).QueryEntry(a)
}

// Update returns a builder for updating this Artifact.
// Note that you need to call Artifact.Unwrap() before calling this method if this Artifact
// was returned from a transaction, and the transaction was committed or rolled back.
func (a *Artifact) Update() *ArtifactUpdateOne {
	return (&ArtifactClient{config: a.config}).UpdateOne(a)
}

// Unwrap unwraps the Artifact entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (a *Artifact) Unwrap() *Artifact {
	tx, ok := a.config.driver.(*txDriver)
	if !ok {
		panic("ent: Artifact is not a transactional entity")
	}
	a.config.driver = tx.drv
	return a
}

// String implements the fmt.Stringer.
func (a *Artifact) String() string {
	var builder strings.Builder
	builder.WriteString("Artifact(")
	builder.WriteString(fmt.Sprintf("id=%v", a.ID))
	builder.WriteString(", name=")
	builder.WriteString(a.Name)
	builder.WriteString(", sha256=")
	builder.WriteString(a.Sha256)
	builder.WriteString(", type=")
	builder.WriteString(fmt.Sprintf("%v", a.Type))
	builder.WriteString(", time=")
	builder.WriteString(a.Time.Format(time.ANSIC))
	builder.WriteByte(')')
	return builder.String()
}

// Artifacts is a parsable slice of Artifact.
type Artifacts []*Artifact

func (a Artifacts) config(cfg config) {
	for _i := range a {
		a[_i].config = cfg
	}
}
