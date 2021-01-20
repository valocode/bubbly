# Testing

This file describes the testing process for Bubbly. There are two sets of tests: "unit" and "integration". *Both* sets are run _locally_, but the latter depends on external components which need setting up prior to running the tests. There is a variety of test-related target in the root `Makefile`, providing extra information about the tests run.

## Unit tests

### Golang

Nothing fancy, if you want to rerun all the unit tests just do `go test ./...` or `make test-unit`

### NodeJS

TODO

## Integration tests

This is where things get a little more interesting as integration tests depend on external components: the Bubbly server and its Store (currently Postgres).

To run integration tests, the database and Bubbly server must be started first: `make dev`. This will take up the terminal and display logs from `postgres` and `bubbly` server containers. When these services are no longer needed, `<Ctrl+C>` would shut them down.

Then, in another terminal, run the integration tests with `make test-integration`.
