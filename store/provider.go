package store

import (
	"github.com/graphql-go/graphql"
	"github.com/valocode/bubbly/api/core"
	"github.com/valocode/bubbly/env"
)

// Provider provides an interface for persisting readiness data.
type provider interface {
	Apply(*bubblySchema) error
	Save(*env.BubblyContext, *schemaGraph, dataTree) error
	ResolveQuery(graph *schemaGraph, params graphql.ResolveParams) (interface{}, error)
	HasTable(core.Table) (bool, error)
	GenerateMigration(ctx *env.BubblyContext, cl Changelog) (migration, error)
	Migrate(migration) error
}
