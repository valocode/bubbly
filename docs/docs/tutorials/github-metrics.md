---
title: Github Repository Metrics
hide_title: false
hide_table_of_contents: false
keywords:
- docs
- bubbly
- tutorials
- guides
---

:::note Under Active Development
The content for this page is *under active development*. It
may be
incorrect and/or
incomplete and is liable to change at any time.
:::

## Infrastructure

To run the tutorial, you'll need the [Bubbly Agent](../introduction/core-concepts#agent) running with its Store connected to your external Postgresql database, as described in the [Getting Started](../getting-started/getting-started.md) section.

## Introductions

The Control Room of _BubbliCorp_ &ndash; our visionary, even if imaginary, start-up &ndash; is bustling with activity. Not all of it is well-planned and some of the better-planned is less-than-better executed, but it's very much alive and genuine, and everyone is doing their best to finish tasks from the project's first Milestone on time.

Speaking of tasks, there are some marked as `HIGH PRIORITY` by _BubbliCorp_'s CFO, Mr. H.G.Blob, whom seasoned _BubbliCorp_'ers have advised you "not to anger under any circumstances". Perhaps, if you built something valuable for the mighty Mr Blob, you might get noticed... they might even give you a free lunch voucher to compliment your unpaid lunch break!

## The task

Since you are new to Bubbly, you pick a task which appears to be self-contained and mentions some familiar names, such as GitHub:

> Me likes to know how our closest competitors are doign! Especially Docker! Me wants know GitHub numbers! Numbers mean money! Money good! Now get bak to work! ~ Signed: H.G.Blob

## Exploration

After having done some research on this topic, you conclude that the [GitHub GraphQL API][gh-graphql-main] can provide the information you need. To communicate with the GraphQL server, you'll need an OAuth token with the right scopes, so you [create a personal access token][gh-token] and run a few queries using both [GraphiQL][gh-graphiql] and [Insomnia][insomnia] apps. You cannot decide yet which one wins your heart, but you get a good feel for how the API works.

You settle on the following query.

[gh-graphql-main]: https://docs.github.com/en/graphql
[gh-token]: https://docs.github.com/en/github/authenticating-to-github/creating-a-personal-access-token
[gh-graphiql]: https://docs.github.com/en/graphql/guides/using-the-explorer
[insomnia]: https://insomnia.rest/

```graphql
query { 
	repository(owner:"docker", name:"compose") {

		owner {
			login
		}
		name

		forkCount
		stargazerCount
		watchers {
			totalCount
		}

		releases(last:1) {
			totalCount
			# TODO:
			#   You also would like to know the total download count
			#   for the latest release, but there is no field for that,
			#   so you note that it can be done by iterating through
			#   `releaseAssets`, and you leave that for the next step.
		}
	}
}
```

You note:

* Some fields are nested, some are not: `forkCount` vs `watchers/totalCount`;
* There is no straightforward way to get the total count of downloads for the _latest_ release, although you have a strong feeling you can solve this later, so you just focus on available fields for now.

## Thinking in Bubbly

Having skim-read the [Core Concepts][core] section, you begin to sketch out the solution using the vocabulary provided by Bubbly Resources. You have a vague idea that you'd like to build a [`pipeline`][tib-p], consisting of [`extract`][tib-e], [`transform`][tib-t], and [`load`][tib-l] stages. And then to [`run`][tib-r] it.

For that to work, you'll also need to define how to store the data by defining and applying a [Bubbly Schema](../schema/schema).

[tib-e]: ../resources/kinds.md#extract
[tib-t]: ../resources/kinds.md#transform
[tib-l]: ../resources/kinds.md#load
[tib-p]: ../resources/kinds.md#pipeline
[tib-r]: ../resources/kinds.md#run

## Extracting the data

It appears that a query to a GraphQL end-point is best described by the `extract` resource. Having played with the [GitHub GraphQL API][graphql-ref] in the previous sections, you already know the query that you'd like to run.

