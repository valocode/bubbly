// Code generated by entc, DO NOT EDIT.

package releasepolicy

const (
	// Label holds the string label denoting the releasepolicy type in the database.
	Label = "release_policy"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldName holds the string denoting the name field in the database.
	FieldName = "name"
	// FieldModule holds the string denoting the module field in the database.
	FieldModule = "module"
	// EdgeOwner holds the string denoting the owner edge name in mutations.
	EdgeOwner = "owner"
	// EdgeViolations holds the string denoting the violations edge name in mutations.
	EdgeViolations = "violations"
	// Table holds the table name of the releasepolicy in the database.
	Table = "release_policy"
	// OwnerTable is the table that holds the owner relation/edge.
	OwnerTable = "release_policy"
	// OwnerInverseTable is the table name for the Organization entity.
	// It exists in this package in order to avoid circular dependency with the "organization" package.
	OwnerInverseTable = "organization"
	// OwnerColumn is the table column denoting the owner relation/edge.
	OwnerColumn = "release_policy_owner"
	// ViolationsTable is the table that holds the violations relation/edge.
	ViolationsTable = "release_policy_violation"
	// ViolationsInverseTable is the table name for the ReleasePolicyViolation entity.
	// It exists in this package in order to avoid circular dependency with the "releasepolicyviolation" package.
	ViolationsInverseTable = "release_policy_violation"
	// ViolationsColumn is the table column denoting the violations relation/edge.
	ViolationsColumn = "release_policy_violation_policy"
)

// Columns holds all SQL columns for releasepolicy fields.
var Columns = []string{
	FieldID,
	FieldName,
	FieldModule,
}

// ForeignKeys holds the SQL foreign-keys that are owned by the "release_policy"
// table and are not defined as standalone fields in the schema.
var ForeignKeys = []string{
	"release_policy_owner",
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
	// NameValidator is a validator for the "name" field. It is called by the builders before save.
	NameValidator func(string) error
	// ModuleValidator is a validator for the "module" field. It is called by the builders before save.
	ModuleValidator func(string) error
)
