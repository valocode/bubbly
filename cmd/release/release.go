package schema

import (
	"github.com/spf13/cobra"

	createCmd "github.com/valocode/bubbly/cmd/release/create"
	evalCmd "github.com/valocode/bubbly/cmd/release/eval"
	listCmd "github.com/valocode/bubbly/cmd/release/list"
	viewCmd "github.com/valocode/bubbly/cmd/release/view"
	"github.com/valocode/bubbly/env"
)

// New creates a new cobra.Command representing "bubbly release"
func New(bCtx *env.BubblyContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "release <command>",
		Short: "Manage your bubbly releases",
		Long:  `Manage your bubbly releases`,
	}

	cmd.AddCommand(createCmd.New(bCtx))
	cmd.AddCommand(evalCmd.New(bCtx))
	cmd.AddCommand(listCmd.New(bCtx))
	cmd.AddCommand(viewCmd.New(bCtx))

	return cmd
}
