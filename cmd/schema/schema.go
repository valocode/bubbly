package schema

import (
	"github.com/spf13/cobra"

	schemaApplyCmd "github.com/verifa/bubbly/cmd/schema/apply"
	"github.com/verifa/bubbly/env"
)

// NewCmdGet creates a new cobra.Command representing "bubbly schema"
func NewCmdSchema(bCtx *env.BubblyContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "schema <command>",
		Short: "Manage your bubbly schema",
		Long:  `Manage your bubbly schema`,
	}

	schemaApplyCmd, _ := schemaApplyCmd.NewCmdApply(bCtx)
	cmd.AddCommand(schemaApplyCmd)

	return cmd
}
