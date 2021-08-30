package create

import (
	"fmt"

	"github.com/valocode/bubbly/client"
	"github.com/valocode/bubbly/env"
	"github.com/valocode/bubbly/release"

	"github.com/spf13/cobra"

	cmdutil "github.com/valocode/bubbly/cmd/util"
)

var (
	cmdLong = cmdutil.LongDesc(
		`
		Create a bubbly release
		`,
	)

	cmdExamples = cmdutil.Examples(
		`
		# Create a bubbly release

		bubbly release create
		`,
	)
)

func New(bCtx *env.BubblyContext) *cobra.Command {

	cmd := &cobra.Command{
		Use:     "create [flags]",
		Short:   "Create a bubbly release",
		Long:    cmdLong + "\n\n",
		Example: cmdExamples,
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			req, err := release.CreateRelease(bCtx)
			if err != nil {
				return err
			}

			if err := client.CreateRelease(bCtx, req); err != nil {
				return err
			}
			fmt.Println("Release Created: ", *req.Release.Name)
			return nil
		},
	}
	return cmd
}
