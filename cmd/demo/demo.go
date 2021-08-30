package demo

import (
	"fmt"

	"github.com/valocode/bubbly/env"
	"github.com/valocode/bubbly/server"
	"github.com/valocode/bubbly/store"

	"github.com/spf13/cobra"

	cmdutil "github.com/valocode/bubbly/cmd/util"
)

var (
	cmdLong = cmdutil.LongDesc(
		`
		Starts the bubbly server in demo mode. The same configuration can be
		achieved using the bubbly server command, making this a convenience
		command for those looking to explore bubbly.

		WARNING: the data is not persisted... you have been warned!

			$ bubbly demo
		`,
	)

	cmdExamples = cmdutil.Examples(
		`
		# Starts the bubbly server in demo mode

		bubbly demo
		`,
	)
)

func New(bCtx *env.BubblyContext) *cobra.Command {

	cmd := &cobra.Command{
		Use:     "demo [flags]",
		Short:   "Start the bubbly server in demo mode",
		Long:    cmdLong + "\n\n",
		Example: cmdExamples,
		RunE: func(cmd *cobra.Command, args []string) error {

			fmt.Println("Initializing store...")
			store, err := store.New(bCtx)
			if err != nil {
				return fmt.Errorf("error initializing store: %w", err)
			}
			fmt.Println("Store initialized: ", bCtx.StoreConfig.Provider.String())

			fmt.Println("")
			fmt.Println("Populating the store with dummy data...")

			fmt.Println("Creating dummy data...")
			if err := store.PopulateStoreWithPolicies("."); err != nil {
				return err
			}
			if err := store.PopulateStoreWithDummyData(); err != nil {
				return err
			}
			fmt.Println("Done!")
			fmt.Println("")

			fmt.Printf("Starting HTTP server on %s:%s\n", bCtx.ServerConfig.Host, bCtx.ServerConfig.Port)
			if err := server.NewWithStore(bCtx, store).Start(); err != nil {
				return err
			}
			return nil
		},
	}

	f := cmd.Flags()

	f.StringVar(
		&bCtx.ServerConfig.Host,
		"host",
		bCtx.ServerConfig.Host,
		"host name for running the server on",
	)
	f.StringVarP(
		&bCtx.ServerConfig.Port,
		"port", "p",
		bCtx.ServerConfig.Port,
		"port to run the server on",
	)

	return cmd
}
