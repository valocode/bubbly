# Bubbly Specification

**Bubbly - Release Readiness in a bubble**

Bubbly provides a declarative approach for defining metrics associated with *release readiness* to be aggregated and queried to objectively answer the important question: *"Are we ready for a release?"*

This specification is more concerned with describing what Bubbly is and how it will be used.
This specification is not concerned with how Bubbly will be implemented.
For that, see the [implementation documentation](./IMPLEMENTATION.md).

## 1. Introduction

### 1.1 Overview

All software development projects create data regarding the release readiness of the software.
Simple examples like running automated tests and security scans as part of Continuous Integration create results that need to be monitored and acted upon, and those results need to be associated with some data source (such as source code, Docker images, or even a production environment).

Tools exist for aggregating this data but we have experienced a need to *re-invent the wheel* for every project because of different needs, different tools and the fact that the tools we use to aggregate this data were never really made for this exact purpose.
Typical databases that we see tailored for this approach include ElasticSearch, PostgreSQL and Prometheus.
This requires integrations to be developed to populate these databases and teams end up maintaining a collection of utlity scripts, often very specific to the team or organisation.

The goal of Bubbly is to drastically minimise the need to re-invent the wheel each time, and minimise the burden of maintaining lots of scripts and utilities to aggregate data into a queryable backend.

### 1.2 Terminology

