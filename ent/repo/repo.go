// Code generated by entc, DO NOT EDIT.

package repo

const (
	// Label holds the string label denoting the repo type in the database.
	Label = "repo"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldName holds the string denoting the name field in the database.
	FieldName = "name"
	// EdgeProject holds the string denoting the project edge name in mutations.
	EdgeProject = "project"
	// EdgeCommits holds the string denoting the commits edge name in mutations.
	EdgeCommits = "commits"
	// EdgeCveRules holds the string denoting the cve_rules edge name in mutations.
	EdgeCveRules = "cve_rules"
	// Table holds the table name of the repo in the database.
	Table = "repo"
	// ProjectTable is the table that holds the project relation/edge.
	ProjectTable = "repo"
	// ProjectInverseTable is the table name for the Project entity.
	// It exists in this package in order to avoid circular dependency with the "project" package.
	ProjectInverseTable = "project"
	// ProjectColumn is the table column denoting the project relation/edge.
	ProjectColumn = "repo_project"
	// CommitsTable is the table that holds the commits relation/edge.
	CommitsTable = "commit"
	// CommitsInverseTable is the table name for the GitCommit entity.
	// It exists in this package in order to avoid circular dependency with the "gitcommit" package.
	CommitsInverseTable = "commit"
	// CommitsColumn is the table column denoting the commits relation/edge.
	CommitsColumn = "git_commit_repo"
	// CveRulesTable is the table that holds the cve_rules relation/edge. The primary key declared below.
	CveRulesTable = "cve_rule_repo"
	// CveRulesInverseTable is the table name for the CVERule entity.
	// It exists in this package in order to avoid circular dependency with the "cverule" package.
	CveRulesInverseTable = "cve_rule"
)

// Columns holds all SQL columns for repo fields.
var Columns = []string{
	FieldID,
	FieldName,
}

// ForeignKeys holds the SQL foreign-keys that are owned by the "repo"
// table and are not defined as standalone fields in the schema.
var ForeignKeys = []string{
	"repo_project",
}

var (
	// CveRulesPrimaryKey and CveRulesColumn2 are the table columns denoting the
	// primary key for the cve_rules relation (M2M).
	CveRulesPrimaryKey = []string{"cve_rule_id", "repo_id"}
)

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
)