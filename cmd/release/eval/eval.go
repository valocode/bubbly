package list

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/valocode/bubbly/bubbly"
	"github.com/valocode/bubbly/cmd/util"
	cmdutil "github.com/valocode/bubbly/cmd/util"
	"github.com/valocode/bubbly/env"
)

var (
	_       cmdutil.Options = (*options)(nil)
	cmdLong                 = util.LongDesc(`
		Evaluate a release criteria and log a release entry with the result.
		The criteria needs to exist in the release and the result will be attached to a release entry which is created as part of this command.

		    $ bubbly release eval CRITERIA

		`)

	cmdExample = util.Examples(`
		# Evaluate the release criteria called unit_test which is defined in the release.
		# A release entry with the name unit_test will be created.
		bubbly release eval unit_test
		`)
)

// options holds everything necessary to run the command.
// Flag values received to the command are loaded into this struct
type options struct {
	cmdutil.Options
	bCtx    *env.BubblyContext
	Command string
	Args    []string
	Release *bubbly.ReleaseSpec

	// flags
	Criteria string
}

// New creates a new cobra command
func New(bCtx *env.BubblyContext) *cobra.Command {
	o := &options{
		Command: "eval",
		bCtx:    bCtx,
	}

	// cmd represents the apply command
	cmd := &cobra.Command{
		Use:     "eval",
		Short:   "evaluate a release criteria",
		Long:    cmdLong + "\n\n",
		Example: cmdExample,
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			o.Args = args

			if err := o.validate(cmd); err != nil {
				return err
			}
			if err := o.resolve(); err != nil {
				return err
			}
			if err := o.run(); err != nil {
				return err
			}

			o.Print()
			return nil
		},
	}

	return cmd
}

// validate checks the cmd options
func (o *options) validate(cmd *cobra.Command) error {
	// Nothing to do
	return nil
}

// resolve resolves args for the command
func (o *options) resolve() error {
	o.Criteria = o.Args[0]
	return nil
}

// run runs the command over the validated options
func (o *options) run() error {
	release, err := bubbly.EvalReleaseCriteria(o.bCtx, o.Criteria)
	if err != nil {
		return err
	}
	o.Release = release
	return nil
}

// Print prints the successful outcome of the cmd
func (o *options) Print() {
	if o.bCtx.CLIConfig.Color {
		color.Green("Release entry successfully logged!")
	} else {
		fmt.Println("Release entry successfully logged!")
	}
}
