// Code generated by entc, DO NOT EDIT.

package ent

import (
	"fmt"
	"strings"

	"entgo.io/ent/dialect/sql"
	"github.com/valocode/bubbly/ent/component"
)

// Component is the model entity for the Component schema.
type Component struct {
	config `json:"-"`
	// ID of the ent.
	ID int `json:"id,omitempty"`
	// Name holds the value of the "name" field.
	Name string `json:"name,omitempty"`
	// Vendor holds the value of the "vendor" field.
	Vendor string `json:"vendor,omitempty"`
	// Version holds the value of the "version" field.
	Version string `json:"version,omitempty"`
	// Description holds the value of the "description" field.
	Description string `json:"description,omitempty"`
	// URL holds the value of the "url" field.
	URL string `json:"url,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the ComponentQuery when eager-loading is set.
	Edges ComponentEdges `json:"edges"`
}

// ComponentEdges holds the relations/edges for other nodes in the graph.
type ComponentEdges struct {
	// Vulnerabilities holds the value of the vulnerabilities edge.
	Vulnerabilities []*Vulnerability `json:"vulnerabilities,omitempty"`
	// Licenses holds the value of the licenses edge.
	Licenses []*License `json:"licenses,omitempty"`
	// Release holds the value of the release edge.
	Release []*Release `json:"release,omitempty"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [3]bool
}

// VulnerabilitiesOrErr returns the Vulnerabilities value or an error if the edge
// was not loaded in eager-loading.
func (e ComponentEdges) VulnerabilitiesOrErr() ([]*Vulnerability, error) {
	if e.loadedTypes[0] {
		return e.Vulnerabilities, nil
	}
	return nil, &NotLoadedError{edge: "vulnerabilities"}
}

// LicensesOrErr returns the Licenses value or an error if the edge
// was not loaded in eager-loading.
func (e ComponentEdges) LicensesOrErr() ([]*License, error) {
	if e.loadedTypes[1] {
		return e.Licenses, nil
	}
	return nil, &NotLoadedError{edge: "licenses"}
}

// ReleaseOrErr returns the Release value or an error if the edge
// was not loaded in eager-loading.
func (e ComponentEdges) ReleaseOrErr() ([]*Release, error) {
	if e.loadedTypes[2] {
		return e.Release, nil
	}
	return nil, &NotLoadedError{edge: "release"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*Component) scanValues(columns []string) ([]interface{}, error) {
	values := make([]interface{}, len(columns))
	for i := range columns {
		switch columns[i] {
		case component.FieldID:
			values[i] = new(sql.NullInt64)
		case component.FieldName, component.FieldVendor, component.FieldVersion, component.FieldDescription, component.FieldURL:
			values[i] = new(sql.NullString)
		default:
			return nil, fmt.Errorf("unexpected column %q for type Component", columns[i])
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the Component fields.
func (c *Component) assignValues(columns []string, values []interface{}) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case component.FieldID:
			value, ok := values[i].(*sql.NullInt64)
			if !ok {
				return fmt.Errorf("unexpected type %T for field id", value)
			}
			c.ID = int(value.Int64)
		case component.FieldName:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field name", values[i])
			} else if value.Valid {
				c.Name = value.String
			}
		case component.FieldVendor:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field vendor", values[i])
			} else if value.Valid {
				c.Vendor = value.String
			}
		case component.FieldVersion:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field version", values[i])
			} else if value.Valid {
				c.Version = value.String
			}
		case component.FieldDescription:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field description", values[i])
			} else if value.Valid {
				c.Description = value.String
			}
		case component.FieldURL:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field url", values[i])
			} else if value.Valid {
				c.URL = value.String
			}
		}
	}
	return nil
}

// QueryVulnerabilities queries the "vulnerabilities" edge of the Component entity.
func (c *Component) QueryVulnerabilities() *VulnerabilityQuery {
	return (&ComponentClient{config: c.config}).QueryVulnerabilities(c)
}

// QueryLicenses queries the "licenses" edge of the Component entity.
func (c *Component) QueryLicenses() *LicenseQuery {
	return (&ComponentClient{config: c.config}).QueryLicenses(c)
}

// QueryRelease queries the "release" edge of the Component entity.
func (c *Component) QueryRelease() *ReleaseQuery {
	return (&ComponentClient{config: c.config}).QueryRelease(c)
}

// Update returns a builder for updating this Component.
// Note that you need to call Component.Unwrap() before calling this method if this Component
// was returned from a transaction, and the transaction was committed or rolled back.
func (c *Component) Update() *ComponentUpdateOne {
	return (&ComponentClient{config: c.config}).UpdateOne(c)
}

// Unwrap unwraps the Component entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (c *Component) Unwrap() *Component {
	tx, ok := c.config.driver.(*txDriver)
	if !ok {
		panic("ent: Component is not a transactional entity")
	}
	c.config.driver = tx.drv
	return c
}

// String implements the fmt.Stringer.
func (c *Component) String() string {
	var builder strings.Builder
	builder.WriteString("Component(")
	builder.WriteString(fmt.Sprintf("id=%v", c.ID))
	builder.WriteString(", name=")
	builder.WriteString(c.Name)
	builder.WriteString(", vendor=")
	builder.WriteString(c.Vendor)
	builder.WriteString(", version=")
	builder.WriteString(c.Version)
	builder.WriteString(", description=")
	builder.WriteString(c.Description)
	builder.WriteString(", url=")
	builder.WriteString(c.URL)
	builder.WriteByte(')')
	return builder.String()
}

// Components is a parsable slice of Component.
type Components []*Component

func (c Components) config(cfg config) {
	for _i := range c {
		c[_i].config = cfg
	}
}
