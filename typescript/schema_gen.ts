// #######################################
// Code is generated using a custom ent extension.
// DO NOT MODIFY.
// Currently it is manually copied over from the bubbly repository where it is generated.
// #######################################

// #######################################
// Adapter
// #######################################
export interface Adapter_Json {
	adapter?: Adapter[];
}

export interface Adapter {
	id?: number;
	name?: string;
	tag?: string;
	module?: string;
	owner?: Organization;
}

export interface Adapter_Relay {
	adapter_connection?: Adapter_Conn;
}

export interface Adapter_Conn {
	totalCount?: number;
	pageInfo?: pageInfo;
	edges?: Adapter_Edge[];
}

export interface Adapter_Edge {
	node?: Adapter;
}

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
	time?: Date;
	metadata?: UNKNOWN_TYPE_schema.Metadata;
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
	metadata?: UNKNOWN_TYPE_schema.Metadata;
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
	time?: Date;
	metadata?: UNKNOWN_TYPE_schema.Metadata;
	release?: Release;
	entry?: ReleaseEntry;
	issues?: CodeIssue[];
	vulnerabilities?: ReleaseVulnerability[];
	components?: ReleaseComponent[];
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
	metadata?: UNKNOWN_TYPE_schema.Metadata;
	owner?: Organization;
	vulnerabilities?: Vulnerability[];
	licenses?: License[];
	uses?: ReleaseComponent[];
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
	uses?: LicenseUse[];
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
// LicenseUse
// #######################################
export interface LicenseUse_Json {
	license_use?: LicenseUse[];
}

export interface LicenseUse {
	id?: number;
	license?: License;
}

export interface LicenseUse_Relay {
	license_use_connection?: LicenseUse_Conn;
}

export interface LicenseUse_Conn {
	totalCount?: number;
	pageInfo?: pageInfo;
	edges?: LicenseUse_Edge[];
}

export interface LicenseUse_Edge {
	node?: LicenseUse;
}

// #######################################
// Organization
// #######################################
export interface Organization_Json {
	organization?: Organization[];
}

export interface Organization {
	id?: number;
	name?: string;
	projects?: Project[];
	repos?: Repo[];
}

export interface Organization_Relay {
	organization_connection?: Organization_Conn;
}

export interface Organization_Conn {
	totalCount?: number;
	pageInfo?: pageInfo;
	edges?: Organization_Edge[];
}

