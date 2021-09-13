---
title: Core Concepts
hide_title: false
hide_table_of_contents: false
description: Core Concepts of bubbly
keywords:
- docs
- bubbly
- core concepts
---

:::note Under Active Development
The content for this page is *under active development*. It
may be
incorrect and/or
incomplete and is liable to change at any time.
:::

## Release

The most important entity in Bubbly's schema is the `Release` entity.
It is the "connection point" for all the data that we want to associate with our Release Readiness decision.
Therefore it is very important to know that a release is just an entity, but at the heart of Bubbly.

Check it out in the [schema](../schema/schema#release)

## Adapters

Getting data in to Bubbly was one of the main challenges we wanted to solve.
Adapters were created to solve this.

Adapters are simply [Rego](https://www.openpolicyagent.org/docs/latest/policy-language/) files, that define some results from an input.
E.g. you can parse a JSON file and convert that into `code_issues` or a `test_run` with `test_cases`.

Check the [documentation](../adapters/adapters.md)

## Policies

Policies are evaluated against each release and form the brains of the Continuous Release Readiness process.
Like Adapters, Policies are simply [Rego](https://www.openpolicyagent.org/docs/latest/policy-language/) files that define some rules.
E.g. you can `require` that some event occurs in the release process, or `deny` certain results if they are not good enough.

Check the [documentation](../policies/policies.md)

## GraphQL API

[GraphQL](https://graphql.org/) is a query language which is very simple to use, and great for getting structured data out from an API.

For example, to get the `head` (which is the latest commit on the main branch) for a repository called `"demo"` we use the following query:

```graphql
{
 release(where: {has_head_of_with: {name: "demo"}}) {
    name
    version
  }
}
```

Check the [documentation](../graphql/graphql.md) for more examples

## Schema

The Bubbly schema has been, and will continue to be, one of the most important parts of the development of Bubbly.
As Bubbly is not trying to replace all the rich data sources you might already have, but instead provide a higher-level decision making capability on the data, the schema has to be "just right" in that it doesn't require too much data, but also captures enough to be useful.

One of the main ways in which we solved this simply was to provide a `Metadata` type, which is basically raw JSON, so that you *can* add extra data to almost all entities.

Check the [documentation](../schema/schema.md) which was auto generated from the real schema.
