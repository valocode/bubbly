# Commands

## Implementing a new command

Every command requires a few things:

### 1. an equivalent `<command>Options` struct

This struct inherits from the `cmdutil.Options` interface:

```golang

type Options interface {
 Validate(cmd *cobra.Command) error
 Resolve(cmd *cobra.Command) error
 Run() error
 Print(cmd *cobra.Command)
}

```

And therefore requires the following functions to be implemented:

- `Validate(cmd *cobra.Command)` checks the inputs to the command are valid and that there is sufficient information run the command.
- `Resolve(cmd *cobra.Command)` populates `<command>Options` attributes from the provided arguments to `cmd`
- `Run()` Run runs the command against the resolved `<command>Options`. Command results should be populated into the `<command>Options.Results` attribute
- `Print(cmd *cobra.Command)` prints outputs from `<command>Options.Results`

### 2. A `NewCmd<command>` function

This is the "heart" of the Command and creates a new `cobra.Command` representing `bubbly <command>`. Here is where you should instantiate your `<command>Options` struct and create your `cobra.Command`.

### 3.  An `init()` function

This binds the command to the `rootCmd`.

# Testing

[gock](https://github.com/h2non/gock) lets us create mock responses to HTTP routes that don't yet exist on the Bubbly server. This allows us to:

1. implement end-to-end tests for command functionality for which there is not yet backend support
2. prototype data formats being sent to and from the bubbly server. For example, `client.DescribeResourceReturn`
