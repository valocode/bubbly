package view

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/valocode/bubbly/bubbly"
	"github.com/valocode/bubbly/bubbly/builtin"
	"github.com/valocode/bubbly/cmd/util"
	cmdutil "github.com/valocode/bubbly/cmd/util"
	"github.com/valocode/bubbly/env"
)

var (
	_       cmdutil.Options = (*options)(nil)
	cmdLong                 = util.LongDesc(`
		View a release

		    $ bubbly release view

		`)

	cmdExample = util.Examples(`
		# View a release
		bubbly release view
		`)
)

// options holds everything necessary to run the command.
// Flag values received to the command are loaded into this struct
type options struct {
	cmdutil.Options
	bCtx    *env.BubblyContext
	Command string
	Args    []string
	Release *builtin.Release

	// flags
}

// New creates a new cobra command
func New(bCtx *env.BubblyContext) *cobra.Command {
	o := &options{
		Command: "view",
		bCtx:    bCtx,
	}

	// cmd represents the apply command
	cmd := &cobra.Command{
		Use:     "view",
		Short:   "view a release",
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
	release, err := bubbly.GetRelease(o.bCtx)
	if err != nil {
		return err
	}
	o.Release = release
	return nil
}

// Print prints the successful outcome of the cmd
func (o *options) Print() {
	fmt.Println("Project: " + o.Release.Project.Id)
	fmt.Println("Name: " + o.Release.Name)
	fmt.Println("Version: " + o.Release.Version)
	fmt.Println("Type: " + o.Release.ReleaseItem[0].Type)
	fmt.Println("Status: " + builtin.ReleaseStatusByStages(*o.Release))
}
