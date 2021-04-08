---
title: Bubbly Schema
hide_title: false
hide_table_of_contents: false
keywords:
- docs
- bubbly
- schema
- migration
---

:::note Under Active Development
The content for this page is *under active development*. It
may be
incorrect and/or
incomplete and is liable to change at any time.
:::

## Introduction

As introduced in our [Core Concepts: Bubbly Schema](../introduction/core-concepts#bubbly-schema),
the Bubbly Schema is a declarative definition of the format you expect of your data. 
In practical terms, it represents a definition of your underlying database schema 
in the Bubbly Language.

Bubbly schemas are very simple to write, and are made up of three things:

1. Tables - define database tables 
2. Fields - define columns for database tables
3. Joins - define relationships between tables


## Specification Reference

The following attributes and blocks are supported:

- `table "<BLOCK LABEL>"`: one or more configuration block representing
  an entity about which you want to store data. Bubbly uses `table` blocks to construct
  the underlying database tables.The label of the block specifies the name of the 
  database table. Within this block, the following attributes and blocks are supported:
    - `field "<BLOCK LABEL>"`: One or more fields representing database columns for a table.
      The label of the block specifies the name of the table column. Within this block, 
      the following attributes are supported:
        - `type`: The data type expected within this database column.
        - `unique`: (Optional) Specify whether all values in this column must be unique. Default: `false`
    - `table "<BLOCK LABEL>"`: (Optional) Zero or more nested `table` configuration blocks. 
      These follow the same specification as the root `table` configuration block.
    - `join "<BLOCK LABEL>"`: (Optional) Zero or more configuration blocks specifying
    joins between database tables. Within this block, the following attributes are supported:
      - `unique`: (Optional) Specify whether all values in this join must be unique. Default: `false`
      :::note
      TODO: Clarify function of a join's `unique` attribute
      :::

## Example

Below is an example Bubbly Schema:

```hcl
table "repo" {
    field "name" {
        type = string
        unique = true
    }
 
    table "repo_version" {
        field "commit" {
            type = string
            unique = true
        }
        field "tag" {
            type = string
        }
        field "branch" {
            type = string
        }
    }
}
 
table "code_issue" {
    join "repo_version" { }
    field "name" {
        type = string
    }
    field "severity" {
        type = string
    }
    field "type" {
        type = string
    }
}
```

This Schema:

1. Creates a `repo` table with the `repo_version` sub-table
    1. This makes the `repo_version` have an _implicit_ join to `repo`
2. Creates a `code_issue` table
    1. This table has as an _explicit_ join to the `repo_version` table

The result: a `code_issue` belongs to a `repo_version` which belongs to a `repo`.

## Applying a Bubbly Schema

Bubbly Schemas are applied in a similar mechanism to which you apply Bubbly Resources:
via the Bubbly CLI. However, to protect users from unintentionally applying Bubbly Schema files
along with their Bubbly Resources, we have separated the schema management to a different command:
`bubbly schema apply`.

At apply-time, Bubbly handles the creating and updating the schema of the
backend database.

### Example:

Imagine you have defined a Bubbly Schema file in your current directory (`./schema.bubbly`).
To apply this schema, run `bubbly schema apply -f ./schema.bubbly`.

## Migrations

As [previously mentioned](../introduction/core-concepts#bubbly-schema), 
Bubbly is capable of evolving with you and your data. As a result, Bubbly supports
the re-applying of Bubbly Schema files, and will update the schema of the backend
database to match the latest Bubbly Schema applied.