[graphql-ref]: https://docs.github.com/en/graphql/reference/objects
[core]: ../introduction/core-concepts.md

Armed with this knowledge, you write your first definition of a Bubbly resource, and save it in a file with the flamboyant name `tutorial1.bubbly`:

```hcl
#
# tutorial1.bubbly
#

resource "extract" "repo_stats" {
	spec {
		type = "graphql"

		source {
			url = "https://api.github.com/graphql"

			bearer_token = env("GH_TOKEN")

			query = <<-EOT
				query { 
					repository(owner:"docker", name:"compose") {

						owner {
							login
						}
						name

						forkCount
						stargazerCount
						watchers {
							totalCount
						}

						releases(last:1) {
							totalCount
						}
					}
				}
			EOT

			format = object({
				repository: object({
					
					owner: object({
						login: string
					}),

					name: string,
					forkCount: number,
					stargazerCount: number,

					watchers: object({
						totalCount: number
					}),

					releases: object({
						totalCount: number
					})
				})
			})
		}
	}
}
```

This Bubbly file is written in standard HCL, using the notation understood by the Bubbly's `extract` resource. The `spec` block 

:::caution TODO

This is standard HCL, spec block for different types of extract, each extract type has different options , bearer token in env variable, graphql query as HCL heredoc, format is HCL type expression (link to detailed guide), mention objects and strings, numbers, {} is map syntax because even format = ... is valid HCL.

../resources/kinds.md#extract provides the spec for the different source types. We definitely need to provide a detailed guide for the format attribute. I suggest adding an admonition stating this is under construction and a full explanation for the format attribute used in extract/repo_stats for this tutorial is coming soon.

:::

## Bubbly Schema

Now that you know the format of your input, you feel ready to formalise the storage requirements. To do so, you define a [Bubbly Schema](../schema/schema.md). A Bubbly Schema is conceptually close to a database schema, in which you describe the data blocks and the relations between them just as they are stored in the database. Your imagination goes wild with possibilities and you pick the name `tutorial1.schema` for the file, containing your latest masterpiece:

```hcl
#
# tutorial1.schema
#

table "repo_stats" {

	field "id" {
		type = string
		unique = true
	}

    field "owner" {
        type = string
    }
    field "repo" {
        type = string
    }

    field "forks" {
        type = number
    }
    field "stargazers" {
        type = number
    }
    field "watchers" {
        type = number
    }

    field "releases" {
        type = number
    }

	# TODO there is more information that you'd like
	#      to save, such as the total download count,
	#      but you leave that for another day...
}
```

This Schema describes the creation of a single table, `repo_stats`, with columns `owner`, `repo`, `forks`, `stargazers`, `watchers` and `releases` mapping directly onto the data from your GraphQL query. Further, the `id` column will store a concatenation of `owner` and `repo`, providing a natural ID for each entry and allowing you to extend your `extract` in the future to monitor the stats for multiple repositories. 

You've picked names that are, perhaps, more convenient than the names given to you in the GraphQL output, and you are already having ideas for more synthetic fields, such as the total number of downloads for the latest release, which is something you'll have Bubbly compute for you later. For now, this example suffices.

You save the Bubbly Schema to a `tutorial1.schema` file and proceed to _apply_ it to the Bubbly _agent_, running on `--host 127.0.0.1` and `--port 8111` (default settings) by running:

```shell
$ bubbly schema apply -f ./tutorial1.schema
```

By default, the agent is configured to use a Postgres database as its _store_.

The agent responds with a confirmation and your heart sings.

```shell
schema at path "./tutorial1.schema" successfully applied
```

## Transforming the data

Now that you have defined the Schema for your data and set up an `extract` for it, it's time to shape the data extracted into a form suitable for loading into the Bubbly Store (and, by extension, into your PostgreSQL database). 

