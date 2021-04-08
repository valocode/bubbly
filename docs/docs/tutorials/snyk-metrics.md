---
title: Snyk Security Scanning Metrics 
hide_title: false 
hide_table_of_contents: false 
keywords:

- docs
- bubbly
- tutorials
- guides

---

:::note Under Active Development The content for this page is *under active development*. It may be incorrect and/or
incomplete and is liable to change at any time.
:::

# Infrastructure

To run the tutorial, you'll need the Bubbly Agent running with its Store connected to your external Postgresql database,
as described in the Getting Started section. In addition, you will need to have `Snyk` installed on your CLI.

# Introduction
Today started off much like any other day for Richard Locke, head of digital security at BubbliCorp; staring at the coffee machine as
those ever important final drops cascaded into the pot. Coffee in hand, Richard Locke sat down at his computer in his bedro..
er office. Richard's eyes grew wide, 20 messages on Slack! That couldn't be good. Either the engineers were arguing over rust
vs go again, or someone deleted a production DB again... Uh oh, all of the messages were from Mr Blob (CFO at BubbliCorp).

Looks like over the weekend, one of BubbliCorp's servers were hacked using an unpatched exploit! "HOW COULD YOU LET THIS
HAPPEN!?", fumed Blob. Richard, tried to explain that there just wasn't enough time allocated to manually inspecting all of
their environments. Richard noted that security only really seemed to matter when things go wrong. Thankfully, one of the devops
engineers had been able to patch the exploit over the night, but the damage had been done. The homepage banner had proudly
proclaimed "BubbliCorp SUX!" for over 8 hours. The company image!

In no uncertain terms, Richard was informed that this had "better not happen again...". The implications were clear. With no
budget to hire an army of interns to track the servers' status, things looked grim.

Richard suddenly remembered hearing good things from the team working with Github integrations who were using Bubbly. Maybe
he could write a pipeline to solve this!? (spoiler, he can)

`TLDR: You can use Bubbly to ingest results from scanning tools like Snyk`

In this example, we will take a look at the output from a popular security scan and static analysis tool,
`Snyk`.

`Snyk` will scan your codebase, and search for any dependencies that include known security vulnerabilities or `CVE`.

You can run `Snyk` locally in your repository with the following command:

`snyk test --json-file-output=./snyk.json`

This will store the output of the command in `json` format, to the path of your choosing. If you open the resulting
file, you will notice results like the following snippet.

```json
{
  "vulnerabilities": [
    {
      "id": "SNYK-GOLANG-GITHUBCOMDGRIJALVAJWTGO-596515",
      "cvssScore": 7.5,
      "identifiers": {
        "CVE": [
          "CVE-2020-26160"
        ],
        "CWE": [
          "CWE-287"
        ]
      },
      "language": "golang",
      "severity": "high",
      "name": "github.com/dgrijalva/jwt-go",
      "version": "3.2.0"
    },
    {
      "id": "SNYK-GOLANG-GITHUBCOMDGRIJALVAJWTGO-596515",
      "cvssScore": 5.1,
      "identifiers": {
        "CVE": [
          "CVE-2020-26162"
        ],
        "CWE": [
          "CWE-287"
        ]
      },
      "language": "golang",
      "severity": "high",
      "name": "github.com/dgrijalva/jwt-go",
      "version": "3.2.0"
    }
  ]
}
```

This lists all the current vulnerabilities that were found in the dependencies in your project. In this example, we are
interested in writing a `bubbly` file that will define a criteria that warns us when `Snyk` has found any
vulnerabilities with a severity marked `high`.

# Updating the Schema

```hcl
table "snyk_vulnerabilities" {
  field "id" {
    type = string
    unique = false
  }
  field "name" {
    type = string
    unique = false
  }
  field "creationTime" {
    type = string
    unique = false
  }
  field "cvssScore" {
    type = string
    unique = false
  }
  field "language" {
    type = string
    unique = false
  }
  field "severity" {
    type = string
    unique = false
  }
  field "version" {
    type = string
    unique = false
  }
  field "cve" {
    type = string
    unique = false
  }
}
```

# Writing the pipeline

In order to successfully get the data to Bubbly, we need the following resources:

## Extract

Here is an example of what an `Extract` resource could look like for our use case. While there are many more fields
available in the `Snyk` output that we examined above, we will only select a few to give a feel for a simple `Bubbly`
file

