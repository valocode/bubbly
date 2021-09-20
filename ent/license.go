// Code generated by entc, DO NOT EDIT.

package ent

import (
	"fmt"
	"strings"

	"entgo.io/ent/dialect/sql"
	"github.com/valocode/bubbly/ent/license"
	"github.com/valocode/bubbly/ent/organization"
	"github.com/valocode/bubbly/ent/spdxlicense"
)

// License is the model entity for the License schema.
type License struct {
	config `json:"-"`
	// ID of the ent.
	ID int `json:"id,omitempty"`
	// LicenseID holds the value of the "license_id" field.
	LicenseID string `json:"license_id,omitempty"`
	// Name holds the value of the "name" field.
	Name string `json:"name,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the LicenseQuery when eager-loading is set.
	Edges         LicenseEdges `json:"edges"`
	license_owner *int
	license_spdx  *int
}

// LicenseEdges holds the relations/edges for other nodes in the graph.
type LicenseEdges struct {
	// Owner holds the value of the owner edge.
	Owner *Organization `json:"owner,omitempty"`
	// Spdx holds the value of the spdx edge.
	Spdx *SPDXLicense `json:"spdx,omitempty"`
	// Components holds the value of the components edge.
	Components []*Component `json:"components,omitempty"`
	// Instances holds the value of the instances edge.
	Instances []*ReleaseLicense `json:"instances,omitempty"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [4]bool
}

// OwnerOrErr returns the Owner value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e LicenseEdges) OwnerOrErr() (*Organization, error) {
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

// SpdxOrErr returns the Spdx value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e LicenseEdges) SpdxOrErr() (*SPDXLicense, error) {
	if e.loadedTypes[1] {
		if e.Spdx == nil {
			// The edge spdx was loaded in eager-loading,
			// but was not found.
			return nil, &NotFoundError{label: spdxlicense.Label}
		}
		return e.Spdx, nil
	}
	return nil, &NotLoadedError{edge: "spdx"}
}

// ComponentsOrErr returns the Components value or an error if the edge
// was not loaded in eager-loading.
func (e LicenseEdges) ComponentsOrErr() ([]*Component, error) {
	if e.loadedTypes[2] {
		return e.Components, nil
	}
	return nil, &NotLoadedError{edge: "components"}
}

// InstancesOrErr returns the Instances value or an error if the edge
// was not loaded in eager-loading.
func (e LicenseEdges) InstancesOrErr() ([]*ReleaseLicense, error) {
	if e.loadedTypes[3] {
		return e.Instances, nil
	}
	return nil, &NotLoadedError{edge: "instances"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*License) scanValues(columns []string) ([]interface{}, error) {
	values := make([]interface{}, len(columns))
	for i := range columns {
		switch columns[i] {
		case license.FieldID:
			values[i] = new(sql.NullInt64)
		case license.FieldLicenseID, license.FieldName:
			values[i] = new(sql.NullString)
		case license.ForeignKeys[0]: // license_owner
			values[i] = new(sql.NullInt64)
		case license.ForeignKeys[1]: // license_spdx
			values[i] = new(sql.NullInt64)
		default:
			return nil, fmt.Errorf("unexpected column %q for type License", columns[i])
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the License fields.
func (l *License) assignValues(columns []string, values []interface{}) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case license.FieldID:
			value, ok := values[i].(*sql.NullInt64)
			if !ok {
				return fmt.Errorf("unexpected type %T for field id", value)
			}
			l.ID = int(value.Int64)
		case license.FieldLicenseID:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field license_id", values[i])
			} else if value.Valid {
				l.LicenseID = value.String
			}
		case license.FieldName:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field name", values[i])
			} else if value.Valid {
				l.Name = value.String
			}
		case license.ForeignKeys[0]:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for edge-field license_owner", value)
			} else if value.Valid {
				l.license_owner = new(int)
				*l.license_owner = int(value.Int64)
			}
		case license.ForeignKeys[1]:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for edge-field license_spdx", value)
			} else if value.Valid {
				l.license_spdx = new(int)
				*l.license_spdx = int(value.Int64)
			}
		}
	}
	return nil
}

// QueryOwner queries the "owner" edge of the License entity.
func (l *License) QueryOwner() *OrganizationQuery {
	return (&LicenseClient{config: l.config}).QueryOwner(l)
}

// QuerySpdx queries the "spdx" edge of the License entity.
func (l *License) QuerySpdx() *SPDXLicenseQuery {
	return (&LicenseClient{config: l.config}).QuerySpdx(l)
}

// QueryComponents queries the "components" edge of the License entity.
func (l *License) QueryComponents() *ComponentQuery {
	return (&LicenseClient{config: l.config}).QueryComponents(l)
}

// QueryInstances queries the "instances" edge of the License entity.
func (l *License) QueryInstances() *ReleaseLicenseQuery {
	return (&LicenseClient{config: l.config}).QueryInstances(l)
}

// Update returns a builder for updating this License.
// Note that you need to call License.Unwrap() before calling this method if this License
// was returned from a transaction, and the transaction was committed or rolled back.
func (l *License) Update() *LicenseUpdateOne {
	return (&LicenseClient{config: l.config}).UpdateOne(l)
}

// Unwrap unwraps the License entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (l *License) Unwrap() *License {
	tx, ok := l.config.driver.(*txDriver)
	if !ok {
		panic("ent: License is not a transactional entity")
	}
	l.config.driver = tx.drv
	return l
}

// String implements the fmt.Stringer.
func (l *License) String() string {
	var builder strings.Builder
	builder.WriteString("License(")
	builder.WriteString(fmt.Sprintf("id=%v", l.ID))
	builder.WriteString(", license_id=")
	builder.WriteString(l.LicenseID)
	builder.WriteString(", name=")
	builder.WriteString(l.Name)
	builder.WriteByte(')')
	return builder.String()
}

// Licenses is a parsable slice of License.
type Licenses []*License

func (l Licenses) config(cfg config) {
	for _i := range l {
		l[_i].config = cfg
	}
}
