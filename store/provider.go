package store

import (
	"github.com/graphql-go/graphql"
	"github.com/valocode/bubbly/api/core"
	"github.com/valocode/bubbly/env"
)

// Provider provides an interface for persisting readiness data.
type provider interface {
	Tenants() ([]string, error)
	CreateTenant(string) error
	Apply(string, *bubblySchema) error
	Migrate(string, *bubblySchema, schemaUpdates) error
	Save(*env.BubblyContext, string, *schemaGraph, dataTree) error
	ResolveQuery(string, *schemaGraph, graphql.ResolveParams) (interface{}, error)
	HasTable(string, core.Table) (bool, error)
}
