---
title: Resource Kinds
hide_title: false
hide_table_of_contents: false
keywords:
- docs
- bubbly
- resources
- kinds
---

:::note Under Active Development
The content for this page is *under active development*. It
may be
incorrect and/or
incomplete and is liable to change at any time.
:::

## `extract`

`extract` resources allow you to extract data from REST, GraphQL, JSON or XML sources.

### Example Usage

#### Github Repository Fork Count

```hcl
resource "extract" "github_fork_count_extract" {
	spec {
		type = "graphql"

		source {
			url = "https://api.github.com/graphql"
			
			query = <<EOT
				query { 
					repository(owner:"docker", name:"compose") {
						owner {
							login
						}
						name
						forkCount
					}
				}
			EOT

			bearer_token = env("GH_TOKEN")

			format = object({
				repository: object({
					owner: object({
						login: string
					}),
					name: string,
					forkCount: number,
				})
			})
		}
	}
}
```

### Specification Reference

An `extract` resource has the following specification:

- `type`: The source type. Options are `rest`, `graphql`, `json`, `xml`, `git` 
- `source`: A single configuration block containing configuration of the chosen source type

#### `graphql` Source

The following attributes and blocks are supported:

- `basic_auth`: (Optional) Configuration block setting the `Authorization` header. It supports the following:
  - `username`: (Optional) Username
  - `password`: (Optional) Password
  - `password_file` (Optional) Path to file containing a password
- `bearer_token`: (Optional) the bearer token to be embedded into the `Authorization` header
- `bearer_token_file`: (Optional) Path to file containing the bearer token to be embedded into the `Authorization` header
- `format`: The format of the response to the GraphQL API query, defined as a `cty.Type`
  :::note
  The content for the `format` attribute is *under active development* and will be
  expanded upon soon.
  :::
- `headers`: (Optional) HTTP headers to set
- `method`: (Optional) Method of data extraction. Options: `GET`, `POST`. Default: `POST`
- `params`: (Optional) URL query parameters which get prepended to the end of the URL
- `query`: the GraphQL query string, wrapped between `<<EOT` and `EOT` per heredoc syntax.
- `timeout`: (Optional) How long (in seconds) the extractor waits before giving up trying to extract data from the given source. Default: 1 second
- `url`: URI of the GraphQL endpoint

#### `rest` Source

The following attributes and blocks are supported:

- `basic_auth`: (Optional) Configuration block setting the `Authorization` header. It supports the following:
  - `username`: (Optional) Username
  - `password`: (Optional) Password
  - `password_file` (Optional) Path to file containing a password
- `bearer_token`: (Optional) the bearer token to be embedded into the `Authorization` header
- `bearer_token_file`: (Optional) Path to file containing the bearer token to be embedded into the `Authorization` header
- `flavour`: (Optional) The expected format of the response. Options: `json`. Default: `json`
- `format`: The format of the response to the REST request, defined as a cty.Type
  :::note
  The content for the `format` attribute is *under active development* and will be
  expanded upon soon.
  :::

- `headers`: (Optional) HTTP headers to set
- `method`: (Optional) Method of data extraction. Options: `GET`, `POST`. Default: `POST`
- `params`: (Optional) URL query parameters which get prepended to the end of the URL
- `timeout`: (Optional) How long (in seconds) the extractor waits before giving up trying to extract data from the given source. Default: 1 second
- `url`: URI of the GraphQL endpoint

#### `git` Source

The following attributes and blocks are supported:

- `directory`: Path to the directory containing the local git repository

#### `json` Source

The following attributes and blocks are supported:

- `file`: Path to the JSON file
- `format`: The format of the raw input data, defined as a cty.Type
  :::note
  The content for the `format` attribute is *under active development* and will be
  expanded upon soon.
  :::

#### `xml` Source

The following attributes and blocks are supported:

- `file`: Path to the XML file
- `format`: The format of the raw input data, defined as a cty.Type
  :::note
  The content for the `format` attribute is *under active development* and will be
  expanded upon soon.
  :::


## `transform`

`transform` resources allow you to transform data extracted from `extract` 
resources into a format compliant with your underlying Bubbly Schema.


