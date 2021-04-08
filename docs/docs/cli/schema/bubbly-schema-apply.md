---
title: bubbly schema apply
sidebar_label: bubbly schema apply
hide_title: false
hide_table_of_contents: false
description: Bubbly CLI - bubbly schema apply
keywords:
- docs
- bubbly
- cli
- schema
- apply
---

### Synopsis

Apply a bubbly schema

    $ bubbly schema apply -f FILENAME



```
bubbly schema apply -f FILENAME [flags]
```

### Examples

```
  # Apply a bubbly schema located in a specific file
  bubbly schema apply -f ./schema.bubbly
```

### Options

```
  -f, --filename string   filename that contains the .bubbly schema file to apply
  -h, --help              help for apply
```

### Options inherited from parent commands

```
      --debug         specify whether to enable debug logging
      --host string   bubbly API server host (default "127.0.0.1")
      --port string   bubbly API server port (default "8111")
```

### SEE ALSO

* [bubbly schema](../bubbly-schema)	 - manage your bubbly schema