package view

import (
	"fmt"
	"strings"

	"github.com/valocode/bubbly/client"
	"github.com/valocode/bubbly/config"
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

func New(bCtx *config.BubblyConfig) *cobra.Command {
	var withAffects bool

	cmd := &cobra.Command{
		Use:     "view name [flags]",
		Short:   "View a Bubbly policy",
		Long:    cmdLong + "\n\n",
		Example: cmdExamples,
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			resp, err := client.GetPolicies(bCtx, &api.ReleasePolicyGetRequest{
				Name:        name,
				WithAffects: withAffects,
			})
			if err != nil {
				return err
			}
			if err != nil {
				return fmt.Errorf("getting policies: %w", err)
			}
			if len(resp.Policies) == 0 {
				fmt.Println("No policies with name: " + name)
				return nil
			}
			p := resp.Policies[0]
			fmt.Println("Name: " + *p.Name)
			fmt.Println("")
			if withAffects {
				fmt.Println("Projects:", strings.Join(p.Affects.Projects, ", "))
				fmt.Println("Repos:", strings.Join(p.Affects.Repos, ", "))
				fmt.Println("")
			}
			fmt.Println("===")
			fmt.Println("")
			fmt.Println(*p.Module)
			return nil
		},
	}
	f := cmd.Flags()
	f.BoolVar(&withAffects, "affects", false, "Include projects, repos that the policy applies to")

	return cmd
}
