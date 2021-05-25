// #######################################
// _SCHEMA
// #######################################
export interface _schema {
	_id?: string;
	tables?: object;
}
export interface _schema_wrap {
	_schema?: _schema[];
}

// #######################################
// _RESOURCE
// #######################################
export interface _resource {
	_id?: string;
	id?: string;
	name?: string;
	kind?: string;
	api_version?: string;
	spec?: string;
	metadata?: object;
	_event?: _event[];
	release_entry?: release_entry[];
}
export interface _resource_wrap {
	_resource?: _resource[];
}

// #######################################
// _EVENT
// #######################################
export interface _event {
	_id?: string;
	status?: string;
	error?: string;
	time?: string;
	_resource?: _resource;
}
export interface _event_wrap {
	_event?: _event[];
}

// #######################################
// RELEASE_ENTRY
// #######################################
export interface release_entry {
	_id?: string;
	name?: string;
	result?: boolean;
	reason?: string;
	release?: release;
	release_criteria?: release_criteria;
	_resource?: _resource;
}
export interface release_entry_wrap {
	release_entry?: release_entry[];
}

// #######################################
// RELEASE
// #######################################
export interface release {
	_id?: string;
	name?: string;
	version?: string;
	project?: project;
	release_input?: release_input[];
	release_entry?: release_entry[];
	release_stage?: release_stage[];
	release_criteria?: release_criteria[];
	code_scan?: code_scan[];
	test_run?: test_run[];
}
export interface release_wrap {
	release?: release[];
}

// #######################################
// PROJECT
// #######################################
export interface project {
	_id?: string;
	id?: string;
	name?: string;
	repo?: repo[];
	release?: release[];
}
export interface project_wrap {
	project?: project[];
}

// #######################################
// REPO
// #######################################
export interface repo {
	_id?: string;
	id?: string;
	name?: string;
	project?: project;
	branch?: branch[];
	commit?: commit[];
}
export interface repo_wrap {
	repo?: repo[];
}

// #######################################
// BRANCH
// #######################################
export interface branch {
	_id?: string;
	name?: string;
	repo?: repo;
	commit?: commit[];
}
export interface branch_wrap {
	branch?: branch[];
}

// #######################################
// COMMIT
// #######################################
export interface commit {
	_id?: string;
	id?: string;
	tag?: string;
	time?: string;
	branch?: branch;
	repo?: repo;
	release_input?: release_input;
}
export interface commit_wrap {
	commit?: commit[];
}

// #######################################
// RELEASE_INPUT
// #######################################
export interface release_input {
	_id?: string;
	type?: string;
	release?: release;
	commit?: commit;
}
export interface release_input_wrap {
	release_input?: release_input[];
}

// #######################################
// RELEASE_STAGE
// #######################################
export interface release_stage {
	_id?: string;
	name?: string;
	release?: release;
	release_criteria?: release_criteria[];
}
export interface release_stage_wrap {
	release_stage?: release_stage[];
}

// #######################################
// RELEASE_CRITERIA
// #######################################
export interface release_criteria {
	_id?: string;
	entry_name?: string;
	release_entry?: release_entry[];
	release_stage?: release_stage;
	release?: release;
}
export interface release_criteria_wrap {
	release_criteria?: release_criteria[];
}

// #######################################
// CODE_SCAN
// #######################################
export interface code_scan {
	_id?: string;
	tool?: string;
	release?: release;
	code_issue?: code_issue[];
}
export interface code_scan_wrap {
	code_scan?: code_scan[];
}

// #######################################
// CODE_ISSUE
// #######################################
export interface code_issue {
	_id?: string;
	id?: string;
	message?: string;
	severity?: string;
	type?: string;
	code_scan?: code_scan;
}
export interface code_issue_wrap {
	code_issue?: code_issue[];
}

// #######################################
// TEST_RUN
// #######################################
export interface test_run {
	_id?: string;
	tool?: string;
	type?: string;
	name?: string;
	elapsed?: number;
	result?: boolean;
	release?: release;
	test_case?: test_case[];
}
export interface test_run_wrap {
	test_run?: test_run[];
}

// #######################################
// TEST_CASE
// #######################################
export interface test_case {
	_id?: string;
	name?: string;
	result?: boolean;
	message?: string;
	test_run?: test_run;
}
export interface test_case_wrap {
	test_case?: test_case[];
}

// #######################################
// ARTIFACT
// #######################################
export interface artifact {
	_id?: string;
	name?: string;
	sha256?: string;
	location?: string;
}
export interface artifact_wrap {
	artifact?: artifact[];
}

