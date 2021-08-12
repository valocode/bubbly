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

			// if !skipAll {
			// 	if !skipCVE {
			// 		fmt.Println("Fetching CVEs from NVD... This will take a few seconds...")
			// 		if err := test.SaveCVEData(client); err != nil {
			// 			log.Fatal("loading CVEs: ", err)
			// 		}
			// 		fmt.Println("Done!")
			// 		fmt.Println("")
			// 	}
			// 	if !skipSPDX {
			// 		fmt.Println("Fetching SPDX licenses from GitHub...")
			// 		if err := test.SaveSPDXData(client); err != nil {
			// 			log.Fatal("loading SPDX: ", err)
			// 		}
			// 		fmt.Println("Done!")
			// 		fmt.Println("")

			// 	}

			fmt.Println("Creating dummy releases...")
			if err := test.CreateDummyData(store); err != nil {
				log.Fatal("creating dummy data: ", err)
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
			if err := server.ListenAndServe(bCtx, store); err != nil {
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
