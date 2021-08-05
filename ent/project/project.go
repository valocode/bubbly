// Code generated by entc, DO NOT EDIT.

package project

const (
	// Label holds the string label denoting the project type in the database.
	Label = "project"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldName holds the string denoting the name field in the database.
	FieldName = "name"
	// EdgeRepos holds the string denoting the repos edge name in mutations.
	EdgeRepos = "repos"
	// EdgeReleases holds the string denoting the releases edge name in mutations.
	EdgeReleases = "releases"
	// EdgeCveRules holds the string denoting the cve_rules edge name in mutations.
	EdgeCveRules = "cve_rules"
	// Table holds the table name of the project in the database.
	Table = "project"
	// ReposTable is the table the holds the repos relation/edge.
	ReposTable = "repo"
	// ReposInverseTable is the table name for the Repo entity.
	// It exists in this package in order to avoid circular dependency with the "repo" package.
	ReposInverseTable = "repo"
	// ReposColumn is the table column denoting the repos relation/edge.
	ReposColumn = "repo_project"
	// ReleasesTable is the table the holds the releases relation/edge.
	ReleasesTable = "release"
	// ReleasesInverseTable is the table name for the Release entity.
	// It exists in this package in order to avoid circular dependency with the "release" package.
	ReleasesInverseTable = "release"
	// ReleasesColumn is the table column denoting the releases relation/edge.
	ReleasesColumn = "release_project"
	// CveRulesTable is the table the holds the cve_rules relation/edge. The primary key declared below.
	CveRulesTable = "cve_rule_project"
	// CveRulesInverseTable is the table name for the CVERule entity.
	// It exists in this package in order to avoid circular dependency with the "cverule" package.
	CveRulesInverseTable = "cve_rule"
)

// Columns holds all SQL columns for project fields.
var Columns = []string{
	FieldID,
	FieldName,
}

var (
	// CveRulesPrimaryKey and CveRulesColumn2 are the table columns denoting the
	// primary key for the cve_rules relation (M2M).
	CveRulesPrimaryKey = []string{"cve_rule_id", "project_id"}
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
)