In other words, your aim is to _transform_ input data from your GraphQL query into data that you can _load_ into the Bubbly Store.

The format of the input data is given by the shape of your GraphQL query, and the format of the store is defined by the Schema that you've just applied to it. The only thing we're missing is a specification for _how_ we make the transformation. Enter, the `transform` resource.

A `transform` is just another kind of resource, so you add the following definition to the file `tutorial1.bubbly`, right after your definition of the `extract` resource named `repo_stats`:

```hcl
#
# tutorial1.bubbly
#

# ...

resource "transform" "repo_stats" {
	spec {
		input "data" {}

		data "repo_stats" {
			fields = {
				"id":  join("/", [
								self.input.data.repository.owner.login, 
								self.input.data.repository.name
							])
				"owner":        self.input.data.repository.owner.login
				"repo":         self.input.data.repository.name
				"forks":        self.input.data.repository.forkCount
				"stargazers":   self.input.data.repository.stargazerCount
				"watchers":     self.input.data.repository.watchers.totalCount
				"releases":     self.input.data.repository.releases.totalCount
			}
		}
	}
}
```

:::caution TODO

Add a more direct comment about the mapping between the `spec.data.repo_stats` fields and input data. E.g., "a `transform`'s data block maps the input data from the `extract/repo_stats` resource into the `repo_stats` database table by explicit field matching..."

:::

You note with delight, how easy it is to pick the values of nested fields as well as to generate values for synthetic field by using standard HCL syntax and functions. You wonder, if this would also be the place where you can sum the download counts for different release asset, and so you make a mental note to investigate later.

## Loading the data into the Store

After the data has been extracted from Github via the `extract/repo_stats` resource, and shaped into a format matching your `tutorial1.schema` by the `transform/repo_stats` resource, the next logical step is to store it somewhere safe. To this end, you add a `load` resource to your configuration, right after the `transform`:

```hcl
#
# tutorial1.bubbly
#

# ...

# ...

resource "load" "repo_stats" {
	spec {
		input "data" {}
		data = self.input.data
	}
}
```

## Checkpoint: ETL

You realise at this point that you had defined all three individual stages of the _extract-transform-load_ pipeline. 

Well done! Feeling proud with your technical prowess, you take a comfort break to do some mild stretching and consume a hot beverage of choice.

Your `tutorial1.bubbly` file, at the highest level, is now of the following structure:

```hcl
#
# tutorial1.bubbly
#

resource "extract" "repo_stats" {
	# ...
}

resource "transform" "repo_stats" {
	# ...
}

resource "load" "repo_stats" {
	# ...
}
```

Now all that is left to do is to link these components in a `pipeline`, and _run_ it with a `run` resource!

## Defining the pipeline

The various parts of the pipeline that you have already defined all come together to form a new resource, of kind `pipeline`, which defines the order in which the resources are applied:

```hcl
#
# tutorial1.bubbly
#

# ... extract/transform/load definitions first ...

resource "pipeline" "repo_stats" {
	spec {

		task "extract" {
			resource = "extract/repo_stats"
		}

		task "transform" {
			resource = "transform/repo_stats"
			input "data" {
				value = self.task.extract.value
			}
		}

		task "load" {
			resource = "load/repo_stats"
			input "data" {
				value = self.task.transform.value
			}
		}
	}
}
```

The `pipeline` resource is made of `task` sub-resources. Each `task` refers to one of the resources you have previously defined, and effectively states, "run the resource referenced by this `task` when running this `pipeline`".

## Defining how the pipeline is run

To tie it all together, you add a `run` resource, whose purpose (as its resource kind would suggest!) is to run the pipeline:

```hcl
#
# tutorial1.bubbly
#

# ... extract/transform/load definitions first ...
# ... pipeline definitions next ...

resource "run" "repo_stats" {
	spec {
		resource = "pipeline/repo_stats"
	}
}
```

