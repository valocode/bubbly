# Architecture

This document gices a rough overview of the bubbly architecture.

Most notable directory mentions are:

1. [ent](./ent) contains everything related to the data access layer and database schema (implemented using [entgo](https://entgo.io/))
   1. The most interesting thing for someone is probably the schema itself, which can be found in the [schema dir](./ent/schema)
2. [adapter](./adapter) contains the implementation for the bubbly adapters
3. [policy](./policy) contains the implementation for the bubbly policies
4. [monitor](./monitor) contains the implementation for *monitoring* external services (such as fetching SPDX licenses, vulnerabilities, etc.)
5. [ui](./ui) contains the frontend built using [SvelteKit](https://kit.svelte.dev/) and TailwindCSS
6. [docs](./docs) contains the documentation for bubbly which gets deployed to [https://docs.bubbly.dev](https://docs.bubbly.dev)
