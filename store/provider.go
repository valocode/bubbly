package store

import (
	"github.com/graphql-go/graphql"
	"github.com/valocode/bubbly/env"
)

// Provider provides an interface for persisting readiness data.
type provider interface {
	Tenants() ([]string, error)
	CreateTenant(string) error
	Close()
	Apply(string, *bubblySchema) error
	Migrate(string, *bubblySchema, schemaUpdates) error
	Save(*env.BubblyContext, string, *SchemaGraph, dataTree) error
	ResolveQuery(string, *SchemaGraph, graphql.ResolveParams) (interface{}, error)
	HasTable(string, string) (bool, error)
}
