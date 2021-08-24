package release

import (
	"github.com/valocode/bubbly/cmd/release/create"
	"github.com/valocode/bubbly/env"

	"github.com/spf13/cobra"
)

func New(bCtx *env.BubblyContext) *cobra.Command {

	cmd := &cobra.Command{
		Use:   "release <command>",
		Short: "Manage bubbly releases",
		Long:  `Manage bubbly releases`,
	}

	cmd.AddCommand(create.New(bCtx))
	return cmd
}
