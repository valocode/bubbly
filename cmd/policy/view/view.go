package view

import (
	"fmt"

	"github.com/valocode/bubbly/client"
	"github.com/valocode/bubbly/env"
	"github.com/valocode/bubbly/store/api"

	"github.com/spf13/cobra"

	cmdutil "github.com/valocode/bubbly/cmd/util"
)

var (
	cmdLong = cmdutil.LongDesc(
		`
		View a Bubbly policy
		`,
	)

	cmdExamples = cmdutil.Examples(
		`
		# View a Bubbly policy

		bubbly policy view
		`,
	)
)

func New(bCtx *env.BubblyConfig) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "view name[:tag] [flags]",
		Short:   "View a Bubbly policy",
		Long:    cmdLong + "\n\n",
		Example: cmdExamples,
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			a, err := client.GetPolicy(bCtx, &api.ReleasePolicyGetRequest{
				Name: &name,
			})
			if err != nil {
				return err
			}
			fmt.Println("Name: " + *a.Name)
			fmt.Println("")
			fmt.Println("===")
			fmt.Println("")
			fmt.Println(*a.Module)
			return nil
		},
	}

	return cmd
}
