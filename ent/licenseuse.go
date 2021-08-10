// Code generated by entc, DO NOT EDIT.

package ent

import (
	"fmt"
	"strings"

	"entgo.io/ent/dialect/sql"
	"github.com/valocode/bubbly/ent/license"
	"github.com/valocode/bubbly/ent/licenseuse"
)

// LicenseUse is the model entity for the LicenseUse schema.
type LicenseUse struct {
	config
	// ID of the ent.
	ID int `json:"id,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the LicenseUseQuery when eager-loading is set.
	Edges               LicenseUseEdges `json:"edges"`
	license_use_license *int
}

// LicenseUseEdges holds the relations/edges for other nodes in the graph.
type LicenseUseEdges struct {
	// License holds the value of the license edge.
	License *License `json:"license,omitempty"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [1]bool
}

// LicenseOrErr returns the License value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e LicenseUseEdges) LicenseOrErr() (*License, error) {
	if e.loadedTypes[0] {
		if e.License == nil {
			// The edge license was loaded in eager-loading,
			// but was not found.
			return nil, &NotFoundError{label: license.Label}
		}
		return e.License, nil
	}
	return nil, &NotLoadedError{edge: "license"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*LicenseUse) scanValues(columns []string) ([]interface{}, error) {
	values := make([]interface{}, len(columns))
	for i := range columns {
		switch columns[i] {
		case licenseuse.FieldID:
			values[i] = new(sql.NullInt64)
		case licenseuse.ForeignKeys[0]: // license_use_license
			values[i] = new(sql.NullInt64)
		default:
			return nil, fmt.Errorf("unexpected column %q for type LicenseUse", columns[i])
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the LicenseUse fields.
func (lu *LicenseUse) assignValues(columns []string, values []interface{}) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case licenseuse.FieldID:
			value, ok := values[i].(*sql.NullInt64)
			if !ok {
				return fmt.Errorf("unexpected type %T for field id", value)
			}
			lu.ID = int(value.Int64)
		case licenseuse.ForeignKeys[0]:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for edge-field license_use_license", value)
			} else if value.Valid {
				lu.license_use_license = new(int)
				*lu.license_use_license = int(value.Int64)
			}
		}
	}
	return nil
}

// QueryLicense queries the "license" edge of the LicenseUse entity.
func (lu *LicenseUse) QueryLicense() *LicenseQuery {
	return (&LicenseUseClient{config: lu.config}).QueryLicense(lu)
}

// Update returns a builder for updating this LicenseUse.
// Note that you need to call LicenseUse.Unwrap() before calling this method if this LicenseUse
// was returned from a transaction, and the transaction was committed or rolled back.
func (lu *LicenseUse) Update() *LicenseUseUpdateOne {
	return (&LicenseUseClient{config: lu.config}).UpdateOne(lu)
}

// Unwrap unwraps the LicenseUse entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (lu *LicenseUse) Unwrap() *LicenseUse {
	tx, ok := lu.config.driver.(*txDriver)
	if !ok {
		panic("ent: LicenseUse is not a transactional entity")
	}
	lu.config.driver = tx.drv
	return lu
}

// String implements the fmt.Stringer.
func (lu *LicenseUse) String() string {
	var builder strings.Builder
	builder.WriteString("LicenseUse(")
	builder.WriteString(fmt.Sprintf("id=%v", lu.ID))
	builder.WriteByte(')')
	return builder.String()
}

// LicenseUses is a parsable slice of LicenseUse.
type LicenseUses []*LicenseUse

func (lu LicenseUses) config(cfg config) {
	for _i := range lu {
		lu[_i].config = cfg
	}
}