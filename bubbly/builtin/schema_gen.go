package builtin

import "github.com/valocode/bubbly/api/core"

// #######################################
// _SCHEMA
// #######################################
type Schema struct {
	DBlock_Table  string                 `json:"_schema,omitempty"`
	DBlock_Policy core.DataBlockPolicy   `json:"-"`
	DBlock_Joins  []string               `json:"-"`
	Tables        map[string]interface{} `json:"tables,omitempty"`
}
type Schema_Wrap struct {
	Schema []Schema `json:"_schema,omitempty"`
}

// #######################################
// _RESOURCE
// #######################################
type Resource struct {
	DBlock_Table  string                 `json:"_resource,omitempty"`
	DBlock_Policy core.DataBlockPolicy   `json:"-"`
	DBlock_Joins  []string               `json:"-"`
	Id            string                 `json:"id,omitempty"`
	Name          string                 `json:"name,omitempty"`
	Kind          string                 `json:"kind,omitempty"`
	ApiVersion    string                 `json:"api_version,omitempty"`
	Spec          string                 `json:"spec,omitempty"`
	Metadata      map[string]interface{} `json:"metadata,omitempty"`
	Event         []Event                `json:"_event,omitempty"`
	ReleaseEntry  []ReleaseEntry         `json:"release_entry,omitempty"`
}
type Resource_Wrap struct {
	Resource []Resource `json:"_resource,omitempty"`
}

// #######################################
// _EVENT
// #######################################
type Event struct {
	DBlock_Table  string               `json:"_event,omitempty"`
	DBlock_Policy core.DataBlockPolicy `json:"-"`
	DBlock_Joins  []string             `json:"-"`
	Status        string               `json:"status,omitempty"`
	Error         string               `json:"error,omitempty"`
	Time          string               `json:"time,omitempty"`
	Resource      *Resource            `json:"_resource,omitempty"`
}
type Event_Wrap struct {
	Event []Event `json:"_event,omitempty"`
}

// #######################################
// RELEASE_ENTRY
// #######################################
type ReleaseEntry struct {
	DBlock_Table    string               `json:"release_entry,omitempty"`
	DBlock_Policy   core.DataBlockPolicy `json:"-"`
	DBlock_Joins    []string             `json:"-"`
	Name            string               `json:"name,omitempty"`
	Result          bool                 `json:"result,omitempty"`
	Reason          string               `json:"reason,omitempty"`
	Release         *Release             `json:"release,omitempty"`
	ReleaseCriteria *ReleaseCriteria     `json:"release_criteria,omitempty"`
	Resource        *Resource            `json:"_resource,omitempty"`
}
type ReleaseEntry_Wrap struct {
	ReleaseEntry []ReleaseEntry `json:"release_entry,omitempty"`
}

// #######################################
// RELEASE
// #######################################
type Release struct {
	DBlock_Table    string               `json:"release,omitempty"`
	DBlock_Policy   core.DataBlockPolicy `json:"-"`
	DBlock_Joins    []string             `json:"-"`
	Name            string               `json:"name,omitempty"`
	Version         string               `json:"version,omitempty"`
	Project         *Project             `json:"project,omitempty"`
	ReleaseItem     []ReleaseItem        `json:"release_item,omitempty"`
	ReleaseEntry    []ReleaseEntry       `json:"release_entry,omitempty"`
	ReleaseStage    []ReleaseStage       `json:"release_stage,omitempty"`
	ReleaseCriteria []ReleaseCriteria    `json:"release_criteria,omitempty"`
	CodeScan        []CodeScan           `json:"code_scan,omitempty"`
	TestRun         []TestRun            `json:"test_run,omitempty"`
}
type Release_Wrap struct {
	Release []Release `json:"release,omitempty"`
}

// #######################################
// PROJECT
// #######################################
type Project struct {
	DBlock_Table  string               `json:"project,omitempty"`
	DBlock_Policy core.DataBlockPolicy `json:"-"`
	DBlock_Joins  []string             `json:"-"`
	Name          string               `json:"name,omitempty"`
	Repo          []Repo               `json:"repo,omitempty"`
	Release       []Release            `json:"release,omitempty"`
}
type Project_Wrap struct {
	Project []Project `json:"project,omitempty"`
}

// #######################################
// REPO
// #######################################
type Repo struct {
	DBlock_Table  string               `json:"repo,omitempty"`
	DBlock_Policy core.DataBlockPolicy `json:"-"`
	DBlock_Joins  []string             `json:"-"`
	Id            string               `json:"id,omitempty"`
	Name          string               `json:"name,omitempty"`
	Project       *Project             `json:"project,omitempty"`
	Branch        []Branch             `json:"branch,omitempty"`
	Commit        []Commit             `json:"commit,omitempty"`
}
type Repo_Wrap struct {
	Repo []Repo `json:"repo,omitempty"`
}

// #######################################
// BRANCH
// #######################################
type Branch struct {
	DBlock_Table  string               `json:"branch,omitempty"`
	DBlock_Policy core.DataBlockPolicy `json:"-"`
	DBlock_Joins  []string             `json:"-"`
	Name          string               `json:"name,omitempty"`
	Repo          *Repo                `json:"repo,omitempty"`
	Commit        []Commit             `json:"commit,omitempty"`
}
type Branch_Wrap struct {
	Branch []Branch `json:"branch,omitempty"`
}

