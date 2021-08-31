// Code generated by entc, DO NOT EDIT.

package ent

import (
	"encoding/json"
	"fmt"
	"strings"

	"entgo.io/ent/dialect/sql"
	schema "github.com/valocode/bubbly/ent/schema/types"
	"github.com/valocode/bubbly/ent/testcase"
	"github.com/valocode/bubbly/ent/testrun"
)

// TestCase is the model entity for the TestCase schema.
type TestCase struct {
	config `json:"-"`
	// ID of the ent.
	ID int `json:"id,omitempty"`
	// Name holds the value of the "name" field.
	Name string `json:"name,omitempty"`
	// Result holds the value of the "result" field.
	Result bool `json:"result,omitempty"`
	// Message holds the value of the "message" field.
	Message string `json:"message,omitempty"`
	// Elapsed holds the value of the "elapsed" field.
	Elapsed float64 `json:"elapsed,omitempty"`
	// Metadata holds the value of the "metadata" field.
	Metadata schema.Metadata `json:"metadata,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the TestCaseQuery when eager-loading is set.
	Edges         TestCaseEdges `json:"edges"`
	test_case_run *int
}

// TestCaseEdges holds the relations/edges for other nodes in the graph.
type TestCaseEdges struct {
	// Run holds the value of the run edge.
	Run *TestRun `json:"run,omitempty"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [1]bool
}

// RunOrErr returns the Run value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e TestCaseEdges) RunOrErr() (*TestRun, error) {
	if e.loadedTypes[0] {
		if e.Run == nil {
			// The edge run was loaded in eager-loading,
			// but was not found.
			return nil, &NotFoundError{label: testrun.Label}
		}
		return e.Run, nil
	}
	return nil, &NotLoadedError{edge: "run"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*TestCase) scanValues(columns []string) ([]interface{}, error) {
	values := make([]interface{}, len(columns))
	for i := range columns {
		switch columns[i] {
		case testcase.FieldMetadata:
			values[i] = new([]byte)
		case testcase.FieldResult:
			values[i] = new(sql.NullBool)
		case testcase.FieldElapsed:
			values[i] = new(sql.NullFloat64)
		case testcase.FieldID:
			values[i] = new(sql.NullInt64)
		case testcase.FieldName, testcase.FieldMessage:
			values[i] = new(sql.NullString)
		case testcase.ForeignKeys[0]: // test_case_run
			values[i] = new(sql.NullInt64)
		default:
			return nil, fmt.Errorf("unexpected column %q for type TestCase", columns[i])
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the TestCase fields.
func (tc *TestCase) assignValues(columns []string, values []interface{}) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case testcase.FieldID:
			value, ok := values[i].(*sql.NullInt64)
			if !ok {
				return fmt.Errorf("unexpected type %T for field id", value)
			}
			tc.ID = int(value.Int64)
		case testcase.FieldName:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field name", values[i])
			} else if value.Valid {
				tc.Name = value.String
			}
		case testcase.FieldResult:
			if value, ok := values[i].(*sql.NullBool); !ok {
				return fmt.Errorf("unexpected type %T for field result", values[i])
			} else if value.Valid {
				tc.Result = value.Bool
			}
		case testcase.FieldMessage:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field message", values[i])
			} else if value.Valid {
				tc.Message = value.String
			}
		case testcase.FieldElapsed:
			if value, ok := values[i].(*sql.NullFloat64); !ok {
				return fmt.Errorf("unexpected type %T for field elapsed", values[i])
			} else if value.Valid {
				tc.Elapsed = value.Float64
			}
		case testcase.FieldMetadata:
			if value, ok := values[i].(*[]byte); !ok {
				return fmt.Errorf("unexpected type %T for field metadata", values[i])
			} else if value != nil && len(*value) > 0 {
				if err := json.Unmarshal(*value, &tc.Metadata); err != nil {
					return fmt.Errorf("unmarshal field metadata: %w", err)
				}
			}
		case testcase.ForeignKeys[0]:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for edge-field test_case_run", value)
			} else if value.Valid {
				tc.test_case_run = new(int)
				*tc.test_case_run = int(value.Int64)
			}
		}
	}
	return nil
}

// QueryRun queries the "run" edge of the TestCase entity.
func (tc *TestCase) QueryRun() *TestRunQuery {
	return (&TestCaseClient{config: tc.config}).QueryRun(tc)
}

// Update returns a builder for updating this TestCase.
// Note that you need to call TestCase.Unwrap() before calling this method if this TestCase
// was returned from a transaction, and the transaction was committed or rolled back.
func (tc *TestCase) Update() *TestCaseUpdateOne {
	return (&TestCaseClient{config: tc.config}).UpdateOne(tc)
}

// Unwrap unwraps the TestCase entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (tc *TestCase) Unwrap() *TestCase {
	tx, ok := tc.config.driver.(*txDriver)
	if !ok {
		panic("ent: TestCase is not a transactional entity")
	}
	tc.config.driver = tx.drv
	return tc
}

// String implements the fmt.Stringer.
func (tc *TestCase) String() string {
	var builder strings.Builder
	builder.WriteString("TestCase(")
	builder.WriteString(fmt.Sprintf("id=%v", tc.ID))
	builder.WriteString(", name=")
	builder.WriteString(tc.Name)
	builder.WriteString(", result=")
	builder.WriteString(fmt.Sprintf("%v", tc.Result))
	builder.WriteString(", message=")
	builder.WriteString(tc.Message)
	builder.WriteString(", elapsed=")
	builder.WriteString(fmt.Sprintf("%v", tc.Elapsed))
	builder.WriteString(", metadata=")
	builder.WriteString(fmt.Sprintf("%v", tc.Metadata))
	builder.WriteByte(')')
	return builder.String()
}

// TestCases is a parsable slice of TestCase.
type TestCases []*TestCase

func (tc TestCases) config(cfg config) {
	for _i := range tc {
		tc[_i].config = cfg
	}
}
