# API

This package contains the api, which defines the resources and structure of the HCL, as well as other things (Work In Progress).

## Structure

There are three main levels of packages:

```bash
api
# The API directory contains the "layer" that is used by other packages to
# interact with the API, such as creating a new resource from a resource block
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

## Packages

### api

The main functionality of the `api` package is to "tie together" the core and the versioned types.

E.g. one of the main methods is `NewResource(*ResourceBlock) *Resource` which returns a `Resource` based on the provided `ResourceBlock`.

The returned `Resource` implements the `Resource` interface but is specific to the `ResourceKind` and `ResourceVersion` specified in the `ResourceBlock`.
E.g. if it specifies `api_version: "v1"` of the resource kind `importer`, then a new instance of `v1.Importer` is returned.

### core

The `core` package defines the very important `ResourceBlock` type which describes the shape of a `resource {...}` block in HCL.

It also defines the interfaces for the different resources, like `importer`, `translator`, `publish`, `pipeline`, etc.

### v1

As mentioned earlier, `core` is coupled to the version of bubbly which is being used.
When a user defines a resource, they specify an `api_version` which specifies the *versioned* type that they want, and those versioned types are defined in a package by the version number.

Thus `v1` is the package containing all the definitions of all api version `v1` resource kinds.

## Resources

Every object in bubbly will be a different kind of `Resource`, including `importer`, `translator`, `publish`, `pipeline`, etc.

Please see more information in [DESIGN.md](../docs/DESIGN.md)
