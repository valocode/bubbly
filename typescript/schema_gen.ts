// #######################################
// Code is generated using a custom ent extension.
// DO NOT MODIFY.
// Currently it is manually copied over from the bubbly repository where it is generated.
// #######################################

// #######################################
// Artifact
// #######################################
export interface Artifact_Json {
	artifact?: Artifact[];
}

export interface Artifact {
	id?: number;
	name?: string;
	sha256?: string;
	type?: ArtifactType;
	release?: Release;
	entry?: ReleaseEntry;
}

export interface Artifact_Relay {
	artifact_connection?: Artifact_Conn;
}

export interface Artifact_Conn {
	totalCount?: number;
	pageInfo?: pageInfo;
	edges?: Artifact_Edge[];
}

export interface Artifact_Edge {
	node?: Artifact;
}

export enum ArtifactType {
	docker = "docker",
	file = "file",
}

// #######################################
// CVE
// #######################################
export interface CVE_Json {
	cve?: CVE[];
}

export interface CVE {
	id?: number;
	cve_id?: string;
	description?: string;
	severity_score?: number;
	severity?: CVESeverity;
	published_data?: Date;
	modified_data?: Date;
	found?: Vulnerability[];
	rules?: CVERule[];
}

export interface CVE_Relay {
	cve_connection?: CVE_Conn;
}

export interface CVE_Conn {
	totalCount?: number;
	pageInfo?: pageInfo;
	edges?: CVE_Edge[];
}

export interface CVE_Edge {
	node?: CVE;
}

export enum CVESeverity {
	None = "None",
	Low = "Low",
	Medium = "Medium",
	High = "High",
	Critical = "Critical",
}

// #######################################
// CVERule
// #######################################
export interface CVERule_Json {
	cve_rule?: CVERule[];
}

export interface CVERule {
	id?: number;
	name?: string;
	cve?: CVE;
	project?: Project[];
	repo?: Repo[];
}

export interface CVERule_Relay {
	cve_rule_connection?: CVERule_Conn;
}

export interface CVERule_Conn {
	totalCount?: number;
	pageInfo?: pageInfo;
	edges?: CVERule_Edge[];
}

export interface CVERule_Edge {
	node?: CVERule;
}

// #######################################
// CVEScan
// #######################################
export interface CVEScan_Json {
	cve_scan?: CVEScan[];
}

export interface CVEScan {
	id?: number;
	tool?: string;
	release?: Release;
	entry?: ReleaseEntry;
	vulnerabilities?: Vulnerability[];
}

export interface CVEScan_Relay {
	cve_scan_connection?: CVEScan_Conn;
}

export interface CVEScan_Conn {
	totalCount?: number;
	pageInfo?: pageInfo;
	edges?: CVEScan_Edge[];
}

export interface CVEScan_Edge {
	node?: CVEScan;
}

// #######################################
// CWE
// #######################################
export interface CWE_Json {
	cwe?: CWE[];
}

export interface CWE {
	id?: number;
	cwe_id?: string;
	description?: string;
	url?: number;
	issues?: CodeIssue[];
}

export interface CWE_Relay {
	cwe_connection?: CWE_Conn;
}

export interface CWE_Conn {
	totalCount?: number;
	pageInfo?: pageInfo;
	edges?: CWE_Edge[];
}

export interface CWE_Edge {
	node?: CWE;
}

// #######################################
// CodeIssue
// #######################################
export interface CodeIssue_Json {
	code_issue?: CodeIssue[];
}

export interface CodeIssue {
	id?: number;
	rule_id?: string;
	message?: string;
	severity?: CodeIssueSeverity;
	type?: CodeIssueType;
	cwe?: CWE[];
	scan?: CodeScan;
}

export interface CodeIssue_Relay {
	code_issue_connection?: CodeIssue_Conn;
}

export interface CodeIssue_Conn {
	totalCount?: number;
	pageInfo?: pageInfo;
	edges?: CodeIssue_Edge[];
}

export interface CodeIssue_Edge {
	node?: CodeIssue;
}

export enum CodeIssueSeverity {
	low = "low",
	medium = "medium",
	high = "high",
}

export enum CodeIssueType {
	style = "style",
	security = "security",
	bug = "bug",
}

// #######################################
// CodeScan
// #######################################
export interface CodeScan_Json {
	code_scan?: CodeScan[];
}

export interface CodeScan {
	id?: number;
	tool?: string;
	release?: Release;
	issues?: CodeIssue[];
	entry?: ReleaseEntry;
}

