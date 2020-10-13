# api

This package contains the api, which defines the resources and structure of the HCL, as well as other things (Work In Progress).

## Structure

There are three main levels of packages:

```bash
api
# The API directory contains the "layer" that is used by other packages to
# interact with the API
├── core
    # Core contains the definitions which are versioned (coupled) with the
    # version of bubbly.
    # These should ideally never change but when they do could cause breaking
    # changes.
    # The core package defines types such as the interface for Resources and
    # the ResourceBlock struct.
├── v1
    # The versioned packages provide the specific versioned API types, such as
    # the different resources (e.g. Importer, Translator, Upload)
└── v2 # just another version
```

## api

The main functionality of the `api` package is to "tie together" the core and the versioned types.

E.g. one of the main methods is `NewResource(*ResourceBlock) *Resource` which returns a `Resource` based on the provided `ResourceBlock`.

The returned `Resource` implements the `Resource` interface but is specific to the `ResourceKind` and `ResourceVersion` specified in the `ResourceBlock`.
E.g. if it specifies `apiVersion: "v1"` of the resource kind `importer`, then a new instance of `v1.Importer` is returned.

## core

// TODO

## v1

// TODO
