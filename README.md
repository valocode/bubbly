
![bubbly-logo](/docs/static/img/bubbly-blue-wide.png)

> Release Readiness in a Bubble.

[![Go Report Card](https://goreportcard.com/badge/github.com/valocode/bubbly)](https://goreportcard.com/report/github.com/valocode/bubbly)
[![Mozilla Public License Version 2.0](https://img.shields.io/github/license/valocode/bubbly?color=brightgreen&label=License)](https://opensource.org/licenses/MPL-2.0)
[![Bubbly on Twitter](https://img.shields.io/badge/Follow-bubblydotdev-blue.svg?style=flat&logo=twitter)](https://twitter.com/intent/follow?screen_name=bubblydotdev)

Bubbly is a release readiness platform helping software teams release compliant software with confidence. Gain visibility into your release process with reports and analytics to lower risk, increase quality, reduce cycle time and drive continuous improvement.

**Release Readiness** is a term that we use to define the state of being ready to release software. This project aims at helping teams help teams define their Release Readiness policies (using the [Rego Policy Language](https://www.openpolicyagent.org/docs/latest/policy-language/)), collect the data to enforce those policies, and measure their performance getting to their Release Readiness goals (whether it be speed or stability).

![bubbly-in-a-bubble](/docs/static/img/bubbly-in-a-bubble.svg)

## Project Status

This project is currently in **ALPHA**

This means that it should work, and if it doesn't please raise an issue. The APIs may change in the coming weeks/months.

If you are interested in getting involved in the project, or want some help getting started, please reach out to us though our website: [https://bubbly.dev](https://bubbly.dev)

## Problem Statement

The problem around Release Readiness is one of data. All the tools used in the software process produce data, and that data should be used to drive the Release Readiness decision.

The real problem is that this data gets scattered all over the place which results in:

1. Broken relationships across the data
2. Non-standard interfaces which make accessing the data hard (sometimes impossible)
3. Overhead maintaining multiple data stores and data bases, ensuring that the data is up to date and accurate
4. Lack of understanding of the data's hierarchy

## Proposed Solution

Bubbly has been built to address the core problems mentioned above. This has been achieved by implementing a simple engine built over the [Rego Policy Language](https://www.openpolicyagent.org/docs/latest/policy-language/) with a versatile database schema, built using [entgo](https://entgo.io/).

A GraphQL and standard REST API are provided to help get data out, as one design goal of bubbly is putting the right data in the right place.

The data in Bubbly can now be used to make automated decisions in pipelines (e.g. check that results are good) and can be used to power dashboards. Any dashboard that supports GraphQL (like [Grafana](https://grafana.com/grafana/)) can be used to visualize the data. Bubbly also ships with its own UI, and this was born out of the need to understand the data hierarchy (something that ordinary dashboards do not do too well with).

Check out the [Core Concepts](https://docs.bubbly.dev/introduction/core-concepts) for more information on things like the adapters, policies and more.

### Goals

The goal of Bubbly was to be a "bridge" between data, and make connecting and accessing data related to release readiness as simple, reliable and fun as possible!

We are not aiming to replace every tool you are already using, in fact, without other tools Bubbly itself is useless. The goal is to complement every other tool by making their data more accessible.

Some ways in which Bubbly could be used are:

* Real-time knowledgebase of all projects, repos and their high-level release readiness state.
* Catalogue of OSS components, licenses, Vulnerabilities and approvals/rejections for those.
* A decision protocol to formally define your Release Readiness Policy, and automatically apply that to all projects.

If you have some ideas or questions for things you want to solve, please raise an issue in the [Discussion](https://github.com/valocode/bubbly/discussions) and we'll check it out :)

### Non-Goals

The things that Bubbly is not aiming to be:

* A testing tool/framework.
* A traditional monitoring tool (think Ops metrics, like Prometheus).
* Another dashboard (like Grafana).
  
## Architecture

See the [ARCHITECTURE.md](./ARCHITECTURE.md) file for the structure of this repository.

## Why Open Source?

Bubbly is open source, licensed under the [Mozilla Public License v2](LICENSE).

We have a few reasons for being open source:

1. The team and company behind Bubbly are all big fans of open source, so it was never really considered not making it open source.
2. We want to focus on building the best product, and we firmly believe making the project OSS will help us achieve this.

## Install

See our [Installation Guide](https://docs.bubbly.dev/getting-started/getting-started#installation).

## Getting Started

See our [Getting Started Guide](https://docs.bubbly.dev/getting-started/getting-started).

## Contributing

See our [Contributing Guide](./docs/docs/contributing/contributing.md).

## License

[Mozilla Public License Version 2.0](LICENSE)
