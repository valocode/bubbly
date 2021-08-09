// Code generated by entc, DO NOT EDIT.

package artifact

import (
	"fmt"
	"io"
	"strconv"

	"entgo.io/ent"
)

const (
	// Label holds the string label denoting the artifact type in the database.
	Label = "artifact"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldName holds the string denoting the name field in the database.
	FieldName = "name"
	// FieldSha256 holds the string denoting the sha256 field in the database.
	FieldSha256 = "sha256"
	// FieldType holds the string denoting the type field in the database.
	FieldType = "type"
	// EdgeRelease holds the string denoting the release edge name in mutations.
	EdgeRelease = "release"
	// EdgeEntry holds the string denoting the entry edge name in mutations.
	EdgeEntry = "entry"
	// Table holds the table name of the artifact in the database.
	Table = "artifact"
	// ReleaseTable is the table that holds the release relation/edge.
	ReleaseTable = "artifact"
	// ReleaseInverseTable is the table name for the Release entity.
	// It exists in this package in order to avoid circular dependency with the "release" package.
	ReleaseInverseTable = "release"
	// ReleaseColumn is the table column denoting the release relation/edge.
	ReleaseColumn = "artifact_release"
	// EntryTable is the table that holds the entry relation/edge.
	EntryTable = "artifact"
	// EntryInverseTable is the table name for the ReleaseEntry entity.
	// It exists in this package in order to avoid circular dependency with the "releaseentry" package.
	EntryInverseTable = "release_entry"
	// EntryColumn is the table column denoting the entry relation/edge.
	EntryColumn = "release_entry_artifact"
)

// Columns holds all SQL columns for artifact fields.
var Columns = []string{
	FieldID,
	FieldName,
	FieldSha256,
	FieldType,
}

// ForeignKeys holds the SQL foreign-keys that are owned by the "artifact"
// table and are not defined as standalone fields in the schema.
var ForeignKeys = []string{
	"artifact_release",
	"release_entry_artifact",
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
	// NameValidator is a validator for the "name" field. It is called by the builders before save.
	NameValidator func(string) error
	// Sha256Validator is a validator for the "sha256" field. It is called by the builders before save.
	Sha256Validator func(string) error
)

// Type defines the type for the "type" enum field.
type Type string

// Type values.
const (
	TypeDocker Type = "docker"
	TypeFile   Type = "file"
)

func (_type Type) String() string {
	return string(_type)
}

// TypeValidator is a validator for the "type" field enum values. It is called by the builders before save.
func TypeValidator(_type Type) error {
	switch _type {
	case TypeDocker, TypeFile:
		return nil
	default:
		return fmt.Errorf("artifact: invalid enum value for type field: %q", _type)
	}
}

// MarshalGQL implements graphql.Marshaler interface.
func (_type Type) MarshalGQL(w io.Writer) {
	io.WriteString(w, strconv.Quote(_type.String()))
}

// UnmarshalGQL implements graphql.Unmarshaler interface.
func (_type *Type) UnmarshalGQL(val interface{}) error {
	str, ok := val.(string)
	if !ok {
		return fmt.Errorf("enum %T must be a string", val)
	}
	*_type = Type(str)
	if err := TypeValidator(*_type); err != nil {
		return fmt.Errorf("%s is not a valid Type", str)
	}
	return nil
}
