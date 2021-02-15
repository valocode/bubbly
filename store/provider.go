package store

import (
	"github.com/graphql-go/graphql"
	"github.com/verifa/bubbly/api/core"
)

// Provider provides an interface for persisting readiness data.
type provider interface {
	Apply(*bubblySchema) error
	Save(*bubblySchema, dataTree) error
	ResolveQuery(graph *schemaGraph, params graphql.ResolveParams) (interface{}, error)
	HasTable(core.Table) (bool, error)
}