When the `run/repo_stats` resource is _applied_, it _runs_ its referenced `pipeline/repo_stats` resource and any child resources referenced by its `tasks`. In practical terms, this means that at runtime of the `run/repo_stats`:

 1.  Data from Github's GraphQL API is extracted
 2. The data is _transformed_ into a Schema-compatible format.
 3. The data is _loaded_ to the Bubbly Store

## Checkpoint: Structure of the Bubbly file

You have a quick look over the result, the `tutorial1.bubbly` file, containing the definitions of the `extract`, `transform`, and `load` resources; all tied together in a `pipeline`, ready to `run`:

```hcl
#
# tutorial1.bubbly
#

# ... extract/transform/load definitions first ...

resource "extract" "..." {
	# ...
}

resource "transform" "..." {
	# ...
}

resource "load" "..." {
	# ...
}

# ... pipeline definitions next ...

resource "pipeline" "..." {

	task "..." {
		# ...
	}

	# ...
}

# ... finally, run definitions for pipelines and other resources ...

resource "run "..." {
	# ...
}
```

## Running the pipeline once

With the `pipeline` and `run` resources defined, you are ready to _apply_ all the resources you defined so far in `tutorial1.bubbly` to the Bubbly Agent.

The environment variable `GH_TOKEN` must be set to the value of the access token.

The `run` resource, when applied, will action the `pipeline` and all resources referred therein:

```shell
$ bubbly apply -f ./tutorial1.bubbly
```

and Bubbly responds with 

```shell
resource(s) at path/directory "./tutorial1.bubbly" applied successfully
```

## Inspecting the results

There are several ways to query the Bubbly Agent for the information it holds in its Store. One of those ways is by sending a GraphQL query to the Bubbly Agent.

The following GraphQL query is based on the terminology you've introduced when you set up the Bubbly Schema:

```graphql
{
	repo_stats {
		id
		forks
		stargazers
	}
}
```

This query can be sent to the Bubbly agent with any GraphQL client, including the ones you've used in the [Exploration](#exploration) section. But this time you feel adventurous and decide to try the command-line option, via `curl`.

```shell
$ curl \
   -X POST \
   -H "Content-Type: application/json" \
   --data '{ "query": "{ repo_stats { id forks stargazers } }" }' \
   http://localhost:8111/api/v1/graphql
```

And the Bubbly agent responds with the data:

```json
{"data":{"repo_stats":[{"forks":3683,"id":"docker/compose","stargazers":22052}]}}
```

Although not the most convenient, it certainly is the most vanilla way to run GraphQL queries. It requires only a small stretch of imagination to conjure an image of a dashboard which displays all sorts of useful information from the Bubbly Store. Of course, in order to be useful, that information has to be updated on regular basis. To that end, Bubbly provides the `interval` feature, described in another tutorial.

## Conclusions

Congratulations! In this tutorial, you have written your first Bubbly data pipeline, which:

 1. Collects only the data you need from a GraphQL API end-point;
 2. Transforms it to match a format specific to your database wants/needs;
 3. Loads it into your PostgreSQL database;
 4. Provides a mighty (and slightly terrifying) C-level officer with a query interface for this data.
 
You've gained some experience points, are very pleased with yourself, and are looking forward to what the next day might bring in this new Bubbly-rific world.

## Links to files for this section

* [tutorial1.schema](tutorial1.schema)
* [tutorial1.bubbly](tutorial1.bubbly)

## Next Steps

:::caution TODO

Link to other tutorials when they are ready.

:::

Having set up a basic pipeline, you already have some ideas for tackling the issue of counting  downloads for GitHub releases &ndash; the GitHub API returns that data in a list, and you realise that you can use HCL features in `transform` resource to sum over that list for you!

You also think it would be useful to set up the pipeline to run at regular _intervals_. This way, you'd always have the latest information in your Bubbly Store.

Finally, you wonder, if there is a way to drive the `extract` process with information extracted from a _different_ data source, joining the result...
