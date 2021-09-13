---
title: Schema
hide_title: false
hide_table_of_contents: false
keywords:
- docs
- bubbly
- schema
---

import Mermaid from '@theme/Mermaid';

## Introduction

The schema is an core part of bubbly - it is needed to store all the relevant Release Readiness data over which policies are run.

The library used to define the schema is [entgo](https://entgo.io/) which is an amazing framework, and this documentation was generated using it :)

## Overview

The following diagram shows the bubbly schema
<Mermaid chart={`
erDiagram
    Adapter {
        int id
        string name
        string tag
        string module
    }
    Artifact {
        int id
        string name
        string sha256
        artifactDOTType type
        timeDOTTime time
        schemaDOTMetadata metadata
    }
    CodeIssue {
        int id
        string rule_id
        string message
        codeissueDOTSeverity severity
        codeissueDOTType type
        schemaDOTMetadata metadata
    }
    CodeScan {
        int id
        string tool
        timeDOTTime time
        schemaDOTMetadata metadata
    }
    Component {
        int id
        string name
        string vendor
        string version
        string description
        string url
        schemaDOTMetadata metadata
    }
    GitCommit {
        int id
        string hash
        string branch
        string tag
        timeDOTTime time
    }
    License {
        int id
        string spdx_id
        string name
        string reference
        string details_url
        bool is_osi_approved
    }
    LicenseUse {
        int id
    }
    Organization {
        int id
        string name
    }
    Project {
        int id
        string name
    }
    Release {
        int id
        string name
        string version
        releaseDOTStatus status
    }
    ReleaseComponent {
        int id
        releasecomponentDOTType type
    }
    ReleaseEntry {
        int id
        releaseentryDOTType type
        timeDOTTime time
    }
    ReleaseLicense {
        int id
    }
    ReleasePolicy {
        int id
        string name
        string module
    }
    ReleasePolicyViolation {
        int id
        string message
        releasepolicyviolationDOTType type
        releasepolicyviolationDOTSeverity severity
    }
    ReleaseVulnerability {
        int id
    }
    Repo {
        int id
        string name
        string default_branch
    }
    TestCase {
        int id
        string name
        bool result
        string message
        float64 elapsed
        schemaDOTMetadata metadata
    }
    TestRun {
        int id
        string tool
        timeDOTTime time
        schemaDOTMetadata metadata
    }
    Vulnerability {
        int id
        string vid
        string summary
        string description
        float64 severity_score
        vulnerabilityDOTSeverity severity
        timeDOTTime published
        timeDOTTime modified
        schemaDOTMetadata metadata
    }
    VulnerabilityReview {
        int id
        string name
        vulnerabilityreviewDOTDecision decision
    }
    	Adapter }o--o| Organization : "owner"
    	Artifact }o--o| Release : "release/artifacts"
    	CodeIssue }o--o| CodeScan : "scan/issues"
    	CodeScan }o--o| Release : "release/code_scans"
    	Component }o--o| Organization : "owner"
    	Component }o--o{ Vulnerability : "vulnerabilities/components"
    	Component }o--o{ License : "licenses/components"
    	GitCommit }o--o| Repo : "repo/commits"
    	GitCommit |o--o| Release : "release/commit"
    	LicenseUse }o--o| License : "license/uses"
    	Project }o--o| Organization : "owner/projects"
    	Release }o--o{ Release : "dependencies/subreleases"
    	ReleaseComponent }o--o| Release : "release/components"
    	ReleaseComponent }o--o{ CodeScan : "scans/components"
    	ReleaseComponent }o--o| Component : "component/uses"
    	ReleaseEntry |o--o| Artifact : "artifact/entry"
    	ReleaseEntry |o--o| CodeScan : "code_scan/entry"
    	ReleaseEntry |o--o| TestRun : "test_run/entry"
    	ReleaseEntry }o--o| Release : "release/log"
    	ReleaseLicense }o--o| License : "license"
    	ReleaseLicense }o--o| ReleaseComponent : "component"
    	ReleaseLicense }o--o| Release : "release"
    	ReleaseLicense |o--o{ CodeScan : "scans"
    	ReleasePolicy }o--o| Organization : "owner"
    	ReleasePolicy }o--o{ Project : "projects/policies"
    	ReleasePolicy }o--o{ Repo : "repos/policies"
    	ReleasePolicyViolation }o--o| ReleasePolicy : "policy/violations"
    	ReleasePolicyViolation }o--o| Release : "release/violations"
    	ReleaseVulnerability }o--o| Vulnerability : "vulnerability/instances"
    	ReleaseVulnerability }o--o| ReleaseComponent : "component/vulnerabilities"
    	ReleaseVulnerability }o--o| Release : "release/vulnerabilities"
    	ReleaseVulnerability }o--o{ VulnerabilityReview : "reviews/instances"
    	ReleaseVulnerability }o--o| CodeScan : "scan/vulnerabilities"
    	Repo }o--o| Organization : "owner/repos"
    	Repo }o--o| Project : "project/repos"
    	Repo |o--o| Release : "head/head_of"
    	TestCase }o--o| TestRun : "run/tests"
    	TestRun }o--o| Release : "release/test_runs"
    	Vulnerability }o--o| Organization : "owner"
    	VulnerabilityReview }o--o| Vulnerability : "vulnerability/reviews"
    	VulnerabilityReview }o--o{ Project : "projects/vulnerability_reviews"
    	VulnerabilityReview }o--o{ Repo : "repos/vulnerability_reviews"
    	VulnerabilityReview }o--o{ Release : "releases/vulnerability_reviews"
`}/>

## Types

The types in the schema (AKA tables in the SQL schema) are listed below.
This list is auto-generated from the ent schema.
### Adapter

#### Fields
- **name** (string)
- **tag** (string)
- **module** (string)

#### Edges
- **owner** (M2O to [Organization](#organization))
### Artifact

#### Fields
- **name** (string)
- **sha256** (string)
- **type** (artifact.Type)
- **time** (time.Time)
- **metadata** (schema.Metadata)

#### Edges
- **release** (M2O to [Release](#release))
- **entry** (O2O to [ReleaseEntry](#releaseentry))
### CodeIssue

#### Fields
- **rule_id** (string)
- **message** (string)
- **severity** (codeissue.Severity)
- **type** (codeissue.Type)
- **metadata** (schema.Metadata)

#### Edges
- **scan** (M2O to [CodeScan](#codescan))
### CodeScan

#### Fields
- **tool** (string)
- **time** (time.Time)
- **metadata** (schema.Metadata)

#### Edges
- **release** (M2O to [Release](#release))
- **entry** (O2O to [ReleaseEntry](#releaseentry))
- **issues** (O2M to [CodeIssue](#codeissue))
- **vulnerabilities** (O2M to [ReleaseVulnerability](#releasevulnerability))
- **components** (M2M to [ReleaseComponent](#releasecomponent))
### Component

#### Fields
- **name** (string)
- **vendor** (string)
- **version** (string)
- **description** (string)
- **url** (string)
- **metadata** (schema.Metadata)

#### Edges
- **owner** (M2O to [Organization](#organization))
- **vulnerabilities** (M2M to [Vulnerability](#vulnerability))
- **licenses** (M2M to [License](#license))
- **uses** (O2M to [ReleaseComponent](#releasecomponent))
### GitCommit

#### Fields
- **hash** (string)
- **branch** (string)
- **tag** (string)
- **time** (time.Time)

#### Edges
- **repo** (M2O to [Repo](#repo))
- **release** (O2O to [Release](#release))
### License

#### Fields
- **spdx_id** (string)
- **name** (string)
- **reference** (string)
- **details_url** (string)
- **is_osi_approved** (bool)

#### Edges
- **components** (M2M to [Component](#component))
- **uses** (O2M to [LicenseUse](#licenseuse))
### LicenseUse

#### Fields

#### Edges
- **license** (M2O to [License](#license))
### Organization

#### Fields
- **name** (string)

#### Edges
- **projects** (O2M to [Project](#project))
- **repos** (O2M to [Repo](#repo))
### Project

#### Fields
- **name** (string)

#### Edges
- **owner** (M2O to [Organization](#organization))
- **repos** (O2M to [Repo](#repo))
- **vulnerability_reviews** (M2M to [VulnerabilityReview](#vulnerabilityreview))
- **policies** (M2M to [ReleasePolicy](#releasepolicy))
### Release

#### Fields
- **name** (string)
- **version** (string)
- **status** (release.Status)

#### Edges
- **subreleases** (M2M to [Release](#release))
- **dependencies** (M2M to [Release](#release))
- **commit** (O2O to [GitCommit](#gitcommit))
- **head_of** (O2O to [Repo](#repo))
- **log** (O2M to [ReleaseEntry](#releaseentry))
- **violations** (O2M to [ReleasePolicyViolation](#releasepolicyviolation))
- **artifacts** (O2M to [Artifact](#artifact))
- **components** (O2M to [ReleaseComponent](#releasecomponent))
- **vulnerabilities** (O2M to [ReleaseVulnerability](#releasevulnerability))
- **code_scans** (O2M to [CodeScan](#codescan))
- **test_runs** (O2M to [TestRun](#testrun))
- **vulnerability_reviews** (M2M to [VulnerabilityReview](#vulnerabilityreview))
### ReleaseComponent

#### Fields
- **type** (releasecomponent.Type)

#### Edges
- **release** (M2O to [Release](#release))
- **scans** (M2M to [CodeScan](#codescan))
- **component** (M2O to [Component](#component))
- **vulnerabilities** (O2M to [ReleaseVulnerability](#releasevulnerability))
### ReleaseEntry

#### Fields
- **type** (releaseentry.Type)
- **time** (time.Time)

#### Edges
- **artifact** (O2O to [Artifact](#artifact))
- **code_scan** (O2O to [CodeScan](#codescan))
- **test_run** (O2O to [TestRun](#testrun))
- **release** (M2O to [Release](#release))
### ReleaseLicense

#### Fields

#### Edges
- **license** (M2O to [License](#license))
- **component** (M2O to [ReleaseComponent](#releasecomponent))
- **release** (M2O to [Release](#release))
- **scans** (O2M to [CodeScan](#codescan))
### ReleasePolicy

#### Fields
- **name** (string)
- **module** (string)

#### Edges
- **owner** (M2O to [Organization](#organization))
- **projects** (M2M to [Project](#project))
- **repos** (M2M to [Repo](#repo))
- **violations** (O2M to [ReleasePolicyViolation](#releasepolicyviolation))
### ReleasePolicyViolation

#### Fields
- **message** (string)
- **type** (releasepolicyviolation.Type)
- **severity** (releasepolicyviolation.Severity)

#### Edges
- **policy** (M2O to [ReleasePolicy](#releasepolicy))
- **release** (M2O to [Release](#release))
### ReleaseVulnerability

#### Fields

#### Edges
- **vulnerability** (M2O to [Vulnerability](#vulnerability))
- **component** (M2O to [ReleaseComponent](#releasecomponent))
- **release** (M2O to [Release](#release))
- **reviews** (M2M to [VulnerabilityReview](#vulnerabilityreview))
- **scan** (M2O to [CodeScan](#codescan))
### Repo

#### Fields
- **name** (string)
- **default_branch** (string)

#### Edges
- **owner** (M2O to [Organization](#organization))
- **project** (M2O to [Project](#project))
- **head** (O2O to [Release](#release))
- **commits** (O2M to [GitCommit](#gitcommit))
- **vulnerability_reviews** (M2M to [VulnerabilityReview](#vulnerabilityreview))
- **policies** (M2M to [ReleasePolicy](#releasepolicy))
### TestCase

#### Fields
- **name** (string)
- **result** (bool)
- **message** (string)
- **elapsed** (float64)
- **metadata** (schema.Metadata)

#### Edges
- **run** (M2O to [TestRun](#testrun))
### TestRun

#### Fields
- **tool** (string)
- **time** (time.Time)
- **metadata** (schema.Metadata)

#### Edges
- **release** (M2O to [Release](#release))
- **entry** (O2O to [ReleaseEntry](#releaseentry))
- **tests** (O2M to [TestCase](#testcase))
### Vulnerability

#### Fields
- **vid** (string)
- **summary** (string)
- **description** (string)
- **severity_score** (float64)
- **severity** (vulnerability.Severity)
- **published** (time.Time)
- **modified** (time.Time)
- **metadata** (schema.Metadata)

#### Edges
- **owner** (M2O to [Organization](#organization))
- **components** (M2M to [Component](#component))
- **reviews** (O2M to [VulnerabilityReview](#vulnerabilityreview))
- **instances** (O2M to [ReleaseVulnerability](#releasevulnerability))
### VulnerabilityReview

#### Fields
- **name** (string)
- **decision** (vulnerabilityreview.Decision)

#### Edges
- **vulnerability** (M2O to [Vulnerability](#vulnerability))
- **projects** (M2M to [Project](#project))
- **repos** (M2M to [Repo](#repo))
- **releases** (M2M to [Release](#release))
- **instances** (M2M to [ReleaseVulnerability](#releasevulnerability))
