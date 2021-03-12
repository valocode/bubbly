---
title: bubbly get
sidebar_label: bubbly get
hide_title: false
hide_table_of_contents: false
description: Bubbly CLI - bubbly get
keywords:
- docs
- bubbly
- cli
- get
---

## Synopsis

Display one or many bubbly resources

```
bubbly get (KIND | ID | all) [flags]
```

### Examples

```
  # Display all bubbly resources stored on the bubbly server
  bubbly get all
  
  # Display all bubbly resources of kind extract
  bubbly get extract
  
  # Display a specific bubbly resource
  bubbly get default/extract/sonarqube
  
  # Display a specific bubbly resource and associated events
  bubbly get default/extract/sonarqube --events
```

### Options

```
  -e, --events   specify whether to display resource events
  -h, --help     help for get
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
