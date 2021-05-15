package builtin

// #######################################
// _SCHEMA
// #######################################
type Schema struct {
	Tables map[string]interface{} `json:"tables"`
}
type Schema_Wrap struct {
	Schema []Schema `json:"_schema"`
}

// #######################################
// _RESOURCE
// #######################################
type Resource struct {
	Id           string                 `json:"id"`
	Name         string                 `json:"name"`
	Kind         string                 `json:"kind"`
	ApiVersion   string                 `json:"api_version"`
	Spec         string                 `json:"spec"`
	Metadata     map[string]interface{} `json:"metadata"`
	Event        []Event                `json:"_event"`
	ReleaseEntry []ReleaseEntry         `json:"release_entry"`
}
type Resource_Wrap struct {
	Resource []Resource `json:"_resource"`
}

// #######################################
// _EVENT
// #######################################
type Event struct {
	Status   string    `json:"status"`
	Error    string    `json:"error"`
	Time     string    `json:"time"`
	Resource *Resource `json:"_resource"`
}
type Event_Wrap struct {
	Event []Event `json:"_event"`
}

// #######################################
// RELEASE_ENTRY
// #######################################
type ReleaseEntry struct {
	Name            string           `json:"name"`
	Result          bool             `json:"result"`
	Reason          string           `json:"reason"`
	Release         *Release         `json:"release"`
	ReleaseCriteria *ReleaseCriteria `json:"release_criteria"`
	Resource        *Resource        `json:"_resource"`
}
type ReleaseEntry_Wrap struct {
	ReleaseEntry []ReleaseEntry `json:"release_entry"`
}

// #######################################
// RELEASE
// #######################################
type Release struct {
	Name            string            `json:"name"`
	Version         string            `json:"version"`
	Project         *Project          `json:"project"`
	ReleaseItem     []ReleaseItem     `json:"release_item"`
	ReleaseEntry    []ReleaseEntry    `json:"release_entry"`
	ReleaseStage    []ReleaseStage    `json:"release_stage"`
	ReleaseCriteria []ReleaseCriteria `json:"release_criteria"`
	CodeScan        []CodeScan        `json:"code_scan"`
	TestRun         []TestRun         `json:"test_run"`
}
type Release_Wrap struct {
	Release []Release `json:"release"`
}

// #######################################
// PROJECT
// #######################################
type Project struct {
	Id      string    `json:"id"`
	Name    string    `json:"name"`
	Repo    []Repo    `json:"repo"`
	Release []Release `json:"release"`
}
type Project_Wrap struct {
	Project []Project `json:"project"`
}

// #######################################
// REPO
// #######################################
type Repo struct {
	Id      string   `json:"id"`
	Name    string   `json:"name"`
	Project *Project `json:"project"`
	Branch  []Branch `json:"branch"`
	Commit  []Commit `json:"commit"`
}
type Repo_Wrap struct {
	Repo []Repo `json:"repo"`
}

// #######################################
// BRANCH
// #######################################
type Branch struct {
	Name   string   `json:"name"`
	Repo   *Repo    `json:"repo"`
	Commit []Commit `json:"commit"`
}
type Branch_Wrap struct {
	Branch []Branch `json:"branch"`
}

// #######################################
// COMMIT
// #######################################
type Commit struct {
	Id          string       `json:"id"`
	Tag         string       `json:"tag"`
	Time        string       `json:"time"`
	Branch      *Branch      `json:"branch"`
	Repo        *Repo        `json:"repo"`
	ReleaseItem *ReleaseItem `json:"release_item"`
}
type Commit_Wrap struct {
	Commit []Commit `json:"commit"`
}

// #######################################
// RELEASE_ITEM
// #######################################
type ReleaseItem struct {
	Type     string    `json:"type"`
	Release  *Release  `json:"release"`
	Commit   *Commit   `json:"commit"`
	Artifact *Artifact `json:"artifact"`
}
type ReleaseItem_Wrap struct {
	ReleaseItem []ReleaseItem `json:"release_item"`
}

// #######################################
// ARTIFACT
// #######################################
type Artifact struct {
	Name        string       `json:"name"`
	Sha256      string       `json:"sha256"`
	Location    string       `json:"location"`
	ReleaseItem *ReleaseItem `json:"release_item"`
}
type Artifact_Wrap struct {
	Artifact []Artifact `json:"artifact"`
}

// #######################################
// RELEASE_STAGE
// #######################################
type ReleaseStage struct {
	Name            string            `json:"name"`
	Release         *Release          `json:"release"`
	ReleaseCriteria []ReleaseCriteria `json:"release_criteria"`
}
type ReleaseStage_Wrap struct {
	ReleaseStage []ReleaseStage `json:"release_stage"`
}

// #######################################
// RELEASE_CRITERIA
// #######################################
type ReleaseCriteria struct {
	EntryName    string         `json:"entry_name"`
	ReleaseEntry []ReleaseEntry `json:"release_entry"`
	ReleaseStage *ReleaseStage  `json:"release_stage"`
	Release      *Release       `json:"release"`
}
type ReleaseCriteria_Wrap struct {
	ReleaseCriteria []ReleaseCriteria `json:"release_criteria"`
}

// #######################################
// CODE_SCAN
// #######################################
type CodeScan struct {
	Tool      string      `json:"tool"`
	Release   *Release    `json:"release"`
	CodeIssue []CodeIssue `json:"code_issue"`
}
type CodeScan_Wrap struct {
	CodeScan []CodeScan `json:"code_scan"`
}

// #######################################
// CODE_ISSUE
// #######################################
type CodeIssue struct {
	Id       string    `json:"id"`
	Message  string    `json:"message"`
	Severity string    `json:"severity"`
	Type     string    `json:"type"`
	CodeScan *CodeScan `json:"code_scan"`
}
type CodeIssue_Wrap struct {
	CodeIssue []CodeIssue `json:"code_issue"`
}

// #######################################
// TEST_RUN
// #######################################
type TestRun struct {
	Tool     string     `json:"tool"`
	Type     string     `json:"type"`
	Name     string     `json:"name"`
	Elapsed  int64      `json:"elapsed"`
	Result   bool       `json:"result"`
	Release  *Release   `json:"release"`
	TestCase []TestCase `json:"test_case"`
}
type TestRun_Wrap struct {
	TestRun []TestRun `json:"test_run"`
}

// #######################################
// TEST_CASE
// #######################################
type TestCase struct {
	Name    string   `json:"name"`
	Result  bool     `json:"result"`
	Message string   `json:"message"`
	TestRun *TestRun `json:"test_run"`
}
type TestCase_Wrap struct {
	TestCase []TestCase `json:"test_case"`
}
