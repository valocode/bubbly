# Unique

Contains tests for testing the unique constraints

Ensure that all the data blocks in the [data.hcl](./data.hcl) files are unique, because the tests will apply each multiple times, and then verify that there is only one of each.

The schema/tables exist in [tables.hcl](./tables.hcl).

## Update test

See files [data_update.hcl](data_update.hcl) and [tables_update.hcl](tables_update.hcl)

Goal: table `tu1` should only be created once for each `tu2`.
How can we make it so that we don't create more `tu1` instances everytime we create/update `tu2`?
