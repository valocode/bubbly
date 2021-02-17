package util

import (
	"fmt"

	"github.com/spf13/cobra"
)

// UsageErrorf can be used as a generic 'you've made a mistake, look at the help documentation'.
func UsageErrorf(cmd *cobra.Command, format string, args ...interface{}) error {
	msg := fmt.Sprintf(format, args...)
	return fmt.Errorf("%s\nSee '%s -h' for help and examples", msg, cmd.CommandPath())
}
