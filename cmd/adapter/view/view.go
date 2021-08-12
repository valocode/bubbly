package view

import (
	"fmt"

	"github.com/valocode/bubbly/adapter"
	"github.com/valocode/bubbly/client"
	"github.com/valocode/bubbly/env"
	"github.com/valocode/bubbly/store/api"

	"github.com/spf13/cobra"

	cmdutil "github.com/valocode/bubbly/cmd/util"
)

var (
	cmdLong = cmdutil.LongDesc(
		`
		View a Bubbly adapter
		`,
	)

	cmdExamples = cmdutil.Examples(
		`
		# View a Bubbly adapter

		bubbly adapter view
		`,
	)
)

func New(bCtx *env.BubblyContext) *cobra.Command {

	cmd := &cobra.Command{
		Use:     "view name[:tag] [flags]",
		Short:   "View a Bubbly adapter",
		Long:    cmdLong + "\n\n",
		Example: cmdExamples,
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name, tag, err := adapter.ParseAdpaterID(args[0])
			if err != nil {
				return err
			}
			a, err := client.GetAdapter(bCtx, &api.AdapterGetRequest{
				Name: &name,
				Tag:  &tag,
			})
			if err != nil {
				return err
			}
			results, err := a.Results.SpecBytes()
			if err != nil {
				return err
			}
			fmt.Println("Name: " + a.Name)
			fmt.Println("Tag: " + a.TagOrDefault())
			fmt.Println("Type: " + a.Operation.Type)
			fmt.Println("Results:")
			fmt.Printf("%s\n", results)
			return nil
		},
	}

	return cmd
}