export interface Organization_Edge {
	node?: Organization;
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
	owner?: Organization;
	repos?: Repo[];
	vulnerability_reviews?: VulnerabilityReview[];
	policies?: ReleasePolicy[];
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
	commit?: GitCommit;
	head_of?: Repo;
	log?: ReleaseEntry[];
	violations?: ReleasePolicyViolation[];
	artifacts?: Artifact[];
	components?: ReleaseComponent[];
	vulnerabilities?: ReleaseVulnerability[];
	code_scans?: CodeScan[];
	test_runs?: TestRun[];
	vulnerability_reviews?: VulnerabilityReview[];
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
// ReleaseComponent
// #######################################
export interface ReleaseComponent_Json {
	release_component?: ReleaseComponent[];
}

export interface ReleaseComponent {
	id?: number;
	type?: ReleaseComponentType;
	release?: Release;
	scans?: CodeScan[];
	component?: Component;
	vulnerabilities?: ReleaseVulnerability[];
}

export interface ReleaseComponent_Relay {
	release_component_connection?: ReleaseComponent_Conn;
}

export interface ReleaseComponent_Conn {
	totalCount?: number;
	pageInfo?: pageInfo;
	edges?: ReleaseComponent_Edge[];
}

export interface ReleaseComponent_Edge {
	node?: ReleaseComponent;
}

export enum ReleaseComponentType {
	embedded = "embedded",
	distributed = "distributed",
	development = "development",
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
	test_run = "test_run",
}

// #######################################
// ReleaseLicense
// #######################################
export interface ReleaseLicense_Json {
	release_license?: ReleaseLicense[];
}

export interface ReleaseLicense {
	id?: number;
	license?: License;
	component?: ReleaseComponent;
	release?: Release;
	scans?: CodeScan[];
}

export interface ReleaseLicense_Relay {
	release_license_connection?: ReleaseLicense_Conn;
}

export interface ReleaseLicense_Conn {
	totalCount?: number;
	pageInfo?: pageInfo;
	edges?: ReleaseLicense_Edge[];
}

export interface ReleaseLicense_Edge {
	node?: ReleaseLicense;
}

// #######################################
// ReleasePolicy
// #######################################
export interface ReleasePolicy_Json {
	release_policy?: ReleasePolicy[];
}

export interface ReleasePolicy {
	id?: number;
	name?: string;
	module?: string;
	owner?: Organization;
	projects?: Project[];
	repos?: Repo[];
	violations?: ReleasePolicyViolation[];
}

export interface ReleasePolicy_Relay {
	release_policy_connection?: ReleasePolicy_Conn;
}

export interface ReleasePolicy_Conn {
	totalCount?: number;
	pageInfo?: pageInfo;
	edges?: ReleasePolicy_Edge[];
}

export interface ReleasePolicy_Edge {
	node?: ReleasePolicy;
}

// #######################################
// ReleasePolicyViolation
// #######################################
export interface ReleasePolicyViolation_Json {
	release_policy_violation?: ReleasePolicyViolation[];
}

export interface ReleasePolicyViolation {
	id?: number;
	message?: string;
	type?: ReleasePolicyViolationType;
	severity?: ReleasePolicyViolationSeverity;
	policy?: ReleasePolicy;
	release?: Release;
}

export interface ReleasePolicyViolation_Relay {
	release_policy_violation_connection?: ReleasePolicyViolation_Conn;
}

export interface ReleasePolicyViolation_Conn {
	totalCount?: number;
	pageInfo?: pageInfo;
	edges?: ReleasePolicyViolation_Edge[];
}

export interface ReleasePolicyViolation_Edge {
	node?: ReleasePolicyViolation;
}

export enum ReleasePolicyViolationType {
	require = "require",
	deny = "deny",
}

export enum ReleasePolicyViolationSeverity {
	suggestion = "suggestion",
	warning = "warning",
	blocking = "blocking",
}

// #######################################
// ReleaseVulnerability
// #######################################
export interface ReleaseVulnerability_Json {
	release_vulnerability?: ReleaseVulnerability[];
}

export interface ReleaseVulnerability {
	id?: number;
	vulnerability?: Vulnerability;
	component?: ReleaseComponent;
	release?: Release;
	reviews?: VulnerabilityReview[];
	scan?: CodeScan;
}

export interface ReleaseVulnerability_Relay {
	release_vulnerability_connection?: ReleaseVulnerability_Conn;
}

export interface ReleaseVulnerability_Conn {
	totalCount?: number;
	pageInfo?: pageInfo;
	edges?: ReleaseVulnerability_Edge[];
}

export interface ReleaseVulnerability_Edge {
	node?: ReleaseVulnerability;
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
	default_branch?: string;
	owner?: Organization;
	project?: Project;
	head?: Release;
	commits?: GitCommit[];
	vulnerability_reviews?: VulnerabilityReview[];
	policies?: ReleasePolicy[];
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
	metadata?: UNKNOWN_TYPE_schema.Metadata;
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
	time?: Date;
	metadata?: UNKNOWN_TYPE_schema.Metadata;
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
	vid?: string;
	summary?: string;
	description?: string;
	severity_score?: number;
	severity?: VulnerabilitySeverity;
	published?: Date;
	modified?: Date;
	metadata?: UNKNOWN_TYPE_schema.Metadata;
	owner?: Organization;
	components?: Component[];
	reviews?: VulnerabilityReview[];
	instances?: ReleaseVulnerability[];
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

export enum VulnerabilitySeverity {
	None = "None",
	Low = "Low",
	Medium = "Medium",
	High = "High",
	Critical = "Critical",
}

// #######################################
// VulnerabilityReview
// #######################################
export interface VulnerabilityReview_Json {
	vulnerability_review?: VulnerabilityReview[];
}

export interface VulnerabilityReview {
	id?: number;
	name?: string;
	decision?: VulnerabilityReviewDecision;
	vulnerability?: Vulnerability;
	projects?: Project[];
	repos?: Repo[];
	releases?: Release[];
	instances?: ReleaseVulnerability[];
}

export interface VulnerabilityReview_Relay {
	vulnerability_review_connection?: VulnerabilityReview_Conn;
}

export interface VulnerabilityReview_Conn {
	totalCount?: number;
	pageInfo?: pageInfo;
	edges?: VulnerabilityReview_Edge[];
}

export interface VulnerabilityReview_Edge {
	node?: VulnerabilityReview;
}

export enum VulnerabilityReviewDecision {
	exploitable = "exploitable",
	in_progress = "in_progress",
	invalid = "invalid",
	mitigated = "mitigated",
	ineffective = "ineffective",
}

export interface pageInfo {
	hasNextPage?: boolean;
	hasPreviousPage?: boolean;
	startCursor?: string;
	endCursor?: string;
}

