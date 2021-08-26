// Code generated by entc, DO NOT EDIT.

package codeissue

import (
	"fmt"
	"io"
	"strconv"
)

const (
	// Label holds the string label denoting the codeissue type in the database.
	Label = "code_issue"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldRuleID holds the string denoting the rule_id field in the database.
	FieldRuleID = "rule_id"
	// FieldMessage holds the string denoting the message field in the database.
	FieldMessage = "message"
	// FieldSeverity holds the string denoting the severity field in the database.
	FieldSeverity = "severity"
	// FieldType holds the string denoting the type field in the database.
	FieldType = "type"
	// EdgeScan holds the string denoting the scan edge name in mutations.
	EdgeScan = "scan"
	// Table holds the table name of the codeissue in the database.
	Table = "code_issue"
	// ScanTable is the table that holds the scan relation/edge.
	ScanTable = "code_issue"
	// ScanInverseTable is the table name for the CodeScan entity.
	// It exists in this package in order to avoid circular dependency with the "codescan" package.
	ScanInverseTable = "code_scan"
	// ScanColumn is the table column denoting the scan relation/edge.
	ScanColumn = "code_issue_scan"
)

// Columns holds all SQL columns for codeissue fields.
var Columns = []string{
	FieldID,
	FieldRuleID,
	FieldMessage,
	FieldSeverity,
	FieldType,
}

// ForeignKeys holds the SQL foreign-keys that are owned by the "code_issue"
// table and are not defined as standalone fields in the schema.
var ForeignKeys = []string{
	"code_issue_scan",
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
	// RuleIDValidator is a validator for the "rule_id" field. It is called by the builders before save.
	RuleIDValidator func(string) error
	// MessageValidator is a validator for the "message" field. It is called by the builders before save.
	MessageValidator func(string) error
)

// Severity defines the type for the "severity" enum field.
type Severity string

// Severity values.
const (
	SeverityLow    Severity = "low"
	SeverityMedium Severity = "medium"
	SeverityHigh   Severity = "high"
)

func (s Severity) String() string {
	return string(s)
}

// SeverityValidator is a validator for the "severity" field enum values. It is called by the builders before save.
func SeverityValidator(s Severity) error {
	switch s {
	case SeverityLow, SeverityMedium, SeverityHigh:
		return nil
	default:
		return fmt.Errorf("codeissue: invalid enum value for severity field: %q", s)
	}
}

// Type defines the type for the "type" enum field.
type Type string

// Type values.
const (
	TypeStyle    Type = "style"
	TypeSecurity Type = "security"
	TypeBug      Type = "bug"
)

func (_type Type) String() string {
	return string(_type)
}

// TypeValidator is a validator for the "type" field enum values. It is called by the builders before save.
func TypeValidator(_type Type) error {
	switch _type {
	case TypeStyle, TypeSecurity, TypeBug:
		return nil
	default:
		return fmt.Errorf("codeissue: invalid enum value for type field: %q", _type)
	}
}

// MarshalGQL implements graphql.Marshaler interface.
func (s Severity) MarshalGQL(w io.Writer) {
	io.WriteString(w, strconv.Quote(s.String()))
}

// UnmarshalGQL implements graphql.Unmarshaler interface.
func (s *Severity) UnmarshalGQL(val interface{}) error {
	str, ok := val.(string)
	if !ok {
		return fmt.Errorf("enum %T must be a string", val)
	}
	*s = Severity(str)
	if err := SeverityValidator(*s); err != nil {
		return fmt.Errorf("%s is not a valid Severity", str)
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
