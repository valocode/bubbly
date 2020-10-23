# Design

This document outlines the high-level design for Bubbly.

## Definitions

Below are some useful definitions that help us communicate more clearly and succintly about the design of bubbly.

**Resource:** bubbly is composed a different `Resource` kinds, that can be applied.

**Data Schema:** the desired structure of the data that we want to collect, consisting of tables and fields.

**Data Source:** is a source of data, such as a JSON file or REST API (e.g. a Jira server could be seen as a Data Source).

**Importer:** is a `ResourceKind` that defines a Data Source and a format for the data, and returns a `cty.Value` representation of the external data.

**Translator:** is a `ResourceKind` that defines a translation or conversion of data from an `Importer` and ouputs a mapping of the imported data to the defined Data Schema.

**Query:** is a `ResourceKind` that defines question to the bubbly server which returns some data.

**Criteria:** is a `ResourceKind` that defines a list of queries and conditions on the individal or combined queries, providing a pass or fail result.

## API Design

Everything in Bubbly is based around a kind of `Resource` that can be *applied*.
Certain resources allow for a list of `input` values to be provded, so that the resource can be reused and keep things DRY.

### Example HCL

```hcl
// Everything is a resource, and all resources have an output.
// The output can be something like:
// output:
//   - status: sucess
//   - value: <the value from the importer>
resource "importer" "junit" {
    api_version = "v1"
    // this is an input to the importer to make it reusable
    spec {
        input "file" {}
        type = "xml"
        source {
            // reference the "self" here to get the input
            file = self.spec.input.file
        }
        format = object({
            testsuites: object({
                duration: number,
                testsuite: list(object({
                    failures: number,
                    name: string,
                    package: string,
                    tests: number,
                    time: number,
                    testcase: list(object({
                        classname: string
                        name: string
                        time: number
                    }))
                }))
            })
        })
    }
}


resource "translator" "junit" {
    api_version = "v1"
    spec {
        input "data" {}
        // this is some crazy dynamic HCL stuff to create the "data" blocks
        dynamic "data" {
            for_each = self.spec.input.data
            iterator = it
            labels = ["test_table", "table_${it.value}"]
            content {
                field "test_field" {
                    value = "This field exists in table_${it.value}"
                }
            }
        }
    }
}

// This is the upload step, just renamed to publish...
resource "publish" "junit" {
    api_version = "v1"
    spec {
        input "data" {}
        // what do we need here?!
        data = self.spec.input.data
    }
}

// How do we tie all the above resources together...? In a pipeline!
// A pipeline is just another reusable resource, and only a pipelineRun
// actually triggers a pipeline to run
resource "pipeline" "junit" {
    api_version = "v1"
    spec {
        input "file" {}
        // Each task in a pipeline has an output, similar to resources,
        // so that task outputs can be referenced
        task "import" {
            resource = resource.importer.junit
            input "xml_file" {
                value = self.spec.input.file
            }
        }
        task "translator" {
            resource = resource.translator.junit
            input "results" {
                // here we reference the output of the task "importer"
                value = self.spec.task.import.value
            }
        }
        task "publish" {
            resource = resource.publish.junit
            input "data" {
                value = self.spec.task.translator.value
            }
        }
    }
}

// A pipelineRun resources executes a pipeline resource.
resource "pipelineRun" "junit" {
    api_version = "v1"
    spec {
        // specify the name of the pipeline
        pipeline = resource.pipeline.junit
        // specify the input required
        input "file" {
            value = "testdata/importer/junit.xml"
        }
    }
}
```
