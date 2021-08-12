package server

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
		Starts the bubbly server. The server exposes the API (REST and GraphQL)
		and initializes the store which connects to the specified database

			$ bubbly server
		`,
	)

	cmdExamples = cmdutil.Examples(
		`
		# Starts the bubbly server

		bubbly server
		`,
	)
)

func New(bCtx *env.BubblyContext) *cobra.Command {

	cmd := &cobra.Command{
		Use:     "server [flags]",
		Short:   "Start the bubbly server",
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
			fmt.Printf("Starting HTTP server on %s:%s\n", bCtx.ServerConfig.Host, bCtx.ServerConfig.Port)
			if err := server.ListenAndServe(bCtx, store); err != nil {
				return err
			}
			return nil
		},
	}

	f := cmd.Flags()

	// agent.AgentDeploymentType's underlying type is a string,
	// so we cast to *string on bind
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
	f.StringVar(
		&bCtx.StoreConfig.PostgresAddr,
		"postgres-addr",
		bCtx.StoreConfig.PostgresAddr,
		"postgres address for the data store",
	)
	f.StringVar(
		&bCtx.StoreConfig.PostgresUser,
		"postgres-username",
		bCtx.StoreConfig.PostgresUser,
		"postgres username for the data store",
	)
	f.StringVar(
		&bCtx.StoreConfig.PostgresPassword,
		"postgres-password",
		bCtx.StoreConfig.PostgresPassword,
		"postgres password for the data store",
	)
	f.StringVar(
		&bCtx.StoreConfig.PostgresDatabase,
		"postgres-database",
		bCtx.StoreConfig.PostgresDatabase,
		"postgres database for the data store",
	)

	return cmd
}