export interface CodeScan_Relay {
	code_scan_connection?: CodeScan_Conn;
}

export interface CodeScan_Conn {
	totalCount?: number;
	pageInfo?: pageInfo;
	edges?: CodeScan_Edge[];
}

export interface CodeScan_Edge {
	node?: CodeScan;
}

// #######################################
// Component
// #######################################
export interface Component_Json {
	component?: Component[];
}

export interface Component {
	id?: number;
	name?: string;
	vendor?: string;
	version?: string;
	description?: string;
	url?: string;
	vulnerabilities?: Vulnerability[];
	licenses?: License[];
	release?: Release[];
}

export interface Component_Relay {
	component_connection?: Component_Conn;
}

export interface Component_Conn {
	totalCount?: number;
	pageInfo?: pageInfo;
	edges?: Component_Edge[];
}

export interface Component_Edge {
	node?: Component;
}

// #######################################
// GitCommit
// #######################################
export interface GitCommit_Json {
	commit?: GitCommit[];
}

export interface GitCommit {
	id?: number;
	hash?: string;
	branch?: string;
	tag?: string;
	time?: Date;
	repo?: Repo;
	release?: Release;
}

export interface GitCommit_Relay {
	commit_connection?: GitCommit_Conn;
}

export interface GitCommit_Conn {
	totalCount?: number;
	pageInfo?: pageInfo;
	edges?: GitCommit_Edge[];
}

export interface GitCommit_Edge {
	node?: GitCommit;
}

// #######################################
// License
// #######################################
export interface License_Json {
	license?: License[];
}

export interface License {
	id?: number;
	spdx_id?: string;
	name?: string;
	reference?: string;
	details_url?: string;
	is_osi_approved?: boolean;
	components?: Component[];
	usages?: LicenseUsage[];
}

export interface License_Relay {
	license_connection?: License_Conn;
}

export interface License_Conn {
	totalCount?: number;
	pageInfo?: pageInfo;
	edges?: License_Edge[];
}

export interface License_Edge {
	node?: License;
}

// #######################################
// LicenseScan
// #######################################
export interface LicenseScan_Json {
	license_scan?: LicenseScan[];
}

export interface LicenseScan {
	id?: number;
	tool?: string;
	release?: Release;
	entry?: ReleaseEntry;
	licenses?: LicenseUsage[];
}

export interface LicenseScan_Relay {
	license_scan_connection?: LicenseScan_Conn;
}

export interface LicenseScan_Conn {
	totalCount?: number;
	pageInfo?: pageInfo;
	edges?: LicenseScan_Edge[];
}

export interface LicenseScan_Edge {
	node?: LicenseScan;
}

// #######################################
// LicenseUsage
// #######################################
export interface LicenseUsage_Json {
	license_usage?: LicenseUsage[];
}

export interface LicenseUsage {
	id?: number;
	license?: License;
	scan?: LicenseScan;
}

export interface LicenseUsage_Relay {
	license_usage_connection?: LicenseUsage_Conn;
}

export interface LicenseUsage_Conn {
	totalCount?: number;
	pageInfo?: pageInfo;
	edges?: LicenseUsage_Edge[];
}

export interface LicenseUsage_Edge {
	node?: LicenseUsage;
}

// #######################################
// Project
// #######################################
export interface Project_Json {
	project?: Project[];
}

export interface Project {
	id?: number;
	name?: string;
	repos?: Repo[];
	releases?: Release[];
	cve_rules?: CVERule[];
}

export interface Project_Relay {
	project_connection?: Project_Conn;
}

export interface Project_Conn {
	totalCount?: number;
	pageInfo?: pageInfo;
	edges?: Project_Edge[];
}

export interface Project_Edge {
	node?: Project;
}

// #######################################
// Release
// #######################################
export interface Release_Json {
	release?: Release[];
}

export interface Release {
	id?: number;
	name?: string;
	version?: string;
	status?: ReleaseStatus;
	subreleases?: Release[];
	dependencies?: Release[];
	project?: Project;
	commit?: GitCommit;
	artifacts?: Artifact[];
	checks?: ReleaseCheck[];
	log?: ReleaseEntry[];
	code_scans?: CodeScan[];
	cve_scans?: CVEScan[];
	license_scans?: LicenseScan[];
	test_runs?: TestRun[];
	components?: Component[];
}

export interface Release_Relay {
	release_connection?: Release_Conn;
}

