package apply

import (
	"fmt"
	"os"
	"path/filepath"

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
		Apply bubbly resources to a bubbly API server 
	`)

	applyExample = normalise.Examples(`
		# Apply the bubbly resources in the file ./main.bubbly
		bubbly apply -f ./main.bubbly

		# Apply the configuration in the directory ./resources
		bubbly apply -f ./resources
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

// NewCmdApply creates a new cobra.Command representing "bubbly apply"
func NewCmdApply(bCtx *env.BubblyContext) (*cobra.Command, *ApplyOptions) {
	o := &ApplyOptions{
		Command:       "apply",
		BubblyContext: bCtx,
	}

	// cmd represents the apply command
	cmd := &cobra.Command{
		Use:     "apply (-f (FILENAME | DIRECTORY)) [flags]",
		Short:   "Apply one or more bubbly resource to a bubbly agent",
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
		"filename or directory that contains the bubbly resources to apply")

	cmd.MarkFlagRequired("filename")

	return cmd, o
}

// Validate checks the ApplyOptions to see if there is sufficient information to
// run the command.
func (o *ApplyOptions) Validate(cmd *cobra.Command) error {
	if len(o.Args) != 0 {
		return cmdutil.UsageErrorf(cmd, "Unexpected args: %v", o.Args)
	}

	// check the file/directory is valid and fail fast if not
	if _, err := os.Stat(o.filename); err != nil {
		return fmt.Errorf(
			`failed to validate file/path to bubbly resources "%s": %w`,
			filepath.FromSlash(o.filename),
			err)

	}

	return nil
}

// Resolve resolves various ApplyOptions attributes from the provided arguments to cmd
func (o *ApplyOptions) Resolve() error {
	return nil
}

// Run runs the apply command over the validated ApplyOptions configuration
func (o *ApplyOptions) Run() error {
	if err := bubbly.Apply(o.BubblyContext, o.filename); err != nil {
		return fmt.Errorf("failed to apply configuration: %w", err)
	}
	return nil
}

// Print prints the successful outcome of applying the resource(s)
func (o *ApplyOptions) Print() {
	successString := fmt.Sprintf(
		`resource(s) at path/directory "%s" applied successfully`,
		filepath.FromSlash(o.filename))

	if o.BubblyContext.CLIConfig.Color {
		color.Green(successString)
	} else {
		fmt.Println(successString)
	}
}
