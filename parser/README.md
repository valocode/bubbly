# HCL Parser

This package contains the implementation for parsing HCL.

The design is such that the HCL is structured into modules, and there is a *"root"* module which is ultimately what should be resolved, and any sub modules are merely designed to keep things DRY (Don't Repeat Yourself) and enable reuse.

## Module

TODO: describe a Module, how it works and how to use it

## Scope

TODO: describe a Scope, how it works and how to use it

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

#### How to use a runtime defined body

Let's take the example of an `importer`:

```hcl
importer "importer_name" {
    type = "json"

    // The source block schema, is defined by the above attribute type...
    // Thus, if we have a json type, this block will look different to a rest
    // type.
    source {
        file = "./my-file.json"
    }
    // Using typeexpr to define a cty.Type, which in this case is a list of
    // strings.
    format = list(string)
}
```

So how do we parse this into a struct? Quite easily...

```go
type HCLBody struct {
    Importers []Importer `hcl:"importer,block"`
}

type Importer struct {
    Name string `hcl:",label"`
    Type string `hcl:"type,attr"`
    // The hcl tag "remain" means that it will not be decoded as part of the
    // decode procedure... i.e. it remains!
    // This means we can post-porcess it however we want :)
    SourceBody hcl.Body `hcl:"source,remain"`

    // hcl.Expressions are also not decoded, so that we can post-process them
    Format hcl.Expression `hcl:"format,attr"`
}

type JSONSource struct {
    File string `hcl:"file,attr"`
}
```

After we decode the `HCLBody`, we can do something like:

```go
hcl := HCLBody{}
// now decode the HCL body into hcl...
gohcl.Decode(body, &hcl)

for _, importer := range hcl.Importers {
    switch ty := importer.Type {
    case "json":
        jsonSource := JSONSource{}
        // now decode the importer.SourceBody into jsonSource...
        gohcl.Decode(importer.SourceBody, &jsonSource)
    default:
        panic(fmt.Sprintf("unknown importer type: %s", ty))
    }
}
```

### 2. Traversals - what are they

// TODO

### 3. EvalContext - how is it used

// TODO