export interface Release_Conn {
	totalCount?: number;
	pageInfo?: pageInfo;
	edges?: Release_Edge[];
}

export interface Release_Edge {
	node?: Release;
}

export enum ReleaseStatus {
	pending = "pending",
	ready = "ready",
	blocked = "blocked",
}

// #######################################
// ReleaseCheck
// #######################################
export interface ReleaseCheck_Json {
	release_check?: ReleaseCheck[];
}

export interface ReleaseCheck {
	id?: number;
	type?: ReleaseCheckType;
	release?: Release;
}

export interface ReleaseCheck_Relay {
	release_check_connection?: ReleaseCheck_Conn;
}

export interface ReleaseCheck_Conn {
	totalCount?: number;
	pageInfo?: pageInfo;
	edges?: ReleaseCheck_Edge[];
}

export interface ReleaseCheck_Edge {
	node?: ReleaseCheck;
}

export enum ReleaseCheckType {
	artifact = "artifact",
	unit_test = "unit_test",
	security = "security",
}

// #######################################
// ReleaseEntry
// #######################################
export interface ReleaseEntry_Json {
	release_entry?: ReleaseEntry[];
}

export interface ReleaseEntry {
	id?: number;
	type?: ReleaseEntryType;
	time?: Date;
	artifact?: Artifact;
	code_scan?: CodeScan;
	test_run?: TestRun;
	cve_scan?: CVEScan;
	license_scan?: LicenseScan;
	release?: Release;
}

export interface ReleaseEntry_Relay {
	release_entry_connection?: ReleaseEntry_Conn;
}

export interface ReleaseEntry_Conn {
	totalCount?: number;
	pageInfo?: pageInfo;
	edges?: ReleaseEntry_Edge[];
}

export interface ReleaseEntry_Edge {
	node?: ReleaseEntry;
}

export enum ReleaseEntryType {
	artifact = "artifact",
	deploy = "deploy",
	code_scan = "code_scan",
	license_scan = "license_scan",
	cve_scan = "cve_scan",
	test_run = "test_run",
}

// #######################################
// Repo
// #######################################
export interface Repo_Json {
	repo?: Repo[];
}

export interface Repo {
	id?: number;
	name?: string;
	project?: Project;
	commits?: GitCommit[];
	cve_rules?: CVERule[];
}

export interface Repo_Relay {
	repo_connection?: Repo_Conn;
}

export interface Repo_Conn {
	totalCount?: number;
	pageInfo?: pageInfo;
	edges?: Repo_Edge[];
}

export interface Repo_Edge {
	node?: Repo;
}

// #######################################
// TestCase
// #######################################
export interface TestCase_Json {
	test_case?: TestCase[];
}

export interface TestCase {
	id?: number;
	name?: string;
	result?: boolean;
	message?: string;
	elapsed?: number;
	run?: TestRun;
}

export interface TestCase_Relay {
	test_case_connection?: TestCase_Conn;
}

export interface TestCase_Conn {
	totalCount?: number;
	pageInfo?: pageInfo;
	edges?: TestCase_Edge[];
}

export interface TestCase_Edge {
	node?: TestCase;
}

// #######################################
// TestRun
// #######################################
export interface TestRun_Json {
	test_run?: TestRun[];
}

export interface TestRun {
	id?: number;
	tool?: string;
	release?: Release;
	entry?: ReleaseEntry;
	tests?: TestCase[];
}

export interface TestRun_Relay {
	test_run_connection?: TestRun_Conn;
}

export interface TestRun_Conn {
	totalCount?: number;
	pageInfo?: pageInfo;
	edges?: TestRun_Edge[];
}

export interface TestRun_Edge {
	node?: TestRun;
}

// #######################################
// Vulnerability
// #######################################
export interface Vulnerability_Json {
	vulnerability?: Vulnerability[];
}

export interface Vulnerability {
	id?: number;
	cve?: CVE;
	scan?: CVEScan;
	component?: Component;
}

export interface Vulnerability_Relay {
	vulnerability_connection?: Vulnerability_Conn;
}

export interface Vulnerability_Conn {
	totalCount?: number;
	pageInfo?: pageInfo;
	edges?: Vulnerability_Edge[];
}

export interface Vulnerability_Edge {
	node?: Vulnerability;
}

export interface pageInfo {
	hasNextPage?: boolean;
	hasPreviousPage?: boolean;
	startCursor?: string;
	endCursor?: string;
}