### Example Usage

#### Github Repository Fork Count

Imagine you have defined a Bubbly Schema with at least the following `repo_stats` table:
```hcl
table "repo_stats" {
    field "owner" {
        type = string
    }
    field "repo" {
        type = string
    }
    field "fork_count" {
        type = number
    }
}
```

We now need a `transform` resource to transform the output from the `github_fork_count_extract`
[example](#extract-github-repository-fork-count) in order to match this schema definition.

Using the Bubbly Language this is as simple as:

```hcl
resource "transform" "github_fork_count_transform" {
	spec {
		input "data" {}

		data "repo_stats" {
			fields = {
				"owner":               self.input.data.repository.owner.login
				"repo":                self.input.data.repository.name
				"fork_count":          self.input.data.repository.forkCount
			}
		}
        # Last few open issues
        dynamic "data" {
          labels = ["issues"]
  
          for_each = self.input.data.repository.open_issues.edges
          iterator = it
  
          content {
            joins = ["repo_stats"]
            fields = {
              "url":    it.value.node.url
              "title":  it.value.node.title
              "state":  it.value.node.state
            }
          }
        }
	}
}
```

### Specification Reference

The following attributes and blocks are supported:

- `input "<BLOCK LABEL>"`: (Optional) one or more configuration block defining inputs to the `resource`. 
  The label represents the locally-scoped name of the input.
  - `description`: (Optional) Description of the input
  - `default`: (Optional) The default value of this input, defined as a `cty.Value`
  - `type`: (Optional) The type of this input, defined as a `cty.Type`
- `data "<BLOCK LABEL>"`: (Optional) one or more configuration block defining the mapping of
input to output data. The block label represents the mapping between the `table` block of the Bubbly Schema
and the transformed data.
  - `fields`: A map of fields to assign input data to.
- `dynamic "<BLOCK LABEL>"`: (Optional) Zero or more repeatable configuration blocks defining iteration over input data. 
  These blocks produce nested `data` blocks on each iteration. The label of the dynamic block specifies what kind of nested
  block to generate.
  - `for_each`: The value to iterate over
  - `iterator`: (Optional) Sets the name of a temporary variable that represents the current 
    element of the value being iterated over
  - `labels`: (Optional) A list of of strings that specifies the block labels,
    in order, to use for each generated block. It is possible to use the temporary 
    iterator variable defined by the `iterator` attribute in this value.
  - `content`: A configuration block defining the body of each generated block. 
    It is possible to use the temporary iterator variable defined by the `iterator` 
    attribute in this value.
  
  :::note 
  The specification of `dynamic` blocks matches that of the Terraform Language, 
  and is therefore well defined in the 
  [Terraform documentation](https://www.terraform.io/docs/language/expressions/dynamic-blocks.html).
  :::

## `load`

`load` resources post the data transformed by `transform` resources to Bubbly.

### Example Usage

#### Github Repository Fork Count

Given a collection of data that has been transformed to satisfy your Bubbly Schema,
we define a `load` resource to _load_ this data to Bubbly like so:

```hcl
resource "load" "github_fork_count_load" {
	spec {
		input "data" {}
		data = self.input.data
	}
}
```

### Specification Reference

The following attributes and blocks are supported:

- `input "<BLOCK LABEL>"`: (Optional) one or more configuration block defining inputs to the `resource`.
  The label represents the locally-scoped name of the input.
  - `description`: (Optional) Description of the input
  - `default`: (Optional) The default value of this input, defined as a `cty.Value`
  - `type`: (Optional) The type of this input, defined as a `cty.Type`
- `data`: An explicit mapping of input data to data to be loaded to Bubbly

## `query`

`query` resources query Bubbly for some subset of data stored in the Bubbly Store.

### Example Usage

#### Github Repository Fork Count

Given a `run` resource has been run, loading Github repository fork count data
to the Bubbly Store, we can write the following `query` resource to query the
Bubbly Store for all repositories and their respective fork counts:

```hcl
resource "query" "github_fork_count_query" {
  spec {
    query = <<EOT
      {
        repo_stats {
          repo
          fork_count
        }
      }
    EOT
  }
}
```

### Specification Reference

- `query`: the GraphQL query string, wrapped between `<<EOT` and `EOT` per heredoc syntax.

## `criteria`

`criteria` resources represent a list of queries, conditions and logical operators 
between these conditions. At runtime, these conditions should all be evaluated to 
return a boolean result. In short, `criteria` resources let you set up gated 
checks against your data to make sure that certain conditions are (always) met.

### Example Usage

:::note Under Active Development
This content is *under active development*. It
may be
incorrect and/or
incomplete and is liable to change at any time.
:::


### Specification Reference

- `query "<BLOCK LABEL>"`: one or more configuration block defining the `query` 
  resources relevant to this `criteria`. The block label represents the name of the
  `query` resource.
- `condition "<BLOCK NAME>"`: one or more configuration block defining the 
  conditions that should be evaluated by this `criteria` resource. The block label represents
  the locally-scoped name of the `condition`, which can be used by the `operation`
- `operation "<BLOCK NAME>"`: A single configuration block defining the logical
  operator(s) to be performed over the value of the `criteria`'s one or 
  more `condition` resources.
  

## `pipeline`

`pipeline` resources encapsulate many resources together with some sequential run order,
taking the output of one resource run and using it as the input for the next.

`pipeline` resources are comprised of a logical grouping of inputs and tasks:
- `tasks` define specific Bubbly resource to be _applied_ and _run_
- `inputs` define content required to run the resources referenced by the `tasks`.
  For example:
    - the directory of a junit test file required by an `extract` resource
    referenced by one of the pipeline's `tasks`
    - the output of another `task` in the pipeline, referenced by `self.task.<TASK NAME>.value`.
    See [`task` behaviour](./kinds#task-behaviour) for more information.
      
### Example Usage

#### Github Repository Fork Count

Given a definition of an `extract` resource which extract Github repository fork count, 
a `transform` resource which transforms that data to comply to your Bubbly Schema and
a `load` resource which posts that data to Bubbly, we define a `pipeline` resource
to define the process of `extract -> transform -> load` as follows:

```hcl
resource "pipeline" "github_fork_count_pipeline" {
	spec {
		task "extract" {
			resource = "extract/`github_fork_count_extract`"
		}
		task "transform" {
			resource = "transform/github_fork_count_transform"
			input "data" {
				value = self.task.extract.value
			}
		}
		task "load" {
			resource = "load/github_fork_count_load"
			input "data" {
				value = self.task.transform.value
			}
		}
	}
}
```

### Specification Reference

The following attributes and blocks are supported:

- `input "<BLOCK LABEL>"`: (Optional) one or more configuration block defining inputs to the `resource`.
  The label represents the locally-scoped name of the input. Within this block, 
  the following attributes are supported:
  - `description`: (Optional) Description of the input
  - `default`: (Optional) The default value of this input, defined as a `cty.Value`
  - `type`: (Optional) The type of this input, defined as a `cty.Type`
- `task "<BLOCK LABEL>"`: One or more configuration block defining a task to be performed 
  by the pipeline at runtime. The label of the block represents the locally-scoped name of
  the task. Within this block, the following attributes and blocks are supported:
    - `resource`: The ID of the resource to be run at task runtime. This ID must
      be of the form `"<kind>/<name>"`.
    - `input "<BLOCK LABEL>"`: (Optional) One or more configuration block defining inputs 
      to the `task`, and by extension to the `resource` referenced by the `task`.
      The label represents the name of the input, and must match the name used by
      the underlying `resource`'s definition. Within this block, the following 
      attributes are supported:
        - `value`: The value of this input, defined as a `cty.Value`

### `task` behaviour

Each task in a `pipeline` has an output available to other `tasks` of the same `pipeline` 
resource, referenced by `<TASK NAME>.value`. As a result, the output from one `task` can
be _piped_ as input to another `task`!

Take another look at the [`github_fork_count_pipeline`](./kinds#github-repository-fork-count-4).
Verbosely, this resource states:
1. Extract the data from Github by _running_ the `extract` resource of name `github_fork_count_extract`
2. Using the output from the `github_fork_count_extract` resource as input, 
   transform the data by _running_ the `transform` resource of name `github_fork_count_transform`
3. Using the output from the `github_fork_count_transform` resource as input,
   load the data to Bubbly by _running_ the `load` resource of name `github_fork_count_load`

## `run`

`run` resources are responsible for _running_ referenced resources. 
After being _applied_ to Bubbly, `run` resource are _immediately run_.

### Example Usage

#### Github Repository Fork Count

Given a `pipeline` resource capable of extracting, transforming and loading Github
repository fork count data, we define a `run` resource to run the `pipeline` resource as follows:

```hcl
resource "run" "github_fork_count_run" {
	spec {
		resource = "pipeline/github_fork_count_pipeline"
	}
}
```

After being _applied_ to Bubbly, the `run/github_fork_count_run` resource will be _run_.
At runtime, all resources referenced by the underlying
pipeline's `tasks` will also be run.

### Specification Reference

- `input "<BLOCK LABEL>"`: (Optional) one or more configuration block defining inputs to the `resource`.
  The label represents the locally-scoped name of the input.
  - `description`: (Optional) Description of the input`
  - `default`: (Optional) The default value of this input, defined as a `cty.Value`
  - `type`: (Optional) The type of this input, defined as a `cty.Type`
- `resource`: The ID of the resource to be _run_ by this `run` resource
- `remote`: (Optional) A single configuration block enabling and defining the
  configuration of _remote running_ of this `run` resource. Use this when you do
  not want to run a resource locally, but instead want the Bubbly Worker to run
  the resource remotely.
  - `interval`: the interval at which to remotely run the `run` resource. 
    This string should be of the format `"XhYmZs"`, such as `"72h3m0.5s"`

  :::note 
  See [`run` resource behaviour](../resources/kinds#run-resource-behaviour-differences-to-other-resources)
  for a more information on the effects this block has on the behaviour of `run` resources
  :::

### `run` Resource Behaviour: Local vs Remote

`run` resources are, by design, flexible in their runtime location and behaviour. 
There are currently three different behaviours:

### 1. Local `run`

The default behaviour for a resource of `run` kind. At time of apply, a local run
be sent to Bubbly to be saved and will then be run on the _same local machine_ running
the `bubbly` executable.

#### Example

```hcl
resource "run" "github_fork_count_run" {
	spec {
		resource = "pipeline/github_fork_count_pipeline"
	}
}
```

### 2. Remote `run`

A `run` resource is _remote_ when its specification includes a `remote{}` block.
At time of apply, the remote run be sent to Bubbly to be saved and is then picked
up by a Bubbly Worker to be run on whichever machine this Worker is running.

While Remote `run` resources are _always_ run by the Bubbly Worker, they have
differing behaviour depending on the configuration of this `remote{}` block.

#### 2.1 _One-Off_ Remote `run`

If the `remote{}` block is empty, the `run` is _one-off_ by default.
The Bubbly Worker will run this resource once,
load the results of the run to the Bubbly Store, then end.

#### Example

```hcl
resource "run" "github_fork_count_run" {
	spec {
		resource = "pipeline/github_fork_count_pipeline"
	}
    remote{}
}
```

#### 2.2. _Continuous_ Remote `run`

:::caution Behaviour Disabled
The _continuous_ remote `run` type has been disabled for stability purposes while
work to expand the Bubbly Worker's functionality is undergone. All _continuous_
`run`s will default to _one-off_ runs in the meantime.
:::

If the `remote{}` block specifies the optional `interval` attribute, the `run` is
_continuous_. The Bubbly Worker will run this resource continuously at this
interval, loading the results of each interval run to the Bubbly Store

#### Example

```hcl
resource "run" "github_fork_count_run" {
	spec {
		resource = "pipeline/github_fork_count_pipeline"
	}
    remote{
      interval = "300s"
    }
}
```

### `run` Resource Behaviour: Differences to Other Resources

Under the hood, Bubbly manages `run` resource kinds slightly differently to other
resources. While at apply time all resources are sent to Bubbly and saved, `run`
resources are also then _immediately run_. 

This nuanced different is key: you can consider `run` resources the actual executors of
your resources. Without them, your resources would be saved in the Bubbly Store but never
actually be run!

Conceptually, this is similar to the behaviour employed by [TektonCD](https://github.com/tektoncd/pipeline/tree/master/docs#tekton-pipelines-entities)
and its `TaskRun` and `PipelineRun` entities, which are themselves 
the executors of TektonCD's `Task` and `Pipeline` entities respectively.

### `run` Resource Behaviour: Using Bubbly's `/api/v1/run/:name` endpoint

:::note Under Active Development
This functionality is *under active development*. It may be incorrect and/or
incomplete and is liable to change at any time.
:::

In some environments, it may be preferential to export the data produced within
your software processes to Bubbly for remote processing.
Bubbly's `/api/v1/run/:name` endpoint supports this use-case, providing the 
ability
to trigger the _run_ of a named `run` resource via a 
[multipart formpost](https://everything.curl.dev/http/multipart) 
(`POST` request with the body formatted with data as a series of parts).

With this, it is possible to embed the input data to the named `run` resource
within the `POST` request.

#### Example

Consider the following example `run` and `extract` resources:

```hcl
resource "extract" "sonarqube" {
  metadata {
    labels = {
      "environment": "prod"
    }
  }
  spec {
    input "file" {}
    type = "json"
    source {
      file = self.input.file
      #format = ...
    }
  }
}

resource "run" "sonarqube_remote" {
  spec {
    remote {}

    resource = "extract/sonarqube"
    input "file" {
      value = "./testdata/sonarqube/sonarqube-example.json"
    }
  }
}
```

When _run_, the `sonarqube_remote` resource will attempt to run the
`extract/sonarqube` resource, which parses the
`./testdata/sonarqube/sonarqube-example.json` and extracts the data according
to its defined `format` field.

`sonarqube_remote`'s ``remote {}` block signifies that this `run` resource should
be run _remotely_ by the Bubbly Worker. The challenge is that its `file` input 
requires that the Bubbly Worker have available the 
`./testdata/sonarqube/sonarqube-example.json` file _locally_. The 
`/api/v1/run/:name`
endpoint allows us to send this file, along with a run request, to Bubbly, via the
following `cURL` commands:

```shell
$ curl \
   -F "file=@</path/to/file.json>;type=application/json;filename=</path/to/file.json" \
   http://<HOST>:<PORT>/api/v1/run/<RUN-NAME>
```

In this case, the following applies:
```shell
$ curl \
  -F "file=@./testdata/sonarqube/sonarqube-example.json;type=application/json;filename=./testdata/sonarqube/sonarqube-example.json" \
  http://<HOST>:<PORT>/api/v1/run/sonarqube_remote
```

In this case, the Bubbly Worker will save the contents to the relative directory
`./testdata/sonarqube/sonarqube-example.json`.

Similarly, the `.json` file can be compressed and sent to Bubbly as a `.zip` file via:

```shell
$ curl \
  -F "file=@./testdata/sonarqube/sonarqube-example.zip;type=application/zip;filename=./testdata/sonarqube/sonarqube-example.zip" \
  http://<HOST>:<PORT>/api/v1/run/sonarqube_remote
```

The Bubbly Worker will uncompress the contents to the relative directory `./testdata/sonarqube/<contents>`. 
Paths defined internally within the `.zip` file are also maintained.

#### Conventions

:::caution

There are several conventions that must be followed in order for Bubbly to handle the
request correctly

:::

To keep this feature well-scoped, the Bubbly API Server
enforces some conventions for the expected format and declaration of this input data:

1. Currently, only one file upload is supported, and must be sourced using the `file` field name.
   If multiple files are required, upload a `.zip` file
2. The `type` field name is required; Bubbly uses this to determine the underlying content type of the file
3. In cases where the path value specified by the `run`'s `input{}` is non-root, the `filename`
   field is _required_, and _must_ match the path specified by the `run` resource's `input`; Bubbly uses this
   to determine the relative path to use when saving the file to the Worker's local filesystem

#### Limitations

Currently, only `application/json` and `application/zip` content types are supported.
