---
title: Gosec static analysis results
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
as described in the Getting Started section. In addition, you will need to have `gosec`

```shell
go get github.com/securego/gosec/v2/cmd/gosec
```

# Introduction

[gosec](https://github.com/securego/gosec) is a powerful tool that scans for errors in your code by leveraging Go's AST.
Gosec will scan the provided directory, and provide results similar to the example below.

```json
{
  "Issues": [
    {
      "severity": "MEDIUM",
      "confidence": "MEDIUM",
      "cwe": {
        "ID": "118",
        "URL": "https://cwe.mitre.org/data/definitions/118.html"
      },
      "rule_id": "G601",
      "details": "Implicit memory aliasing in for loop.",
      "file": "/Users/username/Projects/bubbly/store/store.go",
      "code": "283: \n284: \t\taddImplicitJoins(schema, t.Tables, \u0026t)\n285: \t\t// Clear the child tables\n",
      "line": "284",
      "column": "38"
    }
  ]
}
```

In this example, we will write a `Bubbly` pipeline that will determine if you have any issues with a `HIGH` severity.

# Updating the Bubbly Schema

First, we will need to update our `Bubbly Schema` to support adding the severity of an issue.

```hcl
table "gosec_issues" {
    field "severity" {
        type = string
    }
}

```

# Writing the Pipeline

Now, we will create a new `.bubbly` file that will handle our new pipeline

## Extract

Here, we will simply extract just the severity from the issues.

```hcl
resource "extract" "gosec" {
  spec {
    input "file" {}
    type = "json"
    source {
      file = self.input.file
      format = object({
        Issues: list(object({
          severity: string,
        }))
      })
    }
  }
}
```

## Transform

Now that we have the severity, it is time to transform it into a format that `Bubbly` can handle.

```hcl
resource "transform" "gosec" {
  spec {
    input "data" {}
    dynamic "data" {
      for_each = self.input.data.Issues
      iterator = it
      labels = ["gosec_issues"]
      content {
        fields = {
          "severity": it.value.severity
        }
      }
    }
  }
}
```

## Load

The next step is to load the data!

```hcl
resource "load" "gosec" {
  spec {
    input "data" {}
    data = self.input.data
  }
}
```

## Pipeline

In this step, we will stitch all the data together in a single pipeline.

```hcl
resource "pipeline" "gosec" {
  spec {
    input "gosec_results" {
      default = "./testdata/gosec/results.json"
    }
    task "extract_gosec" {
      resource = "extract/gosec"
      input "file" {
        value = self.input.gosec_results
      }
    }
    task "transform" {
      resource = "transform/gosec"
      input "data" {
        value = self.task.extract_gosec.value
      }
    }
    task "load" {
      resource = "load/gosec"
      input "data" {
        value = self.task.transform.value
      }
    }
  }
}
```

## Run

Next, we will add a simple `run` resource for the pipeline.

```hcl
resource "run" "gosec" {
  spec {
    resource = "pipeline/gosec"
  }
}
```

## Run

We will also create a `run` resource for extracting the results from a specified file. Make sure to change the path to
reflect the location of the file you are looking to upload.

```hcl
resource "run" "gosec/extract" {
  spec {
    resource = "extract/gosec"
    input "file" {
      value = "./testdata/gosec/results.json"
    }
  }
}
```

## Query

Now that the data is uploaded, it is time to start querying the data. Here we will query for just the severity of the
issue.

```hcl
resource "query" "gosec_data" {
  spec {
    query = <<EOT
      {
        gosec_issues {
          severity
        }
      }
    EOT
  }
}
```

## Criteria

Here, we will write a quick criteria, that will pass on the condition that there are no issues with "HIGH" severity

```hcl
resource "criteria" "gosec_status" {
  spec {
    query "gosec_data" {}
    condition "no_high" {
      value = length(regexall("HIGH", self.query.gosec_data.value)) == 0
    }

    operation "release" {
      value = self.condition.no_high.value
    }
  }
}
```

## Run

The final step is to just write a `run` resource for the criteria as follows.

```hcl
resource "run" "gosec_criteria" {
  spec {
    resource = "criteria/gosec_status"
  }
}
```

# Results

In this tutorial, we took a look at what goes into writing a pipeline for ensuring that your code doesn't include any
issues with "HIGH" severity. `Gosec` is a powerful tool, and there are many other interesting ways to inspect your code
quality.
