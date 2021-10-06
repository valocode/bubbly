// Code generated by entc, DO NOT EDIT.

package event

import (
	"fmt"
	"io"
	"strconv"
	"time"
)

const (
	// Label holds the string label denoting the event type in the database.
	Label = "event"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldMessage holds the string denoting the message field in the database.
	FieldMessage = "message"
	// FieldStatus holds the string denoting the status field in the database.
	FieldStatus = "status"
	// FieldType holds the string denoting the type field in the database.
	FieldType = "type"
	// FieldTime holds the string denoting the time field in the database.
	FieldTime = "time"
	// EdgeRelease holds the string denoting the release edge name in mutations.
	EdgeRelease = "release"
	// EdgeRepo holds the string denoting the repo edge name in mutations.
	EdgeRepo = "repo"
	// EdgeProject holds the string denoting the project edge name in mutations.
	EdgeProject = "project"
	// Table holds the table name of the event in the database.
	Table = "event"
	// ReleaseTable is the table that holds the release relation/edge.
	ReleaseTable = "event"
	// ReleaseInverseTable is the table name for the Release entity.
	// It exists in this package in order to avoid circular dependency with the "release" package.
	ReleaseInverseTable = "release"
	// ReleaseColumn is the table column denoting the release relation/edge.
	ReleaseColumn = "event_release"
	// RepoTable is the table that holds the repo relation/edge.
	RepoTable = "event"
	// RepoInverseTable is the table name for the Repo entity.
	// It exists in this package in order to avoid circular dependency with the "repo" package.
	RepoInverseTable = "repo"
	// RepoColumn is the table column denoting the repo relation/edge.
	RepoColumn = "event_repo"
	// ProjectTable is the table that holds the project relation/edge.
	ProjectTable = "event"
	// ProjectInverseTable is the table name for the Project entity.
	// It exists in this package in order to avoid circular dependency with the "project" package.
	ProjectInverseTable = "project"
	// ProjectColumn is the table column denoting the project relation/edge.
	ProjectColumn = "event_project"
)

// Columns holds all SQL columns for event fields.
var Columns = []string{
	FieldID,
	FieldMessage,
	FieldStatus,
	FieldType,
	FieldTime,
}

// ForeignKeys holds the SQL foreign-keys that are owned by the "event"
// table and are not defined as standalone fields in the schema.
var ForeignKeys = []string{
	"event_release",
	"event_repo",
	"event_project",
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

var (
	// DefaultMessage holds the default value on creation for the "message" field.
	DefaultMessage string
	// DefaultTime holds the default value on creation for the "time" field.
	DefaultTime func() time.Time
)

// Status defines the type for the "status" enum field.
type Status string

// StatusOk is the default value of the Status enum.
const DefaultStatus = StatusOk

// Status values.
const (
	StatusOk    Status = "ok"
	StatusError Status = "error"
)

func (s Status) String() string {
	return string(s)
}

// StatusValidator is a validator for the "status" field enum values. It is called by the builders before save.
func StatusValidator(s Status) error {
	switch s {
	case StatusOk, StatusError:
		return nil
	default:
		return fmt.Errorf("event: invalid enum value for status field: %q", s)
	}
}

// Type defines the type for the "type" enum field.
type Type string

// Type values.
const (
	TypeEvaluateRelease Type = "evaluate_release"
	TypeMonitor         Type = "monitor"
)

func (_type Type) String() string {
	return string(_type)
}

// TypeValidator is a validator for the "type" field enum values. It is called by the builders before save.
func TypeValidator(_type Type) error {
	switch _type {
	case TypeEvaluateRelease, TypeMonitor:
		return nil
	default:
		return fmt.Errorf("event: invalid enum value for type field: %q", _type)
	}
}

// MarshalGQL implements graphql.Marshaler interface.
func (s Status) MarshalGQL(w io.Writer) {
	io.WriteString(w, strconv.Quote(s.String()))
}

// UnmarshalGQL implements graphql.Unmarshaler interface.
func (s *Status) UnmarshalGQL(val interface{}) error {
	str, ok := val.(string)
	if !ok {
		return fmt.Errorf("enum %T must be a string", val)
	}
	*s = Status(str)
	if err := StatusValidator(*s); err != nil {
		return fmt.Errorf("%s is not a valid Status", str)
	}
	return nil
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