```hcl
resource "extract" "snyk" {
  api_version = "v1"
  spec {
    input "file" {}
    type = "json"
    source {
      file = self.input.file
      format = object({
        vulnerabilities: list(object{
          id: string
          cvssScore: number,
          identifiers: object({
            CVE: list(string),
            CWE: list(string),
          }),
          language: string,
          severity: string,
          name: string,
          version: string
        })
      })
    }
  }
}
```

Here, we simply layout an `extract` resource, and specify what format we expect the data to look like. In our example,
we are only focused on finding vulnerabilities with a severity of `high`.

## Transform

In this stage, we will write a `transform` resource, that will convert data from the input type (in this case the `Snyk`
json) to a format that matches the `Bubbly Schema`.

```hcl
resource "transform" "snyk" {
  api_version = "v1"
  spec {
    input "data" {}

    dynamic "data" {
      for_each = self.input.data.vulnerabilities
      iterator = it
      labels = [
        "snyk_vulnerabilities"]
      content {
        fields = {
          "id": it.value.id
          "name": it.value.name
          "creationTime": it.value.creationTime
          "cvssScore": it.value.cvssScore
          "language": it.value.language
          "severity": it.value.severity
          "version": it.value.version
          "cve": it.value.identifiers.CVE[0]
        }
      }
    }
  }
}
```

## Load

Now that our data has been transformed, the next step, is to `load` the data into the `Bubbly Server`.

```hcl
resource "load" "snyk" {
  api_version = "v1"
  spec {
    input "data" {}
    data = self.input.data
  }
}
```

## Pipeline

Now that we have defined out ETL phase, it is time to define our pipeline. In this phase, we will need to make sure to
note where our input file is coming from. In our case `./testdata/snyk/snyk.json`. The pipeline is essentially stitching
together the earlier phases into one coherent place.

```hcl
resource "pipeline" "snyk" {
  api_version = "v1"
  spec {
    input "snyk_file" {
      default = "./testdata/snyk/snyk.json"
    }
    task "extract_snyk" {
      resource = "extract/snyk"
      input "file" {
        value = self.input.snyk_file
      }
    }
    task "transform" {
      resource = "transform/snyk"
      input "data" {
        value = self.task.extract_snyk.value
      }
    }
    task "load" {
      resource = "load/snyk"
      input "data" {
        value = self.task.transform.value
      }
    }
  }
}
```

## Run

Now that the pipeline is in place, we need to inform the `API` that we want to run the pipeline resource.

```hcl
resource "run" "snyk" {
  api_version = "v1"
  spec {
    resource = "pipeline/snyk"
  }
}
```

## Run

```hcl
resource "run" "sonarqube/extract" {
  api_version = "v1"
  spec {
    resource = "extract/snyk"
    input "file" {
      value = "./testdata/snyk/snyk.json"
    }
  }
}
```

## Query

Now that we have processed and uploaded our data, it is time to get some insights into our project's status. In our
case, we want the `Bubbly` to inform us if there are any vulnerabilities with a severity of `high`. We will write a
quick query using `Graphql` to fetch the severity of the vulnerabilities.

```hcl
resource "query" "snyk_data" {
  api_version = "v1"
  spec {
    query = <<EOT
            {
                snyk_vulnerabilities {
                    severity
                }
            }
        EOT
  }
}
```

## Criteria

Since we are looking to ensure that there are no high severity vulnerabilities, we will need to write an expression to
ensure this.

```hcl
resource "criteria" "snyk_status" {
  api_version = "v1"
  spec {
    query "snyk_data" {}
    condition "no_high" {
      value = length(regexall("high", self.query.snyk_data.value)) == 0
    }

    operation "release" {
      value = self.condition.no_high.value
    }
  }
}
```

## Run

The final step, is simply to provide a `run resource` to inform `Bubbly` that it should process the criteria.

```hcl
resource "run" "snyk_criteria" {
  api_version = "v1"
  spec {
    resource = "criteria/snyk_status"
  }
}
```

## Checking the output

Now that all of the steps have run, it's time to check the output! In the CLI, we will write:
`bubbly get run/snyk_criteria --events` to see what the status of our criteria is. In our sample case, there is a
vulnerability with a severity of `high`, so we will see the following output:
`run/snyk_criteria  RunFailure       2021-03-23T15:13:03+02:00`

# Conclusion

It looks like now, Richard has a pipeline setup that will inform him and the release team when they have code that relies
on dependencies with severe vulnerabilities. This will give confidence in code quality and security as it moves into production.
It might also keep Richard's job safe for the time being (at least until someone finds those sticky notes with prod passwords...)