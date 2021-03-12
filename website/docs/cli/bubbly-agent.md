---
title: bubbly agent
sidebar_label: bubbly agent
hide_title: false
hide_table_of_contents: false
description: Bubbly CLI - bubbly agent
keywords:
- docs
- bubbly
- cli
- agent
---
  
### Synopsis

Starts a bubbly agent. The agent can be configured to run all components, or only some subset, 
depending on the flags provided.

```
bubbly agent [flags]
```

### Examples

```
  # Starts the bubbly agent with all components (API Server, NATS Server, Store and Worker)
  using application defaults
  
  bubbly agent
  
  # Starts the bubbly agent running only the API Server components
  bubbly agent --api-server
```

### Options

```
      --api-server                   whether to run the api server on this agent
      --data-store                   whether to run the data store on this agent
      --data-store-addr string       address of the data store (default "postgres:5432")
      --data-store-database string   database of the data store (default "bubbly")
      --data-store-password string   password of the data store (default "postgres")
      --data-store-provider string   provider of the bubbly data store (default "postgres")
      --data-store-username string   username of the data store (default "postgres")
      --deployment-type string       the type of agent deployment. Options: single (default "single")
  -h, --help                         help for agent
      --nats-server                  whether to run the NATS Server on this agent (default true)
      --nats-server-addr string      address of the NATS Server (default "localhost:4223")
      --nats-server-http-port int    HTTP Port of the NATS Server (default 8222)
      --nats-server-port int         port of the NATS Server (default 4223)
      --worker                       whether to run a bubbly worker on this agent
```

### Options inherited from parent commands

```
      --debug         specify whether to enable debug logging
      --host string   bubbly API server host (default "127.0.0.1")
      --port string   bubbly API server port (default "8111")
```

### SEE ALSO

* [bubbly](bubbly.md)	 - bubbly: release readiness in a bubble

Find more information: https://bubbly.dev
