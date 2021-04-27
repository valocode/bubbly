# Release Readiness Checklist

For producing a good test of the release readiness checklist we need 3 things:

1. A schema, which can be found in [schema.hcl](./schema.hcl)
2. Some test data, which can be found in [data.hcl](./data.hcl)
3. A GraphQL query, which exists below

## Query

In order to check if a release is complete or not, we have to make a query on the release.
For this we need a specific `release` `_id`, or some other way to uniquely identify the `release`, but this should be available when a user clicks to show a `release`.

Then we get all the `release_item(s)` that belong to a release, and for each of those, we get the `release_entry(s)`, which are the recorded/logged events against that `release_item` (e.g. running unit tests).
Also for each `release_item` we get the `type` (git, artifact or release) and then the `commit`, `artifact` and `release` for that `release_item` but only one of these will be non-null, dependent on the `type`.
If the `release_item.type == "release"` then we also need to perform this query for that release as we cannot recursively resolve queries using GraphQL... But we can get to that in the future.

```graphql
{
    release(_id: "123") {
        release_item {
            type
            release_entry {
                result
            }

            commit {
                id
                tag
                branch
                time
                repo {
                    name
                }
            }
            artifact {
                sha256
            }
            release {
                _id
            }
        }
    }
}
```

## TODO: Connect a `release` criteria with a `release_entry`

What's missing is the link... We define a bunch of `release_item` `criteria`, and we record a bunch of `release_entry(s)` but we how do we check those?!

## TODO: Missing features

There are some missing features that need to be implemented to support this:

1. A table can have `unique` fields, but currently a table cannot be uniquely identified by its `join(s)`... This needs to be added
2. Currently there is no way to specify in a `data {}` block what the save `policy` is... For instance, we may not want to create the data if it doesn't exist, instead it is only a reference.

## Data block `unique`

Adding `unique` to a field gives it a `UNIQUE` constraint in Postgres.
We should be able to specify `unique` on a join, to make the foreign key fields part of the unique constraint.

### Data block `policy`

The `policy` inside a data block can be:

1. `create_update` (default): creates or updates the data block
2. `reference`: does not create or update, but only pulls existing data and errors if a record does not exist
3. `create`: creates only, so if there is a conflict due to unique constraints, then it should error
