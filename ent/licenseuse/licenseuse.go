// Code generated by entc, DO NOT EDIT.

package licenseuse

const (
	// Label holds the string label denoting the licenseuse type in the database.
	Label = "license_use"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// EdgeLicense holds the string denoting the license edge name in mutations.
	EdgeLicense = "license"
	// Table holds the table name of the licenseuse in the database.
	Table = "license_use"
	// LicenseTable is the table that holds the license relation/edge.
	LicenseTable = "license_use"
	// LicenseInverseTable is the table name for the License entity.
	// It exists in this package in order to avoid circular dependency with the "license" package.
	LicenseInverseTable = "license"
	// LicenseColumn is the table column denoting the license relation/edge.
	LicenseColumn = "license_use_license"
)

// Columns holds all SQL columns for licenseuse fields.
var Columns = []string{
	FieldID,
}

// ForeignKeys holds the SQL foreign-keys that are owned by the "license_use"
// table and are not defined as standalone fields in the schema.
var ForeignKeys = []string{
	"license_use_license",
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