// Code generated by entc, DO NOT EDIT.

package testrun

import (
	"time"

	"entgo.io/ent"
)

const (
	// Label holds the string label denoting the testrun type in the database.
	Label = "test_run"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldTool holds the string denoting the tool field in the database.
	FieldTool = "tool"
	// FieldTime holds the string denoting the time field in the database.
	FieldTime = "time"
	// FieldMetadata holds the string denoting the metadata field in the database.
	FieldMetadata = "metadata"
	// EdgeRelease holds the string denoting the release edge name in mutations.
	EdgeRelease = "release"
	// EdgeEntry holds the string denoting the entry edge name in mutations.
	EdgeEntry = "entry"
	// EdgeTests holds the string denoting the tests edge name in mutations.
	EdgeTests = "tests"
	// Table holds the table name of the testrun in the database.
	Table = "test_run"
	// ReleaseTable is the table that holds the release relation/edge.
	ReleaseTable = "test_run"
	// ReleaseInverseTable is the table name for the Release entity.
	// It exists in this package in order to avoid circular dependency with the "release" package.
	ReleaseInverseTable = "release"
	// ReleaseColumn is the table column denoting the release relation/edge.
	ReleaseColumn = "test_run_release"
	// EntryTable is the table that holds the entry relation/edge.
	EntryTable = "test_run"
	// EntryInverseTable is the table name for the ReleaseEntry entity.
	// It exists in this package in order to avoid circular dependency with the "releaseentry" package.
	EntryInverseTable = "release_entry"
	// EntryColumn is the table column denoting the entry relation/edge.
	EntryColumn = "release_entry_test_run"
	// TestsTable is the table that holds the tests relation/edge.
	TestsTable = "test_case"
	// TestsInverseTable is the table name for the TestCase entity.
	// It exists in this package in order to avoid circular dependency with the "testcase" package.
	TestsInverseTable = "test_case"
	// TestsColumn is the table column denoting the tests relation/edge.
	TestsColumn = "test_case_run"
)

// Columns holds all SQL columns for testrun fields.
var Columns = []string{
	FieldID,
	FieldTool,
	FieldTime,
	FieldMetadata,
}

// ForeignKeys holds the SQL foreign-keys that are owned by the "test_run"
// table and are not defined as standalone fields in the schema.
var ForeignKeys = []string{
	"release_entry_test_run",
	"test_run_release",
}

// ValidColumn reports if the column name is valid (part of the table columns).
func ValidColumn(column string) bool {
	for i := range Columns {
		if column == Columns[i] {
			return true
		}
	}
	for i := range ForeignKeys {
		if column == ForeignKeys[i] {
			return true
		}
	}
	return false
}

// Note that the variables below are initialized by the runtime
// package on the initialization of the application. Therefore,
// it should be imported in the main as follows:
//
//	import _ "github.com/valocode/bubbly/ent/runtime"
//
var (
	Hooks [1]ent.Hook
	// ToolValidator is a validator for the "tool" field. It is called by the builders before save.
	ToolValidator func(string) error
	// DefaultTime holds the default value on creation for the "time" field.
	DefaultTime func() time.Time
)
