---
title: bubbly apply
sidebar_label: bubbly apply
hide_title: false
hide_table_of_contents: false
description: Bubbly CLI - bubbly apply
keywords:
- docs
- bubbly
- cli
- apply
---

## Synopsis

Apply bubbly resources to a bubbly API server



```
bubbly apply (-f (FILENAME | DIRECTORY)) [flags]
```

### Examples

```
  # Apply the bubbly resources in the file ./main.bubbly
  bubbly apply -f ./main.bubbly
  
  # Apply the configuration in the directory ./resources
  bubbly apply -f ./resources
```

### Options

```
  -f, --filename string   filename or directory that contains the bubbly resources to apply
  -h, --help              help for apply
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
