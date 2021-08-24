package demo

import (
	"fmt"
	"log"

	"github.com/valocode/bubbly/env"
	"github.com/valocode/bubbly/server"
	"github.com/valocode/bubbly/store"
	"github.com/valocode/bubbly/test"

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

			fmt.Println("Creating dummy releases...")
			data := test.CreateDummyData()

			for _, repos := range data {
				for _, release := range repos.Releases {
					if _, err := store.CreateRelease(release.Release); err != nil {
						log.Fatalf("creating release %s: %s", *release.Release.Release.Name, err.Error())
					}
					for _, scan := range release.CodeScans {
						if _, err := store.SaveCodeScan(scan); err != nil {
							log.Fatalf("saving code scan: %s", err.Error())
						}
					}
					for _, run := range release.TestRuns {
						if _, err := store.SaveTestRun(run); err != nil {
							log.Fatalf("saving test run: %s", err.Error())
						}
					}
				}
			}
			fmt.Println("Done!")
			fmt.Println("")

			fmt.Println("Evaluating releases...")
			if err := test.FailSomeRandomReleases(store); err != nil {
				log.Fatal("evaluating releases: ", err)
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
