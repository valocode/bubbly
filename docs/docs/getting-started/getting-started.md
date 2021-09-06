---
title: Getting Started
hide_title: false
hide_table_of_contents: false
description: Installation of bubbly
keywords:
- docs
- bubbly
- installation
- deployment
- agent
---


## Installation

Bubbly is a single Go binary containing botht the server, client and UI. Run it using Docker or download and run natively.

### Docker

`docker run valocode/bubbly:<version>`

### Homebrew

`brew install valocode/tap/bubbly`

### Releases

Binaries for macOS, Windows and Linux AMD64/ARM are available to download on the
[release page](https://github.com/valocode/bubbly/releases)

## Deployment

The choice of database is important.
Bubbly uses [entgo](https://entgo.io) which has the following [supported dialects](https://entgo.io/docs/dialects/).

Entgo supports SQLite, and be default Bubbly uses an in-memory SQLite database.
So if you just want to spin up Bubbly and try it, don't worry about a database choice for now but do not get angry when your data is not persisted... :)

Once you have chosen your database, run bubbly with your preferred method:

### Docker

If using docker, it's easy to integrate Bubbly with any orchestrator.
Either provide command line arguments or set environment variables to configure your backend database.

### From binary

```bash
# run a demo instance of bubbly, which uses sqlite and has some dummy data so that
# you can immediately see something
bubbly demo

# view the help to know the flags
bubbly server --help

# start the bubbly server
bubbly server

# skip the ui
bubbly server --ui=false
```

Check the [documentation](../cli/bubbly_server)

## Basic Commands

There are a few main actions/commands in Bubbly.

### Creating a Release

Creating a release is done via the command line, so that it can be integrated into CI.

```bash
# view the help for creating a release
bubbly release create --help

# create a release using the local .git directory
bubbly release create
```

This command will check for a `.bubbly.json` file which can be used to override the default release options (name and version).

Check the [documentation](../cli/bubbly_release_create)

### Applying a Policy

Applying a policy is the way to define your Release Readiness policy.

A policy is a set of rules, that either `require` certain events in the release log, or `deny` certain results.
Each time a rule is violated, a violation is created against your release and the violation can be customised by the rule, such as setting the severity(`suggestion`, `warning` or `blocker`).

Below is a policy that looks for `code_issues` with a severity of `high` and creates a violation if there are *any* such code issues.

```rego
package policy

deny[violation] {
 issues := code_issues()
 high_issues := [issue | issues[i].severity == "high"; issue := issues[i]]
 count(high_issues) > 0
 violation = {
  "message": sprintf("%d high issue(s)", [count(high_issues)]),
  "severity": "warning",
 }
}
```

To save this policy, put it in a file called `code_issue_high_severity.rego` and use the Bubbly command line:

```bash
# view the help for bubbly policy command
bubbly policy --help

# save the policy and apply it to the project "default"
bubbly policy save code_issue_high_severity.rego --set-projects default

# alternatively do this in two steps: first save and then associate
bubbly policy save code_issue_high_severity.rego
bubbly policy set code_issue_high_severity --projects default
```

This policy now affects/applies to the project default.
Any releases created within a repository belonging to the project default will inherit this policy.

### Running an adapter

Running an adapter is the way to get results into Bubbly. First you need an adapter to run.

As a short example, here is an adapter for import [gosec](https://github.com/securego/gosec) results.

```rego
package adapter

code_scan[scan] {
 scan := {
  "tool": "gosec",
  "metadata": {"env": {"some_var": "some_value"}},
 }
}

code_issue[issue] {
 some i
 iss := input[_].Issues[i]
 issue := {
  # providing the i is necessary so that we get all unique code_issues
  "i": i,
  "rule_id": iss.rule_id,
  "message": iss.details,
  "severity": lower(iss.severity),
  "type": "security",
 }
}
```

If we save this in a file called `gosec.rego` and then use the bubbly command line we can either run the adapter directly from the file, or save it and run by fetching it from the remote server.

We will need a results file as input also, which can be found here (this is a result from running gosec over bubbly some time back): <https://raw.githubusercontent.com/valocode/bubbly/main/adapter/testdata/adapters/gosec.json>

```bash
# run the adapter by providing the path to it as input
bubbly adapter run ./gosec.rego
```

This command checks for a local git repository and uses the `HEAD` commit to associate the results with.
A release must exist for this commit.

Check the [documentation](../cli/bubbly_adapter_run.md)

### View Release

Either do this through the Bubbly UI or via the command line:

```bash
# view the local release
bubbly release view

# view the local release with policies that apply to it
bubbly release view --policies
```
