package util

import (
	"fmt"

	"github.com/spf13/cobra"
	normalise "github.com/verifa/bubbly/util/normalise"
)

// SuggestBubblyResources returns a suggestion to use the "api-resources" command
// to retrieve a supported list of resources
// TODO: the bubbly server should expose an endpoint for grabbing a list of possible bubbly resource supported by 'bubbly describe <resource>'
func SuggestBubblyResources() string {
	return normalise.LongDesc(fmt.Sprintf("Use 'bubbly api-resources' for a complete list of supported resources."))
}

// UsageErrorf can be used as a generic 'you've made a mistake, look at the help documentation'.
func UsageErrorf(cmd *cobra.Command, format string, args ...interface{}) error {
	msg := fmt.Sprintf(format, args...)
	return fmt.Errorf("%s\nSee '%s -h' for help and examples", msg, cmd.CommandPath())
}
