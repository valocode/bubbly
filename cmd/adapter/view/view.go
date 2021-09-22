package view

import (
	"fmt"

	"github.com/valocode/bubbly/adapter"
	"github.com/valocode/bubbly/client"
	"github.com/valocode/bubbly/config"
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

func New(bCtx *config.BubblyConfig) *cobra.Command {
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
			resp, err := client.GetAdapters(bCtx, &api.AdapterGetRequest{
				Name: name,
				Tag:  tag,
			})
			if err != nil {
				return err
			}
			if len(resp.Adapters) == 0 {
				fmt.Println("No adapters found")
				return nil
			}
			if len(resp.Adapters) > 1 {
				fmt.Println("More than one adapter found with id: ", args[0])
				return nil
			}
			a := resp.Adapters[0]
			fmt.Println("Name: " + *a.Name)
			fmt.Println("Tag: " + *a.Tag)
			fmt.Println("")
			fmt.Println("===")
			fmt.Println("")
			fmt.Println(*a.Module)
			return nil
		},
	}

	return cmd
}
