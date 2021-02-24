# Commands

## Implementing a new command

Every command requires a few things:

### 1. an equivalent `<command>Options` struct

This struct embeds the `cmdutil.Options` interface:

```go
type Options interface {
    Validate(cmd *cobra.Command) error
    Resolve() error
    Run() error
    Print()
}
```

And therefore requires the following functions to be implemented:

- `Validate(cmd *cobra.Command)` checks the inputs to the command are valid 
  and that there is sufficient information run the command.
- `Resolve()` uses flag values in `<command>Options` to resolve other fields 
  needed downstream by the `Run()`
- `Run()` Run "runs" the command. That is, it calls logic external to the `cmd` 
  package using the `<command>Options` to fulfil the core purpose of the 
  command.
- `Print()` prints the final (successful) output from running the command.

### 2. A `NewCmd<command>` function

This is the heart of any Command. It creates a new `cobra.Command` 
representing `bubbly <command>`. Here is where you should instantiate your 
`<command>Options` struct and create your `cobra.Command`.

### 3.  An explicit bind to the `rootCmd` in `cmd/root.go` 

This binds the command to the `rootCmd`.