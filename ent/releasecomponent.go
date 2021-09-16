// Code generated by entc, DO NOT EDIT.

package ent

import (
	"fmt"
	"strings"

	"entgo.io/ent/dialect/sql"
	"github.com/valocode/bubbly/ent/component"
	"github.com/valocode/bubbly/ent/release"
	"github.com/valocode/bubbly/ent/releasecomponent"
)

// ReleaseComponent is the model entity for the ReleaseComponent schema.
type ReleaseComponent struct {
	config `json:"-"`
	// ID of the ent.
	ID int `json:"id,omitempty"`
	// Type holds the value of the "type" field.
	//
	// The type indicates how the component is used in the project,
	// e.g. whether it is embedded into the build (static link) or just
	// distributed (dynamic link) or just a development dependency
	Type releasecomponent.Type `json:"type,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the ReleaseComponentQuery when eager-loading is set.
	Edges                       ReleaseComponentEdges `json:"edges"`
	release_component_release   *int
	release_component_component *int
}

// ReleaseComponentEdges holds the relations/edges for other nodes in the graph.
type ReleaseComponentEdges struct {
	// Release holds the value of the release edge.
	Release *Release `json:"release,omitempty"`
	// Scans holds the value of the scans edge.
	Scans []*CodeScan `json:"scans,omitempty"`
	// Component holds the value of the component edge.
	Component *Component `json:"component,omitempty"`
	// Vulnerabilities holds the value of the vulnerabilities edge.
	Vulnerabilities []*ReleaseVulnerability `json:"vulnerabilities,omitempty"`
	// Licenses holds the value of the licenses edge.
	Licenses []*ReleaseLicense `json:"licenses,omitempty"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [5]bool
}

// ReleaseOrErr returns the Release value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e ReleaseComponentEdges) ReleaseOrErr() (*Release, error) {
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

// ScansOrErr returns the Scans value or an error if the edge
// was not loaded in eager-loading.
func (e ReleaseComponentEdges) ScansOrErr() ([]*CodeScan, error) {
	if e.loadedTypes[1] {
		return e.Scans, nil
	}
	return nil, &NotLoadedError{edge: "scans"}
}

// ComponentOrErr returns the Component value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e ReleaseComponentEdges) ComponentOrErr() (*Component, error) {
	if e.loadedTypes[2] {
		if e.Component == nil {
			// The edge component was loaded in eager-loading,
			// but was not found.
			return nil, &NotFoundError{label: component.Label}
		}
		return e.Component, nil
	}
	return nil, &NotLoadedError{edge: "component"}
}

// VulnerabilitiesOrErr returns the Vulnerabilities value or an error if the edge
// was not loaded in eager-loading.
func (e ReleaseComponentEdges) VulnerabilitiesOrErr() ([]*ReleaseVulnerability, error) {
	if e.loadedTypes[3] {
		return e.Vulnerabilities, nil
	}
	return nil, &NotLoadedError{edge: "vulnerabilities"}
}

// LicensesOrErr returns the Licenses value or an error if the edge
// was not loaded in eager-loading.
func (e ReleaseComponentEdges) LicensesOrErr() ([]*ReleaseLicense, error) {
	if e.loadedTypes[4] {
		return e.Licenses, nil
	}
	return nil, &NotLoadedError{edge: "licenses"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*ReleaseComponent) scanValues(columns []string) ([]interface{}, error) {
	values := make([]interface{}, len(columns))
	for i := range columns {
		switch columns[i] {
		case releasecomponent.FieldID:
			values[i] = new(sql.NullInt64)
		case releasecomponent.FieldType:
			values[i] = new(sql.NullString)
		case releasecomponent.ForeignKeys[0]: // release_component_release
			values[i] = new(sql.NullInt64)
		case releasecomponent.ForeignKeys[1]: // release_component_component
			values[i] = new(sql.NullInt64)
		default:
			return nil, fmt.Errorf("unexpected column %q for type ReleaseComponent", columns[i])
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the ReleaseComponent fields.
func (rc *ReleaseComponent) assignValues(columns []string, values []interface{}) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case releasecomponent.FieldID:
			value, ok := values[i].(*sql.NullInt64)
			if !ok {
				return fmt.Errorf("unexpected type %T for field id", value)
			}
			rc.ID = int(value.Int64)
		case releasecomponent.FieldType:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field type", values[i])
			} else if value.Valid {
				rc.Type = releasecomponent.Type(value.String)
			}
		case releasecomponent.ForeignKeys[0]:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for edge-field release_component_release", value)
			} else if value.Valid {
				rc.release_component_release = new(int)
				*rc.release_component_release = int(value.Int64)
			}
		case releasecomponent.ForeignKeys[1]:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for edge-field release_component_component", value)
			} else if value.Valid {
				rc.release_component_component = new(int)
				*rc.release_component_component = int(value.Int64)
			}
		}
	}
	return nil
}

// QueryRelease queries the "release" edge of the ReleaseComponent entity.
func (rc *ReleaseComponent) QueryRelease() *ReleaseQuery {
	return (&ReleaseComponentClient{config: rc.config}).QueryRelease(rc)
}

// QueryScans queries the "scans" edge of the ReleaseComponent entity.
func (rc *ReleaseComponent) QueryScans() *CodeScanQuery {
	return (&ReleaseComponentClient{config: rc.config}).QueryScans(rc)
}

// QueryComponent queries the "component" edge of the ReleaseComponent entity.
func (rc *ReleaseComponent) QueryComponent() *ComponentQuery {
	return (&ReleaseComponentClient{config: rc.config}).QueryComponent(rc)
}

// QueryVulnerabilities queries the "vulnerabilities" edge of the ReleaseComponent entity.
func (rc *ReleaseComponent) QueryVulnerabilities() *ReleaseVulnerabilityQuery {
	return (&ReleaseComponentClient{config: rc.config}).QueryVulnerabilities(rc)
}

// QueryLicenses queries the "licenses" edge of the ReleaseComponent entity.
func (rc *ReleaseComponent) QueryLicenses() *ReleaseLicenseQuery {
	return (&ReleaseComponentClient{config: rc.config}).QueryLicenses(rc)
}

// Update returns a builder for updating this ReleaseComponent.
// Note that you need to call ReleaseComponent.Unwrap() before calling this method if this ReleaseComponent
// was returned from a transaction, and the transaction was committed or rolled back.
func (rc *ReleaseComponent) Update() *ReleaseComponentUpdateOne {
	return (&ReleaseComponentClient{config: rc.config}).UpdateOne(rc)
}

// Unwrap unwraps the ReleaseComponent entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (rc *ReleaseComponent) Unwrap() *ReleaseComponent {
	tx, ok := rc.config.driver.(*txDriver)
	if !ok {
		panic("ent: ReleaseComponent is not a transactional entity")
	}
	rc.config.driver = tx.drv
	return rc
}

// String implements the fmt.Stringer.
func (rc *ReleaseComponent) String() string {
	var builder strings.Builder
	builder.WriteString("ReleaseComponent(")
	builder.WriteString(fmt.Sprintf("id=%v", rc.ID))
	builder.WriteString(", type=")
	builder.WriteString(fmt.Sprintf("%v", rc.Type))
	builder.WriteByte(')')
	return builder.String()
}

// ReleaseComponents is a parsable slice of ReleaseComponent.
type ReleaseComponents []*ReleaseComponent

func (rc ReleaseComponents) config(cfg config) {
	for _i := range rc {
		rc[_i].config = cfg
	}
}
