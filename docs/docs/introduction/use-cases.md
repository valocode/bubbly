---
title: Use Cases
hide_title: false
hide_table_of_contents: false
description: Use Cases of Bubbly
keywords:
- docs
- bubbly
- use-cases
---

This page describes a few of the use cases that we have applied Bubbly to.
They are documented here for readers to get an idea of *how* bubbly can help, and the kinds of problems that it solves.
If you currently suffer from similar pains then perhaps that suggest you should take Bubbly for a spin ;)

## Results & Policies

### Results

Capturing results from different tools (linters, static code analysis, software composition analysis, different levels of test automation, vulnerabilities, software licenses, etc.) is one of the core use cases for Bubbly.
The idea is not to duplicate your already detailed source of these, but to capture the top-level data to help build useful KPIs and drive the Release Readiness decision.
For example, if you have test results with logs and profiling information somewhere already - great! Don't duplicate that, just pull the information you want to write policies about and bring under Bubbly's hood for trends and analytics.
Similarly with static code analysis results - Bubbly won't allow you to manage state of those or view the annotated source code with the issue. So use your existing tools for that (if you need to).

### Policies

One you collect these results into Bubbly, it now gets interesting to start applying some policies to your release process.
Policies define what is `required` in your release process (e.g. "tool xyz must run") and what is `denied` in your release process (e.g. "code issues with a high severity are not allowed").
When the policy engine runs, violations will be created for each violation of the rules you have defined.

We are fortunate enough to live in a world where someone already solved this problem amazingly, and so we are building on the great work of Open Policy Agent and using their [Rego Policy Language](https://www.openpolicyagent.org/docs/latest/policy-language/).

## OSS License and Security

The topic of OSS License and Security is not one we want to dig into too deeply here, because we will go head first into a rabbit hole that would take a lot of work for us to find the exit.
However, in short, we have found that there is friction between topics like 3rd Party (License) Clearing and teams wanting to adopt DevOps-practices and move quickly and we thought giving this topic extra attention would try to solve this.

The tools in the industry to support these activities have been dominated in the past decade by commercial vendors, such as [BlackDuck](https://www.blackducksoftware.com/) and [Revenera](https://www.revenera.com/protect/products.html) (originally Palamida), with more newcomers like [Snyk](https://snyk.io/), [Whitesource](https://www.whitesourcesoftware.com/), [Fossa](https://fossa.com/)... the list goes on.

There is now a growth in OSS tools to help with these activities, and initiatives like [DoubleOpen](https://www.doubleopen.org/) and projects like [SW360](https://www.eclipse.org/sw360/), [Fossology](https://www.fossology.org/) and [OSS Review Toolkit](https://github.com/oss-review-toolkit/ort) exist, which is amazing.

These are a mix of "scanners" (basically Software Composition Analysis tools that give you a Bill of Material and can report licenses and vulnerabilities) and catalogue management tools, that help you manage your components and licenses.

Bubbly, like in all cases, will serve only as a datastore for results and the Bubbly Adapters make it possible to get data from any tools/format and Bubbly Policies make it possible to write very flexible rules for what is required and denied.

## Continuous Improvement

Another goal of Bubbly from the start was to help teams who hit their "CI/CD plateau"...
From experience we have seen that the initial gain in implementing these Continuous practices is fairly significant, but after a while the question becomes "what next?"
All the studies show that teams continuing to improve will be better performers.

Bubbly profiles your release process, with timestamps and results from tools, and therefore knows:

1. When did the events happen, with durations between process
2. How often do releases fail, and for what reason / violation

From this data, it is possible to derive insights that can help teams identify bottlenecks or common weaknesses in their release process, and therefore understand (justified with data) where to invest in improvements.

:::tip Under Active Development
The use case for Continuous Improvement is something we are continually improving... (pardon the pun).

We have made some initial work but are looking for more real-world use cases to make this something "out of the box".
So let us know if you are interested and the Bubbly team could work together with you :)
:::
