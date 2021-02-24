package apply

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/verifa/bubbly/bubbly"
	cmdutil "github.com/verifa/bubbly/cmd/util"
	"github.com/verifa/bubbly/env"
	"github.com/verifa/bubbly/util/normalise"
)

var (
	_         cmdutil.Options = (*ApplyOptions)(nil)
	applyLong                 = normalise.LongDesc(`
		Apply a bubbly schema

		    $ bubbly schema apply -f FILENAME

		`)

	applyExample = normalise.Examples(`
		# Apply a bubbly schema located in a specific file
		bubbly schema apply -f ./schema.bubbly
		`)
)

// ApplyOptions holds everything necessary to run the command.
// Flag values received to the command are loaded into this struct
type ApplyOptions struct {
	cmdutil.Options
	BubblyContext *env.BubblyContext
	Command       string
	Args          []string

	// flags
	filename string
}

// NewCmdApply creates a new cobra.Command representing "schema apply"
func NewCmdApply(bCtx *env.BubblyContext) (*cobra.Command, *ApplyOptions) {
	o := &ApplyOptions{
		Command:       "apply",
		BubblyContext: bCtx,
	}

	// cmd represents the apply command
	cmd := &cobra.Command{
		Use:     "apply -f FILENAME",
		Short:   "apply a bubbly schema",
		Long:    applyLong + "\n\n",
		Example: applyExample,
		RunE: func(cmd *cobra.Command, args []string) error {
			o.Args = args

			validationError := o.Validate(cmd)

			if validationError != nil {
				return validationError
			}

			resolveError := o.Resolve()

			if resolveError != nil {
				return resolveError
			}

			runError := o.Run()

			if runError != nil {
				return runError
			}

			o.Print()

			return nil
		},
	}

	f := cmd.Flags()

	f.StringVarP(&o.filename,
		"filename",
		"f",
		"",
		"filename that contains the .bubbly schema file to apply")

	cmd.MarkFlagRequired("filename")

	return cmd, o
}

// Validate checks the ApplyOptions to see if there is sufficient information run the command.
func (o *ApplyOptions) Validate(cmd *cobra.Command) error {
	fi, err := os.Stat(o.filename)
	if err != nil {
		return fmt.Errorf(`cannot read from filename "%s"`, o.filename)
	}

	// restrict apply of a schema to a single file to avoid the case
	// of multiple schema file definitions
	if fi.IsDir() {
		return fmt.Errorf("cannot apply a directory")
	}

	return nil
}

// Resolve resolves various ApplyOptions attributes from the provided arguments to cmd
func (o *ApplyOptions) Resolve() error {
	return nil
}

// Run runs the apply command over the validated ApplyOptions configuration
func (o *ApplyOptions) Run() error {
	err := bubbly.ApplySchema(o.BubblyContext, o.filename)

	if err != nil {
		return fmt.Errorf("failed to apply schema: %w", err)
	}
	return nil
}

// Print prints the successful outcome of applying a schema
func (o *ApplyOptions) Print() {
	successString := fmt.Sprintf(
		`schema at path "%s" successfully applied`,
		o.filename)

	if o.BubblyContext.CLIConfig.Color {
		color.Green(successString)
	} else {
		fmt.Println(successString)
	}
}
