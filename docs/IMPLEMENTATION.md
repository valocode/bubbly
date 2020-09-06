# Bubbly Implementation

This document provides details on the implementation of Bubbly.

## Overview

Considerable inspiration has been taken from HashiCorp's approach to developing software and a similar model has been adopted for Bubbly.

Thus, a couple of high-level points:

1. Bubbly will be written in Go and will be shipped as a single binary
2. All configurations will be provided using the [HashiCorp Configuration Language (HCL)](https://github.com/hashicorp/hcl)
3. HashiCorp's [go-memdb](https://github.com/hashicorp/go-memdb) library will be used as an in-memory schemaful database

A lot of inspiration will be drawn from tools like [Vault](http://github.com/hashicorp/vault), [Nomad](http://github.com/hashicorp/nomad) and [Terraform](http://github.com/hashicorp/terraform) to drive the direction of Bubbly.
