// Code generated by entc, DO NOT EDIT.

package ent

import (
	"fmt"
	"strings"

	"entgo.io/ent/dialect/sql"
	"github.com/valocode/bubbly/ent/codeissue"
	"github.com/valocode/bubbly/ent/codescan"
)

// CodeIssue is the model entity for the CodeIssue schema.
type CodeIssue struct {
	config `json:"-"`
	// ID of the ent.
	ID int `json:"id,omitempty"`
	// RuleID holds the value of the "rule_id" field.
	RuleID string `json:"rule_id,omitempty"`
	// Message holds the value of the "message" field.
	Message string `json:"message,omitempty"`
	// Severity holds the value of the "severity" field.
	Severity codeissue.Severity `json:"severity,omitempty"`
	// Type holds the value of the "type" field.
	Type codeissue.Type `json:"type,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the CodeIssueQuery when eager-loading is set.
	Edges           CodeIssueEdges `json:"edges"`
	code_issue_scan *int
}

// CodeIssueEdges holds the relations/edges for other nodes in the graph.
type CodeIssueEdges struct {
	// Cwe holds the value of the cwe edge.
	Cwe []*CWE `json:"cwe,omitempty"`
	// Scan holds the value of the scan edge.
	Scan *CodeScan `json:"scan,omitempty"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [2]bool
}

// CweOrErr returns the Cwe value or an error if the edge
// was not loaded in eager-loading.
func (e CodeIssueEdges) CweOrErr() ([]*CWE, error) {
	if e.loadedTypes[0] {
		return e.Cwe, nil
	}
	return nil, &NotLoadedError{edge: "cwe"}
}

// ScanOrErr returns the Scan value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e CodeIssueEdges) ScanOrErr() (*CodeScan, error) {
	if e.loadedTypes[1] {
		if e.Scan == nil {
			// The edge scan was loaded in eager-loading,
			// but was not found.
			return nil, &NotFoundError{label: codescan.Label}
		}
		return e.Scan, nil
	}
	return nil, &NotLoadedError{edge: "scan"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*CodeIssue) scanValues(columns []string) ([]interface{}, error) {
	values := make([]interface{}, len(columns))
	for i := range columns {
		switch columns[i] {
		case codeissue.FieldID:
			values[i] = new(sql.NullInt64)
		case codeissue.FieldRuleID, codeissue.FieldMessage, codeissue.FieldSeverity, codeissue.FieldType:
			values[i] = new(sql.NullString)
		case codeissue.ForeignKeys[0]: // code_issue_scan
			values[i] = new(sql.NullInt64)
		default:
			return nil, fmt.Errorf("unexpected column %q for type CodeIssue", columns[i])
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the CodeIssue fields.
func (ci *CodeIssue) assignValues(columns []string, values []interface{}) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case codeissue.FieldID:
			value, ok := values[i].(*sql.NullInt64)
			if !ok {
				return fmt.Errorf("unexpected type %T for field id", value)
			}
			ci.ID = int(value.Int64)
		case codeissue.FieldRuleID:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field rule_id", values[i])
			} else if value.Valid {
				ci.RuleID = value.String
			}
		case codeissue.FieldMessage:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field message", values[i])
			} else if value.Valid {
				ci.Message = value.String
			}
		case codeissue.FieldSeverity:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field severity", values[i])
			} else if value.Valid {
				ci.Severity = codeissue.Severity(value.String)
			}
		case codeissue.FieldType:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field type", values[i])
			} else if value.Valid {
				ci.Type = codeissue.Type(value.String)
			}
		case codeissue.ForeignKeys[0]:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for edge-field code_issue_scan", value)
			} else if value.Valid {
				ci.code_issue_scan = new(int)
				*ci.code_issue_scan = int(value.Int64)
			}
		}
	}
	return nil
}

// QueryCwe queries the "cwe" edge of the CodeIssue entity.
func (ci *CodeIssue) QueryCwe() *CWEQuery {
	return (&CodeIssueClient{config: ci.config}).QueryCwe(ci)
}

// QueryScan queries the "scan" edge of the CodeIssue entity.
func (ci *CodeIssue) QueryScan() *CodeScanQuery {
	return (&CodeIssueClient{config: ci.config}).QueryScan(ci)
}

// Update returns a builder for updating this CodeIssue.
// Note that you need to call CodeIssue.Unwrap() before calling this method if this CodeIssue
// was returned from a transaction, and the transaction was committed or rolled back.
func (ci *CodeIssue) Update() *CodeIssueUpdateOne {
	return (&CodeIssueClient{config: ci.config}).UpdateOne(ci)
}

// Unwrap unwraps the CodeIssue entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (ci *CodeIssue) Unwrap() *CodeIssue {
	tx, ok := ci.config.driver.(*txDriver)
	if !ok {
		panic("ent: CodeIssue is not a transactional entity")
	}
	ci.config.driver = tx.drv
	return ci
}

// String implements the fmt.Stringer.
func (ci *CodeIssue) String() string {
	var builder strings.Builder
	builder.WriteString("CodeIssue(")
	builder.WriteString(fmt.Sprintf("id=%v", ci.ID))
	builder.WriteString(", rule_id=")
	builder.WriteString(ci.RuleID)
	builder.WriteString(", message=")
	builder.WriteString(ci.Message)
	builder.WriteString(", severity=")
	builder.WriteString(fmt.Sprintf("%v", ci.Severity))
	builder.WriteString(", type=")
	builder.WriteString(fmt.Sprintf("%v", ci.Type))
	builder.WriteByte(')')
	return builder.String()
}

// CodeIssues is a parsable slice of CodeIssue.
type CodeIssues []*CodeIssue

func (ci CodeIssues) config(cfg config) {
	for _i := range ci {
		ci[_i].config = cfg
	}
}