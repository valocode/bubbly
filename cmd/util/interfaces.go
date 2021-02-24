package util

import (
	"github.com/spf13/cobra"
)

type Options interface {
	Validate(cmd *cobra.Command) error
	Resolve() error
	Run() error
	Print()
}
