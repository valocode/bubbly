package list

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/valocode/bubbly/bubbly"
	cmdutil "github.com/valocode/bubbly/cmd/util"
	"github.com/valocode/bubbly/env"
	"github.com/valocode/bubbly/util/normalise"
)

var (
	_       cmdutil.Options = (*options)(nil)
	cmdLong                 = normalise.LongDesc(`
		Create a release

		    $ bubbly release create

		`)

	cmdExample = normalise.Examples(`
		# Create a release
		bubbly release create
		`)
)

// options holds everything necessary to run the command.
// Flag values received to the command are loaded into this struct
type options struct {
	cmdutil.Options
	BubblyContext *env.BubblyContext
	Command       string
	Args          []string
	Release       *bubbly.ReleaseSpec
	// flags
}

// New creates a new cobra command
func New(bCtx *env.BubblyContext) *cobra.Command {
	o := &options{
		Command:       "create",
		BubblyContext: bCtx,
	}

	// cmd represents the apply command
	cmd := &cobra.Command{
		Use:     "create",
		Short:   "create a release",
		Long:    cmdLong + "\n\n",
		Example: cmdExample,
		Args:    cobra.NoArgs,
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
	return nil
}

// run runs the command over the validated options
func (o *options) run() error {
	release, err := bubbly.CreateRelease(o.BubblyContext)
	if err != nil {
		return err
	}
	o.Release = release
	return nil
}

// Print prints the successful outcome of the cmd
func (o *options) Print() {
	color.Green("Release successfully created!")

	fmt.Println("\n" + o.Release.String())
}
