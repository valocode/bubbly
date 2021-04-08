---
title: Getting Started
hide_title: false
hide_table_of_contents: false
description: Installation of bubbly
keywords:
- docs
- bubbly
- installation
- deployment
- agent
---

:::note Under Active Development
The content for this page is *under active development*. It
may be
incorrect and/or
incomplete and is liable to change at any time.
:::

## Installation

### Homebrew

`brew install valocode/tap/bubbly`

### Releases

Binaries for macOS, Windows and Linux AMD64/ARM are available to download on the
[release page](https://github.com/valocode/bubbly/releases)

## Deployment

### Single-Server

#### Pre-requisites

Bubbly requires an existing postgresql database.

#### Running

Start the [Bubbly Agent](../introduction/core-concepts#agent) in single server mode, 
which starts all of the Agent's features on the same machine: 

`bubbly agent`

The `bubbly agent` command defaults to connecting to a postgres database available on
`postgres:5432`. Use `bubbly agent -h` for command-line flags and `bubbly env` for
environment variables for configuring this.


### Multi-Server

:::note Under Active Development
The content for this page is *under active development*. It
may be
incorrect and/or
incomplete and is liable to change at any time.
:::

In future, Bubbly will support a distributed deployment model in which you can
freely customise the number of instances of each Bubbly Agent feature (as long
as you have at least one of each).

For this, you have two strategies:

1. Native: All communication between the running Bubbly Agents handled by our Bubbly-embedded NATS server.
2. Custom: Roll your own dedicated [NATS server](https://github.com/nats-io/nats-server).
   We recommend this strategy for more demanding deployments.
   
## Next Steps

- [Tutorials: Getting Started](../tutorials/github-metrics) 
  provides a practical tutorial for immediately getting started with Bubbly
