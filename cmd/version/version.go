package version

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/valocode/bubbly/env"
)

// New creates a new Cobra command
func New(bCtx *env.BubblyContext) *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Show the Bubbly version information",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("bubbly version", bCtx.Version.Tag, bCtx.Version.SHA1)
		},
	}
}