// #######################################
// COMMIT
// #######################################
type Commit struct {
	DBlock_Table  string               `json:"commit,omitempty"`
	DBlock_Policy core.DataBlockPolicy `json:"-"`
	DBlock_Joins  []string             `json:"-"`
	Id            string               `json:"id,omitempty"`
	Tag           string               `json:"tag,omitempty"`
	Time          string               `json:"time,omitempty"`
	Branch        *Branch              `json:"branch,omitempty"`
	Repo          *Repo                `json:"repo,omitempty"`
	ReleaseItem   *ReleaseItem         `json:"release_item,omitempty"`
}
type Commit_Wrap struct {
	Commit []Commit `json:"commit,omitempty"`
}

// #######################################
// RELEASE_ITEM
// #######################################
type ReleaseItem struct {
	DBlock_Table  string               `json:"release_item,omitempty"`
	DBlock_Policy core.DataBlockPolicy `json:"-"`
	DBlock_Joins  []string             `json:"-"`
	Type          string               `json:"type,omitempty"`
	Release       *Release             `json:"release,omitempty"`
	Commit        *Commit              `json:"commit,omitempty"`
	Artifact      *Artifact            `json:"artifact,omitempty"`
}
type ReleaseItem_Wrap struct {
	ReleaseItem []ReleaseItem `json:"release_item,omitempty"`
}

// #######################################
// ARTIFACT
// #######################################
type Artifact struct {
	DBlock_Table  string               `json:"artifact,omitempty"`
	DBlock_Policy core.DataBlockPolicy `json:"-"`
	DBlock_Joins  []string             `json:"-"`
	Name          string               `json:"name,omitempty"`
	Sha256        string               `json:"sha256,omitempty"`
	Location      string               `json:"location,omitempty"`
	ReleaseItem   *ReleaseItem         `json:"release_item,omitempty"`
}
type Artifact_Wrap struct {
	Artifact []Artifact `json:"artifact,omitempty"`
}

// #######################################
// RELEASE_STAGE
// #######################################
type ReleaseStage struct {
	DBlock_Table    string               `json:"release_stage,omitempty"`
	DBlock_Policy   core.DataBlockPolicy `json:"-"`
	DBlock_Joins    []string             `json:"-"`
	Name            string               `json:"name,omitempty"`
	Release         *Release             `json:"release,omitempty"`
	ReleaseCriteria []ReleaseCriteria    `json:"release_criteria,omitempty"`
}
type ReleaseStage_Wrap struct {
	ReleaseStage []ReleaseStage `json:"release_stage,omitempty"`
}

// #######################################
// RELEASE_CRITERIA
// #######################################
type ReleaseCriteria struct {
	DBlock_Table  string               `json:"release_criteria,omitempty"`
	DBlock_Policy core.DataBlockPolicy `json:"-"`
	DBlock_Joins  []string             `json:"-"`
	EntryName     string               `json:"entry_name,omitempty"`
	ReleaseEntry  []ReleaseEntry       `json:"release_entry,omitempty"`
	ReleaseStage  *ReleaseStage        `json:"release_stage,omitempty"`
	Release       *Release             `json:"release,omitempty"`
}
type ReleaseCriteria_Wrap struct {
	ReleaseCriteria []ReleaseCriteria `json:"release_criteria,omitempty"`
}

// #######################################
// CODE_SCAN
// #######################################
type CodeScan struct {
	DBlock_Table  string               `json:"code_scan,omitempty"`
	DBlock_Policy core.DataBlockPolicy `json:"-"`
	DBlock_Joins  []string             `json:"-"`
	Tool          string               `json:"tool,omitempty"`
	Release       *Release             `json:"release,omitempty"`
	CodeIssue     []CodeIssue          `json:"code_issue,omitempty"`
}
type CodeScan_Wrap struct {
	CodeScan []CodeScan `json:"code_scan,omitempty"`
}

// #######################################
// CODE_ISSUE
// #######################################
type CodeIssue struct {
	DBlock_Table  string               `json:"code_issue,omitempty"`
	DBlock_Policy core.DataBlockPolicy `json:"-"`
	DBlock_Joins  []string             `json:"-"`
	Id            string               `json:"id,omitempty"`
	Message       string               `json:"message,omitempty"`
	Severity      string               `json:"severity,omitempty"`
	Type          string               `json:"type,omitempty"`
	CodeScan      *CodeScan            `json:"code_scan,omitempty"`
}
type CodeIssue_Wrap struct {
	CodeIssue []CodeIssue `json:"code_issue,omitempty"`
}

// #######################################
// TEST_RUN
// #######################################
type TestRun struct {
	DBlock_Table  string               `json:"test_run,omitempty"`
	DBlock_Policy core.DataBlockPolicy `json:"-"`
	DBlock_Joins  []string             `json:"-"`
	Tool          string               `json:"tool,omitempty"`
	Type          string               `json:"type,omitempty"`
	Name          string               `json:"name,omitempty"`
	Elapsed       int64                `json:"elapsed,omitempty"`
	Result        bool                 `json:"result,omitempty"`
	Release       *Release             `json:"release,omitempty"`
	TestCase      []TestCase           `json:"test_case,omitempty"`
}
type TestRun_Wrap struct {
	TestRun []TestRun `json:"test_run,omitempty"`
}

// #######################################
// TEST_CASE
// #######################################
type TestCase struct {
	DBlock_Table  string               `json:"test_case,omitempty"`
	DBlock_Policy core.DataBlockPolicy `json:"-"`
	DBlock_Joins  []string             `json:"-"`
	Name          string               `json:"name,omitempty"`
	Result        bool                 `json:"result,omitempty"`
	Message       string               `json:"message,omitempty"`
	TestRun       *TestRun             `json:"test_run,omitempty"`
}
type TestCase_Wrap struct {
	TestCase []TestCase `json:"test_case,omitempty"`
}
