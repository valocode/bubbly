package list

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
		List Bubbly policies
		`,
	)

	cmdExamples = cmdutil.Examples(
		`
		# List Bubbly policies

		bubbly policy list
		`,
	)
)

func New(bCtx *env.BubblyContext) *cobra.Command {

	var (
		name string
	)
	cmd := &cobra.Command{
		Use:     "list [flags]",
		Short:   "List Bubbly policies",
		Long:    cmdLong + "\n\n",
		Example: cmdExamples,
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			resp, err := client.GetPolicies(bCtx, &api.ReleasePolicyGetRequest{
				Name: name,
			})
			if err != nil {
				return err
			}
			if err != nil {
				return fmt.Errorf("getting policies: %w", err)
			}

			for _, p := range resp.Policies {
				fmt.Println(*p.Name)
			}
			return nil
		},
	}

	f := cmd.Flags()
	f.StringVarP(&name, "name", "n", "", "The name of the policy to filter by")

	return cmd
}
