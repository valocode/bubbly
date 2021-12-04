// Code generated by entc, DO NOT EDIT.

package ent

import (
	"encoding/json"
	"fmt"
	"strings"

	"entgo.io/ent/dialect/sql"
	"github.com/valocode/bubbly/ent/component"

	"github.com/valocode/bubbly/ent/organization"
	schema "github.com/valocode/bubbly/ent/schema/types"
)

// Component is the model entity for the Component schema.
type Component struct {
	config `json:"-"`
	// ID of the ent.
	ID int `json:"id,omitempty"`
	// Scheme holds the value of the "scheme" field.
	Scheme string `json:"scheme,omitempty"`
	// Namespace holds the value of the "namespace" field.
	Namespace string `json:"namespace,omitempty"`
	// Name holds the value of the "name" field.
	Name string `json:"name,omitempty"`
	// Version holds the value of the "version" field.
	Version string `json:"version,omitempty"`
	// Description holds the value of the "description" field.
	Description string `json:"description,omitempty"`
	// URL holds the value of the "url" field.
	URL string `json:"url,omitempty"`
	// Metadata holds the value of the "metadata" field.
	Metadata schema.Metadata `json:"metadata,omitempty"`
	// Labels holds the value of the "labels" field.
	Labels schema.Labels `json:"labels,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the ComponentQuery when eager-loading is set.
	Edges           ComponentEdges `json:"edges"`
	component_owner *int
}

// ComponentEdges holds the relations/edges for other nodes in the graph.
type ComponentEdges struct {
	// Owner holds the value of the owner edge.
	Owner *Organization `json:"owner,omitempty"`
	// Vulnerabilities holds the value of the vulnerabilities edge.
	Vulnerabilities []*Vulnerability `json:"vulnerabilities,omitempty"`
	// Licenses holds the value of the licenses edge.
	Licenses []*License `json:"licenses,omitempty"`
	// Uses holds the value of the uses edge.
	Uses []*ReleaseComponent `json:"uses,omitempty"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [4]bool
}

// OwnerOrErr returns the Owner value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e ComponentEdges) OwnerOrErr() (*Organization, error) {
	if e.loadedTypes[0] {
		if e.Owner == nil {
			// The edge owner was loaded in eager-loading,
			// but was not found.
			return nil, &NotFoundError{label: organization.Label}
		}
		return e.Owner, nil
	}
	return nil, &NotLoadedError{edge: "owner"}
}

// VulnerabilitiesOrErr returns the Vulnerabilities value or an error if the edge
// was not loaded in eager-loading.
func (e ComponentEdges) VulnerabilitiesOrErr() ([]*Vulnerability, error) {
	if e.loadedTypes[1] {
		return e.Vulnerabilities, nil
	}
	return nil, &NotLoadedError{edge: "vulnerabilities"}
}

// LicensesOrErr returns the Licenses value or an error if the edge
// was not loaded in eager-loading.
func (e ComponentEdges) LicensesOrErr() ([]*License, error) {
	if e.loadedTypes[2] {
		return e.Licenses, nil
	}
	return nil, &NotLoadedError{edge: "licenses"}
}

// UsesOrErr returns the Uses value or an error if the edge
// was not loaded in eager-loading.
func (e ComponentEdges) UsesOrErr() ([]*ReleaseComponent, error) {
	if e.loadedTypes[3] {
		return e.Uses, nil
	}
	return nil, &NotLoadedError{edge: "uses"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*Component) scanValues(columns []string) ([]interface{}, error) {
	values := make([]interface{}, len(columns))
	for i := range columns {
		switch columns[i] {
		case component.FieldMetadata, component.FieldLabels:
			values[i] = new([]byte)
		case component.FieldID:
			values[i] = new(sql.NullInt64)
		case component.FieldScheme, component.FieldNamespace, component.FieldName, component.FieldVersion, component.FieldDescription, component.FieldURL:
			values[i] = new(sql.NullString)
		case component.ForeignKeys[0]: // component_owner
			values[i] = new(sql.NullInt64)
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
		case component.FieldScheme:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field scheme", values[i])
			} else if value.Valid {
				c.Scheme = value.String
			}
		case component.FieldNamespace:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field namespace", values[i])
			} else if value.Valid {
				c.Namespace = value.String
			}
		case component.FieldName:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field name", values[i])
			} else if value.Valid {
				c.Name = value.String
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
		case component.FieldMetadata:
			if value, ok := values[i].(*[]byte); !ok {
				return fmt.Errorf("unexpected type %T for field metadata", values[i])
			} else if value != nil && len(*value) > 0 {
				if err := json.Unmarshal(*value, &c.Metadata); err != nil {
					return fmt.Errorf("unmarshal field metadata: %w", err)
				}
			}
		case component.FieldLabels:
			if value, ok := values[i].(*[]byte); !ok {
				return fmt.Errorf("unexpected type %T for field labels", values[i])
			} else if value != nil && len(*value) > 0 {
				if err := json.Unmarshal(*value, &c.Labels); err != nil {
					return fmt.Errorf("unmarshal field labels: %w", err)
				}
			}
		case component.ForeignKeys[0]:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for edge-field component_owner", value)
			} else if value.Valid {
				c.component_owner = new(int)
				*c.component_owner = int(value.Int64)
			}
		}
	}
	return nil
}

// QueryOwner queries the "owner" edge of the Component entity.
func (c *Component) QueryOwner() *OrganizationQuery {
	return (&ComponentClient{config: c.config}).QueryOwner(c)
}

// QueryVulnerabilities queries the "vulnerabilities" edge of the Component entity.
func (c *Component) QueryVulnerabilities() *VulnerabilityQuery {
	return (&ComponentClient{config: c.config}).QueryVulnerabilities(c)
}

// QueryLicenses queries the "licenses" edge of the Component entity.
func (c *Component) QueryLicenses() *LicenseQuery {
	return (&ComponentClient{config: c.config}).QueryLicenses(c)
}

// QueryUses queries the "uses" edge of the Component entity.
func (c *Component) QueryUses() *ReleaseComponentQuery {
	return (&ComponentClient{config: c.config}).QueryUses(c)
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
	builder.WriteString(", scheme=")
	builder.WriteString(c.Scheme)
	builder.WriteString(", namespace=")
	builder.WriteString(c.Namespace)
	builder.WriteString(", name=")
	builder.WriteString(c.Name)
	builder.WriteString(", version=")
	builder.WriteString(c.Version)
	builder.WriteString(", description=")
	builder.WriteString(c.Description)
	builder.WriteString(", url=")
	builder.WriteString(c.URL)
	builder.WriteString(", metadata=")
	builder.WriteString(fmt.Sprintf("%v", c.Metadata))
	builder.WriteString(", labels=")
	builder.WriteString(fmt.Sprintf("%v", c.Labels))
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
