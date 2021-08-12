package push

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
		Push a Bubbly adapter to a remote server
		`,
	)

	cmdExamples = cmdutil.Examples(
		`
		# Push a Bubbly adapter

		bubbly adapter push
		`,
	)
)

func New(bCtx *env.BubblyContext) *cobra.Command {

	cmd := &cobra.Command{
		Use:     "push <adapter-file> [flags]",
		Short:   "Push a Bubbly adapter",
		Long:    cmdLong + "\n\n",
		Example: cmdExamples,
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			filename := args[0]
			adapter, err := adapter.FromFile(filename)
			if err != nil {
				return err
			}
			model, err := adapter.Model()
			if err != nil {
				return fmt.Errorf("preparing adapter for push: %w", err)
			}
			fmt.Println("results_type: ", *model.ResultsType)
			if err := client.SaveAdapter(bCtx, &api.AdapterSaveRequest{
				AdapterModel: model,
			}); err != nil {
				return err
			}
			return nil
		},
	}
	return cmd
}
