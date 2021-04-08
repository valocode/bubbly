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

## Test Automation

Test automation refers to the typical levels of automated testing companies setup, usually as part of their CI (Continuous Integration) pipelines.

It covers both **functional** and **non-functional** tests (e.g. testing features vs testing stability over time) and the different levels that usually fall into **unit**, **integration** and **system** testing.
Tests such as "User Acceptance Tests" can be a set of tests that cover a number of these types of tests, and this is a good intro to the problem description.

### Problem

The different levels of tests (unit, integration and system) produce different test results, and are usually performed by different test frameworks.
Additionally, functional tests usually product simple pass/fail type results, where are non-functional tests can vary immensely (capturing network speed, memory usage, cpu usage, page refresh, etc).
You might also want to capture other metrics, like code coverage, as well as tracing your test cases back to features or requirements.

There are tools purpose built for test result management, e.g. [Allure](http://allure.qatools.ru/) and [TestRail](https://www.gurock.com/testrail/), and there are other tools that can be customized to do this, e.g. [Jira](https://www.atlassian.com/software/jira) or (please don't do this) [Jenkins](https://www.jenkins.io/).

The challenge is that they are very good at capturing pass/fail style results, and maybe code coverage, but they fall short or become cumbersome when it comes to anything beyond that.
Especially if we have some time-series databases storing non-functional results from test runs, how can we use that data in our results together with pass/fail?

### Solution

With Bubbly you can define schemas to capture the results that you care about, and make the relationships that you care about, and then use the GraphQL API to write powerful queries about the data.

The obvious benefits here are that you would have **one API** to query for your test results.
This includes adding notifications to alert on failure, building dashboards to view your results, and adding some kind of quality gate into your (automated) release process.

This helps unify results and utlimately provide faster feedback to your team.

As with all things Bubbly, we are not suggesting that we replace any of the tools you use already, but we provide the umbrella around them.

E.g. Bubbly works great with [Jira and Xray](https://www.getxray.app/) to store test cases and results, together with [InfluxDB](https://www.influxdata.com/) to store test run profiling data.
We can also capture some more detailed data after the test run, and aggregate that into one place.

## OSS License and Security

The topic of OSS License and Security is not one we want to dig into too deeply here, because we will go head first into a rabbit hole that would take a lot of work for us to find the exit.
However, in short, we have found that there is friction between topics like 3rd Party (License) Clearing and teams wanting to adopt DevOps-practices and move quickly.
There is not so much science behind this (yet) but the friction is caused because the people involved in 3rd Party (and by 3rd party we mean OSS + commercial software) goes from software to management to legal (and back).
And as we all know, cross-team processes are often great ways to slow things down (why did DevOps become such a hot topic?).
The challenge is about tools than can enable processes and for each stakeholder to get what they need.

The tools in the industry to support these activities have been absolutely dominated in the past decade by commercial vendors, such as [BlackDuck](https://www.blackducksoftware.com/) and [Revenera](https://www.revenera.com/protect/products.html) (originally Palamida), with more newcomers like [Snyk](https://snyk.io/), [Whitesource](https://www.whitesourcesoftware.com/), [Fossa](https://fossa.com/)... the list goes on.

There is now a growth in OSS tools to help with these activities, and initiatives like [DoubleOpen](https://www.doubleopen.org/) and projects like [SW360](https://www.eclipse.org/sw360/), [Fossology](https://www.fossology.org/) and [OSS Review Toolkit](https://github.com/oss-review-toolkit/ort) exist, which is amazing.

These are a mix of "scanners" (basically Software Composition Analysis tools that give you a Bill of Material) and catalogue management tools, that help you manage your components and licenses.

### Problem

So what's the problem if there exist both commercial and OSS tools to help with this?
Having worked with companies wanting to streamline this process we found that there is a gap in connecting these tools and using them to drive CI/CD activities.

And this is exactly where Bubbly wants to come in.

### Solution

The Bubbly schema can be used to capture and store all the data you want to collect in order for the stakeholders to be happy:

1. Development teams will be the ones using the libraries to implement features
2. Product owners / managers will want to know that they are compliant and secure
3. Security teams will want an up-to-date list of vulnerabilities
4. If there is a legal team, they will want to have the necessary data for them to approve the usage of OSS

Bubbly is able to collect all this information into one database, and provide the queries necessary to extract the data.

