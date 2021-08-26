package push

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/valocode/bubbly/adapter"
	"github.com/valocode/bubbly/client"
	"github.com/valocode/bubbly/ent"
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

	var (
		name string
		tag  string
	)
	cmd := &cobra.Command{
		Use:     "push <adapter-file> [flags]",
		Short:   "Push a Bubbly adapter",
		Long:    cmdLong + "\n\n",
		Example: cmdExamples,
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			filename := args[0]
			b, err := os.ReadFile(filename)
			if err != nil {
				return err
			}
			if name == "" {
				fname := filepath.Base(filename)
				name = strings.TrimSuffix(fname, filepath.Ext(fname))
			}
			if err := client.SaveAdapter(bCtx, &api.AdapterSaveRequest{
				Adapter: ent.NewAdapterModelCreate().
					SetName(name).
					SetTag(tag).
					SetModule(string(b)),
			}); err != nil {
				return err
			}
			return nil
		},
	}

	f := cmd.PersistentFlags()
	f.StringVar(
		&name,
		"name",
		"",
		`Provide the name of the adapter (default is filename without extension)`,
	)
	f.StringVar(
		&tag,
		"tag",
		adapter.DefaultTag,
		fmt.Sprintf(`Provide the tag to apply to the adapter (default is "%s")`, adapter.DefaultTag),
	)
	return cmd
}