- Bubbly - name of this project
- HCL - [HashiCorp Configuration Language](https://github.com/hashicorp/hcl)
- CI - Continuous Integration
- API - Application Programming Interface

### 1.3 Project Contribution

We consider this problem worth solving because we want to make metrics associated with *release readiness* actionable through an API, so that they can be built into the development process.

Development teams will be able to integrate Bubbly into their release process, whether automated or not, allowing data related to release readiness to be collected and aggregated into a single, queryable interface.
Then, whether through an automated deployment or a manual delivery of software, the API to Bubbly should be able to objectively answer a very important question: *"Are we ready for a release?"*

An additional benefit to aggregating data into a queryable API is that dashboards can be created using tools like Grafana and Kibana to visualize any particular metrics of interest.

## 2. Techical Requirements

The list of technical requirements are very high level, and the Architecture and Design sections will provide more detail on how these requirements are implemented.

Bubbly will provide 3 main requirements:

1. Data conversion - the ability to convert data into different structures/types using declarative HCL
2. Data digestion - the ability to digest and store data in a specified backend database
3. Data querying - the ability to query the stored data for visualization and calculation of the different metrics for release readiness

### 2.1. Data Conversion

This requirement is concerned with being able to take lots of different formats of data and convert them into a defined common structure.
The goal with Bubbly is to make this easy, maintainable and reliable.
As such, the HCL declarative language is used to define the data conversion and bubbly takes care of the heavy lifting and error handling - consider it a reusable library or framework.

### 2.2. Data Digestion

### 2.2. Data Querying

## 3. Architecture

The high-level architecture for Bubbly consist of two parts:

1. The **Bubbly Server** is a long-running backend process that is available to upload new data or return data based on a provided query

2. The **Bubbly Client** is executed in an automated manner (e.g. CI pipeline) or manually by a user to perform operations agains the Bubbly Server

These main components can be illustrated using the following simple diagram:

![Bubbly High Level Architecture](./images/high-level-arch.drawio.svg)

The other relevant parts in this diagram are:

1. The **Client Configs** tell the Bubbly client everything it needs to know about what to do. This includes: the location of the *Bubbly server*, where the *input data* is, what to upload to the *Bubbly server*, etc. All the client configs are defined in HCL files.
2. The **Input Data** is the data which should be uploaded to the Bubbly server and can come from any source, such as a JSON or XML file, or by querying a REST API or by executing a command line tool. Before it is uploaded it needs to be converted into the expected format, which is defined by the *schema*, and the *importers* are used to convert the data.
3. The **Schema** defines the format to store the data in Bubbly. It is not strictly SQL but follows a similar approach of defining tables with relations. The schema is defined using HCL and uploaded to the Bubbly server or provided as a config during startup.
4. The **Importers** define the expected *input data* and how to convert that into a data structure defined by the *schema* so that Bubbly can store understand the *input data*.

All of the configurations, schemas, importers and such are defined as HCL, which provides the necessary high level abstraction to make configuring Bubbly easy, but still provides the necessary flexibility to process lots of different types of data.

The above architecture can be further elaborated using the following sequence diagram:

![Bubbly High Level Sequence Diagram](./images/high-level-sequence-diagram.drawio.svg)

## 4. Design

This section describes some design aspects of bubbly.

## 4.1 Data Model and Schema

To help describe the design of Bubbly we will use the data model illustrated below, which is a subset of the overall model for modelling testing results, such as those from an automated test suite.
Such an example could be adapted to suit any type of test result data, whether it be static analysis of source code, or scanning for 3rd party components, or issues related to infrastructure.

![Bubbly Example Data Model](./images/example-data-model-testing.drawio.svg)

In this example we have a Product, which can have many Projects.
Repo can belong to many Projects, and they have a version (e.g. Git commit) through RepoVersion.
The model for the automated test results is modelled by a TestRun associated with a RepoVersion, which has TestSet which include TestCases.
This simple data model would allow us to store automated test results from a test run that is associated with a verison of source code in a repository.

This data model should be defined as a **schema** using HCL, and the following is an example of such a schema:

```hcl
// flat list
table "product" {
    field "name" { /* ... */ }
}

table "project" {
    field "name" { /* ... */ }
    // relations using something like a foreign key
    field "product_ids" {}
}

// nested relations
table "repo" {
    field "name" { /* ... */ }
    field "project_ids" {}

    table "repo_version" {
        field "name" { /* ... */ }
        field "version" { /* ... */ }
        // repo_id is automatically added as repo_verison is nested under repo
        // field "repo_id" { /* ... */ }
    }
}

table "test_run" {
    field "name" {}
    field "repo_version_id" {}

    table "test_set" {
        field "name" {}

        table "test_case" {
            field "name" {}
            field "status" {}
            field "test_set_id" {}
        }
    }
}
```

This schema should be stored in the Bubbly server so that it can be fetched by the Bubbly client.
It can either be uploaded after the Bubbly server starts, or can be provided as a configuration during startup.

### 4.2 Importers

The purpose of Importers is to import and convert input data, such as report files produced by different tools.

If we continue on the example of test automation, a common result format for automated tests is that produced by the different [xUnit](https://en.wikipedia.org/wiki/XUnit) tools, such as [JUnit](https://junit.org/junit5/), which is an XML file syntax looking something like the following:

```xml
<!-- SOURCE: https://gist.github.com/n1k0/4332371 -->
<?xml version="1.0" encoding="UTF-8"?>
<testsuites duration="50.5">
    <testsuite failures="0" name="Untitled suite in /Users/niko/Sites/casperjs/tests/suites/casper/agent.js" package="tests/suites/casper/agent" tests="3" time="0.256">
        <testcase classname="tests/suites/casper/agent" name="Default user agent matches /CasperJS/" time="0.103"/>
        <testcase classname="tests/suites/casper/agent" name="Default user agent matches /plop/" time="0.146"/>
        <testcase classname="tests/suites/casper/agent" name="Default user agent matches /plop/" time="0.007"/>
    </testsuite>
</testsuites>
```

We could produce an importer using the following HCL specification.

```hcl
importer "xunit_report" {
    type = "xml"
    format = object({
        testsuites: object({
            duration: number,
            testsuite: list(object({
                failures: number,
                name: string,
                package: string,
                tests: number,
                time: number
            }))
        })
    })

    // this is part of bubbly's magic sauce and needs to be created...
    // basically prepare the data for use to follow the schema
    convert {
        // do something like setproduct to make the data ready to process in HCL
        // https://www.terraform.io/docs/configuration/functions/setproduct.html
        // idea is to do as much heavy liftin on the import and make the
        // consumption of the import as easy as possible, and keep things DRY
    }
}
```

### 4.3 Client Configs

The client configs tell the Bubbly client what data to upload to the Bubbly server.

First we would need to provide some of the contextual data, such as the project name, the repository name, the version of the codem and generic things which having nothing to do with parsing data, but will feed into the schema and data model.

In the following short example we set the name of the `repo` to `test-repo`.
The idea is if this `repo` already exists on the bubbly server then it would associate the data we are about to produce with that existing `repo`.
Next we create a `repo_verion`, which may exist if we have already uploaded data against the same `repo_version`.

```hcl
// retrieve the repo called "test-repo"
data "repo" "test_repo" {
    name = "test-repo"
}

// create a repo_version
data "repo_version" "test_version" {
    version = env.GIT_COMMIT
    source_repo_id = data.source_repo.test_repo.id
    metadata {
        branch = env.GIT_BRANCH
    }
}
```

Now that we have a `repo_version` to associate our `test_run` with, we can create the `test_run` together with its associated instances of `test_suite` and `test_case`.
We simply say that we want to use the importer called `xunit_report` and tell it where to get the `input_data` from, which is a glob expression for the XML files under the directory `xunit`, in this example.

```hcl
// create a test_run called test-run-XYZ
data "test_run" "test_run" {
    name = "test-run-${env.JOB_NUMBER}"
    // specify the importer to use
    importer = "xunit_report"
    // specify the input data
    input_data {
        type = "xml"
        source = "./xunit/*.xml"
    }
}
```

### 4.4 Queries

Once a suitable data model has been created and data has been populated using the client configs and importers, it is time to make use of this data.

Queries are defined as HCL and can be sent to the Bubbly server and the Bubbly server will process the query and return the relevant data.

THIS IS A SIGNIFICANT #TODO
