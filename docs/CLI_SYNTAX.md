# How we document our command line syntax

## Placeholder values

Use full word capitilisation to represent a value the user must replace.

_example:_
`bubbly get (TYPE NAME_PREFIX) [flags]`
Replace `TYPE NAME_PREFIX` with the Bubbly resource type and name.

Note: Should `NAME_PREFIX` actually be `[name]`, since it is optional? And lowercase?

## Optional arguments

Place optional arguments in square brackets.

_example:_
`bubbly version [flags]`
Flags are optional.

## Required mutually exclusive arguments

Place required mutually exclusive required arguments inside `()`, separate arguments with vertical bars.

_example:_
`bubbly get ((TYPE NAME_PREFIX) | GENERIC)`

## Repeatable arguments

Ellipsis represent arguments that can appear multiple times.

_example:_
`bubbly explain TYPE[.<field-name>...] [flags]`

## Variable naming

For multi-word variables use dash-case (all lower case with words separated by dashes)

_example:_
`bubbly explain TYPE[.<field-name>...] [flags]`
