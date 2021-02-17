package cmd

import (
	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"
)

var HelpTopics = map[string]map[string]string{
	"environment": {
		"alias": "env",
		"short": "Environment variables that can be used with bubbly",
		"long": heredoc.Doc(`
			# general

			DEBUG: set to true to enable debug logging. Default: false

			# bubbly API server

			BUBBLY_PROTOCOL: specify the bubbly API server protocol (http/https). Default: http

			BUBBLY_HOST: specify the bubbly API server host. Default: localhost
			
			BUBBLY_PORT: specify the bubbly API  server port. Default: 8111

			# bubbly store

			## generic

			BUBBLY_STORE_PROVIDER: specify the database provider of the bubbly store. Default: postgres

			## postgres

			POSTGRES_ADDR: specify the address of the postgres instance. Default: postgres:5432

			POSTGRES_USER: specify the postgres user. Default: postgres

			POSTGRES_PASSWORD: specify the password of the postgres user. Default: postgres

			POSTGRES_DATABASE: specify the postgres database to use. Default: bubbly

			## cockroachdb

			COCKROACH_ADDR: specify the address of the cockroachdb instance. Default: cockroachdb:26257

			COCKROACH_USER: specify the cockroachdb user. Default: root

			COCKROACH_PASSWORD: specify the password of the cockroachdb user. Default admin

			COCKROACH_DATABASE: specify the cockroachdb database to use. Default defaultdb

			# bubbly agent

			## generic

			AGENT_DEPLOYMENT_TYPE: specify the type of bubbly agent deployment. Default: single

			AGENT_API_SERVER_TOGGLE: specify whether to run the API server as a part of the agent. Default: false

			AGENT_DATA_STORE_TOGGLE: specify whether to run the data store as a part of the agent. Default: false

			AGENT_WORKER_TOGGLE: specify whether to run the worker as a part of the agent. Default: false

			AGENT_NATS_SERVER_TOGGLE: specify whether to run a NATS Server as a part of the agent. Default: false

			## NATS Server

			NATS_SERVER_ADDR: specify the address of the NATS server. Default: localhost:4223

			NATS_SERVER_HTTP_PORT: specify the http port address of the NATS server. Default: 8222

			NATS_SERVER_PORT: specify the port address of the NATS server. Default: 4223
		`),
	},
}

func NewHelpTopic(topic string) *cobra.Command {
	cmd := &cobra.Command{
		Aliases: []string{HelpTopics[topic]["alias"]},
		Use:     topic,
		Short:   HelpTopics[topic]["short"],
		Long:    HelpTopics[topic]["long"],
		Hidden:  false,
	}

	return cmd
}
