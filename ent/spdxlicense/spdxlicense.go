// Code generated by entc, DO NOT EDIT.

package spdxlicense

import (
	"entgo.io/ent"
)

const (
	// Label holds the string label denoting the spdxlicense type in the database.
	Label = "spdx_license"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldLicenseID holds the string denoting the license_id field in the database.
	FieldLicenseID = "license_id"
	// FieldName holds the string denoting the name field in the database.
	FieldName = "name"
	// FieldReference holds the string denoting the reference field in the database.
	FieldReference = "reference"
	// FieldDetailsURL holds the string denoting the details_url field in the database.
	FieldDetailsURL = "details_url"
	// FieldIsOsiApproved holds the string denoting the is_osi_approved field in the database.
	FieldIsOsiApproved = "is_osi_approved"
	// Table holds the table name of the spdxlicense in the database.
	Table = "spdx_license"
)

// Columns holds all SQL columns for spdxlicense fields.
var Columns = []string{
	FieldID,
	FieldLicenseID,
	FieldName,
	FieldReference,
	FieldDetailsURL,
	FieldIsOsiApproved,
}

// ValidColumn reports if the column name is valid (part of the table columns).
func ValidColumn(column string) bool {
	for i := range Columns {
		if column == Columns[i] {
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
	// LicenseIDValidator is a validator for the "license_id" field. It is called by the builders before save.
	LicenseIDValidator func(string) error
	// NameValidator is a validator for the "name" field. It is called by the builders before save.
	NameValidator func(string) error
	// DefaultIsOsiApproved holds the default value on creation for the "is_osi_approved" field.
	DefaultIsOsiApproved bool
)
