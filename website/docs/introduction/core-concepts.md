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

## Hashicorp Configuration Language (HCL)

[HCL](https://github.com/hashicorp/hcl) is a toolkit for creating
configuration languages which are human- and machine-friendly.  At its core, the HCL [syntax](https://github.com/hashicorp/hcl#information-model-and-syntax) is made up of _attributes_ and _blocks_, 
and utilisers of HCL offer application-specific _block_ types whose _labels_ represent different logical entities of the application.

### Example: Terraform 

Perhaps the most widely-known utiliser of HCL, [Terraform](https://www.terraform.io/), 
leverages HCL to offer the [Terraform Language](https://www.terraform.io/docs/language/index.html#about-the-terraform-language), 
which provides a declarative definition of infrastructure objects, such as an AWS VPC:

```hcl
resource "aws_vpc" "main" {
  cidr_block = var.base_cidr_block
}
```

- _block_ `resource`: declares this as _some_ infrastructure object 
- _labels_ `aws_vpc` and `main`: uniquely identify the underlying type (`aws_vpc`) and name (`main`) of the infrastructure object
- _body_ (between `{` and `}`): attributes and nested _blocks_ representing the VPC's configuration


### Example: Bubbly 

Similarly, Bubbly leverages HCL to offer the Bubbly Language, a configuration language for 
defining [Bubbly Resources](../resources/overview) and [Bubbly Schemas](../schema/schema) in `.bubbly` files. Just like the Terraform Language 
provides the `resource` block as an abstraction for infrastructure objects, the 
Bubbly Language provides the `resource` block as an abstraction for data pipeline objects
and the `table` block as an abstraction for database tables. 

For example, an extraction of Github repository data from its GraphQL API:

```hcl
resource "extract" "github_extract" {
    api_version = "v1"
    
    spec { 
        type = "graphql"
        source { 
            url = "https://api.github.com/graphql"
            #...	
        }
    }
}
```

- _block_ `resource`: declares this as _some_ data pipeline object
- _labels_ `extract` and `github_extract`: uniquely identify the underlying kind 
  (`extract`) and name (`github_extract`) of the object
- _body_ (between `{` and `}`): attributes and nested _blocks_ representing the 
  object's configuration (in this case, defining how to interact with the Github API)

## Bubbly Schema

The Bubbly Schema represents the desired database structure of the data that Bubbly collects. 
It is expected that as the knowledge and understanding of this data structure is refined, 
the schema will change, and Bubbly has been written to evolve with these changes.

See [Bubbly Schema](../schema/schema) for more information.

## Resources

Bubbly resources are the declarative, highly-customizable building blocks of Bubbly. 
Whether you want to extract data from an external API, load data to Bubbly or query Bubbly for data, 
Bubbly resources are your interface for doing so.

A Resource's lifecycle has 2 parts:

1. Apply: The resource is sent to Bubbly and saved.
2. Run: The action defined by the resource is performed.
    1. A resource is run when:
        1. _applying_ a resource of kind `run`. 
            1. See [`run` resource behaviour](../resources/kinds#run-resource-behaviour-differences-to-other-resources)
           for more information on this kind of resource run.
        2. the resource is referenced by a `task` within a running `pipeline` resource. 
            1. See the [`pipeline` resource](../resources/kinds#pipeline) for more 
               information on `tasks`.
            

See [Resources](../resources/overview) for more information.

## Data Pipelines

Bubbly provides a common framework for data integration 
from your various tools and processes. While this framework is novel, 
the concept is not, being commonly referred to as a _data pipeline_ in data engineering.

A typical data pipeline has three main stages, and so as not to go against the grain 
Bubbly [resource kinds](../resources/kinds) are named according to these same stages:

1. Extraction: pulling data from one or more data sources. 
   Provided by the [`extract`](../resources/kinds#extract) resource kind.
2. Transformation - manipulating the data prior to storage. 
   Provided by the [`transform`](../resources/kinds#transform) resource kind.
3. Loading - uploading the data into the target data store. 
   Provided by the [`load`](../resources/kinds#load) resource kind.

Finally, the [pipeline](../resources/kinds#pipeline) resource kind provides a 
tidy way to connect these (and other) resources together sequentially, 
taking the output of one resource run and using it as the input for the next.

## GraphQL

[GraphQL](https://graphql.org/) is a query language for APIs. 
For a given [Bubbly Schema](#bubbly-schema), Bubbly  automatically generates 
a GraphQL interface exposing all data that has been stored. 
This allows you to write GraphQL queries against data loaded into the Bubbly Store, 
either by providing a GraphQL query directly through the `/api/v1/query` endpoint or 
through the [`query`](../resources/kinds#query) resource kind.

## Agent

Bubbly is distributed as a single executable, which provides client- and server-side 
functionality through varying commands. Users of Hashicorp's [Nomad](https://www.nomadproject.io/)
or [Consul](https://www.consul.io/) may be familiar with this architecture pattern.

The Bubbly agent is the core process of Bubbly. It has a set of optional features 
that it can serve in a plug-and-play style architecture:

- API Server: the user-facing REST API server that exposes all the internal bubbly services (such as the GraphQL interface)
- Store: the Bubbly store handles all the database transactions as well as serving the internal GraphQL interface
- [Worker](#worker): the Bubbly worker is a service that runs resources on the server-side
- NATS Server: Bubbly comes with a built-in NATS server that can be used for simple deployments. 
  It should be replaced with a dedicated NATS deployment for more demanding environments

## Worker

The Bubbly Worker is a feature of the Bubbly agent, but whose responsibilities deserve their own dedicated mention.

The Worker’s purpose is to run _remote_ resources of `run` kind.
A `run` resource is _remote_ when its specification includes a `remote{}` block. 

Remote `run` resources can be run in two ways:
1. A _one-off_ run. A remote run is _one-off_ by default. 
   The Bubbly Worker will run this resource once, 
   load the results of the run to the Bubbly Store, then end. 
2. A _continuous_ run. A remote run is _continuous_ when it has a specified 
   interval period (e.g. minutes/hours). The Bubbly Worker will run this resource
   continuously at this interval, loading the results of each interval run
   to the Bubbly Store

See [Resource Kinds](../resources/kinds#run) for more information about the `run` resource kind.

## UI and Data Visualisation

Originally, Bubbly was envisioned to leverage existing 
dashboard solutions like Grafana and Kibana for its data visualisation. However, these tools do not provide 
clear support for dynamic navigation of hierarchical data. Since much of the data
in Bubbly is hierarchical (according to the Bubbly Schema), we saw the need for a
purpose-built UI to solve this use-case.

Therefore, there are multiple options for Bubbly data visualisation: 

1. Existing dashboard solutions like Grafana, Kibana, which are fully supported by Bubbly's external GraphQL interface
    - these solutions are strongly encouraged in cases where hierarchical navigation is not needed
2. Bubbly's purpose-built web-based UI
    -  designed specifically with the hierarchy navigation use-case in mind

The Bubbly UI is still in its early stages of development and is not part of the `bubbly`
executable.

## External Dependencies

Bubbly requires an external database for the Bubbly Store to connect to. Currently supported are:

- Postgres

There are plans to support more databases in the future, and not just SQL databases but also stores like ElasticSearch.
