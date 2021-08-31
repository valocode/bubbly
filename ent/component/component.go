// Code generated by entc, DO NOT EDIT.

package component

const (
	// Label holds the string label denoting the component type in the database.
	Label = "component"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldName holds the string denoting the name field in the database.
	FieldName = "name"
	// FieldVendor holds the string denoting the vendor field in the database.
	FieldVendor = "vendor"
	// FieldVersion holds the string denoting the version field in the database.
	FieldVersion = "version"
	// FieldDescription holds the string denoting the description field in the database.
	FieldDescription = "description"
	// FieldURL holds the string denoting the url field in the database.
	FieldURL = "url"
	// FieldMetadata holds the string denoting the metadata field in the database.
	FieldMetadata = "metadata"
	// EdgeVulnerabilities holds the string denoting the vulnerabilities edge name in mutations.
	EdgeVulnerabilities = "vulnerabilities"
	// EdgeLicenses holds the string denoting the licenses edge name in mutations.
	EdgeLicenses = "licenses"
	// EdgeUses holds the string denoting the uses edge name in mutations.
	EdgeUses = "uses"
	// Table holds the table name of the component in the database.
	Table = "component"
	// VulnerabilitiesTable is the table that holds the vulnerabilities relation/edge. The primary key declared below.
	VulnerabilitiesTable = "component_vulnerabilities"
	// VulnerabilitiesInverseTable is the table name for the Vulnerability entity.
	// It exists in this package in order to avoid circular dependency with the "vulnerability" package.
	VulnerabilitiesInverseTable = "vulnerability"
	// LicensesTable is the table that holds the licenses relation/edge. The primary key declared below.
	LicensesTable = "component_licenses"
	// LicensesInverseTable is the table name for the License entity.
	// It exists in this package in order to avoid circular dependency with the "license" package.
	LicensesInverseTable = "license"
	// UsesTable is the table that holds the uses relation/edge.
	UsesTable = "release_component"
	// UsesInverseTable is the table name for the ReleaseComponent entity.
	// It exists in this package in order to avoid circular dependency with the "releasecomponent" package.
	UsesInverseTable = "release_component"
	// UsesColumn is the table column denoting the uses relation/edge.
	UsesColumn = "release_component_component"
)

// Columns holds all SQL columns for component fields.
var Columns = []string{
	FieldID,
	FieldName,
	FieldVendor,
	FieldVersion,
	FieldDescription,
	FieldURL,
	FieldMetadata,
}

var (
	// VulnerabilitiesPrimaryKey and VulnerabilitiesColumn2 are the table columns denoting the
	// primary key for the vulnerabilities relation (M2M).
	VulnerabilitiesPrimaryKey = []string{"component_id", "vulnerability_id"}
	// LicensesPrimaryKey and LicensesColumn2 are the table columns denoting the
	// primary key for the licenses relation (M2M).
	LicensesPrimaryKey = []string{"component_id", "license_id"}
)

// ValidColumn reports if the column name is valid (part of the table columns).
func ValidColumn(column string) bool {
	for i := range Columns {
		if column == Columns[i] {
			return true
		}
	}
	return false
}

var (
	// NameValidator is a validator for the "name" field. It is called by the builders before save.
	NameValidator func(string) error
	// DefaultVendor holds the default value on creation for the "vendor" field.
	DefaultVendor string
	// VersionValidator is a validator for the "version" field. It is called by the builders before save.
	VersionValidator func(string) error
)
