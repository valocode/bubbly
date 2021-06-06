// #######################################
// _SCHEMA
// #######################################
export interface _schema {
	tables?: object;
}
export interface _schema_wrap {
	_schema?: _schema[];
}

// #######################################
// _RESOURCE
// #######################################
export interface _resource {
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
	status?: string;
	error?: string;
	time?: Date;
	_resource?: _resource;
}
export interface _event_wrap {
	_event?: _event[];
}

// #######################################
// RELEASE_ENTRY
// #######################################
export interface release_entry {
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
	name?: string;
	version?: string;
	project?: project;
	release_item?: release_item[];
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
// RELEASE_CRITERIA
// #######################################
export interface release_criteria {
	entry_name?: string;
	release_stage?: release_stage;
	release?: release;
	release_entry?: release_entry[];
}
export interface release_criteria_wrap {
	release_criteria?: release_criteria[];
}

// #######################################
// RELEASE_STAGE
// #######################################
export interface release_stage {
	name?: string;
	release?: release;
	release_criteria?: release_criteria[];
}
export interface release_stage_wrap {
	release_stage?: release_stage[];
}

// #######################################
// CODE_SCAN
// #######################################
export interface code_scan {
	tool?: string;
	release?: release;
	lifecycle?: lifecycle;
	code_issue?: code_issue[];
}
export interface code_scan_wrap {
	code_scan?: code_scan[];
}

// #######################################
// LIFECYCLE
// #######################################
export interface lifecycle {
	status?: string;
	lifecycle_entry?: lifecycle_entry[];
	code_scan?: code_scan;
	test_run?: test_run;
}
export interface lifecycle_wrap {
	lifecycle?: lifecycle[];
}

// #######################################
// LIFECYCLE_ENTRY
// #######################################
export interface lifecycle_entry {
	author?: string;
	message?: string;
	signed?: boolean;
	lifecycle?: lifecycle;
}
export interface lifecycle_entry_wrap {
	lifecycle_entry?: lifecycle_entry[];
}

// #######################################
// TEST_RUN
// #######################################
export interface test_run {
	tool?: string;
	type?: string;
	name?: string;
	elapsed?: number;
	result?: boolean;
	test_case?: test_case[];
	release?: release;
	lifecycle?: lifecycle;
}
export interface test_run_wrap {
	test_run?: test_run[];
}

// #######################################
// TEST_CASE
// #######################################
export interface test_case {
	name?: string;
	result?: boolean;
	message?: string;
	test_run?: test_run;
}
export interface test_case_wrap {
	test_case?: test_case[];
}

// #######################################
// CODE_ISSUE
// #######################################
export interface code_issue {
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
// PROJECT
// #######################################
export interface project {
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
	id?: string;
	name?: string;
	branch?: branch[];
	commit?: commit[];
	project?: project;
}
export interface repo_wrap {
	repo?: repo[];
}

// #######################################
// BRANCH
// #######################################
export interface branch {
	name?: string;
	commit?: commit[];
	repo?: repo;
}
export interface branch_wrap {
	branch?: branch[];
}

// #######################################
// COMMIT
// #######################################
export interface commit {
	id?: string;
	tag?: string;
	time?: string;
	branch?: branch;
	repo?: repo;
	release_item?: release_item;
}
export interface commit_wrap {
	commit?: commit[];
}

// #######################################
// RELEASE_ITEM
// #######################################
export interface release_item {
	type?: string;
	commit?: commit;
	artifact?: artifact;
	release?: release;
}
export interface release_item_wrap {
	release_item?: release_item[];
}

// #######################################
// ARTIFACT
// #######################################
export interface artifact {
	name?: string;
	sha256?: string;
	location?: string;
	release_item?: release_item;
}
export interface artifact_wrap {
	artifact?: artifact[];
}

