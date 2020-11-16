# CLI Design Document

- [CLI Design Document](#cli-design-document)
  - [Command Families](#command-families)
  - [Optional Command Families](#optional-command-families)
  - [Family Type Descriptions](#family-type-descriptions)
    - [Declarative](#declarative)
      - [Similar Strategies](#similar-strategies)
    - [Print](#print)
      - [Similar Strategies](#similar-strategies-1)
    - [Serve](#serve)
      - [Similar Strategies](#similar-strategies-2)
- [Commands](#commands)
  - [apply](#apply)
    - [Usage](#usage)
    - [Flags](#flags)
  - [query](#query)
    - [Usage](#usage-1)
  - [delete](#delete)
    - [Usage](#usage-2)
      - [Possible Imperative Extensions](#possible-imperative-extensions)
    - [Flags](#flags-1)
  - [get](#get)
    - [Usage](#usage-3)
    - [Flags](#flags-2)
  - [version](#version)
    - [Usage](#usage-4)
    - [Output](#output)
    - [Flags](#flags-3)
  - [describe](#describe)
    - [Usage](#usage-5)
    - [Flags](#flags-4)
  - [logs](#logs)
  - [explain](#explain)
    - [Usage](#usage-6)
    - [Flags](#flags-5)
  - [api-resources](#api-resources)
    - [Usage](#usage-7)
    - [Flags](#flags-6)
  - [api-versions](#api-versions)
    - [Usage](#usage-8)
  - [diff](#diff)
    - [Usage](#usage-9)
  - [flags](#flags-7)
    - [Usage](#usage-10)
  - [server](#server)
    - [Usage](#usage-11)
    - [Subcommands](#subcommands)
    - [Flags](#flags-8)
- [Global flags](#global-flags)
  - [log levels](#log-levels)
    - [Background](#background)
- [Open questions](#open-questions)
- [References](#references)

## Command Families

| Family      | Purpose                                                             | Description                                                                                                                                         | Commands                                                                                |
| ----------- | ------------------------------------------------------------------- | --------------------------------------------------------------------------------------------------------------------------------------------------- | --------------------------------------------------------------------------------------- |
| Declarative | Declarative Resource Management for Data Ingestion and Querying     | Declaratively manage the data ingestion process and queries that extract subsets of this data                                                       | `apply`, `query`                                                                        |
| Print       | Printing information about Bubbly server, client and Resource state | Transparent access to all aspects of the Bubbly system, such that users can print/debug system and Resource state information from the command line | `status`, `describe`, `version`, `logs`, `diff`, `api-resources`, `api-versions`, `get` |
| Serve       | Running Server(s)                                                   | Running the Bubbly server                                                                                                                           | `server`                                                                                |

Note: Commands of differing families can still make use of other families' features. For example, `get` is still able to make use of declarative Resource declarations to identify the resource to 'get' from a server.

- `bubbly get (-f (FILENAME | DIRECTORY))`

## Optional Command Families

We may wish to avoid implementing imperative commands that are functionally equivalent to declarative ones.

| Family     | Purpose                                           | Description                                                                                                                                           | Commands                   |
| ---------- | ------------------------------------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------- | -------------------------- |
| Imperative | Imperative Resource Management for Data ingestion | Commands to manage Bubbly Resources via CLI. `bubbly create extract junit`, `bubbly delete extract junit`. Similar to `kubectl create deployment app` | `edit`, `create`, `delete` |

## Family Type Descriptions

### Declarative

The preferred approach for managing Resources is through declarative `.hcl` files using the Resource's equivalent `ResourceKind`. These Resources provide a declarative interface for defining data ingestion and querying, and are _applied_ to the Bubbly server with the `bubbly apply` command. This command reads a local (or remote) `.hcl` file structure, identifies the Resource(s), and modifies Bubbly server state to reflect the declared intent outlined by the Resource(s).

#### Similar Strategies

1. `kubectl apply`
2. `terraform apply`
3. `nomad run`

### Print

This family of commands provides transparent access to all aspects of the Bubbly system, such that both system and Resource state information is open, accessible and easily monitorable to Bubbly users via command line.

Collectively, these CLI commands enable the following features:

- _Summarising_ state and information about the Bubbly client, Bubbly server and Bubbly Resources
  - client and server:
    - `bubbly version`
  - server:
    - `bubbly server status [-o simple]`
  - Resources:
    - `bubbly get extract junit [-o simple]`
    - `bubbly get pipelineRun junit -o events`
- Printing complete state and information about Bubbly client, Bubbly server and Bubbly Resources
  - server: `bubbly server status -o detailed`
  - Resources:
    - `bubbly describe extract junit`
    - `bubbly get extract junit -o detailed`
- Understanding of the Bubbly API versions, their respective Resources and the definitions of these Resources:
  - `bubbly api-versions`
  - `bubbly api-resources`
  - `bubbly explain extract.spec.type`
- Printing Resource logs:
  - `bubbly logs pipelineRun junit`
- Visualising Resource changes:
  - `bubbly diff pipelineRun junit`

#### Similar Strategies

1. Almost all commands are inspired by `kubectl`:
   - `get`, `describe`, `explain`, `logs`, `version`, `api-versions`, `api-resources`
2. `terraform plan` and `kubectl diff` inspire `bubbly diff`

### Serve

Commands for running the Bubbly server, enabling a bundling of the Bubbly client and server into the same binary.

#### Similar Strategies

1. `nomad agent -server`
2. `consul agent -server`

# Commands

|   Family    |             Command             |                                                          Description                                                          | Aliases |
| :---------: | :-----------------------------: | :---------------------------------------------------------------------------------------------------------------------------: | :-----: |
| Declarative |         [apply](#apply)         |                  Apply a Bubbly configuration (collection of 1 or more Bubbly Resources) to a Bubbly server                   |    a    |
| Declarative |         [query](#query)         |                 Apply a Bubbly configuration (consisting of  1 or more `query` Resources) to a Bubbly server                  |    q    |
| Declarative |        [delete](#delete)        |                           Delete Resources specified by a Bubbly configuration from a Bubbly server                           |         |
|    Print    |           [get](#get)           |                               Display one or many resources of a specific bubbly resource type                                |         |
|    Print    |       [version](#version)       |                            Show details of the Bubbly client and Bubbly server version information                            |         |
|    Print    |      [describe](#describe)      |                               Show details of a specific bubbly resource or group of resources                                |         |
|    Print    |          [logs](#logs)          |                                        Print the logs for a specific Bubbly resource.                                         |         |
|    Print    |       [explain](#explain)       |                               Describes the fields associated with each supported API resources                               |         |
|    Print    | [api-resources](#api-resources) |                                       Prints the supported API resources of the server                                        |   ar    |
|    Print    |  [api-versions](#api-versions)  |                                        Prints the supported API versions of the server                                        |   av    |
|    Print    |          [diff](#diff)          | Show details of proposed changes to Bubbly resources between the current server configuration and the provided configuration. |         |
|    Print    |         [flags](#flags)         |                                Show details of global flags that can be passed to any command                                 |         |
|    Serve    |        [server](#server)        |                                                    Starts a Bubbly server                                                     |         |

## apply

### Usage

`bubbly apply (-f (FILENAME | DIRECTORY)) [flags]`

### Flags

|     Flag      |       Options        | Default |                                                                         Description                                                                          |
| :-----------: | :------------------: | :-----: | :----------------------------------------------------------------------------------------------------------------------------------------------------------: |
| --server-side |     true, false      |  false  |                                                 If true, the apply runs in the server instead of the client.                                                 |
|   --dry-run   | none, server, client |  none   | Client: Print the `[]DataTable` (or other) object that would be sent to the Bubbly server. Server: Send the object server-side without persisting the object |

## query

### Usage

`bubbly query (-f (FILENAME | DIRECTORY)) [flags]`

## delete

### Usage

`bubbly delete (-f (FILENAME | DIRECTORY)) [flags]`

#### Possible Imperative Extensions

Delete Bubbly Resource(s) by resource type and name selector: `bubbly delete (TYPE [(NAME | --all)]) [options]`

Examples:

- `bubbly delete extracts --all`
- `bubbly delete extract junit`

### Flags

|   Flag    |   Options   | Default |                                           Description                                            |
| :-------: | :---------: | :-----: | :----------------------------------------------------------------------------------------------: |
| --dry-run | true, false |  false  | Print the Resource(s) that would be deleted, without actually deleting it from the Bubbly server |

## get

### Usage

`bubbly get ((TYPE NAME_PREFIX) | GENERIC) [flags]`

| Usage Term | Options                                                                  | Description                                                            |
| ---------- | ------------------------------------------------------------------------ | ---------------------------------------------------------------------- |
| TYPE       | `extract`, `transform`, `load`                                           | Bubbly Resource Type                                                   |
| GENERIC    | `all` (describes all Resources), `system` (describes core Bubbly system) | An object to be described that does not match a specific Resource type |

### Flags

|     Flag      |   Options   | Default |                    Description                    |
| :-----------: | :---------: | :-----: | :-----------------------------------------------: |
|  --failures   | true, false |  false  | Describe only failures in Resources or the system |
| --api-version | v1, v2, ... |         |  Describe only Resources of a given api-version   |

## version

### Usage

`bubbly version [flags]`

### Output

`bubbly version` should provide something akin to the following for both client and server:

```json

version.Info{
    Major:"1",
    Minor:"18",
    GitVersion:"v1.18.2",
    GitCommit:"52c56ce7a8272c798dbc29846288d7cd9fbae032",
    BuildDate:"2020-04-16T23:34:25Z",
    GoVersion:"go1.14.2",
    Compiler:"gc",
    Platform:"darwin/amd64"
}

```

### Flags

|   Flag   |   Options   | Default |                      Description                       |
| :------: | :---------: | :-----: | :----------------------------------------------------: |
| --short  | true, false |  false  | Print just the version number of the client and server |
| --client | true, false |  false  |         Print just client version information          |

## describe

### Usage

`bubbly describe ((TYPE NAME_PREFIX) | GENERIC) [flags]`

| Usage Term | Options                                                                  | Description                                                            |
| ---------- | ------------------------------------------------------------------------ | ---------------------------------------------------------------------- |
| TYPE       | `extract`, `transform`, `load`                                           | Bubbly Resource Type                                                   |
| GENERIC    | `all` (describes all Resources), `system` (describes core Bubbly system) | An object to be described that does not match a specific Resource type |

### Flags

|     Flag      |   Options   | Default |                    Description                    |
| :-----------: | :---------: | :-----: | :-----------------------------------------------: |
|  --failures   | true, false |  false  | Describe only failures in Resources or the system |
| --api-version | v1, v2, ... |         |  Describe only Resources of a given api-version   |

## logs

// TODO / REMOVE

## explain

### Usage

`bubbly explain TYPE[.<field-name>...] [flags]`

Examples:

- `bubbly explain extract`
- `bubbly explain extract.spec.type`

### Flags

|    Flag     |   Options   | Default |              Description               |
| :---------: | :---------: | :-----: | :------------------------------------: |
| --recursive | true, false |  false  | Print the fields of fields recursively |

## api-resources

### Usage

`kubectl api-resources [flags]`

### Flags

|      Flag       |     Options     | Default |                                                 Description                                                 |
| :-------------: | :-------------: | :-----: | :---------------------------------------------------------------------------------------------------------: |
|  --api-version  |   v1, v2, ...   |   all   |                                 Print only Resources of a given API version                                 |
| --resource-type | runner, definer |   all   | Print only Resources of a given functional type. `pipelineRun` is a `runner`, `pipeline` is a `definer`  |

## api-versions

### Usage

`bubbly api-versions [flags]`

## diff

Similar to `terraform plan` or `kubectl diff`, `bubbly diff` provides a safe, convenient way to check that proposed changes to Bubbly resources will have the desired effect.

### Usage

`bubbly diff -f FILENAME [flags]`

## flags

### Usage

`bubbly flags [flags]`

## server

### Usage

`bubbly server [flags]`

### Subcommands

| Subcommand |                          Description                           |
| :--------: | :------------------------------------------------------------: |
|   status   | Display a list of the known Bubbly servers and their statuses. |

### Flags

| Flag  |   Options   | Default |                                          Description                                           |
| :---: | :---------: | :-----: | :--------------------------------------------------------------------------------------------: |
| --dev | true, false |  false  | Start the Bubbly server in development mode (no data persistence). **Not** for production use. |

# Global flags

These are flags that can be passed to any command. Documented by `bubbly flags`.

|          Flag          |           Options            | Default |                                    Description                                     |
| :--------------------: | :--------------------------: | :-----: | :--------------------------------------------------------------------------------: |
|     -f, --filename     |                              |         |             directory or file that contains the configuration to apply             |
|    -R, --recursive     |         true, false          |  false  |             Process the directory used in -f, --filename recursively.              |
|      -o, --output      | json, yaml, simple, detailed | simple  |                                   Output format                                    |
|      -s, --server      |                              |         |                   The address and port of the Bubbly API server                    |
|        --token         |                              |         |                 Token for authentication to the Bubbly API server                  |
|       --log-dir        |                              |         |             Directory to write logs to. If empty, write logs to stderr             |
|       --log-file       |                              |         |               File to write logs to. If empty, write logs to stderr                |
| --match-server-version |                              |  false  |                   Require server version to match client version                   |
|      --log-level       |           0,1,2,3            |    1    | Bubbly output verbosity. Integer increase corresponds to an increase in verbosity. |

## log levels

Simplified version of Kubectl's [strategy](https://kubernetes.io/docs/reference/kubectl/cheatsheet/).

| Log-level | Description                                                                                                                           |                             Examples                              |
| --------- | ------------------------------------------------------------------------------------------------------------------------------------- | :---------------------------------------------------------------: |
| 0         | Generally useful for this to always be visible to a Bubbly server operator / Bubbly client user.                                      |      CLI argument handling, HCL configuration parsing error       |
| 1         | Useful state information about the Bubbly client and server, and log messages that indicate changes to this state. **Default level**. |            Logging HTTP requests and their exit codes.            |
| 2         | Debug level verbosity.                                                                                                                |                                                                   |
| 3         | Trace level verbosity                                                                                                                 | Context to understand the steps leading up to errors and warnings |

### Background

Kubectl actually has its own [logging convention](https://github.com/kubernetes/community/blob/master/contributors/devel/sig-instrumentation/logging.md) which it implements via the [klog Go package](https://github.com/kubernetes/community/blob/master/contributors/devel/sig-instrumentation/logging.md). Interesting read as inspiration for our logging strategy, perhaps.

# Open questions

1. Thoughts on support for imperative actions (e.g. `kubectl create deployment app`)? Do we want to offer **any** kind of imperative actions to our users? I wonder if there's research out there describing the adoption of `kubectl create ...` compared to `kubectl apply` to get a better idea of how our users may wish to interact. At least the documentation even lists its use as "Development Only".
   - AFAIK we will need a way to support deletion (e.g. `kubectl delete deployment app`), but this could also be done declaratively:
     - through a `deletor` Bubbly resource
     - `bubbly delete -f FILENAME`

2. Are both `get` and `describe` overkill?
   - I quite like the distinction in k8s, but perhaps for our purposes both are not needed and only serve to confuse users. It may be better to consolidate into one command and use flags (`-o (simple | detailed)`) for filtering of information.
3. Is `logs` necessary, since we're not really "tracking" resources in the same way `kubectl logs` or `nomad alloc logs` is?
   - `-o events` might be sufficient as a mechanism for retrieving state changes to a resource (e.g. `pipelineRun`) and its associated `ResourceOutput`.

# References

- [tkn](https://github.com/tektoncd/pipeline) for its approach to the use of `runner` resources vs `definer` resources
- [kubectl](https://github.com/kubernetes/kubectl) and [kubectl docs](https://kubectl.docs.kubernetes.io/)
- [Hashicorp Nomad](https://github.com/hashicorp/nomad)
