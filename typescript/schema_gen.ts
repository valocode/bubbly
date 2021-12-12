// #######################################
// Code is generated using a custom ent extension.
// DO NOT MODIFY.
// Currently it is manually copied over from the bubbly repository where it is generated.
// #######################################

// #######################################
// Adapter
// #######################################
// #######################################
// Artifact
// #######################################
export enum ArtifactType {
	docker = "docker",
	file = "file",
}

// #######################################
// CodeIssue
// #######################################
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
// #######################################
// Component
// #######################################
// #######################################
// Event
// #######################################
export enum EventStatus {
	ok = "ok",
	error = "error",
}

export enum EventType {
	evaluate_release = "evaluate_release",
	monitor = "monitor",
}

// #######################################
// GitCommit
// #######################################
// #######################################
// License
// #######################################
// #######################################
// Organization
// #######################################
// #######################################
// Project
// #######################################
// #######################################
// Release
// #######################################
// #######################################
// ReleaseComponent
// #######################################
export enum ReleaseComponentType {
	embedded = "embedded",
	distributed = "distributed",
	development = "development",
}

// #######################################
// ReleaseEntry
// #######################################
export enum ReleaseEntryType {
	artifact = "artifact",
	deploy = "deploy",
	code_scan = "code_scan",
	test_run = "test_run",
}

// #######################################
// ReleaseLicense
// #######################################
// #######################################
// ReleasePolicy
// #######################################
// #######################################
// ReleasePolicyViolation
// #######################################
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
// #######################################
// Repository
// #######################################
// #######################################
// SPDXLicense
// #######################################
// #######################################
// TestCase
// #######################################
// #######################################
// TestRun
// #######################################
// #######################################
// Vulnerability
// #######################################
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
export enum VulnerabilityReviewDecision {
	exploitable = "exploitable",
	in_progress = "in_progress",
	invalid = "invalid",
	mitigated = "mitigated",
	ineffective = "ineffective",
	patched = "patched",
}

