package store

import (
	"github.com/graphql-go/graphql"
	"github.com/verifa/bubbly/api/core"
	"github.com/verifa/bubbly/env"
)

// Provider provides an interface for persisting readiness data.
type provider interface {
	Apply(*bubblySchema) error
	Save(*env.BubblyContext, *bubblySchema, dataTree) error
	ResolveQuery(graph *schemaGraph, params graphql.ResolveParams) (interface{}, error)
	HasTable(core.Table) (bool, error)
	GenerateMigration(ctx *env.BubblyContext, cl Changelog) (migration, error)
	Migrate(migration) error
}
