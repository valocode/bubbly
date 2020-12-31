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

#### HCL - cty library

To do anything semi-advanced with HCL it is absolutely necessary to know about the [cty package](https://github.com/zclconf/go-cty).

The `cty` package provides a dynamic type system (e.g. strings, numbers, objects) so that raw text values in HCL can be converted into some type in memory, that can easily be converted into a corresponding Go-native type or be manipulated using `cty`'s [stdlib](https://github.com/zclconf/go-cty/tree/master/cty/function/stdlib).

For example, consider this HCL:

```hcl
understand_cty {
    name = "juggernaut"
    strength = 10000

    friends = xmen.bad_guys
}
```

How will the text `"juggernaut"` be interpreted, and what will the `type` of `name` be? `cty` takes care of this when parsing, decoding and evaluating HCL.
Similarly with the raw text `10000`, `cty` will interpret this as a `Number`.

Now what about `friends = xmen.bad_guys`? In this case, `friends` should receive a list. But what is the type of `xmen`?
We are trying to access an attribute `bad_guys` which means `xmen` is expected to be an `Object`.

Anyway, this was just to get you thinking and understand what happens under the hood, so ignore this and go read the real documentation, such as the [`cty` types](https://github.com/zclconf/go-cty/blob/master/docs/types.md).
