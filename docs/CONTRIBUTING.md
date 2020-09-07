# Contributing Guide for Bubbly

## Background Learning

### HashiCorp Configuration Language (HCL)

The basics of HCL are so easy there is not much to learn...
However there is some useful introductory stuff on the Terraform docs: <https://www.terraform.io/docs/configuration/syntax.html>

For what we are going to be doing, we need to dig into the more advanced stuff.
A great starting place is a YouTube video highlighting some of the more advanced features:
<https://youtu.be/YZcDTTc_VJM>

To accompany this video is the repository with the source code which can be found here:
<https://github.com/RussellRollins/pet-sounds>

### Radix Trees - memdb

The internal data structure used by Bubbly is implemented using radix trees, using the module `go-memdb`.

To understand Radix trees, this presentation is very useful: [radix-trees](./files/Radix-Txn-MemDB.pdf)


