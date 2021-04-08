
# Bubbly

![bubbly-logo](/docs/static/img/bubbly-blue-wide.png)

Bubbly - Release Readiness in a Bubble

Bubbly emerged from a need that many lean software teams practicing Continuous Integration and Delivery have: the need to be better performers. The idea is quite simple: build a model that can capture all the metrics that you care about, extract data from any tool in your process, and use those metrics to drive better and more automated decisions.

**Release Readiness** is a term that we use to define the state of being ready to release software. What is interesting to consider is how different teams define their release readiness. Often teams have planned to implement some features and run automated quality and security checks throughout the process (static analysis, automated tests, CVE checks, etc.). Ultimately the decision of being release ready comes down to a measure of confidence. So how best to inform this measure of confidence in an efficient and reliable way?

## Bubbly Status

[![Go Report Card](https://goreportcard.com/badge/github.com/valocode/bubbly)](https://goreportcard.com/report/github.com/valocode/bubbly)

### Problem

The problem around Release Readiness is one of data. All the tools used in the software process produce data, and that data should be used to measure attributes such as feature-completeness, quality and security of your product.

The real problem is that this data gets scattered all over the place which results in:

1. Broken relationships across the data
2. Non-standard interfaces which make accessing the data hard (sometimes impossible)
3. Massive overhead maintaining multiple data stores and data bases, ensuring that the data is up to date and accurate
4. Lack of understanding of the data's hierarchy

### Proposed Solution

Bubbly has been built to address the core problems mentioned above. This has been achieved by implementing a very lightweight data platform using a user-defined schema and data pipelines, with everything defined as code using a DSL built on top of HashiCorp's Configuration Language (HCL).

A user defines a schema that models the metrics and Key Performance Indicators (KPIs) that you care about. This could include automated test results, OSS components, CVEs, development metrics, traceability across requirements, and anything else that you can think of

Data pipelines extract the data from a data source (i.e. the original source of the data), transform the data into your schema and load the results into a persistent database.

A GraphQL API is automatically generated based on the user-defined schema, which is the main way of getting data out from Bubbly. What's noteworthy is that the generated GraphQL API is extended to include some cool tricks to make life that litle bit better

The data in Bubbly can now be used to make automated decisions in pipelines (e.g. check that results are good) and can be used to power dashboards. Any dashboard that supports GraphQL (like [Grafana](https://grafana.com/grafana/)) can be used to visualize the data. Bubbly also ships with its own UI, and this was born out of the need to understand the data hierarchy (something that ordinary dashboards do not do too well with).

Check out the [Core Concepts](https://docs.bubbly.dev/introduction/core-concepts) for more information on things like the schema, HCL and data pipelines.

#### Goals

The goal of Bubbly was to be a "bridge" between data, and make connecting and accessing data as simple, reliable and fun as possible! Writing data pipelines feels awesome!

We are not aiming to replace every tool you are already using, in fact, without other tools Bubbly itself is useless. The goal is to complement every other tool by making their data more accessible.

Some ways in which Bubbly could be used are:

* Real-time knowledgebase of all products and their high-level release readiness state.
* Catalogue of OSS components, licenses, CVEs and approvals/rejections for those.
* A decision protocol to formally define your Release Readiness Criteria, and automatically apply that to all products.

If you want to see some use cases that we are currently solving with Bubbly, check out the Use Cases.

#### Non-Goals

The things that Bubbly is not aiming to be:

* A testing tool/framework.
* A traditional monitoring tool (think Ops metrics, like Prometheus).
* Another dashboard (like Grafana).
  
  We like to make the distinction between a "dashboard" which typically has some specific views to visualize data, vs a "knowledgebase" which is more like a web app that has multiple dashboards inside it.

## Architecture

TODO: this is in progress... Check back soon!

## Why Open Source?

Bubbly is open source, licensed under the [Mozilla Public License v2](https://www.mozilla.org/en-US/MPL/2.0/).

We have a few reasons for being open source:

1. The team and company behind Bubbly are all big fans of open source, so it was never really considered not making it open source.
2. We want to focus on building the best product, and we firmly believe making the project OSS will help us achieve this.
3. We want to enable teams to be better. And that's why we are launching a SaaS service for Bubbly to get you up and running in minutes.

## Alternatives

Bubbly was born due to (what we believe is) a gap in the market.
For each project, when we needed to solve this issue we would re-invent something and build that on top of amazing OSS projects.

Hence, the only real alternative that we know of is building it yourself with something like [Elastic](https://www.elastic.co/) and [Elastic Beats](https://www.elastic.co/beats/), or [InfluxDB Telegraf](https://www.influxdata.com/time-series-platform/telegraf/).
You can of course roll your own SQL database, and define a schema (or schema migrator) and implement your own data access layer, which then needs to provide an SDK or a DSL in order to use, which means picking a language and implementing an API or a parser, and setup your own data-driven architecture using a pub/sub model or a message queue, deployed with security and auth tokens, coupled with a CLI... (do you see where this is going?)

Of course it can be a very viable approach to build these things yourself, but what often happens in internal projects that start small, is that they snowball into something that needs to be reused in other projects and features get added on top (not designed from within), and eventually you end up developing a separate product that has it's own release schedules and sophisticated software archtiecture, when really what you should have been focusing on all along is your company's product.

Let Bubbly be our product that you can apply to solve these problems, and if you need extra features or support then let us know!
We might already be working on the feature that you need, or can prioritise it if we know people need it, and of course we are open to being sponsored and helping you to deploy Bubbly.

And as it is open source we really welcome contributions!

## Install

See our [Installation Guide](https://docs.bubbly.dev/getting-started/getting-started#installation)

## Getting Started

See our [Getting Started Guide](https://docs.bubbly.dev/getting-started/getting-started)

## Contributing

See our [Contributing Guide](./docs/docs/contributing/contributing.md)

## License

[Mozilla Public License Version 2.0](https://www.mozilla.org/en-US/MPL/)
