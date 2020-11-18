# HCL Parser

This package contains the implementation for parsing HCL.

The parser takes care of parsing `*.bubbly` files written in HCL (JSON not supported directly) to get the list of resources.
If you want to parse json, e.g. from the bubbly server, the `parser_json.go` file contains what you need.
The reason for handling these differently is because of HCL's underlying representation of the HCL AST, when it is parsed from HCL or JSON.

## Parser

The `Parser` type is the top-level type in this package and should be considered the entrypoint for calling programs.

## Scope

The `Scope` type is a wrapper around our `SymbolTable` implementation, which takes care of the EvalContext for resolving variables in HCL.

The main difference between `Parser` and `Scope` is that `Scope` does not maintain state of the HCL files to decode, it only maintains state of the EvalContext and is capable of producing *nested* EvalContexts, e.g. when a resource needs to reference its `self`, as that scope needs to be contained to just that resource.

## Useful knowledge

### 1. HCL structure

```hcl
my_block "block_label" {
    my_attr = "xyz"
}
```

The whole text in HCL, is a `body`, and the `body` content consists of `blocks` and `attributes`.
In this example, there is a `block` `type` called "my_block", which takes one `label` which is set to "block_label" in this case.
The "my_block" `block` contains a list of `labels` (in this case the list length is 1) and then a nested `body`.
And that nested `body` contains a list of `blocks` and `attributes`, but by our schema, it only contains one attribute called "my_attr".

```go
// HCLBody represents the entire body of HCL in our example
type HCLBody struct {
    // MyBlocks represents the list of my_block blocks in the body of HCLBody
    MyBlocks []struct {
        // Label is the label in the "my_block" block
        Label string `hcl:",label"`
        // MyAttr is the attribute inside the block "my_block"
        MyAttr cty.Value `hcl:"my_attr,attr"`
    } `hcl:"my_block,block"`
}
```

### 2. Traversals - what are they

A traversal in HCL can be seen as a *path* in HCL to a node in the Abstract Syntax Tree (AST).

For example, given the following HCL snippet, the traversal `myblock.block_label.my_attr` would point the the attribute `my_attr` which will be evaluated to a `cty.StringVal("xyz")`.

```hcl
my_block "block_label" {
    my_attr = "xyz"
}
```

### 3. EvalContext - how is it used

Extending our earlier example of HCL to include a `local` value, we could use the traversal `local.my_local.value` to reference the attribute `value` in the local `my_local`.

```hcl
local "my_local" {
    value = "xyz"
}

my_block "block_label" {
    my_attr = local.my_local.value
}
```

How does the expression evaluator in HCL know what the value of `local.my_local.value` is?
It needs to be in the `EvalContext.Variables` map in order to be "known", and this is the purpose of the EvalContext - to store variables/traversals as well as pre-defined functions.
