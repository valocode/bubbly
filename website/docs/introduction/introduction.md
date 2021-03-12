---
title: Introduction
hide_title: false
hide_table_of_contents: false
slug: /
description: Introduction to bubbly
keywords:
- docs
- bubbly
- introduction
---

![bubbly-logo](/img/bubbly-blue-wide.png)

## Bubbly - Release Readiness in a Bubble

Bubbly emerged from a need that many lean software teams practicing Continuous Integration and Delivery have: the need to be better performers. The idea is quite simple: build a model that can capture all the metrics that you care about, extract data from any tool in your process, and use those metrics to drive better and more automated decisions.

**Release Readiness** is a term that we use to define *the state of being ready to release software*.
What is interesting to consider is how different teams define their release readiness.
Often teams have planned to implement some features and run automated quality and security checks throughout the process (static analysis, automated tests, CVE checks, etc.).
Ultimately the decision of being release ready comes down to a measure of confidence.
So how best to inform this measure of confidence in an efficient and reliable way?

## Problem

The problem around Release Readiness is one of data.
All the tools used in the software process produce data, and that data should be used to measure attributes such as feature-completeness, quality and security of your product.

The real problem is that this data gets scattered all over the place which results in:

1. Broken relationships across the data
2. Non-standard interfaces which make accessing the data hard (sometimes impossible)
3. Massive overhead maintaining multiple data stores and data bases, ensuring that the data is up to date and accurate
4. Lack of understanding of the data's hierarchy

## Proposed Solution

Bubbly has been built to address the core problems mentioned above.
This has been achieved by implementing a very lightweight *data platform* using a user-defined schema and data pipelines, with everything defined as code using a DSL built on top of HashiCorp's Configuration Language (HCL).

A user defines a schema that models the metrics and Key Performance Indicators (KPIs) that you care about. This could include automated test results, OSS components, CVEs, development metrics, traceability across requirements, and anything else that you can think of :)

Data pipelines extract the data from a data source (i.e. the original source of the data), transform the data into your schema and load the results into a persistent database.

A GraphQL API is automatically generated based on the user-defined schema, which is the main way of getting data out from Bubbly.
What's noteworthy is that the generated GraphQL API is extended to include some cool tricks to make life that litle bit better :)

The data in Bubbly can now be used to make automated decisions in pipelines (e.g. check that results are good) and can be used to power dashboards.
Any dashboard that supports GraphQL (like [Grafana](https://grafana.com/grafana/)) can be used to visualize the data.
Bubbly also ships with its own UI, and this was born out of the need to understand the data hierarchy (something that ordinary dashboards do not do too well with).

Check out the [Core Concepts](./core-concepts.md) for more information on things like the schema, HCL and data pipelines.

### Goals

The goal of Bubbly was to be a "bridge" between data, and make connecting and accessing data as simple, reliable and fun as possible! Writing data pipelines feels awesome!

We are not aiming to replace every tool you are already using, in fact, without other tools Bubbly itself is useless.
The goal is to complement every other tool by making their data more accessible.

Some ways in which Bubbly could be used are:

1. Real-time knowledgebase of all products and their high-level release readiness state
2. Catalogue of OSS components, licenses, CVEs and approvals/rejections for those
3. A decision protocol to formally define your Release Readiness Criteria, and automatically apply that to all products

If you want to see some use cases that we are currently solving with Bubbly, check out the [Use Cases](./use-cases.md).

### Non-Goals

The things that Bubbly is not aiming to be:

1. A testing tool/framework
2. A traditional monitoring tool (think Ops metrics, like Prometheus)
3. Another dashboard (like Grafana)
   1. We like to make the distinction between a "dashboard" which typically has some specific views to visualize data, vs a "knowledgebase" which is more like a web app that has multiple dashboards inside it

### Architecture

TODO: this is in progress... Check back soon!
