// Code generated by entc, DO NOT EDIT.

package predicate

import (
	"entgo.io/ent/dialect/sql"
)

// Artifact is the predicate function for artifact builders.
type Artifact func(*sql.Selector)

// CVE is the predicate function for cve builders.
type CVE func(*sql.Selector)

// CVERule is the predicate function for cverule builders.
type CVERule func(*sql.Selector)

// CWE is the predicate function for cwe builders.
type CWE func(*sql.Selector)

// CodeIssue is the predicate function for codeissue builders.
type CodeIssue func(*sql.Selector)

// CodeScan is the predicate function for codescan builders.
type CodeScan func(*sql.Selector)

// Component is the predicate function for component builders.
type Component func(*sql.Selector)

// ComponentUse is the predicate function for componentuse builders.
type ComponentUse func(*sql.Selector)

// GitCommit is the predicate function for gitcommit builders.
type GitCommit func(*sql.Selector)

// License is the predicate function for license builders.
type License func(*sql.Selector)

// LicenseUse is the predicate function for licenseuse builders.
type LicenseUse func(*sql.Selector)

// Project is the predicate function for project builders.
type Project func(*sql.Selector)

// Release is the predicate function for release builders.
type Release func(*sql.Selector)

// ReleaseEntry is the predicate function for releaseentry builders.
type ReleaseEntry func(*sql.Selector)

// Repo is the predicate function for repo builders.
type Repo func(*sql.Selector)

// TestCase is the predicate function for testcase builders.
type TestCase func(*sql.Selector)

// TestRun is the predicate function for testrun builders.
type TestRun func(*sql.Selector)

// Vulnerability is the predicate function for vulnerability builders.
type Vulnerability func(*sql.Selector)
